package structs

import "github.com/jinzhu/gorm"

// DraftPick - a draftpick for a draft
type DraftPick struct {
	gorm.Model
	DraftRound     uint
	DraftNumber    uint
	DrafteeID      uint
	TeamID         uint
	Team           string
	OriginalTeamID uint
	OriginalTeam   string
	DraftValue     uint
}
