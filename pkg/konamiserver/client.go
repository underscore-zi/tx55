package konamiserver

import (
	"github.com/google/uuid"
	"net"
	"sync"
	"tx55/pkg/packet"
)

type sequences struct {
	in  uint32
	out uint32
}

type client struct {
	id            string
	seq           sequences
	out           chan packet.Packet
	conn          net.Conn
	server        *Server
	once          sync.Once
	gameClient    GameClient
	activePacket  *packet.Packet
	writerMutex   sync.Mutex
	recentPackets []packet.Packet
}

func newClient(conn net.Conn) *client {
	return &client{
		id:   uuid.New().String(),
		out:  make(chan packet.Packet),
		conn: conn,
		seq:  sequences{1, 1},
	}
}
