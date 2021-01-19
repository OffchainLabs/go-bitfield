package bitfield

import (
	"encoding/binary"
	"math/bits"
)

var _ = Bitfield(&Bitlist64{})

const (
	// wordSize configures how many bits are there in a single element of bitlist array.
	wordSize = uint64(64)
	// wordSizeLog2 allows optimized division by wordSize using right shift (numBits >> wordSizeLog2).
	// Note: log_2(64) = 6.
	wordSizeLog2 = uint64(6)
	// bytesInWord defines how many bytes are there in a single word i.e. wordSize/8.
	bytesInWord = 8
	// bytesInWordLog2 = log_2(8)
	bytesInWordLog2 = 3
	// allBitsSet is a word with all bits set.
	allBitsSet = uint64(0xffffffffffffffff)
)

// Bitlist64 is a bitfield implementation backed by an array of uint64.
type Bitlist64 struct {
	size uint64
	data []uint64
}

// NewBitlist64 creates a new bitlist of size N.
func NewBitlist64(n uint64) *Bitlist64 {
	return &Bitlist64{
		size: n,
		data: make([]uint64, numWordsRequired(n)),
	}
}

// NewBitlist64From creates a new bitlist for a given uint64 array.
func NewBitlist64From(data []uint64) *Bitlist64 {
	return &Bitlist64{
		size: uint64(len(data)) * wordSize,
		data: data,
	}
}

// BitAt returns the bit value at the given index. If the index requested
// exceeds the number of bits in the bitlist, then this method returns false.
func (b *Bitlist64) BitAt(idx uint64) bool {
	// Out of bounds, must be false.
	if idx >= b.size {
		return false
	}

	i := uint64(1 << (idx % wordSize))
	return b.data[idx>>wordSizeLog2]&i == i
}

