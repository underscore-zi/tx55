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
	handlers.Register(AddUserToList{})
}

type AddUserToList struct{}

func (h AddUserToList) Type() types.PacketType {
	return types.ClientAddUserToList
}

func (h AddUserToList) ArgumentTypes() (out []reflect.Type) {
	//out = append(out, reflect.TypeOf(ArgsAddToList{}))
	return
}

func (h AddUserToList) Handle(_ *session.Session, _ *packet.Packet) (out []types.Response, err error) {
	// Not exactly sure how these packets work, the code does something pretty weird and when we don't complete
	// the request as expected the client stops doing any FL/BL actions
	// But when we send all this data, it removes them from the visible list immediately (ingame only)
	// No packet is sent for deletion, so really this just doesn't work. So I'm leaving the code for my own sake
	// and we'll keep the NOP response for now

	out = append(out, ResponseAddUserToList{})
	return
}

func (h AddUserToList) addFriend(s *session.Session, args *ArgsAddToList) (out ResponseAddUserToList, err error) {
	var list []models.Friend
	if tx := s.DB.Where("user_id = ?", s.User.ID).Find(&list); tx.Error != nil {
		out.ErrorCode = handlers.ErrDatabase.Code
		err = tx.Error
		return
	}

	if len(list) >= 16 {
		out.ErrorCode = handlers.ErrInvalidArguments.Code
		err = handlers.ErrInvalidArguments
		return
	}

	for _, entry := range list {
		if types.UserID(entry.FriendID) == args.UserID {
			// User is already in the list
			out.ErrorCode = 0
			out.UserID = args.UserID
			out.ListType = args.ListType
			copy(out.DisplayName[:], "some random name")
			err = nil
			return
		}
	}

	newEntry := models.Friend{
		UserID:   s.User.ID,
		FriendID: uint(args.UserID),
	}
	if tx := s.DB.Create(&newEntry); tx.Error != nil {
		out.ErrorCode = handlers.ErrDatabase.Code
		err = tx.Error
		return
	}

	out.ErrorCode = 0
	out.UserID = args.UserID
	out.ListType = args.ListType
	copy(out.DisplayName[:], "some random name")
	err = nil
	return
}

func (h AddUserToList) addBlocked(s *session.Session, args *ArgsAddToList) (out ResponseAddUserToList, err error) {
	var list []models.Blocked
	if tx := s.DB.Where("user_id = ?", s.User.ID).Find(&list); tx.Error != nil {
		out.ErrorCode = handlers.ErrDatabase.Code
		err = tx.Error
		return
	}

	if len(list) >= 16 {
		out.ErrorCode = handlers.ErrInvalidArguments.Code
		err = handlers.ErrInvalidArguments
		return
	}

	for _, entry := range list {
		if types.UserID(entry.BlockedID) == args.UserID {
			// User is already in the list
			out.ErrorCode = 0
			out.UserID = args.UserID
			out.ListType = args.ListType
			copy(out.DisplayName[:], "some random name")
			err = nil
			return
		}
	}

	newEntry := models.Blocked{
		UserID:    s.User.ID,
		BlockedID: uint(args.UserID),
	}
	if tx := s.DB.Create(&newEntry); tx.Error != nil {
		out.ErrorCode = handlers.ErrDatabase.Code
		err = tx.Error
		return
	}

	out.ErrorCode = 0
	out.UserID = args.UserID
	out.ListType = args.ListType
	copy(out.DisplayName[:], "some random name")
	err = nil
	return
}

func (h AddUserToList) HandleArgs(s *session.Session, args *ArgsAddToList) (out []types.Response, err error) {
	switch args.ListType {
	case types.UserListFriends:
		var res ResponseAddUserToList
		res, err = h.addFriend(s, args)
		out = append(out, res)
	case types.UserListBlocked:
		var res ResponseAddUserToList
		res, err = h.addBlocked(s, args)
		out = append(out, res)
	}
	return out, err
}

// --- Packets ---

type ArgsAddToList struct {
	ListType types.UserListType
	UserID   types.UserID
}

type ResponseAddUserToListEmpty types.ResponseEmpty

func (r ResponseAddUserToListEmpty) Type() types.PacketType { return types.ServerAddUserToList }

type ResponseAddUserToListError types.ResponseErrorCode

func (r ResponseAddUserToListError) Type() types.PacketType { return types.ServerAddUserToList }

type ResponseAddUserToList struct {
	ErrorCode   int32
	UserID      types.UserID
	ListType    types.UserListType
	DisplayName [16]byte // Might be something else
}

func (r ResponseAddUserToList) Type() types.PacketType { return types.ServerAddUserToList }
