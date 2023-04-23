package auth

import (
	"reflect"
	"tx55/pkg/metalgearonline1/handlers"
	"tx55/pkg/metalgearonline1/session"
	"tx55/pkg/metalgearonline1/types"
	"tx55/pkg/packet"
)

func init() {
	handlers.Register(NotificationReadReceiptHandler{})
}

type NotificationReadReceiptHandler struct{}

func (h NotificationReadReceiptHandler) Type() types.PacketType {
	return types.ClientNotificationReadReceipt
}

func (h NotificationReadReceiptHandler) ArgumentTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf(ArgsNotificationRead{}),
	}
}

func (h NotificationReadReceiptHandler) Handle(sess *session.Session, p *packet.Packet) ([]types.Response, error) {
	return nil, handlers.ErrNotImplemented
}

func (h NotificationReadReceiptHandler) HandleArgs(sess *session.Session, args *ArgsNotificationRead) ([]types.Response, error) {
	var out []types.Response
	out = append(out, ResponseUnknownNotification{ErrorCode: 0})
	return out, nil
}

// --- Packets ---

type ArgsNotificationRead struct {
	ReadAt uint32
}
type ResponseUnknownNotification types.ResponseErrorCode

func (r ResponseUnknownNotification) Type() types.PacketType {
	return types.ServerNotificationReadReceipt
}
