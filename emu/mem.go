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

// Segments.
const (
	// KUSEG, 2 GB
	kusegStart = 0x00000000
	kusegEnd   = 0x7FFFFFFF
	kusegSize  = kusegEnd - kusegStart + 1
	// KSEG0, 512 MB
	kseg0Start = 0x80000000
	kseg0End   = 0x9FFFFFFF
	kseg0Size  = kseg0End - kseg0Start + 1
	// KSEG1, 512 MB
	kseg1Start = 0xA0000000
	kseg1End   = 0xBFFFFFFF
	kseg1Size  = kseg1End - kseg1Start + 1
	// KSEG2, 1 GB
	kseg2Start = 0xC0000000
	kseg2End   = 0xFFFFFFFF
	kseg2Size  = kseg2End - kseg2Start + 1
)

// maskSegment returns the corresponding address in KUSEG of the given address.
func maskSegment(addr uint32) uint32 {
	// masks maps from the 3 most significant bits of the address to the mask
	// required to convert the address into a corresponding address in KUSEG.
	masks := [8]uint32{
		// KUSEG: 2048 MB
		//    000, 001, 010, 011
		0xFFFFFFFF, 0xFFFFFFFF, 0xFFFFFFFF, 0xFFFFFFFF, // address already in KUSEG.
		// KSEG0: 512 MB
		//    100
		0x7FFFFFFF, // clear the most-significant bit get KUSEG address.
		// KSEG1: 512 MB
		//    101
		0x1FFFFFFF, // clear the 3 most-significant bits to get KUSEG address.
		// KSEG2: 1024 MB
		//    110, 111
		0xFFFFFFFF, 0xFFFFFFFF, // keep KSEG2 address unmodifed.
	}
	// get the 3 most-significant bits.
	i := addr >> 29
	return addr & masks[i]
}
