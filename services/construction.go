// Copyright (c) 2020 IoTeX Foundation
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package services

import (
	"context"
	"encoding/hex"
	"log"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/cast"

	"github.com/iotexproject/go-pkgs/crypto"
	"github.com/iotexproject/go-pkgs/hash"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"

	ic "github.com/iotexproject/iotex-core-rosetta-gateway/iotex-client"
)

const (
	// TODO config this
	CurveType     = "secp256k1"
	SignatureType = "ecdsa"
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
	if terr := ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier); terr != nil {
		return nil, terr
	}

	tran, err := hex.DecodeString(request.UnsignedTransaction)
	if err != nil {
		// TODO better error
		return nil, ErrUnmarshal
	}
	act := &iotextypes.Action{}
	if err := proto.Unmarshal(tran, act); err != nil {
		return nil, ErrUnmarshal
	}

	if len(request.Signatures) != 1 {
		// TODO better error
		return nil, ErrUnmarshal
	}

	// NOTE  clear out payload
	act.GetCore().GetTransfer().Payload = nil
	rawPub := request.Signatures[0].PublicKey.Bytes
	if btcec.IsCompressedPubKey(rawPub) {
		pubk, err := btcec.ParsePubKey(rawPub, btcec.S256())
		if err != nil {
			terr := ErrInvalidPublicKey
			terr.Message += err.Error()
			return nil, terr
		}
		rawPub = pubk.SerializeUncompressed()
	}

	act.SenderPubKey = rawPub
	rawSig := request.Signatures[0].Bytes
	log.Println(len(rawSig), rawSig)
	if len(rawSig) != 64 && len(rawSig) != 65 {
		terr := ErrInvalidPublicKey
		terr.Message += "nvalid signature format"
		return nil, terr
	}
	if len(rawSig) == 64 {
		rawSig = append(rawSig, 27)
	}
	act.Signature = rawSig

	msg, err := proto.Marshal(act)
	if err != nil {
		// TODO better error
		return nil, ErrUnmarshal
	}
	return &types.ConstructionCombineResponse{
		SignedTransaction: hex.EncodeToString(msg),
	}, nil
}

