package structs

import (
	"math/rand"

	"github.com/jinzhu/gorm"
)

type CollegePlayer struct {
	gorm.Model
	BasePlayer
	PlayerID      int
	TeamID        int
	TeamAbbr      string
	IsRedshirt    bool
	IsRedshirting bool
	HasGraduated  bool
	Stats         []CollegePlayerStats
}

func (c *CollegePlayer) UpdateMinutes(newMinutes int) {
	c.Minutes = newMinutes
}

func (c *CollegePlayer) SetID(id int) {
	c.ID = uint(id)
}

func (cp *CollegePlayer) Progress(attr CollegePlayerProgressions) {
	// cp.Age++
	// cp.Year++
	cp.Ballwork = attr.Ballwork
	cp.Shooting2 = attr.Shooting2
	cp.Shooting3 = attr.Shooting3
	cp.Finishing = attr.Finishing
	cp.Defense = attr.Defense
	cp.Rebounding = attr.Rebounding
	cp.Overall = attr.Overall
}

func (cp *CollegePlayer) GetPotentialGrade() {
	adjust := rand.Intn(20) - 10
	if adjust == 0 {
		test := rand.Intn(2000) - 1000

		if test > 0 {
			adjust += 1
		} else if test < 0 {
			adjust -= 1
		} else {
			adjust = 0
		}
	}
	potential := cp.Potential + adjust
	if potential > 80 {
		cp.PotentialGrade = "A+"
	} else if potential > 70 {
		cp.PotentialGrade = "A"
	} else if potential > 65 {
		cp.PotentialGrade = "A-"
	} else if potential > 60 {
		cp.PotentialGrade = "B+"
	} else if potential > 55 {
		cp.PotentialGrade = "B"
	} else if potential > 50 {
		cp.PotentialGrade = "B-"
	} else if potential > 40 {
		cp.PotentialGrade = "C+"
	} else if potential > 30 {
		cp.PotentialGrade = "C"
	} else if potential > 25 {
		cp.PotentialGrade = "C-"
	} else if potential > 20 {
		cp.PotentialGrade = "D+"
	} else if potential > 15 {
		cp.PotentialGrade = "D"
	} else if potential > 10 {
		cp.PotentialGrade = "D-"
	} else {
		cp.PotentialGrade = "F"
	}
}
