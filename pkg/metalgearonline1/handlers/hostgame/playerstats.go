package hostgame

import (
	"reflect"
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

	var game *models.Game
	game, err = sess.Game()
	if err != nil {
		out = append(out, ResponseHostPlayerStats{ErrorCode: handlers.ErrDatabase.Code})
		return
	}

	rules := game.GameOptions.Rules[game.CurrentRound]

	// Make sure this is actually a player in the current game
	if player, found := game.FindPlayer(sess.DB, uint(args.UserID)); !found {
		out = append(out, ResponseHostPlayerStats{ErrorCode: handlers.ErrInvalidArguments.Code})
		err = handlers.ErrInvalidArguments
		return
	} else {
		player.Kills += args.Stats.Kills
		player.Deaths += args.Stats.Deaths
		player.Seconds += args.Stats.PlayTime
		player.Score += args.Stats.Points
		if tx := sess.DB.Save(player); tx.Error != nil {
			out = append(out, ResponseHostPlayerStats{ErrorCode: handlers.ErrDatabase.Code})
			err = tx.Error
			return
		}
	}

	// We track each stats report four times,
	// For each period (weekly and all-time) we track stats as map specific and all maps to avoid annoying calculations

	var stats models.PlayerStats

	// Weekly, map specific
	stats = models.PlayerStats{}
	if tx := sess.DB.FirstOrInit(&stats, map[string]interface{}{
		"user_id": uint(args.UserID),
		"period":  types.PeriodWeekly,
		"mode":    rules.Mode,
		"map":     rules.Map,
	}); tx.Error != nil {
		out = append(out, ResponseHostPlayerStats{ErrorCode: handlers.ErrDatabase.Code})
		err = tx.Error
		return
	}
	stats.AddStats(args.Stats)
	if tx := sess.DB.Save(&stats); tx.Error != nil {
		out = append(out, ResponseHostPlayerStats{ErrorCode: handlers.ErrDatabase.Code})
		err = tx.Error
		return
	}

	// All time, map specific
	stats = models.PlayerStats{}
	if tx := sess.DB.FirstOrInit(&stats, map[string]interface{}{
		"user_id": uint(args.UserID),
		"period":  types.PeriodAllTime,
		"mode":    rules.Mode,
		"map":     rules.Map,
	}); tx.Error != nil {
		out = append(out, ResponseHostPlayerStats{ErrorCode: handlers.ErrDatabase.Code})
		err = tx.Error
		return
	}
	stats.AddStats(args.Stats)
	if tx := sess.DB.Save(&stats); tx.Error != nil {
		out = append(out, ResponseHostPlayerStats{ErrorCode: handlers.ErrDatabase.Code})
		err = tx.Error
		return
	}

	// Also need to track it in the game player stats
	out = append(out, ResponseHostPlayerStats{ErrorCode: 0})
	return

}

// --- Packets ---

type ArgsHostPlayerStats struct {
	UserID types.UserID
	Stats  types.HostReportedStats
}

type ResponseHostPlayerStats types.ResponseErrorCode

func (r ResponseHostPlayerStats) Type() types.PacketType { return types.ServerHostPlayerStats }
