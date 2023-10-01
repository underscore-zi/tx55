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
	handlers.Register(HostPlayerKickedHandler{})
}

type HostPlayerKickedHandler struct{}

func (h HostPlayerKickedHandler) Type() types.PacketType {
	return types.ClientHostPlayerKicked
}

func (h HostPlayerKickedHandler) ArgumentTypes() (out []reflect.Type) {
	out = append(out, reflect.TypeOf(ArgsHostPlayerKicked{}))
	return
}

func (h HostPlayerKickedHandler) Handle(_ *session.Session, _ *packet.Packet) (out []types.Response, err error) {
	out = append(out, ResponseHostPlayerKicked{ErrorCode: handlers.ErrNotImplemented.Code})
	err = handlers.ErrNotImplemented
	return
}

func (h HostPlayerKickedHandler) HandleArgs(sess *session.Session, args *ArgsHostPlayerKicked) (out []types.Response, err error) {
	if !sess.IsHost() {
		out = append(out, ResponseHostPlayerKicked{ErrorCode: handlers.ErrNotHosting.Code})
		err = handlers.ErrNotHosting
		return
	}

	go sess.GameState.KickPlayer(args.UserID)
	out = append(out, ResponseHostPlayerKicked{ErrorCode: 0, UserID: args.UserID})

	sess.LogEntry().WithFields(logrus.Fields{
		"player_id": args.UserID,
	}).Info("kicked")

	return
}

// --- Packets ---

type ArgsHostPlayerKicked struct {
	UserID types.UserID
}

type ResponseHostPlayerKicked struct {
	ErrorCode int32
	UserID    types.UserID
}

func (r ResponseHostPlayerKicked) Type() types.PacketType { return types.ServerHostPlayerKicked }
