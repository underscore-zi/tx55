package models

import (
	"crypto/md5"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
	"tx55/pkg/metalgearonline1/types"
)

func init() {
	All = append(All, &User{})
}

//goland:noinspection GoUnusedConst
const SALT = "\x84\xbd\xb8\xcf\xad\x46\xdd\x6e\x42\x4a\xe4\xd8\xd2\x6a\x12\xf3"

type User struct {
	gorm.Model
	PreviousUpdatedAt time.Time
	Username          []byte `gorm:"uniqueIndex,size:16"`
	DisplayName       []byte `gorm:"uniqueIndex,size:16"`
	Password          string `gorm:"type:varchar(128)"`
	HasEmblem         bool
	EmblemText        []byte `gorm:"size:16"`
	OverallRank       uint
	WeeklyRank        uint
	VsRating          uint
	VsRatingRank      uint
	Sessions          []Session
	PlayerSettings    PlayerSettings
	Connections       []Connection
	FBList            []UserList
}

// HashPassword will hash the password in the right format for MGO1 and then bcrypt it
func (u *User) HashPassword(password []byte) ([]byte, error) {
	sum := u.Md5Password(password)
	return bcrypt.GenerateFromPassword(sum, bcrypt.DefaultCost)
}

func (u *User) Md5Password(password []byte) []byte {
	hash := md5.New()
	hash.Write(u.Username)
	hash.Write(types.NONCE[:])
	hash.Write(password)
	return hash.Sum(nil)
}

func (u *User) CheckRawPassword(password []byte) bool {
	return u.CheckPassword(u.Md5Password(password))
}

func (u *User) CheckPassword(password []byte) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), password); err != nil {
		return false
	}
	return true
}

func (u *User) hashIfNecessary() error {
	if u.Password == "" {
		return nil
	}
	if len(u.Password) == 60 {
		// Already hashed (in theory)
		return nil
	}

	hash, err := u.HashPassword([]byte(u.Password))
	if err != nil {
		return err
	}

	u.Password = string(hash)
	return nil
}

func (u *User) BeforeCreate(_ *gorm.DB) error {
	if len(u.Password) == 0 {
		return errors.New("missing password")
	}
	return u.hashIfNecessary()
}

func (u *User) BeforeUpdate(_ *gorm.DB) error {
	return u.hashIfNecessary()
}

func (u *User) PlayerOverview() *types.PlayerOverview {
	o := types.PlayerOverview{
		UserID:          types.UserID(u.ID),
		VsRating:        uint32(u.VsRating),
		LastLogin:       uint32(u.UpdatedAt.Unix()),
		BeforeLastLogin: uint32(u.PreviousUpdatedAt.Unix()),
	}

	o.FillDisplayName(u.DisplayName)
	o.FillEmblem(u.HasEmblem, u.EmblemText)
	return &o
}

func (u *User) SharedAccounts(db *gorm.DB) []uint {
	query := "SELECT DISTINCT u2.id AS common_user_id\nFROM users u1\nJOIN connections c1 ON u1.id = c1.user_id\nJOIN connections c2 ON c1.remote_addr = c2.remote_addr AND c1.user_id <> c2.user_id\nJOIN users u2 ON c2.user_id = u2.id\nWHERE u1.id = ?;"
	var out []uint
	db.Raw(query, u.ID).Scan(&out)
	return out
}
