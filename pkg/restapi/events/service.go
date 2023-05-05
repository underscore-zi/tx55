package events

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

type Service struct {
	clients     map[string]*Client
	addCh       chan *Client
	removeCh    chan *Client
	mu          sync.Mutex
	Log         logrus.FieldLogger
	upgrader    websocket.Upgrader
	ValidTokens []string
}

func NewService(logger logrus.FieldLogger, validTokens []string) *Service {
	s := &Service{
		clients:     make(map[string]*Client),
		addCh:       make(chan *Client),
		removeCh:    make(chan *Client),
		ValidTokens: validTokens,
		Log:         logger,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
	return s
}

// CheckToken is a simple utility function that determines if the token is valid for this service
func (s *Service) CheckToken(token string) bool {
	for _, validToken := range s.ValidTokens {
		if token == validToken {
			return true
		}
	}
	return false
}

func (s *Service) broadcast(message []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, client := range s.clients {
		err := client.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			s.Log.WithError(err).WithField("client_id", client.id).Error("broadcast error")
		}
	}
}

func (s *Service) Run() {
	s.Log.Info("Starting Events Service")
	for {
		select {
		case client := <-s.addCh:
			s.mu.Lock()
			s.clients[client.id] = client
			s.mu.Unlock()
			s.Log.Info("Client connected")
		case client := <-s.removeCh:
			s.mu.Lock()
			delete(s.clients, client.id)
			s.mu.Unlock()
			s.Log.Info("Client disconnected")
		}
	}
}
