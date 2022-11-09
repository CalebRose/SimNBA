package managers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
)

func GetPlayerStatsByPlayerId(playerId string) []structs.PlayerStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.PlayerStats

	db.Where("player_id = ?", playerId).Find(playerStats)

	return playerStats
}

func GetPlayerStatsBySeason(playerId string, seasonId string) []structs.PlayerStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.PlayerStats

	db.Where("player_id = ? AND season_id = ?", playerId, seasonId).Find(playerStats)

	return playerStats
}

func GetPlayerStatsInConferenceBySeason(seasonId string, conference string) []structs.PlayerStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.PlayerStats

	db.Where("season_id = ? AND conference = ?", seasonId, conference).Find(playerStats)

	return playerStats
}

func GetPlayerStatsByMatch(matchId string) []structs.CollegePlayerStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.CollegePlayerStats

	db.Where("match_id = ?", matchId).Find(playerStats)

	return playerStats
}

func GetTeamStatsBySeason(teamId string, seasonId string) []structs.PlayerStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.PlayerStats

	db.Where("team_id = ? AND season_id = ?", teamId, seasonId).Find(playerStats)

	return playerStats
}

func GetTeamStatsByMatch(teamId string, matchId string) structs.TeamStats {
	db := dbprovider.GetInstance().GetDB()

	var teamStats structs.TeamStats

	db.Where("team_id = ? AND match_id = ?", teamId, matchId).Find(&teamStats)

	return teamStats
}

func GetPlayerSeasonStatsByPlayerID(playerID string, seasonID string) structs.CollegePlayerSeasonStats {
	db := dbprovider.GetInstance().GetDB()

	var seasonStats structs.CollegePlayerSeasonStats

	err := db.Where("college_player_id = ? AND season_id = ?", playerID, seasonID).Find(&seasonStats)
	if err != nil {
		fmt.Println("Could not find existing record for player... generating new one.")
	}

	return seasonStats
}

func GetTeamSeasonStatsByTeamID(teamID string, seasonID string) structs.TeamSeasonStats {
	db := dbprovider.GetInstance().GetDB()

	var seasonStats structs.TeamSeasonStats

	err := db.Where("team_id = ? AND season_id = ?", teamID, seasonID).Find(&seasonStats)
	if err != nil {
		fmt.Println("Could not find existing record for team... generating new one.")
	}

	return seasonStats
}

func UpdateSeasonStats(ts structs.Timestamp) {
	db := dbprovider.GetInstance().GetDB()

	weekId := strconv.Itoa(int(ts.CollegeWeekID))
	seasonId := strconv.Itoa(int(ts.SeasonID))

	matches := GetMatchesByWeekId(weekId, seasonId)

	for _, match := range matches {
		homeTeamStats := GetTeamStatsByMatch(strconv.Itoa(int(match.HomeTeamID)), strconv.Itoa(int(match.ID)))

		homeSeasonStats := GetTeamSeasonStatsByTeamID(strconv.Itoa(int(match.HomeTeamID)), seasonId)

		homeSeasonStats.AddStatsToSeasonRecord(homeTeamStats)

		err := db.Save(&homeSeasonStats).Error
		if err != nil {
			log.Fatalln("Could not save season stats for " + strconv.Itoa(int(match.HomeTeamID)))
		}

		awayTeamStats := GetTeamStatsByMatch(strconv.Itoa(int(match.AwayTeamID)), strconv.Itoa(int(match.ID)))

		awaySeasonStats := GetTeamSeasonStatsByTeamID(strconv.Itoa(int(match.AwayTeamID)), seasonId)

		awaySeasonStats.AddStatsToSeasonRecord(awayTeamStats)

		err = db.Save(&awaySeasonStats).Error
		if err != nil {
			log.Fatalln("Could not save season stats for " + strconv.Itoa(int(match.AwayTeamID)))
		}

		playerStats := GetPlayerStatsByMatch(strconv.Itoa(int(match.ID)))

		for _, stat := range playerStats {
			playerSeasonStats := GetPlayerSeasonStatsByPlayerID(strconv.Itoa(int(stat.CollegePlayerID)), seasonId)

			playerSeasonStats.AddStatsToSeasonRecord(stat)

			err = db.Save(&playerSeasonStats).Error
			if err != nil {
				log.Fatalln("Could not save season stats for " + strconv.Itoa(int(playerSeasonStats.CollegePlayerID)))
			}
		}
	}
}
