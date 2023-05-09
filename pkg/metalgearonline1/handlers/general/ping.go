package general

import (
	"reflect"
	"tx55/pkg/metalgearonline1/handlers"
	"tx55/pkg/metalgearonline1/session"
	"tx55/pkg/metalgearonline1/types"
	"tx55/pkg/packet"
)

func init() {
	handlers.Register(PingHandler{})
}

type PingHandler struct{}

func (h PingHandler) Type() types.PacketType {
	return types.ClientPing
}

func (h PingHandler) ArgumentTypes() []reflect.Type {
	return []reflect.Type{}
}

func (h PingHandler) Handle(_ *session.Session, _ *packet.Packet) ([]types.Response, error) {
	return []types.Response{ResponsePing{}}, nil
}

type ResponsePing types.ResponseEmpty

func (r ResponsePing) Type() types.PacketType { return types.ServerPing }
