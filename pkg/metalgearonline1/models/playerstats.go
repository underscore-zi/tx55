package models

import (
	"gorm.io/gorm"
	"tx55/pkg/metalgearonline1/types"
)

func init() {
	All = append(All, &PlayerStats{})
}

type PlayerStats struct {
	gorm.Model
	UserID uint `gorm:"index"`
	User   User
	Mode   types.GameMode
	Map    types.GameMap
	Period types.PlayerStatsPeriod

	Rank               uint32
	Kills              int32
	Deaths             int32
	KillStreak         uint16
	DeathStreak        uint16
	Stuns              uint32
	StunsReceived      uint32
	SnakeFrags         uint32
	Points             int32
	Suicides           uint32
	SelfStuns          uint32
	TeamKills          uint32
	TeamStuns          uint32
	RoundsPlayed       uint32
	RoundsNoDeath      uint32
	KerotansForWin     uint32
	KerotansPlaced     uint32
	RadioUses          uint32
	TextChatUses       uint32
	CQCAttacks         uint32
	CQCAttacksReceived uint32
	HeadShots          uint32
	HeadShotsReceived  uint32
	TeamWins           uint32
	KillsWithScorpion  uint32
	KillsWithKnife     uint32
	TimesEaten         uint32
	Rolls              uint32
	InfraredGoggleUses uint32
	PlayTime           uint32
}

func (s *PlayerStats) AddStats(stats types.HostReportedStats) {
	dest := s.ToHostReportedStats()
	dest.AddStats(stats)
	s.FromHostReportedStats(dest)
}

func (s *PlayerStats) ToHostReportedStats() (out types.HostReportedStats) {
	out.Kills = s.Kills
	out.Deaths = s.Deaths
	out.KillStreak = s.KillStreak
	out.DeathStreak = s.DeathStreak
	out.Stuns = s.Stuns
	out.StunsReceived = s.StunsReceived
	out.SnakeFrags = s.SnakeFrags
	out.Points = s.Points
	out.Suicides = s.Suicides
	out.SelfStuns = s.SelfStuns
	out.TeamKills = s.TeamKills
	out.TeamStuns = s.TeamStuns
	out.RoundsPlayed = s.RoundsPlayed
	out.RoundsNoDeath = s.RoundsNoDeath
	out.KerotansForWin = s.KerotansForWin
	out.KerotansPlaced = s.KerotansPlaced
	out.RadioUses = s.RadioUses
	out.TextChatUses = s.TextChatUses
	out.CQCAttacks = s.CQCAttacks
	out.CQCAttacksReceived = s.CQCAttacksReceived
	out.HeadShots = s.HeadShots
	out.HeadShotsReceived = s.HeadShotsReceived
	out.TeamWins = s.TeamWins
	out.KillsWithScorpion = s.KillsWithScorpion
	out.KillsWithKnife = s.KillsWithKnife
	out.TimesEaten = s.TimesEaten
	out.Rolls = s.Rolls
	out.InfraredGoggleUses = s.InfraredGoggleUses
	out.PlayTime = s.PlayTime
	return
}

func (s *PlayerStats) FromHostReportedStats(in types.HostReportedStats) {
	s.Kills = in.Kills
	s.Deaths = in.Deaths
	s.KillStreak = in.KillStreak
	s.DeathStreak = in.DeathStreak
	s.Stuns = in.Stuns
	s.StunsReceived = in.StunsReceived
	s.SnakeFrags = in.SnakeFrags
	s.Points = in.Points
	s.Suicides = in.Suicides
	s.SelfStuns = in.SelfStuns
	s.TeamKills = in.TeamKills
	s.TeamStuns = in.TeamStuns
	s.RoundsPlayed = in.RoundsPlayed
	s.RoundsNoDeath = in.RoundsNoDeath
	s.KerotansForWin = in.KerotansForWin
	s.KerotansPlaced = in.KerotansPlaced
	s.RadioUses = in.RadioUses
	s.TextChatUses = in.TextChatUses
	s.CQCAttacks = in.CQCAttacks
	s.CQCAttacksReceived = in.CQCAttacksReceived
	s.HeadShots = in.HeadShots
	s.HeadShotsReceived = in.HeadShotsReceived
	s.TeamWins = in.TeamWins
	s.KillsWithScorpion = in.KillsWithScorpion
	s.KillsWithKnife = in.KillsWithKnife
	s.TimesEaten = in.TimesEaten
	s.Rolls = in.Rolls
	s.InfraredGoggleUses = in.InfraredGoggleUses
	s.PlayTime = in.PlayTime
	return
}
