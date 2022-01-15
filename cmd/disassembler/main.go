package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

type Opcode struct {
	Mnemonic    string
	Size        uint8
	FirstOp     Operand
	OperandLow  Operand
	OperandHigh Operand
}

type Operand int

const (
	Nil       Operand = iota // no operand
	RegA                     // Register A
	RegB                     // Register B
	RegC                     // Register C
	RegD                     // Register D
	RegE                     // Register E
	RegH                     // Register H
	RegL                     // Register L
	RegM                     // Not a register, but means (HL), and is treated as such
	RegSp                    // Stack Pointer
	Reg0                     // Register used for RST
	Reg1                     // Register used for RST
	Reg2                     // Register used for RST
	Reg3                     // Register used for RST
	Reg4                     // Register used for RST
	Reg5                     // Register used for RST
	Reg6                     // Register used for RST
	Reg7                     // Register used for RST
	RegPsw                   // Register A + Flags
	Addr                     // Address
	Immediate                // Immediate value
)

func (o Operand) IsRegister() bool {
	return o >= RegA && o <= RegPsw
}

var opcodes = map[byte]Opcode{
	0x00: Opcode{Mnemonic: "NOP ", Size: 1},
	0x01: Opcode{"LXI ", 3, RegB, Immediate, Immediate},
	0x02: Opcode{Mnemonic: "STAX", Size: 1, FirstOp: RegB},
	0x03: Opcode{Mnemonic: "INX ", Size: 1, FirstOp: RegB},
	0x04: Opcode{Mnemonic: "INR ", Size: 1, FirstOp: RegB},
	0x05: Opcode{Mnemonic: "DCR ", Size: 1, FirstOp: RegB},
	0x06: Opcode{"MVI ", 2, RegB, Immediate, Nil},
	0x07: Opcode{Mnemonic: "RLC ", Size: 1},

	0x09: Opcode{Mnemonic: "DAD ", Size: 1, FirstOp: RegB},
	0x0a: Opcode{Mnemonic: "LDAX", Size: 1, FirstOp: RegB},
	0x0b: Opcode{Mnemonic: "DCX ", Size: 1, FirstOp: RegB},
	0x0c: Opcode{Mnemonic: "INR ", Size: 1, FirstOp: RegC},
	0x0d: Opcode{Mnemonic: "DCR ", Size: 1, FirstOp: RegC},
	0x0e: Opcode{"MVI ", 2, RegC, Immediate, Nil},
	0x0f: Opcode{Mnemonic: "RRC ", Size: 1},

	0x11: Opcode{"LXI ", 3, RegD, Immediate, Immediate},
	0x12: Opcode{Mnemonic: "STAX", Size: 1, FirstOp: RegD},
	0x13: Opcode{Mnemonic: "INX", Size: 1, FirstOp: RegD},
	0x14: Opcode{Mnemonic: "INR", Size: 1, FirstOp: RegD},
	0x15: Opcode{Mnemonic: "DCR", Size: 1, FirstOp: RegD},
	0x16: Opcode{"MVI ", 2, RegD, Immediate, Nil},
	0x17: Opcode{Mnemonic: "RAL", Size: 1},

	0x19: Opcode{Mnemonic: "DAD ", Size: 1, FirstOp: RegD},
	0x1a: Opcode{Mnemonic: "LDAX", Size: 1, FirstOp: RegD},
	0x1b: Opcode{Mnemonic: "DCX ", Size: 1, FirstOp: RegD},
	0x1c: Opcode{Mnemonic: "INR ", Size: 1, FirstOp: RegD},
	0x1d: Opcode{Mnemonic: "DCR ", Size: 1, FirstOp: RegD},
	0x1e: Opcode{"MVI ", 2, RegE, Immediate, Nil},
	0x1f: Opcode{Mnemonic: "RAR ", Size: 1},

	0x21: Opcode{"LXI ", 3, RegH, Immediate, Immediate},
	0x22: Opcode{"SHLD", 3, Addr, Addr, Nil},
	0x23: Opcode{Mnemonic: "INX", Size: 1, FirstOp: RegH},
	0x24: Opcode{Mnemonic: "INR", Size: 1, FirstOp: RegH},
	0x25: Opcode{Mnemonic: "DCR", Size: 1, FirstOp: RegH},
	0x26: Opcode{"MVI ", 2, RegH, Immediate, Nil},
	0x27: Opcode{Mnemonic: "DAA", Size: 1},

	0x29: Opcode{Mnemonic: "DAD ", Size: 1, FirstOp: RegH},
	0x2a: Opcode{"LHLD", 3, Addr, Addr, Nil},
	0x2b: Opcode{Mnemonic: "DCX ", Size: 1, FirstOp: RegH},
	0x2c: Opcode{Mnemonic: "INR ", Size: 1, FirstOp: RegL},
	0x2d: Opcode{Mnemonic: "DCR ", Size: 1, FirstOp: RegL},
	0x2e: Opcode{"MVI ", 2, RegL, Immediate, Nil},
	0x2f: Opcode{Mnemonic: "CMA ", Size: 1},

	0x31: Opcode{"LXI ", 3, RegSp, Immediate, Immediate},
	0x32: Opcode{"STA ", 3, Addr, Addr, Nil},
	0x33: Opcode{Mnemonic: "INX", Size: 1, FirstOp: RegSp},
	0x34: Opcode{Mnemonic: "INR", Size: 1, FirstOp: RegM},
	0x35: Opcode{Mnemonic: "DCR", Size: 1, FirstOp: RegM},
	0x36: Opcode{"MVI ", 2, RegM, Immediate, Nil},
	0x37: Opcode{Mnemonic: "STC", Size: 1},

	0x39: Opcode{Mnemonic: "DAD ", Size: 1, FirstOp: RegSp},
	0x3a: Opcode{"LDA", 3, Addr, Addr, Nil},
	0x3b: Opcode{Mnemonic: "DCX ", Size: 1, FirstOp: RegSp},
	0x3c: Opcode{Mnemonic: "INR ", Size: 1, FirstOp: RegA},
	0x3d: Opcode{Mnemonic: "DCR ", Size: 1, FirstOp: RegA},
	0x3e: Opcode{"MVI ", 2, RegA, Immediate, Nil},
	0x3f: Opcode{Mnemonic: "CMC ", Size: 1},
	0x40: Opcode{"MOV ", 1, RegB, RegB, Nil},
	0x41: Opcode{"MOV ", 1, RegB, RegC, Nil},
	0x42: Opcode{"MOV ", 1, RegB, RegD, Nil},
	0x43: Opcode{"MOV ", 1, RegB, RegE, Nil},
	0x44: Opcode{"MOV ", 1, RegB, RegH, Nil},
	0x45: Opcode{"MOV ", 1, RegB, RegL, Nil},
	0x46: Opcode{"MOV ", 1, RegB, RegM, Nil},
	0x47: Opcode{"MOV ", 1, RegB, RegA, Nil},
	0x48: Opcode{"MOV ", 1, RegC, RegB, Nil},
	0x49: Opcode{"MOV ", 1, RegC, RegC, Nil},
	0x4a: Opcode{"MOV ", 1, RegC, RegD, Nil},
	0x4b: Opcode{"MOV ", 1, RegC, RegE, Nil},
	0x4c: Opcode{"MOV ", 1, RegC, RegH, Nil},
	0x4d: Opcode{"MOV ", 1, RegC, RegL, Nil},
	0x4e: Opcode{"MOV ", 1, RegC, RegM, Nil},
	0x4f: Opcode{"MOV ", 1, RegC, RegA, Nil},
	0x50: Opcode{"MOV ", 1, RegD, RegB, Nil},
	0x51: Opcode{"MOV ", 1, RegD, RegC, Nil},
	0x52: Opcode{"MOV ", 1, RegD, RegD, Nil},
	0x53: Opcode{"MOV ", 1, RegD, RegE, Nil},
	0x54: Opcode{"MOV ", 1, RegD, RegH, Nil},
	0x55: Opcode{"MOV ", 1, RegD, RegL, Nil},
	0x56: Opcode{"MOV ", 1, RegD, RegM, Nil},
	0x57: Opcode{"MOV ", 1, RegD, RegA, Nil},
	0x58: Opcode{"MOV ", 1, RegE, RegB, Nil},
	0x59: Opcode{"MOV ", 1, RegE, RegC, Nil},
	0x5a: Opcode{"MOV ", 1, RegE, RegD, Nil},
	0x5b: Opcode{"MOV ", 1, RegE, RegE, Nil},
	0x5c: Opcode{"MOV ", 1, RegE, RegH, Nil},
	0x5d: Opcode{"MOV ", 1, RegE, RegL, Nil},
	0x5e: Opcode{"MOV ", 1, RegE, RegM, Nil},
	0x5f: Opcode{"MOV ", 1, RegE, RegA, Nil},
	0x60: Opcode{"MOV ", 1, RegH, RegB, Nil},
	0x61: Opcode{"MOV ", 1, RegH, RegC, Nil},
	0x62: Opcode{"MOV ", 1, RegH, RegD, Nil},
	0x63: Opcode{"MOV ", 1, RegH, RegE, Nil},
	0x64: Opcode{"MOV ", 1, RegH, RegH, Nil},
	0x65: Opcode{"MOV ", 1, RegH, RegL, Nil},
	0x66: Opcode{"MOV ", 1, RegH, RegM, Nil},
	0x67: Opcode{"MOV ", 1, RegH, RegA, Nil},
	0x68: Opcode{"MOV ", 1, RegL, RegB, Nil},
	0x69: Opcode{"MOV ", 1, RegL, RegC, Nil},
	0x6a: Opcode{"MOV ", 1, RegL, RegD, Nil},
	0x6b: Opcode{"MOV ", 1, RegL, RegE, Nil},
	0x6c: Opcode{"MOV ", 1, RegL, RegH, Nil},
	0x6d: Opcode{"MOV ", 1, RegL, RegL, Nil},
	0x6e: Opcode{"MOV ", 1, RegL, RegM, Nil},
	0x6f: Opcode{"MOV ", 1, RegL, RegA, Nil},
	0x70: Opcode{"MOV ", 1, RegM, RegB, Nil},
	0x71: Opcode{"MOV ", 1, RegM, RegC, Nil},
	0x72: Opcode{"MOV ", 1, RegM, RegD, Nil},
	0x73: Opcode{"MOV ", 1, RegM, RegE, Nil},
	0x74: Opcode{"MOV ", 1, RegM, RegH, Nil},
	0x75: Opcode{"MOV ", 1, RegM, RegL, Nil},
	0x76: Opcode{"HLT ", 1, Nil, Nil, Nil},
	0x77: Opcode{"MOV ", 1, RegM, RegA, Nil},
	0x78: Opcode{"MOV ", 1, RegA, RegB, Nil},
	0x79: Opcode{"MOV ", 1, RegA, RegC, Nil},
	0x7a: Opcode{"MOV ", 1, RegA, RegD, Nil},
	0x7b: Opcode{"MOV ", 1, RegA, RegE, Nil},
	0x7c: Opcode{"MOV ", 1, RegA, RegH, Nil},
	0x7d: Opcode{"MOV ", 1, RegA, RegL, Nil},
	0x7e: Opcode{"MOV ", 1, RegA, RegM, Nil},
	0x7f: Opcode{"MOV ", 1, RegA, RegA, Nil},
	0x80: Opcode{"ADD ", 1, RegB, Nil, Nil},
	0x81: Opcode{"ADD ", 1, RegC, Nil, Nil},
	0x82: Opcode{"ADD ", 1, RegD, Nil, Nil},
	0x83: Opcode{"ADD ", 1, RegE, Nil, Nil},
	0x84: Opcode{"ADD ", 1, RegH, Nil, Nil},
	0x85: Opcode{"ADD ", 1, RegL, Nil, Nil},
	0x86: Opcode{"ADD ", 1, RegM, Nil, Nil},
	0x87: Opcode{"ADD ", 1, RegA, Nil, Nil},
	0x88: Opcode{"ADC ", 1, RegB, Nil, Nil},
	0x89: Opcode{"ADC ", 1, RegC, Nil, Nil},
	0x8a: Opcode{"ADC ", 1, RegD, Nil, Nil},
	0x8b: Opcode{"ADC ", 1, RegE, Nil, Nil},
	0x8c: Opcode{"ADC ", 1, RegH, Nil, Nil},
	0x8d: Opcode{"ADC ", 1, RegL, Nil, Nil},
	0x8e: Opcode{"ADC ", 1, RegM, Nil, Nil},
	0x8f: Opcode{"ADC ", 1, RegA, Nil, Nil},
	0x90: Opcode{"SUB ", 1, RegB, Nil, Nil},
	0x91: Opcode{"SUB ", 1, RegC, Nil, Nil},
	0x92: Opcode{"SUB ", 1, RegD, Nil, Nil},
	0x93: Opcode{"SUB ", 1, RegE, Nil, Nil},
	0x94: Opcode{"SUB ", 1, RegH, Nil, Nil},
	0x95: Opcode{"SUB ", 1, RegL, Nil, Nil},
	0x96: Opcode{"SUB ", 1, RegM, Nil, Nil},
	0x97: Opcode{"SUB ", 1, RegA, Nil, Nil},
	0x98: Opcode{"SBB ", 1, RegB, Nil, Nil},
	0x99: Opcode{"SBB ", 1, RegC, Nil, Nil},
	0x9a: Opcode{"SBB ", 1, RegD, Nil, Nil},
	0x9b: Opcode{"SBB ", 1, RegE, Nil, Nil},
	0x9c: Opcode{"SBB ", 1, RegH, Nil, Nil},
	0x9d: Opcode{"SBB ", 1, RegL, Nil, Nil},
	0x9e: Opcode{"SBB ", 1, RegM, Nil, Nil},
	0x9f: Opcode{"SBB ", 1, RegA, Nil, Nil},
	0xa0: Opcode{"ANA ", 1, RegB, Nil, Nil},
	0xa1: Opcode{"ANA ", 1, RegC, Nil, Nil},
	0xa2: Opcode{"ANA ", 1, RegD, Nil, Nil},
	0xa3: Opcode{"ANA ", 1, RegE, Nil, Nil},
	0xa4: Opcode{"ANA ", 1, RegH, Nil, Nil},
	0xa5: Opcode{"ANA ", 1, RegL, Nil, Nil},
	0xa6: Opcode{"ANA ", 1, RegM, Nil, Nil},
	0xa7: Opcode{"ANA ", 1, RegA, Nil, Nil},
	0xa8: Opcode{"XRA ", 1, RegB, Nil, Nil},
	0xa9: Opcode{"XRA ", 1, RegC, Nil, Nil},
	0xaa: Opcode{"XRA ", 1, RegD, Nil, Nil},
	0xab: Opcode{"XRA ", 1, RegE, Nil, Nil},
	0xac: Opcode{"XRA ", 1, RegH, Nil, Nil},
	0xad: Opcode{"XRA ", 1, RegL, Nil, Nil},
	0xae: Opcode{"XRA ", 1, RegM, Nil, Nil},
	0xaf: Opcode{"XRA ", 1, RegA, Nil, Nil},
	0xb0: Opcode{"ORA ", 1, RegB, Nil, Nil},
	0xb1: Opcode{"ORA ", 1, RegC, Nil, Nil},
	0xb2: Opcode{"ORA ", 1, RegD, Nil, Nil},
	0xb3: Opcode{"ORA ", 1, RegE, Nil, Nil},
	0xb4: Opcode{"ORA ", 1, RegH, Nil, Nil},
	0xb5: Opcode{"ORA ", 1, RegL, Nil, Nil},
	0xb6: Opcode{"ORA ", 1, RegM, Nil, Nil},
	0xb7: Opcode{"ORA ", 1, RegA, Nil, Nil},
	0xb8: Opcode{"CMP ", 1, RegB, Nil, Nil},
	0xb9: Opcode{"CMP ", 1, RegC, Nil, Nil},
	0xba: Opcode{"CMP ", 1, RegD, Nil, Nil},
	0xbb: Opcode{"CMP ", 1, RegE, Nil, Nil},
	0xbc: Opcode{"CMP ", 1, RegH, Nil, Nil},
	0xbd: Opcode{"CMP ", 1, RegL, Nil, Nil},
	0xbe: Opcode{"CMP ", 1, RegM, Nil, Nil},
	0xbf: Opcode{"CMP ", 1, RegA, Nil, Nil},
	0xc0: Opcode{"RNZ ", 1, Nil, Nil, Nil},
	0xc1: Opcode{"POP ", 1, RegB, Nil, Nil},
	0xc2: Opcode{"JNZ ", 3, Addr, Addr, Nil},
	0xc3: Opcode{"JMP ", 3, Addr, Addr, Nil},
	0xc4: Opcode{"CNZ ", 3, Addr, Addr, Nil},
	0xc5: Opcode{"PUSH", 1, RegB, Nil, Nil},
	0xc6: Opcode{"ADI ", 2, Immediate, Nil, Nil},
	0xc7: Opcode{"RST ", 1, Reg0, Nil, Nil},
	0xc8: Opcode{"RZ  ", 1, Nil, Nil, Nil},
	0xc9: Opcode{"RET ", 1, Nil, Nil, Nil},
	0xca: Opcode{"JZ  ", 3, Addr, Addr, Nil},

	0xcc: Opcode{"CZ  ", 3, Addr, Addr, Nil},
	0xcd: Opcode{"CALL", 3, Addr, Addr, Nil},
	0xce: Opcode{"ACI ", 2, Immediate, Nil, Nil},
	0xcf: Opcode{"RST ", 1, Reg1, Nil, Nil},
	0xd0: Opcode{"RNC ", 1, Reg1, Nil, Nil},
	0xd1: Opcode{"POP ", 1, RegD, Nil, Nil},
	0xd2: Opcode{"JNC ", 3, Addr, Addr, Nil},
	0xd3: Opcode{"OUT ", 2, Immediate, Nil, Nil},
	0xd4: Opcode{"CNC ", 3, Addr, Addr, Nil},
	0xd5: Opcode{"PUSH", 1, RegD, Nil, Nil},
	0xd6: Opcode{"SUI ", 2, Immediate, Nil, Nil},
	0xd7: Opcode{"RST ", 1, Reg2, Nil, Nil},
	0xd8: Opcode{"RC  ", 1, Nil, Nil, Nil},

	0xda: Opcode{"JC  ", 3, Addr, Addr, Nil},
	0xdb: Opcode{"IN  ", 2, Immediate, Nil, Nil},
	0xdc: Opcode{"CC  ", 3, Addr, Addr, Nil},

	0xde: Opcode{"SBI ", 2, Immediate, Nil, Nil},
	0xdf: Opcode{"RST ", 1, Reg3, Nil, Nil},
	0xe0: Opcode{"RPO ", 1, Nil, Nil, Nil},
	0xe1: Opcode{"POP ", 1, RegH, Nil, Nil},
	0xe2: Opcode{"JPO ", 3, Addr, Addr, Nil},
	0xe3: Opcode{"XTHL", 1, Nil, Nil, Nil},
	0xe4: Opcode{"CPO ", 3, Addr, Addr, Nil},
	0xe5: Opcode{"PUSH", 1, RegH, Nil, Nil},
	0xe6: Opcode{"ANI ", 2, Immediate, Nil, Nil},
	0xe7: Opcode{"RST ", 1, Reg4, Nil, Nil},
	0xe8: Opcode{"RPE ", 1, Nil, Nil, Nil},
	0xe9: Opcode{"PCHL", 1, Nil, Nil, Nil},
	0xea: Opcode{"JPE ", 3, Addr, Addr, Nil},
	0xeb: Opcode{"XCHG", 1, Nil, Nil, Nil},
	0xec: Opcode{"CPE ", 3, Addr, Addr, Nil},

	0xee: Opcode{"XRI ", 2, Immediate, Nil, Nil},
	0xef: Opcode{"RST ", 1, Reg5, Nil, Nil},
	0xf0: Opcode{"RP  ", 1, Reg5, Nil, Nil},
	0xf1: Opcode{"POP ", 1, RegPsw, Nil, Nil},
	0xf2: Opcode{"JP  ", 3, Addr, Addr, Nil},
	0xf3: Opcode{"DI  ", 1, Nil, Nil, Nil},
	0xf4: Opcode{"CP  ", 3, Addr, Addr, Nil},
	0xf5: Opcode{"PUSH", 1, RegPsw, Nil, Nil},
	0xf6: Opcode{"ORI ", 2, Immediate, Nil, Nil},
	0xf7: Opcode{"RST ", 1, Reg6, Nil, Nil},
	0xf8: Opcode{"RM  ", 1, Nil, Nil, Nil},
	0xf9: Opcode{"SPHL", 1, Nil, Nil, Nil},
	0xfa: Opcode{"JM  ", 3, Addr, Addr, Nil},
	0xfb: Opcode{"EI  ", 1, Nil, Nil, Nil},
	0xfc: Opcode{"CM  ", 3, Addr, Addr, Nil},

	0xfe: Opcode{"CPI ", 2, Immediate, Nil, Nil},
	0xff: Opcode{"RST ", 1, Reg7, Nil, Nil},
}

