package lobby

import (
	"reflect"
	"tx55/pkg/metalgearonline1/handlers"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/metalgearonline1/session"
	"tx55/pkg/metalgearonline1/types"
	"tx55/pkg/metalgearonline1/types/bitfield"
	"tx55/pkg/packet"
)

func init() {
	handlers.Register(UpdatePlayerSettingsHandler{})
}

// UpdatePlayerSettingsHandler is called when a user updated their "Online Game Bitfield" settings.
// if a player has invalid settings the server will also try to change their options through
// a call here when they are first joining the lobby
type UpdatePlayerSettingsHandler struct{}

func (h UpdatePlayerSettingsHandler) Type() types.PacketType {
	return types.ClientUpdatePlayerSettings
}

func (h UpdatePlayerSettingsHandler) ArgumentTypes() []reflect.Type {
	return []reflect.Type{reflect.TypeOf(ArgsUpdatePlayerSettings{})}
}

func (h UpdatePlayerSettingsHandler) Handle(sess *session.Session, packet *packet.Packet) ([]types.Response, error) {
	return nil, handlers.ErrNotImplemented
}

func (h UpdatePlayerSettingsHandler) HandleArgs(sess *session.Session, args *ArgsUpdatePlayerSettings) ([]types.Response, error) {
	var out []types.Response
	var settings models.PlayerSettings

	sess.DB.First(&settings, "user_id = ?", sess.User.ID)
	settings.UserID = sess.User.ID
	settings.FromBitfield(args.Settings)
	settings.FKey0 = args.FKeys[0][:]
	settings.FKey1 = args.FKeys[1][:]
	settings.FKey2 = args.FKeys[2][:]
	settings.FKey3 = args.FKeys[3][:]
	settings.FKey4 = args.FKeys[4][:]
	settings.FKey5 = args.FKeys[5][:]
	settings.FKey6 = args.FKeys[6][:]
	settings.FKey7 = args.FKeys[7][:]
	settings.FKey8 = args.FKeys[8][:]
	settings.FKey9 = args.FKeys[9][:]
	settings.FKey10 = args.FKeys[10][:]
	settings.FKey11 = args.FKeys[11][:]

	if tx := sess.DB.Save(&settings); tx.Error != nil {
		out = append(out, ResponseUpdatePlayerSettings{
			ErrorCode: handlers.ErrDatabase.Code,
		})
	} else {
		out = append(out, ResponseUpdatePlayerSettings{
			ErrorCode: 0,
		})
	}

	return out, nil
}

// --- Packets ---

type ArgsUpdatePlayerSettings struct {
	Settings bitfield.PlayerSettings
	FKeys    [12][26]byte
}

type ResponseUpdatePlayerSettings types.ResponseErrorCode

func (r ResponseUpdatePlayerSettings) Type() types.PacketType { return types.ServerPlayerSettings }
