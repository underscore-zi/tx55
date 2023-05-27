package types

import (
	"fmt"
	"tx55/pkg/metalgearonline1/types/bitfield"
)

type PingInfo struct {
	UserID UserID
	Ping   uint32
}

type GameInfo struct {
	ID                GameID
	Name              [16]byte
	Description       [128]byte
	_                 byte
	IsHostOnly        bool
	Rules             [15]GameRules
	_                 [2]byte // Marks the end of the rules
	RedTeam           Team
	BlueTeam          Team
	WeaponRestriction WeaponRestrictions
	MaxPlayers        uint8
	PlayerCount       uint8
	_                 [22]byte // unknown
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
	_                 [7]byte
	Bitfield          bitfield.GameSettings
	AutoBalance       byte
	IdleKickMinutes   uint16
	TeamKillCount     uint16
	_                 uint32             // Original SaveMGO server always sends 80L
	Players           [9]GamePlayerStats `packet:"truncate"`
}

type GamePlayerStats struct {
	UserID      UserID
	DisplayName [16]byte
	Team        Team
	Kills       uint32
	Deaths      uint32
	Score       uint32
	Seconds     uint32
	Ping        uint32
}

// PlayerSpecificSettings is the "Online Game Bitfield" menu (inverted camera, camera speeds, etc)
type PlayerSpecificSettings struct {
	// Settings is a bitfield
	Settings bitfield.PlayerSettings
	FKeys    [12][26]byte
}

type SwitchOrientation bool

const (
	PlayerOrientation SwitchOrientation = false
	CameraOrientation SwitchOrientation = true
)

func (s SwitchOrientation) String() string {
	switch s {
	case PlayerOrientation:
		return "Player Orientation"
	case CameraOrientation:
		return "Camera Orientation"
	default:
		return fmt.Sprintf("Unknown (%t)", s)
	}
}

type GearSwitchMode byte

const (
	GearSwitchToggle    GearSwitchMode = 0
	GearSwitchFlashback GearSwitchMode = 1
	GearSwitchCycle     GearSwitchMode = 2
)

func (g GearSwitchMode) String() string {
	switch g {
	case GearSwitchToggle:
		return "Equipped/Unequipped"
	case GearSwitchFlashback:
		return "Flashback"
	case GearSwitchCycle:
		return "Cycle"
	default:
		return fmt.Sprintf("Unknown (%d)", g)
	}
}
