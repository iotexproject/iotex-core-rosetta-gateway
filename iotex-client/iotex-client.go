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
