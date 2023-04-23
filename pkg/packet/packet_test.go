package packet

import (
	"fmt"
	"testing"
)

func updateHash(data, newHash []byte) {
	for i, b := range newHash {
		data[8+i] = b
	}
}
func printHash(p Packet) {
	fmt.Printf("[]byte{")
	for i, b := range p.Hash() {
		if i != 0 {
			fmt.Printf(", ")
		}
		fmt.Printf("%d", b)
	}
	fmt.Printf("}\n")
}

func TestUnmarshal_Type(t *testing.T) {
	var key = []byte{0, 0, 0, 0}
	p := New()
	marshalled := p.Marshal(key)

	// Manually update the Type bytes
	marshalled[0] = 0xAA
	marshalled[1] = 0xBB

	updateHash(marshalled, []byte{223, 35, 162, 80, 88, 5, 41, 43, 122, 71, 49, 74, 16, 250, 56, 201})

	if err := p.Unmarshal(key, marshalled); err != nil {
		t.Error(err)
	}

	if p.Type() != 0xAABB {
		t.Fail()
	}
}

func TestUnmarshal_Length(t *testing.T) {
	var key = []byte{0, 0, 0, 0}
	p := New()
	marshalled := p.Marshal(key)

	// Zero Length Packet
	if err := p.Unmarshal(key, marshalled); err != nil {
		t.Error(err)
	}

	if p.Length() != 0 {
		t.Errorf("Expected Length(%d) but got Length(%d)", 0, p.Length())
		t.Fail()
	}

	// Length Mismatch
	marshalled = append(marshalled, 0, 0, 0, 0)
	if err := p.Unmarshal(key, marshalled); err != ErrPacketTooLong {
		t.Error("Expected ErrPacketTooLong but got", err)
	}

	// Packet Length Mismatch
	marshalled[3] = 0x04
	updateHash(marshalled, []byte{157, 24, 231, 144, 126, 16, 188, 187, 247, 78, 53, 51, 110, 189, 22, 30})
	if err := p.Unmarshal(key, marshalled); err != nil {
		t.Error("Unexpected error:", err)
	}

	if p.Length() != 4 {
		t.Errorf("Expected Length(%d) but got Length(%d)", 4, p.Length())
	}

}
