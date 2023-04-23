package bitfield

import (
	"testing"
)

func TestBitfield_SetBits(t *testing.T) {
	b := Bitfield{}

	for i := 0; i < 8; i++ {
		b.SetBits(byte(i), 1, 1)
	}
	if b.Data != 0b11111111 {
		t.Errorf("Expected b.Bytes(%d) but got b.Bytes(%d)", 1, b.Data)
	}

	b = Bitfield{}
	for i := 0; i < 8; i += 2 {
		b.SetBits(byte(i), 2, 2)
	}
	if b.Data != 0b10101010 {
		t.Errorf("Expected b.Bytes(%d) but got b.Bytes(%d)", 1, b.Data)
	}

	b = Bitfield{}
	b.SetBits(4, 3, 7)
	if b.Data != 0b01110000 {
		t.Errorf("Expected b.Bytes(%d) but got b.Bytes(%d)", 1, b.Data)
	}
	b.SetBits(1, 3, 7)
	if b.Data != 0b01111110 {
		t.Errorf("Expected b.Bytes(%d) but got b.Bytes(%d)", 1, b.Data)
	}
	b.SetBits(3, 3, 0)
	if b.Data != 0b01000110 {
		t.Errorf("Expected b.Bytes(%d) but got b.Bytes(%d)", 1, b.Data)
	}
}

func TestBitfield_SetBit(t *testing.T) {
	b := Bitfield{}
	b.SetBit(4, true)
	if b.Data != 16 {
		t.Errorf("Expected b.Bytes(%d) but got b.Bytes(%d)", 1, b.Data)
	}
}

func TestBitfield_GetBit(t *testing.T) {
	b := Bitfield{}
	b.Data = 16
	if b.GetBit(4) != true {
		t.Errorf("Expected b.GetBit(4)(%t) but got b.GetBit(4)(%t)", true, b.GetBit(4))
	}
}

func TestBitfield_GetBits(t *testing.T) {
	b := Bitfield{}
	b.Data = 112
	if b.GetBits(4, 3) != 7 {
		t.Errorf("Expected b.GetBits(4, 3)(%d) but got b.GetBits(4, 3)(%d)", 7, b.GetBits(4, 3))
	}
}
