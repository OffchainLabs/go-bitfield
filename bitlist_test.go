package bitfield

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func TestNewBitlist(t *testing.T) {
	tests := []struct {
		size uint64
		want Bitlist
	}{
		{
			size: 0,
			want: Bitlist{0x01},
		},
		{
			size: 1,
			want: Bitlist{0x02},
		},
		{
			size: 2,
			want: Bitlist{0x04},
		},
		{
			size: 3,
			want: Bitlist{0x08},
		},
		{
			size: 8,
			want: Bitlist{0x00, 0x01},
		},
		{
			size: 9,
			want: Bitlist{0x00, 0x02},
		},
	}

	for _, tt := range tests {
		got := NewBitlist(tt.size)
		if !bytes.Equal(got, tt.want) {
			t.Errorf(
				"NewBitlist(%d) = %x, wanted %x",
				tt.size,
				got,
				tt.want,
			)
		}
	}
}

func TestBitlist_Len(t *testing.T) {
	tests := []struct {
		bitlist Bitlist
		want    uint64
	}{
		{
			bitlist: Bitlist{},
			want:    0,
		},
		{
			bitlist: Bitlist{0x00}, // 0b00000000, invalid list
			want:    0,
		},
		{
			bitlist: Bitlist{0x01}, // 0b00000001
			want:    0,
		},
		{
			bitlist: Bitlist{0x02}, // 0b00000010
			want:    1,
		},
		{
			bitlist: Bitlist{0x08}, // 0b00001000
			want:    3,
		},
		{
			bitlist: Bitlist{0x0E}, // 0b00001110
			want:    3,
		},
		{
			bitlist: Bitlist{0x0F}, // 0b00001111
			want:    3,
		},
		{
			bitlist: Bitlist{0x10}, // 0b00010000
			want:    4,
		},
		{
			bitlist: Bitlist{0x00, 0x01}, // 0b00000000, 0b00000001
			want:    8,
		},
		{
			bitlist: Bitlist{0x00, 0x02}, // 0b00000000, 0b00000010
			want:    9,
		},
		{
			bitlist: Bitlist{0x00, 0x02, 0x08}, // 0b00000000, 0b00000010, 0b00001000
			want:    19,
		},
	}

	for _, tt := range tests {
		if tt.bitlist.Len() != tt.want {
			t.Errorf("(%x).Len() = %d, wanted %d", tt.bitlist, tt.bitlist.Len(), tt.want)
		}
	}
}

