package structs

import (
	"github.com/jinzhu/gorm"
)

// PlayerRecruitProfile - The points allocated to one player
type PlayerRecruitProfile struct {
	gorm.Model
	SeasonID               int
	RecruitID              int
	ProfileID              int
	TotalPoints            int
	CurrentWeeksPoints     int
	SpendingCount          int
	Scholarship            bool
	ScholarshipRevoked     bool
	TeamAbbreviation       string
	InterestLevel          string
	InterestLevelThreshold int
	IsSigned               bool
	RemovedFromBoard       bool
	Recruit                Player `gorm:"foreignKey:PlayerID"`
}

func (r *PlayerRecruitProfile) AllocatePoints(points int) {
	r.CurrentWeeksPoints = points
}

func (r *PlayerRecruitProfile) SignPlayer() {
	if r.Scholarship {
		r.IsSigned = true
	}
}

func (r *PlayerRecruitProfile) AllocateTotalPoints(points int) {
	r.TotalPoints += r.CurrentWeeksPoints
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
