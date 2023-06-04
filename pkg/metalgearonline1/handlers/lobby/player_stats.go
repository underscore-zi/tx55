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
