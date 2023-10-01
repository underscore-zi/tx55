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
	handlers.Register(HostNewRoundHandler{})
}

type HostNewRoundHandler struct{}

func (h HostNewRoundHandler) Type() types.PacketType {
	return types.ClientHostNewRound
}

func (h HostNewRoundHandler) ArgumentTypes() (out []reflect.Type) {
	out = append(out, reflect.TypeOf(ArgsHostNewRound{}))
	return
}

func (h HostNewRoundHandler) Handle(_ *session.Session, _ *packet.Packet) (out []types.Response, err error) {
	out = append(out, ResponseHostNewRound{ErrorCode: handlers.ErrNotImplemented.Code})
	err = handlers.ErrNotImplemented
	return
}

func (h HostNewRoundHandler) HandleArgs(sess *session.Session, args *ArgsHostNewRound) (out []types.Response, err error) {
	if !sess.IsHost() {
		out = append(out, ResponseHostNewRound{ErrorCode: handlers.ErrNotHosting.Code})
		err = handlers.ErrNotHosting
		return
	}

	go sess.GameState.NewRound(args.RoundID)
	sess.EventGameNewRound(args.RoundID)
	out = append(out, ResponseHostNewRound{ErrorCode: 0})

	sess.LogEntry().WithFields(logrus.Fields{
		"round_id": args.RoundID,
		"map_id":   sess.GameState.Rules[args.RoundID].Map,
		"map":      sess.GameState.Rules[args.RoundID].Map.String(),
		"mode":     sess.GameState.Rules[args.RoundID].Mode.String(),
		"mode_id":  sess.GameState.Rules[args.RoundID].Mode,
	}).Info("new round")

	return
}

// --- Packets ---

type ArgsHostNewRound struct {
	RoundID byte
}

type ResponseHostNewRound types.ResponseErrorCode

func (r ResponseHostNewRound) Type() types.PacketType { return types.ServerHostNewRound }
