package structs

import "github.com/jinzhu/gorm"

type NBADraftPageResponse struct {
	WarRoom          NBAWarRoom
	DraftablePlayers []NBADraftee
	NBATeams         []NBATeam
	AllDraftPicks    []DraftPick
}

type NBAWarRoom struct {
	gorm.Model
	TeamID         uint
	Team           string
	ScoutingPoints uint
	SpentPoints    uint
	DraftPicks     []DraftPick       `gorm:"foreignKey:TeamID"`
	ScoutProfiles  []ScoutingProfile `gorm:"foreignKey:TeamID"`
}

func (w *NBAWarRoom) ResetSpentPoints() {
	w.SpentPoints = 0
}

func (w *NBAWarRoom) AddToSpentPoints(points uint) {
	w.SpentPoints += points
}

type ScoutingProfile struct {
	gorm.Model
	PlayerID             uint
	TeamID               uint
	ShowShooting2        bool
	ShowShooting3        bool
	ShowFreeThrow        bool
	ShowFinishing        bool
	ShowBallwork         bool
	ShowRebounding       bool
	ShowInteriorDefense  bool
	ShowPerimeterDefense bool
	ShowPotential        bool
	RemovedFromBoard     bool
	ShowCount            uint
	Draftee              NBADraftee `gorm:"foreignKey:PlayerID;references:PlayerID"`
}

func (sp *ScoutingProfile) RevealAttribute(attr string) {
	if attr == "Shooting2" {
		sp.ShowShooting2 = true
	} else if attr == "Shooting3" {
		sp.ShowShooting3 = true
	} else if attr == "FreeThrow" {
		sp.ShowFreeThrow = true
	} else if attr == "Finishing" {
		sp.ShowFinishing = true
	} else if attr == "Ballwork" {
		sp.ShowBallwork = true
	} else if attr == "Rebounding" {
		sp.ShowRebounding = true
	} else if attr == "InteriorDefense" {
		sp.ShowInteriorDefense = true
	} else if attr == "PerimeterDefense" {
		sp.ShowPerimeterDefense = true
	} else if attr == "Potential" {
		sp.ShowPotential = true
	}
	sp.ShowCount++
}

func (sp *ScoutingProfile) RemoveFromBoard() {
	sp.RemovedFromBoard = true
}

func (sp *ScoutingProfile) ReplaceOnBoard() {
	sp.RemovedFromBoard = false
}

type ScoutingProfileDTO struct {
	PlayerID uint
	TeamID   uint
}

type RevealAttributeDTO struct {
	ScoutProfileID uint
	Attribute      string
	Points         uint
	TeamID         uint
}

type ScoutingDataResponse struct {
	DrafteeSeasonStats CollegePlayerSeasonStats
	TeamStandings      CollegeStandings
}

type ExportDraftPicksDTO struct {
	DraftPicks []DraftPick
}
