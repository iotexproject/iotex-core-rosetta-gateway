package services

import (
	"context"
	"testing"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/iotexproject/iotex-core-rosetta-gateway/config"
	"github.com/iotexproject/iotex-core-rosetta-gateway/iotex-client/mock_client"
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
			RosettaVersion: "1.4.10",
		},
		KeepNoneTxAction: false,
	}
}

func TestAccountAPIService_AccountBalance(t *testing.T) {
	var (
		block = &types.BlockIdentifier{
			Hash:  "block1",
			Index: 1,
		}
		amount = &types.Amount{
			Value: "100",
			Currency: &types.Currency{
				Symbol:   "IOTX",
				Decimals: 18,
			},
		}
		networkIdentifier = &types.NetworkIdentifier{
			Blockchain: "IoTeX",
			Network:    "testnet",
		}
		ret = &types.AccountBalanceResponse{
			BlockIdentifier: block,
			Balances:        []*types.Amount{amount},
		}

		require = require.New(t)
		ctrl    = gomock.NewController(t)
		cli     = mock_client.NewMockIoTexClient(ctrl)
		cfg     = testConfig()
	)
	cli.EXPECT().GetConfig().Return(cfg).AnyTimes()
	cli.EXPECT().GetAccount(gomock.Any(), gomock.Eq(int64(0)), gomock.Any()).
		Return(ret, nil).
		AnyTimes()

	clt := NewAccountAPIService(cli)
	resp, typErr := clt.AccountBalance(context.Background(), &types.AccountBalanceRequest{
		NetworkIdentifier: networkIdentifier,
		AccountIdentifier: &types.AccountIdentifier{
			Address: "io1d4c5lp4ea4754wy439g2t99ue7wryu5r2lslh2",
		},
	})
	require.Nil(typErr)
	require.Equal(ret, resp)
}
