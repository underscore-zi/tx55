package lobby

import (
	"reflect"
	"tx55/pkg/metalgearonline1/handlers"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/metalgearonline1/session"
	"tx55/pkg/metalgearonline1/types"
	"tx55/pkg/packet"
)

func init() {
	handlers.Register(PlayerStatsHandler{})
}

type PlayerStatsHandler struct{}

func (h PlayerStatsHandler) Type() types.PacketType {
	return types.ClientGetPlayerStats
}

func (h PlayerStatsHandler) ArgumentTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf(ArgsGetPlayerStats{}),
	}
}

func (h PlayerStatsHandler) Handle(_ *session.Session, _ *packet.Packet) ([]types.Response, error) {
	return nil, handlers.ErrNotImplemented
}

func (h PlayerStatsHandler) HandleArgs(sess *session.Session, args *ArgsGetPlayerStats) ([]types.Response, error) {
	var out []types.Response
	out = append(out, h.playerOverview(sess, args))
	out = append(out, h.playerStats(sess, args)...)
	return out, nil

}

func (h PlayerStatsHandler) playerOverview(sess *session.Session, args *ArgsGetPlayerStats) types.Response {
	var user models.User
	sess.DB.First(&user, args.UserID)
	if user.ID == 0 {
		return ResponsePlayerStatsOverview{
			ErrorCode: uint32(handlers.ErrNotFound.Code),
		}
	}

	overview := ResponsePlayerStatsOverview{
		ErrorCode:     0,
		Overview:      *user.PlayerOverview(),
		OverallRank:   uint32(user.OverallRank),
		WeeklyRank:    uint32(user.WeeklyRank),
		OverallVSRank: uint32(user.VsRatingRank),
	}

	var list []models.UserList
	sess.DB.Where("user_id = ?", args.UserID).Find(&list)
	var indexF, indexB int
	for _, l := range list {
		switch types.UserListType(l.ListType) {
		case types.UserListFriends:
			if indexF < len(overview.FriendsList) {
				overview.FriendsList[indexF] = types.UserID(l.EntryID)
				indexF++
			}
		case types.UserListBlocked:
			if indexB < len(overview.BlockedList) {
				overview.BlockedList[indexB] = types.UserID(l.EntryID)
				indexB++
			}
		}
	}

	// Expand the blocklist to catches alternate accounts
	// we do this by getting all the IDs that have an IP in common with this session's user
	// then if one of the alts is in the blocklist we replace it with the current user's ID
	// This only really matters for the stats request when joining a game, as it'll trip the
	// blocklist check. Normally you can't reach this without a session, but add the check anyway.
	if len(sess.SharedIds) > 0 {
		alts := make(map[types.UserID]bool)
		for _, alt := range sess.SharedIds {
			alts[types.UserID(alt)] = true
		}

		match := -1
		for i, id := range overview.BlockedList {
			if alts[id] {
				match = i
				break
			}
		}
		if match > -1 {
			overview.BlockedList[match] = types.UserID(user.ID)
		}
	}

	return overview
}

func (h PlayerStatsHandler) playerStats(sess *session.Session, args *ArgsGetPlayerStats) []types.Response {
	var out []types.Response
	var l = sess.LogEntry().WithField("player_id", args.UserID)

	all, weekly, err := models.GetPlayerStats(sess.DB, args.UserID)
	if err != nil {
		l.WithError(err).Error("failed to get player stats")
	}
	out = append(out, ResponsePlayerStats{
		Stats: all,
	})
	out = append(out, ResponsePlayerStats{
		Stats: weekly,
	})

	return out
}

// --------------------

type ArgsGetPlayerStats struct {
	UserID types.UserID
}

// --------------------

func (r ResponsePlayerStatsOverview) Type() types.PacketType { return types.ServerPlayerStatsOverview }

type ResponsePlayerStatsOverview struct {
	ErrorCode     uint32
	Overview      types.PlayerOverview
	FriendsList   [16]types.UserID
	BlockedList   [16]types.UserID
	OverallRank   uint32
	WeeklyRank    uint32
	OverallVSRank uint32
}

// --------------------

func (r ResponsePlayerStats) Type() types.PacketType { return types.ServerPlayerStats }

type ResponsePlayerStats struct {
	ErrorCode uint32
	Stats     types.PeriodStats
}
