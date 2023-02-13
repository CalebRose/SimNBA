package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimNBA/managers"
	"github.com/gorilla/mux"
)

// GetStatsPageData - Get Stats Page Data
func GetStatsPageData(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	response := managers.GetStatsPageData()

	json.NewEncoder(w).Encode(response)
}

// GetPlayerStatsBySeason - Get Stats By PlayerId and SeasonId
func GetPlayerStats(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	playerId := vars["playerId"]
	if len(playerId) == 0 {
		panic("User did not provide both a playerId and a Season Id")
	}

	playerStats := managers.GetPlayerStatsByPlayerId(playerId)

	json.NewEncoder(w).Encode(playerStats)
}

// GetPlayerStatsBySeason - Get Stats By PlayerId and SeasonId
func GetPlayerStatsBySeason(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)

	playerId := vars["playerId"]
	seasonId := vars["seasonId"]
	if len(playerId) == 0 || len(seasonId) == 0 {
		panic("User did not provide both a playerId and a Season Id")
	}

	playerStats := managers.GetPlayerStatsBySeason(playerId, seasonId)

	json.NewEncoder(w).Encode(playerStats)
}

// GetPlayerStatsBySeason - Get Stats By PlayerId and SeasonId
func GetPlayerStatsInConferenceBySeason(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)

	seasonId := vars["seasonId"]
	conference := vars["conference"]
	if len(seasonId) == 0 || len(conference) == 0 {
		panic("User did not provide both a playerId and a Season Id")
	}

	playerStats := managers.GetPlayerStatsInConferenceBySeason(seasonId, conference)

	json.NewEncoder(w).Encode(playerStats)
}

// GetPlayerStatsByMatch - Get Player Stats by Match played | NOTE: Will revise this func for later
func GetPlayerStatsByMatch(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	playerId := vars["playerId"]
	matchId := vars["matchId"]
	if len(playerId) == 0 || len(matchId) == 0 {
		panic("User did not provide both a playerId and a Match Id")
	}

	playerStats := managers.GetPlayerStatsByMatch(matchId)

	json.NewEncoder(w).Encode(playerStats)
}

// GetTeamStatsBySeason - Get Stats By PlayerId and SeasonId
func GetTeamStatsBySeason(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	teamId := vars["teamId"]
	seasonId := vars["seasonId"]
	if len(teamId) == 0 || len(seasonId) == 0 {
		panic("User did not provide both a teamId and a Season Id")
	}

	playerStats := managers.GetTeamStatsBySeason(teamId, seasonId)

	json.NewEncoder(w).Encode(playerStats)
}

// GetTeamStatsByMatch - Get Player Stats by Match played
func GetTeamStatsByMatch(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)

	teamId := vars["teamId"]
	matchId := vars["matchId"]
	if len(teamId) == 0 || len(matchId) == 0 {
		panic("User did not provide both a teamId and a Match Id")
	}

	playerStats := managers.GetTeamStatsByMatch(teamId, matchId)

	json.NewEncoder(w).Encode(playerStats)
}

// GetTeamStatsByWeek | To be written
