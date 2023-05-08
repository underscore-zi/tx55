package auth

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"reflect"
	"tx55/pkg/metalgearonline1/handlers"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/metalgearonline1/session"
	"tx55/pkg/metalgearonline1/types"
	"tx55/pkg/packet"
)

func init() {
	handlers.Register(LoginHandler{})
}

var ErrInvalidCredentials int32 = handlers.ErrInvalidArguments.Code
var ErrNotFound int32 = handlers.ErrNotFound.Code
var ErrDatabaseError int32 = handlers.ErrDatabase.Code

type LoginHandler struct{}

func (h LoginHandler) Type() types.PacketType {
	return types.ClientLogin
}

func (h LoginHandler) ArgumentTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf(ArgsLoginCredentials{}),
		reflect.TypeOf(ArgsLoginSession{}),
	}
}

func (h LoginHandler) Handle(_ *session.Session, _ *packet.Packet) ([]types.Response, error) {
	return nil, handlers.ErrUnexpectedArgument
}

func (h LoginHandler) HandleWithCredentials(sess *session.Session, args *ArgsLoginCredentials) ([]types.Response, error) {
	var row models.User
	sess.DB.First(&row, "username = ?", types.BytesToString(args.Username[:]))

	if row.ID == 0 || !row.CheckPassword(args.Password[:]) {
		return []types.Response{ResponseLoginError{ErrorCode: ErrInvalidCredentials}}, nil
	}

	// Only want to update the previous with a login with credentials
	// so if they maybe got TSU rank, it'll last until they disconnect entirely
	sess.DB.Model(row).Updates(map[string]interface{}{
		"previous_updated_at": gorm.Expr("updated_at"),
	})

	sess.Login(&row)

	// Valid login attempt, create a new session ID
	newSession := models.Session{
		UserID: row.ID,
	}

	tx := sess.DB.Model(&newSession).Create(&newSession)
	if tx.Error != nil {
		return []types.Response{ResponseLoginError{ErrorCode: ErrDatabaseError}}, nil
	}

	bs, _ := newSession.ID.MarshalBinary()
	resp := ResponseLogin{ErrorCode: 0}
	resp.SessionID = [16]byte(bs[:16])

	return []types.Response{resp}, nil
}

// HandleWithSession is just a session restore, gets a previously obtained session id and returns the user info
func (h LoginHandler) HandleWithSession(s *session.Session, args *ArgsLoginSession) ([]types.Response, error) {
	var out []types.Response
	id, err := uuid.FromBytes(args.SessionID[:])
	if err != nil {
		return nil, err
	}

	var row models.Session
	tx := s.DB.Model(&models.Session{}).Joins("User").Preload("User.PlayerSettings").Where(id).First(&row)

	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		return []types.Response{ResponseLoginError{ErrorCode: ErrDatabaseError}}, nil
	}

	if tx.RowsAffected == 0 || row.User.ID == 0 {
		return []types.Response{ResponseLoginError{ErrorCode: ErrNotFound}}, nil
	}

	s.Login(&row.User)
	out = append(out, ResponseLogin{
		ErrorCode: 0,
		SessionID: args.SessionID,
	})

	return out, nil
}

// --- Packets ---

func (r ResponseLogin) Type() types.PacketType      { return types.ServerLogin }
func (r ResponseLoginError) Type() types.PacketType { return types.ServerLogin }

type ResponseLoginError types.ResponseErrorCode
type ResponseLogin struct {
	ErrorCode int32
	SessionID [16]byte
}

type ArgsLoginCredentials struct {
	Username [16]byte
	Password [16]byte
}

type ArgsLoginSession struct {
	UserID    uint32
	SessionID [16]byte
}
