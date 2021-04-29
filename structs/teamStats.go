package structs

import "github.com/jinzhu/gorm"

// TeamStats -- Statistics related to the performance of a team during a match
type TeamStats struct {
	gorm.Model
	TeamID             int
	MatchID            int
	SeasonID           int
	Points             int
	Possessions        int
	FGM                int
	FGA                int
	FGPercent          float32
	ThreePointsMade    int
	ThreePointAttempts int
	ThreePointPercent  float32
	FTM                int
	FTA                int
	FTPercent          float32
	Rebounds           int
	OffRebounds        int
	DefRebounds        int
	Assists            int
	Steals             int
	Blocks             int
	TotalTurnovers     int
	LargestLead        int
	FirstHalfScore     int
	SecondHalfScore    int
}
