package packet

import (
	"testing"
)

func updateHash(data, newHash []byte) {
	for i, b := range newHash {
		data[8+i] = b
	}
}

func TestUnmarshal_Type(t *testing.T) {
	var key = []byte{0, 0, 0, 0}
	p := New()
	marshalled := p.Marshal(key)

	// Manually update the Type bytes
	marshalled[0] = 0xAA
	marshalled[1] = 0xBB

	updateHash(marshalled, []byte{223, 35, 162, 80, 88, 5, 41, 43, 122, 71, 49, 74, 16, 250, 56, 201})

	if n, err := p.Unmarshal(key, marshalled); err != nil {
		t.Error(err)
	} else if n != len(marshalled) {
		t.Errorf("Expected to consume %d bytes but consumed %d bytes", len(marshalled), n)
		t.Fail()
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
	if n, err := p.Unmarshal(key, marshalled); err != nil {
		t.Error(err)
	} else if n != len(marshalled) {
		t.Errorf("Expected to consume %d bytes but consumed %d bytes", len(marshalled), n)
		t.Fail()
	}

	if p.Length() != 0 {
		t.Errorf("Expected Length(%d) but got Length(%d)", 0, p.Length())
		t.Fail()
	}

	// Length Mismatch
	marshalled = append(marshalled, 0, 0, 0, 0)
	if n, err := p.Unmarshal(key, marshalled); p.Length() != 0 {
		t.Error("Expected Unmarshall to still return 0 length packet")
	} else if err != nil {
		t.Error("Unexpected error:", err)
	} else if n != len(marshalled)-4 {
		t.Errorf("Expected to consume %d bytes but consumed %d bytes", len(marshalled)-4, n)
		t.Fail()
	}

	// Correct length
	marshalled[3] = 0x04
	updateHash(marshalled, []byte{157, 24, 231, 144, 126, 16, 188, 187, 247, 78, 53, 51, 110, 189, 22, 30})
	if n, err := p.Unmarshal(key, marshalled); err != nil {
		t.Error("Unexpected error:", err)
	} else if n != len(marshalled) {
		t.Errorf("Expected to consume %d bytes but consumed %d bytes", len(marshalled), n)
		t.Fail()
	}

	if p.Length() != 4 {
		t.Errorf("Expected Length(%d) but got Length(%d)", 4, p.Length())
	}

}
