package services

import (
	"context"
	"testing"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"
	"github.com/stretchr/testify/require"

	"github.com/iotexproject/iotex-core-rosetta-gateway/config"
	ic "github.com/iotexproject/iotex-core-rosetta-gateway/iotex-client"
)

func TestValidateNetworkIdentifier(t *testing.T) {
	require := require.New(t)
	var tests = []struct {
		cfg    *config.Config
		ni     *types.NetworkIdentifier
		expect *types.Error
	}{
		{
			&config.Config{
				NetworkIdentifier: config.NetworkIdentifier{
					Blockchain: "IoTeX",
				},
			},
			&types.NetworkIdentifier{
				Blockchain: "UNKNOWN",
			},
			ErrInvalidBlockchain,
		}, {
			&config.Config{
				NetworkIdentifier: config.NetworkIdentifier{
					Blockchain: "IoTeX",
				},
			},
			&types.NetworkIdentifier{
				Blockchain:           "IoTeX",
				SubNetworkIdentifier: &types.SubNetworkIdentifier{},
			},
			ErrInvalidSubnetwork,
		}, {
			&config.Config{
				NetworkIdentifier: config.NetworkIdentifier{
					Network: "testnet",
				},
			},
			&types.NetworkIdentifier{
				Network: "mainnet",
			},
			ErrInvalidNetwork,
		}, {
			&config.Config{
				NetworkIdentifier: config.NetworkIdentifier{
					Network: "testnet",
				},
			},
			nil,
			ErrMissingNID,
		}, {
			&config.Config{
				NetworkIdentifier: config.NetworkIdentifier{
					Blockchain: "IoTeX",
					Network:    "testnet",
				},
			},
			&types.NetworkIdentifier{
				Blockchain: "IoTeX",
				Network:    "testnet",
			},
			nil,
		},
	}
	for _, test := range tests {
		cli, err := ic.NewIoTexClient(test.cfg)
		require.NoError(err)
		typErr := ValidateNetworkIdentifier(context.Background(), cli, test.ni)
		require.Equal(test.expect, typErr)
	}
}

func TestIsSupportedConstructionType(t *testing.T) {
	require := require.New(t)
	var tests = []struct {
		typ    string
		expect bool
	}{
		{iotextypes.TransactionLogType_NATIVE_TRANSFER.String(), true},
		{"OTHERS", false},
	}
	for i, test := range tests {
		isSupported := IsSupportedConstructionType(test.typ)
		require.Equal(test.expect, isSupported, "index:", i)
	}
}
