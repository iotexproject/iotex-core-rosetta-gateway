// Copyright (c) 2020 IoTeX Foundation
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package iotex_client

import (
	"bytes"
	"context"
	"encoding/hex"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/iotexproject/go-pkgs/crypto"
	"github.com/iotexproject/go-pkgs/hash"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"
)

func assertAction(act *iotextypes.Action, operations operationList) operationList {
	oper := &operation{}
	oper.amount = "0"
	switch {
	case act.GetCore().GetTransfer() != nil:
		oper.actionType = Transfer
		oper.amount = act.GetCore().GetTransfer().GetAmount()
		oper.dst = act.GetCore().GetTransfer().GetRecipient()
	case act.GetCore().GetDepositToRewardingFund() != nil:
		oper.actionType = DepositToRewardingFund
		oper.amount = act.GetCore().GetDepositToRewardingFund().GetAmount()
		oper.dst = RewardingAddress
	case act.GetCore().GetClaimFromRewardingFund() != nil:
		oper.actionType = ClaimFromRewardingFund
		oper.amount = act.GetCore().GetClaimFromRewardingFund().GetAmount()
		oper.isPositive = true
		oper.dst = RewardingAddress
	case act.GetCore().GetStakeAddDeposit() != nil:
		oper.actionType = StakeAddDeposit
		oper.amount = act.GetCore().GetStakeAddDeposit().GetAmount()
		oper.dst = StakingAddress
	case act.GetCore().GetStakeCreate() != nil:
		oper.actionType = StakeCreate
		oper.amount = act.GetCore().GetStakeCreate().GetStakedAmount()
		oper.dst = StakingAddress
	//case stakewithdraw already handled before this call
	case act.GetCore().GetCandidateRegister() != nil:
		oper.actionType = CandidateRegister
		oper.amount = act.GetCore().GetCandidateRegister().GetStakedAmount()
		oper.dst = StakingAddress
	}
	if oper.amount != "0" && oper.actionType != "" {
		operations = append(operations, oper)
	}

	return operations
}

func getContractAddress(ctx context.Context, h string, client iotexapi.APIServiceClient) (contractAddr string, err error) {
	// need to get contract address generated of this action hash
	responseReceipt, err := client.GetReceiptByAction(ctx, &iotexapi.GetReceiptByActionRequest{ActionHash: h})
	if err != nil {
		return
	}
	contractAddr = responseReceipt.GetReceiptInfo().GetReceipt().GetContractAddress()
	return
}

func getCaller(act *iotextypes.Action) (callerAddr address.Address, err error) {
	srcPub, err := crypto.BytesToPublicKey(act.GetSenderPubKey())
	if err != nil {
		return
	}
	callerAddr, err = address.FromBytes(srcPub.Hash())
	return
}

func getActionType(topic []byte) string {
	InContractTransfer := common.Hash{}
	BucketWithdrawAmount := hash.BytesToHash256([]byte("withdrawAmount"))
	switch {
	case bytes.Compare(topic, InContractTransfer[:]) == 0:
		return Execution
	case bytes.Compare(topic, BucketWithdrawAmount[:]) == 0:
		return StakeWithdraw
	}
	return ""
}

func fillIndex(ret []*types.Transaction) {
	for _, t := range ret {
		for i, oper := range t.Operations {
			oper.OperationIdentifier.Index = int64(i)
		}
	}
}

func getImplicitTransferLog(ctx context.Context, height int64, client iotexapi.APIServiceClient) (
	transferLogMap map[string][]*iotextypes.ImplicitTransferLog_Transaction, err error) {
	transferLogMap = make(map[string][]*iotextypes.ImplicitTransferLog_Transaction)
	transferLog, err := client.GetImplicitTransferLogByBlockHeight(
		ctx,
		&iotexapi.GetImplicitTransferLogByBlockHeightRequest{BlockHeight: uint64(height)},
	)

	if err == nil && transferLog.GetBlockImplicitTransferLog().GetNumTransactions() != 0 {
		for _, a := range transferLog.GetBlockImplicitTransferLog().GetImplicitTransferLog() {
			h := hex.EncodeToString(a.ActionHash)
			transferLogMap[h] = a.GetTransactions()
		}
	}
	return transferLogMap, nil
}