// ConstructionDerive implements the /construction/derive endpoint.
func (s *constructionAPIService) ConstructionDerive(
	ctx context.Context,
	request *types.ConstructionDeriveRequest,
) (*types.ConstructionDeriveResponse, *types.Error) {
	if terr := ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier); terr != nil {
		return nil, terr
	}

	if len(request.PublicKey.Bytes) == 0 || request.PublicKey.CurveType != CurveType {
		return nil, ErrUnsupportedPublicKeyType
	}

	rawPub := request.PublicKey.Bytes
	if btcec.IsCompressedPubKey(rawPub) {
		pubk, err := btcec.ParsePubKey(rawPub, btcec.S256())
		if err != nil {
			terr := ErrInvalidPublicKey
			terr.Message += err.Error()
			return nil, terr
		}
		rawPub = pubk.SerializeUncompressed()
	}

	pub, err := crypto.BytesToPublicKey(rawPub)
	if err != nil {
		terr := ErrInvalidPublicKey
		terr.Message += err.Error()
		return nil, terr
	}
	addr, err := address.FromBytes(pub.Hash())
	if err != nil {
		terr := ErrInvalidPublicKey
		terr.Message += err.Error()
		return nil, terr
	}
	return &types.ConstructionDeriveResponse{
		Address: addr.String(),
	}, nil
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
	tran, err := hex.DecodeString(request.SignedTransaction)
	if err != nil {
		// TODO better error
		return nil, ErrUnmarshal
	}
	h := hash.Hash256b(tran)

	return &types.TransactionIdentifierResponse{
		TransactionIdentifier: &types.TransactionIdentifier{
			hex.EncodeToString(h[:]),
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
	idRaw, ok := request.Options["sender"]
	if !ok {
		return nil, ErrInvalidAccountAddress
	}
	addr, ok := idRaw.(string)
	if !ok {
		return nil, ErrInvalidAccountAddress
	}
	account, err := s.client.GetAccount(ctx, 0, addr)
	if err != nil {
		log.Println("error get account here ", err)
		return nil, ErrUnableToGetNextNonce
	}
	meta := account.Metadata

	if _, ok := request.Options["gasLimit"]; !ok {
		// need a valid pubkey to estimate, just use one
		rawPrivKey, err := btcec.NewPrivateKey(btcec.S256())
		if err != nil {
			// TODO clean up error
			return nil, ErrUnableToGetNextNonce
		}
		rawPubKey := rawPrivKey.PubKey()
		gasLimit, err := s.client.EstimateGasForAction(ctx, &iotextypes.Action{
			Core: &iotextypes.ActionCore{
				Action: &iotextypes.ActionCore_Transfer{
					Transfer: &iotextypes.Transfer{},
				},
			},
			SenderPubKey: rawPubKey.SerializeUncompressed(),
		})
		if err != nil {
			// TODO clean up error
			log.Println("error get estimate gas here ", err)
			return nil, ErrUnableToGetNextNonce
		}
		meta["gasLimit"] = gasLimit
	}
	if _, ok := request.Options["gasPrice"]; !ok {
		gasPrice, err := s.client.SuggestGasPrice(ctx)
		if err != nil {
			// TODO clean up error
			log.Println("error get SuggestGasPrice gas here ", err)
			return nil, ErrUnableToGetNextNonce
		}
		meta["gasPrice"] = gasPrice
	}
	return &types.ConstructionMetadataResponse{
		Metadata: meta,
	}, nil
}

// ConstructionParse implements the /construction/parse endpoint.
func (s *constructionAPIService) ConstructionParse(
	ctx context.Context,
	request *types.ConstructionParseRequest,
) (*types.ConstructionParseResponse, *types.Error) {
	if terr := ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier); terr != nil {
		return nil, terr
	}
	tran, err := hex.DecodeString(request.Transaction)
	if err != nil {
		return nil, ErrUnableToParseTx
	}

	act := &iotextypes.Action{}
	if err := proto.Unmarshal(tran, act); err != nil {
		return nil, ErrUnableToParseTx
	}

	sender, terr := s.checkIoAction(act, request.Signed)
	if terr != nil {
		return nil, terr
	}
	ops, meta := s.ioActionToOps(sender, act)

	resp := &types.ConstructionParseResponse{
		Operations: ops,
		Metadata:   meta,
	}
	if request.Signed {
		resp.Signers = []string{sender}
	}
	return resp, nil
}

// ConstructionPayloads implements the /construction/payloads endpoint.
func (s *constructionAPIService) ConstructionPayloads(
	ctx context.Context,
	request *types.ConstructionPayloadsRequest,
) (*types.ConstructionPayloadsResponse, *types.Error) {
	if err := ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier); err != nil {
		return nil, err
	}
	log.Printf("%+v", request)
	if err := s.checkOperationAndMeta(request.Operations, request.Metadata, true); err != nil {
		return nil, err
	}

	act := s.opsToIoAction(request.Operations, request.Metadata)
	log.Printf("%+v", act)

	msg, err := proto.Marshal(act)
	if err != nil {
		// TODO better error
		return nil, ErrUnmarshal
	}

	// NOTE  unset payload here
	act.GetCore().GetTransfer().Payload = nil
	core, err := proto.Marshal(act.GetCore())
	if err != nil {
		// TODO better error
		return nil, ErrUnmarshal
	}
	h := hash.Hash256b(core)
	return &types.ConstructionPayloadsResponse{
		UnsignedTransaction: hex.EncodeToString(msg),
		Payloads: []*types.SigningPayload{
			&types.SigningPayload{
				Address:       request.Operations[0].Account.Address,
				Bytes:         h[:],
				SignatureType: SignatureType,
			},
		},
	}, nil
}