// SetBitAt will set the bit at the given index to the given value.
// If the index requested exceeds the number of bits in the bitlist, then this method returns false.
func (b *Bitlist64) SetBitAt(idx uint64, val bool) {
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
func (b *Bitlist64) Len() uint64 {
	return b.size
}

// Bytes returns underlying array of uint64s as an array of bytes.
// The leading zeros in the bitlist will be trimmed to the smallest byte length representation of
// the bitlist. This may produce an empty byte slice if all bits were zero.
func (b *Bitlist64) Bytes() []byte {
	if len(b.data) == 0 {
		return []byte{}
	}

	ret := make([]byte, len(b.data)*bytesInWord)
	for idx, word := range b.data {
		start := idx << bytesInWordLog2
		binary.LittleEndian.PutUint64(ret[start:start+bytesInWord], word)
	}

	// Clear any leading zero bytes.
	allLeadingZeroes := 0
	for i := len(b.data) - 1; i >= 0; i-- {
		leadingZeroes := 0
		if b.data[i] == 0 {
			leadingZeroes = int(wordSize)
		} else {
			leadingZeroes = bits.LeadingZeros64(b.data[i])
		}
		allLeadingZeroes += leadingZeroes
		// If the whole word is 0x0, allow to test the next word, break otherwise.
		if uint64(leadingZeroes) != wordSize {
			break
		}
	}

	return ret[:len(ret)-allLeadingZeroes>>bytesInWordLog2]
}

// Count returns the number of 1s in the bitlist.
func (b *Bitlist64) Count() uint64 {
	c := 0
	for _, bt := range b.data {
		c += bits.OnesCount64(bt)
	}

	return uint64(c)
}

// Contains returns true if the bitlist contains all of the bits from the provided argument
// bitlist i.e. if `b` is a superset of `c`.
// This method will panic if bitlists are not the same length.
func (b *Bitlist64) Contains(c *Bitlist64) bool {
	if b.Len() != c.Len() {
		panic("bitlists are different lengths")
	}

	// To ensure all of the bits in c are present in b, we iterate over every word, combine
	// the words from b and c, then XOR them against b. If the result of this is non-zero, then we
	// are assured that a word in c had bits not present in word in b.
	for idx, word := range b.data {
		if word^(word|c.data[idx]) != 0 {
			return false
		}
	}

	return true
}

// Overlaps returns true if the bitlist contains one of the bits from the provided argument
// bitlist. This method will panic if bitlists are not the same length.
func (b *Bitlist64) Overlaps(c *Bitlist64) bool {
	lenB, lenC := b.Len(), c.Len()
	if lenB != lenC {
		panic("bitlists are different lengths")
	}

	if lenB == 0 || lenC == 0 {
		return false
	}

	// To ensure all of the bits in c are not overlapped in b, we iterate over every word, invert b
	// and xor the word from b and c, then and it against c. If the result is non-zero, then
	// we can be assured that word in c had bits not overlapped in b.
	for idx, word := range b.data {
		if (^word^c.data[idx])&c.data[idx]&allBitsSet != 0 {
			return true
		}
	}

	return false
}

// Or returns the OR result of the two bitfields (union).
// This method will panic if the bitlists are not the same length.
func (b *Bitlist64) Or(c *Bitlist64) *Bitlist64 {
	if b.Len() != c.Len() {
		panic("bitlists are different lengths")
	}

	ret := b.Clone()
	b.NoAllocOr(c, ret)

	return ret
}

// NoAllocOr computes the OR result of the two bitfields (union).
// Result is written into provided variable, so no allocation takes place inside the function.
// This method will panic if the bitlists are not the same length.
func (b *Bitlist64) NoAllocOr(c, ret *Bitlist64) {
	if b.Len() != c.Len() {
		panic("bitlists are different lengths")
	}

	for idx, word := range b.data {
		ret.data[idx] = word | c.data[idx]
	}
}

// And returns the AND result of the two bitfields (intersection).
// This method will panic if the bitlists are not the same length.
func (b *Bitlist64) And(c *Bitlist64) *Bitlist64 {
	if b.Len() != c.Len() {
		panic("bitlists are different lengths")
	}

	ret := b.Clone()
	b.NoAllocAnd(c, ret)

	return ret
}

// NoAllocAnd computes the AND result of the two bitfields (intersection).
// Result is written into provided variable, so no allocation takes place inside the function.
// This method will panic if the bitlists are not the same length.
func (b *Bitlist64) NoAllocAnd(c, ret *Bitlist64) {
	if b.Len() != c.Len() {
		panic("bitlists are different lengths")
	}

	for idx, word := range b.data {
		ret.data[idx] = word & c.data[idx]
	}
}

// Xor returns the XOR result of the two bitfields (symmetric difference).
// This method will panic if the bitlists are not the same length.
func (b *Bitlist64) Xor(c *Bitlist64) *Bitlist64 {
	if b.Len() != c.Len() {
		panic("bitlists are different lengths")
	}

	ret := b.Clone()
	b.NoAllocXor(c, ret)

	return ret
}

// NoAllocXor returns the XOR result of the two bitfields (symmetric difference).
// Result is written into provided variable, so no allocation takes place inside the function.
// This method will panic if the bitlists are not the same length.
func (b *Bitlist64) NoAllocXor(c, ret *Bitlist64) {
	if b.Len() != c.Len() {
		panic("bitlists are different lengths")
	}

	for idx, word := range b.data {
		ret.data[idx] = word ^ c.data[idx]
	}
}

// Not returns the NOT result of the bitfield (complement).
func (b *Bitlist64) Not() *Bitlist64 {
	if b.Len() == 0 {
		return b
	}

	ret := b.Clone()
	b.NoAllocNot(ret)

	return ret
}

// NoAllocNot returns the NOT result of the bitfield (complement).
// Result is written into provided variable, so no allocation takes place inside the function.
func (b *Bitlist64) NoAllocNot(ret *Bitlist64) {
	if b.Len() == 0 {
		return
	}

	for idx, word := range b.data {
		ret.data[idx] = ^word
	}
}

// BitIndices returns list of bit indexes of bitlist where value is set to true.
func (b *Bitlist64) BitIndices() []int {
	indices := make([]int, b.Count())
	b.NoAllocBitIndices(indices)

	return indices
}

// NoAllocBitIndices returns list of bit indexes of bitlist where value is set to true.
// No allocation happens inside the function, so number of returned indexes is capped by the capacity
// of the ret param.
//
// Expected usage pattern:
//
// b := NewBitlist64(n)
// indices := make([]int, b.Count())
// b.NoAllocBitIndices(indices)
func (b *Bitlist64) NoAllocBitIndices(ret []int) {
	capacity := cap(ret)
	k := 0
	processWord := func(idx int, word uint64) uint64 {
		// Push index of the first non-zero bit.
		ret[k] = (idx << wordSizeLog2) + bits.TrailingZeros64(word)
		k++
		if k == capacity {
			return 0
		}
		// Clear less significant (rightmost) non-zero bit, and iterate.
		// Consider the following bitlist, b := 0001.1001.0011.0000
		// The `(^word) + 1` clears all bits till the word's non-zero bit i.e. `(^word)` == 1110.0110.1100.1111,
		// then `(^word) + 1` == 1110.0110.1101.0000.
		// The `word & ((^word) + 1)` clears all bits, except the one that was set to 1 in the original word i.e.
		// `word & ((^word) + 1)` == 0000.0000.0001.0000.
		// Now, XOR this with the original word to remove the rightmost bit.
		return word ^ (word & ((^word) + 1))
	}

	for idx, word := range b.data {
		for word != 0 {
			word = processWord(idx, word)
		}
	}
}

// Clone safely copies a given bitlist.
func (b *Bitlist64) Clone() *Bitlist64 {
	c := NewBitlist64(b.size)
	if b.data != nil {
		copy(c.data, b.data)
	}
	return c
}

// numWordsRequired calculates how many words are required to hold bitlist of n bits.
func numWordsRequired(n uint64) int {
	return int((n + (wordSize - 1)) >> wordSizeLog2)
}
