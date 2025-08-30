package bitfield

import (
	"bytes"
	"reflect"
	"testing"
)

func TestBitvector2_Len(t *testing.T) {
	bvs := []Bitvector2{
		{},
		{0x01},
		{0x02},
		{0x03},
	}

	for _, bv := range bvs {
		if bv.Len() != 2 {
			t.Errorf("(%x).Len() = %d, wanted %d", bv, bv.Len(), 2)
		}
	}
}

func TestBitvector2_BitAt(t *testing.T) {
	tests := []struct {
		bitlist Bitvector2
		idx     uint64
		want    bool
	}{
		{
			bitlist: Bitvector2{0x01}, // 0b00000001
			idx:     55,               // Out of bounds
			want:    false,
		},
		{
			bitlist: Bitvector2{0x01}, // 0b00000001
			idx:     0,                //          ^
			want:    true,
		},
		{
			bitlist: Bitvector2{},
			idx:     0,
			want:    false,
		},
		{
			bitlist: Bitvector2{0xFF, 0xFF},
			idx:     0,
			want:    false,
		},
		{
			bitlist: Bitvector2{0x02}, // 0b00000010
			idx:     0,                //          ^
			want:    false,
		},
		{
			bitlist: Bitvector2{0x02}, // 0b00000010
			idx:     1,                //         ^
			want:    true,
		},
		{
			bitlist: Bitvector2{0x03}, // 0b00000011
			idx:     0,                //          ^
			want:    true,
		},
		{
			bitlist: Bitvector2{0x03}, // 0b00000011
			idx:     1,                //         ^
			want:    true,
		},
		{
			bitlist: Bitvector2{0x03}, // 0b00000011
			idx:     2,                // Out of bounds
			want:    false,
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

func TestBitvector2_SetBitAt(t *testing.T) {
	tests := []struct {
		bitvector Bitvector2
		idx       uint64
		val       bool
		want      Bitvector2
	}{
		{
			bitvector: Bitvector2{0x01}, // 0b00000001
			idx:       0,                //          ^
			val:       true,
			want:      Bitvector2{0x01}, // 0b00000001
		},
		{
			bitvector: Bitvector2{0x02}, // 0b00000010
			idx:       0,                //          ^
			val:       true,
			want:      Bitvector2{0x03}, // 0b00000011
		},
		{
			bitvector: Bitvector2{0x00}, // 0b00000000
			idx:       1,                //         ^
			val:       true,
			want:      Bitvector2{0x02}, // 0b00000010
		},
		{
			bitvector: Bitvector2{}, // 0b00000000
			idx:       1,            //     ^
			val:       true,
			want:      Bitvector2{}, // 0b00000000
		},
		{
			bitvector: Bitvector2{0x00}, // 0b00000000
			idx:       2,                // Out of bounds
			val:       true,
			want:      Bitvector2{0x00}, // 0b00000000
		},
		{
			bitvector: Bitvector2{0x03}, // 0b00000011
			idx:       0,                //          ^
			val:       true,
			want:      Bitvector2{0x03}, // 0b00000011
		},
		{
			bitvector: Bitvector2{0x03}, // 0b00000011
			idx:       0,                //          ^
			val:       false,
			want:      Bitvector2{0x02}, // 0b00000010
		},
		{
			bitvector: Bitvector2{0x03}, // 0b00000011
			idx:       1,                //         ^
			val:       false,
			want:      Bitvector2{0x01}, // 0b00000001
		},
	}

	for _, tt := range tests {
		original := make(Bitvector2, len(tt.bitvector))
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

func TestBitvector2_Count(t *testing.T) {
	tests := []struct {
		bitvector Bitvector2
		want      uint64
	}{
		{
			bitvector: Bitvector2{},
			want:      0,
		},
		{
			bitvector: Bitvector2{0x01}, // 0b00000001
			want:      1,
		},
		{
			bitvector: Bitvector2{0x02}, // 0b00000010
			want:      1,
		},
		{
			bitvector: Bitvector2{0x03}, // 0b00000011
			want:      2,
		},
		{
			bitvector: Bitvector2{0xFF}, // 0b11111111
			want:      2,
		},
		{
			bitvector: Bitvector2{0xFC}, // 0b11111100
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

func TestBitvector2_Bytes(t *testing.T) {
	tests := []struct {
		bitvector Bitvector2
		want      []byte
	}{
		{
			bitvector: Bitvector2{},
			want:      []byte{},
		},
		{
			bitvector: Bitvector2{0x00}, // 0b00000000
			want:      []byte{0x00},     // 0b00000000
		},
		{
			bitvector: Bitvector2{0x01}, // 0b00000001
			want:      []byte{0x01},     // 0b00000001
		},
		{
			bitvector: Bitvector2{0x02}, // 0b00000010
			want:      []byte{0x02},     // 0b00000010
		},
		{
			bitvector: Bitvector2{0x03}, // 0b00000011
			want:      []byte{0x03},     // 0b00000011
		},
		{
			bitvector: Bitvector2{0xFF}, // 0b11111111
			want:      []byte{0x03},     // 0b00000011
		},
		{
			bitvector: Bitvector2{0xFC}, // 0b11111100
			want:      []byte{0x00},     // 0b00000000
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

func TestBitvector2_Shift(t *testing.T) {
	tests := []struct {
		bitvector Bitvector2
		shift     int
		want      Bitvector2
	}{
		{
			bitvector: Bitvector2{},
			shift:     1,
			want:      Bitvector2{},
		},
		{
			bitvector: Bitvector2{0x01},
			shift:     1,
			want:      Bitvector2{0x02},
		},
		{
			bitvector: Bitvector2{0x02},
			shift:     1,
			want:      Bitvector2{0x00},
		},
		{
			bitvector: Bitvector2{0x03},
			shift:     1,
			want:      Bitvector2{0x02},
		},
		{
			bitvector: Bitvector2{0x02},
			shift:     -1,
			want:      Bitvector2{0x01},
		},
		{
			bitvector: Bitvector2{0x03},
			shift:     -1,
			want:      Bitvector2{0x01},
		},
		{
			bitvector: Bitvector2{0x03},
			shift:     2,
			want:      Bitvector2{0x00},
		},
		{
			bitvector: Bitvector2{0x03},
			shift:     -2,
			want:      Bitvector2{0x00},
		},
		{
			bitvector: Bitvector2{0x03},
			shift:     8,
			want:      Bitvector2{0x00},
		},
		{
			bitvector: Bitvector2{0x03},
			shift:     -256,
			want:      Bitvector2{0x00},
		},
	}

	for _, tt := range tests {
		original := make(Bitvector2, len(tt.bitvector))
		copy(original, tt.bitvector)

		tt.bitvector.Shift(tt.shift)
		if !bytes.Equal(tt.bitvector, tt.want) {
			t.Errorf(
				"(%x).Shift(%d) = %x, wanted %x",
				original,
				tt.shift,
				tt.bitvector,
				tt.want,
			)
		}
	}
}

func TestBitvector2_BitIndices(t *testing.T) {
	tests := []struct {
		a    Bitvector2
		want []int
	}{
		{
			a:    Bitvector2{0b01},
			want: []int{0},
		},
		{
			a:    Bitvector2{0b10},
			want: []int{1},
		},
		{
			a:    Bitvector2{0b11},
			want: []int{0, 1},
		},
		{
			a:    Bitvector2{0b00},
			want: []int{},
		},
	}

	for _, tt := range tests {
		if !reflect.DeepEqual(tt.a.BitIndices(), tt.want) {
			t.Errorf(
				"(%0.8b).BitIndices() = %x, wanted %x",
				tt.a,
				tt.a.BitIndices(),
				tt.want,
			)
		}
	}
}