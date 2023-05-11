package admin

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var AllModels = []interface{}{&User{}, &Role{}}

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"uniqueIndex"`
	Password string `json:"-"`
	RoleID   uint   `json:"role_id"`
	Role     Role   `json:"role" gorm:"foreignKey:RoleID"`
}

type Role struct {
	ID             uint   `json:"id" gorm:"primaryKey"`
	Name           string `json:"name" gorm:"uniqueIndex"`
	AllPrivileges  bool   `json:"all_privileges"`
	UpdateProfiles bool   `json:"update_profiles"`
	ReadIPs        bool   `json:"read_ips"`
	ReadBans       bool   `json:"read_bans"`
	UpdateBans     bool   `json:"update_bans"`
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

	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashed)
	return nil
}

func (u *User) BeforeCreate(_ *gorm.DB) error {
	if u.Password == "" {
		return errors.New("missing password")
	}
	return u.hashIfNecessary()
}

func (u *User) BeforeUpdate(_ *gorm.DB) error {
	return u.hashIfNecessary()
}
