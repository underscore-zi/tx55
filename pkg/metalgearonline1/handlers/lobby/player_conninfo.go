package lobby

import (
	"reflect"
	"time"
	"tx55/pkg/metalgearonline1/handlers"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/metalgearonline1/session"
	"tx55/pkg/metalgearonline1/types"
	"tx55/pkg/packet"
)

func init() {
	handlers.Register(ReportConnectionInfo{})
}

type ReportConnectionInfo struct{}

func (h ReportConnectionInfo) Type() types.PacketType {
	return types.ClientReportConnectionInfo
}

func (h ReportConnectionInfo) ArgumentTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf(ArgsReportConnectionInfo{}),
	}
}

func (h ReportConnectionInfo) Handle(_ *session.Session, _ *packet.Packet) ([]types.Response, error) {
	return nil, handlers.ErrNotImplemented
}

func (h ReportConnectionInfo) HandleArgs(sess *session.Session, args *ArgsReportConnectionInfo) ([]types.Response, error) {
	tx := sess.DB.Model(&models.Connection{})
	tx = tx.Where("user_id = ?", sess.User.ID)
	tx = tx.Where("remote_addr = ?", sess.IP)
	tx = tx.Where("remote_port = ?", args.RemotePort)
	tx = tx.Where("local_addr = ?", types.BytesToString(args.LocalAddr[:]))
	tx = tx.Where("local_port = ?", args.LocalPort)

	var conn models.Connection
	tx = tx.First(&conn)

	if conn.ID == 0 {
		conn.UserID = sess.User.ID
		conn.RemoteAddr = sess.IP
		conn.RemotePort = args.RemotePort
		conn.LocalAddr = types.BytesToString(args.LocalAddr[:])
		conn.LocalPort = args.LocalPort
		tx = sess.DB.Create(&conn)
	} else {
		tx = sess.DB.Model(&conn).Update("updated_at", time.Now())
	}

	if tx.Error != nil {
		return []types.Response{ResponseReportConnectionInfo{ErrorCode: handlers.ErrDatabase.Code}}, handlers.ErrDatabase
	}

	sess.ActiveConnection = conn

	return []types.Response{ResponseReportConnectionInfo{ErrorCode: 0}}, nil
}

// ----------------------

type ArgsReportConnectionInfo struct {
	RemotePort uint16
	LocalAddr  [16]byte
	LocalPort  uint16
	Unknown    uint16
}

// ----------------------

type ResponseReportConnectionInfo types.ResponseErrorCode

func (r ResponseReportConnectionInfo) Type() types.PacketType {
	return types.ServerReportConnectionInfo
}
