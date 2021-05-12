package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CalebRose/SimNBA/structs"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// GetPlayerStatsBySeason - Get Stats By PlayerId and SeasonId
func GetPlayerStats(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(c["db"], c["cs"])
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}

	defer db.Close()

	vars := mux.Vars(r)
	playerId := vars["playerId"]
	if len(playerId) == 0 {
		panic("User did not provide both a playerId and a Season Id")
	}

	var playerStats []structs.PlayerStats

	db.Where("player_id = ?", playerId).Find(playerStats)
	json.NewEncoder(w).Encode(playerStats)
}

// GetPlayerStatsBySeason - Get Stats By PlayerId and SeasonId
func GetPlayerStatsBySeason(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(c["db"], c["cs"])
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}

	defer db.Close()

	vars := mux.Vars(r)
	playerId := vars["playerId"]
	seasonId := vars["seasonId"]
	if len(playerId) == 0 || len(seasonId) == 0 {
		panic("User did not provide both a playerId and a Season Id")
	}

	var playerStats []structs.PlayerStats

	db.Where("player_id = ? AND season_id = ?", playerId, seasonId).Find(playerStats)
	json.NewEncoder(w).Encode(playerStats)
}

// GetPlayerStatsBySeason - Get Stats By PlayerId and SeasonId
func GetPlayerStatsInConferenceBySeason(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(c["db"], c["cs"])
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}

	defer db.Close()

	vars := mux.Vars(r)
	seasonId := vars["seasonId"]
	conference := vars["conference"]
	if len(seasonId) == 0 || len(conference) == 0 {
		panic("User did not provide both a playerId and a Season Id")
	}

	var playerStats []structs.PlayerStats

	// Get Teams, preload players, 

	db.Where("season_id = ? AND conference = ?", seasonId, conference).Find(playerStats)
	json.NewEncoder(w).Encode(playerStats)
}

// GetPlayerStatsByMatch - Get Player Stats by Match played
func GetPlayerStatsByMatch(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(c["db"], c["cs"])
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}

	defer db.Close()

	vars := mux.Vars(r)
	playerId := vars["playerId"]
	matchId := vars["matchId"]
	if len(playerId) == 0 || len(matchId) == 0 {
		panic("User did not provide both a playerId and a Match Id")
	}

	var playerStats []structs.PlayerStats

	db.Where("player_id = ? AND match_id = ?", playerId, matchId).Find(playerStats)
	json.NewEncoder(w).Encode(playerStats)
}

// GetTeamStatsBySeason - Get Stats By PlayerId and SeasonId
func GetTeamStatsBySeason(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(c["db"], c["cs"])
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}

	defer db.Close()

	vars := mux.Vars(r)
	teamId := vars["teamId"]
	seasonId := vars["seasonId"]
	if len(teamId) == 0 || len(seasonId) == 0 {
		panic("User did not provide both a teamId and a Season Id")
	}

	var playerStats []structs.PlayerStats

	db.Where("team_id = ? AND season_id = ?", teamId, seasonId).Find(playerStats)
	json.NewEncoder(w).Encode(playerStats)
}

// GetTeamStatsByMatch - Get Player Stats by Match played
func GetTeamStatsByMatch(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(c["db"], c["cs"])
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}

	defer db.Close()

	vars := mux.Vars(r)
	teamId := vars["teamId"]
	matchId := vars["matchId"]
	if len(teamId) == 0 || len(matchId) == 0 {
		panic("User did not provide both a teamId and a Match Id")
	}

	var playerStats []structs.PlayerStats

	db.Where("team_id = ? AND match_id = ?", teamId, matchId).Find(playerStats)
	json.NewEncoder(w).Encode(playerStats)
}

// GetTeamStatsByWeek
