// Package emu implements emulation of the PlayStation I video game console.
package emu

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mewkiz/pkg/term"
	"github.com/mewmew/mips"
	"github.com/pkg/errors"
)

// TODO: Remove debugging once the package matures.
var (
	// dbg represents a logger with the "emu:" prefix, which logs debug messages
	// to standard error.
	dbg = log.New(os.Stderr, term.MagentaBold("emu:")+" ", 0)
	// warn represents a logger with the "emu:" prefix, which logs warning
	// messages to standard error.
	warn = log.New(os.Stderr, term.RedBold("emu:")+" ", 0)
)

// A CPU holds the state of a central processing unit.
type CPU struct {
	// Program counter.
	PC uint32
	// Regs is the general purpose register bank of the CPU; where index 0 is
	// $zero and index 31 is $ra.
	Regs [32]uint32
	// Memory maps.
	Mems []Mem
	// Next instruction to be executed.
	DelaySlot mips.Inst
	// Co-processor 0 unit.
	CO0 *CO0
}

// NewCPU returns a new CPU state, as initialized after reset.
func NewCPU(mems ...Mem) *CPU {
	const (
		// Entry point of the PlayStation 1 BIOS.
		entryPoint = biosBase
	)
	// NOP instruction initially placed in the delay slot.
	nop, err := mips.Decode([]byte{0, 0, 0, 0})
	if err != nil {
		panic(errors.Wrap(err, "unable to decode NOP instruction"))
	}
	cpu := &CPU{
		PC:        entryPoint,
		Mems:      mems,
		DelaySlot: nop,
		CO0:       NewCO0(),
	}
	// Init registers with garbage data.
	for i := range cpu.Regs {
		cpu.Regs[i] = 0xDEADC0DE
	}
	// $zero is always 0.
	cpu.Regs[cpu.regIndex(mips.ZERO)] = 0
	return cpu
}

// Step steps a single instruction, performing one fetch, decode and execute
// cycle of the CPU.
func (cpu *CPU) Step() {
	// Use instruction from delay slot in order to handle jump, branch and memory
	// operations which unconditionally execute the instruction directly
	// succedding them.
	inst := cpu.DelaySlot
	dbg.Printf("inst (addr: %08X): %v\n", cpu.PC, inst)

	// Fetch
	bits := cpu.LoadUint32(cpu.PC)
	cpu.PC += 4

	// Decode
	cpu.DelaySlot = cpu.Decode(bits)

	// Execute
	cpu.Execute(inst)
}

// Isolate cache if bit in SR register of CO0 is set.
const isolateCacheMask = 0x00010000

// LoadUint32 loads a 32-bit unsigned integer from the given address.
func (cpu *CPU) LoadUint32(addr uint32) uint32 {
	if addr%4 != 0 {
		panic(fmt.Errorf("unaligned access of memory at address 0x%08X", addr))
	}
	for _, mem := range cpu.Mems {
		if offset, ok := mem.Range().Contains(addr); ok {
			return mem.LoadUint32(offset)
		}
	}
	panic(fmt.Errorf("unable to load value from address 0x%08X", addr))
}

// StoreUint32 stores v to the given address.
func (cpu *CPU) StoreUint32(addr, v uint32) {
	if addr%4 != 0 {
		panic(fmt.Errorf("unaligned access of memory at address 0x%08X", addr))
	}
	if cpu.CO0.Reg(mips.SR)&isolateCacheMask != 0 {
		// Cache is isolated, ignore write.
		warn.Println("write with isolate cache not yet implemented")
		return
	}
	// TODO: Remove debug output.
	dbg.Printf("store at: 0x%08X (%d)\n", addr, v)
	for _, mem := range cpu.Mems {
		if offset, ok := mem.Range().Contains(addr); ok {
			mem.StoreUint32(offset, v)
			return
		}
	}
	panic(fmt.Errorf("unable to store value at address 0x%08X; no memory map found for region", addr))
}

// Decode decodes the given MIPS I instruction bit pattern.
func (*CPU) Decode(bits uint32) mips.Inst {
	var src [4]byte
	binary.LittleEndian.PutUint32(src[:], bits)
	inst, err := mips.Decode(src[:])
	if err != nil {
		panic(errors.Wrapf(err, "unable to decode instruction (0x%08X)", bits))
	}
	return inst
}

