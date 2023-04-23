package types

import "fmt"

type PacketType uint16
type UserID uint32

// Response types should also be serializable by the binary package, so all types need to be fixed size
type Response interface {
	Type() PacketType
}

const (
	ClientDisconnect PacketType = 0x0003
	ServerDisconnect PacketType = 0x0004
	ClientPing       PacketType = 0x0005
	ServerPing       PacketType = 0x0005

	// --- Lobby Service Commands ---

	ClientGetLobbyList   PacketType = 0x2005
	ServerLobbyListStart PacketType = 0x2002
	ServerLobbyListEntry PacketType = 0x2003
	ServerLobbyListEnd   PacketType = 0x2004

	ClientGetNewsList   PacketType = 0x2008
	ServerNewsListStart PacketType = 0x2009
	ServerNewsListEntry PacketType = 0x200A
	ServerNewsListEnd   PacketType = 0x200B

	// --- Auth Service Commands ---

	ClientGetNonce PacketType = 0x3001
	ServerNonce    PacketType = 0x3002

	ClientLogin PacketType = 0x3003
	ServerLogin PacketType = 0x3004

	ClientGetSessionInfo PacketType = 0x3040
	ServerSessionInfo    PacketType = 0x3041

	ClientGetNotifications   PacketType = 0x3042
	ServerNotificationsStart PacketType = 0x3043
	ServerNotificationsEntry PacketType = 0x3044
	ServerNotificationsEnd   PacketType = 0x3045

	ClientNotificationReadReceipt PacketType = 0x3046
	ServerNotificationReadReceipt PacketType = 0x3047

	// --- Game Service Commands ---

	ClientJoinLobby        PacketType = 0x4100
	ServerPersonalOverview PacketType = 0x4101
	ServerPersonalStats    PacketType = 0x4104

	ClientGetPlayerStats      PacketType = 0x4102
	ServerPlayerStatsOverview PacketType = 0x4103
	ServerPlayerStats         PacketType = 0x4105

	ClientUpdatePlayerSettings PacketType = 0x4110
	ServerPlayerSettings       PacketType = 0x4111

	ClientGetFilteredGameList PacketType = 0x4112

	ClientGetCreateGameSettings PacketType = 0x4304
	ServerCreateGameSettings    PacketType = 0x4305

	ClientGetGameList   PacketType = 0x4300
	ServerGameListStart PacketType = 0x4301
	ServerGameListEntry PacketType = 0x4302
	ServerGameListEnd   PacketType = 0x4303

	ClientCreateGame PacketType = 0x4310
	ServerCreateGame PacketType = 0x4311

	ClientGetGameInfo PacketType = 0x4312
	ServerGameInfo    PacketType = 0x4313

	ClientHostReadyToCreate PacketType = 0x4316
	ServerHostReadyToCreate PacketType = 0x4317

	ClientGetHostInfo PacketType = 0x4320
	ServerHostInfo    PacketType = 0x4321

	ClientPlayerFailedToJoinHost PacketType = 0x4322
	ServerPlayerFailedToJoinHost PacketType = 0x4323

	ClientHostPlayerJoin PacketType = 0x4340
	ServerHostPlayerJoin PacketType = 0x4341

	ClientHostPlayerLeave PacketType = 0x4342
	ServerHostPlayerLeave PacketType = 0x4343

	ClientHostPlayerJoinTeam PacketType = 0x4344
	ServerHostPlayerJoinTeam PacketType = 0x4345

	ClientHostPlayerKicked PacketType = 0x4346
	ServerHostPlayerKicked PacketType = 0x4347

	ClientHostQuitGame PacketType = 0x4380
	ServerHostQuitGame PacketType = 0x4381

	ClientHostPlayerStats PacketType = 0x4390
	ServerHostPlayerStats PacketType = 0x4391

	ClientHostNewRound PacketType = 0x4392
	ServerHostNewRound PacketType = 0x4393

	ClientHost4394 PacketType = 0x4394
	ServerHost4394 PacketType = 0x4395

	ClientHostPingInformation PacketType = 0x4398
	ServerHostPingInformation PacketType = 0x4399

	ClientGetUserList   PacketType = 0x4580
	ServerUserListStart PacketType = 0x4581
	ServerUserListEntry PacketType = 0x4582
	ServerUserListEnd   PacketType = 0x4583

	ClientAddUserToList      PacketType = 0x4500
	ServerAddUserToList      PacketType = 0x4502
	ClientRemoveUserFromList PacketType = 0x4510
	ServerRemoveUserFromList PacketType = 0x4512

	ClientReportConnectionInfo PacketType = 0x4700
	ServerReportConnectionInfo PacketType = 0x4701
)

type ResponseErrorCode struct {
	ErrorCode int32
}

type ResponseEmpty struct{}

