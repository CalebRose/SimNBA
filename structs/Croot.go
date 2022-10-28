package structs

import (
	"sort"
)

type Croot struct {
	ID             uint
	PlayerID       uint
	TeamID         uint
	College        string
	FirstName      string
	LastName       string
	Position       string
	Archetype      string
	Height         string
	Stars          int
	Shooting2      int
	Shooting3      int
	Finishing      int
	Ballwork       int
	Rebounding     int
	Defense        int
	PotentialGrade string
	Personality    string
	RecruitingBias string
	AcademicBias   string
	WorkEthic      string
	State          string
	Country        string
	IsSigned       bool
	OverallGrade   string
	TotalRank      float64
	SigningStatus  string
	LeadingTeams   []LeadingTeams
}

type LeadingTeams struct {
	TeamName string
	TeamAbbr string
	Odds     float64
}

// Sorting Funcs
type ByLeadingPoints []LeadingTeams

func (rp ByLeadingPoints) Len() int      { return len(rp) }
func (rp ByLeadingPoints) Swap(i, j int) { rp[i], rp[j] = rp[j], rp[i] }
func (rp ByLeadingPoints) Less(i, j int) bool {
	return rp[i].Odds > rp[j].Odds
}

func (c *Croot) Map(r Recruit) {
	c.ID = r.ID
	c.PlayerID = r.PlayerID
	c.TeamID = r.TeamID
	c.FirstName = r.FirstName
	c.LastName = r.LastName
	c.Position = r.Position
	c.Height = r.Height
	c.Stars = r.Stars
	c.Shooting2 = r.Shooting2
	c.Shooting3 = r.Shooting3
	c.Finishing = r.Finishing
	c.Ballwork = r.Ballwork
	c.Rebounding = r.Rebounding
	c.Defense = r.Defense
	c.PotentialGrade = r.PotentialGrade
	c.Personality = r.Personality
	c.RecruitingBias = r.RecruitingBias
	c.AcademicBias = r.AcademicBias
	c.WorkEthic = r.WorkEthic
	c.State = r.State
	c.Country = r.Country
	c.College = r.TeamAbbr
	c.IsSigned = r.IsSigned
	c.SigningStatus = r.SigningStatus

	mod := r.TopRankModifier
	if mod == 0 {
		mod = 1
	}
	c.TotalRank = (r.RivalsRank + r.ESPNRank + r.Rank247) / r.TopRankModifier

	var totalPoints float64 = 0
	var runningThreshold float64 = 0

	sortedProfiles := r.RecruitProfiles

	sort.Sort(ByPoints(sortedProfiles))

	for _, recruitProfile := range sortedProfiles {
		if !recruitProfile.Scholarship && r.TeamAbbr == "" {
			continue
		}
		if runningThreshold == 0 {
			runningThreshold = float64(recruitProfile.TotalPoints) * 0.5
		}

		if float64(recruitProfile.TotalPoints) >= runningThreshold {
			totalPoints += float64(recruitProfile.TotalPoints)
		}

	}

	for i := 0; i < len(sortedProfiles); i++ {
		if !sortedProfiles[i].Scholarship && r.TeamAbbr == "" {
			continue
		}
		var odds float64 = 0

		if float64(sortedProfiles[i].TotalPoints) >= runningThreshold && runningThreshold > 0 {
			odds = float64(sortedProfiles[i].TotalPoints) / totalPoints
		}
		leadingTeam := LeadingTeams{
			TeamAbbr: r.RecruitProfiles[i].TeamAbbreviation,
			Odds:     odds,
		}
		c.LeadingTeams = append(c.LeadingTeams, leadingTeam)
	}
	sort.Sort(ByLeadingPoints(c.LeadingTeams))
}

func (c *Croot) SetOverallGrade(grade string) {
	c.OverallGrade = grade
}

type ByCrootRank []Croot

func (c ByCrootRank) Len() int      { return len(c) }
func (c ByCrootRank) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c ByCrootRank) Less(i, j int) bool {
	return c[i].TotalRank > c[j].TotalRank || c[i].Stars > c[j].Stars
}
