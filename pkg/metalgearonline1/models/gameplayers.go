package models

import (
	"gorm.io/gorm"
	"tx55/pkg/metalgearonline1/types"
)

func init() {
	All = append(All, &GamePlayers{})
}

type GamePlayers struct {
	gorm.Model
	UserID    uint `gorm:"index"`
	User      User
	GameID    uint `gorm:"index"`
	Game      Game
	Team      types.Team
	Kills     uint32
	Deaths    uint32
	Score     uint32
	Seconds   uint32
	Ping      uint32
	WasKicked bool
}

func (g GamePlayers) GamePlayerStats() (out types.GamePlayerStats) {
	out.UserID = types.UserID(g.UserID)
	copy(out.DisplayName[:], g.User.DisplayName)
	out.Team = g.Team
	out.Kills = g.Kills
	out.Deaths = g.Deaths
	out.Score = g.Score
	out.Seconds = g.Seconds
	out.Ping = g.Ping
	return
}
