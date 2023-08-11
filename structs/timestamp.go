package structs

import "github.com/jinzhu/gorm"

// Timestamp - The Global Timestamp for the Season
type Timestamp struct {
	gorm.Model
	SeasonID                  uint
	Season                    int
	CollegeWeekID             uint
	NBAWeekID                 uint
	CollegeWeek               int
	NBAWeek                   int
	GamesARan                 bool
	GamesBRan                 bool
	GamesCRan                 bool
	GamesDRan                 bool
	CollegePollRan            bool
	RecruitingSynced          bool
	GMActionsComplete         bool
	IsRecruitingLocked        bool
	AIBoardsCreated           bool
	AIPointAllocationComplete bool
	IsOffSeason               bool
	IsNBAOffseason            bool
	NBAPreseason              bool
	IsFreeAgencyLocked        bool
	IsDraftTime               bool
	Y1Capspace                float64
	Y2Capspace                float64
	Y3Capspace                float64
	Y4Capspace                float64
	Y5Capspace                float64
	FreeAgencyRound           uint
}

func (t *Timestamp) MoveUpWeekCollege() {
	t.CollegeWeekID++
	t.CollegeWeek++
}

func (t *Timestamp) MoveUpWeekNBA() {
	t.NBAWeekID++
}

func (t *Timestamp) ToggleGames(matchType string) {
	if matchType == "A" {
		t.GamesARan = true
	} else if matchType == "B" {
		t.GamesBRan = true
	} else if matchType == "C" {
		t.GamesCRan = true
	} else if matchType == "D" {
		t.GamesDRan = true
	}
}

func (t *Timestamp) ToggleGamesARan() {
	t.GamesARan = !t.GamesARan
}

func (t *Timestamp) ToggleGamesBRan() {
	t.GamesBRan = !t.GamesBRan
}

func (t *Timestamp) ToggleGamesCRan() {
	t.GamesCRan = !t.GamesCRan
}

func (t *Timestamp) ToggleRecruiting() {
	t.RecruitingSynced = !t.RecruitingSynced
}

func (t *Timestamp) ToggleAIRecruiting() {
	t.AIPointAllocationComplete = !t.AIPointAllocationComplete
}

func (t *Timestamp) ToggleLockRecruiting() {
	t.IsRecruitingLocked = !t.IsRecruitingLocked
}

func (t *Timestamp) ToggleGMActions() {
	t.GMActionsComplete = !t.GMActionsComplete
}

func (t *Timestamp) TogglePollRan() {
	t.CollegePollRan = !t.CollegePollRan
}

func (t *Timestamp) ToggleFALock() {
	t.IsFreeAgencyLocked = !t.IsFreeAgencyLocked
}

func (t *Timestamp) SyncToNextWeek() {
	t.MoveUpWeekCollege()
	t.MoveUpWeekNBA()
	// Reset Toggles
	t.GamesARan = false
	t.GamesBRan = false
	t.GamesCRan = false
	t.GamesDRan = false
	t.RecruitingSynced = false
	t.AIPointAllocationComplete = false
	t.GMActionsComplete = false
	t.TogglePollRan()
}

func (t *Timestamp) MoveUpFreeAgencyRound() {
	t.FreeAgencyRound++
	if t.FreeAgencyRound > 10 {
		t.FreeAgencyRound = 0
		t.IsFreeAgencyLocked = true
		// t.IsDraftTime = true
	}
}

func (t *Timestamp) ToggleDraftTime() {
	t.IsDraftTime = !t.IsDraftTime
	// t.IsNBAOffseason = false
}
