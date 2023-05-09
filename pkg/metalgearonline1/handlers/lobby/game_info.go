package lobby

import (
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

// GetGameInfoHandler is called a couple of times in game from the game list page.
// - Hovering over a game in the game list, after a couple seconds it'll show a little box
// - Selecting Player Info or Host Info from the menu in the game list
type GetGameInfoHandler struct{}

func (h GetGameInfoHandler) Type() types.PacketType {
	return types.ClientGetGameInfo
}

func (h GetGameInfoHandler) ArgumentTypes() []reflect.Type {
	return []reflect.Type{reflect.TypeOf(ArgsGetGameInfo{})}
}

func (h GetGameInfoHandler) Handle(_ *session.Session, _ *packet.Packet) ([]types.Response, error) {
	return nil, handlers.ErrNotImplemented
}

func (h GetGameInfoHandler) HandlerArgs(sess *session.Session, args *ArgsGetGameInfo) (out []types.Response, err error) {
	var info ResponseGameInfo
	var game models.Game
	var players []models.GamePlayers

	if err = sess.DB.Joins("GameOptions").First(&game, "games.id = ?", uint(args.GameID)).Error; err != nil {
		out = append(out, ResponseGameInfo{ErrorCode: handlers.ErrDatabase.Code})
		return
	}

	if err = sess.DB.Joins("User").Find(&players, "game_id = ?", uint(args.GameID)).Error; err != nil {
		out = append(out, ResponseGameInfo{ErrorCode: handlers.ErrDatabase.Code})
		return
	}

	info.Info = game.GameOptions.GameInfo()
	info.Info.PlayerCount = byte(len(players))
	for i, player := range players {
		info.Info.Players[i] = player.GamePlayerStats()
	}

	return []types.Response{info}, nil

}

type ArgsGetGameInfo struct {
	GameID types.GameID
}

type ResponseGameInfo struct {
	ErrorCode int32
	Info      types.GameInfo
}

func (r ResponseGameInfo) Type() types.PacketType { return types.ServerGameInfo }
