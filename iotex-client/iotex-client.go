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
	"fmt"
	"math/big"
	"sync"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"

	"github.com/iotexproject/go-pkgs/hash"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"

	"github.com/iotexproject/iotex-core-rosetta-gateway/config"
)

const (
	rewardingProtocolID      = "rewarding"
	stakingProtocolID        = "staking"
	availableBalanceMethodID = "AvailableBalance"
	Transfer                 = "transfer"
	Execution                = "execution"
	DepositToRewardingFund   = "depositToRewardingFund"
	ClaimFromRewardingFund   = "claimFromRewardingFund"
	StakeCreate              = "stakeCreate"
	StakeWithdraw            = "stakeWithdraw"
	StakeAddDeposit          = "stakeAddDeposit"
	CandidateRegister        = "candidateRegister"
	StatusSuccess            = "success"
	StatusFail               = "fail"
	ActionTypeFee            = "fee"
	// NonceKey is the name of the key in the Metadata map inside a
	// ConstructionMetadataResponse that specifies the next valid nonce.
	NonceKey = "nonce"
)

var (
	RewardingAddress string
	StakingAddress   string
)

func init() {
	h := hash.Hash160b([]byte(rewardingProtocolID))
	addr, _ := address.FromBytes(h[:])
	RewardingAddress = addr.String()
	h = hash.Hash160b([]byte(stakingProtocolID))
	addr, _ = address.FromBytes(h[:])
	StakingAddress = addr.String()
}

