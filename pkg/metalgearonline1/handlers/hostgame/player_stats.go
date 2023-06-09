package hostgame

import (
	"fmt"
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

	go func() { _ = h.updatePlayerStats(sess, uint(args.UserID), args.Stats) }()
	out = append(out, ResponseHostPlayerStats{ErrorCode: 0})
	return
}

func (h HostPlayerStatsHandler) updatePlayerStats(sess *session.Session, UserID uint, stats types.HostReportedStats) error {
	if _, found := sess.GameState.Players[types.UserID(UserID)]; !found {
		return fmt.Errorf("player %d not found in game", UserID)
	}

	l := sess.LogEntry().WithFields(logrus.Fields{
		"stats_user_id": UserID,
		"mode":          sess.GameState.Rules[sess.GameState.CurrentRound].Mode,
		"map":           sess.GameState.Rules[sess.GameState.CurrentRound].Map,
		"kills":         stats.Kills,
		"deaths":        stats.Deaths,
		"points":        stats.Points,
		"vs_rating":     stats.VsRating,
	})

	// Update the overview stats for the game listing
	updates := map[string]interface{}{
		"kills":  gorm.Expr("kills + ?", stats.Kills),
		"deaths": gorm.Expr("deaths + ?", stats.Deaths),
		"score":  gorm.Expr("score + ?", stats.Points),
	}
	q := sess.DB.Model(&models.GamePlayers{}).Where("game_id = ? AND user_id = ?", sess.GameState.GameID, UserID)
	if rowCount := q.Updates(updates).RowsAffected; rowCount != 1 {
		l.Error("Attempting to update stats for a player that is not in the game")
		return nil
	}

	// If stats are disabled for this game don't process the the rest
	if !sess.GameState.CollectStats {
		return nil
	}

	hasZeroRating := stats.VsRating == 0
	hasNegativeValues := stats.Points < 0 || stats.Kills < 0 || stats.Deaths < 0
	hasNoPlayTime := stats.PlayTime == 0
	hasJoinedTeam := sess.GameState.Players[types.UserID(UserID)].After(sess.GameState.RoundStart)

	if hasZeroRating || hasNegativeValues {
		// Not sure why these happen, but they happen rarely enough that we just won't save these stats
		l.Debug("Ignoring stats with zero rating or negative values")
		return nil
	}

	if hasNoPlayTime {
		l.WithFields(logrus.Fields{
			"play_time":   stats.PlayTime,
			"updated_at":  sess.GameState.Players[types.UserID(UserID)],
			"round_start": sess.GameState.RoundStart,
		}).Warn("Stats with no play time!")
		return nil
	}

	if !hasJoinedTeam && stats.Points == 0 {
		// Trying to detect when someone just spectates the entire round so we don't count it as a round played
		l.Debug("Ignoring stats with no points and no team join")
		return nil
	}

	sess.DB.Model(&models.User{
		Model: gorm.Model{ID: UserID},
	}).Update("vs_rating", stats.VsRating)

	updates = map[string]interface{}{
		"kills":                gorm.Expr("kills + ?", stats.Kills),
		"deaths":               gorm.Expr("deaths + ?", stats.Deaths),
		"stuns":                gorm.Expr("stuns + ?", stats.Stuns),
		"stuns_received":       gorm.Expr("stuns_received + ?", stats.StunsReceived),
		"snake_frags":          gorm.Expr("snake_frags + ?", stats.SnakeFrags),
		"points":               gorm.Expr("points + ?", stats.Points),
		"suicides":             gorm.Expr("suicides + ?", stats.Suicides),
		"self_stuns":           gorm.Expr("self_stuns + ?", stats.SelfStuns),
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
	}

	// sqlite uses MAX(...) whereas others reserve MAX(...) for aggregates
	switch strings.ToLower(sess.DB.Dialector.Name()) {
	case "sqlite3":
		fallthrough
	case "sqlite":
		updates["kill_streak"] = gorm.Expr("MAX(kill_streak, ?)", stats.KillStreak)
		updates["death_streak"] = gorm.Expr("MAX(death_streak, ?)", stats.DeathStreak)
	case "mssql":
		fallthrough
	case "postgres":
		fallthrough
	case "mysql":
		updates["kill_streak"] = gorm.Expr("GREATEST(kill_streak, ?)", stats.KillStreak)
		updates["death_streak"] = gorm.Expr("GREATEST(death_streak, ?)", stats.DeathStreak)
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
		l.Info("Creating new stats for user")

		newStats := models.PlayerStats{
			UserID: UserID,
			Mode:   currentRules.Mode,
			Map:    currentRules.Map,
			Period: types.PeriodWeekly,
		}
		newStats.FromHostReportedStats(stats)
		sess.DB.Create(&newStats)

		newStats = models.PlayerStats{
			UserID: UserID,
			Mode:   currentRules.Mode,
			Map:    currentRules.Map,
			Period: types.PeriodAllTime,
		}
		newStats.FromHostReportedStats(stats)
		sess.DB.Create(&newStats)
	}

	l.Info("Updated user stats")
	return nil
}

// --- Packets ---

type ArgsHostPlayerStats struct {
	UserID types.UserID
	Stats  types.HostReportedStats
}

type ResponseHostPlayerStats types.ResponseErrorCode

func (r ResponseHostPlayerStats) Type() types.PacketType { return types.ServerHostPlayerStats }