// ConstructionPreprocess implements the /construction/preprocess endpoint.
func (s *constructionAPIService) ConstructionPreprocess(
	ctx context.Context,
	request *types.ConstructionPreprocessRequest,
) (*types.ConstructionPreprocessResponse, *types.Error) {
	if err := ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier); err != nil {
		return nil, err
	}
	if err := s.checkOperationAndMeta(request.Operations, request.Metadata, false); err != nil {
		return nil, err
	}

	options := make(map[string]interface{})
	options["amount"] = request.Operations[1].Amount.Value
	options["symbol"] = request.Operations[1].Amount.Currency.Symbol
	options["decimals"] = request.Operations[1].Amount.Currency.Decimals
	options["sender"] = request.Operations[0].Account.Address
	options["recipient"] = request.Operations[1].Account.Address
	// TODO it is unclear where these meta data should be
	if request.Metadata["gasLimit"] != nil {
		options["gasLimit"] = request.Metadata["gasLimit"]
	}
	if request.Metadata["gasPrice"] != nil {
		options["gasPrice"] = request.Metadata["gasPrice"]
	}

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
	tran, err := hex.DecodeString(request.SignedTransaction)
	if err != nil {
		log.Println("hex", err)
		return nil, ErrUnableToSubmitTx
	}

	act := &iotextypes.Action{}
	if err := proto.Unmarshal(tran, act); err != nil {
		log.Println("proto", err)
		return nil, ErrUnableToSubmitTx
	}

	log.Printf("%+v", act)
	txID, err := s.client.SubmitTx(ctx, act)
	if err != nil {
		log.Println("grpc", err)
		return nil, ErrUnableToSubmitTx
	}
	log.Println("hash", txID)

	return &types.TransactionIdentifierResponse{
		TransactionIdentifier: &types.TransactionIdentifier{
			Hash: txID,
		},
	}, nil
}

func (s *constructionAPIService) opsToIoAction(ops []*types.Operation, meta map[string]interface{}) *iotextypes.Action {
	gasPrice := cast.ToUint64(meta["gasPrice"])
	return &iotextypes.Action{
		Core: &iotextypes.ActionCore{
			Action: &iotextypes.ActionCore_Transfer{
				Transfer: &iotextypes.Transfer{
					Amount:    ops[1].Amount.Value,
					Recipient: ops[1].Account.Address,
					// NOTE use payload to pass sender address, if user need to use payload,
					// need to marshal payload with sender address and real payload
					Payload: []byte(ops[0].Account.Address),
				},
			},
			GasLimit: cast.ToUint64(meta["gasLimit"]),
			GasPrice: new(big.Int).SetUint64(gasPrice).String(),
			Nonce:    cast.ToUint64(meta["nonce"]),
		},
	}
}

func (s *constructionAPIService) ioActionToOps(sender string, act *iotextypes.Action) ([]*types.Operation, map[string]interface{}) {
	meta := make(map[string]interface{})
	meta["nonce"] = act.GetCore().GetNonce()
	meta["gasLimit"] = act.GetCore().GetGasLimit()
	gasPrice, _ := new(big.Int).SetString(act.GetCore().GetGasPrice(), 10)
	meta["gasPrice"] = gasPrice.Uint64()

	ops := []*types.Operation{
		&types.Operation{
			OperationIdentifier: &types.OperationIdentifier{
				Index: 0,
			},
			Type: "NATIVE_TRANSFER",
			Account: &types.AccountIdentifier{
				Address: sender,
			},
			Amount: &types.Amount{
				Value: "-" + act.GetCore().GetTransfer().GetAmount(),
				Currency: &types.Currency{
					Symbol:   s.client.GetConfig().Currency.Symbol,
					Decimals: s.client.GetConfig().Currency.Decimals,
				},
			},
		},
		&types.Operation{
			OperationIdentifier: &types.OperationIdentifier{
				Index: 1,
			},
			RelatedOperations: []*types.OperationIdentifier{&types.OperationIdentifier{
				Index: 0,
			}},
			Type: "NATIVE_TRANSFER",
			Account: &types.AccountIdentifier{
				Address: act.GetCore().GetTransfer().GetRecipient(),
			},
			Amount: &types.Amount{
				Value: act.GetCore().GetTransfer().GetAmount(),
				Currency: &types.Currency{
					Symbol:   s.client.GetConfig().Currency.Symbol,
					Decimals: s.client.GetConfig().Currency.Decimals,
				},
			},
		},
	}

	return ops, meta
}

