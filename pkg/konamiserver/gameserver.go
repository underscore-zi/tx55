package konamiserver

import (
	"net"
	"tx55/pkg/packet"
)

type GameClientFactory func(string) GameClient

type GameClient interface {
	// OnConnected is called when a client connects and is passed the remote address
	OnConnected(conn net.Conn, out chan packet.Packet) error
	OnDisconnected()
	OnPacket(*packet.Packet, chan packet.Packet) error
}
