package main

import (
	"testing"
)

func TestFlags(t *testing.T) {
	var f Flags

	var flags = []Flag{FlagZ, FlagS, FlagP, FlagCy, FlagAc}
	for _, fl := range flags {
		if f.IsSet(fl) {
			t.Errorf("IsSet(%v) = true. None has been set", fl)
		}
	}

	f.Set(FlagCy)
	if !f.IsSet(FlagCy) {
		t.Errorf("IsSet(FlagCy) = false after f.Set(FlagCy)")
	}

	f.Unset(FlagCy)
	for _, fl := range flags {
		if f.IsSet(fl) {
			t.Errorf("IsSet(%v) = true. After f.Unset(FlagCy)", fl)
		}
	}
}
