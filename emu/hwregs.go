package emu

import (
	"fmt"
)

// Ensure that HWRegs implements the Mem interface.
var _ Mem = (*HWRegs)(nil)

// Hardware registers memory map.
const (
	hwRegsBase = 0x1F801000
	hwRegsLen  = 8 * 1024 // 8 KB
)

// HWRegs is a memory map the hardware registers.
type HWRegs struct {
}

// NewHWRegs returns a new memory map the hardware registers.
func NewHWRegs() *HWRegs {
	return &HWRegs{}
}

// LoadUint32 loads a 32-bit unsigned integer from the given offset of the
// memory region.
func (hwregs *HWRegs) LoadUint32(offset uint32) uint32 {
	panic("emu.HWRegs.LoadUint32: not yet implemented")
}

// StoreUint32 stores v to the given offset of the memory region.
func (hwregs *HWRegs) StoreUint32(offset, v uint32) {
	switch offset {
	// Memory control.
	case 0x0000:
		// Expansion 1 base address.
		if v != expansionRegion1Base {
			panic(fmt.Errorf("invalid write to expansion 1 base address; expected 0x%08X, got 0x%08X", expansionRegion1Base, v))
		}
		return
	case 0x0004:
		// Expansion 2 base address.
		if v != expansionRegion2Base {
			panic(fmt.Errorf("invalid write to expansion 2 base address; expected 0x%08X, got 0x%08X", expansionRegion2Base, v))
		}
		return

	// RAM configuration register.
	case 0x0060:
		warn.Printf("support for write to RAM configuration register at offset 0x%04X not yet implemented", offset)
		return
	}

	// Memory control.
	switch {
	case memControlBaseOffset <= offset && offset < memControlBaseOffset+memControlLen:
		// TODO: Implement support.
		warn.Printf("support for write to memory control register at offset 0x%04X not yet implemented", offset)
		return
	}

	// TODO: Implement support.
	panic(fmt.Errorf("support for write to hardware register at offset 0x%04X not yet implemented", offset))
	// nothing to do.
}

// Range returns the address range of the memory region.
func (hwregs *HWRegs) Range() Range {
	return Range{hwRegsBase, hwRegsLen}
}

const (
	memControlBaseOffset = 0x0
	memControlLen        = 36 // 36 B
)

const (
	expansionRegion1Base = 0x1F000000
	expansionRegion1Len  = 8192 * 1024 // 8192 K
)

const (
	expansionRegion2Base = 0x1F802000
)
