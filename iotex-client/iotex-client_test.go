package iotex_client

import (
	"context"
	"encoding/hex"
	"errors"
	"math/rand"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/iotexproject/go-pkgs/hash"
	"github.com/iotexproject/iotex-core/action"
	"github.com/iotexproject/iotex-core/blockchain/block"
	"github.com/iotexproject/iotex-core/pkg/unit"
	"github.com/iotexproject/iotex-core/test/identityset"
	"github.com/iotexproject/iotex-core/testutil"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"github.com/iotexproject/iotex-proto/golang/iotexapi/mock_iotexapi"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/iotexproject/iotex-core-rosetta-gateway/config"
)

func testServerAddr() string { return "127.0.0.1:14014" }

func testConfig() *config.Config {
	return &config.Config{
		NetworkIdentifier: config.NetworkIdentifier{
			Blockchain: "IoTeX",
			Network:    "testnet",
		},
		Currency: config.Currency{
			Symbol:   "IOTX",
			Decimals: 18,
		},
		Server: config.Server{
			Port:           "8080",
			Endpoint:       testServerAddr(),
			SecureEndpoint: false,
			RosettaVersion: "1.3.5",
		},
		KeepNoneTxAction: false,
	}
}

func testChain() []*iotextypes.BlockMeta {
	now := time.Now()
	return []*iotextypes.BlockMeta{
		{
			Hash:      "genesis hash",
			Height:    1,
			Timestamp: &timestamp.Timestamp{Seconds: now.Unix()},
		}, {
			Hash:      "hash 2",
			Height:    2,
			Timestamp: &timestamp.Timestamp{Seconds: now.Unix()},
		}, {
			Hash:      "hash 3",
			Height:    3,
			Timestamp: &timestamp.Timestamp{Seconds: now.Unix()},
		}, {
			Hash:      "hash 4",
			Height:    4,
			Timestamp: &timestamp.Timestamp{Seconds: now.Unix()},
		},
	}
}

func testChainMeta(chain []*iotextypes.BlockMeta) *iotextypes.ChainMeta {
	return &iotextypes.ChainMeta{
		Height: uint64(len(chain)),
	}
}

func testAccountMeta() *iotextypes.AccountMeta {
	return &iotextypes.AccountMeta{
		Address:      "address",
		Balance:      unit.ConvertIotxToRau(1000).String(),
		Nonce:        rand.Uint64(),
		PendingNonce: rand.Uint64(),
		NumActions:   rand.Uint64(),
	}
}

func testBlockIdentifier(chain []*iotextypes.BlockMeta) *iotextypes.BlockIdentifier {
	return &iotextypes.BlockIdentifier{
		Hash:   "block identifier hash",
		Height: uint64(rand.Intn(len(chain))),
	}
}

func testServerMeta() *iotextypes.ServerMeta {
	return &iotextypes.ServerMeta{
		PackageVersion:  "1.1.0",
		PackageCommitID: "commit id",
		GitStatus:       "ok",
		GoVersion:       "1.14.6",
		BuildTime:       time.Now().String(),
	}
}

func testTransactionLogs() *iotextypes.TransactionLogs {
	return &iotextypes.TransactionLogs{
		Logs: []*iotextypes.TransactionLog{
			{
				ActionHash:      []byte("f8cdb02b5f0d219d4aeebc47317a5639f74acac66bca44806040a841838bf2d9"),
				NumTransactions: 2,
				Transactions: []*iotextypes.TransactionLog_Transaction{
					{
						Topic:     []byte("topic 1"),
						Amount:    unit.ConvertIotxToRau(rand.Int63n(1)).String(),
						Sender:    "sender",
						Recipient: "recipient",
						Type:      0,
					}, {
						Topic:     []byte("topic 1"),
						Amount:    unit.ConvertIotxToRau(rand.Int63n(2)).String(),
						Sender:    "sender",
						Recipient: "recipient",
						Type:      0,
					},
				},
			}, {
				ActionHash:      []byte("7b1c902f5dadfc9d719faa20a7d2185c53d3f251380d18fb2857c690e9365d51"),
				NumTransactions: 1,
				Transactions: []*iotextypes.TransactionLog_Transaction{
					{
						Topic:     []byte("topic 1"),
						Amount:    unit.ConvertIotxToRau(rand.Int63n(2)).String(),
						Sender:    "sender",
						Recipient: "recipient",
						Type:      0,
					},
				},
			},
		},
	}
}

func testTransactionLog() *iotextypes.TransactionLog {
	return &iotextypes.TransactionLog{
		ActionHash:      []byte("f8cdb02b5f0d219d4aeebc47317a5639f74acac66bca44806040a841838bf2d9"),
		NumTransactions: 2,
		Transactions: []*iotextypes.TransactionLog_Transaction{
			{
				Topic:     []byte("topic 1"),
				Amount:    unit.ConvertIotxToRau(rand.Int63n(1)).String(),
				Sender:    "sender",
				Recipient: "recipient",
				Type:      0,
			},
		},
	}
}

