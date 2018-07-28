package emu

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
)

// Ensure that BIOS implements the Mem interface.
var _ Mem = (*BIOS)(nil)

// BIOS memory map.
const (
	biosBase = 0x1FC00000
	biosLen  = 512 * 1024 // 512 KB
)

// A BIOS is a basic input/output system.
type BIOS struct {
	// Contents of BIOS.
	Data []byte
}

// LoadBIOS loads the BIOS at the given path.
func LoadBIOS(path string) (*BIOS, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &BIOS{Data: data}, nil
}

// LoadUint32 loads a 32-bit unsigned integer from the given offset of the
// memory region.
func (bios *BIOS) LoadUint32(offset uint32) uint32 {
	return binary.LittleEndian.Uint32(bios.Data[offset:])
}

// StoreUint32 stores v to the given offset of the memory region.
func (bios *BIOS) StoreUint32(offset, v uint32) {
	panic(fmt.Errorf("invalid write to offset 0x%08X of BIOS; read-only memory", offset))
}

// Range returns the address range of the memory region.
func (bios *BIOS) Range() Range {
	return Range{biosBase, biosLen}
}
