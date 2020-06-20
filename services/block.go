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

	blk, err := s.client.GetBlock(ctx, height)
	if err != nil {
		return nil, ErrUnableToGetBlk
	}
	txns, err := s.client.GetTransactions(ctx, height)
	if err != nil {
		return nil, ErrUnableToGetBlk
	}
	tblk := &types.Block{
		BlockIdentifier: &types.BlockIdentifier{
			Index: blk.Height,
			Hash:  blk.Hash,
		},
		ParentBlockIdentifier: &types.BlockIdentifier{
			Index: blk.ParentHeight,
			Hash:  blk.ParentHash,
		},
		Timestamp:    blk.Timestamp,
		Transactions: txns,
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
