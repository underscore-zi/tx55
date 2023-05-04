package hostgame

import (
	"reflect"
	"tx55/pkg/metalgearonline1/handlers"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/metalgearonline1/session"
	"tx55/pkg/metalgearonline1/types"
	"tx55/pkg/packet"
)

func init() {
	handlers.Register(CreateGameSettingsHandler{})
}

type CreateGameSettingsHandler struct{}

func (h CreateGameSettingsHandler) Type() types.PacketType {
	return types.ClientGetCreateGameSettings
}

func (h CreateGameSettingsHandler) ArgumentTypes() []reflect.Type {
	return []reflect.Type{}
}

func (h CreateGameSettingsHandler) Handle(sess *session.Session, packet *packet.Packet) ([]types.Response, error) {
	var latest models.GameOptions
	if tx := sess.DB.Model(&models.GameOptions{}).Where("user_id = ?", sess.User.ID).Order("id desc").First(&latest); tx.Error != nil {
		sess.Log.WithError(tx.Error).WithFields(sess.LogFields()).Error("failed to get latest game options")
	}
	return []types.Response{ResponseCreateGameSettings{Options: latest.CreateGameOptions()}}, nil
}

// --- Packets ---

func (r ResponseCreateGameSettings) Type() types.PacketType { return types.ServerCreateGameSettings }

type ResponseCreateGameSettings struct {
	ErrorCode uint32
	Options   types.CreateGameOptions
}
