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
// Note: we don't implement this, since we already return all transactions
// in the /mempool endpoint reponse above.
func (s *memPoolAPIService) MempoolTransaction(
	ctx context.Context,
	in *types.MempoolTransactionRequest,
	) (*types.MempoolTransactionResponse, *types.Error) {
	return nil, ErrNotImplemented
}