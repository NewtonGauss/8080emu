package main

// rlc rotates A left
// |7|6|5|4|3|2|1|0| and Cy|c| -> |6|5|4|3|2|1|0|7| and Cy|b7|
func (s *State) rlc() {
	var x byte = s.regA
	var bit7 byte = (x & 0b1000_0000) >> 7
	s.regA = (bit7) | (x << 1)
	s.flags.SetValue(FlagCy, 1 == bit7)
}

// rlc rotates A left through Carry
// |7|6|5|4|3|2|1|0| and Cy|c| -> |6|5|4|3|2|1|0|c| and Cy|b7|
func (s *State) ral() {
	var x byte = s.regA

	var cy byte = 0
	if s.flags.IsSet(FlagCy) {
		cy = 1
	}

	s.regA = (cy) | (x << 1)

	var bit7 byte = (x & 0b1000_0000) >> 7
	s.flags.SetValue(FlagCy, 1 == bit7)
}

// rlc rotates A right
// |7|6|5|4|3|2|1|0| and Cy|c| -> |0|7|6|5|4|3|2|1| and Cy|b0|
func (s *State) rrc() {
	var x byte = s.regA
	var bit0 byte = x & 1
	s.regA = (bit0 << 7) | (x >> 1)
	s.flags.SetValue(FlagCy, 1 == bit0)
}

// rlc rotates A right through Carry
// |7|6|5|4|3|2|1|0| and Cy|c| -> |c|7|6|5|4|3|2|1| and Cy|b0|
func (s *State) rar() {
	var x byte = s.regA

	var cy byte = 0
	if s.flags.IsSet(FlagCy) {
		cy = 1
	}

	s.regA = (cy << 7) | (x >> 1)
	var bit0 byte = x & 1
	s.flags.SetValue(FlagCy, 1 == bit0)
}
