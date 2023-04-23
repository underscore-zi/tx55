package types

import "tx55/pkg/metalgearonline1/types/bitfield"

type GameMode byte

const (
	ModeDeathmatch GameMode = iota
	ModeTeamDeathmatch
	ModeRescue
	ModeCapture
	ModeSneaking
)

type GameMap byte

const (
	MapLostForest GameMap = iota + 1
	MapGhostFactory
	MapCityUnderSiege
	MapKillhouseA
	MapKillhouseB
	MapKillhouseC
	MapSvyatogornyjEast
	MapMountaintop
	MapGraninyGorkiLab
	MapPillboxPurgatory
	MapHighIce
	MapBrownTown
	MapAll = 0xFF
)

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

type VSRatingRestriction byte

const (
	VSRatingRestrictionNone  VSRatingRestriction = 0
	VSRatingRestrictionBelow VSRatingRestriction = 2
	VSRatingRestrictionAbove VSRatingRestriction = 1
)

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