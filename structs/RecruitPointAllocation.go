package structs

import "github.com/jinzhu/gorm"

type RecruitPointAllocation struct {
	gorm.Model
	RecruitID          uint
	TeamProfileID      int
	RecruitProfileID   int
	WeekID             int
	Points             float64
	RESAffectedPoints  float64
	AffinityOneApplied bool
	AffinityTwoApplied bool
	CaughtCheating     bool
}

func (rpa *RecruitPointAllocation) UpdatePointsSpent(points float64, res float64) {
	rpa.Points = points
	rpa.RESAffectedPoints = res
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
