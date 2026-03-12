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
	Potential             bool
	// RecruitPoints          []RecruitPointAllocation `gorm:"foreignKey:RecruitProfileID"`
}

func (r *RecruitPlayerProfile) AllocatePoints(points int) {
	r.CurrentWeeksPoints = uint8(points)
}

func (r *RecruitPlayerProfile) SignPlayer() {
	if r.Scholarship {
		r.IsSigned = true
	}
}

func (r *RecruitPlayerProfile) LockPlayer() {
	r.IsLocked = true
}

func (r *RecruitPlayerProfile) AllocateTotalPoints(points float64) {
	r.TotalPoints += points
	r.PreviouslySpentPoints = r.CurrentWeeksPoints
	r.CurrentWeeksPoints = 0
}

func (r *RecruitPlayerProfile) ResetTotalPoints() {
	r.TotalPoints = 0
}

func (r *RecruitPlayerProfile) ToggleTotalMax() {
	r.TeamReachedMax = true
}

func (r *RecruitPlayerProfile) ToggleScholarship(reward bool, revoke bool) {
	if r.Scholarship {
		r.RevokeScholarship()
		return
	}
	if !r.ScholarshipRevoked {
		r.Scholarship = true
	}
}

func (r *RecruitPlayerProfile) RevokeScholarship() {
	r.Scholarship = false
	r.ScholarshipRevoked = true
}

func (r *RecruitPlayerProfile) RemoveRecruitFromBoard() {
	r.RemovedFromBoard = true
	r.CurrentWeeksPoints = 0
}

func (r *RecruitPlayerProfile) ReplaceRecruitToBoard() {
	r.RemovedFromBoard = false
}

func (rp *RecruitPlayerProfile) ApplyScoutingAttribute(attr string) {
	if attr == "Agility" {
		rp.Agility = true
	}
	if attr == "Inside Shooting" {
		rp.InsideShooting = true
	}
	if attr == "Midrange Shooting" {
		rp.MidrangeShooting = true
	}
	if attr == "Three Point Shooting" {
		rp.ThreePointShooting = true
	}
	if attr == "Free Throw" {
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
	if attr == "Interior Defense" {
		rp.InteriorDefense = true
	}
	if attr == "Perimeter Defense" {
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
