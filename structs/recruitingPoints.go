package structs

import (
	"github.com/jinzhu/gorm"
)

// RecruitingPoints - The points allocated to one player
type RecruitingPoints struct {
	gorm.Model
	SeasonID               int
	PlayerID               int
	ProfileID              int
	TotalPointsSpent       int
	CurrentPointsSpent     int
	SpendingCount          int
	Scholarship            bool
	Team                   string
	InterestLevel          string
	InterestLevelThreshold int
	Signed                 bool
	RemovedFromBoard       bool
	Recruit                Player `gorm:"foreignKey:PlayerID"`
}

func (r *RecruitingPoints) AllocatePoints(points int) {
	r.CurrentPointsSpent = points
}

func (r *RecruitingPoints) SignPlayer() {
	if r.Scholarship {
		r.Signed = true
	}
}

func (r *RecruitingPoints) AllocateTotalPoints(points int) {
	r.TotalPointsSpent += r.CurrentPointsSpent
	r.CurrentPointsSpent = 0
}

func (r *RecruitingPoints) AllocateScholarship() {
	r.Scholarship = true
}

func (r *RecruitingPoints) RevokeScholarship() {
	r.Scholarship = false
}

func (r *RecruitingPoints) RemoveRecruitFromBoard() {
	r.RemovedFromBoard = true
}

func (r *RecruitingPoints) ReplaceRecruitToBoard() {
	r.RemovedFromBoard = false
}
