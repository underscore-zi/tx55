package hostgame

import (
	"reflect"
	"tx55/pkg/metalgearonline1/handlers"
	"tx55/pkg/metalgearonline1/session"
	"tx55/pkg/metalgearonline1/types"
	"tx55/pkg/packet"
)

func init() {
	handlers.Register(HostQuitHandler{})
}

type HostQuitHandler struct{}

func (h HostQuitHandler) Type() types.PacketType {
	return types.ClientHostQuitGame
}

func (h HostQuitHandler) ArgumentTypes() []reflect.Type {
	return []reflect.Type{}
}

func (h HostQuitHandler) Handle(sess *session.Session, _ *packet.Packet) ([]types.Response, error) {
	sess.StopHosting()
	return []types.Response{
		ResponseHostQuit{ErrorCode: 0},
	}, nil
}

type ResponseHostQuit types.ResponseErrorCode

func (r ResponseHostQuit) Type() types.PacketType { return types.ServerHostQuitGame }
