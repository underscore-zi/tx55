package admin

import "golang.org/x/crypto/bcrypt"

var AllModels = []interface{}{&User{}, &Role{}}

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"uniqueIndex"`
	Password string `json:"-"`
	RoleID   uint   `json:"role_id"`
	Role     Role   `json:"role" gorm:"foreignKey:RoleID"`
}

func (u *User) CheckPassword(password []byte) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), password); err != nil {
		return false
	}
	return true
}

type Role struct {
	ID            uint   `json:"id" gorm:"primaryKey"`
	Name          string `json:"name" gorm:"uniqueIndex"`
	AllPrivileges bool   `json:"all_privileges"`
}
