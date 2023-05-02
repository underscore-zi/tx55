package types

// --- End Packets ---

type PlayerStatsPeriod uint32

const (
	PeriodAllTime PlayerStatsPeriod = 0
	PeriodWeekly  PlayerStatsPeriod = 1
	PeriodArchive PlayerStatsPeriod = 2
)

func (p PlayerStatsPeriod) String() string {
	switch p {
	case PeriodAllTime:
		return "All Time"
	case PeriodWeekly:
		return "Weekly"
	case PeriodArchive:
		return "Archive"
	default:
		return "Unknown"
	}
}

type PeriodStats struct {
	Period         PlayerStatsPeriod
	Deathmatch     GameTypeStatsWithRank
	TeamDeathmatch GameTypeStatsWithRank
	Rescue         GameTypeStatsWithRank
	Capture        GameTypeStatsWithRank
	Sneaking       GameTypeStatsWithRank
}

type HostReportedStats struct {
	Kills              int32
	Deaths             int32
	KillStreak         uint16
	DeathStreak        uint16
	Stuns              uint32
	StunsReceived      uint32
	SnakeFrags         uint32
	Points             int32
	Unknown1           uint32
	Unknown2           uint32
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
	// Unknown3 on the Java server was labeled VSRating but that is unconfirmed since VSRating is sent in the overview
	// packet and is not ever shown to be game-mode specific
	Unknown3 uint32
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
