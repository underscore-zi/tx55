package events

import "github.com/gorilla/websocket"

type Client struct {
	id      string
	conn    *websocket.Conn
	service *Service
}
