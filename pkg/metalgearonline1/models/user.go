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

const SALT = "\x84\xbd\xb8\xcf\xad\x46\xdd\x6e\x42\x4a\xe4\xd8\xd2\x6a\x12\xf3"

type User struct {
	gorm.Model
	Username       []byte `gorm:"uniqueIndex,size:16"`
	DisplayName    []byte `gorm:"uniqueIndex,size:16"`
	Password       string `gorm:"type:varchar(128)"`
	HasEmblem      bool
	EmblemText     []byte `gorm:"size:16"`
	Sessions       []Session
	PlayerSettings PlayerSettings
	Connections    []Connection
	Friends        []Friend
	Blocked        []Blocked
}

// HashPassword will hash the password in the right format for MGO1 and then bcrypt it
func (u *User) HashPassword(password []byte) ([]byte, error) {
	hash := md5.New()
	hash.Write([]byte(u.Username))
	hash.Write(types.NONCE[:])
	hash.Write(password)
	sum := hash.Sum(nil)

	return bcrypt.GenerateFromPassword(sum, bcrypt.DefaultCost)
}

func (u *User) CheckPassword(password []byte) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return false
	}
	return true
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Password == "" {
		return errors.New("missing password")
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

func (u *User) PlayerOverview() *types.PlayerOverview {
	o := types.PlayerOverview{
		UserID:      types.UserID(u.ID),
		VSRating:    0,
		CurrentTime: uint32(time.Now().Unix()),
		LastLogin:   uint32(time.Now().Add(-time.Hour * 24 * 31).Unix()),
	}

	o.FillDisplayName(u.DisplayName)
	o.FillEmblem(u.HasEmblem, u.EmblemText)
	return &o
}
