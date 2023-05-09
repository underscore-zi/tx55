package lobby

import (
	"gorm.io/gorm"
	"reflect"
	"tx55/pkg/metalgearonline1/handlers"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/metalgearonline1/session"
	"tx55/pkg/metalgearonline1/types"
	"tx55/pkg/packet"
)

func init() {
	handlers.Register(GetUserListHandler{})
}

type GetUserListHandler struct{}

func (h GetUserListHandler) Type() types.PacketType {
	return types.ClientGetUserList
}

func (h GetUserListHandler) ArgumentTypes() []reflect.Type {
	return []reflect.Type{reflect.TypeOf(ArgsGetUserList{})}
}

func (h GetUserListHandler) Handle(_ *session.Session, _ *packet.Packet) ([]types.Response, error) {
	return nil, handlers.ErrNotImplemented
}

func (h GetUserListHandler) HandleArgs(sess *session.Session, args *ArgsGetUserList) (out []types.Response, err error) {
	out = append(out, ResponseUserListStart{})

	var list []models.UserList
	err = sess.DB.Where("user_id = ? AND list_type = ?", sess.User.ID, args.ListType).Joins("Entry").Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		out = append(out, ResponseUserListEnd{})
		return
	}

	for _, entry := range list {
		newEntry := types.UserListEntry{
			UserID: types.UserID(entry.Entry.ID),
		}
		copy(newEntry.DisplayName[:], entry.Entry.DisplayName[:])
		out = append(out, ResponseUserListEntry{User: newEntry})
	}

	out = append(out, ResponseUserListEnd{})
	return out, nil
}

// --- Packets ---
type ArgsGetUserList struct {
	ListType types.UserListType
}

type ResponseUserListStart types.ResponseEmpty

func (r ResponseUserListStart) Type() types.PacketType { return types.ServerUserListStart }

type ResponseUserListEnd types.ResponseEmpty

func (r ResponseUserListEnd) Type() types.PacketType { return types.ServerUserListEnd }

type ResponseUserListEntry struct {
	User types.UserListEntry
}

func (r ResponseUserListEntry) Type() types.PacketType { return types.ServerUserListEntry }
