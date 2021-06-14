package structs

import "github.com/jinzhu/gorm"

// RecruitingProfile - The profile for a team for recruiting
type RecruitingProfile struct {
	gorm.Model
	TeamID                int
	Team                  string
	State                 string
	Region                string
	ScholarshipsAvailable int
	WeeklyPoints          int
	BonusPoints           int
	SpentPoints           int
	Recruits              []RecruitingPoints `gorm:"foreignKey:ProfileID"`
}

func (r *RecruitingProfile) SubtractScholarshipsAvailable() {
	r.ScholarshipsAvailable--
}

func (r *RecruitingProfile) ReallocateScholarship() {
	r.ScholarshipsAvailable++
}

func (r *RecruitingProfile) AllocateSpentPoints(points int) {
	r.SpentPoints = points
}
