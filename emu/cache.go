package emu

// Ensure that CacheControl implements the Mem interface.
var _ Mem = (*CacheControl)(nil)

// Hardware registers memory map.
const (
	cacheControlBase = 0xFFFE0130
	cacheControlLen  = 4 // 4 B
)

// CacheControl is a memory map of the cache control.
type CacheControl struct {
}

// NewCacheControl returns a new memory map of the cache control.
func NewCacheControl() *CacheControl {
	return &CacheControl{}
}

// LoadUint32 loads a 32-bit unsigned integer from the given offset of the
// memory region.
func (cache *CacheControl) LoadUint32(offset uint32) uint32 {
	panic("emu.CacheControl.LoadUint32: not yet implemented")
}

// StoreUint32 stores v to the given offset of the memory region.
func (cache *CacheControl) StoreUint32(offset, v uint32) {
	warn.Printf("support for write to cache control at offset 0x%04X not yet implemented", offset)
}

// Range returns the address range of the memory region.
func (cache *CacheControl) Range() Range {
	return Range{cacheControlBase, cacheControlLen}
}
