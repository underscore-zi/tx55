package bitfield

import "fmt"

type Bitfield struct {
	Data byte
}

// SetBits can be used to treat a set of multiple bits as a single value
// Position is zero-based index of the first bit to set starting with the least significant bit
// Size is how many bits larger to include
func (b *Bitfield) SetBits(position, size byte, value byte) {
	if value >= (1 << size) {
		panic(fmt.Sprintf("SetBits: value %d out of range for size %d", value, size))
	}

	// Create a mask with the specified bits set to 1
	mask := byte(1<<size-1) << position

	// Clear the bits to be modified
	b.Data &^= mask

	// Set the bits to the specified value
	b.Data |= (value << position) & mask
}

// SetBit can be used to set a single bit
func (b *Bitfield) SetBit(bit byte, value bool) {
	if value == false {
		b.Data &^= 1 << bit
	} else {
		b.Data |= 1 << bit
	}
}

// GetBit can be used to get a single bit
func (b *Bitfield) GetBit(bit byte) bool {
	return b.Data&(1<<bit) != 0
}

// GetBits can get a set of multiple bits as a single value
func (b *Bitfield) GetBits(position, size byte) byte {
	return (b.Data >> position) & ((1 << size) - 1)
}
