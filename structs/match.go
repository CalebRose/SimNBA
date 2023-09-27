package structs

import "github.com/jinzhu/gorm"

// MatchPageResponse - The Data Structure for the Schedule Page
type MatchPageResponse struct {
	CBBGames []Match
	NBAGames []NBAMatch
	ISLGames []NBAMatch
}

// Match - The Data Structure for a Game
type Match struct {
	gorm.Model
	MatchName              string // For Post-Season matchups
	WeekID                 uint
	Week                   uint
	SeasonID               uint
	HomeTeamID             uint
	HomeTeam               string
	HomeTeamCoach          string
	HomeTeamWin            bool
	HomeTeamRank           uint
	AwayTeamID             uint
	AwayTeam               string
	AwayTeamCoach          string
	AwayTeamWin            bool
	AwayTeamRank           uint
	MatchOfWeek            string
	HomeTeamScore          int
	AwayTeamScore          int
	TimeSlot               string
	Arena                  string
	City                   string
	State                  string
	NextGameID             uint
	NextGameHOA            string
	IsNeutralSite          bool
	IsNBAMatch             bool
	IsConference           bool
	IsConferenceTournament bool
	IsNITGame              bool
	IsPlayoffGame          bool
	IsNationalChampionship bool
	IsRivalryGame          bool
	IsInvitational         bool
	GameComplete           bool
}

func (m *Match) UpdateScore(homeTeamScore int, awayTeamScore int) {
	m.HomeTeamScore = homeTeamScore
	m.AwayTeamScore = awayTeamScore
	if m.HomeTeamScore > m.AwayTeamScore {
		m.HomeTeamWin = true
	} else {
		m.AwayTeamWin = true
	}
	m.GameComplete = true
}

func (m *Match) UpdateCoach(TeamID int, Username string) {
	if m.HomeTeamID == uint(TeamID) {
		m.HomeTeamCoach = Username
	} else if m.AwayTeamID == uint(TeamID) {
		m.AwayTeamCoach = Username
	}
}

func (m *Match) AddTeam(isHome bool, id, rank uint, team, coach string) {
	if isHome {
		m.HomeTeam = team
		m.HomeTeamID = id
		m.HomeTeamRank = rank
		m.HomeTeamCoach = coach
	} else {
		m.AwayTeam = team
		m.AwayTeamID = id
		m.AwayTeamRank = rank
		m.AwayTeamCoach = coach
	}
}

func (m *Match) AssignRank(id, rank uint) {
	isHome := id == m.HomeTeamID
	if isHome {
		m.HomeTeamRank = rank
	} else {
		m.AwayTeamRank = rank
	}
}
