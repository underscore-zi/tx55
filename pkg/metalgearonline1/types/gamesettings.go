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
	ModeInvalid        GameMode = 255
)

type GameModeString string

const (
	ModeDeathmatchString          GameModeString = "deathmatch"
	ModeDeathmatchShortString     GameModeString = "dm"
	ModeTeamDeathmatchString      GameModeString = "team deathmatch"
	ModeTeamDeathmatchShortString GameModeString = "tdm"
	ModeRescueString              GameModeString = "rescue"
	ModeRescueShortString         GameModeString = "res"
	ModeCaptureString             GameModeString = "capture"
	ModeCaptureShortString        GameModeString = "cap"
	ModeSneakingString            GameModeString = "sneaking"
	ModeSneakingShortString       GameModeString = "sne"
	ModeOverallString             GameModeString = "overall"
	ModeOverallShortString        GameModeString = "all"
)

func (g GameMode) String() GameModeString {
	switch g {
	case ModeDeathmatch:
		return ModeDeathmatchString
	case ModeTeamDeathmatch:
		return ModeTeamDeathmatchString
	case ModeRescue:
		return ModeRescueString
	case ModeCapture:
		return ModeCaptureString
	case ModeSneaking:
		return ModeSneakingString
	case ModeOverall:
		return ModeOverallString
	default:
		return "Unknown"
	}
}

func (g GameMode) ShortString() GameModeString {
	switch g {
	case ModeDeathmatch:
		return ModeDeathmatchShortString
	case ModeTeamDeathmatch:
		return ModeTeamDeathmatchShortString
	case ModeRescue:
		return ModeRescueShortString
	case ModeCapture:
		return ModeCaptureShortString
	case ModeSneaking:
		return ModeSneakingShortString
	case ModeOverall:
		return ModeOverallShortString
	default:
		return "Unknown"
	}
}

func (s GameModeString) GameMode() GameMode {
	switch s {
	case ModeDeathmatchString, ModeDeathmatchShortString:
		return ModeDeathmatch
	case ModeTeamDeathmatchString, ModeTeamDeathmatchShortString:
		return ModeTeamDeathmatch
	case ModeRescueString, ModeRescueShortString:
		return ModeRescue
	case ModeCaptureString, ModeCaptureShortString:
		return ModeCapture
	case ModeSneakingString, ModeSneakingShortString:
		return ModeSneaking
	case ModeOverallString, ModeOverallShortString:
		return ModeOverall
	default:
		return ModeInvalid
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

type GameMapString string

const (
	MapLostForestString       GameMapString = "lost forest"
	MapGhostFactoryString     GameMapString = "ghost factory"
	MapCityUnderSiegeString   GameMapString = "city under siege"
	MapKillhouseAString       GameMapString = "killhouse a"
	MapKillhouseBString       GameMapString = "killhouse b"
	MapKillhouseCString       GameMapString = "killhouse c"
	MapSvyatogornyjEastString GameMapString = "svyatogornyj east"
	MapMountainTopString      GameMapString = "mountaintop"
	MapGraninyGorkiLabString  GameMapString = "graniny gorki lab"
	MapPillboxPurgatoryString GameMapString = "pillbox purgatory"
	MapHighIceString          GameMapString = "high ice"
	MapBrownTownString        GameMapString = "brown town"
	MapAllString              GameMapString = "all"
)

func (gm GameMap) String() GameMapString {
	switch gm {
	case MapLostForest:
		return MapLostForestString
	case MapGhostFactory:
		return MapGhostFactoryString
	case MapCityUnderSiege:
		return MapCityUnderSiegeString
	case MapKillhouseA:
		return MapKillhouseAString
	case MapKillhouseB:
		return MapKillhouseBString
	case MapKillhouseC:
		return MapKillhouseCString
	case MapSvyatogornyjEast:
		return MapSvyatogornyjEastString
	case MapMountainTop:
		return MapMountainTopString
	case MapGraninyGorkiLab:
		return MapGraninyGorkiLabString
	case MapPillboxPurgatory:
		return MapPillboxPurgatoryString
	case MapHighIce:
		return MapHighIceString
	case MapBrownTown:
		return MapBrownTownString
	case MapAll:
		return MapAllString
	default:
		return ""
	}
}

func (gms GameMapString) GameMap() GameMap {
	switch gms {
	case MapLostForestString:
		return MapLostForest
	case MapGhostFactoryString:
		return MapGhostFactory
	case MapCityUnderSiegeString:
		return MapCityUnderSiege
	case MapKillhouseAString:
		return MapKillhouseA
	case MapKillhouseBString:
		return MapKillhouseB
	case MapKillhouseCString:
		return MapKillhouseC
	case MapSvyatogornyjEastString:
		return MapSvyatogornyjEast
	case MapMountainTopString:
		return MapMountainTop
	case MapGraninyGorkiLabString:
		return MapGraninyGorkiLab
	case MapPillboxPurgatoryString:
		return MapPillboxPurgatory
	case MapHighIceString:
		return MapHighIce
	case MapBrownTownString:
		return MapBrownTown
	case MapAllString:
		return MapAll
	default:
		return 0
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
	// The game sends this as `3` when off, so you need to check the bitfield
	TeamKillCount uint16
}
