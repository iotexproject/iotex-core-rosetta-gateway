// Copyright (c) 2020 IoTeX Foundation
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

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
