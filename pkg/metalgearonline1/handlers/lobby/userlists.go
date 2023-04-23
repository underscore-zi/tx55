package lobby

import (
	"gorm.io/gorm/clause"
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

func (h GetUserListHandler) Handle(sess *session.Session, packet *packet.Packet) ([]types.Response, error) {
	return nil, handlers.ErrNotImplemented
}

func (h GetUserListHandler) HandleArgs(sess *session.Session, args *ArgsGetUserList) (out []types.Response, err error) {
	out = append(out, ResponseUserListStart{})
	switch args.ListType {
	case types.UserListFriends:
		var list []models.Friend
		if tx := sess.DB.Where("user_id = ?", sess.User.ID).Preload(clause.Associations).Find(&list); tx.Error != nil {
			out = append(out, ResponseUserListEnd{})
			err = tx.Error
			return
		}
		for _, entry := range list {
			newEntry := types.UserListEntry{
				UserID: types.UserID(entry.FriendID),
			}
			copy(newEntry.DisplayName[:], entry.Friend.DisplayName[:])
			out = append(out, ResponseUserListEntry{User: newEntry})
		}
	case types.UserListBlocked:
		var list []models.Blocked
		if tx := sess.DB.Where("user_id = ?", sess.User.ID).Preload(clause.Associations).Find(&list); tx.Error != nil {
			out = append(out, ResponseUserListEnd{})
			err = tx.Error
			return
		}
		for _, entry := range list {
			newEntry := types.UserListEntry{
				UserID: types.UserID(entry.BlockedID),
			}
			copy(newEntry.DisplayName[:], entry.Blocked.DisplayName[:])
			out = append(out, ResponseUserListEntry{User: newEntry})
		}
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
