package structs

import "github.com/jinzhu/gorm"

// TeamRecruitingProfile - The profile for a team for recruiting
type TeamRecruitingProfile struct {
	gorm.Model
	TeamID                uint
	TeamAbbr              string
	State                 string
	Region                string
	ScholarshipsAvailable int
	WeeklyPoints          int
	BonusPoints           int
	SpentPoints           int
	TotalCommitments      int
	RecruitClassSize      int
	IsAI                  bool
	AIBehavior            string
	AIQuality             string
	ESPNScore             float64
	RivalsScore           float64
	Rank247Score          float64
	CompositeScore        float64
	CaughtCheating        bool
	Recruits              []PlayerRecruitProfile `gorm:"foreignKey:ProfileID"`
}

func (r *TeamRecruitingProfile) SubtractScholarshipsAvailable() {
	r.ScholarshipsAvailable--
}

func (r *TeamRecruitingProfile) ReallocateScholarship() {
	r.ScholarshipsAvailable++
}

func (r *TeamRecruitingProfile) AllocateSpentPoints(points int) {
	r.SpentPoints = points
}
