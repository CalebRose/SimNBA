package structs

import (
	"strings"

	"github.com/jinzhu/gorm"
)

type CollegeStandings struct {
	gorm.Model
	TeamID                  uint
	TeamName                string
	TeamAbbr                string
	SeasonID                uint
	Season                  int
	ConferenceID            uint
	ConferenceName          string
	PostSeasonStatus        string
	IsConferenceChampion    bool
	InvitationalParticipant bool
	Invitational            string
	InvitationalChampion    bool
	Rank                    uint
	BaseStandings
}

func (cs *CollegeStandings) UpdateCollegeStandings(game Match) {
	isAway := cs.TeamID == game.AwayTeamID
	winner := (!isAway && game.HomeTeamWin) || (isAway && game.AwayTeamWin)
	if winner {
		cs.TotalWins += 1
		if isAway {
			cs.AwayWins += 1
			if game.HomeTeamRank > 0 && !game.IsPlayoffGame {
				cs.RankedWins += 1
			}
		} else {
			cs.HomeWins += 1
		}
		if game.IsConference {
			cs.ConferenceWins += 1
		}
		cs.Streak += 1
		if game.IsInvitational && strings.Contains(game.MatchName, "Finals") && !strings.Contains(game.MatchName, "Semifinals") {
			cs.InvitationalChampion = true
		}
		if game.IsConferenceTournament && strings.Contains(game.MatchName, "Finals") && !strings.Contains(game.MatchName, "Semifinals") {
			cs.PostSeasonStatus = "Conference Champion"
			cs.IsConferenceChampion = true
		}
		if game.IsPlayoffGame {
			cs.PostSeasonStatus = game.MatchName
		}
		if game.IsNITGame {
			cs.PostSeasonStatus = game.MatchName
		}
		if game.IsCBIGame {
			cs.PostSeasonStatus = game.MatchName
		}
		if game.IsPlayoffGame && game.IsNationalChampionship {
			cs.PostSeasonStatus = "National Champion"
		}
	} else {
		cs.TotalLosses += 1
		cs.Streak = 0
		if isAway && game.HomeTeamRank > 0 && !game.IsPlayoffGame {
			cs.RankedLosses += 1
		}
		if !isAway && game.AwayTeamRank > 0 && !game.IsPlayoffGame {
			cs.RankedLosses += 1
		}
		if game.IsConference {
			cs.ConferenceLosses += 1
		}
		if game.IsPlayoffGame {
			cs.PostSeasonStatus = game.MatchName
		}
		if game.IsNationalChampionship {
			cs.PostSeasonStatus = "National Champion Runner-Up"
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

func (cs *CollegeStandings) RegressCollegeStandings(game Match) {
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

func (cs *CollegeStandings) UpdateCoach(coach string) {
	cs.Coach = coach
}

func (cs *CollegeStandings) AssignRank(rank int) {
	cs.Rank = uint(rank)
}

func (cs *BaseStandings) MaskGames(wins, losses, confWins, confLosses int) {
	cs.TotalWins = wins
	cs.TotalLosses = losses
	cs.ConferenceWins = confWins
	cs.ConferenceLosses = confLosses
}