type (
	// IoTexClient is the IoTex blockchain client interface.
	IoTexClient interface {
		// GetChainID returns the network chain context, derived from the
		// genesis document.
		GetChainID(ctx context.Context) (string, error)

		// GetBlock returns the IoTex block at given height.
		GetBlock(ctx context.Context, height int64) (*types.Block, error)

		// GetLatestBlock returns latest IoTex block.
		GetLatestBlock(ctx context.Context) (*types.Block, error)

		// GetGenesisBlock returns the IoTex genesis block.
		GetGenesisBlock(ctx context.Context) (*types.Block, error)

		// GetAccount returns the IoTex staking account for given owner address
		// at given height.
		GetAccount(ctx context.Context, height int64, owner string) (*types.AccountBalanceResponse, error)

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

	// sort is not necessary,add operation according the sequence from core
	addressAmount struct {
		address    string
		amount     string
		actionType string
	}
	addressAmountList []*addressAmount

	operation struct {
		src        string
		dst        string
		amount     string
		actionType string
		isPositive bool
	}
	operationList []*operation

	// grpcIoTexClient is an implementation of IoTexClient using gRPC.
	grpcIoTexClient struct {
		sync.RWMutex

		grpcConn *grpc.ClientConn
		client   iotexapi.APIServiceClient
		cfg      *config.Config
	}
)

// NewIoTexClient returns an implementation of IoTexClient
func NewIoTexClient(cfg *config.Config) (cli IoTexClient, err error) {
	gcli := &grpcIoTexClient{cfg: cfg}
	if err = gcli.connect(); err != nil {
		return
	}
	cli = gcli
	return
}

func (c *grpcIoTexClient) GetChainID(ctx context.Context) (string, error) {
	return c.cfg.NetworkIdentifier.Network, nil
}

func (c *grpcIoTexClient) GetBlock(ctx context.Context, height int64) (ret *types.Block, err error) {
	if err = c.connect(); err != nil {
		return
	}
	return c.getBlock(ctx, height)
}

func (c *grpcIoTexClient) GetLatestBlock(ctx context.Context) (*types.Block, error) {
	if err := c.connect(); err != nil {
		return nil, err
	}
	return c.getLatestBlock(ctx)
}

func (c *grpcIoTexClient) getLatestBlock(ctx context.Context) (*types.Block, error) {
	res, err := c.client.GetChainMeta(context.Background(), &iotexapi.GetChainMetaRequest{})
	if err != nil {
		return nil, err
	}
	return c.getBlock(ctx, int64(res.ChainMeta.Height))
}

func (c *grpcIoTexClient) GetGenesisBlock(ctx context.Context) (*types.Block, error) {
	if err := c.connect(); err != nil {
		return nil, err
	}
	return c.getBlock(ctx, 1)
}

func (c *grpcIoTexClient) GetAccount(ctx context.Context, height int64, owner string) (ret *types.AccountBalanceResponse, err error) {
	if err = c.connect(); err != nil {
		return
	}

	if owner == RewardingAddress {
		return c.getRewardingAccount(ctx, height)
	}

	request := &iotexapi.GetAccountRequest{Address: owner}
	resp, err := c.client.GetAccount(ctx, request)
	if err != nil {
		return
	}
	acc := resp.GetAccountMeta()
	blk := resp.GetBlockIdentifier()
	ret = &types.AccountBalanceResponse{
		BlockIdentifier: &types.BlockIdentifier{
			Index: int64(blk.GetHeight()),
			Hash:  blk.GetHash(),
		},
		Balances: []*types.Amount{{
			Value: acc.GetBalance(),
			Currency: &types.Currency{
				Symbol:   c.cfg.Currency.Symbol,
				Decimals: c.cfg.Currency.Decimals,
				Metadata: nil,
			}}},
		Metadata: &map[string]interface{}{NonceKey: acc.GetPendingNonce()},
	}
	return
}

func (c *grpcIoTexClient) getRewardingAccount(ctx context.Context, height int64) (ret *types.AccountBalanceResponse, err error) {
	// call readState
	out, err := c.client.ReadState(context.Background(), &iotexapi.ReadStateRequest{
		ProtocolID: []byte(rewardingProtocolID),
		MethodName: []byte(availableBalanceMethodID),
		Arguments:  nil,
	})
	if err != nil {
		return
	}
	val, ok := big.NewInt(0).SetString(string(out.Data), 10)
	if !ok {
		err = errors.New("balance convert error")
		return
	}
	blk, err := c.getLatestBlock(ctx)
	if err != nil {
		return
	}

	ret = &types.AccountBalanceResponse{
		BlockIdentifier: blk.BlockIdentifier,
		Balances: []*types.Amount{{
			Value: val.String(),
			Currency: &types.Currency{
				Symbol:   c.cfg.Currency.Symbol,
				Decimals: c.cfg.Currency.Decimals,
				Metadata: nil,
			},
		},
		},
		Metadata: &map[string]interface{}{NonceKey: 0},
	}
	return
}

func (c *grpcIoTexClient) GetTransactions(ctx context.Context, height int64) (ret []*types.Transaction, err error) {
	ret = make([]*types.Transaction, 0)
	if err = c.connect(); err != nil {
		return
	}
	actionMap, receiptMap, hashSlice, err := c.getRawBlock(ctx, height)
	if err != nil {
		return
	}
	// get ImplicitTransferLog by height,if log is not exist,the err will be nil
	transferLogMap, err := getImplicitTransferLog(ctx, height, c.client)
	if err != nil {
		return
	}
	for _, h := range hashSlice {
		// grantReward action,gas fee and amount both 0
		if actionMap[h].GetCore().GetGrantReward() != nil {
			continue
		}
		r, ok := receiptMap[h]
		if !ok {
			err = errors.New(fmt.Sprintf("failed find receipt:%s", h))
			return
		}
		decode, err := c.decodeAction(ctx, h, actionMap[h], r, transferLogMap[h], height)
		if err != nil {
			return nil, err
		}
		if decode != nil {
			ret = append(ret, decode)
		}
	}
	fillIndex(ret)
	return
}

func (c *grpcIoTexClient) SubmitTx(ctx context.Context, tx *iotextypes.Action) (txid string, err error) {
	if err = c.connect(); err != nil {
		return
	}
	ret, err := c.client.SendAction(ctx, &iotexapi.SendActionRequest{Action: tx})
	if err != nil {
		return
	}
	txid = ret.ActionHash
	return
}

func (c *grpcIoTexClient) GetStatus(ctx context.Context) (*iotexapi.GetChainMetaResponse, error) {
	err := c.connect()
	if err != nil {
		return nil, err
	}
	return c.client.GetChainMeta(context.Background(), &iotexapi.GetChainMetaRequest{})
}

func (c *grpcIoTexClient) GetVersion(ctx context.Context) (*iotexapi.GetServerMetaResponse, error) {
	err := c.connect()
	if err != nil {
		return nil, err
	}
	return c.client.GetServerMeta(ctx, &iotexapi.GetServerMetaRequest{})
}

func (c *grpcIoTexClient) GetConfig() *config.Config {
	return c.cfg
}

func (c *grpcIoTexClient) connect() (err error) {
	c.Lock()
	defer c.Unlock()
	// Check if the existing connection is good.
	if c.grpcConn != nil && c.grpcConn.GetState() != connectivity.Shutdown {
		return
	}
	opts := []grpc.DialOption{}
	if c.cfg.Server.SecureEndpoint {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	c.grpcConn, err = grpc.Dial(c.cfg.Server.Endpoint, opts...)
	c.client = iotexapi.NewAPIServiceClient(c.grpcConn)
	return err
}

func (c *grpcIoTexClient) getBlock(ctx context.Context, height int64) (ret *types.Block, err error) {
	var parentHeight uint64
	if height <= 1 {
		parentHeight = 1
	} else {
		parentHeight = uint64(height) - 1
	}
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
	resp, err := c.client.GetBlockMetas(ctx, request)
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
	ret = &types.Block{
		BlockIdentifier: &types.BlockIdentifier{
			Index: int64(blk.Height),
			Hash:  blk.Hash,
		},
		ParentBlockIdentifier: &types.BlockIdentifier{
			Index: int64(parentHeight),
			Hash:  parentBlk.Hash,
		},
		Timestamp: blk.Timestamp.Seconds * 1e3, // ms,
	}
	return
}

func (c *grpcIoTexClient) getRawBlock(ctx context.Context, height int64) (actionMap map[string]*iotextypes.Action, receiptMap map[string]*iotextypes.Receipt, hashSlice []string, err error) {
	getRawBlocksRes, err := c.client.GetRawBlocks(ctx, &iotexapi.GetRawBlocksRequest{
		StartHeight:  uint64(height),
		Count:        1,
		WithReceipts: true,
	})
	if err != nil || len(getRawBlocksRes.GetBlocks()) != 1 {
		return
	}

	actionMap = make(map[string]*iotextypes.Action)
	receiptMap = make(map[string]*iotextypes.Receipt)
	// hashSlice for fixed sequence,b/c map is unordered
	hashSlice = make([]string, 0)
	blk := getRawBlocksRes.GetBlocks()[0]
	for _, act := range blk.GetBlock().GetBody().GetActions() {
		var pro []byte
		pro, err = proto.Marshal(act)
		if err != nil {
			return
		}
		hashArray := hash.Hash256b(pro)
		h := hex.EncodeToString(hashArray[:])
		actionMap[h] = act
		hashSlice = append(hashSlice, h)
	}
	for _, receipt := range blk.GetReceipts() {
		receiptMap[hex.EncodeToString(receipt.ActHash)] = receipt
	}
	return
}

func (c *grpcIoTexClient) decodeAction(ctx context.Context, h string, act *iotextypes.Action,
	receipt *iotextypes.Receipt, transferLogs []*iotextypes.ImplicitTransferLog_Transaction,
	height int64) (ret *types.Transaction, err error) {
	ret = &types.Transaction{TransactionIdentifier: &types.TransactionIdentifier{h}}
	callerAddr, err := getCaller(act)
	if err != nil {
		return
	}
	oper, err := c.getGasFee(act, receipt, height)
	if err != nil {
		return
	}
	oper.src = callerAddr.String()
	operations := make(operationList, 0)
	operations = append(operations, oper)

	if receipt.Status == uint64(iotextypes.ReceiptStatus_Success) {
		operations, err = c.prepareOperations(ctx, act, h, transferLogs, operations)
		if err != nil {
			return
		}
	}
	for _, oper := range operations {
		err = c.handleOperations(ret, oper, callerAddr.String())
		if err != nil {
			return
		}
	}
	return
}

func (c *grpcIoTexClient) prepareOperations(ctx context.Context, act *iotextypes.Action, h string, transferLogs []*iotextypes.ImplicitTransferLog_Transaction, operations operationList) (operationList, error) {
	if transferLogs != nil {
		// handle implicit transfer log first
		for _, t := range transferLogs {
			operations = append(operations, &operation{
				src:        t.GetSender(),
				dst:        t.GetRecipient(),
				amount:     t.GetAmount(),
				actionType: getActionType(t.GetTopic()),
				isPositive: false,
			})
		}
	}

	// get execution action operations,this is still needed for those in case of there're no implicit log
	if act.GetCore().GetExecution() != nil && transferLogs == nil {
		oper, err := c.getExecution(ctx, act, h)
		if err != nil {
			return operations, err
		}
		operations = append(operations, oper)
	}
	// get general action operations
	operations = assertAction(act, operations)
	return operations, nil
}

func (c *grpcIoTexClient) handleOperations(ret *types.Transaction, oper *operation, caller string) error {
	senderAmountWithSign := oper.amount
	dstAmountWithSign := oper.amount
	if oper.amount != "0" {
		if !oper.isPositive {
			senderAmountWithSign = "-" + senderAmountWithSign
		} else {
			dstAmountWithSign = "-" + dstAmountWithSign
		}
	}
	if oper.src == "" {
		oper.src = caller
	}
	aal := addressAmountList{{oper.src, senderAmountWithSign, oper.actionType}}
	if oper.dst != "" {
		aal = append(aal, &addressAmount{oper.dst, dstAmountWithSign, oper.actionType})
	}
	return c.addOperation(ret, aal, StatusSuccess)
}

func (c *grpcIoTexClient) getExecution(ctx context.Context, act *iotextypes.Action,
	h string) (ret *operation, err error) {
	ret = &operation{}
	ret.dst = act.GetCore().GetExecution().GetContract()
	if ret.dst == "" {
		ret.dst, err = getContractAddress(ctx, h, c.client)
		if err != nil {
			return
		}
	}
	ret.amount = act.GetCore().GetExecution().GetAmount()
	ret.actionType = Execution
	return
}

func (c *grpcIoTexClient) getGasFee(act *iotextypes.Action, receipt *iotextypes.Receipt,
	height int64) (oper *operation, err error) {
	oper = &operation{}
	gasConsumed := new(big.Int).SetUint64(receipt.GetGasConsumed())
	gasPrice, ok := new(big.Int).SetString(act.GetCore().GetGasPrice(), 10)
	if !ok {
		err = errors.New("convert gas price error")
		return
	}
	gasFee := gasPrice.Mul(gasPrice, gasConsumed)
	amount := gasFee.String()
	if height < c.cfg.PacificBlockHeight && act.GetCore().GetExecution() != nil {
		// before PacificBlockHeight Execution pay double gas fee
		// and only one share go to reward address
		amount = gasFee.Mul(gasFee, big.NewInt(2)).String()
	}
	oper.amount = amount
	oper.actionType = ActionTypeFee
	oper.dst = RewardingAddress
	return
}

func (c *grpcIoTexClient) addOperation(ret *types.Transaction, amountList addressAmountList, status string) error {
	var oper []*types.Operation
	for _, s := range amountList {
		oper = append(oper, &types.Operation{
			OperationIdentifier: &types.OperationIdentifier{
				NetworkIndex: nil,
			},
			RelatedOperations: nil,
			Type:              s.actionType,
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
	}
	ret.Operations = append(ret.Operations, oper...)
	return nil
}
