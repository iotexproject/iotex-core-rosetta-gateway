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
	acc, err := s.client.GetAccount(ctx, 0, request.AccountIdentifier.Address)
	if err != nil {
		return nil, ErrUnableToGetAccount
	}
	blk, err := s.client.GetLatestBlock(ctx)
	if err != nil {
		return nil, ErrUnableToGetBlk
	}

	md := make(map[string]interface{})
	md[NonceKey] = acc.Nonce

	resp := &types.AccountBalanceResponse{
		BlockIdentifier: &types.BlockIdentifier{
			Index: blk.Height,
			Hash:  blk.Hash,
		},
		Balances: []*types.Amount{
			&types.Amount{
				Value:    acc.Balance,
				Currency: IoTexCurrency,
			},
		},
		Metadata: &md,
	}
	return resp, nil
}