func testActions() []*iotextypes.Action {
	senderPubKey, _ := hex.DecodeString("04403d3c0dbd3270ddfc248c3df1f9aafd60f1d8e7456961c9ef26292262cc68f0ea9690263bef9e197a38f06026814fc70912c2b98d2e90a68f8ddc5328180a01")
	signature, _ := hex.DecodeString("010203040506070809")
	return []*iotextypes.Action{
		{
			Core:         &iotextypes.ActionCore{
				Version:  1,
				Nonce:    10,
				GasLimit: 20010,
				GasPrice: "11000000000000000000",
				Action:   &iotextypes.ActionCore_Transfer{
					Transfer: &iotextypes.Transfer{
						Amount:    "1010000000000000000000",
						Recipient: "io1jh0ekmccywfkmj7e8qsuzsupnlk3w5337hjjg2",
						Payload:   nil,
					},
				},
			},
			SenderPubKey: senderPubKey,
			Signature:    signature,
		},
	}
}


func newMockServer(t *testing.T) (svr iotexapi.APIServiceServer, cli IoTexClient) {
	require := require.New(t)
	service := mock_iotexapi.NewMockAPIServiceServer(gomock.NewController(t))
	server := grpc.NewServer()
	iotexapi.RegisterAPIServiceServer(server, service)
	listener, err := net.Listen("tcp", testServerAddr())
	require.NoError(err)
	go func() {
		err := server.Serve(listener)
		if err != nil && err != grpc.ErrServerStopped {
			panic(err)
		}
	}()
	cli, err = NewIoTexClient(testConfig())
	require.NoError(err)

	chain := testChain()
	chainMeta := testChainMeta(chain)
	accountMeta := testAccountMeta()
	blockIdentifier := testBlockIdentifier(chain)
	serverMeta := testServerMeta()
	transactionLogs := testTransactionLogs()
	transactionLog := testTransactionLog()

	service.EXPECT().
		GetBlockMetas(gomock.Any(), gomock.AssignableToTypeOf(&iotexapi.GetBlockMetasRequest{})).
		DoAndReturn(func(ctx context.Context, req *iotexapi.GetBlockMetasRequest) (*iotexapi.GetBlockMetasResponse, error) {
			query := req.GetByIndex()
			if query == nil {
				return nil, errors.New("unsupported query method")
			}
			return &iotexapi.GetBlockMetasResponse{
				Total:    uint64(len(chain)),
				BlkMetas: chain[query.Start-1 : query.Start+query.Count-1],
			}, nil
		}).
		AnyTimes()
	service.EXPECT().
		GetChainMeta(gomock.Any(), gomock.AssignableToTypeOf(&iotexapi.GetChainMetaRequest{})).
		Return(&iotexapi.GetChainMetaResponse{ChainMeta: chainMeta}, nil).
		AnyTimes()
	service.EXPECT().
		GetAccount(gomock.Any(), gomock.AssignableToTypeOf(&iotexapi.GetAccountRequest{})).
		Return(&iotexapi.GetAccountResponse{
			AccountMeta:     accountMeta,
			BlockIdentifier: blockIdentifier,
		}, nil).
		AnyTimes()
	service.EXPECT().
		GetServerMeta(gomock.Any(), gomock.AssignableToTypeOf(&iotexapi.GetServerMetaRequest{})).
		Return(&iotexapi.GetServerMetaResponse{ServerMeta: serverMeta}, nil).
		AnyTimes()
	service.EXPECT().
		SendAction(gomock.Any(), gomock.AssignableToTypeOf(&iotexapi.SendActionRequest{})).
		DoAndReturn(func(ctx context.Context, req *iotexapi.SendActionRequest) (*iotexapi.SendActionResponse, error) {
			return &iotexapi.SendActionResponse{
				ActionHash: hex.EncodeToString(hash.ZeroHash256[:]),
			}, nil
		}).
		AnyTimes()
	service.EXPECT().
		GetActPoolActions(gomock.Any(), gomock.AssignableToTypeOf(&iotexapi.GetActPoolActionsRequest{})).
		DoAndReturn(func(ctx context.Context, req *iotexapi.GetActPoolActionsRequest) (*iotexapi.GetActPoolActionsResponse, error) {
			return &iotexapi.GetActPoolActionsResponse{
				Actions: testActions(),
			}, nil
		}).
		AnyTimes()

	topics := []hash.Hash256{
		hash.Hash256b([]byte("test")),
		hash.Hash256b([]byte("Pacific")),
		hash.Hash256b([]byte("Aleutian")),
	}
	testLog := &action.Log{
		Address:     "1",
		Data:        []byte("cd07d8a74179e032f030d9244"),
		BlockHeight: 1,
		ActionHash:  hash.ZeroHash256,
		Index:       1,
	}
	testLog.Topics = topics
	testLog.NotFixTopicCopyBug = true
	receipt := &action.Receipt{
		Status:          1,
		BlockHeight:     1,
		ActionHash:      hash.ZeroHash256,
		GasConsumed:     1,
		ContractAddress: "test",
	}
	receipt.AddLogs(testLog)
	ra := block.NewRunnableActionsBuilder().Build()
	blk, err := block.NewBuilder(ra).
		SetHeight(1).
		SetTimestamp(testutil.TimestampNow()).
		SetReceipts([]*action.Receipt{receipt}).
		SetPrevBlockHash(hash.ZeroHash256).
		SignAndBuild(identityset.PrivateKey(29))
	require.NoError(err)
	blkInfo := &iotexapi.BlockInfo{Block: blk.ConvertToBlockPb()}
	for _, receipt := range blk.Receipts {
		blkInfo.Receipts = append(blkInfo.Receipts, receipt.ConvertToReceiptPb())
	}
	service.EXPECT().
		GetRawBlocks(gomock.Any(), gomock.AssignableToTypeOf(&iotexapi.GetRawBlocksRequest{})).
		Return(&iotexapi.GetRawBlocksResponse{Blocks: []*iotexapi.BlockInfo{blkInfo}}, nil).
		AnyTimes()
	service.EXPECT().
		GetTransactionLogByBlockHeight(gomock.Any(), gomock.AssignableToTypeOf(&iotexapi.GetTransactionLogByBlockHeightRequest{})).
		Return(&iotexapi.GetTransactionLogByBlockHeightResponse{
			TransactionLogs: transactionLogs,
			BlockIdentifier: blockIdentifier,
		}, nil).
		AnyTimes()
	service.EXPECT().
		GetTransactionLogByActionHash(gomock.Any(), gomock.AssignableToTypeOf(&iotexapi.GetTransactionLogByActionHashRequest{})).
		Return(&iotexapi.GetTransactionLogByActionHashResponse{
			TransactionLog: transactionLog,
		}, nil).
		AnyTimes()
	t.Cleanup(server.Stop)
	return service, cli
}

