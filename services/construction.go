// Copyright (c) 2020 IoTeX Foundation
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package services

import (
	"context"
	"encoding/hex"

	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/gogo/protobuf/proto"

	"github.com/iotexproject/iotex-proto/golang/iotextypes"

	ic "github.com/iotexproject/iotex-core-rosetta-gateway/iotex-client"
)

// OptionsIDKey is the name of the key in the Options map inside a
// ConstructionMetadataRequest that specifies the account ID.
const OptionsIDKey = "id"

// NonceKey is the name of the key in the Metadata map inside a
// ConstructionMetadataResponse that specifies the next valid nonce.
const NonceKey = "nonce"

type constructionAPIService struct {
	client ic.IoTexClient
}

// NewConstructionAPIService creates a new instance of an ConstructionAPIService.
func NewConstructionAPIService(client ic.IoTexClient) server.ConstructionAPIServicer {
	return &constructionAPIService{
		client: client,
	}
}

// ConstructionMetadata implements the /construction/metadata endpoint.
func (s *constructionAPIService) ConstructionMetadata(
	ctx context.Context,
	request *types.ConstructionMetadataRequest,
) (*types.ConstructionMetadataResponse, *types.Error) {
	terr := ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier)
	if terr != nil {
		return nil, terr
	}

	// Get the account ID field from the Options object.
	if request.Options == nil {
		return nil, ErrInvalidAccountAddress
	}
	idRaw, ok := (*request.Options)[OptionsIDKey]
	if !ok {
		return nil, ErrInvalidAccountAddress
	}
	idString, ok := idRaw.(string)
	if !ok {
		return nil, ErrInvalidAccountAddress
	}

	acc,_, err := s.client.GetAccount(ctx, 0, idString)
	if err != nil {
		return nil, ErrUnableToGetNextNonce
	}

	// Return next nonce that should be used to sign transactions for given account.
	md := make(map[string]interface{})
	md[NonceKey] = acc.Nonce

	resp := &types.ConstructionMetadataResponse{
		Metadata: &md,
	}

	return resp, nil
}

// ConstructionSubmit implements the /construction/submit endpoint.
func (s *constructionAPIService) ConstructionSubmit(
	ctx context.Context,
	request *types.ConstructionSubmitRequest,
) (*types.ConstructionSubmitResponse, *types.Error) {
	terr := ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier)
	if terr != nil {
		return nil, terr
	}
	act := iotextypes.Action{}
	tran, err := hex.DecodeString(request.SignedTransaction)
	if err != nil {
		return nil, ErrUnableToSubmitTx
	}
	err = proto.Unmarshal(tran, &act)
	if err != nil {
		return nil, ErrUnableToSubmitTx
	}
	txID, err := s.client.SubmitTx(ctx, &act)
	if err != nil {
		return nil, ErrUnableToSubmitTx
	}

	resp := &types.ConstructionSubmitResponse{
		TransactionIdentifier: &types.TransactionIdentifier{
			Hash: txID,
		},
	}

	return resp, nil
}
