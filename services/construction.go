// Copyright (c) 2020 IoTeX Foundation
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package services

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/protobuf/proto"

	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"

	ic "github.com/iotexproject/iotex-core-rosetta-gateway/iotex-client"
)

const (
	// OptionsIDKey is the name of the key in the Options map inside a
	// ConstructionMetadataRequest that specifies the account ID.
	OptionsIDKey = "id"
	CurveType    = "secp256k1"
)

type constructionAPIService struct {
	client ic.IoTexClient
}

// NewConstructionAPIService creates a new instance of an ConstructionAPIService.
func NewConstructionAPIService(client ic.IoTexClient) server.ConstructionAPIServicer {
	return &constructionAPIService{
		client: client,
	}
}

// ConstructionCombine implements the /construction/combine endpoint.
func (s *constructionAPIService) ConstructionCombine(
	ctx context.Context,
	request *types.ConstructionCombineRequest,
) (*types.ConstructionCombineResponse, *types.Error) {
	terr := ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier)
	if terr != nil {
		return nil, terr
	}
	ac := &iotextypes.ActionCore{}
	if err := proto.Unmarshal([]byte(request.UnsignedTransaction), ac); err != nil {
		return nil, ErrUnmarshal
	}
	sealed := &iotextypes.Action{
		Core:         ac,
		SenderPubKey: request.Signatures[0].PublicKey.Bytes,
		Signature:    request.Signatures[0].Bytes,
	}
	return &types.ConstructionCombineResponse{
		SignedTransaction: sealed.String(),
	}, nil
}

// ConstructionDerive implements the /construction/derive endpoint.
func (s *constructionAPIService) ConstructionDerive(
	ctx context.Context,
	request *types.ConstructionDeriveRequest,
) (*types.ConstructionDeriveResponse, *types.Error) {
	terr := ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier)
	if terr != nil {
		return nil, terr
	}
	//todo sm2 if need
	if request.PublicKey.Bytes != nil && request.PublicKey.CurveType == CurveType {
		addr, err := address.FromBytes(request.PublicKey.Bytes)
		if err != nil {
			return nil, ErrInvalidPublicKey
		}
		meta, typesErr := s.metadata(ctx, 0, addr.String())
		if typesErr != nil {
			return nil, typesErr
		}
		return &types.ConstructionDeriveResponse{Address: addr.String(), Metadata: meta.Metadata}, nil
	}
	return nil, ErrUnsupportedPublicKeyType
}

// ConstructionHash implements the /construction/hash endpoint.
func (s *constructionAPIService) ConstructionHash(
	ctx context.Context,
	request *types.ConstructionHashRequest,
) (*types.TransactionIdentifierResponse, *types.Error) {
	terr := ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier)
	if terr != nil {
		return nil, terr
	}
	digest := crypto.Keccak256([]byte(request.SignedTransaction))
	return &types.TransactionIdentifierResponse{
		TransactionIdentifier: &types.TransactionIdentifier{
			hex.EncodeToString(digest),
		},
	}, nil
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
	idRaw, ok := request.Options[OptionsIDKey]
	if !ok {
		return nil, ErrInvalidAccountAddress
	}
	idString, ok := idRaw.(string)
	if !ok {
		return nil, ErrInvalidAccountAddress
	}

	return s.metadata(ctx, 0, idString)
}

func (s *constructionAPIService) metadata(ctx context.Context, height int64, addr string) (*types.ConstructionMetadataResponse, *types.Error) {
	accout, err := s.client.GetAccount(ctx, height, addr)
	if err != nil {
		return nil, ErrUnableToGetNextNonce
	}

	// Return next nonce that should be used to sign transactions for given account.
	resp := &types.ConstructionMetadataResponse{
		Metadata: accout.Metadata,
	}

	return resp, nil
}

