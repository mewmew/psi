package emu

// A Mem is a region of memory.
type Mem interface {
	// LoadUint32 loads a 32-bit unsigned integer from the given offset of the
	// memory region.
	LoadUint32(offset uint32) uint32
	// StoreUint32 stores v to the given offset of the memory region.
	StoreUint32(offset, v uint32)
	// Range returns the address range of the memory region.
	Range() Range
}

// A Range is a region of memory.
type Range struct {
	// Base address.
	Base uint32
	// Length in bytes.
	Len uint32
}

// Contains reports whether the region of memory contains the given address, and
// if so returns the corresponding relative offset.
func (r Range) Contains(addr uint32) (offset uint32, ok bool) {
	if r.Base <= addr && addr < r.Base+r.Len {
		return addr - r.Base, true
	}
	return 0, false
}
