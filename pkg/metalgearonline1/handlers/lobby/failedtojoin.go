package lobby

import (
	"reflect"
	"tx55/pkg/metalgearonline1/handlers"
	"tx55/pkg/metalgearonline1/session"
	"tx55/pkg/metalgearonline1/types"
	"tx55/pkg/packet"
)

func init() {
	handlers.Register(FailedToJoinHostHandler{})
}

type FailedToJoinHostHandler struct{}

func (h FailedToJoinHostHandler) Type() types.PacketType {
	return types.ClientPlayerFailedToJoinHost
}

func (h FailedToJoinHostHandler) ArgumentTypes() (out []reflect.Type) {
	return
}

func (h FailedToJoinHostHandler) Handle(_ *session.Session, _ *packet.Packet) (out []types.Response, err error) {
	// It might not be a bad idea to delete these hosts when this happens but since we should catch the disconnect
	// event and delete them anyway, I'm not going to worry about it for now.
	out = append(out, ResponsePlayerReadyToJoin{ErrorCode: 0})
	return
}

// --- Packets ---

type ResponsePlayerReadyToJoin types.ResponseErrorCode

func (r ResponsePlayerReadyToJoin) Type() types.PacketType { return types.ServerPlayerFailedToJoinHost }
