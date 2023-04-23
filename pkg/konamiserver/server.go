package konamiserver

import (
	"github.com/sirupsen/logrus"
	"net"
)

type Config struct {
	Address       string
	Key           []byte
	ClientFactory GameClientFactory
	Log           logrus.FieldLogger
}

type Server struct {
	clients      map[string]*client
	listener     net.Listener
	Config       Config
	beforeHooks  map[uint16][]HookFunc
	afterHooks   map[uint16][]HookFunc
	outputHooks  map[uint16][]HookFunc
	Debug        bool
	DebugPackets []uint16
	Log          logrus.FieldLogger
}

// Start will start the server and block until the server is stopped
func (s *Server) Start() error {
	var err error
	s.listener, err = net.Listen("tcp", s.Config.Address)
	if err != nil {
		return err
	}
	return s.mainLoop()
}

func (s *Server) Stop() {
	// Should end up killing the mainLoops but I haven't confirmed that
	_ = s.listener.Close()
}

func NewServer(config Config) *Server {
	return &Server{
		clients:     make(map[string]*client),
		Config:      config,
		beforeHooks: make(map[uint16][]HookFunc),
		afterHooks:  make(map[uint16][]HookFunc),
		outputHooks: make(map[uint16][]HookFunc),
		Log:         logrus.New(),
	}
}
