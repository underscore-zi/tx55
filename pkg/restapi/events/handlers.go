package events

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"net/http"
)

// AcceptGinWebsocket is the gin compatible request handler that upgrades to a websocket connection
func (s *Service) AcceptGinWebsocket(c *gin.Context) {
	s.HandleWebSocketConnection(c.Writer, c.Request)
}

// HandleWebSocketConnection is a more generic request handler that upgrades to a websocket connection
func (s *Service) HandleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Log.WithError(err).Error("Upgrade error")
		return
	}
	defer func() { _ = conn.Close() }()

	clientID := uuid.New().String()
	client := &Client{
		id:      clientID,
		conn:    conn,
		service: s,
	}

	s.addCh <- client
	defer func() {
		s.removeCh <- client
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.Log.WithError(err).Error("ReadMessage error")
			break
		}
		s.Log.Info(string(msg))
	}
}

// PostNewEvent is the endpoint the game server should point to when it wants to broadcast an event
func (s *Service) PostNewEvent(c *gin.Context) {
	token := c.Param("token")
	if !s.CheckToken(token) {
		c.JSON(403, "bad token")
		return
	}

	bs, err := io.ReadAll(c.Request.Body)
	if err != nil {
		s.Log.WithError(err).Error("failed to read body")
		c.JSON(400, "body error")
		return
	}
	s.broadcast(bs)
	c.JSON(200, "ok")
}
