package structs

import "github.com/jinzhu/gorm"

type GlobalPlayer struct {
	gorm.Model
	RecruitID       uint
	CollegePlayerID uint
	NBAPlayerID     uint
}

func (gp *GlobalPlayer) SetID(id uint) {
	gp.ID = id
}
