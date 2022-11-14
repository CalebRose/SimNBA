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
	RecruitingSynced          bool
	GMActionsComplete         bool
	IsRecruitingLocked        bool
	AIBoardsCreated           bool
	AIPointAllocationComplete bool
	IsOffSeason               bool
}

func (t *Timestamp) MoveUpWeekCollege() {
	t.CollegeWeekID++
	t.CollegeWeek++
}

func (t *Timestamp) MoveUpWeekNBA() {
	t.NBAWeekID++
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

func (t *Timestamp) SyncToNextWeek() {
	t.MoveUpWeekCollege()
	t.MoveUpWeekNBA()
	// Reset Toggles
	// t.ToggleGamesARan()
	// t.ToggleGamesBRan()
	// t.ToggleGamesCRan()
	t.GamesARan = false
	t.GamesBRan = false
	t.RecruitingSynced = false
	t.AIPointAllocationComplete = false
	// t.ToggleGMActions()
}
