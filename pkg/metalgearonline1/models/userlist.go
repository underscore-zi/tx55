package models

func init() {
	All = append(All, &UserList{})
}

type UserList struct {
	ID       uint
	UserID   uint
	User     User
	EntryID  uint
	Entry    User `gorm:"foreignKey:EntryID"`
	ListType byte
}
