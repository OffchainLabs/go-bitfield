package bitfield

import (
	"bytes"
	"fmt"
	"testing"
)

// Recover within each subtest so the unfixed implementation reports every
// regression instead of stopping the suite at the first out-of-range access.
func failOnBitvectorPanic(t *testing.T) {
	t.Helper()
	if r := recover(); r != nil {
		t.Errorf("unexpected panic: %v", r)
	}
}

func TestBitvector8_BinaryOpsWrongLength(t *testing.T) {
	operands := []struct {
		name string
		bits Bitvector8
	}{
		{"nil", nil},
		{"empty", Bitvector8{}},
		{"valid", Bitvector8{0x01}},
		{"oversized", Bitvector8{0x01, 0x02}},
	}
	for _, a := range operands {
		for _, b := range operands {
			if len(a.bits) == 1 && len(b.bits) == 1 {
				continue // Valid operations are covered by the existing tests.
			}
			for _, op := range []string{"Contains", "Overlaps", "Or"} {
				t.Run(a.name+"/"+b.name+"/"+op, func(t *testing.T) {
					defer failOnBitvectorPanic(t)
					switch op {
					case "Contains":
						if got, err := a.bits.Contains(b.bits); got || err != ErrWrongLen {
							t.Errorf("Contains = (%v, %v), want (false, %v)", got, err, ErrWrongLen)
						}
					case "Overlaps":
						if got, err := a.bits.Overlaps(b.bits); got || err != ErrWrongLen {
							t.Errorf("Overlaps = (%v, %v), want (false, %v)", got, err, ErrWrongLen)
						}
					case "Or":
						if got, err := a.bits.Or(b.bits); got != nil || err != ErrWrongLen {
							t.Errorf("Or = (%x, %v), want (nil, %v)", got, err, ErrWrongLen)
						}
					}
				})
			}
		}
	}
}

func TestBitvector128_BinaryOpsDifferentLengths(t *testing.T) {
	for _, lenA := range []int{0, 1, 15, 16, 17} {
		for _, lenB := range []int{0, 1, 15, 16, 17} {
			if lenA == lenB {
				continue
			}
			// Zeros prevent early returns from masking an unsafe later access.
			a, b := make(Bitvector128, lenA), make(Bitvector128, lenB)
			for _, op := range []string{"Contains", "Overlaps", "Or"} {
				t.Run(fmt.Sprintf("len%d/len%d/%s", lenA, lenB, op), func(t *testing.T) {
					defer failOnBitvectorPanic(t)
					switch op {
					case "Contains":
						if got, err := a.Contains(b); got || err != ErrBitvectorDifferentLength {
							t.Errorf("Contains = (%v, %v), want (false, %v)", got, err, ErrBitvectorDifferentLength)
						}
					case "Overlaps":
						if got, err := a.Overlaps(b); got || err != ErrBitvectorDifferentLength {
							t.Errorf("Overlaps = (%v, %v), want (false, %v)", got, err, ErrBitvectorDifferentLength)
						}
					case "Or":
						if got, err := a.Or(b); got != nil || err != ErrBitvectorDifferentLength {
							t.Errorf("Or = (%x, %v), want (nil, %v)", got, err, ErrBitvectorDifferentLength)
						}
					}
				})
			}
		}
	}
}

func TestBitvector128_BinaryOpsEqualLengthsRemainSupported(t *testing.T) {
	// This fix must not impose a new strict 16-byte length contract.
	for _, size := range []int{0, 1, 15, 16, 17} {
		t.Run(fmt.Sprintf("len%d", size), func(t *testing.T) {
			a := Bitvector128(bytes.Repeat([]byte{0x03}, size))
			b := Bitvector128(bytes.Repeat([]byte{0x01}, size))
			if got, err := a.Contains(b); !got || err != nil {
				t.Errorf("Contains = (%v, %v), want (true, nil)", got, err)
			}
			if got, err := a.Overlaps(b); got != (size > 0) || err != nil {
				t.Errorf("Overlaps = (%v, %v), want (%v, nil)", got, err, size > 0)
			}
			if got, err := a.Or(b); !bytes.Equal(got, a) || err != nil {
				t.Errorf("Or = (%x, %v), want (%x, nil)", got, err, a)
			}
		})
	}
}

func TestBitvector_ShiftShortSlices(t *testing.T) {
	vectors := []struct {
		name  string
		shift func([]byte, int)
	}{
		{"64", func(b []byte, n int) { Bitvector64(b).Shift(n) }},
		{"128", func(b []byte, n int) { Bitvector128(b).Shift(n) }},
		{"256", func(b []byte, n int) { Bitvector256(b).Shift(n) }},
		{"512", func(b []byte, n int) { Bitvector512(b).Shift(n) }},
	}
	for _, vector := range vectors {
		for size := -1; size <= 8; size++ {
			for _, shift := range []int{-1, 0, 1} {
				t.Run(fmt.Sprintf("%s/len%d/shift%d", vector.name, size, shift), func(t *testing.T) {
					defer failOnBitvectorPanic(t)
					var b []byte // size -1 exercises nil separately from empty.
					if size >= 0 {
						b = bytes.Repeat([]byte{0x02}, size)
					}
					want := append([]byte(nil), b...)
					if size == 8 {
						// The eight-byte boundary must still perform the shift.
						if shift == -1 {
							want = bytes.Repeat([]byte{0x01}, 8)
						} else if shift == 1 {
							want = bytes.Repeat([]byte{0x04}, 8)
						}
					}
					vector.shift(b, shift)
					if !bytes.Equal(b, want) {
						t.Errorf("Shift(%d) = %x, want %x", shift, b, want)
					}
				})
			}
		}
	}
}
