package structs

import (
	"sort"
)

type Croot struct {
	ID                 uint
	PlayerID           uint
	TeamID             uint
	College            string
	FirstName          string
	LastName           string
	Position           string
	Archetype          string
	Height             uint8
	Weight             uint8
	Stars              uint8
	MidRangeShooting   string
	ThreePointShooting string
	FreeThrow          string
	InsideShooting     string
	Ballwork           string
	Agility            string
	Stealing           string
	Blocking           string
	Rebounding         string
	InteriorDefense    string
	PerimeterDefense   string
	PotentialGrade     string
	Personality        string
	RecruitingBias     string
	AcademicBias       string
	WorkEthic          string
	State              string
	Country            string
	ESPNRank           float64
	RivalsRank         float64
	Rank247            float64
	IsSigned           bool
	OverallGrade       string
	TotalRank          float64
	SigningStatus      string
	IsCustomCroot      bool
	CreatedFor         string
	RelativeID         uint8
	Notes              string
	LeadingTeams       []LeadingTeams
}

type LeadingTeams struct {
	TeamName    string
	TeamAbbr    string
	TeamID      uint
	Odds        float64
	Scholarship bool
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
	c.MidRangeShooting = attributeMapper(r.MidRangeShooting, 1)
	c.ThreePointShooting = attributeMapper(r.ThreePointShooting, 1)
	c.InsideShooting = attributeMapper(r.InsideShooting, 1)
	c.FreeThrow = attributeMapper(r.FreeThrow, 1)
	c.Ballwork = attributeMapper(r.Ballwork, 1)
	c.Rebounding = attributeMapper(r.Rebounding, 1)
	c.InteriorDefense = attributeMapper(r.InteriorDefense, 1)
	c.PerimeterDefense = attributeMapper(r.PerimeterDefense, 1)
	c.Agility = attributeMapper(r.Agility, 1)
	c.Blocking = attributeMapper(r.Blocking, 1)
	c.Stealing = attributeMapper(r.Stealing, 1)
	c.PotentialGrade = r.PotentialGrade
	c.Personality = r.Personality
	c.RecruitingBias = r.RecruitingBias
	c.AcademicBias = r.AcademicBias
	c.WorkEthic = r.WorkEthic
	c.State = r.State
	c.Country = r.Country
	c.College = r.Team
	c.IsSigned = r.IsSigned
	c.SigningStatus = r.SigningStatus
	c.ESPNRank = r.ESPNRank
	c.RivalsRank = r.RivalsRank
	c.Rank247 = r.Rank247
	c.IsCustomCroot = r.IsCustomCroot
	c.CreatedFor = r.CreatedFor
	c.RelativeID = r.RelativeID
	c.Notes = r.Notes

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
		if recruitProfile.TeamReachedMax {
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
		if sortedProfiles[i].TeamReachedMax || sortedProfiles[i].RemovedFromBoard {
			continue
		}
		var odds float64 = 0

		if float64(sortedProfiles[i].TotalPoints) >= runningThreshold && runningThreshold > 0 {
			odds = float64(sortedProfiles[i].TotalPoints) / totalPoints
		}
		leadingTeam := LeadingTeams{
			TeamAbbr:    r.RecruitProfiles[i].TeamAbbreviation,
			TeamID:      r.RecruitProfiles[i].ProfileID,
			Odds:        odds,
			Scholarship: r.RecruitProfiles[i].Scholarship,
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
	return c[i].Stars > c[j].Stars && c[i].TotalRank > c[j].TotalRank
}

func attributeMapper(attr, year uint8) string {
	if year < 3 {
		if attr > 23 {
			return "A"
		}
		if attr > 18 {
			return "B"
		}
		if attr > 10 {
			return "C"
		}
		if attr > 5 {
			return "D"
		}
		return "F"
	}
	if attr > 29 {
		return "A+"
	}
	if attr > 26 {
		return "A"
	}
	if attr > 23 {
		return "A-"
	}
	if attr > 20 {
		return "B+"
	}
	if attr > 17 {
		return "B"
	}
	if attr > 14 {
		return "B-"
	}
	if attr > 11 {
		return "C+"
	}
	if attr > 8 {
		return "C"
	}
	if attr > 5 {
		return "C-"
	}
	if attr > 2 {
		return "D"
	}
	return "F"
}
