package services

import (
	"context"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/golang/mock/gomock"
	"github.com/iotexproject/iotex-core-rosetta-gateway/iotex-client/mock_client"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMemPoolAPIService_Mempool(t *testing.T) {
	var (
		cfg     = testConfig()
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
