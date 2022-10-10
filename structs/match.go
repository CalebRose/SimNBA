package structs

import "github.com/jinzhu/gorm"

// Match - The Data Structure for a Game
type Match struct {
	gorm.Model
	MatchName                string // For Post-Season matchups
	WeekID                   uint
	SeasonID                 uint
	HomeTeamID               uint
	HomeTeam                 string
	HomeTeamCoach            string
	HomeTeamWin              bool
	AwayTeamID               uint
	AwayTeam                 string
	AwayTeamCoach            string
	AwayTeamWin              bool
	MatchOfWeek              string
	HomeTeamScore            int
	AwayTeamScore            int
	TimeSlot                 string
	Stadium                  string
	City                     string
	State                    string
	IsNeutral                bool
	IsNBAMatch               bool
	IsConference             bool
	IsConferenceChampionship bool
	IsBowlGame               bool
	IsPlayoffGame            bool
	IsNationalChampionship   bool
	IsRivalryGame            bool
	GameComplete             bool
}
