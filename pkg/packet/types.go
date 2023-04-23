package packet

import (
	"errors"
)

var ErrPacketTooShort = errors.New("packet too short")
var ErrPacketHashMismatch = errors.New("packet hash mismatch")

type RawHeader struct {
	Cmd uint16
	Len uint16
	Seq uint32
	Md5 [16]byte
}

type Packet interface {
	Unmarshal(key []byte, data []byte) (int, error)
	Marshal(key []byte) []byte

	SetType(uint16)
	SetSequence(uint32)
	SetData([]byte)
	SetDataFrom(v any) error

	Type() uint16
	Length() uint16
	Hash() [16]byte
	Sequence() uint32
	Data() []byte
	DataInto(any) error

	Header() RawHeader
}
