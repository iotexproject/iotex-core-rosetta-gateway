package services

import (
	"math/big"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"
)

func checkTransferOps(ops []*types.Operation, currency *types.Currency) *types.Error {
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
	if symbol != currency.Symbol || decimals != currency.Decimals {
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
	return nil
}

func ioTransferToOps(sender string, transfer *iotextypes.Transfer, currency *types.Currency) []*types.Operation {
	typ := iotextypes.TransactionLogType_NATIVE_TRANSFER.String()
	return []*types.Operation{
		&types.Operation{
			OperationIdentifier: &types.OperationIdentifier{
				Index: 0,
			},
			Type: typ,
			Account: &types.AccountIdentifier{
				Address: sender,
			},
			Amount: &types.Amount{
				Value:    "-" + transfer.GetAmount(),
				Currency: currency,
			},
		},
		&types.Operation{
			OperationIdentifier: &types.OperationIdentifier{
				Index: 1,
			},
			RelatedOperations: []*types.OperationIdentifier{&types.OperationIdentifier{
				Index: 0,
			}},
			Type: typ,
			Account: &types.AccountIdentifier{
				Address: transfer.GetRecipient(),
			},
			Amount: &types.Amount{
				Value:    transfer.GetAmount(),
				Currency: currency,
			},
		},
	}
}

func opsToIoTransfer(ops []*types.Operation) *iotextypes.Transfer {
	return &iotextypes.Transfer{
		Amount:    ops[1].Amount.Value,
		Recipient: ops[1].Account.Address,
	}
}
