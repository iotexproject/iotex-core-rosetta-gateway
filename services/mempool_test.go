package services

import (
	"context"
	"testing"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/golang/mock/gomock"
	"github.com/iotexproject/iotex-core-rosetta-gateway/iotex-client/mock_client"
	"github.com/stretchr/testify/require"
)

func TestMemPoolAPIService_Mempool(t *testing.T) {
	var (
		cfg               = testConfig()
		networkIdentifier = &types.NetworkIdentifier{
			Blockchain: "IoTeX",
			Network:    "testnet",
		}
		tis = []*types.TransactionIdentifier{
			{Hash: "322884fb04663019be6fb461d9453827487eafdd57b4de3bd89a7d77c9bf8395"},
		}
		require = require.New(t)
		ctrl    = gomock.NewController(t)
		cli     = mock_client.NewMockIoTexClient(ctrl)
		clt     = NewMemPoolAPIService(cli)
	)
	cli.EXPECT().GetConfig().Return(cfg).AnyTimes()
	cli.EXPECT().GetMemPool(gomock.Any(), gomock.Any()).
		Return(tis, nil).
		AnyTimes()

	resp, typErr := clt.Mempool(context.Background(), &types.NetworkRequest{
		NetworkIdentifier: networkIdentifier,
	})
	require.Nil(typErr)
	require.Equal(tis, resp.TransactionIdentifiers)
}

func TestMemPoolAPIService_MempoolTransaction(t *testing.T) {
	var (
		cfg               = testConfig()
		networkIdentifier = &types.NetworkIdentifier{
			Blockchain: "IoTeX",
			Network:    "testnet",
		}
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
		tis = &types.TransactionIdentifier{
			Hash: "322884fb04663019be6fb461d9453827487eafdd57b4de3bd89a7d77c9bf8395",
		}
		require = require.New(t)
		ctrl    = gomock.NewController(t)
		cli     = mock_client.NewMockIoTexClient(ctrl)
		clt     = NewMemPoolAPIService(cli)
	)
	cli.EXPECT().GetConfig().Return(cfg).AnyTimes()
	cli.EXPECT().GetMemPoolTransaction(gomock.Any(), gomock.Any()).
		Return(tx, nil).
		AnyTimes()

	resp, typErr := clt.MempoolTransaction(context.Background(), &types.MempoolTransactionRequest{
		NetworkIdentifier:     networkIdentifier,
		TransactionIdentifier: tis,
	})

	require.Nil(typErr)
	require.Equal(tx, resp.Transaction)
}
