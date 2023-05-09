package managers

import (
	"fmt"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
)

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

func GetMatchResultsByMatchID(matchId string) structs.MatchResultsResponse {
	match := GetMatchByMatchId(matchId)

	homePlayers := GetCollegePlayersWithMatchStatsByTeamId(strconv.Itoa(int(match.HomeTeamID)), matchId)
	awayPlayers := GetCollegePlayersWithMatchStatsByTeamId(strconv.Itoa(int(match.AwayTeamID)), matchId)
	homeStats := GetCBBTeamStatsByMatch(strconv.Itoa(int(match.HomeTeamID)), matchId)
	awayStats := GetCBBTeamStatsByMatch(strconv.Itoa(int(match.AwayTeamID)), matchId)

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
