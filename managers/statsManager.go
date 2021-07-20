package managers

import (
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

func GetPlayerStatsByMatch(matchId string) []structs.PlayerStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.PlayerStats

	db.Where("match_id = ?", matchId).Find(playerStats)

	return playerStats
}

func GetTeamStatsBySeason(teamId string, seasonId string) []structs.PlayerStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.PlayerStats

	db.Where("team_id = ? AND season_id = ?", teamId, seasonId).Find(playerStats)

	return playerStats
}

func GetTeamStatsByMatch(teamId string, matchId string) []structs.PlayerStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.PlayerStats

	db.Where("team_id = ? AND match_id = ?", teamId, matchId).Find(playerStats)

	return playerStats
}
