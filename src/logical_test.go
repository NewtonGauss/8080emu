package main

import "testing"

func TestAnd(t *testing.T) {
	var table = []Pair{
		Pair{ // ANA B
			State{pc: 0, regA: 0b01010101, regB: 0b10101010, mem: [65536]byte{0xa0}, flags: 0b00001000}, // with carry, should unset
			State{pc: 1, regA: 0b00000000, regB: 0b10101010, mem: [65536]byte{0xa0}, flags: 0b00000101},
		},
		Pair{ // ANA B
			State{pc: 0, regA: 0b11010101, regB: 0b10001111, mem: [65536]byte{0xa0}, flags: 0b00001000}, // with carry, should unset
			State{pc: 1, regA: 0b10000101, regB: 0b10001111, mem: [65536]byte{0xa0}, flags: 0b00000010},
		},
	}
	doTest(t, table)
}

func TestXor(t *testing.T) {
	var table = []Pair{
		Pair{ // XRA B
			State{pc: 0, regA: 0b10101010, regB: 0b00001111, mem: [65536]byte{0xa8}, flags: 0b00001000}, // with carry, should unset
			State{pc: 1, regA: 0b10100101, regB: 0b00001111, mem: [65536]byte{0xa8}, flags: 0b00000110},
		},
	}
	doTest(t, table)
}

func TestOr(t *testing.T) {
	var table = []Pair{
		Pair{ // ORA B
			State{pc: 0, regA: 0b10101010, regB: 0b00001111, mem: [65536]byte{0xb0}, flags: 0b00001000}, // with carry, should unset
			State{pc: 1, regA: 0b10101111, regB: 0b00001111, mem: [65536]byte{0xb0}, flags: 0b00000110},
		},
	}
	doTest(t, table)
}

func TestCmp(t *testing.T) {
	var table = []Pair{
		Pair{ // CMP B
			State{pc: 0, regA: 0xA, regB: 0x5, mem: [65536]byte{0xb8}},
			State{pc: 1, regA: 0xA, regB: 0x5, mem: [65536]byte{0xb8}, flags: 0b00000100},
		},
		Pair{ // CMP C ... A = -0xb, C = 0x5
			State{pc: 0, regA: 0b11100101, regC: 0x5, mem: [65536]byte{0xb9}},
			State{pc: 1, regA: 0b11100101, regC: 0x5, mem: [65536]byte{0xb9}, flags: 0b00000010},
		},
	}
	doTest(t, table)
}