var registers = map[Operand]string{
	RegA:   "A",
	RegB:   "B",
	RegC:   "C",
	RegD:   "D",
	RegE:   "E",
	RegH:   "H",
	RegL:   "L",
	RegM:   "M",
	RegSp:  "SP",
	Reg0:   "0",
	Reg1:   "1",
	Reg2:   "2",
	Reg3:   "3",
	Reg4:   "4",
	Reg5:   "5",
	Reg6:   "6",
	Reg7:   "7",
	RegPsw: "PSW",
}

func main() {
	var f = os.Stdin
	var buf, err = io.ReadAll(f)

	if err != nil {
		log.Fatalf("Error reading stdin: %v", err)
	}

	var pc int
	var sz = len(buf)
	for pc < sz {
		fmt.Printf("%04x ", pc)

		var instr, opsz = disassemble(pc, buf)
		fmt.Println(instr)

		pc += int(opsz)
	}
}

func disassemble(pc int, buf []byte) (string, uint8) {
	var opcode, ok = opcodes[buf[pc]]
	if !ok {
		log.Fatalf("Not recognized operation code: %x", buf[pc])
	}

	switch opcode.Size {
	case 1:
		return disassembleSize1(opcode, buf[pc]), opcode.Size
	case 2:
		return disassembleSize2(opcode, buf[pc], buf[pc+1]), opcode.Size
	case 3:
		return disassembleSize3(opcode, buf[pc], buf[pc+1], buf[pc+2]), opcode.Size
	}

	panic(fmt.Sprintf("Bad size on opcode %x", buf[pc]))
}

