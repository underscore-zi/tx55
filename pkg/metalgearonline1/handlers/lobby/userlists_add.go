package lobby

import (
	"errors"
	"gorm.io/gorm"
	"reflect"
	"tx55/pkg/metalgearonline1/handlers"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/metalgearonline1/session"
	"tx55/pkg/metalgearonline1/types"
	"tx55/pkg/packet"
)

func init() {
	handlers.Register(AddUserToList{})
}

type AddUserToList struct{}

func (h AddUserToList) Type() types.PacketType {
	return types.ClientAddUserToList
}

func (h AddUserToList) ArgumentTypes() (out []reflect.Type) {
	out = append(out, reflect.TypeOf(ArgsAddToList{}))
	return
}

func (h AddUserToList) Handle(_ *session.Session, _ *packet.Packet) (out []types.Response, err error) {
	err = handlers.ErrNotImplemented
	out = append(out, ResponseAddUserToListError{ErrorCode: handlers.ErrNotImplemented.Code})
	return
}

func (h AddUserToList) HandleArgs(s *session.Session, args *ArgsAddToList) (out []types.Response, err error) {
	var list []models.UserList

	if tx := s.DB.Where("user_id = ? AND list_type = ?", s.User.ID, args.ListType).Find(&list); tx.Error != nil {
		out = append(out, ResponseAddUserToListError{ErrorCode: handlers.ErrDatabase.Code})
		err = tx.Error
		return
	}

	// This limit comes from in-game
	if len(list) >= 16 {
		out = append(out, ResponseAddUserToListError{ErrorCode: handlers.ErrInvalidArguments.Code})
		err = errors.New("list is full")
		return
	}

	// Make sure its legit user (and get the display name for later)
	var user models.User
	user.ID = uint(args.UserID)
	if tx := s.DB.First(&user); tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			out = append(out, ResponseAddUserToListError{ErrorCode: handlers.ErrNotFound.Code})
			err = handlers.ErrNotFound
		} else {
			out = append(out, ResponseAddUserToListError{ErrorCode: handlers.ErrDatabase.Code})
			err = tx.Error
		}
		return
	}

	err = s.DB.Model(&models.UserList{}).Save(&models.UserList{
		UserID:   s.User.ID,
		ListType: byte(args.ListType),
		EntryID:  uint(args.UserID),
	}).Error

	if err != nil {
		out = append(out, ResponseAddUserToListError{ErrorCode: handlers.ErrDatabase.Code})
		return
	}

	res := ResponseAddUserToList{
		ErrorCode: 0,
		UserID:    args.UserID,
		ListType:  args.ListType,
	}
	copy(res.DisplayName[:], user.DisplayName)
	out = append(out, res)

	return out, err
}

// --- Packets ---

type ArgsAddToList struct {
	ListType types.UserListType
	UserID   types.UserID
}

type ResponseAddUserToListError types.ResponseErrorCode

func (r ResponseAddUserToListError) Type() types.PacketType { return types.ServerAddUserToList }

type ResponseAddUserToList struct {
	ErrorCode   int32
	UserID      types.UserID
	ListType    types.UserListType
	DisplayName [16]byte // Might be something else
}

func (r ResponseAddUserToList) Type() types.PacketType { return types.ServerAddUserToList }