func (s *constructionAPIService) checkOperationAndMeta(ops []*types.Operation, meta map[string]interface{}, mustMeta bool) *types.Error {
	terr := ErrConstructionCheck
	if len(ops) != 2 {
		terr.Message += "operation numbers are no expected"
		return terr
	}

	// check amount
	if ops[0].Amount.Value != "-"+ops[1].Amount.Value {
		terr.Message += "amount value don't match"
		return terr
	}
	amountStr := ops[1].Amount.Value
	_, ok := new(big.Int).SetString(amountStr, 10)
	if !ok {
		terr.Message += "amount value is invalid"
		return terr
	}

	// check currency
	symbol := ops[1].Amount.Currency.Symbol
	decimals := ops[1].Amount.Currency.Decimals
	if symbol != s.client.GetConfig().Currency.Symbol ||
		decimals != s.client.GetConfig().Currency.Decimals {
		terr.Message += "invalid currency"
		return terr
	}

	// check address
	_, err := address.FromString(ops[0].Account.Address)
	if err != nil {
		terr.Message += "invalid sender address"
		return terr
	}
	_, err = address.FromString(ops[1].Account.Address)
	if err != nil {
		terr.Message += "invalid recipient address"
		return terr
	}

	// check metadata exists
	if mustMeta {
		if meta["gasLimit"] == nil || meta["gasPrice"] == nil || meta["nonce"] == nil {
			terr.Message += "metadata not complete"
			return terr
		}
	}

	// check gas
	if meta["gasLimit"] != nil {
		if _, err := cast.ToUint64E(meta["gasLimit"]); err != nil {
			terr.Message += "invalid gas limit"
			return terr
		}
	}
	if meta["gasPrice"] != nil {
		if _, err := cast.ToUint64E(meta["gasPrice"]); err != nil {
			terr.Message += "invalid gas price"
			return terr
		}
	}
	if meta["nonce"] != nil {
		if _, err := cast.ToUint64E(meta["nonce"]); err != nil {
			terr.Message += "invalid nonce"
			return terr
		}
	}
	return nil
}

func (s *constructionAPIService) checkIoAction(act *iotextypes.Action, signed bool) (sender string, terr *types.Error) {
	terr = ErrConstructionCheck
	if _, ok := new(big.Int).SetString(act.GetCore().GetGasPrice(), 10); !ok {
		terr.Message += "invalid gas price"
		return "", terr
	}

	if !signed {
		// NOTE use payload to pass sender address
		return string(act.GetCore().GetTransfer().GetPayload()), nil
	}
	act.GetCore().GetTransfer().Payload = nil

	// check pubkey and address
	if len(act.GetSenderPubKey()) == 0 {
		terr.Message += "invalid pub key"
		return "", terr
	}

	pub, err := crypto.BytesToPublicKey(act.GetSenderPubKey())
	if err != nil {
		terr.Message += "invalid pub key"
		return "", terr
	}
	senderAddr, err := address.FromBytes(pub.Hash())
	if err != nil {
		terr.Message += "invalid io address"
		return "", terr
	}
	sender = senderAddr.String()

	core, err := proto.Marshal(act.GetCore())
	if err != nil {
		// TODO better error
		return "", ErrUnmarshal
	}
	h := hash.Hash256b(core)
	if !pub.Verify(h[:], act.GetSignature()) {
		terr.Message += "invalid signature"
		return "", terr
	}
	return sender, nil
}

func actionType(pbAct *iotextypes.ActionCore) string {
	switch {
	case pbAct.GetTransfer() != nil:
		return "NATIVE_TRANSFER"
	default:
		return ""
	}
}
