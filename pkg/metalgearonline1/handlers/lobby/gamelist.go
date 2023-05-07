package lobby

import (
	"gorm.io/gorm/clause"
	"reflect"
	"tx55/pkg/metalgearonline1/handlers"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/metalgearonline1/session"
	"tx55/pkg/metalgearonline1/types"
	"tx55/pkg/metalgearonline1/types/bitfield"
	"tx55/pkg/packet"
)

const maxGamesPerPacket = 12

func init() {
	handlers.Register(GetGameListHandler{})
}

// GetGameListHandler is called when a user select "Join Game" and it displays the game list
// It uses a "list" response pattern. It has a dedicated start and end packet. In between those
// you can send multiple packets containing an array of actual content. Or you can bundle it all
// up into a single packet with multiple entries.
type GetGameListHandler struct{}

func (h GetGameListHandler) Type() types.PacketType {
	return types.ClientGetGameList
}

func (h GetGameListHandler) ArgumentTypes() []reflect.Type {
	return []reflect.Type{}
}

func (h GetGameListHandler) getGames(s *session.Session) ([]ResponseGameListEntry, error) {
	var out []ResponseGameListEntry
	var games []models.Game
	if tx := s.DB.Preload(clause.Associations).Where("lobby_id = ?", s.LobbyID).Find(&games); tx.Error != nil {
		return out, tx.Error
	}
	for _, game := range games {
		entry := ResponseGameListEntry{
			ID:                  uint32(game.ID),
			HasPassword:         game.GameOptions.HasPassword,
			IsHostOnly:          game.GameOptions.IsHostOnly,
			CurrentRules:        game.GameOptions.Rules[game.CurrentRound],
			WeaponRestriction:   game.GameOptions.WeaponRestriction,
			MaxPlayers:          game.GameOptions.MaxPlayers,
			PlayerCount:         byte(len(game.Players)),
			Options:             game.GameOptions.Bitfield,
			VSRatingRestriction: game.GameOptions.RatingRestriction,
			VSRating:            game.GameOptions.Rating,
		}
		if len(game.Players) > 0 {
			entry.Ping = game.Players[0].Ping
		}
		copy(entry.Name[:], game.GameOptions.Name)
		out = append(out, entry)
	}

	return out, nil
}

func (h GetGameListHandler) Handle(sess *session.Session, packet *packet.Packet) ([]types.Response, error) {
	var out []types.Response
	out = append(out, ResponseGameListStart{})
	games, err := h.getGames(sess)
	if err == nil {
		for i := 0; i < len(games); i += maxGamesPerPacket {
			end := i + maxGamesPerPacket
			if end > len(games) {
				end = len(games)
			}
			resp := ResponseGameList{}
			copy(resp.Games[:], games[i:end])
			out = append(out, resp)
		}
		out = append(out, ResponseGameListEnd{})
	}
	return out, err
}

// --- Packets ---

type ResponseGameListStart types.ResponseEmpty

func (r ResponseGameListStart) Type() types.PacketType { return types.ServerGameListStart }

type ResponseGameListEnd types.ResponseEmpty

func (r ResponseGameListEnd) Type() types.PacketType { return types.ServerGameListEnd }

type ResponseGameListEntry struct {
	ID                  uint32
	Name                [16]byte
	HasPassword         bool
	IsHostOnly          bool
	CurrentRules        types.GameRules
	WeaponRestriction   types.WeaponRestrictions
	MaxPlayers          uint8
	Options             bitfield.GameSettings
	PlayerCount         uint8
	Ping                uint32
	FriendOrBlocked     uint8
	VSRatingRestriction types.VSRatingRestriction
	VSRating            uint32
	// Based on a leaked struct from MGO2, but these are not displayed so won't implement them
	WinStreak uint16
	WinnerID  types.UserID
}

func (r ResponseGameList) Type() types.PacketType { return types.ServerGameListEntry }

type ResponseGameList struct {
	Games [maxGamesPerPacket]ResponseGameListEntry `packet:"truncate"`
}
