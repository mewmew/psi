package emu

// Ensure that SPU implements the Mem interface.
// TODO: re-enable assertion.
//var _ Mem = (*SPU)(nil)

// SPU memory map.
const (
	spuBase = 0x1F801C00
	spuLen  = 640 // 640 B
)

// A SPU is a sound processing unit.
type SPU struct {
}

// NewSPU returns a new SPU memory region.
func NewSPU() *SPU {
	return &SPU{}
}

// StoreUint16 stores v to the given offset of the memory region.
func (spu *SPU) StoreUint16(offset uint32, v uint16) {
	warn.Printf("support for write to SPU at offset 0x%04X not yet implemented", offset)
	// TODO: Implement SPU write support.
}

// Range returns the address range of the memory region.
func (spu *SPU) Range() Range {
	return Range{spuBase, spuLen}
}
