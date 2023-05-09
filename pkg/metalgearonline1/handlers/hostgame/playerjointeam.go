package hostgame

import (
	"reflect"
	"tx55/pkg/metalgearonline1/handlers"
	"tx55/pkg/metalgearonline1/session"
	"tx55/pkg/metalgearonline1/types"
	"tx55/pkg/packet"
)

func init() {
	handlers.Register(HostPlayerJoinTeam{})
}

type HostPlayerJoinTeam struct{}

func (h HostPlayerJoinTeam) Type() types.PacketType {
	return types.ClientHostPlayerJoinTeam
}

func (h HostPlayerJoinTeam) ArgumentTypes() (out []reflect.Type) {
	out = append(out, reflect.TypeOf(ArgsHostPlayerJoinTeam{}))
	return
}

func (h HostPlayerJoinTeam) Handle(_ *session.Session, _ *packet.Packet) (out []types.Response, err error) {
	out = append(out, ResponseHostPlayerJoinTeam{ErrorCode: handlers.ErrNotImplemented.Code})
	err = handlers.ErrNotImplemented
	return
}

func (h HostPlayerJoinTeam) HandleArgs(sess *session.Session, args *ArgsHostPlayerJoinTeam) (out []types.Response, err error) {
	if !sess.IsHost() {
		out = append(out, ResponseHostPlayerJoinTeam{ErrorCode: handlers.ErrNotHosting.Code})
		err = handlers.ErrNotHosting
		return
	}

	go sess.GameState.JoinTeam(args.UserID, args.TeamID)
	out = append(out, ResponseHostPlayerJoinTeam{ErrorCode: 0, UserID: args.UserID})
	return
}

// --- Packets ---

type ArgsHostPlayerJoinTeam struct {
	UserID types.UserID
	TeamID types.Team
}

type ResponseHostPlayerJoinTeam struct {
	ErrorCode int32
	UserID    types.UserID
}

func (r ResponseHostPlayerJoinTeam) Type() types.PacketType { return types.ServerHostPlayerJoinTeam }
