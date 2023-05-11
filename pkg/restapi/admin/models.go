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

type Privilege string

const (
	PrivNone           Privilege = ""
	PrivAll            Privilege = "all_privileges"
	PrivUpdateProfiles Privilege = "update_profiles"
)

func (u *User) HasPrivilege(p Privilege) bool {
	if u.Role.AllPrivileges {
		return true
	}

	switch p {
	case PrivNone:
		return true
	case PrivAll:
		return u.Role.AllPrivileges
	case PrivUpdateProfiles:
		return u.Role.UpdateProfiles
	}

	return false
}
