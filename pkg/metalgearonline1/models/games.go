package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"tx55/pkg/metalgearonline1/types"
)

func init() {
	All = append(All, &Game{})
}

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

func (g *Game) CheckPassword(password [16]byte) bool {
	pw := types.BytesToString(password[:])
	expectedPw := types.BytesToString(g.GameOptions.Password[:])
	return pw == expectedPw
}
