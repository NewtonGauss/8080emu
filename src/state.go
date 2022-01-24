package main

import "fmt"

type Flags byte

type Flag byte

// do not change the order of flags: Z, S, P, Cy, Ac
const (
	FlagZ  Flag = iota // Zero flag: set to 1 when the result is zero
	FlagS              // Sign flag: set to 1 when bit 7 is set
	FlagP              // Parity flag: set to 1 when the result has even parity
	FlagCy             // Carry flag: set to 1 when the result has a carry
	FlagAc             // Aux carry flag: used for BCD math (Unimplemented)
)

func (f *Flags) SetValue(bit Flag, set bool) {
	if set {
		*f |= Flags(1 << bit)
	} else {
		*f &= ^Flags(1 << bit)
	}
}

func (f *Flags) Set(bit Flag) {
	f.SetValue(bit, true)
}

func (f *Flags) Unset(bit Flag) {
	f.SetValue(bit, false)
}

func (f *Flags) IsSet(bit Flag) bool {
	return *f&Flags(1<<bit) != 0
}

type State struct {
	regA      byte
	regB      byte
	regC      byte
	regD      byte
	regE      byte
	regH      byte
	regL      byte
	sp        uint16
	pc        uint16
	mem       [65536]byte // 32 K = 2^16 = 65536
	flags     Flags
	intEnable byte
}

