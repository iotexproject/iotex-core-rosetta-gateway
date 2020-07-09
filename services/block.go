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

type blockAPIService struct {
	client ic.IoTexClient
}

// NewBlockAPIService creates a new instance of an AccountAPIService.
func NewBlockAPIService(client ic.IoTexClient) server.BlockAPIServicer {
	return &blockAPIService{
		client: client,
	}
}

func (s *blockAPIService) Block(
	ctx context.Context,
	request *types.BlockRequest,
) (*types.BlockResponse, *types.Error) {
	terr := ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier)
	if terr != nil {
		return nil, terr
	}
	var height int64

	if request.BlockIdentifier != nil {
		if request.BlockIdentifier.Index != nil {
			height = *request.BlockIdentifier.Index
		} else if request.BlockIdentifier.Hash != nil {
			return nil, ErrMustQueryByIndex
		}
	}

	tblk, err := s.client.GetBlock(ctx, height)
	if err != nil {
		return nil, ErrUnableToGetBlk
	}
	tblk.Transactions, err = s.client.GetTransactions(ctx, height)
	if err != nil {
		return nil, ErrUnableToGetBlk
	}

	resp := &types.BlockResponse{
		Block: tblk,
	}

	return resp, nil
}

// BlockTransaction implements the /block/transaction endpoint.
// Note: we don't implement this, since we already return all transactions
// in the /block endpoint reponse above.
func (s *blockAPIService) BlockTransaction(
	ctx context.Context,
	request *types.BlockTransactionRequest,
) (*types.BlockTransactionResponse, *types.Error) {
	return nil, ErrNotImplemented
}
