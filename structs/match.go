package structs

import "github.com/jinzhu/gorm"

// Match - The Data Structure for a Game
type Match struct {
	gorm.Model
	WeekID        uint
	SeasonID      uint
	HomeTeamID    uint
	HomeTeam      string
	AwayTeamID    uint
	AwayTeam      string
	MatchOfWeek   string
	HomeTeamScore int
	AwayTeamScore int
	IsNeutral     bool
	IsNBAMatch    bool
	IsConference  bool
	IsDivisional  bool
}
