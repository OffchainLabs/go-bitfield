package bitfield

import "math/bits"

const (
	// wordSize configures how many bits are there in a single element of bitlist array.
	wordSize = uint64(64)
	// wordSizeLog2 allows optimized division by wordSize using right shift (numBits >> wordSizeLog2).
	// Note: log_2(64) = 6.
	wordSizeLog2 = uint64(6)
)

// Bitlist is a bitfield implementation backed by an array of uint64.
type Bitlist struct {
	size uint64
	data []uint64
}

// NewBitlist creates a new bitlist of size N.
func NewBitlist(n uint64) *Bitlist {
	numWords := numWordsRequired(n)
	return &Bitlist{
		size: n,
		data: make([]uint64, numWords, numWords),
	}
}

// NewBitlistFrom creates a new bitlist for a given uint64 array.
func NewBitlistFrom(data []uint64) *Bitlist {
	return &Bitlist{
		size: uint64(len(data)) * wordSize,
		data: data,
	}
}

// BitAt returns the bit value at the given index. If the index requested
// exceeds the number of bits in the bitlist, then this method returns false.
func (b *Bitlist) BitAt(idx uint64) bool {
	// Out of bounds, must be false.
	if idx >= b.size {
		return false
	}

	i := uint64(1 << (idx % wordSize))
	return b.data[idx>>wordSizeLog2]&i == i
}

// SetBitAt will set the bit at the given index to the given value. If the index
// requested exceeds the number of bits in the bitlist, then this method returns
// false.
func (b *Bitlist) SetBitAt(idx uint64, val bool) {
	// Out of bounds, do nothing.
	if idx >= b.size {
		return
	}

	bit := uint64(1 << (idx % wordSize))
	if val {
		b.data[idx>>wordSizeLog2] |= bit
	} else {
		b.data[idx>>wordSizeLog2] &^= bit
	}
}

// Len returns the number of bits in a bitlist (note that underlying array can be bigger).
func (b *Bitlist) Len() uint64 {
	return b.size
}

// Count returns the number of 1s in the bitlist.
func (b *Bitlist) Count() uint64 {
	c := 0
	for _, bt := range b.data {
		c += bits.OnesCount64(bt)
	}
	return uint64(c)
}

// numWordsRequired calculates how many words are required to hold bitlist of n bits.
func numWordsRequired(n uint64) int {
	return int((n + (wordSize - 1)) >> wordSizeLog2)
}
