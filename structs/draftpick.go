package structs

import "gorm.io/gorm"

// DraftPick - a draftpick for a draft
type DraftPick struct {
	gorm.Model
	SeasonID               uint
	Season                 uint
	DraftRound             uint
	DraftNumber            uint
	DrafteeID              uint
	TeamID                 uint
	Team                   string
	OriginalTeamID         uint
	OriginalTeam           string
	PreviousTeamID         uint
	PreviousTeam           string
	DraftValue             uint
	Notes                  string
	SelectedPlayerID       uint
	SelectedPlayerName     string
	SelectedPlayerPosition string
}

func (p *DraftPick) AssignDraftNumber(num uint) {
	p.DraftNumber = num
}

func (p *DraftPick) TradePick(id uint, team string) {
	p.PreviousTeamID = p.TeamID
	p.PreviousTeam = p.Team
	p.TeamID = id
	p.Team = team
	if p.PreviousTeamID == p.OriginalTeamID {
		p.Notes = "From " + p.OriginalTeam
	} else {
		p.Notes = "From " + p.PreviousTeam + " via " + p.OriginalTeam
	}
}

type DraftLottery struct {
	ID            uint
	Team          string
	Chances       []uint
	CurrentChance uint
	Selection     uint
}

func (dl *DraftLottery) ApplyCurrentChance(pick int) {
	dl.CurrentChance = dl.Chances[pick]
}

// Sorting Funcs
type ByDraftChance []DraftLottery

func (fo ByDraftChance) Len() int      { return len(fo) }
func (fo ByDraftChance) Swap(i, j int) { fo[i], fo[j] = fo[j], fo[i] }
func (fo ByDraftChance) Less(i, j int) bool {
	return fo[i].CurrentChance < fo[j].CurrentChance
}

type ByDraftNumber []DraftPick

func (fo ByDraftNumber) Len() int      { return len(fo) }
func (fo ByDraftNumber) Swap(i, j int) { fo[i], fo[j] = fo[j], fo[i] }
func (fo ByDraftNumber) Less(i, j int) bool {
	return fo[i].DraftNumber < fo[j].DraftNumber
}
