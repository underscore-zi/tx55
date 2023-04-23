package types

// --- End Packets ---

type PlayerStatsPeriod uint32

const (
	PeriodAllTime PlayerStatsPeriod = 0
	PeriodWeekly  PlayerStatsPeriod = 1
)

type PeriodStats struct {
	Period         PlayerStatsPeriod
	Deathmatch     GameTypeStatsWithRank
	TeamDeathmatch GameTypeStatsWithRank
	Rescue         GameTypeStatsWithRank
	Capture        GameTypeStatsWithRank
	Sneaking       GameTypeStatsWithRank
}

type HostReportedStats struct {
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

func (s *HostReportedStats) AddStats(stats HostReportedStats) {
	s.Kills += stats.Kills
	s.Deaths += stats.Deaths

	if s.KillStreak < stats.KillStreak {
		s.KillStreak = stats.KillStreak
	}

	if s.DeathStreak < stats.DeathStreak {
		s.DeathStreak = stats.DeathStreak
	}

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

type GameTypeStatsWithRank struct {
	Stats HostReportedStats `json:"stats"`
	Rank  uint32            `json:"rank"`
}

type PlayerOverview struct {
	UserID      UserID
	DisplayName [16]byte
	Emblem      uint16
	U1          uint16
	EmblemText  [16]byte
	U2          [4]uint16
	VSRating    uint32
	CurrentTime uint32
	LastLogin   uint32
}

func (p *PlayerOverview) FillDisplayName(displayName []byte) {
	copy(p.DisplayName[:], displayName)
}

func (p *PlayerOverview) FillEmblem(hasEmblem bool, emblemText []byte) {
	if hasEmblem {
		p.Emblem = 1
		copy(p.EmblemText[:], emblemText)
	}
}
