package bitfield

import (
	"bytes"
	"reflect"
	"testing"
)

func TestBitvector16_Len(t *testing.T) {
	bvs := []Bitvector16{
		{},
		{0x01, 0x00},
		{0x01, 0x02},
		{0x0F, 0x0F},
	}
	for _, bv := range bvs {
		if bv.Len() != 16 {
			t.Errorf("(%x).Len() = %d, wanted %d", bv, bv.Len(), 16)
		}
	}
}

func TestBitvector16_BitAt(t *testing.T) {
	tests := []struct {
		bitlist Bitvector16
		idx     uint64
		want    bool
	}{
		{
			bitlist: Bitvector16{0x01, 0x23, 0xE2}, // wrong length
			idx:     0,
			want:    false,
		},
		{
			bitlist: Bitvector16{0xFF, 0xFF}, // 0b11111111 11111111
			idx:     16,                      // Out of bounds
			want:    false,
		},
		{
			bitlist: Bitvector16{0x01, 0x00}, // 0b00000001 00000000
			idx:     0,                       //          ^
			want:    true,
		},
		{
			bitlist: Bitvector16{},
			idx:     0,
			want:    false,
		},
		{
			bitlist: Bitvector16{0xFF, 0xFF, 0xFF},
			idx:     0,
			want:    false,
		},
		{
			bitlist: Bitvector16{0x0E, 0x00}, // 0b00001110 00000000
			idx:     0,                       //          ^
			want:    false,
		},
		{
			bitlist: Bitvector16{0x0E, 0x00}, // 0b00001110 00000000
			idx:     1,                       //         ^
			want:    true,
		},
		{
			bitlist: Bitvector16{0x0E, 0x00}, // 0b00001110 00000000
			idx:     3,                       //       ^
			want:    true,
		},
		{
			bitlist: Bitvector16{0x00, 0x01}, // 0b00000000 00000001
			idx:     8,                       //                  ^
			want:    true,
		},
		{
			bitlist: Bitvector16{0x00, 0x80}, // 0b00000000 10000000
			idx:     15,                      //           ^
			want:    true,
		},
		{
			bitlist: Bitvector16{0x9E, 0x00}, // 0b10011110 00000000
			idx:     7,                       //   ^
			want:    true,
		},
	}

	for _, tt := range tests {
		if tt.bitlist.BitAt(tt.idx) != tt.want {
			t.Errorf(
				"(%x).BitAt(%d) = %t, wanted %t",
				tt.bitlist,
				tt.idx,
				tt.bitlist.BitAt(tt.idx),
				tt.want,
			)
		}
	}
}

func TestBitvector16_SetBitAt(t *testing.T) {
	tests := []struct {
		bitvector Bitvector16
		idx       uint64
		val       bool
		want      Bitvector16
	}{
		{
			bitvector: Bitvector16{0x01, 0x00}, // 0b00000001 00000000
			idx:       0,                       //          ^
			val:       true,
			want:      Bitvector16{0x01, 0x00}, // 0b00000001 00000000
		},
		{
			bitvector: Bitvector16{0x02, 0x00}, // 0b00000010 00000000
			idx:       0,                       //          ^
			val:       true,
			want:      Bitvector16{0x03, 0x00}, // 0b00000011 00000000
		},
		{
			bitvector: Bitvector16{0x00, 0x00}, // 0b00000000 00000000
			idx:       1,                       //         ^
			val:       true,
			want:      Bitvector16{0x02, 0x00}, // 0b00000010 00000000
		},
		{
			bitvector: Bitvector16{0x00, 0x00}, // 0b00000000 00000000
			idx:       8,                       //                  ^
			val:       true,
			want:      Bitvector16{0x00, 0x01}, // 0b00000000 00000001
		},
		{
			bitvector: Bitvector16{0x00, 0x00}, // 0b00000000 00000000
			idx:       15,                      //           ^
			val:       true,
			want:      Bitvector16{0x00, 0x80}, // 0b00000000 10000000
		},
		{
			bitvector: Bitvector16{0xFF, 0xFF}, // 0b11111111 11111111
			idx:       9,                       //                 ^
			val:       false,
			want:      Bitvector16{0xFF, 0xFD}, // 0b11111111 11111101
		},
		{
			bitvector: Bitvector16{}, // wrong length
			idx:       5,
			val:       true,
			want:      Bitvector16{},
		},
		{
			bitvector: Bitvector16{0x0F, 0x00}, // 0b00001111 00000000
			idx:       0,                       //          ^
			val:       true,
			want:      Bitvector16{0x0F, 0x00}, // 0b00001111 00000000
		},
		{
			bitvector: Bitvector16{0x0F, 0x00}, // 0b00001111 00000000
			idx:       0,                       //          ^
			val:       false,
			want:      Bitvector16{0x0E, 0x00}, // 0b00001110 00000000
		},
		{
			bitvector: Bitvector16{0x00, 0x00}, // Out of bound
			idx:       16,
			val:       true,
			want:      Bitvector16{0x00, 0x00},
		},
	}

	for _, tt := range tests {
		original := make(Bitvector16, len(tt.bitvector))
		copy(original, tt.bitvector)

		tt.bitvector.SetBitAt(tt.idx, tt.val)
		if !bytes.Equal(tt.bitvector, tt.want) {
			t.Errorf(
				"(%x).SetBitAt(%d, %t) = %x, wanted %x",
				original,
				tt.idx,
				tt.val,
				tt.bitvector,
				tt.want,
			)
		}
	}
}

