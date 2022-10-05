package structs

import "github.com/jinzhu/gorm"

type GlobalPlayer struct {
	gorm.Model
	RecruitID       int
	CollegePlayerID int
	NBAPlayerID     int
}
