package structs

import "gorm.io/gorm"

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

func (r *TeamRecruitingProfile) ResetSpentPoints() {
	r.SpentPoints = 0
}

func (r *TeamRecruitingProfile) SubtractScholarshipsAvailable() {
	if r.ScholarshipsAvailable > 0 {
		r.ScholarshipsAvailable--
	}
}

func (r *TeamRecruitingProfile) ReallocateScholarship() {
	if r.ScholarshipsAvailable < 15 {
		r.ScholarshipsAvailable++
	}
}

func (r *TeamRecruitingProfile) ResetScholarshipCount() {
	r.ScholarshipsAvailable = 15
}

func (r *TeamRecruitingProfile) AllocateSpentPoints(points int) {
	r.SpentPoints = points
}

func (r *TeamRecruitingProfile) AIAllocateSpentPoints(points int) {
	r.SpentPoints += points
}

func (r *TeamRecruitingProfile) ResetWeeklyPoints(points int) {
	r.WeeklyPoints = points
}

func (r *TeamRecruitingProfile) AddRecruitsToProfile(croots []PlayerRecruitProfile) {
	r.Recruits = croots
}

func (r *TeamRecruitingProfile) AssignRivalsRank(score float64) {
	r.RivalsScore = score
}

func (r *TeamRecruitingProfile) Assign247Rank(score float64) {
	r.Rank247Score = score
}

func (r *TeamRecruitingProfile) AssignESPNRank(score float64) {
	r.ESPNScore = score
}

func (r *TeamRecruitingProfile) AssignCompositeRank(score float64) {
	r.CompositeScore = score
}

func (r *TeamRecruitingProfile) UpdateTotalSignedRecruits(num int) {
	r.TotalCommitments = num
}

func (r *TeamRecruitingProfile) IncreaseCommitCount() {
	r.TotalCommitments++
}

func (r *TeamRecruitingProfile) ApplyCaughtCheating() {
	r.CaughtCheating = true
}

func (r *TeamRecruitingProfile) ToggleAIBehavior(val bool) {
	r.IsAI = val
}