func TestBitvector16_Count(t *testing.T) {
	tests := []struct {
		bitvector Bitvector16
		want      uint64
	}{
		{
			bitvector: Bitvector16{},
			want:      0,
		},
		{
			bitvector: Bitvector16{0x01, 0x00}, // 0b00000001 00000000
			want:      1,
		},
		{
			bitvector: Bitvector16{0x03, 0x00}, // 0b00000011 00000000
			want:      2,
		},
		{
			bitvector: Bitvector16{0x07, 0x00}, // 0b00000111 00000000
			want:      3,
		},
		{
			bitvector: Bitvector16{0x0F, 0x0F}, // 0b00001111 00001111
			want:      8,
		},
		{
			bitvector: Bitvector16{0xFF, 0xFF}, // 0b11111111 11111111
			want:      16,
		},
		{
			bitvector: Bitvector16{0xF0, 0x0F}, // 0b11110000 00001111
			want:      8,
		},
		{
			bitvector: Bitvector16{0x00, 0x00, 0xFF}, // extra bytes are ignored
			want:      0,
		},
	}

	for _, tt := range tests {
		if tt.bitvector.Count() != tt.want {
			t.Errorf(
				"(%x).Count() = %d, wanted %d",
				tt.bitvector,
				tt.bitvector.Count(),
				tt.want,
			)
		}
	}
}

func TestBitvector16_Bytes(t *testing.T) {
	tests := []struct {
		bitvector Bitvector16
		want      []byte
	}{
		{
			bitvector: Bitvector16{},
			want:      []byte{},
		},
		{
			bitvector: Bitvector16{0x00, 0x00},
			want:      []byte{0x00, 0x00},
		},
		{
			bitvector: Bitvector16{0x01, 0x00},
			want:      []byte{0x01, 0x00},
		},
		{
			bitvector: Bitvector16{0x12, 0x34},
			want:      []byte{0x12, 0x34},
		},
		{
			bitvector: Bitvector16{0xFF, 0xFF},
			want:      []byte{0xFF, 0xFF},
		},
		{
			bitvector: Bitvector16{0xF0, 0x0F, 0xAA}, // extra bytes truncated
			want:      []byte{0xF0, 0x0F},
		},
		{
			bitvector: Bitvector16{0x01}, // shorter than expected, returned as-is
			want:      []byte{0x01},
		},
	}

	for _, tt := range tests {
		if !bytes.Equal(tt.bitvector.Bytes(), tt.want) {
			t.Errorf(
				"(%x).Bytes() = %x, wanted %x",
				tt.bitvector,
				tt.bitvector.Bytes(),
				tt.want,
			)
		}
	}
}

func TestBitvector16_BitIndices(t *testing.T) {
	tests := []struct {
		a    Bitvector16
		want []int
	}{
		{
			a:    Bitvector16{0b1001, 0b0},
			want: []int{0, 3},
		},
		{
			a:    Bitvector16{0b1000, 0b0},
			want: []int{3},
		},
		{
			a:    Bitvector16{0b10, 0b1},
			want: []int{1, 8},
		},
		{
			a:    Bitvector16{0b11111111, 0b11111111},
			want: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		},
		{
			a:    Bitvector16{0b0, 0b00000011},
			want: []int{8, 9},
		},
		{
			a:    Bitvector16{0b0, 0b10000000, 0b1}, // extra bytes ignored
			want: []int{15},
		},
	}

	for _, tt := range tests {
		if !reflect.DeepEqual(tt.a.BitIndices(), tt.want) {
			t.Errorf(
				"(%0.16b).BitIndices() = %x, wanted %x",
				tt.a,
				tt.a.BitIndices(),
				tt.want,
			)
		}
	}
}

