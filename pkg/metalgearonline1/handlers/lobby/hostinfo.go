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
	handlers.Register(GetHostInfoHandler{})
}

type GetHostInfoHandler struct{}

func (h GetHostInfoHandler) Type() types.PacketType {
	return types.ClientGetHostInfo
}

func (h GetHostInfoHandler) ArgumentTypes() (out []reflect.Type) {
	out = append(out, reflect.TypeOf(ArgsGetHostInfo{}))
	return
}

func (h GetHostInfoHandler) Handle(_ *session.Session, _ *packet.Packet) (out []types.Response, err error) {
	out = append(out, ResponseGetHostInfo{ErrorCode: handlers.ErrNotImplemented.Code})
	err = handlers.ErrNotImplemented
	return
}

func (h GetHostInfoHandler) HandleArgs(sess *session.Session, args *ArgsGetHostInfo) (out []types.Response, err error) {
	if args.GameID == 0 {
		out = append(out, ResponseGetHostInfo{ErrorCode: handlers.ErrInvalidArguments.Code})
		err = handlers.ErrInvalidArguments
		return
	}

	game := models.Game{}
	game.ID = uint(args.GameID)
	if tx := sess.DB.First(&game); tx.Error != nil {
		out = append(out, ResponseGetHostInfo{ErrorCode: handlers.ErrDatabase.Code})
		err = handlers.ErrDatabase
		return
	}
	if !game.CheckPassword(args.Password) {
		out = append(out, ResponseGetHostInfo{ErrorCode: handlers.ErrInvalidArguments.Code})
		err = handlers.ErrInvalidArguments
		return
	}

	hostConn := models.Connection{}
	hostConn.ID = game.ConnectionID
	if tx := sess.DB.First(&hostConn); tx.Error != nil {
		out = append(out, ResponseGetHostInfo{ErrorCode: handlers.ErrDatabase.Code})
		err = handlers.ErrDatabase
		return
	}

	hostInfo := ResponseGetHostInfo{
		ErrorCode:  0,
		RemotePort: hostConn.RemotePort,
		LocalPort:  hostConn.LocalPort,
	}

	copy(hostInfo.RemoteAddr[:], hostConn.RemoteAddr)
	copy(hostInfo.LocalAddr[:], hostConn.LocalAddr)

	out = append(out, hostInfo)
	return
}

// --- Packets ---

type ArgsGetHostInfo struct {
	GameID   types.GameID
	Password [16]byte `packet:"truncate"`
}

type ResponseGetHostInfo struct {
	ErrorCode  int32
	RemoteAddr [16]byte
	RemotePort uint16
	LocalAddr  [16]byte
	LocalPort  uint16
}

func (r ResponseGetHostInfo) Type() types.PacketType { return types.ServerHostInfo }
