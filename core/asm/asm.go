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

// Provides support for dealing with EVM assembly instructions (e.g., disassembling them).
package asm

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/core/vm"
)

// 디스어셈블된 EVM 명령어를 위한 이터레이터
// Iterator for disassembled EVM instructions
type instructionIterator struct {
	code    []byte
	pc      uint64
	arg     []byte
	op      vm.OpCode
	error   error
	started bool
}

// Create a new instruction iterator.
// 새 명령어 반복자를 생성합니다.
func NewInstructionIterator(code []byte) *instructionIterator {
	it := new(instructionIterator)
	it.code = code
	return it
}

// Returns true if there is a next instruction and moves on.
// 다음 명령어가 있으면 참을 반환하고 계속 진행합니다.
func (it *instructionIterator) Next() bool {
	if it.error != nil || uint64(len(it.code)) <= it.pc {
		// 이전에 오류 또는 끝에 도달했습니다.
		// We previously reached an error or the end.
		return false
	}

	if it.started {
		// Since the iteration has been already started we move to the next instruction.
		// 반복이 이미 시작되었으므로 다음 명령어로 이동합니다.
		if it.arg != nil {
			it.pc += uint64(len(it.arg))
		}
		it.pc++
	} else {
		// We start the iteration from the first instruction.
		// 첫 번째 명령어부터 반복을 시작합니다.
		it.started = true
	}

	if uint64(len(it.code)) <= it.pc {
		// We reached the end.
		return false
	}

	it.op = vm.OpCode(it.code[it.pc])
	if it.op.IsPush() {
		a := uint64(it.op) - uint64(vm.PUSH1) + 1
		u := it.pc + 1 + a
		if uint64(len(it.code)) <= it.pc || uint64(len(it.code)) < u {
			it.error = fmt.Errorf("incomplete push instruction at %v", it.pc)
			return false
		}
		it.arg = it.code[it.pc+1 : u]
	} else {
		it.arg = nil
	}
	return true
}

// 발생했을 수 있는 모든 오류를 반환합니다.
// Returns any error that may have been encountered.
func (it *instructionIterator) Error() error {
	return it.error
}

// Returns the PC of the current instruction.
// 현재 명령의 PC를 반환합니다.
func (it *instructionIterator) PC() uint64 {
	return it.pc
}

// Returns the opcode of the current instruction.
// 현재 명령어의 연산 코드를 반환합니다.
func (it *instructionIterator) Op() vm.OpCode {
	return it.op
}

// Returns the argument of the current instruction.
// 현재 명령어의 인수를 반환합니다.
func (it *instructionIterator) Arg() []byte {
	return it.arg
}

// Pretty-print all disassembled EVM instructions to stdout.
// 분해된 모든 EVM 인스트럭션을 stdout에 예쁘게 인쇄합니다.
func PrintDisassembled(code string) error {
	script, err := hex.DecodeString(code)
	if err != nil {
		return err
	}

	it := NewInstructionIterator(script)
	for it.Next() {
		if it.Arg() != nil && 0 < len(it.Arg()) {
			fmt.Printf("%06v: %v 0x%x\n", it.PC(), it.Op(), it.Arg())
		} else {
			fmt.Printf("%06v: %v\n", it.PC(), it.Op())
		}
	}
	return it.Error()
}

// Return all disassembled EVM instructions in human-readable format.
// 분해된 모든 EVM 명령어를 사람이 읽을 수 있는 형식으로 반환합니다.
func Disassemble(script []byte) ([]string, error) {
	instrs := make([]string, 0)

	it := NewInstructionIterator(script)
	for it.Next() {
		if it.Arg() != nil && 0 < len(it.Arg()) {
			instrs = append(instrs, fmt.Sprintf("%06v: %v 0x%x\n", it.PC(), it.Op(), it.Arg()))
		} else {
			instrs = append(instrs, fmt.Sprintf("%06v: %v\n", it.PC(), it.Op()))
		}
	}
	if err := it.Error(); err != nil {
		return nil, err
	}
	return instrs, nil
}
