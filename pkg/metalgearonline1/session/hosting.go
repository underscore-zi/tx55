package session

import (
	"gorm.io/gorm"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/metalgearonline1/types"
)

func (s *Session) StartHosting(id types.GameID, rules [15]types.GameRules) {
	s.isHost = true
	s.GameState = &HostSession{
		GameID:       id,
		Rules:        rules,
		CurrentRound: 0,
		Players: map[types.UserID]bool{
			types.UserID(s.User.ID): true,
		},
		ParentSession: s,
	}
}

func (s *Session) StopHosting() {
	s.GameState.StopGame()
	s.EventGameDeleted()
	s.GameState = nil
}

func (hs *HostSession) AddPlayer(id types.UserID) {
	md := hs.ParentSession.DB.Model(&models.Game{
		Model: gorm.Model{ID: uint(hs.GameID)},
	})
	if err := md.Association("Players").Append(&models.GamePlayers{UserID: uint(id)}); err != nil {
		hs.ParentSession.Log.WithError(err).Error("Failed to add player to game")
	}

	hs.Lock.Lock()
	defer hs.Lock.Unlock()
	hs.Players[id] = true
}

func (hs *HostSession) RemovePlayer(id types.UserID) {
	md := hs.ParentSession.DB.Model(&models.GamePlayers{})
	if err := md.Delete(&models.GamePlayers{}, "game_id = ? AND user_id = ?", hs.GameID, id).Error; err != nil {
		hs.ParentSession.Log.WithError(err).Error("Failed to remove player from game")
	}

	hs.Lock.Lock()
	defer hs.Lock.Unlock()
	delete(hs.Players, id)
}

func (hs *HostSession) JoinTeam(id types.UserID, team types.Team) {
	player := models.GamePlayers{
		UserID: uint(id),
		GameID: uint(hs.GameID),
	}
	hs.ParentSession.DB.Model(&player).Where(&player).Update("team", team)
}

func (hs *HostSession) KickPlayer(id types.UserID) {
	player := models.GamePlayers{
		UserID: uint(id),
		GameID: uint(hs.GameID),
	}
	hs.ParentSession.DB.Model(&player).Where(&player).Update("was_kicked", true)
	// Removing the player will happen when the host sends the player left message
}

func (hs *HostSession) StopGame() {
	if err := hs.ParentSession.DB.Where("game_id = ?", hs.GameID).Delete(&models.GamePlayers{}).Error; err != nil {
		hs.ParentSession.Log.WithError(err).Error("Failed to remove players from game")
	}

	if err := hs.ParentSession.DB.Delete(&models.Game{}, uint(hs.GameID)).Error; err != nil {
		hs.ParentSession.Log.WithError(err).Error("Failed to remove game")
	}
}
