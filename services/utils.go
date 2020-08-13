package services

import "encoding/json"

type senderAddrPayload struct {
	SenderAddress string
	RealPayload   []byte
}

func marshalSenderAddrPayload(addr string, payload []byte) ([]byte, error) {
	pl := senderAddrPayload{
		SenderAddress: addr,
		RealPayload:   payload,
	}
	return json.Marshal(pl)
}

func unmarshalSenderAddrPayload(payload []byte) (string, []byte, error) {
	pl := senderAddrPayload{}
	if err := json.Unmarshal(payload, &pl); err != nil {
		return "", nil, err
	}
	return pl.SenderAddress, pl.RealPayload, nil
}
