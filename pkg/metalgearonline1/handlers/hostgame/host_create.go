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
	handlers.Register(CreateGameHandler{})
}

type CreateGameHandler struct{}

func (h CreateGameHandler) Type() types.PacketType {
	return types.ClientCreateGame
}

func (h CreateGameHandler) ArgumentTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf(types.CreateGameOptions{}),
	}
}

func (h CreateGameHandler) Handle(_ *session.Session, _ *packet.Packet) ([]types.Response, error) {
	return nil, handlers.ErrNotImplemented
}

func (h CreateGameHandler) HandleArgs(s *session.Session, args *types.CreateGameOptions) ([]types.Response, error) {
	// Can't host a game until after connection info has been reported
	if s.ActiveConnection.ID == 0 {
		return []types.Response{ResponseCreateGame{ErrorCode: handlers.ErrInvalidArguments.Code}}, nil
	}

	// Can't host a game if you're already hosting
	if s.IsHost() {
		return []types.Response{ResponseCreateGame{ErrorCode: handlers.ErrInvalidArguments.Code}}, nil
	}

	opts := models.GameOptions{
		UserID: s.User.ID,
	}
	opts.FromCreateGameOptions(args)

	if tx := s.DB.Model(&opts).Create(&opts); tx.Error != nil {
		return []types.Response{ResponseCreateGame{ErrorCode: handlers.ErrDatabase.Code}}, tx.Error
	}

	newGame := models.Game{
		LobbyID:       uint(s.LobbyID),
		UserID:        s.User.ID,
		ConnectionID:  s.ActiveConnection.ID,
		GameOptionsID: opts.ID,
	}
	if tx := s.DB.Create(&newGame); tx.Error != nil {
		return []types.Response{ResponseCreateGame{ErrorCode: handlers.ErrDatabase.Code}}, tx.Error
	}

	s.StartHosting(types.GameID(newGame.ID), args.Rules)
	s.EventGameCreated(newGame.ID, args)

	s.GameState.AddPlayer(types.UserID(s.User.ID))
	return []types.Response{ResponseCreateGame{ErrorCode: 0}}, nil
}

// --------------------

type ResponseCreateGame types.ResponseErrorCode

func (r ResponseCreateGame) Type() types.PacketType { return types.ServerCreateGame }
