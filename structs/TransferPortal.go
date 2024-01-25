package structs

import "gorm.io/gorm"

type CollegePromise struct {
	gorm.Model
	TeamID          uint
	CollegePlayerID uint
	PromiseType     string // Minutes, Wins, March Madness, Conf Championship, Final Four, National Championship
	PromiseWeight   string // The impact the promise will have on their decision
	Benchmark       int    // The value that must be met
	PromiseMade     bool   // The player has agreed to the premise of the promise
	IsFullfilled    bool
}

// Player Profile For the Transfer Portal?
type TransferPortalProfile struct {
	gorm.Model
	SeasonID              uint
	CollegePlayerID       uint
	ProfileID             uint
	TeamAbbreviation      string
	TotalPoints           float64
	CurrentWeeksPoints    int
	PreviouslySpentPoints int
	SpendingCount         int
	RemovedFromBoard      bool
	CollegePlayer         CollegePlayer `gorm:"foreignKey:CollegePlayerID"`
	Promise               CollegePromise
}

type TransferPortalResponse struct {
	Team         TeamRecruitingProfile
	TeamBoard    []TransferPortalProfile
	TeamPromises []CollegePromise // List of all promises
	Players      []CollegePlayer  // List of all Transfer Portal Players
}
