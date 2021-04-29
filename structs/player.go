package structs

import (
	"github.com/jinzhu/gorm"
)

// Player - The NBA player for the sim
type Player struct {
	gorm.Model
	FirstName             string
	LastName              string
	Position              string
	Year                  int
	State                 string
	Country               string
	Stars                 int
	Height                string
	TeamID                int
	Shooting              int
	Finishing             int
	Ballwork              int
	Rebounding            int
	Defense               int
	PotentialGrade        int
	ProPotentialGrade     int
	Stamina               int
	PlaytimeExpectations  int
	MinutesA              int
	MinutesB              int
	MinutesC              int
	Overall               int
	UninterestedThreshold int
	LowInterestThreshold  int
	MedInterestThreshold  int
	HighInterestThreshold int
	ReadyToSignThreshold  int
	IsNBA                 bool
	IsRedshirt            bool
	IsRedshirting         bool
	Contracts             []Contract         `gorm:"foreignKey:PlayerID"`
	RecruitingPoints      []RecruitingPoints `gorm:"foreignKey:PlayerID"`
	Stats                 []PlayerStats      `gorm:"foreignKey:PlayerID"`
}

func (p *Player) SetRedshirtingStatus() {
	if p.IsRedshirt == false && p.IsRedshirting == false {
		p.IsRedshirting = true
	}
}

func (p *Player) SetRedshirtStatus() {
	if p.IsRedshirting == true {
		p.IsRedshirting = false
		p.IsRedshirt = true
	}
}
