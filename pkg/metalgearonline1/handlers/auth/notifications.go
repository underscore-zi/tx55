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
	handlers.Register(GetNotificationsHandler{})
}

type GetNotificationsHandler struct{}

func (h GetNotificationsHandler) Type() types.PacketType {
	return types.ClientGetNotifications
}

func (h GetNotificationsHandler) ArgumentTypes() []reflect.Type {
	return []reflect.Type{}
}

func (h GetNotificationsHandler) Handle(sess *session.Session, _ *packet.Packet) ([]types.Response, error) {
	var out []types.Response
	out = append(out, ResponseNotificationsStart{})

	var notifs []models.Notification
	if tx := sess.DB.Where("user_id = ?", sess.User.ID).Find(&notifs); tx.Error != nil {
		sess.Log.WithFields(sess.LogFields()).WithError(tx.Error).Error("failed to get notifications")
		return out, tx.Error
	}

	for _, n := range notifs {
		notif := ResponseNotificationEntry{
			ID:        uint32(n.ID),
			Important: n.IsImportant,
			HasRead:   n.HasRead,
			TimeStr:   [19]byte{},
			Title:     [64]byte{},
			Body:      [900]byte{},
		}
		copy(notif.TimeStr[:], []byte(n.CreatedAt.Format("2006-01-02 15:04:05")))
		copy(notif.Title[:], []byte(n.Title))
		copy(notif.Body[:], []byte(n.Body))
		out = append(out, notif)
	}

	out = append(out, ResponseNotificationsEnd{})
	return out, nil
}

// --- Packets ---
type ResponseNotificationsStart types.ResponseErrorCode

func (r ResponseNotificationsStart) Type() types.PacketType { return types.ServerNotificationsStart }

type ResponseNotificationsEnd types.ResponseErrorCode

func (r ResponseNotificationsEnd) Type() types.PacketType { return types.ServerNotificationsEnd }

type ResponseNotificationEntry struct {
	ID        uint32
	Important bool
	HasRead   bool
	TimeStr   [19]byte
	Title     [64]byte
	Body      [900]byte `packet:"truncate"`
}

func (r ResponseNotificationEntry) Type() types.PacketType { return types.ServerNotificationsEntry }
