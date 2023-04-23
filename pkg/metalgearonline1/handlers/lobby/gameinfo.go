package lobby

import (
	"gorm.io/gorm/clause"
	"reflect"
	"tx55/pkg/metalgearonline1/handlers"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/metalgearonline1/session"
	"tx55/pkg/metalgearonline1/types"
	"tx55/pkg/packet"
)

func init() {
	handlers.Register(GetGameInfoHandler{})
}

// GetGameInfoHandler is called a couple times in game from the game list page.
// - Hovering over a game in the game list, after a couple seconds it'll show a little box
// - Selecting Player Info or Host Info from the menu in the game list
type GetGameInfoHandler struct{}

func (h GetGameInfoHandler) Type() types.PacketType {
	return types.ClientGetGameInfo
}

func (h GetGameInfoHandler) ArgumentTypes() []reflect.Type {
	return []reflect.Type{reflect.TypeOf(ArgsGetGameInfo{})}
}

func (h GetGameInfoHandler) Handle(sess *session.Session, packet *packet.Packet) ([]types.Response, error) {
	return nil, handlers.ErrNotImplemented
}

func (h GetGameInfoHandler) HandlerArgs(sess *session.Session, args *ArgsGetGameInfo) (out []types.Response, err error) {
	info := ResponseGameInfo{}

	var game models.Game
	game.ID = uint(args.GameID)
	if tx := sess.DB.Preload(clause.Associations).Preload("Players.User").First(&game); tx.Error != nil {
		err = tx.Error
		out = append(out, ResponseGameInfo{ErrorCode: 1})
		return
	}

	info.Info = game.GameOptions.GameInfo()

	for i, player := range game.Players {
		info.Info.Players[i] = player.GamePlayerStats()
	}
	info.Info.PlayerCount = byte(len(game.Players))

	return []types.Response{info}, nil

}

type ArgsGetGameInfo struct {
	GameID types.GameID
}

// --- Packets ---
type ResponseGameInfo struct {
	ErrorCode uint32
	Info      types.GameInfo
}

func (r ResponseGameInfo) Type() types.PacketType { return types.ServerGameInfo }
