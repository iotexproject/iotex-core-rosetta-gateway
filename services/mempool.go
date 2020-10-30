package services

import (
	"context"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
	ic "github.com/iotexproject/iotex-core-rosetta-gateway/iotex-client"
)


type memPoolAPIService struct {
	client ic.IoTexClient
}

// NewMemPoolAPIService creates a new instance of an MemPoolAPIService.
func NewMemPoolAPIService(client ic.IoTexClient) server.MempoolAPIServicer {
	return &memPoolAPIService{
		client: client,
	}
}

// MemPool implements the /mempool endpoint.
func (s *memPoolAPIService) Mempool (
	ctx context.Context,
	request *types.NetworkRequest,
) (*types.MempoolResponse, *types.Error) {
	terr := ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier)
	if terr != nil {
		return nil, terr
	}

	memPool, err := s.client.GetMemPool(ctx, []string{})
	if err != nil {
		return nil, ErrUnableToGetMemPool
	}

	return &types.MempoolResponse{
		TransactionIdentifiers: memPool,
	}, nil
}

// MempoolTransaction implements the /mempool/transaction endpoint.
func (s *memPoolAPIService) MempoolTransaction(
	ctx context.Context,
	request *types.MempoolTransactionRequest,
	) (*types.MempoolTransactionResponse, *types.Error) {
	terr := ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier)
	if terr != nil {
		return nil, terr
	}


	trans, err := s.client.GetMemPoolTransaction(ctx, request.TransactionIdentifier.Hash)
	if err != nil {
		return nil, ErrUnableToGetMemPoolTx
	}

	return &types.MempoolTransactionResponse{
		Transaction: trans,
	}, nil
}