package metalgearonline1

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"net"
	"tx55/pkg/metalgearonline1/handlers"
	"tx55/pkg/metalgearonline1/session"
	"tx55/pkg/metalgearonline1/types"
	"tx55/pkg/packet"
)

type GameClient struct {
	Session *session.Session
	Server  *GameServer
	out     chan packet.Packet
	conn    net.Conn
}

func (c *GameClient) OnConnected(conn net.Conn, out chan packet.Packet) (err error) {
	c.conn = conn
	c.out = out
	c.Session.IP = conn.RemoteAddr().(*net.TCPAddr).IP.String()

	if c.Server.IsBannedIP(c.Session.IP) {
		c.Session.LogEntry().Info("Banned IP, closing connection")
		return errors.New("banned ip")
	}

	return nil
}

func (c *GameClient) OnDisconnected() {
	if c.Session.User != nil {
		c.Session.LogEntry().Info("Disconnected")
	}

	if c.Session.IsHost() {
		c.Session.StopHosting()
	}

	c.Server.DeleteSession(c.Session.ID)
	return
}

func (c *GameClient) OnPacket(p *packet.Packet, out chan packet.Packet) error {
	cmd := types.PacketType((*p).Type())
	entry := c.Session.LogEntry().WithField("cmd", cmd.String())

	// Ping packets are too frequent to log
	if cmd > 0x0010 && cmd != types.ClientHostPingInformation && c.Session.IP != "127.0.0.1" {
		entry.Info("Received packet")
	} else {
		entry.Debug("Received packet")
	}

	replies, err := handlers.Handle(c.Session, p)
	entry.WithField("replies", len(replies))
	if err != nil {
		entry.WithError(err).Error("handler error")
	} else if cmd > 0x0010 && cmd != types.ClientHostPingInformation {
		entry.Debug("handler success")
	}

	if len(replies) > 0 {
		for i, reply := range replies {
			p, err := types.ToPacket(reply)

			if err != nil {
				entry.WithFields(log.Fields{
					"reply_index": i,
					"len":         p.Length(),
				}).WithError(err).Error("failed to marshal reply")
			} else {
				if cmd > 0x100 && c.Session.IP != "127.0.0.1" {
					entry.Debug("sending")
				}
				out <- p
			}
		}
	}

	return nil
}
