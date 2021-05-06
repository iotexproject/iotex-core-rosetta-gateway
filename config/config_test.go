// Copyright (c) 2020 IoTeX Foundation
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	r := require.New(t)

	cfg, err := New("../docker/deploy/etc/iotex-rosetta/config.yaml")
	r.NoError(err)
	r.Equal("IoTeX", cfg.NetworkIdentifier.Blockchain)
	r.Equal("mainnet", cfg.NetworkIdentifier.Network)
	r.EqualValues(4689, cfg.NetworkIdentifier.EvmNetworkID)
	r.Equal("IOTX", cfg.Currency.Symbol)
	r.EqualValues(18, cfg.Currency.Decimals)
	r.Equal("8080", cfg.Server.Port)
	r.Equal("127.0.0.1:14014", cfg.Server.Endpoint)
	r.Equal(false, cfg.Server.SecureEndpoint)
	r.Equal("1.4.2", cfg.Server.RosettaVersion)
	r.Equal(false, cfg.KeepNoneTxAction)
}
