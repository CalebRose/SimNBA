package structs

import "github.com/jinzhu/gorm"

// TeamStats -- Statistics related to the performance of a team during a match
type TeamStats struct {
	gorm.Model
	TeamID                    uint
	MatchID                   uint
	SeasonID                  uint
	WeekID                    uint
	Points                    int
	Possessions               int
	FGM                       int
	FGA                       int
	FGPercent                 float64
	ThreePointsMade           int
	ThreePointAttempts        int
	ThreePointPercent         float64
	FTM                       int
	FTA                       int
	FTPercent                 float64
	Rebounds                  int
	OffRebounds               int
	DefRebounds               int
	Assists                   int
	Steals                    int
	Blocks                    int
	TotalTurnovers            int
	LargestLead               int
	FirstHalfScore            int
	SecondHalfScore           int
	Fouls                     int
	PointsAgainst             int
	FGMAgainst                int
	FGAAgainst                int
	FGPercentAgainst          float64
	ThreePointsMadeAgainst    int
	ThreePointAttemptsAgainst int
	ThreePointPercentAgainst  float64
	FTMAgainst                int
	FTAAgainst                int
	FTPercentAgainst          float64
	ReboundsAllowed           int
	OffReboundsAllowed        int
	DefReboundsAllowed        int
	AssistsAllowed            int
	StealsAllowed             int
	BlocksAllowed             int
	TurnoversAllowed          int
}
