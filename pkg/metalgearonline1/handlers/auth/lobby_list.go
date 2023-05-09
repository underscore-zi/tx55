package auth

import (
	"reflect"
	"tx55/pkg/metalgearonline1/handlers"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/metalgearonline1/session"
	"tx55/pkg/metalgearonline1/types"
	"tx55/pkg/packet"
)

const (
	// MaxLobbies is the maximum number of lobbies the game will display
	MaxLobbies = 18
)

func init() {
	handlers.Register(GetLobbyListHandler{})
}

// GetLobbyListHandler is called just before hte game presents the actual lobby list. Either on first connection
// or after disconnecting/backing out of a lobby it will be called again.
type GetLobbyListHandler struct{}

func (h GetLobbyListHandler) Type() types.PacketType {
	return types.ClientGetLobbyList
}

func (h GetLobbyListHandler) ArgumentTypes() []reflect.Type {
	return []reflect.Type{}
}

func (h GetLobbyListHandler) Handle(sess *session.Session, _ *packet.Packet) ([]types.Response, error) {
	var out []types.Response
	out = append(out, ResponseLobbyListStart{})

	var list ResponseLobbyList
	copy(list.Lobbies[:], h.getLobbies(sess))
	out = append(out, list)

	out = append(out, ResponseLobbyListEnd{})
	return out, nil
}

func (h GetLobbyListHandler) getLobbies(sess *session.Session) []LobbyListEntry {
	var out []LobbyListEntry
	var lobbies []models.Lobby
	_ = sess.DB.Find(&lobbies)

	for _, lobby := range lobbies {
		newLobby := LobbyListEntry{
			ID:   lobby.ID,
			Type: lobby.Type,
			// TODO: Calculate player count
			Port:    lobby.Port,
			Players: lobby.Players,
			GID:     uint16(lobby.ID),
		}
		copy(newLobby.Name[:], lobby.Name)
		copy(newLobby.IP[:], lobby.IP)
		out = append(out, newLobby)
	}

	return out
}

// --- Packets ---
type ResponseLobbyListStart types.ResponseEmpty

func (r ResponseLobbyListStart) Type() types.PacketType { return types.ServerLobbyListStart }

type ResponseLobbyListEnd types.ResponseEmpty

func (r ResponseLobbyListEnd) Type() types.PacketType { return types.ServerLobbyListEnd }

type ResponseLobbyList struct {
	Lobbies [MaxLobbies]LobbyListEntry `packet:"truncate"`
}

func (r ResponseLobbyList) Type() types.PacketType { return types.ServerLobbyListEntry }

type LobbyListEntry struct {
	ID      uint32
	Type    types.LobbyType
	Name    [16]byte
	IP      [15]byte
	Port    uint16
	Players uint16
	GID     uint16
}