func disassembleSize1(opcode Opcode, instr byte) string {
	var header = fmt.Sprintf("%02x       %s", instr, opcode.Mnemonic)
	if opcode.FirstOp.IsRegister() && opcode.OperandLow.IsRegister() {
		return fmt.Sprintf("%s   %s, %s", header, registers[opcode.FirstOp], registers[opcode.OperandLow])
	} else if opcode.FirstOp.IsRegister() {
		return fmt.Sprintf("%s   %s", header, registers[opcode.FirstOp])
	} else {
		return fmt.Sprintf("%02x       %s", instr, opcode.Mnemonic)
	}
}

func disassembleSize2(opcode Opcode, instr, operand byte) string {
	var header = fmt.Sprintf("%02x %02x    %s", instr, operand, opcode.Mnemonic)
	if opcode.FirstOp.IsRegister() {
		if opcode.OperandLow != Immediate {
			log.Fatalf("disassembleSize2: the operand must be an immediate value")
		}

		return fmt.Sprintf("%s   %s, #0x%02x", header, registers[opcode.FirstOp], operand)
	} else if opcode.FirstOp == Immediate {
		return fmt.Sprintf("%s   #0x%02x", header, operand)
	}

	panic(fmt.Sprintf("disassembleSize2: unknown operation: %v", opcode))
}

func disassembleSize3(opcode Opcode, instr, low, high byte) string {
	var header = fmt.Sprintf("%02x %02x %02x %s", instr, low, high, opcode.Mnemonic)
	if opcode.FirstOp.IsRegister() && opcode.OperandLow == Immediate && opcode.OperandHigh == Immediate {
		return fmt.Sprintf("%s   %s, #0x%02x%02x", header, registers[opcode.FirstOp], low, high)
	} else if opcode.FirstOp == Addr && opcode.OperandLow == Addr {
		return fmt.Sprintf("%s   $%02x%02x", header, high, low)
	}

	panic(fmt.Sprintf("disassembleSize3: unknown operation: %v", opcode))
}
