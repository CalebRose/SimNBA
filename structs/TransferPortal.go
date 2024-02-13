package structs

import "gorm.io/gorm"

type CollegePromise struct {
	gorm.Model
	TeamID          uint
	CollegePlayerID uint
	PromiseType     string // Minutes (at least minimum), Wins (varies), March Madness (Medium), Conf Championship (High), Final Four (Very High), National Championship (very High), Gameplan Fit (medium), Adjust Gameplan (Low)
	PromiseWeight   string // The impact the promise will have on their decision. Low, Medium, High
	Benchmark       int    // The value that must be met. For wins & minutes
	PromiseMade     bool   // The player has agreed to the premise of the promise
	IsFullfilled    bool   // If the promise was accomplished
	IsActive        bool   //
}

func (p *CollegePromise) Reactivate(promtype, weight string, benchmark int) {
	p.IsActive = true
	p.PromiseType = promtype
	p.PromiseWeight = weight
	p.Benchmark = benchmark
}

func (p *CollegePromise) UpdatePromise(promtype, weight string, benchmark int) {
	p.PromiseType = promtype
	p.PromiseWeight = weight
	p.Benchmark = benchmark
}

func (p *CollegePromise) Deactivate() {
	p.IsActive = false
}

func (p *CollegePromise) MakePromise() {
	p.PromiseMade = true
}

func (p *CollegePromise) FulfillPromise() {
	p.IsFullfilled = true
}

// Player Profile For the Transfer Portal?
type TransferPortalProfile struct {
	gorm.Model
	SeasonID              uint
	CollegePlayerID       uint
	ProfileID             uint
	PromiseID             uint
	TeamAbbreviation      string
	TotalPoints           float64
	CurrentWeeksPoints    int
	PreviouslySpentPoints int
	SpendingCount         int
	RemovedFromBoard      bool
	CollegePlayer         CollegePlayer  `gorm:"foreignKey:CollegePlayerID"`
	Promise               CollegePromise `gorm:"foreignKey:PromiseID"`
}

type TransferPortalResponse struct {
	Team         TeamRecruitingProfile
	TeamBoard    []TransferPortalProfile
	TeamPromises []CollegePromise // List of all promises
	Players      []CollegePlayer  // List of all Transfer Portal Players
}
