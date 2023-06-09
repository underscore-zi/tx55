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

var ErrInvalidCredentials = handlers.ErrInvalidArguments.Code
var ErrNotFound = handlers.ErrNotFound.Code
var ErrDatabaseError = handlers.ErrDatabase.Code
var ErrBanned = handlers.ErrBanned.Code

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
	if tx := sess.DB.First(&row, "username LIKE ?", types.BytesToString(args.Username[:])); tx.Error != nil {
		if tx.Error != gorm.ErrRecordNotFound {
			sess.LogEntry().WithError(tx.Error).Error("Failed to query user")
			return []types.Response{ResponseLoginError{ErrorCode: ErrDatabaseError}}, nil
		}
	}

	if row.ID == 0 || !row.CheckPassword(args.Password[:]) {
		return []types.Response{ResponseLoginError{ErrorCode: ErrInvalidCredentials}}, nil
	}

	if err := sess.DB.First(&models.Ban{}, "user_id = ? and (type = ? or type = ?) and expires_at > NOW()", row.ID, models.UserBan, models.IPBan).Error; err == nil {
		// We do the ban check here, but it does technically allow a user who is already connected to a lobby
		// to stay connected to that lobby. I think it's a fair trade-off to avoid having to do an extra query
		// on every connection since the game connects/reconnects often
		return []types.Response{ResponseLoginError{ErrorCode: ErrBanned}}, nil
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
