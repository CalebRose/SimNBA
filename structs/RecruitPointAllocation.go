package structs

import "github.com/jinzhu/gorm"

type RecruitPointAllocation struct {
	gorm.Model
	RecruitID           uint
	TeamProfileID       uint
	RecruitProfileID    uint
	WeekID              uint
	Points              float64
	BonusAffectedPoints float64
	AffinityOneApplied  bool
	AffinityTwoApplied  bool
	CaughtCheating      bool
}

func (rpa *RecruitPointAllocation) UpdatePointsSpent(points float64, curr float64) {
	rpa.Points = points
	rpa.BonusAffectedPoints = curr
}

func (rpa *RecruitPointAllocation) ApplyAffinityOne() {
	rpa.AffinityOneApplied = true
}

func (rpa *RecruitPointAllocation) ApplyAffinityTwo() {
	rpa.AffinityTwoApplied = true
}

func (rpa *RecruitPointAllocation) ApplyCaughtCheating() {
	rpa.CaughtCheating = true
}
