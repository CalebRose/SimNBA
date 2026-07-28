package structs

import "github.com/jinzhu/gorm"

type NBADraftPageResponse struct {
	WarRoom          NBAWarRoom
	DraftablePlayers []NBADraftee
	NBATeams         []NBATeam
	AllDraftPicks    [2][]DraftPick
	CollegeTeams     []Team
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
	ShowAgility          bool
	ShowBallwork         bool
	ShowRebounding       bool
	ShowStealing         bool
	ShowBlocking         bool
	ShowInteriorDefense  bool
	ShowPerimeterDefense bool
	ShowPotential        bool
	RemovedFromBoard     bool
	ShowCount            uint
	Draftee              NBADraftee `gorm:"foreignKey:PlayerID;references:PlayerID"`
}

func (sp *ScoutingProfile) RevealAttribute(attr string) {
	switch attr {
	case "ShowShooting2":
		sp.ShowShooting2 = true
	case "ShowShooting3":
		sp.ShowShooting3 = true
	case "ShowFreeThrow":
		sp.ShowFreeThrow = true
	case "ShowFinishing":
		sp.ShowFinishing = true
	case "ShowAgility":
		sp.ShowAgility = true
	case "ShowBallwork":
		sp.ShowBallwork = true
	case "ShowRebounding":
		sp.ShowRebounding = true
	case "ShowStealing":
		sp.ShowStealing = true
	case "ShowBlocking":
		sp.ShowBlocking = true
	case "ShowInteriorDefense":
		sp.ShowInteriorDefense = true
	case "ShowPerimeterDefense":
		sp.ShowPerimeterDefense = true
	case "ShowPotential":
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
