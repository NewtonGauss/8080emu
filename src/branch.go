package main

func (s *State) jmpOnFlag(f Flag, set bool) {
	if s.flags.IsSet(f) == set {
		s.jmp()
	} else {
		s.pc += 3
	}
}

func (s *State) jmp() {
	s.pc = pairTo16(s.mem[s.pc+1], s.mem[s.pc+2])
}

func (s *State) callOnFlag(f Flag, set bool) {
	if s.flags.IsSet(f) == set {
		s.call()
	} else {
		s.pc += 3
	}
}

func (s *State) call() {
	var ret = s.pc + 3

	// Push the return address into the stack. Stack goes "down"
	s.mem[s.sp-1] = byte(ret >> 8)
	s.mem[s.sp-2] = byte(ret)
	s.sp += 2

	s.jmp()
}

func (s *State) rst(addr uint16) {
	var ret = s.pc + 3

	s.mem[s.sp-1] = byte(ret >> 8)
	s.mem[s.sp-2] = byte(ret)
	s.sp += 2

	s.pc = addr
}

func (s *State) retOnFlag(f Flag, set bool) {
	if s.flags.IsSet(f) == set {
		s.ret()
	} else {
		s.pc += 1
	}
}

func (s *State) ret() {
	s.pc = pairTo16(s.mem[s.sp+1], s.mem[s.sp])
	s.sp += 2
}
