package structs

import "github.com/jinzhu/gorm"

type NBAPlayer struct {
	gorm.Model
	BasePlayer
	PlayerID            int
	TeamID              int
	TeamAbbr            string
	IsNBA               bool
	IsSuperMaxQualified bool
	IsFreeAgent         bool
	CollegeID           int
	College             string
	DraftPickID         int
	DraftPick           string
	DraftedTeamID       int
	DraftedTeamAbbr     string
	// Contract
	// NBAPlayerStats
}
