package structs

import (
	"database/sql"

	"gorm.io/gorm"
)

type CollegePromise struct {
	gorm.Model
	TeamID          uint
	CollegePlayerID uint
	PromiseType     string // Minutes (at least minimum), Wins (varies), March Madness (Medium), Conf Championship (High), Final Four (Very High), National Championship (very High), Gameplan Fit (medium), Adjust Gameplan (Low), Play Game In State
	PromiseWeight   string // The impact the promise will have on their decision. Low, Medium, High
	Benchmark       int    // The value that must be met. For wins & minutes
	BenchmarkStr    string // Needed value for benchmarks that are a string
	PromiseMade     bool   // The player has agreed to the premise of the promise
	IsFullfilled    bool   // If the promise was accomplished
	IsActive        bool   // Whether the promise is active
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

type TransferPortalBoardDto struct {
	Profiles []TransferPortalProfile
}

// Player Profile For the Transfer Portal?
type TransferPortalProfile struct {
	gorm.Model
	SeasonID              uint
	CollegePlayerID       uint
	ProfileID             uint
	PromiseID             sql.NullInt64
	TeamAbbreviation      string
	TotalPoints           float64
	CurrentWeeksPoints    int
	PreviouslySpentPoints int
	SpendingCount         int
	RemovedFromBoard      bool
	RolledOnPromise       bool
	LockProfile           bool
	IsSigned              bool
	Recruiter             string
	CollegePlayer         CollegePlayer  `gorm:"foreignKey:CollegePlayerID"`
	Promise               CollegePromise `gorm:"foreignKey:PromiseID"`
}

func (p *TransferPortalProfile) Reactivate() {
	p.RemovedFromBoard = false
}

func (p *TransferPortalProfile) SignPlayer() {
	p.IsSigned = true
	p.LockProfile = true
	p.CurrentWeeksPoints = 0
}

func (p *TransferPortalProfile) Lock() {
	p.CurrentWeeksPoints = 0
	p.LockProfile = true
}

func (p *TransferPortalProfile) Deactivate() {
	p.RemovedFromBoard = true
	p.CurrentWeeksPoints = 0
}

func (p *TransferPortalProfile) AllocatePoints(points int) {
	p.CurrentWeeksPoints = points
}

func (p *TransferPortalProfile) AddPointsToTotal(multiplier float64) {
	p.TotalPoints += (float64(p.CurrentWeeksPoints) * multiplier)
	if p.CurrentWeeksPoints == 0 {
		p.SpendingCount = 0
	} else {
		p.SpendingCount += 1
	}
}

func (p *TransferPortalProfile) AssignPromise(id uint) {
	p.PromiseID = sql.NullInt64{Valid: true, Int64: int64(id)}
}
func (p *TransferPortalProfile) ToggleRolledOnPromise() {
	p.RolledOnPromise = true
}

// Player Profile For the Transfer Portal?
type TransferPortalProfileResponse struct {
	ID                    uint
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
	RolledOnPromise       bool
	LockProfile           bool
	IsSigned              bool
	Recruiter             string
	CollegePlayer         TransferPlayerResponse `gorm:"foreignKey:CollegePlayerID"`
	Promise               CollegePromise         `gorm:"foreignKey:PromiseID"`
}

type TransferPortalResponse struct {
	Team         TeamRecruitingProfile
	TeamBoard    []TransferPortalProfileResponse
	TeamPromises []CollegePromise         // List of all promises
	Players      []TransferPlayerResponse // List of all Transfer Portal Players
	TeamList     []Team
}

// UpdateTransferPortalBoard - Data Transfer Object from UI to API
type UpdateTransferPortalBoard struct {
	Profile SimTeamBoardResponse
	Players []TransferPortalProfileResponse
	TeamID  int
}
