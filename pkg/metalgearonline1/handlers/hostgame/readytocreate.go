package hostgame

import (
	"reflect"
	"tx55/pkg/metalgearonline1/handlers"
	"tx55/pkg/metalgearonline1/session"
	"tx55/pkg/metalgearonline1/types"
	"tx55/pkg/packet"
)

func init() {
	handlers.Register(HostReadyToCreateHandler{})
}

type HostReadyToCreateHandler struct{}

func (h HostReadyToCreateHandler) Type() types.PacketType {
	return types.ClientHostReadyToCreate
}

func (h HostReadyToCreateHandler) ArgumentTypes() []reflect.Type {
	return []reflect.Type{}
}

func (h HostReadyToCreateHandler) Handle(sess *session.Session, packet *packet.Packet) ([]types.Response, error) {
	return []types.Response{ResponseHostReadyToCreate{}}, nil
}

type ResponseHostReadyToCreate types.ResponseEmpty

func (r ResponseHostReadyToCreate) Type() types.PacketType {
	return types.ServerHostReadyToCreate
}
