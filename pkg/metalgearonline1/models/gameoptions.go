package models

import (
	"tx55/pkg/metalgearonline1/types"
	"tx55/pkg/metalgearonline1/types/bitfield"
)

func init() {
	All = append(All, &GameOptions{})
}

// GameOptions holds the latest game options a user has set
type GameOptions struct {
	ID                uint `gorm:"primaryKey"`
	UserID            uint `gorm:"index"`
	User              User
	Name              []byte
	Description       []byte
	HasPassword       bool
	Password          []byte
	IsHostOnly        bool
	Rules             []types.GameRules `gorm:"serializer:json"`
	RedTeam           types.Team
	BlueTeam          types.Team
	WeaponRestriction types.WeaponRestrictions
	MaxPlayers        uint8
	RatingRestriction types.VSRatingRestriction
	Rating            uint32
	SneMinutes        uint32
	SneRounds         uint32
	CapMinutes        uint32
	CapRounds         uint32
	ResMinutes        uint32
	ResRounds         uint32
	TDMMinutes        uint32
	TDMRounds         uint32
	TDMTickets        uint32
	DMMinutes         uint32
	Bitfield          bitfield.GameSettings `gorm:"serializer:json"`
	AutoBalance       uint8
	IdleKickMinutes   uint16
	TeamKillCount     uint16
}

// GameInfo fills in most of the types.GameInfo struct, but PlayerCount and Players must be filled in externally
func (g *GameOptions) GameInfo() (out types.GameInfo) {
	out.ID = types.GameID(g.ID)
	copy(out.Name[:], g.Name[:])
	copy(out.Description[:], g.Description[:])
	out.IsHostOnly = g.IsHostOnly
	copy(out.Rules[:], g.Rules[:])
	out.RedTeam = g.RedTeam
	out.BlueTeam = g.BlueTeam
	out.WeaponRestriction = g.WeaponRestriction
	out.MaxPlayers = g.MaxPlayers
	out.RatingRestriction = g.RatingRestriction
	out.Rating = g.Rating
	out.SneMinutes = g.SneMinutes
	out.SneRounds = g.SneRounds
	out.CapMinutes = g.CapMinutes
	out.CapRounds = g.CapRounds
	out.ResMinutes = g.ResMinutes
	out.ResRounds = g.ResRounds
	out.TDMMinutes = g.TDMMinutes
	out.TDMRounds = g.TDMRounds
	out.TDMTickets = g.TDMTickets
	out.DMMinutes = g.DMMinutes
	out.Bitfield = g.Bitfield
	out.AutoBalance = g.AutoBalance
	out.IdleKickMinutes = g.IdleKickMinutes
	out.TeamKillCount = g.TeamKillCount
	return
}

func (g *GameOptions) CreateGameOptions() (out types.CreateGameOptions) {
	copy(out.Name[:], g.Name[:])
	copy(out.Description[:], g.Description[:])
	out.HasPassword = g.HasPassword
	copy(out.Password[:], g.Password[:])
	out.IsHostOnly = g.IsHostOnly
	copy(out.Rules[:], g.Rules[:])
	out.RedTeam = g.RedTeam
	out.BlueTeam = g.BlueTeam
	out.WeaponRestriction = g.WeaponRestriction
	out.MaxPlayers = g.MaxPlayers
	out.RatingRestriction = g.RatingRestriction
	out.Rating = g.Rating
	out.SneMinutes = g.SneMinutes
	out.SneRounds = g.SneRounds
	out.CapMinutes = g.CapMinutes
	out.CapRounds = g.CapRounds
	out.ResMinutes = g.ResMinutes
	out.ResRounds = g.ResRounds
	out.TDMMinutes = g.TDMMinutes
	out.TDMRounds = g.TDMRounds
	out.TDMTickets = g.TDMTickets
	out.DMMinutes = g.DMMinutes
	out.Bitfield = g.Bitfield
	out.AutoBalance = g.AutoBalance
	out.IdleKickMinutes = g.IdleKickMinutes
	out.TeamKillCount = g.TeamKillCount
	return
}

func (g *GameOptions) FromCreateGameOptions(opts *types.CreateGameOptions) {
	g.Name = opts.Name[:]
	g.Description = opts.Description[:]
	g.HasPassword = opts.HasPassword
	g.Password = opts.Password[:]
	g.IsHostOnly = opts.IsHostOnly
	g.Rules = opts.Rules[:]
	g.RedTeam = opts.RedTeam
	g.BlueTeam = opts.BlueTeam
	g.WeaponRestriction = opts.WeaponRestriction
	g.MaxPlayers = opts.MaxPlayers
	g.RatingRestriction = opts.RatingRestriction
	g.Rating = opts.Rating
	g.SneMinutes = opts.SneMinutes
	g.SneRounds = opts.SneRounds
	g.CapMinutes = opts.CapMinutes
	g.CapRounds = opts.CapRounds
	g.ResMinutes = opts.ResMinutes
	g.ResRounds = opts.ResRounds
	g.TDMMinutes = opts.TDMMinutes
	g.TDMRounds = opts.TDMRounds
	g.TDMTickets = opts.TDMTickets
	g.DMMinutes = opts.DMMinutes
	g.Bitfield = opts.Bitfield
	g.AutoBalance = opts.AutoBalance
	g.IdleKickMinutes = opts.IdleKickMinutes
	g.TeamKillCount = opts.TeamKillCount
	return
}
