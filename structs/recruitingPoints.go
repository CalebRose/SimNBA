package structs

import (
	"fmt"

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
	InterestLevel          string
	InterestLevelThreshold int
	Signed                 bool
	Recruit                Player `gorm:"foreignKey:PlayerID"`
}

func (r *RecruitingPoints) AllocatePoints(points int) {
	if r.Scholarship == true {
		r.CurrentPointsSpent = points
	} else {
		fmt.Println("Cannot allocate points without offering a scholarship")
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
