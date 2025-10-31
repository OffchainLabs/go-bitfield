package bitfield

import (
	"math/bits"
)

var _ = Bitfield(Bitvector2{})

// Bitvector2 is a bitfield with a known size of 2. There is no length bit
// present in the underlying byte array.
type Bitvector2 []byte

const bitvector2ByteSize = 1
const bitvector2BitSize = 2

// NewBitvector2 creates a new bitvector of size 2.
func NewBitvector2() Bitvector2 {
	byteArray := [bitvector2ByteSize]byte{}
	return byteArray[:]
}

// BitAt returns the bit value at the given index. If the index requested
// exceeds the number of bits in the bitvector, then this method returns false.
func (b Bitvector2) BitAt(idx uint64) bool {
	// Out of bounds, must be false.
	if idx >= b.Len() || len(b) != bitvector2ByteSize {
		return false
	}

	i := uint8(1 << idx)
	return b[0]&i == i

}

// SetBitAt will set the bit at the given index to the given value. If the index
// requested exceeds the number of bits in the bitvector, then this method returns
// false.
func (b Bitvector2) SetBitAt(idx uint64, val bool) {
	// Out of bounds, do nothing.
	if idx >= b.Len() || len(b) != bitvector2ByteSize {
		return
	}

	bit := uint8(1 << idx)
	if val {
		b[0] |= bit
	} else {
		b[0] &^= bit
	}
}

// Len returns the number of bits in the bitvector.
func (b Bitvector2) Len() uint64 {
	return bitvector2BitSize
}

// Count returns the number of 1s in the bitvector.
func (b Bitvector2) Count() uint64 {
	if len(b) == 0 {
		return 0
	}
	return uint64(bits.OnesCount8(b.Bytes()[0]))
}

// Bytes returns the bytes data representing the bitvector2. This method
// bitmasks the underlying data to ensure that it is an accurate representation.
func (b Bitvector2) Bytes() []byte {
	if len(b) == 0 {
		return []byte{}
	}
	return []byte{b[0] & 0x03}
}

// Shift bitvector by i. If i >= 0, perform left shift, otherwise right shift.
func (b Bitvector2) Shift(i int) {
	if len(b) == 0 {
		return
	}

	// Shifting greater than 2 bits is pointless and can have unexpected behavior.
	if i > 2 {
		i = 2
	} else if i < -2 {
		i = -2
	}

	if i >= 0 {
		b[0] <<= uint8(i)
	} else {
		b[0] >>= uint8(i * -1)
	}
	b[0] &= 0x03
}

// BitIndices returns the list of indices that are set to 1.
func (b Bitvector2) BitIndices() []int {
	indices := make([]int, 0, 2)
	if len(b) != bitvector2ByteSize {
		return indices
	}

	bt := b[0]
	for j := 0; j < bitvector2BitSize; j++ {
		bit := byte(1 << uint(j))
		if bt&bit == bit {
			indices = append(indices, j)
		}
	}

	return indices
}