// Execute executes the given instruction on the CPU.
func (cpu *CPU) Execute(inst mips.Inst) {
	switch inst.Op {
	case mips.LUI:
		// LUI     $t, immediate
		t := inst.Args[0].(mips.Reg)
		i := inst.Args[1].(mips.Imm)
		cpu.SetReg(t, i.Imm<<16)
	case mips.ORI:
		// ORI     $t, $s, immediate
		t := inst.Args[0].(mips.Reg)
		s := inst.Args[1].(mips.Reg)
		i := inst.Args[2].(mips.Imm)
		cpu.SetReg(t, cpu.Reg(s)|i.Imm)
	case mips.SW:
		// SW      $t, offset($s)
		t := inst.Args[0].(mips.Reg)
		m := inst.Args[1].(mips.Mem)
		addr := cpu.Reg(m.Base) + uint32(m.Offset)
		v := cpu.Reg(t)
		cpu.StoreUint32(addr, v)
	case mips.SLL:
		// SLL     $d, $t, shift
		d := inst.Args[0].(mips.Reg)
		t := inst.Args[1].(mips.Reg)
		shift := inst.Args[2].(mips.Imm)
		cpu.SetReg(d, cpu.Reg(t)<<shift.Imm)
	case mips.ADDIU:
		// ADDIU   $t, $s, immediate
		t := inst.Args[0].(mips.Reg)
		s := inst.Args[1].(mips.Reg)
		i := inst.Args[2].(mips.Imm)
		cpu.SetReg(t, cpu.Reg(s)+i.Imm)
	case mips.ADDI:
		// ADDI    $t, $s, immediate
		t := inst.Args[0].(mips.Reg)
		s := inst.Args[1].(mips.Reg)
		i := inst.Args[2].(mips.Imm)
		// TODO: raise exception on overflow.
		cpu.SetReg(t, cpu.Reg(s)+i.Imm)
	case mips.SLT:
		// SLT     $d, $s, $t
		d := inst.Args[0].(mips.Reg)
		s := inst.Args[1].(mips.Reg)
		t := inst.Args[2].(mips.Reg)
		if cpu.Reg(s) < cpu.Reg(t) {
			cpu.SetReg(d, 1)
		} else {
			cpu.SetReg(d, 0)
		}
	case mips.J:
		// J       target
		target := inst.Args[0].(mips.Imm)
		cpu.PC = cpu.PC&0xF0000000 + target.Imm
	case mips.OR:
		// OR      $d, $s, $t
		d := inst.Args[0].(mips.Reg)
		s := inst.Args[1].(mips.Reg)
		t := inst.Args[2].(mips.Reg)
		cpu.SetReg(d, cpu.Reg(s)|cpu.Reg(t))
	case mips.MTC0:
		// MTC0    $t, $d
		t := inst.Args[0].(mips.Reg)
		d := inst.Args[1].(mips.Reg)
		cpu.CO0.SetReg(d, cpu.Reg(t))

	case mips.BNE:
		// BNE     $s, $t, offset
		s := inst.Args[0].(mips.Reg)
		t := inst.Args[1].(mips.Reg)
		offset := inst.Args[2].(mips.PCRel)
		if cpu.Reg(s) != cpu.Reg(t) {
			cpu.PC += uint32(offset)
		}
	case mips.LW:
		// LW      $t, offset($s)
		if cpu.CO0.Reg(mips.SR)&isolateCacheMask != 0 {
			// Cache is isolated, ignore read.
			warn.Println("read with isolate cache not yet implemented")
			return
		}
		t := inst.Args[0].(mips.Reg)
		m := inst.Args[1].(mips.Mem)
		addr := cpu.Reg(m.Base) + uint32(m.Offset)
		cpu.SetReg(t, cpu.LoadUint32(addr))
		// TODO: Figure out how to handle load delay slot.
	case mips.SLTU:
		// SLTU    $d, $s, $t
		d := inst.Args[0].(mips.Reg)
		s := inst.Args[1].(mips.Reg)
		t := inst.Args[2].(mips.Reg)
		if cpu.Reg(s) < cpu.Reg(t) {
			cpu.SetReg(d, 1)
		} else {
			cpu.SetReg(d, 0)
		}
	// TODO: Check difference bewteen SLTU and SLT.
	case mips.ADDU:
		// ADDU    $d, $s, $t
		d := inst.Args[0].(mips.Reg)
		s := inst.Args[1].(mips.Reg)
		t := inst.Args[2].(mips.Reg)
		cpu.SetReg(d, cpu.Reg(s)+cpu.Reg(t))
	default:
		panic(fmt.Errorf("support for instruction opcode %q not yet implemented", inst.Op))
	}
}

// String returns the string representation of the CPU state.
func (cpu *CPU) String() string {
	buf := &strings.Builder{}
	// Print registers.
	maxRegNameLen := len(mips.ZERO.String())
	buf.WriteString("; CPU registers\n")
	for r := mips.ZERO; r <= mips.RA; r++ {
		fmt.Fprintf(buf, "%-*s = 0x%08X\n", maxRegNameLen, r, cpu.Reg(r))
	}
	buf.WriteString("; Program counter\n")
	fmt.Fprintf(buf, "%-*s = 0x%08X\n", maxRegNameLen, mips.PC, cpu.PC)
	return buf.String()
}

// Reg returns the contents of the given register.
func (cpu *CPU) Reg(r mips.Reg) uint32 {
	return cpu.Regs[cpu.regIndex(r)]
}

// SetReg sets the contents of the given register to v; taking precaution to
// always keep $zero at 0.
func (cpu *CPU) SetReg(r mips.Reg, v uint32) {
	// $zero is always 0.
	if r == mips.ZERO {
		return
	}
	cpu.Regs[cpu.regIndex(r)] = v
}

// regIndex returns the index of the CPU register, starting at 0.
func (cpu *CPU) regIndex(r mips.Reg) int {
	cpu.validateReg(r)
	return int(r - mips.ZERO)
}

// validateReg validates the given CPU register.
func (cpu *CPU) validateReg(r mips.Reg) {
	switch {
	case mips.ZERO <= r && r <= mips.RA:
		// valid CPU register.
	default:
		panic(fmt.Errorf("invalid register %v; not present on CPU", r))
	}
}
