package iotex_client

import (
	"strings"
)

type (
	addressAmount struct {
		address string
		amount  string
	}
	addressAmountList []*addressAmount
)

func (l addressAmountList) Len() int      { return len(l) }
func (l addressAmountList) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l addressAmountList) Less(i, j int) bool {
	return strings.Compare(l[i].address, l[j].address) == 1
}
