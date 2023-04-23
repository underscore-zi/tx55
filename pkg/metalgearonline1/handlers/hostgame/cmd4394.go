package hostgame

import (
	"reflect"
	"tx55/pkg/metalgearonline1/handlers"
	"tx55/pkg/metalgearonline1/session"
	"tx55/pkg/metalgearonline1/types"
	"tx55/pkg/packet"
)

func init() {
	handlers.Register(HostCmd4394Handler{})
}

type HostCmd4394Handler struct{}

func (h HostCmd4394Handler) Type() types.PacketType {
	return types.ClientHost4394
}

func (h HostCmd4394Handler) ArgumentTypes() (out []reflect.Type) {
	out = append(out, reflect.TypeOf(ArgsHostNewRound{}))
	return
}

func (h HostCmd4394Handler) Handle(_ *session.Session, _ *packet.Packet) (out []types.Response, err error) {
	// Not sure what this is, so...we just send back a success code and hope all is well
	out = append(out, ResponseHostCmd4394{ErrorCode: 0})
	return
}

// --- Packets ---

type ResponseHostCmd4394 types.ResponseErrorCode

func (r ResponseHostCmd4394) Type() types.PacketType { return types.ServerHost4394 }
