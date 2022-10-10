package managers

import (
	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
)

func GetMatchesByTeamIdAndSeasonId(teamId string, seasonId string) []structs.Match {
	db := dbprovider.GetInstance().GetDB()

	var teamMatches []structs.Match

	db.Where("team_id = ? AND season_id = ?", teamId, seasonId).Find(teamMatches)

	return teamMatches
}

func GetMatchByMatchId(matchId string) structs.Match {
	db := dbprovider.GetInstance().GetDB()

	var match structs.Match

	db.Where("id = ?", matchId).Find(match)

	return match
}

func GetMatchesByWeekId(weekId string) []structs.Match {
	db := dbprovider.GetInstance().GetDB()

	var teamMatches []structs.Match

	db.Where("week_id = ?", weekId).Find(teamMatches)

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
