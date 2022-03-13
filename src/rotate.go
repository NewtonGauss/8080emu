package main

func (s *State) rlc() {
}

func (s *State) rrc() {
	var x byte = s.regA
	var bit0 byte = x & 1
	s.regA = (bit0 << 7) | (x >> 1)
	s.flags.SetValue(FlagCy, 1 == (x&1))
}

func (s *State) ral() { panic("Unimplemented") }

func (s *State) rar() { panic("Unimplemented") }
