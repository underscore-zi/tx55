package hostgame

import (
	"gorm.io/gorm"
	"reflect"
	"tx55/pkg/metalgearonline1/handlers"
	"tx55/pkg/metalgearonline1/models"
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

	game := &models.Game{Model: gorm.Model{ID: uint(sess.GameState.GameID)}}
	tx := sess.DB.Model(game).Update("current_round", args.RoundID)
	if tx.Error != nil {
		out = append(out, ResponseHostNewRound{ErrorCode: handlers.ErrDatabase.Code})
		err = tx.Error
		return
	}

	sess.GameState.CurrentRound = args.RoundID
	out = append(out, ResponseHostNewRound{ErrorCode: 0})
	return
}

// --- Packets ---

type ArgsHostNewRound struct {
	RoundID byte
}

type ResponseHostNewRound types.ResponseErrorCode

func (r ResponseHostNewRound) Type() types.PacketType { return types.ServerHostNewRound }