var instrSz = map[byte]byte{
	0x00: 1,
	0x01: 3,
	0x02: 1,
	0x03: 1,
	0x04: 1,
	0x05: 1,
	0x06: 2,
	0x07: 1,
	0x08: 1,
	0x09: 1,
	0x0a: 1,
	0x0b: 1,
	0x0c: 1,
	0x0d: 1,
	0x0e: 2,
	0x0f: 1,
	0x10: 1,
	0x11: 3,
	0x12: 1,
	0x13: 1,
	0x14: 1,
	0x15: 1,
	0x16: 2,
	0x17: 1,
	0x18: 1,
	0x19: 1,
	0x1a: 1,
	0x1b: 1,
	0x1c: 1,
	0x1d: 1,
	0x1e: 2,
	0x1f: 1,
	0x20: 1,
	0x21: 3,
	0x22: 3,
	0x23: 1,
	0x24: 1,
	0x25: 1,
	0x26: 2,
	0x27: 1,
	0x28: 1,
	0x29: 1,
	0x2a: 3,
	0x2b: 1,
	0x2c: 1,
	0x2d: 1,
	0x2e: 2,
	0x2f: 1,
	0x30: 1,
	0x31: 3,
	0x32: 3,
	0x33: 1,
	0x34: 1,
	0x35: 1,
	0x36: 2,
	0x37: 1,
	0x38: 1,
	0x39: 1,
	0x3a: 3,
	0x3b: 1,
	0x3c: 1,
	0x3d: 1,
	0x3e: 2,
	0x3f: 1,
	0x40: 1,
	0x41: 1,
	0x42: 1,
	0x43: 1,
	0x44: 1,
	0x45: 1,
	0x46: 1,
	0x47: 1,
	0x48: 1,
	0x49: 1,
	0x4a: 1,
	0x4b: 1,
	0x4c: 1,
	0x4d: 1,
	0x4e: 1,
	0x4f: 1,
	0x50: 1,
	0x51: 1,
	0x52: 1,
	0x53: 1,
	0x54: 1,
	0x55: 1,
	0x56: 1,
	0x57: 1,
	0x58: 1,
	0x59: 1,
	0x5a: 1,
	0x5b: 1,
	0x5c: 1,
	0x5d: 1,
	0x5e: 1,
	0x5f: 1,
	0x60: 1,
	0x61: 1,
	0x62: 1,
	0x63: 1,
	0x64: 1,
	0x65: 1,
	0x66: 1,
	0x67: 1,
	0x68: 1,
	0x69: 1,
	0x6a: 1,
	0x6b: 1,
	0x6c: 1,
	0x6d: 1,
	0x6e: 1,
	0x6f: 1,
	0x70: 1,
	0x71: 1,
	0x72: 1,
	0x73: 1,
	0x74: 1,
	0x75: 1,
	0x76: 1,
	0x77: 1,
	0x78: 1,
	0x79: 1,
	0x7a: 1,
	0x7b: 1,
	0x7c: 1,
	0x7d: 1,
	0x7e: 1,
	0x7f: 1,
	0x80: 1,
	0x81: 1,
	0x82: 1,
	0x83: 1,
	0x84: 1,
	0x85: 1,
	0x86: 1,
	0x87: 1,
	0x88: 1,
	0x89: 1,
	0x8a: 1,
	0x8b: 1,
	0x8c: 1,
	0x8d: 1,
	0x8e: 1,
	0x8f: 1,
	0x90: 1,
	0x91: 1,
	0x92: 1,
	0x93: 1,
	0x94: 1,
	0x95: 1,
	0x96: 1,
	0x97: 1,
	0x98: 1,
	0x99: 1,
	0x9a: 1,
	0x9b: 1,
	0x9c: 1,
	0x9d: 1,
	0x9e: 1,
	0x9f: 1,
	0xa0: 1,
	0xa1: 1,
	0xa2: 1,
	0xa3: 1,
	0xa4: 1,
	0xa5: 1,
	0xa6: 1,
	0xa7: 1,
	0xa8: 1,
	0xa9: 1,
	0xaa: 1,
	0xab: 1,
	0xac: 1,
	0xad: 1,
	0xae: 1,
	0xaf: 1,
	0xb0: 1,
	0xb1: 1,
	0xb2: 1,
	0xb3: 1,
	0xb4: 1,
	0xb5: 1,
	0xb6: 1,
	0xb7: 1,
	0xb8: 1,
	0xb9: 1,
	0xba: 1,
	0xbb: 1,
	0xbc: 1,
	0xbd: 1,
	0xbe: 1,
	0xbf: 1,
	0xc0: 0,
	0xc1: 1,
	0xc2: 0, // JNZ ... I'll advance the pc manually in case Zero is set
	0xc3: 0, // JMP do not advance pc
	0xc4: 0,
	0xc5: 1,
	0xc6: 2,
	0xc7: 0,
	0xc8: 0,
	0xc9: 0, // RET do not advance pc
	0xca: 0, // JZ ... I'll advance the pc manually in case Zero is not set

	0xcc: 0, // CZ ... I'll advance the pc manually
	0xcd: 0, // CALL do not advance pc
	0xce: 2,
	0xcf: 0,
	0xd0: 0,
	0xd1: 1,
	0xd2: 0, // JNC ... I'll advance the pc manually in case Cy is set
	0xd3: 2,
	0xd4: 0,
	0xd5: 1,
	0xd6: 2,
	0xd7: 0,
	0xd8: 0,

	0xda: 0, // JC ... I'll advance the pc manually in case Cy is not set
	0xdb: 2,
	0xdc: 0,

	0xde: 2,
	0xdf: 0,
	0xe0: 0,
	0xe1: 1,
	0xe2: 0, // JPO ... I'll advance the pc manually in case Parity is even
	0xe3: 1,
	0xe4: 0,
	0xe5: 1,
	0xe6: 2,
	0xe7: 0,
	0xe8: 0,
	0xe9: 0,
	0xea: 0, // JPE ... I'll advance the pc manually in case Parity is odd
	0xeb: 1,
	0xec: 0,

	0xee: 2,
	0xef: 0,
	0xf0: 0,
	0xf1: 1,
	0xf2: 0, // JP ... I'll advance the pc manually in case the result is negative
	0xf3: 1,
	0xf4: 0,
	0xf5: 1,
	0xf6: 2,
	0xf7: 0,
	0xf8: 0,
	0xf9: 1,
	0xfa: 0, // JM ... I'll advance the pc manually in case the result is positive
	0xfb: 1,
	0xfc: 0,

	0xfe: 2,
	0xff: 0,
}

