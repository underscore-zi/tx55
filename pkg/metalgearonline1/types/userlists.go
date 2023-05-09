package types

type GameID uint32
type LobbyID uint16

type UserListType byte

const (
	UserListFriends UserListType = 0
	UserListBlocked UserListType = 1
)

type UserListEntry struct {
	UserID      UserID
	DisplayName [16]byte
	LobbyID     LobbyID
	LobbyName   [16]byte
	GameID      GameID
	GameName    [16]byte
}

type LobbyType uint32

//goland:noinspection GoUnusedConst,GoUnusedConst
const (
	LobbyTypeGate    LobbyType = 0
	LobbyTypeAccount LobbyType = 1
	LobbyTypeGame    LobbyType = 2
)
