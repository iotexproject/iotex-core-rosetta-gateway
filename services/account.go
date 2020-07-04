// Copyright (c) 2020 IoTeX Foundation
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package services

import (
	"context"

	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"

	ic "github.com/iotexproject/iotex-core-rosetta-gateway/iotex-client"
)

type accountAPIService struct {
	client ic.IoTexClient
}

// NewAccountAPIService creates a new instance of an AccountAPIService.
func NewAccountAPIService(client ic.IoTexClient) server.AccountAPIServicer {
	return &accountAPIService{
		client: client,
	}
}

// AccountBalance implements the /account/balance endpoint.
func (s *accountAPIService) AccountBalance(
	ctx context.Context,
	request *types.AccountBalanceRequest,
) (*types.AccountBalanceResponse, *types.Error) {
	terr := ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier)
	if terr != nil {
		return nil, terr
	}
	// TODO fix this when we have archive mode
	acc,blkIdentifier, err := s.client.GetAccount(ctx, 0, request.AccountIdentifier.Address)
	if err != nil {
		return nil, ErrUnableToGetAccount
	}

	md := make(map[string]interface{})
	md[NonceKey] = acc.Nonce

	resp := &types.AccountBalanceResponse{
		BlockIdentifier:blkIdentifier,
		Balances: []*types.Amount{
			&types.Amount{
				Value: acc.Balance,
				Currency: &types.Currency{
					Symbol:   s.client.GetConfig().Currency.Symbol,
					Decimals: s.client.GetConfig().Currency.Decimals,
					Metadata: nil,
				},
			},
		},
		Metadata: &md,
	}
	return resp, nil
}
