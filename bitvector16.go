package bitfield

import (
	"math/bits"
)

var _ = Bitfield(Bitvector16{})

// Bitvector16 is a bitfield with a fixed defined size of 16. There is no length bit
// present in the underlying byte array.
type Bitvector16 []byte

const bitvector16ByteSize = 2
const bitvector16BitSize = bitvector16ByteSize * 8

// NewBitvector16 creates a new bitvector of size 16.
func NewBitvector16() Bitvector16 {
	byteArray := [bitvector16ByteSize]byte{}
	return byteArray[:]
}

// BitAt returns the bit value at the given index. If the index requested
// exceeds the number of bits in the bitvector, then this method returns false.
func (b Bitvector16) BitAt(idx uint64) bool {
	// Out of bounds or incorrect bitvector byte size, must be false.
	if idx >= b.Len() || len(b) != bitvector16ByteSize {
		return false
	}

	i := uint8(1 << (idx % 8))
	return b[idx/8]&i == i
}

// SetBitAt will set the bit at the given index to the given value. If the index
// requested exceeds the number of bits in the bitvector, then this method returns
// false.
func (b Bitvector16) SetBitAt(idx uint64, val bool) {
	// Out of bounds, do nothing.
	if idx >= b.Len() || len(b) != bitvector16ByteSize {
		return
	}

	bit := uint8(1 << (idx % 8))
	if val {
		b[idx/8] |= bit
	} else {
		b[idx/8] &^= bit
	}
}

// Len returns the number of bits in the bitvector.
func (b Bitvector16) Len() uint64 {
	return bitvector16BitSize
}

// Count returns the number of 1s in the bitvector.
func (b Bitvector16) Count() uint64 {
	if len(b) == 0 {
		return 0
	}
	c := 0
	for i, bt := range b {
		if i >= bitvector16ByteSize {
			break
		}
		c += bits.OnesCount8(bt)
	}
	return uint64(c)
}

// Bytes returns the bytes data representing the Bitvector16.
func (b Bitvector16) Bytes() []byte {
	if len(b) == 0 {
		return []byte{}
	}
	ln := min(len(b), bitvector16ByteSize)
	ret := make([]byte, ln)
	copy(ret, b[:ln])
	return ret[:]
}

// BitIndices returns the list of indices that are set to 1.
func (b Bitvector16) BitIndices() []int {
	indices := make([]int, 0, bitvector16BitSize)
	for i, bt := range b {
		if i >= bitvector16ByteSize {
			break
		}
		for j := 0; j < 8; j++ {
			bit := byte(1 << uint(j))
			if bt&bit == bit {
				indices = append(indices, i*8+j)
			}
		}
	}

	return indices
}

// Contains returns true if the bitlist contains all of the bits from the provided argument
// bitlist. This method will return an error if bitlists are not the same length or not `bitvector16BitSize`.
func (b Bitvector16) Contains(c Bitvector16) (bool, error) {
	if b.Len() != c.Len() {
		return false, ErrBitvectorDifferentLength
	}
	if len(b) != bitvector16ByteSize || len(c) != bitvector16ByteSize {
		return false, ErrWrongLen
	}

	// Combine each byte from b and c, then XOR against b. If the result of this is non-zero for any
	// byte, then we are assured that a byte in c had bits not present in b.
	for i := 0; i < bitvector16ByteSize; i++ {
		if b[i]^(b[i]|c[i]) != 0 {
			return false, nil
		}
	}
	return true, nil
}

// Overlaps returns true if the bitlist contains one of the bits from the provided argument
// bitlist. This method will return an error if bitlists are not the same length.
func (b Bitvector16) Overlaps(c Bitvector16) (bool, error) {
	if b.Len() != c.Len() {
		return false, ErrBitvectorDifferentLength
	}
	if len(b) != bitvector16ByteSize || len(c) != bitvector16ByteSize {
		return false, ErrWrongLen
	}

	// Invert b and xor each byte from b and c, then and it against c. If the result is non-zero for
	// any byte, then we can be assured that a byte in c had bits overlapped in b.
	mask := uint8(0xFF)
	for i := 0; i < bitvector16ByteSize; i++ {
		if (^b[i]^c[i])&c[i]&mask != 0 {
			return true, nil
		}
	}
	return false, nil
}

// Or returns the OR result of the two bitfields. This method will return an error if the bitlists are not the same length.
func (b Bitvector16) Or(c Bitvector16) (Bitvector16, error) {
	if b.Len() != c.Len() {
		return nil, ErrBitvectorDifferentLength
	}
	if len(b) != bitvector16ByteSize || len(c) != bitvector16ByteSize {
		return nil, ErrWrongLen
	}

	ret := make([]byte, bitvector16ByteSize)
	for i := 0; i < bitvector16ByteSize; i++ {
		ret[i] = b[i] | c[i]
	}
	return ret, nil
}
