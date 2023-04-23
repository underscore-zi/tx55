package hostgame

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"reflect"
	"strings"
	"tx55/pkg/metalgearonline1/handlers"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/metalgearonline1/session"
	"tx55/pkg/metalgearonline1/types"
	"tx55/pkg/packet"
)

func init() {
	handlers.Register(HostPlayerStatsHandler{})
}

type HostPlayerStatsHandler struct{}

func (h HostPlayerStatsHandler) Type() types.PacketType {
	return types.ClientHostPlayerStats
}

func (h HostPlayerStatsHandler) ArgumentTypes() (out []reflect.Type) {
	out = append(out, reflect.TypeOf(ArgsHostPlayerStats{}))
	return
}

func (h HostPlayerStatsHandler) Handle(_ *session.Session, _ *packet.Packet) (out []types.Response, err error) {
	out = append(out, ResponseHostPlayerStats{ErrorCode: handlers.ErrNotImplemented.Code})
	err = handlers.ErrNotImplemented
	return
}

func (h HostPlayerStatsHandler) HandleArgs(sess *session.Session, args *ArgsHostPlayerStats) (out []types.Response, err error) {
	if !sess.IsHost() {
		out = append(out, ResponseHostPlayerStats{ErrorCode: handlers.ErrNotHosting.Code})
		err = handlers.ErrNotHosting
		return
	}

	// Ensure the player is actually in the game
	if _, found := sess.GameState.Players[args.UserID]; !found {
		out = append(out, ResponseHostPlayerStats{ErrorCode: handlers.ErrNotFound.Code})
		err = handlers.ErrNotFound
		return
	}

	// Do the update, or create the stats as needed
	if err = h.updatePlayerStats(sess, uint(args.UserID), args.Stats); err != nil {
		out = append(out, ResponseHostPlayerStats{ErrorCode: handlers.ErrDatabase.Code})
		return
	}

	out = append(out, ResponseHostPlayerStats{ErrorCode: 0})
	return
}

func (h HostPlayerStatsHandler) updatePlayerStats(sess *session.Session, UserID uint, stats types.HostReportedStats) error {
	updates := map[string]interface{}{
		"kills":                gorm.Expr("kills + ?", stats.Kills),
		"deaths":               gorm.Expr("deaths + ?", stats.Deaths),
		"stuns":                gorm.Expr("stuns + ?", stats.Stuns),
		"stuns_received":       gorm.Expr("stuns_received + ?", stats.StunsReceived),
		"snake_frags":          gorm.Expr("snake_frags + ?", stats.SnakeFrags),
		"points":               gorm.Expr("points + ?", stats.Points),
		"unknown1":             gorm.Expr("unknown1 + ?", stats.Unknown1),
		"unknown2":             gorm.Expr("unknown2 + ?", stats.Unknown2),
		"team_kills":           gorm.Expr("team_kills + ?", stats.TeamKills),
		"team_stuns":           gorm.Expr("team_stuns + ?", stats.TeamStuns),
		"rounds_played":        gorm.Expr("rounds_played + ?", stats.RoundsPlayed),
		"rounds_no_death":      gorm.Expr("rounds_no_death + ?", stats.RoundsNoDeath),
		"kerotans_for_win":     gorm.Expr("kerotans_for_win + ?", stats.KerotansForWin),
		"kerotans_placed":      gorm.Expr("kerotans_placed + ?", stats.KerotansPlaced),
		"radio_uses":           gorm.Expr("radio_uses + ?", stats.RadioUses),
		"text_chat_uses":       gorm.Expr("text_chat_uses + ?", stats.TextChatUses),
		"cqc_attacks":          gorm.Expr("cqc_attacks + ?", stats.CQCAttacks),
		"cqc_attacks_received": gorm.Expr("cqc_attacks_received + ?", stats.CQCAttacksReceived),
		"head_shots":           gorm.Expr("head_shots + ?", stats.HeadShots),
		"head_shots_received":  gorm.Expr("head_shots_received + ?", stats.HeadShotsReceived),
		"team_wins":            gorm.Expr("team_wins + ?", stats.TeamWins),
		"kills_with_scorpion":  gorm.Expr("kills_with_scorpion + ?", stats.KillsWithScorpion),
		"kills_with_knife":     gorm.Expr("kills_with_knife + ?", stats.KillsWithKnife),
		"times_eaten":          gorm.Expr("times_eaten + ?", stats.TimesEaten),
		"rolls":                gorm.Expr("rolls + ?", stats.Rolls),
		"infrared_goggle_uses": gorm.Expr("infrared_goggle_uses + ?", stats.InfraredGoggleUses),
		"play_time":            gorm.Expr("play_time + ?", stats.PlayTime),
		"unknown3":             gorm.Expr("unknown3 + ?", stats.Unknown3),
	}

	// sqlite uses MAX(...) whereas others reserve MAX(...) for aggregates
	switch strings.ToLower(sess.DB.Dialector.Name()) {
	case "sqlite3":
		fallthrough
	case "sqlite":
		updates["kill_streak"] = gorm.Expr("MAX(kill_streak, ?)", stats.KillStreak)
		updates["death_streak"] = gorm.Expr("MAX(death_streak, ?)", stats.KillStreak)
	case "mssql":
		fallthrough
	case "postgres":
		fallthrough
	case "mysql":
		updates["kill_streak"] = gorm.Expr("GREATEST(kill_streak, ?)", stats.KillStreak)
		updates["death_streak"] = gorm.Expr("GREATEST(death_streak, ?)", stats.KillStreak)
	default:
		panic("unknown dialect: " + sess.DB.Dialector.Name())
	}

	currentRules := sess.GameState.Rules[sess.GameState.CurrentRound]
	tx := sess.DB.Model(&models.PlayerStats{})
	tx = tx.Where("user_id = ? AND mode = ? AND map = ?", UserID, currentRules.Mode, currentRules.Map)
	tx = tx.Updates(updates)

	if tx.Error != nil {
		// I don't believe we can get a ErrRecordNotFound with an UPDATE
		return tx.Error
	}

	if tx.RowsAffected < 2 {
		sess.Log.WithFields(logrus.Fields{
			"user_id": UserID,
			"map":     currentRules.Map,
			"mode":    currentRules.Mode,
		}).Info("Creating new stats for user")

		newStats := models.PlayerStats{
			UserID: uint(UserID),
			Mode:   currentRules.Mode,
			Map:    currentRules.Map,
			Period: types.PeriodWeekly,
		}
		newStats.FromHostReportedStats(stats)
		sess.DB.Create(&newStats)

		newStats = models.PlayerStats{
			UserID: uint(UserID),
			Mode:   currentRules.Mode,
			Map:    currentRules.Map,
			Period: types.PeriodAllTime,
		}
		newStats.FromHostReportedStats(stats)
		sess.DB.Create(&newStats)
	}
	return nil
}

// --- Packets ---

type ArgsHostPlayerStats struct {
	UserID types.UserID
	Stats  types.HostReportedStats
}

type ResponseHostPlayerStats types.ResponseErrorCode

func (r ResponseHostPlayerStats) Type() types.PacketType { return types.ServerHostPlayerStats }
