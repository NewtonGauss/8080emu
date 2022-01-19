package main

func (s *State) add(x byte) {
	// I do the math with 8 bits more, to capture the carry
	var result uint16 = uint16(s.regA) + uint16(x)
	s.setFlags(result)
	s.regA = byte(result)
}

func (s *State) addCy(x byte) {
	var cy byte
	if s.flags.IsSet(FlagCy) {
		cy = 1
	}
	s.add(x + cy)
}

// dad performs the DAD instruction
// DAD B, DAD D, DAD H, DAD SP
// Add x to HL and just sets Carry, but no other flag.
func (s *State) dad(x uint16) {
	var result uint32 = uint32(s.hl()) + uint32(x)
	s.flags.SetValue(FlagCy, result > 0xffff)
	s.regH = byte(result >> 8)
	s.regL = byte(result)
}

func (s *State) inc(x *byte) {
	*x = *x + 1
	s.setFlagsNoCy(*x)
}

func (s *State) inx(x, y *byte) {
	var result = pairTo16(*x, *y) + 1
	*x = byte(result >> 8)
	*y = byte(result)
}

func (s *State) dec(x *byte) {
	*x = *x - 1
	s.setFlagsNoCy(*x)
}

func (s *State) dcx(x, y *byte) {
	var result = pairTo16(*x, *y) - 1
	*x = byte(result >> 8)
	*y = byte(result)
}

func (s *State) sub(x byte) {
	var result uint16 = uint16(s.regA) - uint16(x)
	s.setFlags(result)
	s.regA = byte(result)
}

func (s *State) subCy(x byte) {
	var cy byte
	if s.flags.IsSet(FlagCy) {
		cy = 1
	}
	// a + because inside `sub`, it will be regA - (x + cy) = regA - x - cy
	s.sub(x + cy)
}

func (s *State) setFlags(result uint16) {
	s.setFlagsNoCy(byte(result))
	// Carry flag
	s.flags.SetValue(FlagCy, result > 0xff)
}

func (s *State) setFlagsNoCy(result byte) {
	// Zero flag
	s.flags.SetValue(FlagZ, (result&0xff) == 0)

	// Sign flag
	s.flags.SetValue(FlagS, (result&0b10000000) != 0)

	// Parity flag
	s.flags.SetValue(FlagP, isParityEven(result))
}

func isParityEven(x byte) bool {
	var p int
	for x != 0 {
		x &= x - 1
		p++
	}
	return p%2 == 0
}
