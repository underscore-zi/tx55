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
	handlers.Register(JoinHandler{})
}

// JoinHandler is called when a client selects a game lobby to join
// In response is sends the player two pieces of information:
// - First it sends an "overview" packet containing things like the display name, user id and such
// - Then it sends the players stats in two packets (one for alltime one for weekly stats)
type JoinHandler struct{}

func (h JoinHandler) Type() types.PacketType {
	return types.ClientJoinLobby
}

func (h JoinHandler) ArgumentTypes() []reflect.Type {
	return []reflect.Type{}
}

func (h JoinHandler) Handle(sess *session.Session, packet *packet.Packet) ([]types.Response, error) {
	var out []types.Response
	out = append(out, h.overviewPacket(sess))
	out = append(out, h.myStatsPackets(sess)...)
	return out, nil
}

func (h JoinHandler) overviewPacket(sess *session.Session) types.Response {
	out := ResponsePersonalOverview{
		Overview:    *sess.User.PlayerOverview(),
		Options:     sess.User.PlayerSettings.PlayerSpecificSettings(),
		FriendsList: [16]types.UserID{},
		BlockList:   [16]types.UserID{},
	}
	return out
}

func (h JoinHandler) myStatsPackets(sess *session.Session) []types.Response {
	var out []types.Response

	allTime, weekly, err := models.GetPlayerStats(sess.DB, types.UserID(sess.User.ID))
	if err != nil {
		return nil
	}

	out = append(out, ResponsePersonalStats(allTime), ResponsePersonalStats(weekly))
	return out
}

// --------------------

type ResponsePersonalOverview struct {
	Overview      types.PlayerOverview
	Options       types.PlayerSpecificSettings
	Magic2        [32]byte // All read as one group (ReadN)
	FriendsList   [16]types.UserID
	BlockList     [16]types.UserID
	OverallRank   uint32
	WeeklyRank    uint32
	OverallVSRank uint32
}

func (r ResponsePersonalOverview) Type() types.PacketType { return types.ServerPersonalOverview }

// --------------------

type ResponsePersonalStats types.PeriodStats

func (r ResponsePersonalStats) Type() types.PacketType { return types.ServerPersonalStats }
