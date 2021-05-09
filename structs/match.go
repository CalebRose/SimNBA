package structs

import "github.com/jinzhu/gorm"

// Match - The Data Structure for a Game
type Match struct {
	gorm.Model
	WeekID        int
	SeasonID      int
	HomeTeamID    int
	HomeTeam      string
	AwayTeamID    int
	AwayTeam      string
	MatchOfWeek   string
	HomeTeamScore int
	AwayTeamScore int
	IsNeutral     bool
	IsNBAMatch    bool
	IsConference  bool
	IsDivisional  bool
}
