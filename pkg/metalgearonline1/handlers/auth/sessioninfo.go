package auth

import (
	"reflect"
	"tx55/pkg/metalgearonline1/handlers"
	"tx55/pkg/metalgearonline1/session"
	"tx55/pkg/metalgearonline1/types"
	"tx55/pkg/packet"
)

func init() {
	handlers.Register(SessionInfoHandler{})
}

type SessionInfoHandler struct{}

func (h SessionInfoHandler) Type() types.PacketType {
	return types.ClientGetSessionInfo
}

func (h SessionInfoHandler) ArgumentTypes() []reflect.Type {
	return []reflect.Type{}
}

func (h SessionInfoHandler) Handle(sess *session.Session, _ *packet.Packet) ([]types.Response, error) {
	res := ResponseSessionInfo{
		ErrorCode: 0,
		UserID:    types.UserID(sess.User.ID),
	}
	copy(res.DisplayName[:], sess.User.DisplayName)

	return []types.Response{res}, nil
}

// --- Packets ---
type ResponseSessionInfo struct {
	ErrorCode   uint32
	UserID      types.UserID
	DisplayName [16]byte
}

func (r ResponseSessionInfo) Type() types.PacketType { return types.ServerSessionInfo }
