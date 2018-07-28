package emu

import (
	"fmt"
	"strings"

	"github.com/mewmew/mips"
)

// A CO0 holds the state of a co-processor 0 unit.
type CO0 struct {
	// Register bank of co-processor 0.
	Regs [32]uint32
}

// NewCO0 returns a new co-processor 0 state, as initialized after reset.
func NewCO0() *CO0 {
	co0 := &CO0{}
	// Init registers with garbage data.
	for i := range co0.Regs {
		co0.Regs[i] = 0xDEADC0DE
	}
	// TODO: Verify the correct start values of CO0 registers.
	co0.SetReg(mips.SR, 0)
	return co0
}

// String returns the string representation of the CO0 state.
func (co0 *CO0) String() string {
	buf := &strings.Builder{}
	// Print registers.
	var maxRegNameLen int
	for r := mips.CO0Reg0; r <= mips.CO0Reg31; r++ {
		n := len(r.String())
		if n > maxRegNameLen {
			maxRegNameLen = n
		}
	}
	buf.WriteString("; CO0 registers\n")
	for r := mips.CO0Reg0; r <= mips.CO0Reg31; r++ {
		fmt.Fprintf(buf, "%-*s = 0x%08X\n", maxRegNameLen, r, co0.Reg(r))
	}
	return buf.String()
}

// Reg returns the contents of the given register.
func (co0 *CO0) Reg(r mips.Reg) uint32 {
	return co0.Regs[co0.regIndex(r)]
}

// SetReg sets the contents of the given register to v; taking precaution to
// always keep $zero at 0.
func (co0 *CO0) SetReg(r mips.Reg, v uint32) {
	co0.Regs[co0.regIndex(r)] = v
}

// regIndex returns the index of the CO0 register, starting at 0.
func (co0 *CO0) regIndex(r mips.Reg) int {
	co0.validateReg(r)
	return int(r - mips.CO0Reg0)
}

// validateReg validates the given CO0 register.
func (co0 *CO0) validateReg(r mips.Reg) {
	switch {
	case mips.CO0Reg0 <= r && r <= mips.CO0Reg31:
		// valid CO0 register.
	default:
		panic(fmt.Errorf("invalid register %v; not present on CO0", r))
	}
}
