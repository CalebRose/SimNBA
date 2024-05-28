package structs

import "github.com/jinzhu/gorm"

type NBADraftee struct {
	gorm.Model
	BasePlayer
	PlayerID              uint
	CollegeID             uint
	College               string
	DraftPickID           uint
	DraftPick             string
	DraftedTeamID         uint
	DraftedTeamAbbr       string
	PrimeAge              int
	StandingReach         string
	VerticalLeap          float64
	LaneAgility           float64
	MaxVerticalLeap       float64
	ThreeQuarterSprint    float64
	ShuttleRun            float64
	WingSpan              string
	Shooting2Grade        string
	Shooting3Grade        string
	FreeThrowGrade        string
	FinishingGrade        string
	BallworkGrade         string
	ReboundingGrade       string
	InteriorDefenseGrade  string
	PerimeterDefenseGrade string
	OverallGrade          string
	Prediction            int
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
	n.Archetype = cp.Archetype
	n.Height = cp.Height
	n.Age = cp.Age
	n.Stars = cp.Stars
	n.Overall = cp.Overall
	n.Shooting2 = cp.Shooting2
	n.Shooting3 = cp.Shooting3
	n.FreeThrow = cp.FreeThrow
	n.Finishing = cp.Finishing
	n.Ballwork = cp.Ballwork
	n.Rebounding = cp.Rebounding
	n.InteriorDefense = cp.InteriorDefense
	n.PerimeterDefense = cp.PerimeterDefense
	n.Stamina = cp.Stamina
	n.ProPotentialGrade = cp.ProPotentialGrade
	n.Potential = cp.Potential
	n.PotentialGrade = cp.PotentialGrade
	n.FreeAgency = cp.FreeAgency
	n.Personality = cp.Personality
	n.RecruitingBias = cp.RecruitingBias
	n.WorkEthic = cp.WorkEthic
	n.AcademicBias = cp.AcademicBias
	n.SpecBallwork = cp.SpecBallwork
	n.SpecFinishing = cp.SpecFinishing
	n.SpecFreeThrow = cp.SpecFreeThrow
	n.SpecInteriorDefense = cp.SpecInteriorDefense
	n.SpecPerimeterDefense = cp.SpecPerimeterDefense
	n.SpecRebounding = cp.SpecRebounding
	n.SpecShooting2 = cp.SpecShooting2
	n.SpecShooting3 = cp.SpecShooting3
}

func (n *NBADraftee) MapInternational(cp NBAPlayer) {
	n.ID = cp.ID
	n.PlayerID = cp.PlayerID
	n.College = cp.TeamAbbr
	n.State = cp.State
	n.Country = cp.Country
	n.FirstName = cp.FirstName
	n.LastName = cp.LastName
	n.Position = cp.Position
	n.Archetype = cp.Archetype
	n.Height = cp.Height
	n.Age = cp.Age
	n.Stars = cp.Stars
	n.Overall = cp.Overall
	n.Shooting2 = cp.Shooting2
	n.Shooting3 = cp.Shooting3
	n.FreeThrow = cp.FreeThrow
	n.Finishing = cp.Finishing
	n.Ballwork = cp.Ballwork
	n.Rebounding = cp.Rebounding
	n.InteriorDefense = cp.InteriorDefense
	n.PerimeterDefense = cp.PerimeterDefense
	n.Stamina = cp.Stamina
	n.ProPotentialGrade = cp.ProPotentialGrade
	n.Potential = cp.Potential
	n.PotentialGrade = cp.PotentialGrade
	n.FreeAgency = cp.FreeAgency
	n.Personality = cp.Personality
	n.RecruitingBias = cp.RecruitingBias
	n.WorkEthic = cp.WorkEthic
	n.AcademicBias = cp.AcademicBias
	n.SpecBallwork = cp.SpecBallwork
	n.SpecFinishing = cp.SpecFinishing
	n.SpecFreeThrow = cp.SpecFreeThrow
	n.SpecInteriorDefense = cp.SpecInteriorDefense
	n.SpecPerimeterDefense = cp.SpecPerimeterDefense
	n.SpecRebounding = cp.SpecRebounding
	n.SpecShooting2 = cp.SpecShooting2
	n.SpecShooting3 = cp.SpecShooting3
}

func (n *NBADraftee) AssignPrimeAge(age int) {
	n.PrimeAge = age
}

func (n *NBADraftee) AssignProPotentialGrade(potential int) {
	n.ProPotentialGrade = potential
}

func (n *NBADraftee) ApplyGrades(s2, s3, ft, fn, bw, rb, id, pd, ov string) {
	n.Shooting2Grade = s2
	n.Shooting3Grade = s3
	n.FreeThrowGrade = ft
	n.FinishingGrade = fn
	n.BallworkGrade = bw
	n.ReboundingGrade = rb
	n.InteriorDefenseGrade = id
	n.PerimeterDefenseGrade = pd
	n.OverallGrade = ov
}

func (n *NBADraftee) PredictRound(round int) {
	n.Prediction = round
}

func (n *NBADraftee) AssignDraftedTeam(DraftPick string, PickID, TeamID uint, Abbr string) {
	n.DraftPick = DraftPick
	n.DraftPickID = PickID
	n.DraftedTeamID = TeamID
	n.DraftedTeamAbbr = Abbr
}
