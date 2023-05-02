package restapi

import (
	"time"
	"tx55/pkg/metalgearonline1/types"
)

// GET /api/v1/user/:user_id => UserJSON
// GET /api/v1/lobby/list => []LobbyJSON
// GET /api/v1/game/list => []GameJSON
// GET /api/v1/game/:game_id => GameJSON
// GET /api/v1/stats/:user_id => []PlayerStatsJSON

type ResponseJSON struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type UserJSON struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`

	DisplayName string `json:"display_name"`
	HasEmblem   bool   `json:"has_emblem"`
	EmblemText  string `json:"emblem_text"`
	OverallRank uint   `json:"overall_rank"`
	WeeklyRank  uint   `json:"weekly_rank"`
}

type LobbyJSON struct {
	ID      uint32 `json:"id"`
	Name    string `json:"name"`
	Players uint16 `json:"players"`
}

type GameJSON struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`

	LobbyID      uint              `json:"lobby_id"`
	UserID       uint              `json:"user_id"`
	Options      GameOptionsJSON   `json:"options"`
	Players      []GamePlayersJSON `json:"players"`
	CurrentRound byte              `json:"current_round"`
}

type GameRuleJSON struct {
	Map        types.GameMap  `json:"map"`
	MapString  string         `json:"map_string"`
	Mode       types.GameMode `json:"mode"`
	ModeString string         `json:"mode_string"`
}

type GamePlayersJSON struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`

	UserID     uint       `json:"user_id"`
	Team       types.Team `json:"team"`
	TeamString string     `json:"team_string"`
	Kills      uint32     `json:"kills"`
	Deaths     uint32     `json:"deaths"`
	Score      uint32     `json:"score"`
	Seconds    uint32     `json:"seconds"`
	Ping       uint32     `json:"ping"`
	WasKicked  bool       `json:"was_kicked"`
}

type GameOptionsJSON struct {
	Name              string                    `json:"name"`
	Description       string                    `json:"description"`
	HasPassword       bool                      `json:"has_password"`
	IsHostOnly        bool                      `json:"is_host_only"`
	Rules             []GameRuleJSON            `json:"rules"`
	RedTeam           types.Team                `json:"red_team"`
	BlueTeam          types.Team                `json:"blue_team"`
	WeaponRestriction types.WeaponRestrictions  `json:"weapon_restriction"`
	MaxPlayers        uint8                     `json:"max_players"`
	RatingRestriction types.VSRatingRestriction `json:"rating_restriction"`
	Rating            uint32                    `json:"rating"`
	SneMinutes        uint32                    `json:"sne_minutes"`
	SneRounds         uint32                    `json:"sne_rounds"`
	CapMinutes        uint32                    `json:"cap_minutes"`
	CapRounds         uint32                    `json:"cap_rounds"`
	ResMinutes        uint32                    `json:"res_minutes"`
	ResRounds         uint32                    `json:"res_rounds"`
	TDMMinutes        uint32                    `json:"tdm_minutes"`
	TDMRounds         uint32                    `json:"tdm_rounds"`
	TDMTickets        uint32                    `json:"tdm_tickets"`
	DMMinutes         uint32                    `json:"dm_minutes"`

	IdleKick         bool   `json:"idle_kick"`
	IdleKickMinutes  uint16 `json:"idle_kick_minutes"`
	TeamKillKick     bool   `json:"team_kill_kick"`
	TeamKillCount    uint16 `json:"team_kill_count"`
	AutoBalanced     bool   `json:"auto_balanced"`
	AutoBalanceCount uint8  `json:"auto_balance_count"`
	UniqueCharacters bool   `json:"unique_characters"`
	RumbleRoses      bool   `json:"rumble_roses"`
	Ghosts           bool   `json:"ghosts"`
	FriendFire       bool   `json:"friend_fire"`
	HasVoiceChat     bool   `json:"has_voice_chat"`
}

type PlayerStatsJSON struct {
	UserID    uint                    `json:"user_id"`
	UpdatedAt time.Time               `json:"updated_at"`
	Period    types.PlayerStatsPeriod `json:"period"`
	Mode      types.GameMode          `json:"mode"`
	Map       types.GameMap           `json:"map"`

	// Rank will be the rank in the mode for the period. Though stats are broken down by map also, the rank value will only consider mode.
	Rank uint32 `json:"rank"`

	Kills         int32  `json:"kills"`
	Deaths        int32  `json:"deaths"`
	KillStreak    uint16 `json:"kill_streak"`
	DeathStreak   uint16 `json:"death_streak"`
	Stuns         uint32 `json:"stuns"`
	StunsReceived uint32 `json:"stuns_received"`
	SnakeFrags    uint32 `json:"snake_frags"`
	Points        int32  `json:"points"`
	TeamKills     uint32 `json:"team_kills"`
	TeamStuns     uint32 `json:"team_stuns"`
	RoundsPlayed  uint32 `json:"rounds_played"`
	RoundsNoDeath uint32 `json:"rounds_no_death"`

	// KerotansForWin is the number of Gakos resuced when Mode is rescue
	KerotansForWin uint32 `json:"kerotans_for_win"`
	// KerotansPlaced is the number of goals as snake when Mode is sneaking
	KerotansPlaced uint32 `json:"kerotans_placed"`

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
}