// ConstructionParse implements the /construction/parse endpoint.
func (s *constructionAPIService) ConstructionParse(
	ctx context.Context,
	request *types.ConstructionParseRequest,
) (*types.ConstructionParseResponse, *types.Error) {
	terr := ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier)
	if terr != nil {
		return nil, terr
	}
	act := iotextypes.Action{}
	tran, err := hex.DecodeString(request.Transaction)
	if err != nil {
		return nil, ErrUnableToParseTx
	}
	err = proto.Unmarshal(tran, &act)
	if err != nil {
		return nil, ErrUnableToParseTx
	}
	sender, err := address.FromBytes(act.GetSenderPubKey())
	if err != nil {
		return nil, ErrInvalidPublicKey
	}
	recipient := act.GetCore().GetTransfer().GetRecipient()
	amount := act.GetCore().GetTransfer().GetAmount()
	currency := &types.Currency{
		Symbol:   s.client.GetConfig().Currency.Symbol,
		Decimals: s.client.GetConfig().Currency.Decimals,
		Metadata: nil,
	}
	metadata := make(map[string]interface{})
	metadata["gasLimit"] = fmt.Sprintf("%d", act.GetCore().GetGasLimit())
	metadata["gasPrice"] = act.GetCore().GetGasPrice()
	return &types.ConstructionParseResponse{
		Operations: []*types.Operation{
			{
				OperationIdentifier: &types.OperationIdentifier{Index: 0},
				RelatedOperations:   nil,
				Type:                actionType(act.GetCore()),
				Status:              "SUCCESS", // todo check what does this mean
				Account:             &types.AccountIdentifier{Address: sender.String()},
				Amount:              &types.Amount{Value: "-" + amount, Currency: currency},
				CoinChange:          &types.CoinChange{},
				Metadata:            nil,
			},
			{
				OperationIdentifier: &types.OperationIdentifier{Index: 1},
				RelatedOperations:   []*types.OperationIdentifier{{Index: 0}},
				Type:                actionType(act.GetCore()),
				Status:              "SUCCESS", // todo check what does this mean
				Account:             &types.AccountIdentifier{Address: recipient},
				Amount:              &types.Amount{Value: "+" + amount, Currency: currency},
				CoinChange:          &types.CoinChange{},
				Metadata:            metadata,
			},
		},
		Signers: []string{sender.String()},
	}, nil
}

// ConstructionPayloads implements the /construction/payloads endpoint.
func (s *constructionAPIService) ConstructionPayloads(
	ctx context.Context,
	request *types.ConstructionPayloadsRequest,
) (*types.ConstructionPayloadsResponse, *types.Error) {
	terr := ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier)
	if terr != nil {
		return nil, terr
	}
	return nil, nil
}

// ConstructionPreprocess implements the /construction/preprocess endpoint.
func (s *constructionAPIService) ConstructionPreprocess(
	ctx context.Context,
	request *types.ConstructionPreprocessRequest,
) (*types.ConstructionPreprocessResponse, *types.Error) {
	terr := ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier)
	if terr != nil {
		return nil, terr
	}
	options := make(map[string]interface{})
	options["amount"] = request.Operations[1].Amount.Value
	options["decimals"] = request.Operations[1].Amount.Currency.Decimals
	options["fromAddr"] = request.Operations[0].Account.Address
	options["gasLimit"] = request.Operations[1].Metadata["gasLimit"]
	options["gasPrice"] = request.Operations[1].Metadata["gasPrice"]
	options["payer"] = request.Operations[0].Account.Address
	options["symbol"] = request.Operations[1].Amount.Currency.Symbol
	options["toAddr"] = request.Operations[1].Account.Address

	return &types.ConstructionPreprocessResponse{
		Options: options,
	}, nil
}

// ConstructionSubmit implements the /construction/submit endpoint.
func (s *constructionAPIService) ConstructionSubmit(
	ctx context.Context,
	request *types.ConstructionSubmitRequest,
) (*types.TransactionIdentifierResponse, *types.Error) {
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

	resp := &types.TransactionIdentifierResponse{
		TransactionIdentifier: &types.TransactionIdentifier{
			Hash: txID,
		},
	}

	return resp, nil
}

func actionType(pbAct *iotextypes.ActionCore) string {
	switch {
	case pbAct.GetTransfer() != nil:
		return "transfer"
	default:
		return ""
	}
}
