package session

import (
	"errors"
	"gorm.io/gorm/clause"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/metalgearonline1/types"
)

var ErrNotHosting = errors.New("user is not currently a host")

func (s *Session) StartHosting(id types.GameID, rules [15]types.GameRules) {
	s.isHost = true
	s.GameState = &HostSession{
		GameID:       id,
		Rules:        rules,
		CurrentRound: 0,
		Players: map[types.UserID]bool{
			types.UserID(s.User.ID): true,
		},
	}
}

func (s *Session) StopHosting() {
	if game, err := s.Game(); err == nil {
		if err = game.Stop(s.DB); err != nil {
			s.Log.WithError(err).Error("Failed to stop game")
		}
	}
	s.GameState = nil
}

func (s *Session) Game() (*models.Game, error) {
	if !s.IsHost() {
		return nil, ErrNotHosting
	}

	var game models.Game
	game.ID = uint(s.GameState.GameID)

	if tx := s.DB.Preload(clause.Associations).Find(&game); tx.Error != nil {
		return nil, tx.Error
	}

	return &game, nil
}

func (hs *HostSession) AddPlayer(id types.UserID) {
	hs.Lock.Lock()
	defer hs.Lock.Unlock()
	hs.Players[id] = true
}

func (hs *HostSession) RemovePlayer(id types.UserID) {
	hs.Lock.Lock()
	defer hs.Lock.Unlock()
	delete(hs.Players, id)
}
