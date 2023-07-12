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

package asm

import (
	"encoding/hex"
	"testing"
)

// Tests disassembling the instructions for valid evm code
func TestInstructionIteratorValid(t *testing.T) {
	cnt := 0
	// DecodeString은 16진수 문자열 s로 표시되는 바이트를 반환합니다.
	// DecodeString은 src에 16진수 문자만 포함되고 길이가 짝수일 것으로 예상합니다.
	// 입력이 잘못된 경우, DecodeString은 오류 이전에 디코딩된 바이트를 반환합니다.

	script, _ := hex.DecodeString("61000000")

	it := NewInstructionIterator(script)
	for it.Next() {
		cnt++
	}

	if err := it.Error(); err != nil {
		t.Errorf("Expected 2, but encountered error %v instead.", err)
	}
	if cnt != 2 {
		t.Errorf("Expected 2, but got %v instead.", cnt)
	}
}

// Tests disassembling the instructions for invalid evm code
func TestInstructionIteratorInvalid(t *testing.T) {
	cnt := 0
	script, _ := hex.DecodeString("6100")

	it := NewInstructionIterator(script)
	for it.Next() {
		cnt++
	}

	if it.Error() == nil {
		t.Errorf("Expected an error, but got %v instead.", cnt)
	}
}

// Tests disassembling the instructions for empty evm code
func TestInstructionIteratorEmpty(t *testing.T) {
	cnt := 0
	script, _ := hex.DecodeString("")

	it := NewInstructionIterator(script)
	for it.Next() {
		cnt++
	}

	if err := it.Error(); err != nil {
		t.Errorf("Expected 0, but encountered error %v instead.", err)
	}
	if cnt != 0 {
		t.Errorf("Expected 0, but got %v instead.", cnt)
	}
}
