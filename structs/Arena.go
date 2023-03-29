package structs

import "gorm.io/gorm"

type Arena struct {
	gorm.Model
	ArenaName string
	City      string
	State     string
	Country   string
	Capacity  uint
	HomeTeam  string
}

func (a *Arena) AssignID(id uint) {
	a.ID = id
}
