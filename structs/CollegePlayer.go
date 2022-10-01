package structs

import "github.com/jinzhu/gorm"

type CollegePlayer struct {
	gorm.Model
	BasePlayer
	PlayerID      int
	TeamID        int
	TeamAbbr      string
	IsRedshirt    bool
	IsRedshirting bool
	HasGraduated  bool
	// CollegePlayerStats
}
