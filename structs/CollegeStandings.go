package structs

import (
	"fmt"
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
	PreseasonRank           uint16
	ToucanRank              uint16
	KenPomRank              uint16
	KenPomRating            float32
	RPIRank                 uint16
	RPIRating               float32
	SOS                     float32
	SOR                     float32
	Q1Wins                  uint8
	Q2Wins                  uint8
	Q3Wins                  uint8
	Q4Wins                  uint8
	Q1Losses                uint8
	Q2Losses                uint8
	Q3Losses                uint8
	Q4Losses                uint8
	QuadrantRating          float32
	ConferenceStrengthAdj   float32
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

func (cs *CollegeStandings) GetWinPercentage() float32 {
	totalGames := cs.TotalWins + cs.TotalLosses
	if totalGames == 0 {
		return 0.0
	}
	adjustedWins := float32(cs.TotalWins) * 0.5
	return adjustedWins / float32(totalGames)
}

// GetRPIDisplay returns RPI as a decimal for display purposes
func (cs *CollegeStandings) GetRPIDisplay() string {
	return fmt.Sprintf("%.3f", cs.RPIRating)
}

// GetSOSDisplay returns SOS as a decimal for display purposes
func (cs *CollegeStandings) GetSOSDisplay() string {
	return fmt.Sprintf("%.3f", cs.SOS)
}

// GetSORDisplay returns SOR as a decimal for display purposes
func (cs *CollegeStandings) GetSORDisplay() string {
	return fmt.Sprintf("%.3f", cs.SOR)
}

// GetQualityRecord returns a formatted string showing quality wins
// Q1W-Q1L, Q2W-Q2L, Q3W-Q3L, Q4W-Q4L
func (cs *CollegeStandings) GetQualityRecord() string {
	return fmt.Sprintf("Q1: %d-%d, Q2: %d-%d, Q3: %d-%d, Q4: %d-%d", cs.Q1Wins, cs.Q1Losses, cs.Q2Wins, cs.Q2Losses, cs.Q3Wins, cs.Q3Losses, cs.Q4Wins, cs.Q4Losses)
}

// IsRanked returns true if the team is in the top 25 Toucan rankings
func (cs *CollegeStandings) IsRanked() bool {
	return cs.ToucanRank <= 25 && cs.ToucanRank > 0
}
