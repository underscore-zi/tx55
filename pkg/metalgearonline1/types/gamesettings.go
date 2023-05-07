package types

import "tx55/pkg/metalgearonline1/types/bitfield"

type GameMode byte

const (
	ModeDeathmatch     GameMode = iota
	ModeTeamDeathmatch GameMode = 1
	ModeRescue         GameMode = 2
	ModeCapture        GameMode = 3
	ModeSneaking       GameMode = 4
	ModeOverall        GameMode = 5 // Only used by the Rankings request
)

func (g GameMode) String() string {
	switch g {
	case ModeDeathmatch:
		return "Deathmatch"
	case ModeTeamDeathmatch:
		return "Team Deathmatch"
	case ModeRescue:
		return "Rescue"
	case ModeCapture:
		return "Capture"
	case ModeSneaking:
		return "Sneaking"
	default:
		return "Unknown"
	}
}

type GameMap byte

const (
	MapLostForest GameMap = iota + 1
	MapGhostFactory
	MapCityUnderSiege
	MapKillhouseA
	MapKillhouseB
	MapKillhouseC
	MapSvyatogornyjEast
	MapMountainTop
	MapGraninyGorkiLab
	MapPillboxPurgatory
	MapHighIce
	MapBrownTown
	MapAll = 0xFF
)

func (g GameMap) String() string {
	switch g {
	case MapLostForest:
		return "Lost Forest"
	case MapGhostFactory:
		return "Ghost Factory"
	case MapCityUnderSiege:
		return "City Under Siege"
	case MapKillhouseA:
		return "Killhouse A"
	case MapKillhouseB:
		return "Killhouse B"
	case MapKillhouseC:
		return "Killhouse C"
	case MapSvyatogornyjEast:
		return "Svyatogornyj East"
	case MapMountainTop:
		return "Mountaintop"
	case MapGraninyGorkiLab:
		return "Graniny Gorki Lab"
	case MapPillboxPurgatory:
		return "Pillbox Purgatory"
	case MapHighIce:
		return "High Ice"
	case MapBrownTown:
		return "Brown Town"
	default:
		return "Unknown"
	}
}

type Team byte

const (
	TeamGRU Team = iota
	TeamKGB
	TeamOcelot
)

const (
	TeamRed       = 0
	TeamBlue      = 1
	TeamSpectator = 0xFE
)

func (t Team) UniformString() string {
	switch t {
	case TeamGRU:
		return "GRU"
	case TeamKGB:
		return "KGB"
	case TeamOcelot:
		return "Ocelot"
	case TeamSpectator:
		return "Spectator"
	default:
		return "Unknown"
	}
}

func (t Team) ColorString() string {
	switch t {
	case TeamRed:
		return "Red"
	case TeamBlue:
		return "Blue"
	case TeamSpectator:
		return "Spectator"
	default:
		return "Unknown"
	}
}

type WeaponRestrictions byte

const (
	WeaponRestrictionsNone WeaponRestrictions = iota
	WeaponRestrictionsNoPrimary
	WeaponRestrictionsNoSecondary
	WeaponRestrictionsNoSupport
	WeaponRestrictionsPrimaryOnly
	WeaponRestrictionsSecondaryOnly
	WeaponRestrictionsSupportOnly
	WeaponRestrictionsKnifeOnly
)

func (w WeaponRestrictions) String() string {
	switch w {
	case WeaponRestrictionsNone:
		return "None"
	case WeaponRestrictionsNoPrimary:
		return "No Primary"
	case WeaponRestrictionsNoSecondary:
		return "No Secondary"
	case WeaponRestrictionsNoSupport:
		return "No Support"
	case WeaponRestrictionsPrimaryOnly:
		return "Primary Only"
	case WeaponRestrictionsSecondaryOnly:
		return "Secondary Only"
	case WeaponRestrictionsSupportOnly:
		return "Support Only"
	case WeaponRestrictionsKnifeOnly:
		return "Knife Only"
	default:
		return "Unknown"
	}
}

type VSRatingRestriction byte

const (
	VSRatingRestrictionNone  VSRatingRestriction = 0
	VSRatingRestrictionBelow VSRatingRestriction = 2
	VSRatingRestrictionAbove VSRatingRestriction = 1
)

func (v VSRatingRestriction) String() string {
	switch v {
	case VSRatingRestrictionNone:
		return "None"
	case VSRatingRestrictionBelow:
		return "Below"
	case VSRatingRestrictionAbove:
		return "Above"
	default:
		return "Unknown"
	}
}

type GameRules struct {
	Mode GameMode
	Map  GameMap
}

type CreateGameOptions struct {
	Name              [16]byte
	Description       [75]byte
	_                 [53]byte
	HasPassword       bool
	Password          [16]byte
	IsHostOnly        bool
	Rules             [15]GameRules
	_                 [2]byte // I think to mark the end of rules
	RedTeam           Team
	BlueTeam          Team
	WeaponRestriction WeaponRestrictions
	MaxPlayers        uint8
	_                 [12]byte // unknown
	RatingRestriction VSRatingRestriction
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
	Bitfield          bitfield.GameSettings
	AutoBalance       uint8
	IdleKickMinutes   uint16
	// The game sends this as `3` when off, so you need to check the bitoptions
	TeamKillCount uint16
}