func (s *State) ExecInstruction() {
	var opcode = s.mem[s.pc]
	switch opcode {
	case 0x00:
		/* NOP */
	case 0x01: // LXI B, D16
		s.regC, s.regB = s.mem[s.pc+1], s.mem[s.pc+2]

	case 0x03: // INX B
		s.inx(&s.regB, &s.regC)
	case 0x04: // INR B
		s.inc(&s.regB)
	case 0x05: // DCR B
		s.dec(&s.regB)

	case 0x09: // DAD B
		s.dad(pairTo16(s.regB, s.regC))

	case 0x0b: // DCX B
		s.dcx(&s.regB, &s.regC)
	case 0x0C: // INR C
		s.inc(&s.regC)
	case 0x0D: // DCR C
		s.dec(&s.regC)

	case 0x13: // INX D
		s.inx(&s.regD, &s.regE)
	case 0x14: // INR D
		s.inc(&s.regD)
	case 0x15: // DCR D
		s.dec(&s.regD)

	case 0x19: // DAD D
		s.dad(pairTo16(s.regD, s.regE))

	case 0x1b: // DCX D
		s.dcx(&s.regD, &s.regE)
	case 0x1c: // INR E
		s.inc(&s.regE)
	case 0x1d: // DCR E
		s.dec(&s.regE)

	case 0x23: // INX H
		s.inx(&s.regH, &s.regL)
	case 0x24: // INR H
		s.inc(&s.regH)
	case 0x25: // DCR H
		s.dec(&s.regH)

	case 0x29: // DAD H
		s.dad(pairTo16(s.regH, s.regL))

	case 0x2b: // DCX H
		s.dcx(&s.regH, &s.regL)
	case 0x2c: // INR L
		s.inc(&s.regL)
	case 0x2d: // DCR L
		s.dec(&s.regL)

	case 0x33: // INX SP
		s.sp++
	case 0x34: // INR M
		s.inc(&s.mem[s.hl()])
	case 0x35: // DCR M
		s.dec(&s.mem[s.hl()])

	case 0x39: // DAD SP
		s.dad(s.sp)

	case 0x3b: // DCX SP
		s.sp--
	case 0x3c: // INR A
		s.inc(&s.regA)
	case 0x3d: // DCR A
		s.dec(&s.regA)

	case 0x41: // MOV B, C
		s.regB = s.regC
	case 0x42: // MOV B, D
		s.regB = s.regD
	case 0x43: // MOV B, E
		s.regB = s.regE

	case 0x80: // ADD B
		s.add(s.regB)
	case 0x81: // ADD C
		s.add(s.regC)
	case 0x82: // ADD D
		s.add(s.regD)
	case 0x83: // ADD E
		s.add(s.regE)
	case 0x84: // ADD H
		s.add(s.regH)
	case 0x85: // ADD L
		s.add(s.regL)
	case 0x86: // ADD M
		s.add(s.mem[s.hl()])
	case 0x87: // ADD A
		s.add(s.regA)
	case 0x88: // ADC B
		s.addCy(s.regB)
	case 0x89: // ADC C
		s.addCy(s.regC)
	case 0x8a: // ADC D
		s.addCy(s.regD)
	case 0x8b: // ADC E
		s.addCy(s.regE)
	case 0x8c: // ADC H
		s.addCy(s.regH)
	case 0x8d: // ADC L
		s.addCy(s.regL)
	case 0x8e: // ADC M
		s.addCy(s.mem[s.hl()])
	case 0x8f: // ADC A
		s.addCy(s.regA)
	case 0x90: // SUB B
		s.sub(s.regB)
	case 0x91: // SUB C
		s.sub(s.regC)
	case 0x92: // SUB D
		s.sub(s.regD)
	case 0x93: // SUB E
		s.sub(s.regE)
	case 0x94: // SUB H
		s.sub(s.regH)
	case 0x95: // SUB L
		s.sub(s.regL)
	case 0x96: // SUB M
		s.sub(s.mem[s.hl()])
	case 0x97: // SUB A
		s.sub(s.regA)
	case 0x98: // SBB B
		s.subCy(s.regB)
	case 0x99: // SBB C
		s.subCy(s.regC)
	case 0x9a: // SBB D
		s.subCy(s.regD)
	case 0x9b: // SBB E
		s.subCy(s.regE)
	case 0x9c: // SBB H
		s.subCy(s.regH)
	case 0x9d: // SBB L
		s.subCy(s.regL)
	case 0x9e: // SBB M
		s.subCy(s.mem[s.hl()])
	case 0x9f: // SBB A
		s.subCy(s.regA)
	case 0xa0: // ANA B
		s.and(s.regB)
	case 0xa1: // ANA C
		s.and(s.regC)
	case 0xa2: // ANA D
		s.and(s.regD)
	case 0xa3: // ANA E
		s.and(s.regE)
	case 0xa4: // ANA H
		s.and(s.regH)
	case 0xa5: // ANA L
		s.and(s.regL)
	case 0xa6: // ANA M
		s.and(s.mem[s.hl()])
	case 0xa7: // ANA A
		s.and(s.regA)
	case 0xa8: // XRA B
		s.xor(s.regB)
	case 0xa9: // XRA C
		s.xor(s.regC)
	case 0xaa: // XRA D
		s.xor(s.regD)
	case 0xab: // XRA E
		s.xor(s.regE)
	case 0xac: // XRA H
		s.xor(s.regH)
	case 0xad: // XRA L
		s.xor(s.regL)
	case 0xae: // XRA M
		s.xor(s.mem[s.hl()])
	case 0xaf: // XRA A
		s.xor(s.regA)
	case 0xb0: // ORA B
		s.or(s.regB)
	case 0xb1: // ORA C
		s.or(s.regC)
	case 0xb2: // ORA D
		s.or(s.regD)
	case 0xb3: // ORA E
		s.or(s.regE)
	case 0xb4: // ORA H
		s.or(s.regH)
	case 0xb5: // ORA L
		s.or(s.regL)
	case 0xb6: // ORA M
		s.or(s.mem[s.hl()])
	case 0xb7: // ORA A
		s.or(s.regA)

	case 0xc0: // RNZ
		s.retOnFlag(FlagZ, false)

	case 0xc2: // JNZ addr
		s.jmpOnFlag(FlagZ, false)
	case 0xc3: // JMP addr
		s.jmp()
	case 0xc4: // CNZ addr
		s.callOnFlag(FlagZ, false)
	case 0xc6: // ADI D8
		s.add(s.mem[s.pc+1])
	case 0xc7: // RST 0
		s.rst(0x00)
	case 0xc8: // RZ
		s.retOnFlag(FlagZ, true)
	case 0xc9: // RET
		s.ret()
	case 0xca: // JZ addr
		s.jmpOnFlag(FlagZ, true)

	case 0xcc: // CZ addr
		s.callOnFlag(FlagZ, true)
	case 0xcd: // CALL addr
		s.call()
	case 0xce: // ACI D8
		s.addCy(s.mem[s.pc+1])
	case 0xcf: // RST 1
		s.rst(0x08)
	case 0xd0: // RNC
		s.retOnFlag(FlagCy, false)

	case 0xd2: // JNC addr
		s.jmpOnFlag(FlagCy, false)

	case 0xd4: // CNC addr
		s.callOnFlag(FlagCy, false)

	case 0xd6: // SUI D8
		s.sub(s.mem[s.pc+1])
	case 0xd7: // RST 2
		s.rst(0x10)
	case 0xd8: // RC
		s.retOnFlag(FlagCy, true)

	case 0xda: // JC addr
		s.jmpOnFlag(FlagCy, true)

	case 0xdc: // CC addr
		s.callOnFlag(FlagCy, true)

	case 0xde: // SBI D8
		s.subCy(s.mem[s.pc+1])
	case 0xdf: // RST 3
		s.rst(0x18)
	case 0xe0: // RPO addr
		s.retOnFlag(FlagP, false)

	case 0xe2: // JPO addr
		s.jmpOnFlag(FlagP, false)

	case 0xe4: // CPO addr
		s.callOnFlag(FlagP, false)

	case 0xe7: // RST 4
		s.rst(0x20)

	case 0xe8: // RPE
		s.retOnFlag(FlagP, true)

	case 0xe9: // PCHL
		s.pc = s.hl()
	case 0xea: // JPE addr
		s.jmpOnFlag(FlagP, true)

	case 0xec: // CPE addr
		s.callOnFlag(FlagP, true)

	case 0xef: // RST 5
		s.rst(0x28)
	case 0xf0: // JP addr
		s.retOnFlag(FlagS, false)

	case 0xf2: // JP addr
		s.jmpOnFlag(FlagS, false)

	case 0xf4: // CP addr
		s.callOnFlag(FlagS, false)

	case 0xf7: // RST 6
		s.rst(0x30)
	case 0xf8: // RM
		s.retOnFlag(FlagS, true)

	case 0xfa: // JM addr ... Jump minus = Jump if negative = Jump when the Sign is set
		s.jmpOnFlag(FlagS, true)

	case 0xfc: // CM addr
		s.callOnFlag(FlagS, true)

	case 0xff: // RST 7
		s.rst(0x38)
	default:
		panic(fmt.Sprintf("unimplemented instruction: 0x%02x", opcode))
	}

	s.pc += uint16(instrSz[opcode])
}

// hl returns the value (usually used as an address) formed by hl
func (s *State) hl() uint16 {
	return pairTo16(s.regH, s.regL)
}

// pairTo16 returns a uint16 formed by (xy)
func pairTo16(x, y byte) uint16 {
	return (uint16(x) << 8) | uint16(y)
}

func (s *State) and(x byte) {
	s.regA &= x
	s.setFlagsNoCy(s.regA)
	s.flags.Unset(FlagCy) // AND unsets carry
}

func (s *State) xor(x byte) {
	s.regA ^= x
	s.setFlagsNoCy(s.regA)
	s.flags.Unset(FlagCy)
}

func (s *State) or(x byte) {
	s.regA |= x
	s.setFlagsNoCy(s.regA)
	s.flags.Unset(FlagCy)
}
