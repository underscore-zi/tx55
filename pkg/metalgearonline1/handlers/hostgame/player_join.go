package hostgame

import (
	"github.com/sirupsen/logrus"
	"reflect"
	"tx55/pkg/metalgearonline1/handlers"
	"tx55/pkg/metalgearonline1/session"
	"tx55/pkg/metalgearonline1/types"
	"tx55/pkg/packet"
)

func init() {
	handlers.Register(HostPlayerJoinHandler{})
}

type HostPlayerJoinHandler struct{}

func (h HostPlayerJoinHandler) Type() types.PacketType {
	return types.ClientHostPlayerJoin
}

func (h HostPlayerJoinHandler) ArgumentTypes() (out []reflect.Type) {
	out = append(out, reflect.TypeOf(ArgsHostPlayerJoin{}))
	return
}

func (h HostPlayerJoinHandler) Handle(_ *session.Session, _ *packet.Packet) (out []types.Response, err error) {
	out = append(out, ResponseHostPlayerJoin{ErrorCode: handlers.ErrNotImplemented.Code})
	err = handlers.ErrNotImplemented
	return
}

func (h HostPlayerJoinHandler) HandleArgs(sess *session.Session, args *ArgsHostPlayerJoin) (out []types.Response, err error) {
	if !sess.IsHost() {
		out = append(out, ResponseHostPlayerJoin{ErrorCode: handlers.ErrNotHosting.Code})
		err = handlers.ErrNotHosting
		return
	}

	go sess.GameState.AddPlayer(args.UserID)
	sess.EventGamePlayerJoined(args.UserID)
	out = append(out, ResponseHostPlayerJoin{ErrorCode: 0, UserID: args.UserID})

	sess.LogEntry().WithFields(logrus.Fields{
		"player_id": args.UserID,
	}).Info("connected to host")

	return
}

// --- Packets ---

type ArgsHostPlayerJoin struct {
	UserID types.UserID
}

type ResponseHostPlayerJoin struct {
	ErrorCode int32
	UserID    types.UserID
}

func (r ResponseHostPlayerJoin) Type() types.PacketType { return types.ServerHostPlayerJoin }
