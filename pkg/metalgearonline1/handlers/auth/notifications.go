package auth

import (
	"fmt"
	"reflect"
	"time"
	"tx55/pkg/metalgearonline1/handlers"
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

func (h GetNotificationsHandler) Handle(sess *session.Session, packet *packet.Packet) ([]types.Response, error) {
	var out []types.Response
	out = append(out, ResponseNotificationsStart{})

	notifTime := time.Now().Add(-time.Hour * 24 * 60)

	notif := ResponseNotificationEntry{
		Time:  uint32(notifTime.Unix()),
		IsNew: 1,
	}

	copy(notif.TimeStr[:], []byte(notifTime.Format("2006-01-02 15:04:05")))
	copy(notif.Title[:], []byte("Test Notification"))
	copy(notif.Body[:], []byte("This is a test notification."))
	//out = append(out, notif)
	out = append(out, ResponseNotificationsEnd{})
	fmt.Println("sending notf")
	return out, nil
}

// --- Packets ---
type ResponseNotificationsStart types.ResponseErrorCode

func (r ResponseNotificationsStart) Type() types.PacketType { return types.ServerNotificationsStart }

type ResponseNotificationsEnd types.ResponseErrorCode

func (r ResponseNotificationsEnd) Type() types.PacketType { return types.ServerNotificationsEnd }

type ResponseNotificationEntry struct {
	Time    uint32
	IsNew   uint8
	U3      uint8
	TimeStr [19]byte
	Title   [64]byte
	Body    [900]byte `packet:"truncate"`
}

func (r ResponseNotificationEntry) Type() types.PacketType { return types.ServerNotificationsEntry }
