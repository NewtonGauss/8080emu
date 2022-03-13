package main

import "testing"

func TestRotate(t *testing.T) {
	var table = []Pair{
		Pair{ // RLC
			State{pc: 0, regA: 0b00000001, mem: [65536]byte{0x07}},
			State{pc: 1, regA: 0b00000010, mem: [65536]byte{0x07}},
		},
		Pair{ // RLC
			State{pc: 0, regA: 0b10000000, mem: [65536]byte{0x07}},
			State{pc: 1, regA: 0b00000001, mem: [65536]byte{0x07}, flags: 0b00001000},
		},
		Pair{ // RRC
			State{pc: 0, regA: 0b10000000, mem: [65536]byte{0x0f}},
			State{pc: 1, regA: 0b01000000, mem: [65536]byte{0x0f}},
		},
		Pair{ // RRC
			State{pc: 0, regA: 0b00000001, mem: [65536]byte{0x0f}},
			State{pc: 1, regA: 0b10000000, mem: [65536]byte{0x0f}, flags: 0b00001000},
		},
		Pair{ // RAL
			State{pc: 0, regA: 0b10000000, mem: [65536]byte{0x1f}, flags: 0b00001000},
			State{pc: 1, regA: 0b00000001, mem: [65536]byte{0x1f}, flags: 0b00001000},
		},
		Pair{ // RAL
			State{pc: 0, regA: 0b00000001, mem: [65536]byte{0x1f}},
			State{pc: 1, regA: 0b00000010, mem: [65536]byte{0x1f}},
		},
		Pair{ // RAR
			State{pc: 0, regA: 0b10000000, mem: [65536]byte{0x1f}, flags: 0b00001000},
			State{pc: 1, regA: 0b11000000, mem: [65536]byte{0x1f}},
		},
		Pair{ // RAR
			State{pc: 0, regA: 0b00000001, mem: [65536]byte{0x1f}},
			State{pc: 1, regA: 0b00000000, mem: [65536]byte{0x1f}, flags: 0b00001000},
		},
	}
	doTest(t, table)
}
