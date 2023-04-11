package structs

import "github.com/jinzhu/gorm"

// DraftPick - a draftpick for a draft
type DraftPick struct {
	gorm.Model
	SeasonID       uint
	Season         uint
	DraftRound     uint
	DraftNumber    uint
	DrafteeID      uint
	TeamID         uint
	Team           string
	OriginalTeamID uint
	OriginalTeam   string
	PreviousTeamID uint
	PreviousTeam   string
	DraftValue     uint
	Notes          string
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
