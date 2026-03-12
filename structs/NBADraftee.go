package structs

import "github.com/jinzhu/gorm"

type NBADraftee struct {
	gorm.Model
	BasePlayer
	PlayerID                uint
	CollegeID               uint
	College                 string
	DraftPickID             uint
	DraftPick               string
	DraftedTeamID           uint
	DraftedTeam             string
	PrimeAge                int
	StandingReach           string
	VerticalLeap            float64
	LaneAgility             float64
	MaxVerticalLeap         float64
	ThreeQuarterSprint      float64
	ShuttleRun              float64
	WingSpan                string
	MidrangeShootingGrade   string
	ThreePointShootingGrade string
	FreeThrowGrade          string
	InsideShootingGrade     string
	BallworkGrade           string
	AgilityGrade            string
	StealingGrade           string
	BlockingGrade           string
	ReboundingGrade         string
	InteriorDefenseGrade    string
	PerimeterDefenseGrade   string
	OverallGrade            string
	Prediction              int
	IsInternational         bool
}

func (n *NBADraftee) Map(cp CollegePlayer) {
	n.ID = cp.ID
	n.PlayerID = cp.PlayerID
	n.CollegeID = cp.TeamID
	n.College = cp.Team
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
	n.MidRangeShooting = cp.MidRangeShooting
	n.ThreePointShooting = cp.ThreePointShooting
	n.FreeThrow = cp.FreeThrow
	n.InsideShooting = cp.InsideShooting
	n.Ballwork = cp.Ballwork
	n.Agility = cp.Agility
	n.Stealing = cp.Stealing
	n.Blocking = cp.Blocking
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
	n.SpecInsideShooting = cp.SpecInsideShooting
	n.SpecFreeThrow = cp.SpecFreeThrow
	n.SpecInteriorDefense = cp.SpecInteriorDefense
	n.SpecPerimeterDefense = cp.SpecPerimeterDefense
	n.SpecRebounding = cp.SpecRebounding
	n.SpecMidRangeShooting = cp.SpecMidRangeShooting
	n.SpecThreePointShooting = cp.SpecThreePointShooting
	n.SpecAgility = cp.SpecAgility
	n.SpecStealing = cp.SpecStealing
	n.SpecBlocking = cp.SpecBlocking
}

func (n *NBADraftee) MapInternational(cp NBAPlayer) {
	n.ID = cp.ID
	n.PlayerID = cp.PlayerID
	n.CollegeID = cp.TeamID
	n.College = cp.Team
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
	n.MidRangeShooting = cp.MidRangeShooting
	n.ThreePointShooting = cp.ThreePointShooting
	n.FreeThrow = cp.FreeThrow
	n.InsideShooting = cp.InsideShooting
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
	n.SpecInsideShooting = cp.SpecInsideShooting
	n.SpecFreeThrow = cp.SpecFreeThrow
	n.SpecInteriorDefense = cp.SpecInteriorDefense
	n.SpecPerimeterDefense = cp.SpecPerimeterDefense
	n.SpecRebounding = cp.SpecRebounding
	n.SpecMidRangeShooting = cp.SpecMidRangeShooting
	n.SpecThreePointShooting = cp.SpecThreePointShooting
	n.IsInternational = true
}

func (n *NBADraftee) AssignPrimeAge(age int) {
	n.PrimeAge = age
}

func (n *NBADraftee) AssignProPotentialGrade(potential uint8) {
	n.ProPotentialGrade = potential
}

func (n *NBADraftee) ApplyGrades(s2, s3, ft, fn, bw, rb, id, pd, ov string) {
	n.MidrangeShootingGrade = s2
	n.ThreePointShootingGrade = s3
	n.FreeThrowGrade = ft
	n.InsideShootingGrade = fn
	n.BallworkGrade = bw
	n.ReboundingGrade = rb
	n.InteriorDefenseGrade = id
	n.PerimeterDefenseGrade = pd
	n.OverallGrade = ov
}

func (n *NBADraftee) PredictRound(round int) {
	n.Prediction = round
}

func (n *NBADraftee) AssignDraftedTeam(DraftPick string, PickID, TeamID uint, Team string) {
	n.DraftPick = DraftPick
	n.DraftPickID = PickID
	n.DraftedTeamID = TeamID
	n.DraftedTeam = Team
}
