package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"math/big"

	"github.com/gogo/protobuf/proto"
	"github.com/iotexproject/go-pkgs/crypto"
	"github.com/pkg/errors"

	"github.com/iotexproject/iotex-antenna-go/v2/account"
	"github.com/iotexproject/iotex-core/action"
)

var (
	recipient  string
	privateKey string
	nonce      uint64
	gasLimit   uint64
	amount     string
	gasPrice   string
)

func init() {
	flag.StringVar(&recipient, "r", "", "recipient's address")
	flag.StringVar(&privateKey, "p", "", "sender's private key in hex string format")
	flag.Uint64Var(&nonce, "n", 1, "sender's nonce")
	flag.StringVar(&amount, "a", "0", "send amount")
	flag.Uint64Var(&gasLimit, "gasLimit", 1000000, "gas limit")
	flag.StringVar(&gasPrice, "gasPrice", "1000000000000", "gas price")
}
func main() {
	flag.Parse()

	fmt.Println("::", recipient, privateKey, nonce, gasLimit, amount, gasPrice)

	pri, err := account.HexStringToAccount(privateKey)
	if err != nil {
		panic("privatekey error")
	}
	gasprice, ok := new(big.Int).SetString(gasPrice, 10)
	if !ok {
		panic("gas price error")
	}
	amo, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		panic("amount error")
	}
	signed, err := signedTransfer(recipient, pri.PrivateKey(), nonce, amo, gasLimit, gasprice)
	if err != nil {
		panic("sign error ")
	}
	m, err := proto.Marshal(signed.Proto())
	if err != nil {
		panic("marshal error ")
	}
	fmt.Println(hex.EncodeToString(m))
}

func signedTransfer(recipientAddr string, senderPriKey crypto.PrivateKey, nonce uint64, amount *big.Int, gasLimit uint64, gasPrice *big.Int) (action.SealedEnvelope, error) {
	transfer, err := action.NewTransfer(nonce, amount, recipientAddr, nil, gasLimit, gasPrice)
	if err != nil {
		return action.SealedEnvelope{}, err
	}
	bd := &action.EnvelopeBuilder{}
	elp := bd.SetNonce(nonce).
		SetGasPrice(gasPrice).
		SetGasLimit(gasLimit).
		SetAction(transfer).Build()
	selp, err := action.Sign(elp, senderPriKey)
	if err != nil {
		return action.SealedEnvelope{}, errors.Wrapf(err, "failed to sign transfer %v", elp)
	}
	return selp, nil
}
