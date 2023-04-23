package session

import (
	"errors"
	"gorm.io/gorm/clause"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/metalgearonline1/types"
)

var ErrNotHosting = errors.New("user is not currently a host")

func (s *Session) StartHosting(id types.GameID) {
	s.isHost = true
	s.GameID = id
}

func (s *Session) StopHosting() {
	if game, err := s.Game(); err == nil {
		if err = game.Stop(s.DB); err != nil {
			s.Log.WithError(err).Error("Failed to stop game")
		}
	}
	s.isHost = false
	s.GameID = 0
}

func (s *Session) Game() (*models.Game, error) {
	if !s.isHost {
		return nil, ErrNotHosting
	}

	var game models.Game
	game.ID = uint(s.GameID)

	if tx := s.DB.Preload(clause.Associations).Find(&game); tx.Error != nil {
		return nil, tx.Error
	}

	return &game, nil
}