func TestBitlist_BitAt(t *testing.T) {
	tests := []struct {
		bitlist Bitlist
		idx     uint64
		want    bool
	}{
		{
			bitlist: Bitlist{},
			idx:     0,
			want:    false,
		},
		{
			bitlist: Bitlist{0x01}, // 0b00000001
			idx:     55,            // Out of bounds
			want:    false,
		},
		{
			bitlist: Bitlist{0x01}, // 0b00000001
			idx:     0,             //          ^ (length bit)
			want:    false,
		},
		{
			bitlist: Bitlist{0x0E}, // 0b00001110
			idx:     0,             //          ^
			want:    false,
		},
		{
			bitlist: Bitlist{0x0E}, // 0b00001110
			idx:     1,             //         ^
			want:    true,
		},
		{
			bitlist: Bitlist{0x0E}, // 0b00001110
			idx:     3,             //       ^
			want:    false,
		},
		{
			bitlist: Bitlist{0x0E}, // 0b00001110
			idx:     4,             //       ^ (length bit)
			want:    false,
		},
		{
			bitlist: Bitlist{0xFF, 0x0F}, // 0b11111111, 0b00001111
			idx:     4,                   //      ^
			want:    true,
		},
		{
			bitlist: Bitlist{0xFF, 0x0F}, // 0b11111111, 0b00001111
			idx:     12,                  //                  ^
			want:    false,
		},
		{
			bitlist: Bitlist{0xFF, 0x0F}, // 0b11111111, 0b00001111
			idx:     11,                  //                   ^ (length bit)
			want:    false,
		},
		{
			bitlist: Bitlist{0x00, 0x0F}, // 0b00000000, 0b00001111
			idx:     10,                  //                    ^
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

func TestBitlist_SetBitAt(t *testing.T) {
	tests := []struct {
		bitlist Bitlist
		idx     uint64
		val     bool
		want    Bitlist
	}{
		{
			bitlist: Bitlist{},
			idx:     0,
			val:     true,
			want:    Bitlist{},
		},
		{
			bitlist: Bitlist{0x01}, // 0b00000001
			idx:     0,             //          ^
			val:     true,
			want:    Bitlist{0x01}, // 0b00000001
		},
		{
			bitlist: Bitlist{0x02}, // 0b00000010
			idx:     0,             //          ^
			val:     true,
			want:    Bitlist{0x03}, // 0b00000011
		},
		{
			bitlist: Bitlist{0x10}, // 0b00010000
			idx:     0,             //          ^
			val:     true,
			want:    Bitlist{0x11}, // 0b00010001
		},
		{
			bitlist: Bitlist{0x10}, // 0b00010000
			idx:     0,             //          ^
			val:     true,
			want:    Bitlist{0x11}, // 0b00010001
		},
		{
			bitlist: Bitlist{0x10}, // 0b00010000
			idx:     64,            // Out of bounds
			val:     true,
			want:    Bitlist{0x10}, // 0b00010001
		},
		{
			bitlist: Bitlist{0x1F}, // 0b00011111
			idx:     0,             //          ^
			val:     true,
			want:    Bitlist{0x1F}, // 0b00011111
		},
		{
			bitlist: Bitlist{0x1F}, // 0b00011111
			idx:     1,             //         ^
			val:     false,
			want:    Bitlist{0x1D}, // 0b00011101
		},
		{
			bitlist: Bitlist{0x1F}, // 0b00011111
			idx:     4,             //      ^ (length bit)
			val:     false,
			want:    Bitlist{0x1F}, // 0b00011111
		},
		{
			bitlist: Bitlist{0x1F}, // 0b00011111
			idx:     64,            // Out of bounds
			val:     false,
			want:    Bitlist{0x1F}, // 0b00011111
		},
		{
			bitlist: Bitlist{0x1F, 0x01}, // 0b00011111, 0b00000001
			idx:     0,                   //          ^
			val:     true,
			want:    Bitlist{0x1F, 0x01}, // 0b00011111, 0b00000001
		},
		{
			bitlist: Bitlist{0x1F, 0x01}, // 0b00011111, 0b00000001
			idx:     0,                   //          ^
			val:     false,
			want:    Bitlist{0x1E, 0x01}, // 0b00011110, 0b00000001
		},
		{
			bitlist: Bitlist{0x00, 0x10}, // 0b00000000, 0b00010000
			idx:     8,                   //                      ^
			val:     true,
			want:    Bitlist{0x00, 0x11}, // 0b00000000, 0b00010001
		},
		{
			bitlist: Bitlist{0x00, 0x11}, // 0b00000000, 0b00010001
			idx:     8,                   //                      ^
			val:     false,
			want:    Bitlist{0x00, 0x10}, // 0b00000000, 0b00010000
		},
	}

	for _, tt := range tests {
		original := make(Bitlist, len(tt.bitlist))
		copy(original, tt.bitlist)

		tt.bitlist.SetBitAt(tt.idx, tt.val)
		if !bytes.Equal(tt.bitlist, tt.want) {
			t.Errorf(
				"(%x).SetBitAt(%d, %t) = %x, wanted %x",
				original,
				tt.idx,
				tt.val,
				tt.bitlist,
				tt.want,
			)
		}
	}
}

func TestBitlist_Bytes(t *testing.T) {
	tests := []struct {
		bitlist Bitlist
		want    []byte
	}{
		{
			bitlist: Bitlist{},
			want:    []byte{},
		},
		{
			bitlist: Bitlist{0x00},
			want:    []byte{},
		},
		{
			bitlist: Bitlist{0x01},
			want:    []byte{},
		},
		{
			bitlist: Bitlist{0x02},
			want:    []byte{},
		},
		{
			bitlist: Bitlist{0x03},
			want:    []byte{0x01},
		},
		{
			bitlist: Bitlist{0x12},
			want:    []byte{0x02},
		},
		{
			bitlist: Bitlist{0x02, 0x01},
			want:    []byte{0x02},
		},
		{
			bitlist: Bitlist{0x02, 0x02},
			want:    []byte{0x02},
		},
		{
			bitlist: Bitlist{0x02, 0x01},
			want:    []byte{0x02},
		},
		{
			bitlist: Bitlist{0x02, 0x03},
			want:    []byte{0x02, 0x01},
		},
		{
			bitlist: Bitlist{0x01, 0x00, 0x08},
			want:    []byte{0x01},
		}, {
			bitlist: Bitlist{0x00, 0x00, 0x02},
			want:    []byte{},
		},
	}

	for _, tt := range tests {
		got := tt.bitlist.Bytes()
		if !bytes.Equal(got, tt.want) {
			t.Errorf(
				"(%x).Bytes() = %x, wanted %x",
				tt.bitlist,
				got,
				tt.want,
			)
		}
	}
}

func TestBitlist_BytesNoTrim(t *testing.T) {
	tests := []struct {
		bitlist Bitlist
		want    []byte
	}{
		{
			bitlist: Bitlist{},
			want:    []byte{},
		},
		{
			bitlist: Bitlist{0x00},
			want:    []byte{},
		},
		{
			bitlist: Bitlist{0x01},
			want:    []byte{},
		},
		{
			bitlist: Bitlist{0x02},
			want:    []byte{0x00},
		},
		{
			bitlist: Bitlist{0x03},
			want:    []byte{0x01},
		},
		{
			bitlist: Bitlist{0x12},
			want:    []byte{0x02},
		},
		{
			bitlist: Bitlist{0x02, 0x01},
			want:    []byte{0x02},
		},
		{
			bitlist: Bitlist{0x02, 0x02},
			want:    []byte{0x02, 0x00},
		},
		{
			bitlist: Bitlist{0x02, 0x03},
			want:    []byte{0x02, 0x01},
		},
		{
			bitlist: Bitlist{0x01, 0x00, 0x08},
			want:    []byte{0x01, 0x00, 0x00},
		},
		{
			bitlist: Bitlist{0x00, 0x00, 0x02},
			want:    []byte{0x00, 0x00, 0x00},
		},
		{
			bitlist: Bitlist{0x00, 0x00, 0x01},
			want:    []byte{0x00, 0x00},
		},
	}

	for _, tt := range tests {
		got := tt.bitlist.BytesNoTrim()
		if !bytes.Equal(got, tt.want) {
			t.Errorf("(%#x).BytesNoTrim() = %#v, wanted %#v", tt.bitlist, got, tt.want)
		}
	}
}

func TestBitlist_ToBitlist64(t *testing.T) {
	tests := []struct {
		size            uint64
		selectedIndices []uint64
	}{
		{
			size:            0,
			selectedIndices: []uint64{},
		},
		{
			size:            1,
			selectedIndices: []uint64{0},
		},
		{
			size:            2,
			selectedIndices: []uint64{0, 1},
		},
		{
			size:            7,
			selectedIndices: []uint64{0, 1, 6},
		},
		{
			size:            8,
			selectedIndices: []uint64{0, 1, 6, 7},
		},
		{
			size:            9,
			selectedIndices: []uint64{3, 4},
		},
		{
			size:            60,
			selectedIndices: []uint64{0, 2, 50},
		},
		{
			size:            64,
			selectedIndices: []uint64{0, 2, 63},
		},
		{
			size:            69,
			selectedIndices: []uint64{0, 2, 63, 67},
		},
		{
			size:            128,
			selectedIndices: []uint64{0, 2, 63, 67, 120},
		},
		{
			size:            128,
			selectedIndices: []uint64{0, 2, 63, 67, 90, 100, 120, 126, 127},
		},
		{
			size:            192,
			selectedIndices: []uint64{0, 2, 63, 67, 90, 100, 120, 126, 127, 150, 170},
		},
	}

	selectIndices := func(b Bitfield, indices []uint64) Bitfield {
		for _, idx := range indices {
			b.SetBitAt(idx, true)
		}
		return b
	}
	createBitlist64 := func(n uint64, indices []uint64) *Bitlist64 {
		return (selectIndices(NewBitlist64(n), indices)).(*Bitlist64)
	}
	createBitlist := func(n uint64, indices []uint64) Bitlist {
		return (selectIndices(NewBitlist(n), indices)).(Bitlist)
	}

	for _, tt := range tests {
		source := createBitlist(tt.size, tt.selectedIndices)
		wanted := createBitlist64(tt.size, tt.selectedIndices)
		t.Run(fmt.Sprintf("size:%d,indices:%v", tt.size, tt.selectedIndices), func(t *testing.T) {
			// Convert to Bitlist64.
			got, err := source.ToBitlist64()
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(got, wanted) {
				t.Errorf("ToBitlist64(%#x) = %#b, wanted %#b", source, got, wanted)
			}

			// Now convert back, and compare to the original.
			gotSource := got.ToBitlist()
			if !reflect.DeepEqual(gotSource, source) {
				t.Errorf("ToBitlist64(%#x).ToBitlist() = %+v, wanted %+v", source, gotSource, source)
			}

			// Make sure that both Bitlist and Bitlist64 Bytes() are equal.
			if !bytes.Equal(source.Bytes(), got.Bytes()) {
				t.Errorf("original.Bytes() != converted.Bytes() (%#x != %#x)", source.Bytes(), got.Bytes())
			}
		})
	}
}

func TestBitlist_Count(t *testing.T) {
	tests := []struct {
		bitlist Bitlist
		want    uint64
	}{
		{
			bitlist: Bitlist{},
			want:    0,
		},
		{
			bitlist: Bitlist{0x00}, // 0b00000000, invalid list
			want:    0,
		},
		{
			bitlist: Bitlist{0x01}, // 0b00000001
			want:    0,
		},
		{
			bitlist: Bitlist{0x03}, // 0b00000011
			want:    1,
		},
		{
			bitlist: Bitlist{0x0F}, // 0b00001111
			want:    3,
		},
		{
			bitlist: Bitlist{0x0F, 0x01}, // 0b00001111, 0b00000001
			want:    4,
		},
		{
			bitlist: Bitlist{0x0F, 0x03}, // 0b00001111, 0b00000011
			want:    5,
		},
	}

	for _, tt := range tests {
		if tt.bitlist.Count() != tt.want {
			t.Errorf(
				"(%x).Count() = %d, wanted %d",
				tt.bitlist,
				tt.bitlist.Count(),
				tt.want,
			)
		}
	}
}

func TestBitlist_Contains(t *testing.T) {
	tests := []struct {
		a    Bitlist
		b    Bitlist
		want bool
	}{
		{
			a:    Bitlist{0x02}, // 0b00000010
			b:    Bitlist{0x03}, // 0b00000011
			want: false,
		},
		{
			a:    Bitlist{0x03}, // 0b00000011
			b:    Bitlist{0x03}, // 0b00000011
			want: true,
		},
		{
			a:    Bitlist{0x13}, // 0b00010011
			b:    Bitlist{0x15}, // 0b00010101
			want: false,
		},
		{
			a:    Bitlist{0x1F}, // 0b00011111
			b:    Bitlist{0x13}, // 0b00010011
			want: true,
		},
		{
			a:    Bitlist{0x1F}, // 0b00011111
			b:    Bitlist{0x13}, // 0b00010011
			want: true,
		},
		{
			a:    Bitlist{0x1F, 0x03}, // 0b00011111, 0b00000011
			b:    Bitlist{0x13, 0x02}, // 0b00010011, 0b00000010
			want: true,
		},
		{
			a:    Bitlist{0x1F, 0x01}, // 0b00011111, 0b00000001
			b:    Bitlist{0x93, 0x01}, // 0b10010011, 0b00000001
			want: false,
		},
		{
			a:    Bitlist{0xFF, 0x02}, // 0b11111111, 0x00000010
			b:    Bitlist{0x13, 0x03}, // 0b00010011, 0x00000011
			want: false,
		},
		{
			a:    Bitlist{0xFF, 0x85}, // 0b11111111, 0x10000111
			b:    Bitlist{0x13, 0x8F}, // 0b00010011, 0x10001111
			want: false,
		},
		{
			a:    Bitlist{0xFF, 0x8F}, // 0b11111111, 0x10001111
			b:    Bitlist{0x13, 0x83}, // 0b00010011, 0x10000011
			want: true,
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

func TestBitlist_Overlaps(t *testing.T) {
	tests := []struct {
		a    Bitlist
		b    Bitlist
		want bool
	}{
		{
			a:    Bitlist{0x06}, // 0b00000110
			b:    Bitlist{0x05}, // 0b00000101
			want: false,
		},
		{
			a:    Bitlist{0x32}, // 0b00110010
			b:    Bitlist{0x21}, // 0b00100001
			want: false,
		},
		{
			a:    Bitlist{0x41}, // 0b00100001
			b:    Bitlist{0x40}, // 0b00100000
			want: false,
		},
		{
			a:    Bitlist{0x1F}, // 0b00011111
			b:    Bitlist{0x11}, // 0b00010001
			want: true,
		},
		{
			a:    Bitlist{0xFF, 0x85}, // 0b11111111, 0b10000111
			b:    Bitlist{0x13, 0x8F}, // 0b00010011, 0b10001111
			want: true,
		},
		{
			a:    Bitlist{0x01, 0x40}, // 0b00000001, 0b01000000
			b:    Bitlist{0x00, 0x40}, // 0b00000010, 0b01000000
			want: false,
		},
		{
			a:    Bitlist{0x01, 0x40}, // 0b00000001, 0b01000000
			b:    Bitlist{0x00, 0x40}, // 0b00000010, 0b01000000
			want: false,
		},
		{
			a:    Bitlist{0x01, 0x40}, // 0b00000001, 0b01000000
			b:    Bitlist{0x02, 0x40}, // 0b00000010, 0b01000000
			want: false,
		},
		{
			a:    Bitlist{0x01, 0x40}, // 0b00000001, 0b01000000
			b:    Bitlist{0x03, 0x40}, // 0b00000011, 0b01000000
			want: true,
		},
		{
			a:    Bitlist{0x01, 0x01, 0x01}, // 0b00000001, 0b00000001, 0b00000001
			b:    Bitlist{0x02, 0x00, 0x01}, // 0b00000010, 0b00000000, 0b00000001
			want: false,
		},
	}

	for _, tt := range tests {
		if result, err := tt.a.Overlaps(tt.b); result != tt.want || err != nil {
			t.Errorf(
				"(%x).Overlaps(%x) = %t, %v, wanted %t",
				tt.a,
				tt.b,
				result,
				err,
				tt.want,
			)
		}
	}
}

func TestBitlist_Or(t *testing.T) {
	tests := []struct {
		a    Bitlist
		b    Bitlist
		want Bitlist
	}{
		{
			a:    Bitlist{0x02}, // 0b00000010
			b:    Bitlist{0x03}, // 0b00000011
			want: Bitlist{0x03}, // 0b00000011
		},
		{
			a:    Bitlist{0x03}, // 0b00000011
			b:    Bitlist{0x03}, // 0b00000011
			want: Bitlist{0x03}, // 0b00000011
		},
		{
			a:    Bitlist{0x13}, // 0b00010011
			b:    Bitlist{0x15}, // 0b00010101
			want: Bitlist{0x17}, // 0b00010111
		},
		{
			a:    Bitlist{0x1F}, // 0b00011111
			b:    Bitlist{0x13}, // 0b00010011
			want: Bitlist{0x1F}, // 0b00011111
		},
		{
			a:    Bitlist{0x1F, 0x03}, // 0b00011111, 0b00000011
			b:    Bitlist{0x13, 0x02}, // 0b00010011, 0b00000010
			want: Bitlist{0x1F, 0x03}, // 0b00011111, 0b00000011
		},
		{
			a:    Bitlist{0x1F, 0x01}, // 0b00011111, 0b00000001
			b:    Bitlist{0x93, 0x01}, // 0b10010011, 0b00000001
			want: Bitlist{0x9F, 0x01}, // 0b00011111, 0b00000001
		},
		{
			a:    Bitlist{0xFF, 0x02}, // 0b11111111, 0x00000010
			b:    Bitlist{0x13, 0x03}, // 0b00010011, 0x00000011
			want: Bitlist{0xFF, 0x03}, // 0b11111111, 0x00000011
		},
		{
			a:    Bitlist{0xFF, 0x85}, // 0b11111111, 0x10000111
			b:    Bitlist{0x13, 0x8F}, // 0b00010011, 0x10001111
			want: Bitlist{0xFF, 0x8F}, // 0b11111111, 0x10001111
		},
	}

	for _, tt := range tests {
		t.Run("Or()", func(t *testing.T) {
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
		})

		t.Run("NoAllocOr()", func(t *testing.T) {
			for _, tt := range tests {
				res := Bitlist(bytes.Clone(tt.a))
				// Make sure that no existing bits set interfere with operation. This is done to simulate
				// the case when res variable is already populated from the previous run.
				for i := uint64(0); i < res.Len(); i += 10 {
					res.SetBitAt(i, true)
				}
				tt.a.NoAllocOr(tt.b, res)
				if !bytes.Equal(res, tt.want) {
					t.Errorf("(%+v).NoAllocOr(%+v) = %+v, wanted %x", tt.a, tt.b, res, tt.want)
				}
			}
		})
	}
}

func TestBitlist_And(t *testing.T) {
	tests := []struct {
		a    Bitlist
		b    Bitlist
		want Bitlist
	}{
		{
			a:    Bitlist{0x02}, // 0b00000010
			b:    Bitlist{0x03}, // 0b00000011
			want: Bitlist{0x02}, // 0b00000010
		},
		{
			a:    Bitlist{0x03}, // 0b00000011
			b:    Bitlist{0x03}, // 0b00000011
			want: Bitlist{0x03}, // 0b00000011
		},
		{
			a:    Bitlist{0x13}, // 0b00010011
			b:    Bitlist{0x15}, // 0b00010101
			want: Bitlist{0x11}, // 0b00010001
		},
		{
			a:    Bitlist{0x1F}, // 0b00011111
			b:    Bitlist{0x13}, // 0b00010011
			want: Bitlist{0x13}, // 0b00010011
		},
		{
			a:    Bitlist{0x1F, 0x03}, // 0b00011111, 0b00000011
			b:    Bitlist{0x13, 0x02}, // 0b00010011, 0b00000010
			want: Bitlist{0x13, 0x02}, // 0b00010011, 0b00000010
		},
		{
			a:    Bitlist{0x9F, 0x01}, // 0b10011111, 0b00000001
			b:    Bitlist{0x93, 0x01}, // 0b10010011, 0b00000001
			want: Bitlist{0x93, 0x01}, // 0b10010011, 0b00000001
		},
		{
			a:    Bitlist{0xFF, 0x02}, // 0b11111111, 0x00000010
			b:    Bitlist{0x13, 0x03}, // 0b00010011, 0x00000011
			want: Bitlist{0x13, 0x02}, // 0b00010011, 0x00000010
		},
		{
			a:    Bitlist{0xFF, 0x87}, // 0b11111111, 0x10000111
			b:    Bitlist{0x13, 0x8F}, // 0b00010011, 0x10001111
			want: Bitlist{0x13, 0x87}, // 0b00010011, 0x10000111
		},
	}

	for _, tt := range tests {
		if got, err := tt.a.And(tt.b); !bytes.Equal(got, tt.want) || err != nil {
			t.Errorf(
				"(%x).And(%x) = %x, %v, wanted %x",
				tt.a,
				tt.b,
				got,
				err,
				tt.want,
			)
		}
	}
}

func TestBitlist_Xor(t *testing.T) {
	tests := []struct {
		a    Bitlist
		b    Bitlist
		want Bitlist
	}{
		{
			a:    Bitlist{0x02}, // 0b00000010
			b:    Bitlist{0x03}, // 0b00000011
			want: Bitlist{0x03}, // 0b00000011
		},
		{
			a:    Bitlist{0x03}, // 0b00000011
			b:    Bitlist{0x03}, // 0b00000011
			want: Bitlist{0x02}, // 0b00000010
		},
		{
			a:    Bitlist{0x13}, // 0b00010011
			b:    Bitlist{0x15}, // 0b00010101
			want: Bitlist{0x16}, // 0b00010110
		},
		{
			a:    Bitlist{0x1F}, // 0b00011111
			b:    Bitlist{0x13}, // 0b00010011
			want: Bitlist{0x1c}, // 0b00011100
		},
		{
			a:    Bitlist{0x1F, 0x03}, // 0b00011111, 0b00000011
			b:    Bitlist{0x13, 0x02}, // 0b00010011, 0b00000010
			want: Bitlist{0x0c, 0x03}, // 0b00001100, 0b00000011
		},
		{
			a:    Bitlist{0x9F, 0x01}, // 0b10011111, 0b00000001
			b:    Bitlist{0x93, 0x01}, // 0b10010011, 0b00000001
			want: Bitlist{0x0c, 0x01}, // 0b00001100, 0b00000001
		},
		{
			a:    Bitlist{0xFF, 0x02}, // 0b11111111, 0x00000010
			b:    Bitlist{0x13, 0x03}, // 0b00010011, 0x00000011
			want: Bitlist{0xec, 0x03}, // 0b11101100, 0x00000011
		},
		{
			a:    Bitlist{0xFF, 0x87}, // 0b11111111, 0x10000111
			b:    Bitlist{0x13, 0x8F}, // 0b00010011, 0x10001111
			want: Bitlist{0xec, 0x88}, // 0b11101100, 0x10001000
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("(%x).Xor(%x)", tt.a, tt.b), func(t *testing.T) {
			if got, err := tt.a.Xor(tt.b); !bytes.Equal(got, tt.want) || err != nil {
				t.Errorf(
					"(%x).Xor(%x) = %x, %v, wanted %x",
					tt.a,
					tt.b,
					got,
					err,
					tt.want,
				)
			}
		})
	}
}

// Regression test for bitlists whose logical lengths match even though their
// backing slices do not.  This can happen with empty or otherwise
// non-canonical bitlists, and previously caused the binary operations to
// either panic or silently accept incompatible operands.
func TestBitlist_BinaryOperationsRejectDifferentBackingLengths(t *testing.T) {
	short := Bitlist{}
	long := Bitlist{0x00}
	if short.Len() != long.Len() {
		t.Fatalf("test setup has different logical lengths: %d != %d", short.Len(), long.Len())
	}

	operations := map[string]func(Bitlist, Bitlist) (Bitlist, error){
		"And": func(a, b Bitlist) (Bitlist, error) { return a.And(b) },
		"Or":  func(a, b Bitlist) (Bitlist, error) { return a.Or(b) },
		"Xor": func(a, b Bitlist) (Bitlist, error) { return a.Xor(b) },
	}
	for name, operation := range operations {
		operation := operation
		t.Run(name, func(t *testing.T) {
			for _, operands := range [][2]Bitlist{{short, long}, {long, short}} {
				if got, err := operation(operands[0], operands[1]); err != ErrBitlistDifferentLength || got != nil {
					t.Errorf("operation(%x, %x) = %x, %v, wanted nil, %v", operands[0], operands[1], got, err, ErrBitlistDifferentLength)
				}
			}
		})
	}

	t.Run("Contains", func(t *testing.T) {
		for _, operands := range [][2]Bitlist{{short, long}, {long, short}} {
			if got, err := operands[0].Contains(operands[1]); err != ErrBitlistDifferentLength || got {
				t.Errorf("(%x).Contains(%x) = %t, %v, wanted false, %v", operands[0], operands[1], got, err, ErrBitlistDifferentLength)
			}
		}
	})

	t.Run("NoAllocOr operands", func(t *testing.T) {
		if err := long.NoAllocOr(short, long); err != ErrBitlistDifferentLength {
			t.Fatalf("NoAllocOr returned %v, wanted %v", err, ErrBitlistDifferentLength)
		}
	})
	t.Run("NoAllocOr result", func(t *testing.T) {
		if err := long.NoAllocOr(long, short); err != ErrBitlistDifferentLength {
			t.Fatalf("NoAllocOr returned %v, wanted %v", err, ErrBitlistDifferentLength)
		}
	})
}

func TestBitlist_XorEmpty(t *testing.T) {
	got, err := (Bitlist{}).Xor(Bitlist{})
	if err != nil {
		t.Fatalf("Xor returned unexpected error: %v", err)
	}
	if got == nil || len(got) != 0 {
		t.Fatalf("Xor returned %#v, wanted a non-nil empty bitlist", got)
	}
}

func TestBitlist_Not(t *testing.T) {
	tests := []struct {
		a    Bitlist
		want Bitlist
	}{
		{
			a:    Bitlist{0x01}, // 0b00000001
			want: Bitlist{0x01}, // 0b00000001
		},
		{
			a:    Bitlist{0x02}, // 0b00000010
			want: Bitlist{0x03}, // 0b00000011
		},
		{
			a:    Bitlist{0x03}, // 0b00000011
			want: Bitlist{0x02}, // 0b00000010
		},
		{
			a:    Bitlist{0x05}, // 0b00000101
			want: Bitlist{0x06}, // 0b00000110
		},
		{
			a:    Bitlist{0x06}, // 0b00000110
			want: Bitlist{0x05}, // 0b00000101
		},
		{
			a:    Bitlist{0x83}, // 0b10000011
			want: Bitlist{0xfc}, // 0b11111100
		},
		{
			a:    Bitlist{0x13}, // 0b00010011
			want: Bitlist{0x1c}, // 0b00011100
		},
		{
			a:    Bitlist{0x1F}, // 0b00011111
			want: Bitlist{0x10}, // 0b00010000
		},
		{
			a:    Bitlist{0x1F, 0x03}, // 0b00011111, 0b00000011
			want: Bitlist{0xe0, 0x02}, // 0b11100000, 0b00000010
		},
		{
			a:    Bitlist{0x9F, 0x01}, // 0b10011111, 0b00000001
			want: Bitlist{0x60, 0x01}, // 0b01100000, 0b00000001
		},
		{
			a:    Bitlist{0xFF, 0x02}, // 0b11111111, 0x00000010
			want: Bitlist{0x00, 0x03}, // 0b00000000, 0x00000011
		},
		{
			a:    Bitlist{0xFF, 0x87}, // 0b11111111, 0x10000111
			want: Bitlist{0x00, 0xf8}, // 0b00000000, 0x11111000
		},
		{
			a:    Bitlist{0xFF, 0x07}, // 0b11111111, 0x00000111
			want: Bitlist{0x00, 0x04}, // 0b00000000, 0x00000100
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("(%#x).Not()", tt.a), func(t *testing.T) {
			if !bytes.Equal(tt.a.Not(), tt.want) {
				t.Errorf(
					"(%x).Not() = %x, wanted %x",
					tt.a,
					tt.a.Not(),
					tt.want,
				)
			}
		})
	}
}

func TestBitlist_BitIndices(t *testing.T) {
	tests := []struct {
		a    Bitlist
		want []int
	}{
		{
			a:    Bitlist{0b10010},
			want: []int{1},
		},
		{
			a:    Bitlist{0b10000},
			want: []int{},
		},
		{
			a:    Bitlist{0b10, 0b1},
			want: []int{1},
		},
		{
			a:    Bitlist{0b11111111, 0b11},
			want: []int{0, 1, 2, 3, 4, 5, 6, 7, 8},
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
