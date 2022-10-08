package structs

import (
	"math/rand"

	"github.com/jinzhu/gorm"
)

type CollegePlayer struct {
	gorm.Model
	BasePlayer
	PlayerID      uint
	TeamID        uint
	TeamAbbr      string
	IsRedshirt    bool
	IsRedshirting bool
	HasGraduated  bool
	HasProgressed bool
	Stats         []CollegePlayerStats `gorm:"foreignKey:CollegePlayerID"`
}

func (c *CollegePlayer) UpdateMinutes(newMinutes int) {
	c.Minutes = newMinutes
}

func (c *CollegePlayer) SetID(id uint) {
	c.ID = id
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
	cp.HasProgressed = true
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

func (cp *CollegePlayer) MapFromRecruit(r Recruit, t Team) {
	cp.ID = r.ID
	cp.TeamID = t.ID
	cp.TeamAbbr = t.Abbr
	cp.PlayerID = r.PlayerID
	cp.State = r.State
	cp.Country = r.Country
	cp.Year = r.Age - 18
	cp.IsRedshirt = false
	cp.IsRedshirting = false
	cp.HasGraduated = false
	cp.Age = r.Age + 1
	cp.FirstName = r.FirstName
	cp.LastName = r.LastName
	cp.Position = r.Position
	cp.Height = r.Height
	cp.Stars = r.Stars
	cp.Overall = r.Overall
	cp.Shooting2 = r.Shooting2
	cp.Shooting3 = r.Shooting3
	cp.Finishing = r.Finishing
	cp.Ballwork = r.Ballwork
	cp.Rebounding = r.Rebounding
	cp.Defense = r.Defense
	cp.Stamina = r.Stamina
	cp.Potential = r.Potential
	cp.ProPotentialGrade = r.ProPotentialGrade
	cp.PotentialGrade = r.PotentialGrade
	cp.FreeAgency = r.FreeAgency
	cp.Personality = r.Personality
	cp.RecruitingBias = r.RecruitingBias
	cp.WorkEthic = r.WorkEthic
	cp.AcademicBias = r.AcademicBias
}

func (cp *CollegePlayer) GraduatePlayer() {
	cp.HasGraduated = true
}

func (p *CollegePlayer) SetRedshirtStatus() {
	if p.IsRedshirting {
		p.IsRedshirting = false
		p.IsRedshirt = true
	}
}

func (p *CollegePlayer) SetProgressionStatus() {
	p.HasProgressed = true
}

func (p *CollegePlayer) FixAge() {
	p.Year++
	p.Age = p.Year + 18
	if p.IsRedshirt {
		p.Age++
	}
}
