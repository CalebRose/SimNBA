package managers

import (
	"fmt"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
)

func GetMatchesForTimeslot() []structs.MatchResponse {
	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.SeasonID))
	weekID := strconv.Itoa(int(ts.CollegeWeekID))
	nbaWeekID := strconv.Itoa(int(ts.NBAWeekID))

	matchType := ""

	if !ts.GamesARan {
		matchType = "A"
	} else if !ts.GamesBRan {
		matchType = "B"
	} else if !ts.GamesCRan {
		matchType = "C"
	} else if !ts.GamesDRan {
		matchType = "D"
	} else {
		panic("STOP")
	}

	matchesList := []structs.MatchResponse{}

	// Get College Matches
	collegeMatches := GetMatchesByWeekId(weekID, seasonID, matchType)

	for _, c := range collegeMatches {
		if c.GameComplete {
			continue
		}

		match := structs.MatchResponse{
			MatchName:              c.MatchName,
			ID:                     c.ID,
			WeekID:                 c.WeekID,
			Week:                   c.Week,
			SeasonID:               c.SeasonID,
			HomeTeamID:             c.HomeTeamID,
			HomeTeam:               c.HomeTeam,
			AwayTeamID:             c.AwayTeamID,
			AwayTeam:               c.AwayTeam,
			MatchOfWeek:            c.MatchOfWeek,
			Arena:                  c.Arena,
			City:                   c.City,
			State:                  c.State,
			IsNeutralSite:          c.IsNeutralSite,
			IsNBAMatch:             false,
			IsConference:           c.IsConference,
			IsConferenceTournament: c.IsConferenceTournament,
			IsNILGame:              c.IsNILGame,
			IsPlayoffGame:          c.IsPlayoffGame,
			IsNationalChampionship: c.IsNationalChampionship,
			IsRivalryGame:          c.IsRivalryGame,
		}

		matchesList = append(matchesList, match)
	}

	// Get Professional Matches
	nbaMatches := GetNBATeamMatchesByMatchType(nbaWeekID, seasonID, matchType)

	for _, n := range nbaMatches {
		if n.GameComplete {
			continue
		}

		match := structs.MatchResponse{
			MatchName:              n.MatchName,
			ID:                     n.ID,
			WeekID:                 n.WeekID,
			Week:                   n.Week,
			SeasonID:               n.SeasonID,
			HomeTeamID:             n.HomeTeamID,
			HomeTeam:               n.HomeTeam,
			AwayTeamID:             n.AwayTeamID,
			AwayTeam:               n.AwayTeam,
			MatchOfWeek:            n.MatchOfWeek,
			Arena:                  n.Arena,
			City:                   n.City,
			State:                  n.State,
			IsNeutralSite:          n.IsNeutralSite,
			IsNBAMatch:             true,
			IsConference:           n.IsConference,
			IsConferenceTournament: n.IsConferenceTournament,
			IsNILGame:              n.IsNILGame,
			IsPlayoffGame:          n.IsPlayoffGame,
			IsNationalChampionship: n.IsTheFinals,
			IsRivalryGame:          n.IsRivalryGame,
		}

		matchesList = append(matchesList, match)
	}

	return matchesList
}

func GetMatchesByTeamIdAndSeasonId(teamId string, seasonId string) []structs.Match {
	db := dbprovider.GetInstance().GetDB()

	var teamMatches []structs.Match

	db.Where("(home_team_id = ? OR away_team_id = ?) AND season_id = ?", teamId, teamId, seasonId).Find(&teamMatches)

	return teamMatches
}

func GetMatchesBySeasonID(seasonId string) []structs.Match {
	db := dbprovider.GetInstance().GetDB()

	var teamMatches []structs.Match

	db.Where("season_id = ?", seasonId).Find(&teamMatches)

	return teamMatches
}

func GetMatchByMatchId(matchId string) structs.Match {
	db := dbprovider.GetInstance().GetDB()

	var match structs.Match

	err := db.Where("id = ?", matchId).Find(&match).Error
	if err != nil {
		fmt.Println(err.Error())
	}

	return match
}

