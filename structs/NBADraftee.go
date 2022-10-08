package structs

import "github.com/jinzhu/gorm"

type NBADraftee struct {
	gorm.Model
	BasePlayer
	PlayerID        uint
	CollegeID       uint
	College         string
	DraftPickID     uint
	DraftPick       string
	DraftedTeamID   uint
	DraftedTeamAbbr string
	PrimeAge        int
}

func (n *NBADraftee) Map(cp CollegePlayer) {
	n.ID = cp.ID
	n.PlayerID = cp.PlayerID
	n.College = cp.TeamAbbr
	n.State = cp.State
	n.Country = cp.Country
	n.FirstName = cp.FirstName
	n.LastName = cp.LastName
	n.Position = cp.Position
	n.Height = cp.Height
	n.Age = cp.Age
	n.Stars = cp.Stars
	n.Overall = cp.Overall
	n.Shooting2 = cp.Shooting2
	n.Shooting3 = cp.Shooting3
	n.Finishing = cp.Finishing
	n.Ballwork = cp.Ballwork
	n.Rebounding = cp.Rebounding
	n.Defense = cp.Defense
	n.Stamina = cp.Stamina
	n.Stamina = cp.Stamina
	n.ProPotentialGrade = cp.ProPotentialGrade
	n.Potential = cp.Potential
	n.PotentialGrade = cp.PotentialGrade
	n.FreeAgency = cp.FreeAgency
	n.Personality = cp.Personality
	n.RecruitingBias = cp.RecruitingBias
	n.WorkEthic = cp.WorkEthic
	n.AcademicBias = cp.AcademicBias
}

func (n *NBADraftee) AssignPrimeAge(age int) {
	n.PrimeAge = age
}
