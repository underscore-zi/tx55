package lobby

import (
	"bufio"
	"encoding/hex"
	"os"
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

func (h PlayerStatsHandler) Handle(sess *session.Session, packet *packet.Packet) ([]types.Response, error) {
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
		OverallVSRank: 0,
	}

	return overview
}

func ReadHexBytesFromFile(filename string) []byte {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	var out []byte
	for _, line := range lines {
		b, err := hex.DecodeString(line)
		if err != nil {
			panic(err)
		}
		out = append(out, b...)
	}

	return out
}

func (h PlayerStatsHandler) playerStats(sess *session.Session, args *ArgsGetPlayerStats) []types.Response {
	var out []types.Response

	all, weekly, err := models.GetPlayerStats(sess.DB, args.UserID)
	if err != nil {
		return nil
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
