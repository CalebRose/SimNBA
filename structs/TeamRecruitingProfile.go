package structs

import "gorm.io/gorm"

// TeamRecruitingProfile - The profile for a team for recruiting
type TeamRecruitingProfile struct {
	gorm.Model
	TeamID                  uint
	TeamAbbr                string
	State                   string
	Region                  string
	ScholarshipsAvailable   int
	WeeklyPoints            int
	BonusPoints             int
	SpentPoints             int
	TotalCommitments        int
	RecruitClassSize        int
	PortalReputation        int // A value between 1-100 signifying the coach's reputation and behavior in the transfer portal.
	IsAI                    bool
	AIBehavior              string // Aggressive, Normal, Conservative -- will be for determining how likely they'll generate a good coach
	AIQuality               string // Blue Blood, P6, Cinderella, Mid-Major
	AIValue                 string // Star, Talent, Potential
	AIPrestige              string // A new column which will be used to determine how likely a school/boosters will fire a coach pending on how they do. "Very High", "High", "Average", "Low", "Very Low"
	AIAttribute1            string
	AIAttribute2            string
	ESPNScore               float64
	RivalsScore             float64
	Rank247Score            float64
	CompositeScore          float64
	AIMinThreshold          int
	AIMaxThreshold          int
	AIStarMin               int
	AIStarMax               int
	AIAutoOfferscholarships bool
	OffensiveScheme         string
	DefensiveScheme         string
	Recruiter               string
	CaughtCheating          bool
	Recruits                []PlayerRecruitProfile `gorm:"foreignKey:ProfileID"`
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

func (r *TeamRecruitingProfile) SetClassSize(size int) {
	r.RecruitClassSize = size
}

func (r *TeamRecruitingProfile) SetNewBehaviors(value, attr1, attr2 string) {
	r.AIValue = value
	r.AIAttribute1 = attr1
	r.AIAttribute2 = attr2
}

func (r *TeamRecruitingProfile) UpdateAIBehavior(isAi, autoOffer bool, starMax, starMin, min, max int, offScheme, defScheme string) {
	r.IsAI = isAi
	r.AIAutoOfferscholarships = autoOffer
	r.AIStarMax = starMax
	r.AIStarMin = starMin
	r.AIMinThreshold = min
	r.AIMaxThreshold = max
	r.OffensiveScheme = offScheme
	r.DefensiveScheme = defScheme
}

func (r *TeamRecruitingProfile) AssignRecruiter(name string) {
	r.Recruiter = name
}

func (r *TeamRecruitingProfile) IncreaseClassSize() {
	r.RecruitClassSize += 1
}

func (r *TeamRecruitingProfile) AssignBonusPoints(bonus int) {
	r.BonusPoints = 0
}
