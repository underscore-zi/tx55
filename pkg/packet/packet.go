package packet

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"errors"
	"reflect"
)

type packet struct {
	RawHeader
	Payload []byte
}

func (p *packet) SetType(c uint16) {
	p.Cmd = c
}

func (p *packet) SetSequence(s uint32) {
	p.Seq = s
}

func (p *packet) SetData(d []byte) {
	p.Payload = d
}

func (p *packet) Type() uint16 {
	return p.Cmd
}

func (p *packet) Length() uint16 {
	return uint16(len(p.Payload))
}

func (p *packet) Sequence() uint32 {
	return p.Seq
}

func (p *packet) Data() []byte {
	return p.Payload
}

func (p *packet) Hash() [16]byte {
	// Hash payload is: cmd(2) + length(2) + seq(4) + data(?) = 8 + ?
	var payload = make([]byte, 8+len(p.Payload))
	binary.BigEndian.PutUint16(payload[0:2], p.Cmd)
	binary.BigEndian.PutUint16(payload[2:4], uint16(len(p.Payload)))
	binary.BigEndian.PutUint32(payload[4:8], p.Seq)
	copy(payload[8:], p.Payload)

	hash := md5.New()
	hash.Write(payload)
	var result [16]byte
	copy(result[:], hash.Sum(nil))

	return result
}

func (p *packet) Unmarshal(key []byte, data []byte) error {
	if len(data) < 24 {
		return ErrPacketTooShort
	}

	// decrypt data
	for i := 0; i < len(data); i++ {
		data[i] ^= key[i%len(key)]
	}

	err := binary.Read(bytes.NewBuffer(data), binary.BigEndian, &p.RawHeader)
	if err != nil {
		return err
	}
	p.Payload = data[24:]

	// Check Length
	if uint16(len(p.Payload)) < p.RawHeader.Len {
		return ErrPacketTooShort
	} else if uint16(len(p.Payload)) > p.RawHeader.Len {
		// TODO: This tends to happen when we have multiple packets at once, should handle that better
		return ErrPacketTooLong
	}

	// Check hash
	var hash [16]byte
	copy(hash[:], data[8:24])
	if p.Hash() != hash {
		return ErrPacketHashMismatch
	}

	return nil
}

func (p *packet) Marshal(key []byte) []byte {
	// If the full packet is being marshalled, we take ownership over the hash and length values
	p.RawHeader.Md5 = p.Hash()
	p.RawHeader.Len = uint16(len(p.Payload))

	var buf bytes.Buffer
	buf.Write(p.RawHeader.Bytes())
	buf.Write(p.Payload)

	bytes := buf.Bytes()

	for i := 0; i < len(bytes); i++ {
		bytes[i] ^= key[i%len(key)]
	}

	return bytes
}

func (p *packet) Header() RawHeader {
	return p.RawHeader
}

// DataInto parses the packet's data into the given pointer
func (p *packet) DataInto(v any) error {
	return binary.Read(bytes.NewBuffer(p.Data()), binary.BigEndian, v)
}

// ErrInvalidTruncationTarget tag must only be used on [N]byte fields, or the struct must implement the IsZero() interface
var ErrInvalidTruncationTarget = errors.New("unable to truncate packet data")

func (p *packet) truncateField(value reflect.Value) (int, error) {
	switch value.Kind() {
	case reflect.Array:
		for i := value.Len() - 1; i >= 0; i-- {
			if !value.Index(i).IsZero() {
				remove := value.Len() - i - 1
				byteCount := binary.Size(value.Index(0).Interface())
				return remove * byteCount, nil
			}
		}
	}
	return 0, ErrInvalidTruncationTarget
}

// valueForTruncate tries to find the highest level truncation target, but it will recurse if necessary
// this is so truncation can support nested structs
func (p *packet) valueForTruncation(v any) *reflect.Value {
	typeOf := reflect.TypeOf(v)
	if typeOf.Kind() != reflect.Struct {
		return nil
	}

	fieldCount := typeOf.NumField()
	if fieldCount <= 0 {
		return nil
	}

	lastField := typeOf.Field(fieldCount - 1)
	lastValue := reflect.ValueOf(v).Field(fieldCount - 1)

	if lastField.Tag.Get("packet") == "truncate" {
		return &lastValue
	} else if lastField.Type.Kind() == reflect.Struct {
		nested := p.valueForTruncation(lastValue.Interface())
		if nested != nil {
			return nested
		} else {
			return &lastValue
		}
	}
	return nil
}

// truncate will return the number of bytes that should be removed from the end of the packet
// this can be controlled through the `packet:"truncate"` tag
func (p *packet) truncate(v any) (int, error) {
	if val := p.valueForTruncation(v); val != nil {
		return p.truncateField(*val)
	}
	return 0, nil
}

// SetDataFrom parses the data from the given struct into the packet's data
func (p *packet) SetDataFrom(v any) error {
	tmp := &bytes.Buffer{}
	if v != nil {
		if err := binary.Write(tmp, binary.BigEndian, v); err != nil {
			return err
		}
	}
	bs := tmp.Bytes()
	toTrim, _ := p.truncate(v)
	if toTrim > 0 {
		bs = bs[:len(bs)-toTrim]
	}

	p.SetData(bs)
	return nil
}

func New() Packet {
	return &packet{}
}

func Unmarshal(key, data []byte) (Packet, error) {
	p := &packet{}
	err := p.Unmarshal(key, data)
	return p, err
}
