package structs

import "github.com/jinzhu/gorm"

// NBAMatch - The Data Structure for a Game
type NBAMatch struct {
	gorm.Model
	MatchName              string // For Post-Season matchups
	WeekID                 uint
	Week                   uint
	SeasonID               uint
	SeriesID               uint
	HomeTeamID             uint
	HomeTeam               string
	HomeTeamCoach          string
	HomeTeamWin            bool
	AwayTeamID             uint
	AwayTeam               string
	AwayTeamCoach          string
	AwayTeamWin            bool
	MatchOfWeek            string
	HomeTeamScore          int
	AwayTeamScore          int
	NextGameID             uint
	NextGameHOA            string
	TimeSlot               string
	Arena                  string
	City                   string
	State                  string
	Country                string
	IsNeutralSite          bool
	IsConference           bool
	IsDivisional           bool
	IsConferenceTournament bool
	IsInternational        bool
	IsPlayoffGame          bool
	IsTheFinals            bool
	IsRivalryGame          bool
	GameComplete           bool
}

func (m *NBAMatch) UpdateScore(homeTeamScore int, awayTeamScore int) {
	m.HomeTeamScore = homeTeamScore
	m.AwayTeamScore = awayTeamScore
	if m.HomeTeamScore > m.AwayTeamScore {
		m.HomeTeamWin = true
	} else {
		m.AwayTeamWin = true
	}
	m.GameComplete = true
}

func (m *NBAMatch) UpdateCoach(TeamID int, Username string) {
	if m.HomeTeamID == uint(TeamID) {
		m.HomeTeamCoach = Username
	} else if m.AwayTeamID == uint(TeamID) {
		m.AwayTeamCoach = Username
	}
}

type NBASeries struct {
	gorm.Model
	SeriesName      string // For Post-Season matchups
	SeasonID        uint
	HomeTeamID      uint
	HomeTeam        string
	HomeTeamCoach   string
	HomeTeamWins    uint
	HomeTeamWin     bool
	HomeTeamRank    uint
	AwayTeamID      uint
	AwayTeam        string
	AwayTeamCoach   string
	AwayTeamWins    uint
	AwayTeamWin     bool
	AwayTeamRank    uint
	GameCount       uint
	NextSeriesID    uint
	NextSeriesHOA   string
	IsInternational bool
	IsPlayoffGame   bool
	IsTheFinals     bool
	SeriesComplete  bool
}
