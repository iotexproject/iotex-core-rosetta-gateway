// Copyright (c) 2020 IoTeX Foundation
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package iotex_client

import (
	"context"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"sync"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/gogo/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"

	"github.com/iotexproject/go-pkgs/crypto"
	"github.com/iotexproject/go-pkgs/hash"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"

	"github.com/iotexproject/iotex-core-rosetta-gateway/config"
)

const (
	Transfer               = "transfer"
	Execution              = "execution"
	DepositToRewardingFund = "depositToRewardingFund"
	ClaimFromRewardingFund = "claimFromRewardingFund"
	StakeCreate            = "stakeCreate"
	StakeWithdraw          = "stakeWithdraw"
	StakeAddDeposit        = "stakeAddDeposit"
	CandidateRegister      = "candidateRegister"
	StatusSuccess          = "success"
	StatusFail             = "fail"
	ActionTypeFee          = "fee"
)

type (
	// IoTexClient is the IoTex blockchain client interface.
	IoTexClient interface {
		// GetChainID returns the network chain context, derived from the
		// genesis document.
		GetChainID(ctx context.Context) (string, error)

		// GetBlock returns the IoTex block at given height.
		GetBlock(ctx context.Context, height int64) (*IoTexBlock, error)

		// GetLatestBlock returns latest IoTex block.
		GetLatestBlock(ctx context.Context) (*IoTexBlock, error)

		// GetGenesisBlock returns the IoTex genesis block.
		GetGenesisBlock(ctx context.Context) (*IoTexBlock, error)

		// GetAccount returns the IoTex staking account for given owner address
		// at given height.
		GetAccount(ctx context.Context, height int64, owner string) (*Account, error)

		// SubmitTx submits the given encoded transaction to the node.
		SubmitTx(ctx context.Context, tx *iotextypes.Action) (txid string, err error)

		// GetStatus returns the status overview of the node.
		GetStatus(ctx context.Context) (*iotexapi.GetChainMetaResponse, error)

		// GetVersion returns the server's version.
		GetVersion(ctx context.Context) (*iotexapi.GetServerMetaResponse, error)

		// GetTransactions returns transactions of the block.
		GetTransactions(ctx context.Context, height int64) ([]*types.Transaction, error)

		// GetConfig returns the config.
		GetConfig() *config.Config
	}

	// IoTexBlock is the IoTex blockchain's block.
	IoTexBlock struct {
		Height       int64  // Block height.
		Hash         string // Block hash.
		Timestamp    int64  // UNIX time, converted to milliseconds.
		ParentHeight int64  // Height of parent block.
		ParentHash   string // Hash of parent block.
	}

	Account struct {
		Nonce   uint64
		Balance string
	}

	// grpcIoTexClient is an implementation of IoTexClient using gRPC.
	grpcIoTexClient struct {
		sync.RWMutex

		endpoint string
		grpcConn *grpc.ClientConn
		cfg      *config.Config
	}
)

