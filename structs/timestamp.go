package structs

import "github.com/jinzhu/gorm"

// Timestamp - The Global Timestamp for the Season
type Timestamp struct {
	gorm.Model
	SeasonID                      uint
	Season                        int
	CollegeWeekID                 uint
	NBAWeekID                     uint
	CollegeWeek                   int
	NBAWeek                       int
	GamesARan                     bool
	GamesBRan                     bool
	GamesCRan                     bool
	GamesDRan                     bool
	CollegePollRan                bool
	RecruitingSynced              bool
	GMActionsComplete             bool
	IsRecruitingLocked            bool
	AIBoardsCreated               bool
	AIPointAllocationComplete     bool
	IsOffSeason                   bool
	IsNBAOffseason                bool
	NBAPreseason                  bool
	IsFreeAgencyLocked            bool
	IsDraftTime                   bool
	ProgressedCollegePlayers      bool
	ProgressedProfessionalPlayers bool
	CollegeSeasonOver             bool
	NBASeasonOver                 bool
	CrootsGenerated               bool
	Y1Capspace                    float64
	Y2Capspace                    float64
	Y3Capspace                    float64
	Y4Capspace                    float64
	Y5Capspace                    float64
	FreeAgencyRound               uint
	RunCron                       bool
	TransferPortalPhase           uint
	TransferPortalRound           uint
}

func (t *Timestamp) MoveUpASeason() {
	t.SeasonID++
	t.Season++
	t.CollegeSeasonOver = false
	t.NBASeasonOver = false
	t.CollegeWeek = 0
	t.CollegeWeekID += 1
	t.NBAWeek = 0
	t.NBAWeekID += 1
	t.TransferPortalPhase = 0
	t.CrootsGenerated = false
	t.Y1Capspace = t.Y2Capspace
	t.Y2Capspace = t.Y3Capspace
	t.Y3Capspace = t.Y4Capspace
	t.Y4Capspace = t.Y5Capspace
	t.Y5Capspace = t.Y5Capspace + 5
	t.FreeAgencyRound = 1
}

func (t *Timestamp) MoveUpWeekCollege() {
	t.CollegeWeekID++
	t.CollegeWeek++
}

func (t *Timestamp) MoveUpWeekNBA() {
	t.NBAWeekID++
	t.NBAWeek++
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
	t.IsRecruitingLocked = false
}

func (t *Timestamp) ToggleAIRecruiting() {
	t.AIPointAllocationComplete = !t.AIPointAllocationComplete
}

func (t *Timestamp) ToggleLockRecruiting() {
	t.IsRecruitingLocked = !t.IsRecruitingLocked
}

func (t *Timestamp) ToggleGeneratedCroots() {
	t.CrootsGenerated = !t.CrootsGenerated
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

func (t *Timestamp) ToggleCollegeProgression() {
	t.ProgressedCollegePlayers = !t.ProgressedCollegePlayers
}

func (t *Timestamp) ToggleProfessionalProgression() {
	t.ProgressedProfessionalPlayers = !t.ProgressedProfessionalPlayers
	t.IsFreeAgencyLocked = false
	t.IsDraftTime = true
}

func (t *Timestamp) SyncToNextWeek() {
	if t.CollegeWeek < 21 {
		t.MoveUpWeekCollege()
	}
	if !t.IsNBAOffseason {
		t.MoveUpWeekNBA()
		t.GMActionsComplete = false
	}
	if !t.IsOffSeason || t.CollegeWeek < 21 {
		t.RecruitingSynced = false
		t.AIPointAllocationComplete = false
		t.TogglePollRan()
	}

	if !t.IsOffSeason || !t.IsNBAOffseason {
		// Reset Toggles
		t.GamesARan = false
		t.GamesBRan = false
		t.GamesCRan = false
		t.GamesDRan = false
	}

	if t.CollegeSeasonOver && t.NBASeasonOver {
		t.MoveUpASeason()
	}
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

func (t *Timestamp) EndTheCollegeSeason() {
	t.IsOffSeason = true
	t.TransferPortalPhase = 1
	t.CollegeSeasonOver = true
}

func (t *Timestamp) ClosePortal() {
	t.TransferPortalPhase = 0
}

func (t *Timestamp) EnactPromisePhase() {
	t.TransferPortalPhase = 2
}

func (t *Timestamp) EnactPortalPhase() {
	t.TransferPortalPhase = 3
}

func (t *Timestamp) IncrementTransferPortalRound() {
	t.IsRecruitingLocked = false
	if t.TransferPortalRound < 10 {
		t.TransferPortalRound += 1
	}
}

func (t *Timestamp) EndTheProfessionalSeason() {
	t.IsNBAOffseason = true
	t.FreeAgencyRound = 1
	t.IsDraftTime = false
	t.IsFreeAgencyLocked = true
	t.NBASeasonOver = true
}