func TestGrpcIoTexClient_GetChainID(t *testing.T) {
	require := require.New(t)
	_, cli := newMockServer(t)
	block, err := cli.GetChainID(context.Background())
	require.NoError(err)
	require.Equal(testConfig().NetworkIdentifier.Network, block)
}

func TestIoTexClient_GetBlock(t *testing.T) {
	var (
		require   = require.New(t)
		height    = int64(2)
		chain     = testChain()
		blk       = chain[height-1]
		preBlk    = chain[height-2]
		expectBlk = genBlock(preBlk, blk)
	)
	svr, cli := newMockServer(t)
	blockMetas, err := svr.GetBlockMetas(context.Background(), &iotexapi.GetBlockMetasRequest{
		Lookup: &iotexapi.GetBlockMetasRequest_ByIndex{
			ByIndex: &iotexapi.GetBlockMetasByIndexRequest{
				Start: uint64(height) - 1,
				Count: 2,
			},
		},
	})
	require.NoError(err)
	if len(blockMetas.BlkMetas) != 2 {
		require.Error(errors.New("unexpect blocks meta"))
	}
	block, err := cli.GetBlock(context.Background(), height)
	require.NoError(err)
	require.Equal(expectBlk, block)
}

func TestGrpcIoTexClient_GetLatestBlock(t *testing.T) {
	var (
		require    = require.New(t)
		chain      = testChain()
		lastBlk    = chain[len(chain)-1]
		preLastBlk = chain[len(chain)-2]
		expectBlk  = genBlock(preLastBlk, lastBlk)
	)
	_, cli := newMockServer(t)
	block, err := cli.GetLatestBlock(context.Background())
	require.NoError(err)
	require.Equal(expectBlk, block)
}

func TestGrpcIoTexClient_GetGenesisBlock(t *testing.T) {
	var (
		require   = require.New(t)
		chain     = testChain()
		genesis   = chain[0]
		expectBlk = genBlock(genesis, genesis)
	)
	_, cli := newMockServer(t)
	block, err := cli.GetGenesisBlock(context.Background())
	require.NoError(err)
	require.Equal(expectBlk, block)
}

func TestGrpcIoTexClient_GetAccount(t *testing.T) {
	require := require.New(t)
	_, cli := newMockServer(t)
	_, err := cli.GetAccount(context.Background(), 0, "")
	require.NoError(err)
}

