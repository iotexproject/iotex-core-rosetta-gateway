package services

import (
	"context"
	"testing"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/iotexproject/iotex-core-rosetta-gateway/iotex-client/mock_client"
)

func TestBlockAPIService_Block(t *testing.T) {
	var (
		cfg     = testConfig()
		account = &types.AccountIdentifier{
			Address: "test account address",
		}
		amount = &types.Amount{
			Value: "1000",
			Currency: &types.Currency{
				Symbol:   "IOTX",
				Decimals: 18,
			},
		}
		tx = &types.Transaction{
			TransactionIdentifier: &types.TransactionIdentifier{
				Hash: "blah",
			},
			Operations: []*types.Operation{
				{
					OperationIdentifier: &types.OperationIdentifier{
						Index: int64(0),
					},
					Type:    "PAYMENT",
					Status:  types.String("SUCCESS"),
					Account: account,
					Amount:  amount,
				},
				{
					OperationIdentifier: &types.OperationIdentifier{
						Index: int64(1),
					},
					RelatedOperations: []*types.OperationIdentifier{
						{
							Index: int64(0),
						},
					},
					Type:    "PAYMENT",
					Status:  types.String("SUCCESS"),
					Account: account,
					Amount:  amount,
				},
			},
		}
		txs   = []*types.Transaction{tx}
		block = &types.Block{
			BlockIdentifier: &types.BlockIdentifier{
				Index: 100,
				Hash:  "block 100",
			},
			ParentBlockIdentifier: &types.BlockIdentifier{
				Index: 99,
				Hash:  "block 99",
			},
			Timestamp:    1000,
			Transactions: txs,
		}
		networkIdentifier = &types.NetworkIdentifier{
			Blockchain: "IoTeX",
			Network:    "testnet",
		}

		require = require.New(t)
		ctrl    = gomock.NewController(t)
		cli     = mock_client.NewMockIoTexClient(ctrl)
		clt     = NewBlockAPIService(cli)
		ret     = &types.BlockResponse{Block: block}
	)
	cli.EXPECT().GetConfig().Return(cfg).AnyTimes()
	cli.EXPECT().GetBlock(gomock.Any(), gomock.Any()).
		Return(block, nil).
		AnyTimes()
	cli.EXPECT().GetTransactions(gomock.Any(), gomock.Any()).
		Return(txs, nil).
		AnyTimes()

	resp, typErr := clt.Block(context.Background(), &types.BlockRequest{
		NetworkIdentifier: networkIdentifier,
		BlockIdentifier: &types.PartialBlockIdentifier{
			Index: &block.BlockIdentifier.Index,
		},
	})
	require.Nil(typErr)
	require.Equal(ret, resp)
}

func TestBlockAPIService_BlockTransaction(t *testing.T) {
	var (
		cfg     = testConfig()
		account = &types.AccountIdentifier{
			Address: "test account address",
		}
		amount = &types.Amount{
			Value: "1000",
			Currency: &types.Currency{
				Symbol:   "IOTX",
				Decimals: 18,
			},
		}
		tx = &types.Transaction{
			TransactionIdentifier: &types.TransactionIdentifier{
				Hash: "blah",
			},
			Operations: []*types.Operation{
				{
					OperationIdentifier: &types.OperationIdentifier{
						Index: int64(0),
					},
					Type:    "PAYMENT",
					Status:  types.String("SUCCESS"),
					Account: account,
					Amount:  amount,
				},
				{
					OperationIdentifier: &types.OperationIdentifier{
						Index: int64(1),
					},
					RelatedOperations: []*types.OperationIdentifier{
						{
							Index: int64(0),
						},
					},
					Type:    "PAYMENT",
					Status:  types.String("SUCCESS"),
					Account: account,
					Amount:  amount,
				},
			},
		}
		blockIdentifier = &types.BlockIdentifier{
			Index: 100,
			Hash:  "block 100",
		}
		networkIdentifier = &types.NetworkIdentifier{
			Blockchain: "IoTeX",
			Network:    "testnet",
		}

		require = require.New(t)
		ctrl    = gomock.NewController(t)
		cli     = mock_client.NewMockIoTexClient(ctrl)
		clt     = NewBlockAPIService(cli)
	)
	cli.EXPECT().GetConfig().Return(cfg).AnyTimes()
	cli.EXPECT().GetBlockTransaction(gomock.Any(), gomock.Any()).
		Return(tx, nil).
		AnyTimes()

	resp, typErr := clt.BlockTransaction(context.Background(), &types.BlockTransactionRequest{
		NetworkIdentifier:     networkIdentifier,
		BlockIdentifier:       blockIdentifier,
		TransactionIdentifier: tx.TransactionIdentifier,
	})
	require.Nil(typErr)
	require.Equal(tx, resp.Transaction)
}