func (p PacketType) String() string {
	switch p {
	case ClientDisconnect:
		return "ClientDisconnect"
	case ServerDisconnect:
		return "ServerDisconnect"
	case ClientPing:
		// Shared with ServerPing
		return "Ping"
	case ClientGetLobbyList:
		return "ClientGetLobbyList"
	case ServerLobbyListStart:
		return "ServerLobbyListStart"
	case ServerLobbyListEntry:
		return "ServerLobbyListEntry"
	case ServerLobbyListEnd:
		return "ServerLobbyListEnd"
	case ClientGetNewsList:
		return "ClientGetNewsList"
	case ServerNewsListStart:
		return "ServerNewsListStart"
	case ServerNewsListEntry:
		return "ServerNewsListEntry"
	case ServerNewsListEnd:
		return "ServerNewsListEnd"
	case ClientGetNonce:
		return "ClientGetNonce"
	case ServerNonce:
		return "ServerNonce"
	case ClientLogin:
		return "ClientLogin"
	case ServerLogin:
		return "ServerLogin"
	case ClientGetSessionInfo:
		return "ClientGetSessionInfo"
	case ServerSessionInfo:
		return "ServerSessionInfo"
	case ClientGetNotifications:
		return "ClientGetNotifications"
	case ServerNotificationsStart:
		return "ServerNotificationsStart"
	case ServerNotificationsEntry:
		return "ServerNotificationsEntry"
	case ServerNotificationsEnd:
		return "ServerNotificationsEnd"
	case ClientJoinLobby:
		return "ClientJoinLobby"
	case ServerPersonalOverview:
		return "ServerPersonalOverview"
	case ServerPersonalStats:
		return "ServerPersonalStats"
	case ClientGetPlayerStats:
		return "ClientGetPlayerStats"
	case ServerPlayerStatsOverview:
		return "ServerPlayerStatsOverview"
	case ServerPlayerStats:
		return "ServerPlayerStats"
	case ClientUpdatePlayerSettings:
		return "ClientUpdatePlayerSettings"
	case ServerPlayerSettings:
		return "ServerPlayerSettings"
	case ClientGetFilteredGameList:
		return "ClientGetFilteredGameList"
	case ClientGetCreateGameSettings:
		return "ClientGetCreateGameSettings"
	case ServerCreateGameSettings:
		return "ServerCreateGameSettings"
	case ClientGetGameList:
		return "ClientGetGameList"
	case ServerGameListStart:
		return "ServerGameListStart"
	case ServerGameListEntry:
		return "ServerGameListEntry"
	case ServerGameListEnd:
		return "ServerGameListEnd"
	case ClientCreateGame:
		return "ClientCreateGame"
	case ServerCreateGame:
		return "ServerCreateGame"
	case ClientGetGameInfo:
		return "ClientGetGameInfo"
	case ServerGameInfo:
		return "ServerGameInfo"
	case ClientHostReadyToCreate:
		return "ClientHostReadyToCreate"
	case ServerHostReadyToCreate:
		return "ServerHostReadyToCreate"
	case ClientGetHostInfo:
		return "ClientGetHostInfo"
	case ServerHostInfo:
		return "ServerHostInfo"
	case ClientPlayerFailedToJoinHost:
		return "ClientPlayerFailedToJoinHost"
	case ServerPlayerFailedToJoinHost:
		return "ServerPlayerFailedToJoinHost"
	case ClientHostPlayerJoin:
		return "ClientHostPlayerJoin"
	case ServerHostPlayerJoin:
		return "ServerHostPlayerJoin"
	case ClientHostPlayerLeave:
		return "ClientHostPlayerLeave"
	case ServerHostPlayerLeave:
		return "ServerHostPlayerLeave"
	case ClientHostPlayerJoinTeam:
		return "ClientHostPlayerJoinTeam"
	case ServerHostPlayerJoinTeam:
		return "ServerHostPlayerJoinTeam"
	case ClientHostPlayerKicked:
		return "ClientHostPlayerKicked"
	case ServerHostPlayerKicked:
		return "ServerHostPlayerKicked"
	case ClientHostQuitGame:
		return "ClientHostQuitGame"
	case ServerHostQuitGame:
		return "ServerHostQuitGame"
	case ClientHostPlayerStats:
		return "ClientHostPlayerStats"
	case ServerHostPlayerStats:
		return "ServerHostPlayerStats"
	case ClientHostNewRound:
		return "ClientHostNewRound"
	case ServerHostNewRound:
		return "ServerHostNewRound"
	case ClientHost4394:
		return "ClientHost4394"
	case ServerHost4394:
		return "ServerHost4394"
	case ClientHostPingInformation:
		return "ClientHostPingInformation"
	case ServerHostPingInformation:
		return "ServerHostPingInformation"
	case ClientGetUserList:
		return "ClientGetUserList"
	case ServerUserListStart:
		return "ServerUserListStart"
	case ServerUserListEntry:
		return "ServerUserListEntry"
	case ServerUserListEnd:
		return "ServerUserListEnd"
	case ClientAddUserToList:
		return "ClientAddUserToList"
	case ServerAddUserToList:
		return "ServerAddUserToList"
	case ClientRemoveUserFromList:
		return "ClientRemoveUserFromList"
	case ServerRemoveUserFromList:
		return "ServerRemoveUserFromList"
	case ClientReportConnectionInfo:
		return "ClientReportConnectionInfo"
	case ServerReportConnectionInfo:
		return "ServerReportConnectionInfo"
	default:
		return fmt.Sprintf("PacketType(0x%04x)", uint(p))
	}
}
