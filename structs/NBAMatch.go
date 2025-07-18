package structs

import "gorm.io/gorm"

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
	HomeTeamRank           uint
	AwayTeamID             uint
	AwayTeam               string
	AwayTeamCoach          string
	AwayTeamWin            bool
	AwayTeamRank           uint
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

func (m *NBAMatch) AddWeekData(id, week uint, timeslot string) {
	m.WeekID = id
	m.Week = week
	m.TimeSlot = timeslot
	m.MatchOfWeek = timeslot
}

func (m *NBAMatch) HideScore() {
	m.HomeTeamScore = 0
	m.AwayTeamScore = 0
	m.HomeTeamWin = false
	m.AwayTeamWin = false
}

func (m *NBAMatch) ResetScore() {
	m.HideScore()
	m.GameComplete = false
}

func (m *NBAMatch) SwapTeams() {
	tempID := m.HomeTeamID
	tempTeam := m.HomeTeam
	tempCoach := m.HomeTeamCoach
	tempRank := m.HomeTeamRank
	m.HomeTeamID = m.AwayTeamID
	m.HomeTeam = m.AwayTeam
	m.HomeTeamCoach = m.AwayTeamCoach
	m.HomeTeamRank = m.AwayTeamRank
	m.AwayTeamID = tempID
	m.AwayTeamCoach = tempCoach
	m.AwayTeam = tempTeam
	m.AwayTeamRank = tempRank
}

func (m *NBAMatch) AssignArena(arena, city, state string) {
	m.Arena = arena
	m.City = city
	m.State = state
}

func (m *NBAMatch) AddTeam(isHome bool, id, rank uint, team, coach, arena, city, state string) {
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
	if !m.IsNeutralSite && isHome {
		m.Arena = arena
		m.City = city
		m.State = state
	}
}

func (m *NBAMatch) RollbackMatch(ht, htc, at, atc, city, state, arena string, htid, atid uint) {
	m.HomeTeam = ht
	m.HomeTeamID = htid
	m.HomeTeamCoach = htc
	m.AwayTeam = at
	m.AwayTeamCoach = atc
	m.AwayTeamID = atid
	m.City = city
	m.State = state
	m.Arena = arena
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

func (s *NBASeries) AddTeam(isHome bool, id, rank uint, team, coach string) {
	if isHome {
		s.HomeTeam = team
		s.HomeTeamID = id
		s.HomeTeamRank = rank
		s.HomeTeamCoach = coach
	} else {
		s.AwayTeam = team
		s.AwayTeamID = id
		s.AwayTeamRank = rank
		s.AwayTeamCoach = coach
	}
	if s.HomeTeamID > 0 && s.AwayTeamID > 0 && s.HomeTeamRank > s.AwayTeamRank {
		tempID := s.AwayTeamID
		temp := s.AwayTeam
		tempC := s.AwayTeamCoach
		tempR := s.AwayTeamRank
		s.AwayTeamID = s.HomeTeamID
		s.AwayTeam = s.HomeTeam
		s.AwayTeamCoach = s.HomeTeamCoach
		s.AwayTeamRank = s.HomeTeamRank
		s.HomeTeamID = tempID
		s.HomeTeam = temp
		s.HomeTeamCoach = tempC
		s.HomeTeamRank = tempR
	}
	s.GameCount = 1
}

func (s *NBASeries) UpdateWinCount(id int) {
	if id == int(s.HomeTeamID) {
		s.HomeTeamWins += 1
	} else {
		s.AwayTeamWins += 1
	}
	if s.GameCount < 7 {
		s.GameCount += 1
	}
	if s.HomeTeamWins > 3 {
		s.HomeTeamWin = true
		s.SeriesComplete = true
	}
	if s.AwayTeamWins > 3 {
		s.AwayTeamWin = true
		s.SeriesComplete = true
	}

}
