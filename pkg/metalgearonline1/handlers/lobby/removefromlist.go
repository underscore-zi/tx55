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
	handlers.Register(RemoveUserFromList{})
}

type RemoveUserFromList struct{}

func (h RemoveUserFromList) Type() types.PacketType {
	return types.ClientRemoveUserFromList
}

func (h RemoveUserFromList) ArgumentTypes() (out []reflect.Type) {
	out = append(out, reflect.TypeOf(ArgsRemoveFromList{}))
	return
}

func (h RemoveUserFromList) Handle(_ *session.Session, _ *packet.Packet) (out []types.Response, err error) {
	out = append(out, ResponseRemoveUserFromListError{ErrorCode: handlers.ErrNotImplemented.Code})
	err = handlers.ErrNotImplemented
	return
}

func (h RemoveUserFromList) HandleArgs(s *session.Session, args *ArgsRemoveFromList) (out []types.Response, err error) {
	var list []models.UserList
	if tx := s.DB.Where("user_id = ? AND list_type = ?", s.User.ID, args.ListType).Find(&list); tx.Error != nil {
		out = append(out, ResponseAddUserToListError{ErrorCode: handlers.ErrDatabase.Code})
		err = tx.Error
		return
	}

	tx := s.DB.Delete(&models.UserList{}, "user_id = ? AND list_type = ? AND entry_id = ?", s.User.ID, args.ListType, args.UserID)
	if tx.Error != nil {
		out = append(out, ResponseAddUserToListError{ErrorCode: handlers.ErrDatabase.Code})
		err = tx.Error
		return
	}

	out = append(out, ResponseRemoveUserFromList{
		ErrorCode: 0,
		ListType:  args.ListType,
		UserID:    args.UserID,
	})

	return out, err
}

// --- Packets ---

type ArgsRemoveFromList struct {
	ListType types.UserListType
	UserID   types.UserID
}

type ResponseRemoveUserFromListError types.ResponseErrorCode

func (r ResponseRemoveUserFromListError) Type() types.PacketType {
	return types.ServerRemoveUserFromList
}

type ResponseRemoveUserFromList struct {
	ErrorCode int32
	ListType  types.UserListType
	UserID    types.UserID
}

func (r ResponseRemoveUserFromList) Type() types.PacketType {
	return types.ServerRemoveUserFromList
}
