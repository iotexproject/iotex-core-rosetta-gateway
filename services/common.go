// Copyright (c) 2020 IoTeX Foundation
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package services

import (
	"context"

	"github.com/coinbase/rosetta-sdk-go/types"

	ic "github.com/iotexproject/iotex-core-rosetta-gateway/iotex-client"
)

// IoTexCurrency is the currency used on the IoTex blockchain.
var IoTexCurrency = &types.Currency{
	Symbol:   "IoTex",
	Decimals: 18,
}

// ValidateNetworkIdentifier validates the network identifier.
func ValidateNetworkIdentifier(ctx context.Context, c ic.IoTexClient, ni *types.NetworkIdentifier) *types.Error {
	if ni != nil {
		cfg := c.GetConfig()
		if ni.Blockchain != cfg.Network_identifier.Blockchain {
			return ErrInvalidBlockchain
		}
		if ni.SubNetworkIdentifier != nil {
			return ErrInvalidSubnetwork
		}
		if ni.Network != cfg.Network_identifier.Network {
			return ErrInvalidNetwork
		}
	} else {
		return ErrMissingNID
	}
	return nil
}
