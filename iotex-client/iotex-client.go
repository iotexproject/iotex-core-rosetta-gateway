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
	"log"
	"math/big"
	"strconv"
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

		SuggestGasPrice(ctx context.Context) (uint64, error)

		EstimateGasForAction(ctx context.Context, action *iotextypes.Action) (uint64, error)

		GetBlockTransaction(ctx context.Context, actionHash string) (*types.Transaction, error)

		GetMemPool(ctx context.Context, actionHashes []string) ([]*types.TransactionIdentifier, error)

		GetMemPoolTransaction(ctx context.Context, h string) (*types.Transaction, error)
	}
)

type (
	// grpcIoTexClient is an implementation of IoTexClient using gRPC.
	grpcIoTexClient struct {
		sync.RWMutex

		grpcConn *grpc.ClientConn
		client   iotexapi.APIServiceClient
		cfg      *config.Config
	}
)

type (
	addressAmount struct {
		senderAddr   string
		dstAddr      string
		senderAmount string
		dstAmount    string
		actionType   string
	}
	addressAmountList []*addressAmount
)

// NewIoTexClient returns an implementation of IoTexClient
func NewIoTexClient(cfg *config.Config) (cli IoTexClient, err error) {
	return &grpcIoTexClient{cfg: cfg}, nil
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

func (c *grpcIoTexClient) SuggestGasPrice(ctx context.Context) (uint64, error) {
	if err := c.connect(); err != nil {
		return 0, err
	}
	response, err := c.client.SuggestGasPrice(ctx, &iotexapi.SuggestGasPriceRequest{})
	return response.GetGasPrice(), err
}

func (c *grpcIoTexClient) EstimateGasForAction(ctx context.Context, action *iotextypes.Action) (uint64, error) {
	if err := c.connect(); err != nil {
		return 0, err
	}
	response, err := c.client.EstimateGasForAction(ctx, &iotexapi.EstimateGasForActionRequest{Action: action})
	return response.GetGas(), err
}

func (c *grpcIoTexClient) GetAccount(ctx context.Context, height int64, owner string) (ret *types.AccountBalanceResponse, err error) {
	if err = c.connect(); err != nil {
		return
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
		Metadata: map[string]interface{}{NonceKey: acc.GetPendingNonce()},
	}
	return
}

func (c *grpcIoTexClient) GetTransactions(ctx context.Context, height int64) (ret []*types.Transaction, err error) {
	ret = make([]*types.Transaction, 0)
	if err = c.connect(); err != nil {
		return
	}
	actionMap, _, hashSlice, err := c.getRawBlock(ctx, height)
	if err != nil {
		return
	}
	// get TransactionLog by height,if log is not exist,the err will be nil
	transferLogMap, err := getTransactionLog(ctx, height, c.client)
	if err != nil {
		return
	}
	for _, h := range hashSlice {
		if transferLogMap[h] != nil {
			transaction := c.packTransaction(h, transferLogMap[h])
			ret = append(ret, transaction)
		} else if c.cfg.KeepNoneTxAction {
			ret = append(ret, c.genNoneTxActTransaction(h, actionMap[h]))
		} else if rec, err := c.client.GetReceiptByAction(ctx, &iotexapi.GetReceiptByActionRequest{ActionHash: h}); err == nil && rec.ReceiptInfo.Receipt.GetStatus() != 1 {
			transaction, err := c.packActionToTransaction(actionMap[h], h, StatusFail)
			if err != nil {
				continue
			}
			ret = append(ret, transaction)
		}
	}
	ret = fillIndex(ret)
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

func genBlock(parentBlk, blk *iotextypes.BlockMeta) *types.Block {
	return &types.Block{
		BlockIdentifier: &types.BlockIdentifier{
			Index: int64(blk.Height),
			Hash:  blk.Hash,
		},
		ParentBlockIdentifier: &types.BlockIdentifier{
			Index: int64(parentBlk.Height),
			Hash:  parentBlk.Hash,
		},
		Timestamp: blk.Timestamp.Seconds * 1e3, // ms,
	}
}

func (c *grpcIoTexClient) getBlock(ctx context.Context, height int64) (ret *types.Block, err error) {
	parentHeight := uint64(height) - 1
	count := uint64(2)
	if parentHeight <= 0 {
		count = 1
		parentHeight = 1
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
	parentBlk := resp.BlkMetas[0]
	blk := resp.BlkMetas[0]
	if len(resp.BlkMetas) == 2 {
		blk = resp.BlkMetas[1]
	}
	ret = genBlock(parentBlk, blk)
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

func (c *grpcIoTexClient) genNoneTxActTransaction(h string, act *iotextypes.Action) *types.Transaction {
	callerAddr, err := getCaller(act)
	if err != nil {
		log.Fatalln("failed to get action caller", err)
	}
	// gen an empty gas
	tx := &iotextypes.TransactionLog_Transaction{
		Type:      iotextypes.TransactionLogType_GAS_FEE,
		Sender:    callerAddr.String(),
		Recipient: address.RewardingPoolAddr,
		Amount:    "0",
	}
	return c.packTransaction(h, []*iotextypes.TransactionLog_Transaction{tx})
}

func (c *grpcIoTexClient) packTransaction(h string, transferLogs []*iotextypes.TransactionLog_Transaction) *types.Transaction {
	ret := &types.Transaction{TransactionIdentifier: &types.TransactionIdentifier{h}}
	ret.Operations = make([]*types.Operation, 0, len(transferLogs))
	for _, t := range transferLogs {
		ops := c.covertToOperations(t)
		ret.Operations = append(ret.Operations, ops...)
	}
	return ret
}

func (c *grpcIoTexClient) covertToOperations(s *iotextypes.TransactionLog_Transaction) []*types.Operation {
	ops := make([]*types.Operation, 0, 2)
	// sender
	if s.GetSender() != "" {
		senderAmount := "-" + s.GetAmount()
		if s.GetAmount() == "0" {
			senderAmount = s.GetAmount()
		}
		ops = append(ops, &types.Operation{
			OperationIdentifier: &types.OperationIdentifier{
				NetworkIndex: nil,
			},
			RelatedOperations: nil,
			Type:              s.GetType().String(),
			Status:            StatusSuccess,
			Account: &types.AccountIdentifier{
				Address:    s.GetSender(),
				SubAccount: nil,
				Metadata:   nil,
			},
			Amount: &types.Amount{
				Value: senderAmount,
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

	// recipient
	if s.GetRecipient() != "" {
		ops = append(ops, &types.Operation{
			OperationIdentifier: &types.OperationIdentifier{
				NetworkIndex: nil,
			},
			RelatedOperations: nil,
			Type:              s.GetType().String(),
			Status:            StatusSuccess,
			Account: &types.AccountIdentifier{
				Address:    s.GetRecipient(),
				SubAccount: nil,
				Metadata:   nil,
			},
			Amount: &types.Amount{
				Value: s.GetAmount(),
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

	return ops
}

func (c *grpcIoTexClient) GetBlockTransaction(ctx context.Context, actionHash string) (ret *types.Transaction, err error) {
	if err = c.connect(); err != nil {
		return
	}
	return c.getBlockTransaction(ctx, actionHash)
}

func (c *grpcIoTexClient) getBlockTransaction(ctx context.Context, actionHash string) (ret *types.Transaction, err error) {
	request := &iotexapi.GetTransactionLogByActionHashRequest{
		ActionHash: actionHash,
	}

	resp, err := c.client.GetTransactionLogByActionHash(ctx, request)
	if err != nil {
		return nil, err
	}
	if resp.TransactionLog == nil {
		return nil, errors.New("not found")
	}
	ret = c.packTransaction(hex.EncodeToString(resp.TransactionLog.ActionHash), resp.TransactionLog.Transactions)

	rec, err := c.client.GetReceiptByAction(context.Background(), &iotexapi.GetReceiptByActionRequest{ActionHash: actionHash})
	if err == nil && rec.ReceiptInfo.Receipt.GetStatus() != 1 {
		oprs := c.covertAddressAmountsToOperations(addressAmountList{&addressAmount{
			senderAddr:   rec.ReceiptInfo.Receipt.ContractAddress,
			dstAddr:      address.RewardingPoolAddr,
			senderAmount: "-" + strconv.Itoa(int(rec.ReceiptInfo.Receipt.GasConsumed)),
			dstAmount:    strconv.Itoa(int(rec.ReceiptInfo.Receipt.GasConsumed)),
			actionType:   iotextypes.TransactionLogType_GAS_FEE.String(),
		}}, StatusFail)
		ret.Operations = append(ret.Operations, oprs...)
	}
	return
}

func (c *grpcIoTexClient) GetMemPool(ctx context.Context, actionHashes []string) (ret []*types.TransactionIdentifier, err error) {
	if err = c.connect(); err != nil {
		return
	}
	return c.getMemPool(ctx, actionHashes)
}

func (c *grpcIoTexClient) getMemPool(ctx context.Context, actionHashes []string) (ret []*types.TransactionIdentifier, err error) {
	request := &iotexapi.GetActPoolActionsRequest{
		ActionHashes: actionHashes,
	}

	resp, err := c.client.GetActPoolActions(ctx, request)
	if err != nil {
		return nil, err
	}
	for _, act := range resp.Actions {
		byteAct, err := proto.Marshal(act)
		if err != nil {
			return nil, err
		}
		h := hash.Hash256b(byteAct)
		ret = append(ret, &types.TransactionIdentifier{
			Hash: hex.EncodeToString(h[:]),
		})
	}

	return
}

func (c *grpcIoTexClient) GetMemPoolTransaction(ctx context.Context, h string) (ret *types.Transaction, err error) {
	if err = c.connect(); err != nil {
		return
	}
	return c.getMemPoolTransaction(ctx, h)
}

func (c *grpcIoTexClient) getMemPoolTransaction(ctx context.Context, h string) (ret *types.Transaction, err error) {
	request := &iotexapi.GetActPoolActionsRequest{
		ActionHashes: []string{h},
	}

	acts, err := c.client.GetActPoolActions(ctx, request)
	if err != nil {
		return nil, err
	}
	if acts.Actions == nil || len(acts.Actions) < 1 {
		return nil, errors.New("action not found")
	}
	return c.packActionToTransaction(acts.Actions[0], h, StatusSuccess)
}

func (c *grpcIoTexClient) packActionToTransaction(act *iotextypes.Action, h, status string) (ret *types.Transaction, err error) {
	var aal addressAmountList
	aal, err = packActionToAddressAmounts(act)
	if err != nil {
		return
	}

	ret = &types.Transaction{
		TransactionIdentifier: &types.TransactionIdentifier{Hash: h},
		Operations:            c.covertAddressAmountsToOperations(aal, status),
	}
	return
}

func (c *grpcIoTexClient) covertAddressAmountsToOperations(amountList addressAmountList, status string) (ret []*types.Operation) {
	var index int64 = 0
	for _, aa := range amountList {
		sender := c.genOperation(aa.senderAddr, status, aa.senderAmount, aa.actionType, index)
		index++
		dst := c.genOperation(aa.dstAddr, status, aa.dstAmount, aa.actionType, index)
		index++
		ret = append(ret, sender, dst)
	}
	return ret
}

func (c *grpcIoTexClient) genOperation(addr, status, amount, actType string, index int64) *types.Operation {
	return &types.Operation{
		OperationIdentifier: &types.OperationIdentifier{
			Index:        index,
			NetworkIndex: nil,
		},
		RelatedOperations: nil,
		Type:              actType,
		Status:            status,
		Account: &types.AccountIdentifier{
			Address:    addr,
			SubAccount: nil,
			Metadata:   nil,
		},
		Amount: &types.Amount{
			Value: amount,
			Currency: &types.Currency{
				Symbol:   c.cfg.Currency.Symbol,
				Decimals: c.cfg.Currency.Decimals,
				Metadata: nil,
			},
			Metadata: nil,
		},
		Metadata: nil,
	}
}

func packActionToAddressAmounts(act *iotextypes.Action) (aal addressAmountList, err error) {
	amount := "0"
	senderSign := "-"
	actionType := ""
	dst := ""
	callerAddr, err := getCaller(act)
	if err != nil {
		return aal, err
	}

	switch {
	case act.GetCore().GetTransfer() != nil:
		actionType = iotextypes.TransactionLogType_NATIVE_TRANSFER.String()
		amount = act.GetCore().GetTransfer().GetAmount()
		dst = act.GetCore().GetTransfer().GetRecipient()
	case act.GetCore().GetDepositToRewardingFund() != nil:
		actionType = iotextypes.TransactionLogType_DEPOSIT_TO_REWARDING_FUND.String()
		amount = act.GetCore().GetDepositToRewardingFund().GetAmount()
		dst = address.RewardingPoolAddr
	case act.GetCore().GetClaimFromRewardingFund() != nil:
		actionType = iotextypes.TransactionLogType_CLAIM_FROM_REWARDING_FUND.String()
		amount = act.GetCore().GetClaimFromRewardingFund().GetAmount()
		senderSign = "+"
		dst = address.RewardingPoolAddr
	case act.GetCore().GetStakeAddDeposit() != nil:
		actionType = iotextypes.TransactionLogType_DEPOSIT_TO_BUCKET.String()
		amount = act.GetCore().GetStakeAddDeposit().GetAmount()
		dst = address.StakingBucketPoolAddr
	case act.GetCore().GetStakeCreate() != nil:
		actionType = iotextypes.TransactionLogType_CREATE_BUCKET.String()
		amount = act.GetCore().GetStakeCreate().GetStakedAmount()
		dst = address.StakingBucketPoolAddr
	case act.GetCore().GetCandidateRegister() != nil:
		actionType = iotextypes.TransactionLogType_CANDIDATE_SELF_STAKE.String()
		amount = act.GetCore().GetCandidateRegister().GetStakedAmount()
		dst = address.StakingBucketPoolAddr
	case act.GetCore().GetExecution() != nil:
		actionType = iotextypes.TransactionLogType_IN_CONTRACT_TRANSFER.String()
		amount = act.GetCore().GetExecution().GetAmount()
		dst = address.StakingBucketPoolAddr
	}

	senderAmountWithSign := amount
	dstAmountWithSign := amount
	if senderSign == "-" {
		senderAmountWithSign = senderSign + amount
	} else {
		dstAmountWithSign = "-" + amount
	}

	fee := new(big.Int).SetUint64(act.GetCore().GetGasLimit())
	return addressAmountList{
		&addressAmount{
			senderAddr:   callerAddr.String(),
			dstAddr:      address.RewardingPoolAddr,
			senderAmount: "-" + fee.String(),
			dstAmount:    fee.String(),
			actionType:   iotextypes.TransactionLogType_GAS_FEE.String(),
		}, &addressAmount{
			senderAddr:   callerAddr.String(),
			dstAddr:      dst,
			senderAmount: senderAmountWithSign,
			dstAmount:    dstAmountWithSign,
			actionType:   actionType,
		},
	}, nil
}