// NewIoTexClient returns an implementation of IoTexClient
func NewIoTexClient(cfg *config.Config) (cli IoTexClient, err error) {
	opts := []grpc.DialOption{}
	if cfg.Server.SecureEndpoint {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	grpc, err := grpc.Dial(cfg.Server.Endpoint, opts...)
	if err != nil {
		return
	}
	cli = &grpcIoTexClient{grpcConn: grpc, cfg: cfg}
	return
}

func (c *grpcIoTexClient) GetChainID(ctx context.Context) (string, error) {
	return c.cfg.NetworkIdentifier.Network, nil
}

func (c *grpcIoTexClient) GetBlock(ctx context.Context, height int64) (ret *IoTexBlock, err error) {
	return c.getBlock(ctx, height)
}

func (c *grpcIoTexClient) GetLatestBlock(ctx context.Context) (*IoTexBlock, error) {
	err := c.reconnect()
	if err != nil {
		return nil, err
	}
	client := iotexapi.NewAPIServiceClient(c.grpcConn)
	res, err := client.GetChainMeta(context.Background(), &iotexapi.GetChainMetaRequest{})
	if err != nil {
		return nil, err
	}
	return c.getBlock(ctx, int64(res.ChainMeta.Height))
}

func (c *grpcIoTexClient) GetGenesisBlock(ctx context.Context) (*IoTexBlock, error) {
	return c.getBlock(ctx, 1)
}

func (c *grpcIoTexClient) GetAccount(ctx context.Context, height int64, owner string) (ret *Account, err error) {
	err = c.reconnect()
	if err != nil {
		return
	}
	client := iotexapi.NewAPIServiceClient(c.grpcConn)
	request := &iotexapi.GetAccountRequest{Address: owner}
	resp, err := client.GetAccount(ctx, request)
	if err != nil {
		return nil, err
	}
	ret = &Account{
		Nonce:   resp.AccountMeta.Nonce,
		Balance: resp.AccountMeta.Balance,
	}
	return
}

func (c *grpcIoTexClient) GetTransactions(ctx context.Context, height int64) (ret []*types.Transaction, err error) {
	client := iotexapi.NewAPIServiceClient(c.grpcConn)
	getRawBlocksRes, err := client.GetRawBlocks(context.Background(), &iotexapi.GetRawBlocksRequest{
		StartHeight:  uint64(height),
		Count:        1,
		WithReceipts: true,
	})
	if err != nil || len(getRawBlocksRes.GetBlocks()) != 1 {
		return
	}
	ret = make([]*types.Transaction, 0)
	actionMap := make(map[hash.Hash256]*iotextypes.Action)
	receiptMap := make(map[hash.Hash256]*iotextypes.Receipt)
	// hashSlice for fixed sequence,b/c map is unordered
	hashSlice := make([]hash.Hash256, 0)
	blk := getRawBlocksRes.GetBlocks()[0]
	for _, act := range blk.GetBlock().GetBody().GetActions() {
		proto, err := proto.Marshal(act)
		if err != nil {
			return nil, err
		}
		actionMap[hash.Hash256b(proto)] = act
		hashSlice = append(hashSlice, hash.Hash256b(proto))
	}
	for _, receipt := range blk.GetReceipts() {
		receiptMap[hash.BytesToHash256(receipt.ActHash)] = receipt
	}
	for _, h := range hashSlice {
		act := actionMap[h]
		r, ok := receiptMap[h]
		if !ok {
			err = errors.New(fmt.Sprintf("failed find receipt:%s", hex.EncodeToString(h[:])))
			return
		}
		decode, err := c.decodeAction(ctx, act, h, r, client)
		if err != nil {
			// change to continue or return when systemlog is enabled in testnet
			// TODO change it back
			//return nil, err
			continue
		}
		if decode != nil {
			ret = append(ret, decode)
		}
	}
	return
}

func (c *grpcIoTexClient) SubmitTx(ctx context.Context, tx *iotextypes.Action) (txid string, err error) {
	err = c.reconnect()
	if err != nil {
		return
	}
	client := iotexapi.NewAPIServiceClient(c.grpcConn)
	ret, err := client.SendAction(ctx, &iotexapi.SendActionRequest{Action: tx})
	if err != nil {
		return
	}
	txid = ret.ActionHash
	return
}

func (c *grpcIoTexClient) GetStatus(ctx context.Context) (*iotexapi.GetChainMetaResponse, error) {
	err := c.reconnect()
	if err != nil {
		return nil, err
	}
	client := iotexapi.NewAPIServiceClient(c.grpcConn)
	return client.GetChainMeta(context.Background(), &iotexapi.GetChainMetaRequest{})
}

func (c *grpcIoTexClient) GetVersion(ctx context.Context) (*iotexapi.GetServerMetaResponse, error) {
	err := c.reconnect()
	if err != nil {
		return nil, err
	}
	client := iotexapi.NewAPIServiceClient(c.grpcConn)
	return client.GetServerMeta(ctx, &iotexapi.GetServerMetaRequest{})
}

func (c *grpcIoTexClient) GetConfig() *config.Config {
	return c.cfg
}

func (c *grpcIoTexClient) getBlock(ctx context.Context, height int64) (ret *IoTexBlock, err error) {
	err = c.reconnect()
	if err != nil {
		return
	}
	var parentHeight uint64
	if height <= 1 {
		parentHeight = 1
	} else {
		parentHeight = uint64(height) - 1
	}
	client := iotexapi.NewAPIServiceClient(c.grpcConn)
	count := uint64(2)
	if parentHeight == uint64(height) {
		count = 1
	}
	request := &iotexapi.GetBlockMetasRequest{
		Lookup: &iotexapi.GetBlockMetasRequest_ByIndex{
			ByIndex: &iotexapi.GetBlockMetasByIndexRequest{
				Start: parentHeight,
				Count: count,
			},
		},
	}
	resp, err := client.GetBlockMetas(ctx, request)
	if err != nil {
		return nil, err
	}
	if len(resp.BlkMetas) == 0 {
		return nil, errors.New("not found")
	}
	var blk, parentBlk *iotextypes.BlockMeta
	if len(resp.BlkMetas) == 2 {
		blk = resp.BlkMetas[1]
		parentBlk = resp.BlkMetas[0]
	} else {
		blk = resp.BlkMetas[0]
		parentBlk = resp.BlkMetas[0]
	}
	ret = &IoTexBlock{
		Height:       int64(blk.Height),
		Hash:         blk.Hash,
		Timestamp:    blk.Timestamp.Seconds * 1e3, // ms
		ParentHeight: int64(parentHeight),
		ParentHash:   parentBlk.Hash,
	}
	return
}

func (c *grpcIoTexClient) reconnect() (err error) {
	c.Lock()
	defer c.Unlock()
	// Check if the existing connection is good.
	if c.grpcConn != nil && c.grpcConn.GetState() != connectivity.Shutdown {
		return
	}
	c.grpcConn, err = grpc.Dial(c.endpoint, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	return err
}

func (c *grpcIoTexClient) decodeAction(ctx context.Context, act *iotextypes.Action, h hash.Hash256, receipt *iotextypes.Receipt, client iotexapi.APIServiceClient) (ret *types.Transaction, err error) {
	srcPub, err := crypto.BytesToPublicKey(act.GetSenderPubKey())
	if err != nil {
		return
	}
	callerAddr, err := address.FromBytes(srcPub.Hash())
	if err != nil {
		return
	}
	ret, status, err := c.gasFeeAndStatus(callerAddr, act, h, receipt)
	if err != nil {
		return
	}

	if act.GetCore().GetExecution() != nil {
		// TODO test when testnet enable systemlog
		err = c.handleExecution(ctx, ret, status, hex.EncodeToString(h[:]), client)
		return
	}

	amount, senderSign, actionType, dst, err := assertAction(act)
	if err != nil {
		return
	}
	if amount == "" || actionType == "" {
		return
	}
	senderAmountWithSign := amount
	dstAmountWithSign := amount
	if amount != "0" {
		if senderSign == "-" {
			senderAmountWithSign = senderSign + amount
		} else {
			dstAmountWithSign = "-" + amount
		}
	}

	src := []*addressAmount{{address: callerAddr.String(), amount: senderAmountWithSign}}
	var dstAll []*addressAmount
	if dst != "" {
		dstAll = []*addressAmount{{address: dst, amount: dstAmountWithSign}}
	}
	err = c.packTransaction(ret, src, dstAll, actionType, status)
	return
}

func (c *grpcIoTexClient) handleExecution(ctx context.Context, ret *types.Transaction, status, hash string, client iotexapi.APIServiceClient) (err error) {
	request := &iotexapi.GetEvmTransfersByActionHashRequest{
		ActionHash: hash,
	}
	resp, err := client.GetEvmTransfersByActionHash(ctx, request)
	if err != nil {
		return
	}
	var src, dst addressAmountList
	for _, transfer := range resp.GetActionEvmTransfers().GetEvmTransfers() {
		amount := new(big.Int).SetBytes(transfer.Amount)
		amountStr := amount.String()
		if amount.Sign() != 0 {
			amountStr = "-" + amount.String()
		}
		src = append(src, &addressAmount{
			address: transfer.From,
			amount:  amountStr,
		})
		dst = append(dst, &addressAmount{
			address: transfer.To,
			amount:  new(big.Int).SetBytes(transfer.Amount).String(),
		})
	}
	return c.packTransaction(ret, src, dst, Execution, status)
}

func (c *grpcIoTexClient) gasFeeAndStatus(callerAddr address.Address, act *iotextypes.Action, h hash.Hash256, receipt *iotextypes.Receipt) (ret *types.Transaction, status string, err error) {
	status = StatusSuccess
	if receipt.GetStatus() != 1 {
		status = StatusFail
	}
	gasConsumed := new(big.Int).SetUint64(receipt.GetGasConsumed())
	gasPrice, ok := new(big.Int).SetString(act.GetCore().GetGasPrice(), 10)
	if !ok {
		err = errors.New("convert gas price error")
		return
	}
	gasFee := gasPrice.Mul(gasPrice, gasConsumed)
	amount := gasFee.String()
	// if gasFee is not 0
	if gasFee.Sign() == 1 {
		amount = "-" + amount
	}

	sender := addressAmountList{{address: callerAddr.String(), amount: amount}}
	var oper []*types.Operation
	_, oper, err = c.addOperation(sender, ActionTypeFee, status, 0, oper)
	if err != nil {
		return
	}
	ret = &types.Transaction{
		TransactionIdentifier: &types.TransactionIdentifier{
			hex.EncodeToString(h[:]),
		},
		Operations: oper,
		Metadata:   nil,
	}
	return
}

func (c *grpcIoTexClient) packTransaction(ret *types.Transaction, src, dst addressAmountList, actionType, status string) (err error) {
	sort.Sort(src)
	sort.Sort(dst)
	var oper []*types.Operation
	endIndex, oper, err := c.addOperation(src, actionType, status, 1, oper)
	if err != nil {
		return
	}
	_, oper, err = c.addOperation(dst, actionType, status, endIndex, oper)
	if err != nil {
		return
	}
	ret.Operations = append(ret.Operations, oper...)
	return
}

func (c *grpcIoTexClient) addOperation(l addressAmountList, actionType, status string, startIndex int64, oper []*types.Operation) (int64, []*types.Operation, error) {
	for _, s := range l {
		oper = append(oper, &types.Operation{
			OperationIdentifier: &types.OperationIdentifier{
				Index:        startIndex,
				NetworkIndex: nil,
			},
			RelatedOperations: nil,
			Type:              actionType,
			Status:            status,
			Account: &types.AccountIdentifier{
				Address:    s.address,
				SubAccount: nil,
				Metadata:   nil,
			},
			Amount: &types.Amount{
				Value: s.amount,
				Currency: &types.Currency{
					Symbol:   c.cfg.Currency.Symbol,
					Decimals: c.cfg.Currency.Decimals,
					Metadata: nil,
				},
				Metadata: nil,
			},
			Metadata: nil,
		})
		startIndex++
	}
	return startIndex, oper, nil
}

func assertAction(act *iotextypes.Action) (amount, senderSign, actionType, dst string, err error) {
	amount = "0"
	senderSign = "-"
	switch {
	case act.GetCore().GetTransfer() != nil:
		actionType = Transfer
		amount = act.GetCore().GetTransfer().GetAmount()
		dst = act.GetCore().GetTransfer().GetRecipient()
	case act.GetCore().GetDepositToRewardingFund() != nil:
		actionType = DepositToRewardingFund
		amount = act.GetCore().GetDepositToRewardingFund().GetAmount()
	case act.GetCore().GetClaimFromRewardingFund() != nil:
		actionType = ClaimFromRewardingFund
		amount = act.GetCore().GetClaimFromRewardingFund().GetAmount()
		senderSign = "+"
	case act.GetCore().GetStakeAddDeposit() != nil:
		actionType = StakeAddDeposit
		amount = act.GetCore().GetStakeAddDeposit().GetAmount()
	case act.GetCore().GetStakeCreate() != nil:
		actionType = StakeCreate
		amount = act.GetCore().GetStakeCreate().GetStakedAmount()
	case act.GetCore().GetStakeWithdraw() != nil:
		// TODO need to add amount when it's available on iotex-core
		actionType = StakeWithdraw
	case act.GetCore().GetCandidateRegister() != nil:
		actionType = CandidateRegister
		amount = act.GetCore().GetCandidateRegister().GetStakedAmount()
	}
	return
}
