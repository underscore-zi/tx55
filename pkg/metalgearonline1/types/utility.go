package types

import (
	"bytes"
	"tx55/pkg/packet"
)

func (r RawPacket) Type() PacketType { return 0 }

type RawPacket struct {
	PacketType PacketType
	Data       []byte
}

func ToPacket(r Response) (packet.Packet, error) {
	p := packet.New()

	switch v := r.(type) {
	case RawPacket:
		p.SetType(uint16(v.PacketType))
		if err := p.SetDataFrom(v.Data); err != nil {
			return p, err
		}
	default:
		p.SetType(uint16(r.Type()))
		if err := p.SetDataFrom(r); err != nil {
			return p, err
		}
	}
	return p, nil
}

func BytesToString(b []byte) string {
	idx := bytes.IndexByte(b, 0)
	if idx == -1 {
		return string(b)
	} else {
		return string(b[:idx])
	}
}

var NONCE = [16]byte{
	0x84, 0xbd, 0xb8, 0xcf, 0xad, 0x46, 0xdd, 0x6e,
	0x42, 0x4a, 0xe4, 0xd8, 0xd2, 0x6a, 0x12, 0xf3,
}

// XORKEY is the key used to encrypt/decrypt packets, can't be easily configured on the PS2
var XORKEY = []byte{0x5a, 0x70, 0x85, 0xaf}
