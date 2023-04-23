package testclient

import (
	"fmt"
	"net"
	"tx55/pkg/metalgearonline1/types"
	"tx55/pkg/packet"
)

func (c *TestClient) ExpectErrorCode(code int32) error {
	if p, err := c.Receive(); err != nil {
		return err
	} else {
		var payload types.ResponseErrorCode
		if err := (*p).DataInto(&payload); err != nil {
			return err
		}
		if payload.ErrorCode != code {
			return fmt.Errorf("expected error code %d, got %d", code, payload.ErrorCode)
		}
	}
	return nil
}

func (c *TestClient) Connect(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	c.conn = conn
	return nil
}

func (c *TestClient) Close() error {
	return c.conn.Close()
}

func (c *TestClient) Send(cmd types.PacketType, v any) error {
	p := packet.New()
	p.SetType(uint16(cmd))
	err := p.SetDataFrom(v)
	if err != nil {
		return err
	}

	c.seq++
	p.SetSequence(c.seq)
	bs := p.Marshal(c.Key)
	_, err = c.conn.Write(bs)
	return err
}

// Doesn't work well for large/multipart packets
func (c *TestClient) Receive() (*packet.Packet, error) {
	bs := make([]byte, 1024)
	n, err := c.conn.Read(bs)
	if err != nil {
		return nil, err
	}

	_, p, err := packet.Unmarshal(c.Key, bs[:n])
	return &p, err
}
