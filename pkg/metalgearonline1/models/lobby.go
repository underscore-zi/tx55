package models

import "tx55/pkg/metalgearonline1/types"

func init() {
	All = append(All, &Lobby{})
}

type Lobby struct {
	ID      uint32
	Name    string `gorm:"type:varchar(16)"`
	Type    types.LobbyType
	IP      string `gorm:"type:varchar(15)"`
	Port    uint16
	Players uint16
}
