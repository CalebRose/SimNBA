package structs

import (
	"gorm.io/gorm"
)

// PlayerRecruitProfile - The points allocated to one player
type PlayerRecruitProfile struct {
	gorm.Model
	SeasonID              uint
	RecruitID             uint
	ProfileID             uint
	TotalPoints           float64
	AdjustedPoints        float64
	CurrentWeeksPoints    int
	PreviouslySpentPoints int
	SpendingCount         int
	Scholarship           bool
	ScholarshipRevoked    bool
	TeamAbbreviation      string
	InterestLevel         string
	RecruitModifier       int
	IsSigned              bool
	IsLocked              bool
	RemovedFromBoard      bool
	Recruit               Recruit `gorm:"foreignKey:RecruitID"`
	// RecruitPoints          []RecruitPointAllocation `gorm:"foreignKey:RecruitProfileID"`
}

func (r *PlayerRecruitProfile) AllocatePoints(points int) {
	r.CurrentWeeksPoints = points
}

func (r *PlayerRecruitProfile) SignPlayer() {
	if r.Scholarship {
		r.IsSigned = true
	}
}

func (r *PlayerRecruitProfile) LockPlayer() {
	r.IsLocked = true
}

func (r *PlayerRecruitProfile) AllocateTotalPoints(points float64) {
	r.TotalPoints += points
	r.PreviouslySpentPoints = r.CurrentWeeksPoints
	r.CurrentWeeksPoints = 0
}

func (r *PlayerRecruitProfile) AllocateScholarship() {
	r.Scholarship = true
}

func (r *PlayerRecruitProfile) RevokeScholarship() {
	r.Scholarship = false
}

func (r *PlayerRecruitProfile) RemoveRecruitFromBoard() {
	r.RemovedFromBoard = true
}

func (r *PlayerRecruitProfile) ReplaceRecruitToBoard() {
	r.RemovedFromBoard = false
}

// Sorting Funcs
type ByPoints []PlayerRecruitProfile

func (rp ByPoints) Len() int      { return len(rp) }
func (rp ByPoints) Swap(i, j int) { rp[i], rp[j] = rp[j], rp[i] }
func (rp ByPoints) Less(i, j int) bool {
	return rp[i].TotalPoints > rp[j].TotalPoints
}
