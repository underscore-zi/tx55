package auth

import (
	"reflect"
	"tx55/pkg/metalgearonline1/handlers"
	"tx55/pkg/metalgearonline1/models"
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

func (h NotificationReadReceiptHandler) Handle(_ *session.Session, _ *packet.Packet) ([]types.Response, error) {
	return nil, handlers.ErrNotImplemented
}

func (h NotificationReadReceiptHandler) HandleArgs(sess *session.Session, args *ArgsNotificationRead) ([]types.Response, error) {
	var out []types.Response

	q := sess.DB.Model(&models.Notification{}).Where("id = ? and user_id=?", args.ID, sess.User.ID)
	q.Update("has_read", true)

	out = append(out, ResponseUnknownNotification{ErrorCode: 0})
	return out, nil
}

// --- Packets ---

type ArgsNotificationRead struct {
	ID uint32
}
type ResponseUnknownNotification types.ResponseErrorCode

func (r ResponseUnknownNotification) Type() types.PacketType {
	return types.ServerNotificationReadReceipt
}
