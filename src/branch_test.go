package main

import "testing"

func TestJump(t *testing.T) {
	var table = []Pair{
		Pair{ // JNZ addr when Z is set
			State{pc: 0, mem: [65536]byte{0xc2, 0xff, 0xf0}, flags: 0b00000001},
			State{pc: 3, mem: [65536]byte{0xc2, 0xff, 0xf0}, flags: 0b00000001},
		},
		Pair{ // JNZ addr when Z is not set
			State{pc: 0, mem: [65536]byte{0xc3, 0xff, 0xf0}},
			State{pc: 0xfff0, mem: [65536]byte{0xc3, 0xff, 0xf0}},
		},
		Pair{ // JMP addr
			State{pc: 0, mem: [65536]byte{0xc3, 0xff, 0xf0}},
			State{pc: 0xfff0, mem: [65536]byte{0xc3, 0xff, 0xf0}},
		},
		Pair{ // JZ addr when Z is set
			State{pc: 0, mem: [65536]byte{0xca, 0xff, 0xf0}, flags: 0b00000001},
			State{pc: 0xfff0, mem: [65536]byte{0xca, 0xff, 0xf0}, flags: 0b00000001},
		},
		Pair{ // JZ addr when Z is not set
			State{pc: 0, mem: [65536]byte{0xca, 0xff, 0xf0}},
			State{pc: 3, mem: [65536]byte{0xca, 0xff, 0xf0}},
		},
		Pair{ // JNC addr when Cy is set
			State{pc: 0, mem: [65536]byte{0xd2, 0xff, 0xf0}, flags: 0b00001000},
			State{pc: 3, mem: [65536]byte{0xd2, 0xff, 0xf0}, flags: 0b00001000},
		},
		Pair{ // JNC addr when Cy is not set
			State{pc: 0, mem: [65536]byte{0xd2, 0xff, 0xf0}},
			State{pc: 0xfff0, mem: [65536]byte{0xd2, 0xff, 0xf0}},
		},
		Pair{ // JC addr when Cy is set
			State{pc: 0, mem: [65536]byte{0xda, 0xff, 0xf0}, flags: 0b00001000},
			State{pc: 0xfff0, mem: [65536]byte{0xda, 0xff, 0xf0}, flags: 0b00001000},
		},
		Pair{ // JC addr when Cy is not set
			State{pc: 0, mem: [65536]byte{0xda, 0xff, 0xf0}},
			State{pc: 3, mem: [65536]byte{0xda, 0xff, 0xf0}},
		},
		Pair{ // JPO addr when P is set
			State{pc: 0, mem: [65536]byte{0xe2, 0xff, 0xf0}, flags: 0b00000100},
			State{pc: 3, mem: [65536]byte{0xe2, 0xff, 0xf0}, flags: 0b00000100},
		},
		Pair{ // JPO addr when P is not set
			State{pc: 0, mem: [65536]byte{0xe2, 0xff, 0xf0}},
			State{pc: 0xfff0, mem: [65536]byte{0xe2, 0xff, 0xf0}},
		},
		Pair{ // JPE addr when P is set
			State{pc: 0, mem: [65536]byte{0xea, 0xff, 0xf0}, flags: 0b00000100},
			State{pc: 0xfff0, mem: [65536]byte{0xea, 0xff, 0xf0}, flags: 0b00000100},
		},
		Pair{ // JPE addr when P is not set
			State{pc: 0, mem: [65536]byte{0xea, 0xff, 0xf0}},
			State{pc: 3, mem: [65536]byte{0xea, 0xff, 0xf0}},
		},
		Pair{ // JP addr when S is set
			State{pc: 0, mem: [65536]byte{0xf2, 0xff, 0xf0}, flags: 0b00000010},
			State{pc: 3, mem: [65536]byte{0xf2, 0xff, 0xf0}, flags: 0b00000010},
		},
		Pair{ // JP addr when S is not set
			State{pc: 0, mem: [65536]byte{0xf2, 0xff, 0xf0}},
			State{pc: 0xfff0, mem: [65536]byte{0xf2, 0xff, 0xf0}},
		},
		Pair{ // JM addr when S is set
			State{pc: 0, mem: [65536]byte{0xfa, 0xff, 0xf0}, flags: 0b00000010},
			State{pc: 0xfff0, mem: [65536]byte{0xfa, 0xff, 0xf0}, flags: 0b00000010},
		},
		Pair{ // JM addr when S is not set
			State{pc: 0, mem: [65536]byte{0xfa, 0xff, 0xf0}},
			State{pc: 3, mem: [65536]byte{0xfa, 0xff, 0xf0}},
		},
	}
	doTest(t, table)
}

func TestCall(t *testing.T) {
	var table = []Pair{}
	doTest(t, table)
}

func doTest(t *testing.T, table []Pair) {
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
