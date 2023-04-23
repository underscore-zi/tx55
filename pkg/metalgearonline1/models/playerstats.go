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

	Rank               uint32 `json:"rank"`
	Kills              uint32 `json:"kills"`
	Deaths             uint32 `json:"deaths"`
	KillStreak         uint16 `json:"kill_streak"`
	DeathStreak        uint16 `json:"death_streak"`
	Stuns              uint32 `json:"stuns"`
	StunsReceived      uint32 `json:"stuns_received"`
	SnakeFrags         uint32 `json:"snake_frags"`
	Points             uint32 `json:"points"`
	Unknown1           uint32 `json:"unknown1"`
	Unknown2           uint32 `json:"unknown2"`
	TeamKills          uint32 `json:"team_kills"`
	TeamStuns          uint32 `json:"team_stuns"`
	RoundsPlayed       uint32 `json:"rounds_played"`
	RoundsNoDeath      uint32 `json:"rounds_no_death"`
	KerotansForWin     uint32 `json:"kerotans_for_win"`
	KerotansPlaced     uint32 `json:"kerotans_placed"`
	RadioUses          uint32 `json:"radio_uses"`
	TextChatUses       uint32 `json:"text_chat_uses"`
	CQCAttacks         uint32 `json:"cqc_attacks"`
	CQCAttacksReceived uint32 `json:"cqc_attacks_received"`
	HeadShots          uint32 `json:"head_shots"`
	HeadShotsReceived  uint32 `json:"head_shots_received"`
	TeamWins           uint32 `json:"team_wins"`
	KillsWithScorpion  uint32 `json:"kills_with_scorpion"`
	KillsWithKnife     uint32 `json:"kills_with_knife"`
	TimesEaten         uint32 `json:"times_eaten"`
	Rolls              uint32 `json:"rolls"`
	InfraredGoggleUses uint32 `json:"infrared_goggle_uses"`
	PlayTime           uint32 `json:"play_time"`
	Unknown3           uint32 `json:"unknown3"`
}

func (s *PlayerStats) AddStats(stats types.HostReportedStats) {
	s.Kills += stats.Kills
	s.Deaths += stats.Deaths
	s.KillStreak += stats.KillStreak
	s.DeathStreak += stats.DeathStreak
	s.Stuns += stats.Stuns
	s.StunsReceived += stats.StunsReceived
	s.SnakeFrags += stats.SnakeFrags
	s.Points += stats.Points
	s.Unknown1 += stats.Unknown1
	s.Unknown2 += stats.Unknown2
	s.TeamKills += stats.TeamKills
	s.TeamStuns += stats.TeamStuns
	s.RoundsPlayed += stats.RoundsPlayed
	s.RoundsNoDeath += stats.RoundsNoDeath
	s.KerotansForWin += stats.KerotansForWin
	s.KerotansPlaced += stats.KerotansPlaced
	s.RadioUses += stats.RadioUses
	s.TextChatUses += stats.TextChatUses
	s.CQCAttacks += stats.CQCAttacks
	s.CQCAttacksReceived += stats.CQCAttacksReceived
	s.HeadShots += stats.HeadShots
	s.HeadShotsReceived += stats.HeadShotsReceived
	s.TeamWins += stats.TeamWins
	s.KillsWithScorpion += stats.KillsWithScorpion
	s.KillsWithKnife += stats.KillsWithKnife
	s.TimesEaten += stats.TimesEaten
	s.Rolls += stats.Rolls
	s.InfraredGoggleUses += stats.InfraredGoggleUses
	s.PlayTime += stats.PlayTime
	s.Unknown3 += stats.Unknown3
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
	out.Unknown1 = s.Unknown1
	out.Unknown2 = s.Unknown2
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
	out.Unknown3 = s.Unknown3
	return
}
