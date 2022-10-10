package structs

import "github.com/jinzhu/gorm"

type CollegeStandings struct {
	gorm.Model
	TeamID               uint
	TeamName             string
	SeasonID             uint
	Season               int
	ConferenceID         uint
	PostSeasonStatus     string
	IsConferenceChampion bool
	BaseStandings
}

func (cs *CollegeStandings) UpdateCollegeStandings(game Match) {
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
