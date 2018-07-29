// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package math

import (
	"bytes"
	"math/big"
	"testing"
)

func TestPaddedBigBytes(t *testing.T) {
	tests := []struct {
		num    *big.Int
		n      int
		result []byte
	}{
		{num: big.NewInt(0), n: 4, result: []byte{0, 0, 0, 0}},
		{num: big.NewInt(1), n: 4, result: []byte{0, 0, 0, 1}},
		{num: big.NewInt(512), n: 4, result: []byte{0, 0, 2, 0}},
		{num: BigPow(2, 32), n: 4, result: []byte{1, 0, 0, 0, 0}},
	}
	for _, test := range tests {
		if result := PaddedBigBytes(test.num, test.n); !bytes.Equal(result, test.result) {
			t.Errorf("PaddedBigBytes(%d, %d) = %v, want %v", test.num, test.n, result, test.result)
		}
	}
}
