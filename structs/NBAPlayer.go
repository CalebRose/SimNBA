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
	// Contract
	// NBAPlayerStats
}
