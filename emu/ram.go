package emu

import (
	"encoding/binary"
)

// Ensure that RAM implements the Mem interface.
var _ Mem = (*RAM)(nil)

// RAM memory map.
const (
	ramBase = 0xA0000000
	ramLen  = 2048 * 1024 // 2048 KB
)

// A RAM is a random access memory.
type RAM struct {
	// Contents of RAM.
	Data []byte
}

// NewRAM returns a new RAM memory region.
func NewRAM() *RAM {
	data := make([]byte, ramLen)
	// Init RAM with garbage data.
	for i := range data {
		data[i] = 0xCA
	}
	return &RAM{Data: data}
}

// LoadUint32 loads a 32-bit unsigned integer from the given offset of the
// memory region.
func (ram *RAM) LoadUint32(offset uint32) uint32 {
	return binary.LittleEndian.Uint32(ram.Data[offset:])
}

// StoreUint32 stores v to the given offset of the memory region.
func (ram *RAM) StoreUint32(offset, v uint32) {
	binary.LittleEndian.PutUint32(ram.Data[offset:], v)
}

// Range returns the address range of the memory region.
func (ram *RAM) Range() Range {
	return Range{ramBase, ramLen}
}
