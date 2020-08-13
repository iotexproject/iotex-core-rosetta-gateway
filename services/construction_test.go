// Copyright (c) 2020 IoTeX Foundation
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package services

import (
	"context"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/stretchr/testify/require"

	"github.com/iotexproject/iotex-core/action"

	"github.com/iotexproject/iotex-core-rosetta-gateway/config"
	ic "github.com/iotexproject/iotex-core-rosetta-gateway/iotex-client"
)

var (
	gasLimit = uint64(1000000)
	gasPrice = big.NewInt(0).SetUint64(1e9)
	nonce    = uint64(1)
)

func TestConstructionCombine(t *testing.T) {
	require := require.New(t)
	client, err := ic.NewIoTexClient(defaultCfg())
	require.NoError(err)

	// Start the server.
	api := NewConstructionAPIService(client)

	transfer, err := action.NewTransfer(nonce, big.NewInt(1), "io1l9vaqmanwj47tlrpv6etf3pwq0s0snsq4vxke2", []byte("payload"), gasLimit, gasPrice)
	require.NoError(err)
	bd := &action.EnvelopeBuilder{}
	elp := bd.SetNonce(nonce).
		SetGasPrice(gasPrice).
		SetGasLimit(gasLimit).
		SetAction(transfer).Build()

	req := &types.ConstructionCombineRequest{
		NetworkIdentifier:   defaultNetwork(),
		UnsignedTransaction: hex.EncodeToString(elp.Serialize()),
		Signatures: []*types.Signature{{

			//SigningPayload *SigningPayload `json:"signing_payload"`
			//PublicKey      *PublicKey      `json:"public_key"`
			//SignatureType  SignatureType   `json:"signature_type"`
			//Bytes          []byte
		},
		},
	}
	resp, typesErr := api.ConstructionCombine(context.Background(), req)
	require.True(typesErr.Code == 0)

}

func defaultCfg() *config.Config {
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
			Endpoint:       "api.testnet.iotex.one:443",
			SecureEndpoint: true,
			RosettaVersion: "1.4.1",
		},
	}
}

func defaultNetwork() *types.NetworkIdentifier {
	return &types.NetworkIdentifier{
		Blockchain: "IoTeX",
		Network:    "testnet",
	}
}
