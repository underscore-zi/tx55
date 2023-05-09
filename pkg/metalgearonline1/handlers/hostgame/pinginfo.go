package hostgame

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"tx55/pkg/metalgearonline1/handlers"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/metalgearonline1/session"
	"tx55/pkg/metalgearonline1/types"
	"tx55/pkg/packet"
)

func init() {
	handlers.Register(HostPingInfoHandler{})
}

type HostPingInfoHandler struct{}

func (h HostPingInfoHandler) Type() types.PacketType {
	return types.ClientHostPingInformation
}

func (h HostPingInfoHandler) ArgumentTypes() (out []reflect.Type) {
	return
}

func (h HostPingInfoHandler) parseArgs(packet *packet.Packet) (args ArgsHostPingInfo) {
	// Since this packet has a dynamic size we can't use binary.Read, instead we do it ourselves
	buf := bytes.NewBuffer((*packet).Data())
	args.GamePing = binary.BigEndian.Uint32(buf.Next(4))

	for i := 0; i < len(args.Pings) && buf.Len() >= 8; i++ {
		args.Pings[i] = types.PingInfo{
			UserID: types.UserID(binary.BigEndian.Uint32(buf.Next(4))),
			Ping:   binary.BigEndian.Uint32(buf.Next(4)),
		}
	}
	return
}

func (h HostPingInfoHandler) Handle(sess *session.Session, packet *packet.Packet) (out []types.Response, err error) {
	if !sess.IsHost() {
		out = append(out, ResponseHostPingInfo{ErrorCode: handlers.ErrNotHosting.Code})
		err = handlers.ErrNotHosting
		return
	}

	args := h.parseArgs(packet)
	for _, ping := range args.Pings {
		if ping.UserID == 0 {
			break
		}
		tx := sess.DB.Model(&models.GamePlayers{}).Where("game_id = ? AND user_id = ?", sess.GameState.GameID, ping.UserID).Update("ping", ping.Ping)
		if tx.Error != nil {
			out = append(out, ResponseHostPingInfo{ErrorCode: handlers.ErrDatabase.Code})
			err = tx.Error
			return
		}

		if tx.RowsAffected == 0 {
			out = append(out, ResponseHostPingInfo{ErrorCode: handlers.ErrNotFound.Code})
			err = handlers.ErrNotFound
			return
		}
	}

	out = append(out, ResponseHostPingInfo{ErrorCode: 0})

	return
}

// --- Packets ---

type ArgsHostPingInfo struct {
	GamePing uint32
	Pings    [9]types.PingInfo `packet:"truncate"`
}

type ResponseHostPingInfo types.ResponseErrorCode

func (r ResponseHostPingInfo) Type() types.PacketType { return types.ServerHostPingInformation }