func TestGrpcIoTexClient_SubmitTx(t *testing.T) {
	require := require.New(t)
	_, cli := newMockServer(t)
	tx, err := cli.SubmitTx(context.Background(), &iotextypes.Action{})
	require.NoError(err)
	require.Equal(hex.EncodeToString(hash.ZeroHash256[:]), tx)
}

func TestGrpcIoTexClient_GetStatus(t *testing.T) {
	require := require.New(t)
	svr, cli := newMockServer(t)
	expect, err := svr.GetChainMeta(context.Background(), &iotexapi.GetChainMetaRequest{})
	require.NoError(err)
	status, err := cli.GetStatus(context.Background())
	require.NoError(err)
	// require.Equal(chainMeta, status.GetChainMeta())
	require.Equal(expect.GetChainMeta().GetHeight(), status.GetChainMeta().GetHeight())
	require.Equal(expect.GetChainMeta().GetNumActions(), status.GetChainMeta().GetNumActions())
	require.Equal(expect.GetChainMeta().GetTps(), status.GetChainMeta().GetTps())
	require.Equal(expect.GetChainMeta().GetEpoch(), status.GetChainMeta().GetEpoch())
	require.Equal(expect.GetChainMeta().GetTpsFloat(), status.GetChainMeta().GetTpsFloat())
}

func TestGrpcIoTexClient_GetVersion(t *testing.T) {
	require := require.New(t)
	svr, cli := newMockServer(t)
	resp, err := svr.GetServerMeta(context.Background(), &iotexapi.GetServerMetaRequest{})
	serverMeta := resp.GetServerMeta()
	require.NoError(err)
	version, err := cli.GetVersion(context.Background())
	require.NoError(err)
	// require.Equal(serverMeta, version.GetServerMeta())
	require.Equal(serverMeta.GetPackageVersion(), version.GetServerMeta().GetPackageVersion())
	require.Equal(serverMeta.GetPackageCommitID(), version.GetServerMeta().GetPackageCommitID())
	require.Equal(serverMeta.GetGitStatus(), version.GetServerMeta().GetGitStatus())
	require.Equal(serverMeta.GetGoVersion(), version.GetServerMeta().GetGoVersion())
	require.Equal(serverMeta.GetBuildTime(), version.GetServerMeta().GetBuildTime())
}

func TestGrpcIoTexClient_GetTransactions(t *testing.T) {
	require := require.New(t)
	_, cli := newMockServer(t)
	// TODO: fill transactions data into mock server, and check it.
	transactions, err := cli.GetTransactions(context.Background(), 2)
	require.NoError(err)
	require.Equal([]*types.Transaction{}, transactions)
}

func TestGrpcIoTexClient_GetConfig(t *testing.T) {
	_, cli := newMockServer(t)
	config := cli.GetConfig()
	require.Equal(t, testConfig(), config)
}

func TestGrpcIoTexClient_GetBlockTransaction(t *testing.T) {
	require := require.New(t)
	_, cli := newMockServer(t)
	transaction, err := cli.GetBlockTransaction(context.Background(), "")
	tx := testTransactionLog()
	require.NoError(err)
	require.Equal(hex.EncodeToString(tx.ActionHash), transaction.TransactionIdentifier.Hash)
	require.Equal(tx.Transactions[0].Type.String(), transaction.Operations[0].Type)
	require.Equal(tx.Transactions[0].Amount, transaction.Operations[0].Amount.Value)
	require.Equal(tx.Transactions[0].Sender, transaction.Operations[0].Account.Address)
}

func TestGrpcIoTexClient_GetMemPool(t *testing.T) {
	require := require.New(t)
	_, cli := newMockServer(t)
	acts, err := cli.GetMemPool(context.Background(), []string{})
	require.NoError(err)
	require.Equal(true, len(acts) > 0)

	acts, err = cli.GetMemPool(context.Background(), []string{"322884fb04663019be6fb461d9453827487eafdd57b4de3bd89a7d77c9bf8395"})
	require.NoError(err)
	require.Equal(true, len(acts) > 0)
}

func TestGrpcIoTexClient_GetMemPoolTransaction(t *testing.T) {
	require := require.New(t)
	_, cli := newMockServer(t)

	trans, err := cli.GetMemPoolTransaction(context.Background(), "322884fb04663019be6fb461d9453827487eafdd57b4de3bd89a7d77c9bf8395")
	require.NoError(err)
	for i := range trans.Operations {
		act := testActions()[0].GetCore()
		oper := trans.Operations[i]
		if oper.Type == ActionTypeFee && act.GetTransfer().Recipient == oper.Account.Address {
			require.Equal(strconv.Itoa(int(act.GasLimit)), oper.Amount.Value)
		}
		if oper.Type == Transfer && act.GetTransfer().Recipient == oper.Account.Address {
			require.Equal(act.GetTransfer().Amount, oper.Amount.Value)
		}
	}
}
