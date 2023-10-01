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
	handlers.Register(HostPlayerLeaveHandler{})
}

type HostPlayerLeaveHandler struct{}

func (h HostPlayerLeaveHandler) Type() types.PacketType {
	return types.ClientHostPlayerLeave
}

func (h HostPlayerLeaveHandler) ArgumentTypes() (out []reflect.Type) {
	out = append(out, reflect.TypeOf(ArgsHostPlayerLeave{}))
	return
}

func (h HostPlayerLeaveHandler) Handle(_ *session.Session, _ *packet.Packet) (out []types.Response, err error) {
	out = append(out, ResponseHostPlayerLeave{ErrorCode: handlers.ErrNotImplemented.Code})
	err = handlers.ErrNotImplemented
	return
}

func (h HostPlayerLeaveHandler) HandleArgs(sess *session.Session, args *ArgsHostPlayerLeave) (out []types.Response, err error) {
	if !sess.IsHost() {
		out = append(out, ResponseHostPlayerLeave{ErrorCode: handlers.ErrNotHosting.Code})
		err = handlers.ErrNotHosting
		return
	}

	go sess.GameState.RemovePlayer(args.UserID)
	sess.EventGamePlayerLeft(args.UserID)
	out = append(out, ResponseHostPlayerLeave{ErrorCode: 0, UserID: args.UserID})

	sess.LogEntry().WithFields(logrus.Fields{
		"player_id": args.UserID,
	}).Info("quit game")

	return
}

// --- Packets ---

type ArgsHostPlayerLeave struct {
	UserID types.UserID
}

// ResponseHostPlayerLeave is a working response but looking at the client code, it might read more data
// after the UserID, not sure how that data is used though
type ResponseHostPlayerLeave struct {
	ErrorCode int32
	UserID    types.UserID
}

func (r ResponseHostPlayerLeave) Type() types.PacketType { return types.ServerHostPlayerLeave }
