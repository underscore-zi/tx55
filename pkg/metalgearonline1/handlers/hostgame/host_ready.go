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
	handlers.Register(HostReadyToCreateHandler{})
}

type HostReadyToCreateHandler struct{}

func (h HostReadyToCreateHandler) Type() types.PacketType {
	return types.ClientHostReadyToCreate
}

func (h HostReadyToCreateHandler) ArgumentTypes() []reflect.Type {
	return []reflect.Type{}
}

func (h HostReadyToCreateHandler) Handle(sess *session.Session, _ *packet.Packet) ([]types.Response, error) {
	sess.LogEntry().WithFields(logrus.Fields{
		"round_id": 0,
		"map_id":   sess.GameState.Rules[0].Map,
		"map":      sess.GameState.Rules[0].Map.String(),
		"mode":     sess.GameState.Rules[0].Mode.String(),
		"mode_id":  sess.GameState.Rules[0].Mode,
	}).Info("host ready")

	return []types.Response{ResponseHostReadyToCreate{}}, nil
}

type ResponseHostReadyToCreate types.ResponseEmpty

func (r ResponseHostReadyToCreate) Type() types.PacketType {
	return types.ServerHostReadyToCreate
}
