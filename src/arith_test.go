package main

import "testing"

type Pair struct {
	init State // initial state
	exp  State // expected state
}

func TestArithmetic(t *testing.T) {
	var table = []Pair{
		Pair{ // INX B ... does not change flags
			State{regB: 0xff, regC: 0x00, pc: 0, mem: [65536]byte{0x03}, flags: 0b00000101},
			State{regB: 0xff, regC: 0x01, pc: 1, mem: [65536]byte{0x03}, flags: 0b00000101},
		},
		Pair{ // INR B ... does not affect Cy
			State{regB: 0xff, pc: 0, mem: [65536]byte{0x04}},
			State{regB: 0x00, pc: 1, mem: [65536]byte{0x04}, flags: 0b00000101},
		},
		Pair{ // DCR B
			State{regB: 0x01, pc: 0, mem: [65536]byte{0x05}},
			State{regB: 0x00, pc: 1, mem: [65536]byte{0x05}, flags: 0b00000101},
		},
		Pair{ // DAD B
			State{regB: 0x0f, regC: 0x0f, regH: 0x00, regL: 0x01, pc: 0, mem: [65536]byte{0x09}},
			State{regB: 0x0f, regC: 0x0f, regH: 0x0f, regL: 0x10, pc: 1, mem: [65536]byte{0x09}, flags: 0b00000000},
		},
		Pair{ // DAD D ... 0x0000 is even, but DAD does not affect Parity Bit. Cy is set.
			State{regD: 0xff, regE: 0x00, regH: 0x01, regL: 0x00, pc: 0, mem: [65536]byte{0x19}},
			State{regD: 0xff, regE: 0x00, regH: 0x00, regL: 0x00, pc: 1, mem: [65536]byte{0x19}, flags: 0b00001000},
		},
		Pair{ // DCX H
			State{regH: 0x01, regL: 0x00, pc: 0, mem: [65536]byte{0x2b}},
			State{regH: 0x00, regL: 0xff, pc: 1, mem: [65536]byte{0x2b}},
		},
		Pair{ // INX SP
			State{sp: 0x00ff, pc: 0, mem: [65536]byte{0x33}},
			State{sp: 0x0100, pc: 1, mem: [65536]byte{0x33}},
		},
		Pair{ // INR M
			State{regH: 0xff, regL: 0x00, pc: 0, mem: [65536]byte{0: 0x34, 0xff00: 2}},
			State{regH: 0xff, regL: 0x00, pc: 1, mem: [65536]byte{0: 0x34, 0xff00: 3}, flags: 0b00000100},
		},
		Pair{ // DCR M ... Does not affect Cy
			State{regH: 0xff, regL: 0x00, pc: 0, mem: [65536]byte{0: 0x35, 0xff00: 0}},
			State{regH: 0xff, regL: 0x00, pc: 1, mem: [65536]byte{0: 0x35, 0xff00: 0xff}, flags: 0b00000110},
		},
		Pair{ // DAD SP 0x0101 + 0x00FF = 0x0200
			State{regH: 0x01, regL: 0x01, sp: 0x00ff, pc: 0, mem: [65536]byte{0x39}},
			State{regH: 0x02, regL: 0x00, sp: 0x00ff, pc: 1, mem: [65536]byte{0x39}, flags: 0b00000000},
		},
		Pair{ // DCX SP
			State{sp: 0x00ff, pc: 0, mem: [65536]byte{0x3b}},
			State{sp: 0x00fe, pc: 1, mem: [65536]byte{0x3b}},
		},
		Pair{ // ADD B
			State{regA: 1, regB: 2, pc: 0, mem: [65536]byte{0x80}},
			State{regA: 3, regB: 2, pc: 1, mem: [65536]byte{0x80}, flags: 0b00000100},
		},
		Pair{ // ADD C ... 1 + (-1) -> Carry + Parity + Zero
			State{regA: 1, regC: 0b11111111, pc: 0, mem: [65536]byte{0x81}},
			State{regA: 0, regC: 0b11111111, pc: 1, mem: [65536]byte{0x81}, flags: 0b00001101},
		},
		Pair{ // ADD D ... 0 + (-2) -> Sign. No Parity
			State{regA: 0, regD: 0b11111110, pc: 0, mem: [65536]byte{0x82}},
			State{regA: 0b11111110, regD: 0b11111110, pc: 1, mem: [65536]byte{0x82}, flags: 0b00000010},
		},
		Pair{ // ADD M ... M = (HL)
			State{regA: 1, regH: 0xff, regL: 0x00, pc: 0, mem: [65536]byte{0: 0x86, 0xff00: 2}},
			State{regA: 3, regH: 0xff, regL: 0x00, pc: 1, mem: [65536]byte{0: 0x86, 0xff00: 2}, flags: 0b00000100},
		},
		Pair{ // ADD A
			State{regA: 1, pc: 0, mem: [65536]byte{0x87}},
			State{regA: 2, pc: 1, mem: [65536]byte{0x87}, flags: 0b00000000},
		},
		Pair{ // ADC B ... Should add 1 to regA (the carry)
			State{regA: 1, regB: 0, pc: 0, mem: [65536]byte{0x88}, flags: 0b00001000},
			State{regA: 2, regB: 0, pc: 1, mem: [65536]byte{0x88}, flags: 0b00000000},
		},
		Pair{ // ADC C ... Should not add 1 to regA (carry is not set)
			State{regA: 1, regC: 0, pc: 0, mem: [65536]byte{0x88}, flags: 0b00000000},
			State{regA: 1, regC: 0, pc: 1, mem: [65536]byte{0x88}, flags: 0b00000000},
		},
		Pair{ // SUB B regA - regB = 1 - 1 = 0
			State{regA: 1, regB: 1, pc: 0, mem: [65536]byte{0x90}},
			State{regA: 0, regB: 1, pc: 1, mem: [65536]byte{0x90}, flags: 0b00000101},
		},
		Pair{ // SUB C regA - regC = 1 - 2 = -1 = 0b11111111
			State{regA: 1, regC: 2, pc: 0, mem: [65536]byte{0x91}},
			State{regA: 0b11111111, regC: 2, pc: 1, mem: [65536]byte{0x91}, flags: 0b00001110},
		},
		Pair{ // SBB B regA - regB - Cy = 1 - 1 - 1 = -1 = 0b11111111 = 0xff
			State{regA: 1, regB: 1, pc: 0, mem: [65536]byte{0x98}, flags: 0b00001000},
			State{regA: 0xff, regB: 1, pc: 1, mem: [65536]byte{0x98}, flags: 0b00001110},
		},
		Pair{ // SBB C regA - regC - Cy = 1 - 0 - 1 = 0
			State{regA: 1, regC: 0, pc: 0, mem: [65536]byte{0x99}, flags: 0b00001000},
			State{regA: 0, regC: 0, pc: 1, mem: [65536]byte{0x99}, flags: 0b00000101},
		},
		Pair{ // SBB A regA - regA - Cy = 1 - 1 - 0 = 0
			State{regA: 1, pc: 0, mem: [65536]byte{0x9f}, flags: 0b00000000},
			State{regA: 0, pc: 1, mem: [65536]byte{0x9f}, flags: 0b00000101},
		},
		Pair{ // ADI D8
			State{regA: 1, pc: 0, mem: [65536]byte{0xC6, 2}},
			State{regA: 3, pc: 2, mem: [65536]byte{0xC6, 2}, flags: 0b00000100},
		},
		Pair{ // ACI D8 ... regA + 2 + Cy = 1 + 2 + 1 = 4
			State{regA: 1, pc: 0, mem: [65536]byte{0xce, 2}, flags: 0b00001000},
			State{regA: 4, pc: 2, mem: [65536]byte{0xce, 2}, flags: 0b00000000},
		},
		Pair{ // SUI D8 ... 1 - 2 = -1 = 0b11111111 = 0xff
			State{regA: 1, pc: 0, mem: [65536]byte{0xd6, 2}},
			State{regA: 0xff, pc: 2, mem: [65536]byte{0xd6, 2}, flags: 0b00001110},
		},
		Pair{ // SBI D8 ... A - 2 - Cy = 1 - 2 - 1= -2 = 0b11111110 = 0xfe
			State{regA: 1, pc: 0, mem: [65536]byte{0xde, 2}, flags: 0b00001000},
			State{regA: 0xfe, pc: 2, mem: [65536]byte{0xde, 2}, flags: 0b00001010},
		},
	}
	for _, test := range table {
		var env = test.init
		var opcode = env.mem[0]
		env.ExecInstruction()
		if env.regA != test.exp.regA {
			t.Errorf("[0x%02x] env.regA = %v, expected %v", opcode, env.regA, test.exp.regA)
		}
		if env.regB != test.exp.regB {
			t.Errorf("[0x%02x] env.regB = %v, expected %v", opcode, env.regB, test.exp.regB)
		}
		if env.regC != test.exp.regC {
			t.Errorf("[0x%02x] env.regC = %v, expected %v", opcode, env.regC, test.exp.regC)
		}
		if env.regD != test.exp.regD {
			t.Errorf("[0x%02x] env.regD = %v, expected %v", opcode, env.regD, test.exp.regD)
		}
		if env.regE != test.exp.regE {
			t.Errorf("[0x%02x] env.regE = %v, expected %v", opcode, env.regE, test.exp.regE)
		}
		if env.regH != test.exp.regH {
			t.Errorf("[0x%02x] env.regH = %v, expected %v", opcode, env.regH, test.exp.regH)
		}
		if env.regL != test.exp.regL {
			t.Errorf("[0x%02x] env.regL = %v, expected %v", opcode, env.regL, test.exp.regL)
		}
		if env.pc != test.exp.pc {
			t.Errorf("[0x%02x] env.pc = %v, expected %v", opcode, env.pc, test.exp.pc)
		}
		for i := 0; i < len(env.mem); i++ {
			if env.mem[i] != test.exp.mem[i] {
				t.Errorf("[0x%02x] env.mem[%d] = 0x%02x, expected 0x%02x", opcode, i, env.mem[i], test.exp.mem[i])
			}
		}
		if env.flags != test.exp.flags {
			t.Errorf("[0x%02x] env.flags = %.8b, expected %.8b", opcode, env.flags, test.exp.flags)
		}
	}
}
