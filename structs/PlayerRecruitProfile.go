package structs

import (
	"gorm.io/gorm"
)

// RecruitPlayerProfile - The points allocated to one player
type RecruitPlayerProfile struct {
	gorm.Model
	SeasonID              uint
	RecruitID             uint
	ProfileID             uint
	TotalPoints           float64
	AdjustedPoints        float64
	CurrentWeeksPoints    uint8
	PreviouslySpentPoints uint8
	SpendingCount         uint8
	Scholarship           bool
	ScholarshipRevoked    bool
	TeamAbbreviation      string
	InterestLevel         string
	RecruitModifier       int
	IsSigned              bool
	IsLocked              bool
	HasStateBonus         bool
	HasRegionBonus        bool
	RemovedFromBoard      bool
	TeamReachedMax        bool
	Modifier              float32
	Agility               bool
	InsideShooting        bool
	MidrangeShooting      bool
	ThreePointShooting    bool
	FreeThrow             bool
	Ballwork              bool
	Stealing              bool
	Rebounding            bool
	Blocking              bool
	InteriorDefense       bool
	PerimeterDefense      bool
	BasketballIQ          bool
	Potential             bool
	// RecruitPoints          []RecruitPointAllocation `gorm:"foreignKey:RecruitProfileID"`
}

func (rp *RecruitPlayerProfile) AllocatePoints(points int) {
	rp.CurrentWeeksPoints = uint8(points)
}

func (rp *RecruitPlayerProfile) SignPlayer() {
	if rp.Scholarship {
		rp.IsSigned = true
	}
}

func (rp *RecruitPlayerProfile) LockPlayer() {
	rp.IsLocked = true
}

func (rp *RecruitPlayerProfile) AllocateTotalPoints(points float64) {
	// If user spends points on a recruit
	if points > 0 {
		rp.TotalPoints += points
		if rp.SpendingCount < 5 && points >= 1 {
			rp.SpendingCount++
			// In the event that someone tries to exploit the consistency system with a value between 0.00001 and 0.99999
		} else if points > 0 && points < 1 {
			rp.SpendingCount = 0
		}
	} else {
		rp.TotalPoints = 0

		rp.SpendingCount = 0
	}
	rp.PreviouslySpentPoints = rp.CurrentWeeksPoints
	rp.CurrentWeeksPoints = 0
}

func (rp *RecruitPlayerProfile) ResetTotalPoints() {
	rp.TotalPoints = 0
}

func (rp *RecruitPlayerProfile) ToggleTotalMax() {
	rp.TeamReachedMax = true
}

func (rp *RecruitPlayerProfile) ToggleScholarship(reward bool, revoke bool) {
	if rp.Scholarship {
		rp.RevokeScholarship()
		return
	}
	if !rp.ScholarshipRevoked {
		rp.Scholarship = true
	}
}

func (rp *RecruitPlayerProfile) RevokeScholarship() {
	rp.Scholarship = false
	rp.ScholarshipRevoked = true
}

func (rp *RecruitPlayerProfile) RemoveRecruitFromBoard() {
	rp.RemovedFromBoard = true
	rp.CurrentWeeksPoints = 0
}

func (rp *RecruitPlayerProfile) ReplaceRecruitToBoard() {
	rp.RemovedFromBoard = false
}

func (rp *RecruitPlayerProfile) ApplyScoutingAttribute(attr string) {
	if attr == "Agility" {
		rp.Agility = true
	}
	if attr == "InsideShooting" {
		rp.InsideShooting = true
	}
	if attr == "MidRangeShooting" {
		rp.MidrangeShooting = true
	}
	if attr == "ThreePointShooting" {
		rp.ThreePointShooting = true
	}
	if attr == "FreeThrow" {
		rp.FreeThrow = true
	}
	if attr == "Ballwork" {
		rp.Ballwork = true
	}
	if attr == "Stealing" {
		rp.Stealing = true
	}
	if attr == "Rebounding" {
		rp.Rebounding = true
	}
	if attr == "Blocking" {
		rp.Blocking = true
	}
	if attr == "BasketballIQ" {
		rp.BasketballIQ = true
	}
	if attr == "InteriorDefense" {
		rp.InteriorDefense = true
	}
	if attr == "PerimeterDefense" {
		rp.PerimeterDefense = true
	}
	if attr == "Potential" {
		rp.Potential = true
	}
}

// Sorting Funcs
type ByPoints []RecruitPlayerProfile

func (rp ByPoints) Len() int      { return len(rp) }
func (rp ByPoints) Swap(i, j int) { rp[i], rp[j] = rp[j], rp[i] }
func (rp ByPoints) Less(i, j int) bool {
	return rp[i].TotalPoints > rp[j].TotalPoints
}

type ScoutAttributeDTO struct {
	ProfileID uint
	RecruitID uint
	Attribute string
}

type RecruitingOdds struct {
	Odds          int
	IsCloseToHome bool
	IsRegional    bool
}
