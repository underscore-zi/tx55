package models

func init() {
	All = append(All, &Friend{}, &Blocked{})
}

type Friend struct {
	ID       uint
	UserID   uint
	User     User
	FriendID uint
	Friend   User `gorm:"foreignKey:FriendID"`
}

type Blocked struct {
	ID        uint
	UserID    uint
	User      User
	BlockedID uint
	Blocked   User `gorm:"foreignKey:BlockedID"`
}