func TestBitvector16_Contains(t *testing.T) {
	tests := []struct {
		a    Bitvector16
		b    Bitvector16
		want bool
	}{
		{
			a:    Bitvector16{0x02, 0x00}, // 0b00000010 00000000
			b:    Bitvector16{0x03, 0x00}, // 0b00000011 00000000
			want: false,
		},
		{
			a:    Bitvector16{0x03, 0x00}, // 0b00000011 00000000
			b:    Bitvector16{0x03, 0x00}, // 0b00000011 00000000
			want: true,
		},
		{
			a:    Bitvector16{0x13, 0xAA}, // 0b00010011 10101010
			b:    Bitvector16{0x15, 0xAA}, // 0b00010101 10101010
			want: false,
		},
		{
			a:    Bitvector16{0x1F, 0xFF}, // 0b00011111 11111111
			b:    Bitvector16{0x13, 0xAA}, // 0b00010011 10101010
			want: true,
		},
		{
			a:    Bitvector16{0xFF, 0x0F}, // 0b11111111 00001111
			b:    Bitvector16{0xFF, 0x10}, // 0b11111111 00010000 — second byte differs
			want: false,
		},
	}

	for _, tt := range tests {
		if got, err := tt.a.Contains(tt.b); got != tt.want || err != nil {
			t.Errorf(
				"(%x).Contains(%x) = %t, %v, wanted %t",
				tt.a,
				tt.b,
				got,
				err,
				tt.want,
			)
		}
	}
}

func TestBitvector16_Overlaps(t *testing.T) {
	tests := []struct {
		a    Bitvector16
		b    Bitvector16
		want bool
	}{
		{
			a:    Bitvector16{0x06, 0x00}, // 0b00000110 00000000
			b:    Bitvector16{0x01, 0x00}, // 0b00000001 00000000
			want: false,
		},
		{
			a:    Bitvector16{0x06, 0x00}, // 0b00000110 00000000
			b:    Bitvector16{0x05, 0x00}, // 0b00000101 00000000
			want: true,
		},
		{
			a:    Bitvector16{0x1A, 0x00}, // 0b00011010 00000000
			b:    Bitvector16{0x25, 0x00}, // 0b00100101 00000000
			want: false,
		},
		{
			a:    Bitvector16{0x1F, 0x00}, // 0b00011111 00000000
			b:    Bitvector16{0x11, 0x00}, // 0b00010001 00000000
			want: true,
		},
		{
			a:    Bitvector16{0x00, 0x80}, // 0b00000000 10000000
			b:    Bitvector16{0x00, 0x80}, // 0b00000000 10000000 — overlap on second byte
			want: true,
		},
		{
			a:    Bitvector16{0xFF, 0x00}, // 0b11111111 00000000
			b:    Bitvector16{0x00, 0xFF}, // 0b00000000 11111111
			want: false,
		},
	}

	for _, tt := range tests {
		if got, err := tt.a.Overlaps(tt.b); got != tt.want || err != nil {
			t.Errorf(
				"(%x).Overlaps(%x) = %t, %v, wanted %t",
				tt.a,
				tt.b,
				got,
				err,
				tt.want,
			)
		}
	}
}

func TestBitVector16_Or(t *testing.T) {
	tests := []struct {
		a    Bitvector16
		b    Bitvector16
		want Bitvector16
	}{
		{
			a:    Bitvector16{0x02, 0x00}, // 0b00000010 00000000
			b:    Bitvector16{0x03, 0x00}, // 0b00000011 00000000
			want: Bitvector16{0x03, 0x00}, // 0b00000011 00000000
		},
		{
			a:    Bitvector16{0x03, 0x00}, // 0b00000011 00000000
			b:    Bitvector16{0x03, 0x00}, // 0b00000011 00000000
			want: Bitvector16{0x03, 0x00}, // 0b00000011 00000000
		},
		{
			a:    Bitvector16{0x13, 0xA0}, // 0b00010011 10100000
			b:    Bitvector16{0x15, 0x0A}, // 0b00010101 00001010
			want: Bitvector16{0x17, 0xAA}, // 0b00010111 10101010
		},
		{
			a:    Bitvector16{0x1F, 0xF0}, // 0b00011111 11110000
			b:    Bitvector16{0x13, 0x0F}, // 0b00010011 00001111
			want: Bitvector16{0x1F, 0xFF}, // 0b00011111 11111111
		},
	}

	for _, tt := range tests {
		if got, err := tt.a.Or(tt.b); !bytes.Equal(got, tt.want) || err != nil {
			t.Errorf(
				"(%x).Or(%x) = %x, %v, wanted %x",
				tt.a,
				tt.b,
				got,
				err,
				tt.want,
			)
		}
	}
}