func GetNBAMatchByMatchId(matchId string) structs.NBAMatch {
	db := dbprovider.GetInstance().GetDB()

	var match structs.NBAMatch

	err := db.Where("id = ?", matchId).Find(&match).Error
	if err != nil {
		fmt.Println(err.Error())
	}

	return match
}

func GetMatchResultsByMatchID(matchId string) structs.MatchResultsResponse {
	match := GetMatchByMatchId(matchId)

	homePlayers := GetCollegePlayersWithMatchStatsByTeamId(strconv.Itoa(int(match.HomeTeamID)), matchId)
	awayPlayers := GetCollegePlayersWithMatchStatsByTeamId(strconv.Itoa(int(match.AwayTeamID)), matchId)
	homeStats := GetCBBTeamResultsByMatch(strconv.Itoa(int(match.HomeTeamID)), matchId)
	awayStats := GetCBBTeamResultsByMatch(strconv.Itoa(int(match.AwayTeamID)), matchId)

	return structs.MatchResultsResponse{
		HomePlayers: homePlayers,
		AwayPlayers: awayPlayers,
		HomeStats:   homeStats,
		AwayStats:   awayStats,
	}
}

func GetNBAMatchResultsByMatchID(matchId string) structs.MatchResultsResponse {
	match := GetNBAMatchByMatchId(matchId)

	homePlayers := GetNBAPlayersWithMatchStatsByTeamId(strconv.Itoa(int(match.HomeTeamID)), matchId)
	awayPlayers := GetNBAPlayersWithMatchStatsByTeamId(strconv.Itoa(int(match.AwayTeamID)), matchId)
	homeStats := GetNBATeamResultsByMatch(strconv.Itoa(int(match.HomeTeamID)), matchId)
	awayStats := GetNBATeamResultsByMatch(strconv.Itoa(int(match.AwayTeamID)), matchId)

	return structs.MatchResultsResponse{
		HomePlayers: homePlayers,
		AwayPlayers: awayPlayers,
		HomeStats:   homeStats,
		AwayStats:   awayStats,
	}
}

func GetTeamMatchesByWeekId(weekId, seasonID, matchType, teamID string) []structs.Match {
	db := dbprovider.GetInstance().GetDB()

	var teamMatches []structs.Match

	db.Where("week_id = ? AND season_id = ? AND match_of_week = ? AND (home_team_id = ? OR away_team_id = ?)", weekId, seasonID, matchType, teamID, teamID).Find(&teamMatches)

	return teamMatches
}

func GetNBATeamMatchesByMatchType(weekId, seasonID, matchType string) []structs.NBAMatch {
	db := dbprovider.GetInstance().GetDB()

	var teamMatches []structs.NBAMatch

	db.Where("week_id = ? AND season_id = ? AND match_of_week = ?", weekId, seasonID, matchType).Find(&teamMatches)

	return teamMatches
}

func GetNBATeamMatchesByWeekId(weekId, seasonID, matchType, teamID string) []structs.NBAMatch {
	db := dbprovider.GetInstance().GetDB()

	var teamMatches []structs.NBAMatch

	db.Where("week_id = ? AND season_id = ? AND match_of_week = ? AND (home_team_id = ? OR away_team_id = ?)", weekId, seasonID, matchType, teamID, teamID).Find(&teamMatches)

	return teamMatches
}

func GetMatchesByWeekId(weekId string, seasonID string, matchType string) []structs.Match {
	db := dbprovider.GetInstance().GetDB()

	var teamMatches []structs.Match

	db.Where("week_id = ? AND season_id = ? AND match_of_week = ?", weekId, seasonID, matchType).Find(&teamMatches)

	return teamMatches
}

func GetUpcomingMatchesByTeamIdAndSeasonId(teamId string, seasonId string) []structs.Match {
	db := dbprovider.GetInstance().GetDB()

	timeStamp := GetTimestamp()

	var teamMatches []structs.Match

	db.Where("team_id = ? AND season_id = ? AND week_id > ?", teamId, seasonId, timeStamp.CollegeWeekID).Find(teamMatches)

	return teamMatches
}

// SAVE
