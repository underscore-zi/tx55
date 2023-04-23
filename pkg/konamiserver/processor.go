package konamiserver

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"strings"
	"time"
	"tx55/pkg/packet"
)

func (s *Server) mainLoop() error {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}
		if err = conn.(*net.TCPConn).SetKeepAlive(true); err != nil {
			return err
		}
		if err = conn.(*net.TCPConn).SetKeepAlivePeriod(time.Minute); err != nil {
			return err
		}

		c := newClient(conn)
		c.server = s
		c.gameClient = s.Config.ClientFactory(c.id)
		s.clients[c.id] = c

		if err = c.gameClient.OnConnected(conn, c.out); err != nil {
			c.once.Do(c.cleanup)
			continue
		}

		go c.reader()
		go c.writer()
	}
}

func (c *client) cleanup() {
	_ = c.conn.Close()
	delete(c.server.clients, c.id)

	go c.gameClient.OnDisconnected()
}

func asHex(v any) string {
	return fmt.Sprintf("%x", v)
}

func (c *client) reader() {
	defer c.once.Do(c.cleanup)

	var lookAside *packet.Packet
	for {
		buf := make([]byte, 1024)
		n, connError := c.conn.Read(buf)
		if n > 0 {
			p, err := packet.Unmarshal(c.server.Config.Key, buf[:n])

			if p.Sequence() == c.seq.in {
				c.seq.in++
			} else {
				c.server.Log.WithFields(log.Fields{
					"expected": c.seq.in,
					"packet":   p.Sequence(),
				}).Error("sequence mismatch")
				return
			}

			if lookAside != nil {
				// We have a long packet, so make sure this is the same type and append the data
				// then we reset `err` by unmarshalling it again so processing can look normal
				// finally nil the lookAside buffer, if we need it again it'll get set again
				if (*lookAside).Type() != p.Type() {
					c.server.Log.WithFields(log.Fields{
						"lookaside_type": asHex((*lookAside).Type()),
						"packet_type":    asHex(p.Type()),
					}).Error("lookaside type mismatch")
				} else {
					(*lookAside).SetData(append((*lookAside).Data(), p.Data()...))
					err = p.Unmarshal(c.server.Config.Key, (*lookAside).Marshal(c.server.Config.Key))
				}
				lookAside = nil
			}

			if err == packet.ErrPacketTooShort {
				lookAside = &p
				continue
			} else if err != nil {
				c.server.Log.WithError(err).Error("unmarshal error")
				continue
			}

			c.server.Log.WithFields(log.Fields{
				"cmd": asHex(p.Type()),
				"seq": p.Sequence(),
				"len": len(p.Data()),
			}).Debug("packet received")

			c.dumpPacket(PacketIn, &p)
			go c.dispatch(&p)
		}

		if connError != nil {
			if connError == io.EOF {
				// Client disconnected
			} else {
				// Some other error
			}
			return
		}
	}
}

func (c *client) writer() {
	defer c.once.Do(c.cleanup)
	for {
		select {
		case msg, ok := <-c.out:
			if !ok {
				// Channel closed
				return
			}
			msg.SetSequence(c.seq.out)
			bytes := msg.Marshal(c.server.Config.Key)
			_, err := c.conn.Write(bytes)

			c.server.Log.WithFields(log.Fields{
				"cmd": asHex(msg.Type()),
				"seq": msg.Sequence(),
				"len": len(msg.Data()),
			}).Debug("packet sent")

			c.dumpPacket(PacketOut, &msg)
			c.seq.out++

			if err != nil {
				if err == io.EOF {
					// Client disconnected
				} else {
					// Some other error
				}
				return
			}
		}
	}
}

func (c *client) hexDump(data []byte) string {
	var hexDumpSB strings.Builder
	offset := 0
	for offset < len(data) {
		hexDumpSB.WriteString(fmt.Sprintf("%08x: ", offset))
		// Print hex bytes
		for i := 0; i < 16; i++ {
			if i+offset < len(data) {
				hexDumpSB.WriteString(fmt.Sprintf("%02x ", data[offset+i]))
			} else {
				hexDumpSB.WriteString("   ")
			}
			if i%8 == 7 {
				hexDumpSB.WriteString(" ")
			}
		}
		// Print ASCII representation
		hexDumpSB.WriteString(" ")
		for i := 0; i < 16 && i+offset < len(data); i++ {
			b := data[offset+i]
			if b >= 32 && b <= 126 {
				hexDumpSB.WriteByte(b)
			} else {
				hexDumpSB.WriteString(".")
			}
		}
		hexDumpSB.WriteString("\n")
		offset += 16
	}
	return hexDumpSB.String()
}

type PacketDirection bool

const PacketIn PacketDirection = true
const PacketOut PacketDirection = false

func (c *client) dumpPacket(direction PacketDirection, p *packet.Packet) {
	if !c.server.Debug {
		return
	}

	if len(c.server.DebugPackets) > 0 {
		found := false
		for _, v := range c.server.DebugPackets {
			if v == (*p).Type() {
				found = true
				break
			}
		}
		if !found {
			return
		}
	}

	if (*p).Type() < 0x1000 {
		return
	}

	switch direction {
	case PacketIn:
		fmt.Printf("--> [%04x] Seq: %d Len: %d\n", (*p).Type(), (*p).Sequence(), (*p).Length())
		if (*p).Type() != 0x0005 && (*p).Type() != 0x0003 {
			c.recentPackets = append(c.recentPackets, *p)
			if len(c.recentPackets) > 10 {
				c.recentPackets = c.recentPackets[1:]
			}
		}
		for _, line := range strings.Split(c.hexDump((*p).Data()), "\n") {
			if len(line) > 0 {
				fmt.Println("-->       ", line)
			}
		}
	case PacketOut:
		fmt.Printf("<-- [%04x] Seq: %d Len: %d\n", (*p).Type(), (*p).Sequence(), (*p).Length())
		for _, line := range strings.Split(c.hexDump((*p).Data()), "\n") {
			if len(line) > 0 {
				fmt.Println("<--       ", line)
			}
		}
	}
}
