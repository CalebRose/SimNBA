package structs

import (
	"github.com/jinzhu/gorm"
)

// Player - The NBA player for the sim
type Player struct {
	gorm.Model
	FirstName            string
	LastName             string
	Position             string
	Year                 int
	State                string
	Country              string
	Stars                int
	Height               string
	TeamID               int
	Shooting             int
	Finishing            int
	Ballwork             int
	Rebounding           int
	Defense              int
	PotentialGrade       int
	Stamina              int
	PlaytimeExpectations int
	Minutes              int
	Overall              int
	IsNBA                bool
	IsRedshirt           bool
	IsRedshirting        bool
	Contracts            []Contract         `gorm:"foreignKey:PlayerID"`
	RecruitingPoints     []RecruitingPoints `gorm:"foreignKey:PlayerID"`
}
