package models

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"tx55/pkg/metalgearonline1/types"
)

func init() {
	All = append(All, &Game{})
}

var ErrNotFound = errors.New("not found")

type Game struct {
	gorm.Model
	LobbyID       uint
	UserID        uint
	User          User
	ConnectionID  uint
	Connection    Connection
	GameOptionsID uint
	GameOptions   GameOptions
	Players       []GamePlayers
	CurrentRound  byte
}

func (g *Game) Refresh(db *gorm.DB) error {
	tx := db.Model(g).Preload(clause.Associations).First(g)
	return tx.Error
}

func (g *Game) AddPlayer(db *gorm.DB, UserID uint) error {
	return db.Model(g).Association("Players").Append(&GamePlayers{UserID: UserID})
}

func (g *Game) RemovePlayer(db *gorm.DB, UserID uint) error {
	if player, found := g.FindPlayer(db, UserID); found && !player.DeletedAt.Valid {
		tx := db.Delete(&GamePlayers{}, player.ID)
		return tx.Error
	}
	return ErrNotFound
}

func (g *Game) KickPlayer(db *gorm.DB, UserID uint) error {
	if player, found := g.FindPlayer(db, UserID); found {
		player.WasKicked = true
		tx := db.Save(player)
		return tx.Error
		// They should be removed when the host alerts that the player left
	}
	return ErrNotFound
}

func (g *Game) FindPlayer(db *gorm.DB, UserID uint) (*GamePlayers, bool) {
	for i, player := range g.Players {
		if player.UserID == UserID {
			return &g.Players[i], true
		}
	}
	// Didn't find it in the active player list, resort to an unscoped query to find players who left
	var gp GamePlayers
	if tx := db.Unscoped().Where("user_id = ? AND game_id = ?", UserID, g.ID).First(&gp); tx.Error == nil {
		return &gp, true
	}
	return nil, false
}

func (g *Game) Stop(db *gorm.DB) error {
	db.Model(&GamePlayers{}).Where("game_id = ?", g.ID).Delete(&GamePlayers{})
	tx := db.Delete(&Game{}, g.ID)
	return tx.Error
}

func (g *Game) CheckPassword(password [16]byte) bool {
	pw := types.BytesToString(password[:])
	expectedPw := types.BytesToString(g.GameOptions.Password[:])
	return pw == expectedPw
}
