package restapi

import (
	"strings"
	"time"
	"tx55/pkg/metalgearonline1/types"
)

// GET /api/v1/user/:user_id => UserJSON
// GET /api/v1/user/:user_id/stats => []PlayerStatsJSON
// GET /api/v1/user/:user_id/games => []GamePlayedJSON
// GET /api/v1/lobby/list => []LobbyJSON
// GET /api/v1/game/list => []GameJSON
// GET /api/v1/game/:game_id => GameJSON

// See `PeriodParam` for :period values and `GameModeParam` for optional mode values
// GET /api/v1/rankings/:period?mode=GameModeParam => []RankingEntryJSON
// GET /api/v1/ranking/:period/:page?mode=GameModeParam => []RankingEntryJSON

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
	User       *UserJSON  `json:"user"`
	Team       types.Team `json:"team"`
	TeamString string     `json:"team_string"`
	Kills      uint32     `json:"kills"`
	Deaths     uint32     `json:"deaths"`
	Score      uint32     `json:"score"`
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
	UserID       uint                    `json:"user_id"`
	UpdatedAt    time.Time               `json:"updated_at"`
	Period       types.PlayerStatsPeriod `json:"period"`
	PeriodString string                  `json:"period_string"`
	Mode         types.GameMode          `json:"mode"`
	ModeString   string                  `json:"mode_string"`
	Map          types.GameMap           `json:"map"`
	MapString    string                  `json:"map_string"`

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

	RadioUses          uint32        `json:"radio_uses"`
	TextChatUses       uint32        `json:"text_chat_uses"`
	CQCAttacks         uint32        `json:"cqc_attacks"`
	CQCAttacksReceived uint32        `json:"cqc_attacks_received"`
	HeadShots          uint32        `json:"head_shots"`
	HeadShotsReceived  uint32        `json:"head_shots_received"`
	TeamWins           uint32        `json:"team_wins"`
	KillsWithScorpion  uint32        `json:"kills_with_scorpion"`
	KillsWithKnife     uint32        `json:"kills_with_knife"`
	TimesEaten         uint32        `json:"times_eaten"`
	Rolls              uint32        `json:"rolls"`
	InfraredGoggleUses time.Duration `json:"infrared_goggle_uses"`
	PlayTime           time.Duration `json:"play_time"`
}

type RankingEntryJSON struct {
	Rank        uint   `json:"rank"`
	UserID      uint   `json:"user_id"`
	DisplayName string `json:"display_name"`
	Points      uint   `json:"points"`
}

type GamePlayedJSON struct {
	GameID          uint   `json:"game_id"`
	GameName        string `json:"game_name"`
	GameHasPassword bool   `json:"game_has_password"`
	GameHostID      uint   `json:"game_host_id"`
	WasHost         bool   `json:"was_host"`

	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at"`
	WasKicked bool      `json:"was_kicked"`
	Points    uint32    `json:"points"`
	Kills     uint32    `json:"kills"`
	Deaths    uint32    `json:"deaths"`
}

type NewsJSON struct {
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
}

// Params

type GameModeParam string

var GameModeParams = map[string]types.GameMode{
	"deathmatch":      types.ModeDeathmatch,
	"dm":              types.ModeDeathmatch,
	"team-deathmatch": types.ModeTeamDeathmatch,
	"tdm":             types.ModeTeamDeathmatch,
	"capture":         types.ModeCapture,
	"cap":             types.ModeCapture,
	"rescue":          types.ModeRescue,
	"res":             types.ModeRescue,
	"sneaking":        types.ModeSneaking,
	"sne":             types.ModeSneaking,
}

func (m GameModeParam) GameMode() (types.GameMode, bool) {
	v, found := GameModeParams[strings.ToLower(string(m))]
	return v, found
}

type PeriodParam string

var PeriodParams = map[string]types.PlayerStatsPeriod{
	"weekly":   types.PeriodWeekly,
	"week":     types.PeriodWeekly,
	"all-time": types.PeriodAllTime,
	"all":      types.PeriodAllTime,
	"archive":  types.PeriodArchive,
}

func (s PeriodParam) PlayerStatsPeriod() (types.PlayerStatsPeriod, bool) {
	v, found := PeriodParams[strings.ToLower(string(s))]
	return v, found
}

type GameMapParam string

var GameMapParams = map[string]types.GameMap{
	"lost-forest":       types.MapLostForest,
	"lfor":              types.MapLostForest,
	"ghost-factory":     types.MapGhostFactory,
	"gfact":             types.MapGhostFactory,
	"cus":               types.MapCityUnderSiege,
	"city-under-siege":  types.MapCityUnderSiege,
	"kha":               types.MapKillhouseA,
	"killhouse-a":       types.MapKillhouseA,
	"khb":               types.MapKillhouseB,
	"killhouse-b":       types.MapKillhouseB,
	"khc":               types.MapKillhouseC,
	"killhouse-c":       types.MapKillhouseC,
	"seast":             types.MapSvyatogornyjEast,
	"svyatogornyj-east": types.MapSvyatogornyjEast,
	"mtn":               types.MapMountainTop,
	"mtop:":             types.MapMountainTop,
	"mountain-top":      types.MapMountainTop,
	"ggl":               types.MapGraninyGorkiLab,
	"graniny-gorki-lab": types.MapGraninyGorkiLab,
	"pbp":               types.MapPillboxPurgatory,
	"pillbox-purgatory": types.MapPillboxPurgatory,
	"hice":              types.MapHighIce,
	"high-ice":          types.MapHighIce,
	"btown":             types.MapBrownTown,
	"brown-town":        types.MapBrownTown,
}

func (m GameMapParam) GameMap() (types.GameMap, bool) {
	v, found := GameMapParams[strings.ToLower(string(m))]
	return v, found
}
