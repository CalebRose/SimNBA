package structs

import "gorm.io/gorm"

type NBAStandings struct {
	gorm.Model
	TeamID               uint
	TeamName             string
	TeamAbbr             string
	SeasonID             uint
	Season               int
	LeagueID             uint
	League               string
	ConferenceID         uint
	ConferenceName       string
	DivisionID           uint
	DivisionName         string
	PostSeasonStatus     string
	IsConferenceChampion bool
	BaseStandings
}

func (cs *NBAStandings) AssignID(id uint) {
	cs.ID = id
}

func (cs *NBAStandings) UpdateNBAStandings(game Match) {
	isAway := cs.TeamID == game.AwayTeamID
	winner := (!isAway && game.HomeTeamWin) || (isAway && game.AwayTeamWin)
	if winner {
		cs.TotalWins += 1
		if isAway {
			cs.AwayWins += 1
		} else {
			cs.HomeWins += 1
		}
		if game.IsConference {
			cs.ConferenceWins += 1
		}
		cs.Streak += 1
	} else {
		cs.TotalLosses += 1
		cs.Streak = 0
		if game.IsConference {
			cs.ConferenceLosses += 1
		}
	}
	if isAway {
		cs.PointsFor += game.AwayTeamScore
		cs.PointsAgainst += game.HomeTeamScore
	} else {
		cs.PointsFor += game.HomeTeamScore
		cs.PointsAgainst += game.AwayTeamScore
	}
}

func (cs *NBAStandings) RegressNBAStandings(game Match) {
	isAway := cs.TeamID == game.AwayTeamID
	winner := (!isAway && game.HomeTeamWin) || (isAway && game.AwayTeamWin)
	if winner {
		cs.TotalWins -= 1
		if isAway {
			cs.AwayWins -= 1
		} else {
			cs.HomeWins -= 1
		}
		if game.IsConference {
			cs.ConferenceWins -= 1
		}
		cs.Streak -= 1
	} else {
		cs.TotalLosses -= 1
		cs.Streak = 0
		if game.IsConference {
			cs.ConferenceLosses -= 1
		}
	}
	if isAway {
		cs.PointsFor -= game.AwayTeamScore
		cs.PointsAgainst -= game.HomeTeamScore
	} else {
		cs.PointsFor -= game.HomeTeamScore
		cs.PointsAgainst -= game.AwayTeamScore
	}
}

func (cs *NBAStandings) UpdateCoach(coach string) {
	cs.Coach = coach
}
