package main

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

func (s *State) cmp(x byte) {
	var result uint16 = uint16(s.regA) + (^uint16(x) + 1)
	s.setFlags(result)
}
