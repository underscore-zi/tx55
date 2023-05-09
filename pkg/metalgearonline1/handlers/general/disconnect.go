package general

import (
	"reflect"
	"tx55/pkg/metalgearonline1/handlers"
	"tx55/pkg/metalgearonline1/session"
	"tx55/pkg/metalgearonline1/types"
	"tx55/pkg/packet"
)

func init() {
	handlers.Register(DisconnectHandler{})
}

type DisconnectHandler struct{}

func (h DisconnectHandler) Type() types.PacketType {
	return types.ClientDisconnect
}

func (h DisconnectHandler) ArgumentTypes() []reflect.Type {
	return []reflect.Type{}
}

func (h DisconnectHandler) Handle(_ *session.Session, _ *packet.Packet) ([]types.Response, error) {
	return nil, nil
}
