package structs

import (
	"sort"
)

type CollegePlayerResponse struct {
	FirstName             string
	LastName              string
	Archetype             string
	Position              string
	Age                   int
	Year                  int
	State                 string
	Country               string
	Stars                 int
	Height                string
	PotentialGrade        string
	Shooting2Grade        string
	Shooting3Grade        string
	FreeThrowGrade        string
	FinishingGrade        string
	BallworkGrade         string
	ReboundingGrade       string
	InteriorDefenseGrade  string
	PerimeterDefenseGrade string
	Stamina               int
	PlaytimeExpectations  int
	InsideProportion      float64
	MidRangeProportion    float64
	ThreePointProportion  float64
	PositionOne           string
	PositionTwo           string
	PositionThree         string
	P1Minutes             int
	P2Minutes             int
	P3Minutes             int
	Minutes               int
	Potential             int
	OverallGrade          string
	Personality           string
	RecruitingBias        string
	RecruitingBiasValue   string
	WorkEthic             string
	AcademicBias          string
	PlayerID              uint
	TeamID                uint
	TeamAbbr              string
	IsRedshirting         bool
	IsRedshirt            bool
	PreviousTeamID        uint
	PreviousTeam          string
	TransferStatus        int    // 1 == Intends, 2 == Is Transferring
	TransferLikeliness    string // Low, Medium, High
	LegacyID              uint   // Either a legacy school or a legacy coach
	SeasonStats           CollegePlayerSeasonStats
	Stats                 CollegePlayerStats
}

type NBAPlayerResponse struct {
	FirstName             string
	LastName              string
	Position              string
	Age                   int
	Year                  int
	State                 string
	Country               string
	Stars                 int
	Height                string
	PotentialGrade        string
	Shooting2Grade        string
	Shooting3Grade        string
	FreeThrowGrade        string
	FinishingGrade        string
	BallworkGrade         string
	ReboundingGrade       string
	InteriorDefenseGrade  string
	PerimeterDefenseGrade string
	Stamina               int
	PlaytimeExpectations  int
	Minutes               int
	Potential             int
	OverallGrade          string
	Personality           string
	RecruitingBias        string
	WorkEthic             string
	AcademicBias          string
	PlayerID              uint
	TeamID                uint
	TeamAbbr              string
	IsRedshirting         bool
	IsRedshirt            bool
	SeasonStats           NBAPlayerSeasonStats
	Stats                 NBAPlayerStats
}

type TransferPlayerResponse struct {
	FirstName            string
	LastName             string
	Archetype            string
	Position             string
	Age                  int
	Year                 int
	State                string
	Country              string
	Stars                int
	Height               string
	PotentialGrade       string
	Overall              string
	Shooting2            string
	Shooting3            string
	FreeThrow            string
	Finishing            string
	Ballwork             string
	Rebounding           string
	InteriorDefense      string
	PerimeterDefense     string
	Stamina              int
	PlaytimeExpectations int
	Minutes              int
	OverallGrade         string
	Personality          string
	RecruitingBias       string
	RecruitingBiasValue  string
	WorkEthic            string
	AcademicBias         string
	PlayerID             uint
	TeamID               uint
	TeamAbbr             string
	IsRedshirting        bool
	IsRedshirt           bool
	PreviousTeamID       uint
	PreviousTeam         string
	TransferStatus       int    // 1 == Intends, 2 == Is Transferring
	TransferLikeliness   string // Low, Medium, High
	LegacyID             uint   // Either a legacy school or a legacy coach
	SeasonStats          CollegePlayerSeasonStats
	Stats                CollegePlayerStats
	LeadingTeams         []LeadingTeams
}

func (c *TransferPlayerResponse) Map(r CollegePlayer, ovr string) {
	c.PlayerID = r.PlayerID
	c.TeamID = r.TeamID
	c.FirstName = r.FirstName
	c.LastName = r.LastName
	c.Position = r.Position
	c.Height = r.Height
	c.Stars = r.Stars
	c.Shooting2 = attributeMapper(r.Shooting2)
	c.Shooting3 = attributeMapper(r.Shooting3)
	c.Finishing = attributeMapper(r.Finishing)
	c.FreeThrow = attributeMapper(r.FreeThrow)
	c.Ballwork = attributeMapper(r.Ballwork)
	c.Rebounding = attributeMapper(r.Rebounding)
	c.InteriorDefense = attributeMapper(r.InteriorDefense)
	c.PerimeterDefense = attributeMapper(r.PerimeterDefense)
	c.OverallGrade = ovr
	c.PotentialGrade = r.PotentialGrade
	c.Personality = r.Personality
	c.RecruitingBias = r.RecruitingBias
	c.AcademicBias = r.AcademicBias
	c.WorkEthic = r.WorkEthic
	c.State = r.State
	c.Country = r.Country
	c.TeamAbbr = r.TeamAbbr
	c.PreviousTeam = r.PreviousTeam
	c.PreviousTeamID = r.PreviousTeamID

	var totalPoints float64 = 0
	var runningThreshold float64 = 0

	sortedProfiles := r.Profiles

	sort.Slice(sortedProfiles, func(i, j int) bool {
		return sortedProfiles[i].TotalPoints > sortedProfiles[j].TotalPoints
	})
	for _, recruitProfile := range sortedProfiles {
		if recruitProfile.RemovedFromBoard {
			continue
		}
		if runningThreshold == 0 {
			runningThreshold = float64(recruitProfile.TotalPoints) * 0.66
		}

		if float64(recruitProfile.TotalPoints) >= runningThreshold {
			totalPoints += float64(recruitProfile.TotalPoints)
		}

	}

	for i := 0; i < len(sortedProfiles); i++ {
		if sortedProfiles[i].RemovedFromBoard {
			continue
		}
		var odds float64 = 0

		if float64(sortedProfiles[i].TotalPoints) >= runningThreshold && runningThreshold > 0 {
			odds = float64(sortedProfiles[i].TotalPoints) / totalPoints
		}
		leadingTeam := LeadingTeams{
			TeamAbbr: r.Profiles[i].TeamAbbreviation,
			Odds:     odds,
		}
		c.LeadingTeams = append(c.LeadingTeams, leadingTeam)
	}
	sort.Sort(ByLeadingPoints(c.LeadingTeams))
}
