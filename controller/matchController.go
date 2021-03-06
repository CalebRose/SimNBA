package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimNBA/managers"
	"github.com/gorilla/mux"
)

func GetMatchesByTeamIdAndSeasonId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamId := vars["teamId"]
	seasonId := vars["seasonId"]
	if len(teamId) == 0 || len(seasonId) == 0 {
		panic("User did not provide both a teamId and a Season Id")
	}

	teamMatches := managers.GetMatchesByTeamIdAndSeasonId(teamId, seasonId)

	json.NewEncoder(w).Encode(teamMatches)
}

func GetMatchByMatchId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	matchId := vars["matchId"]
	if len(matchId) == 0 {
		panic("User did not provide a matchId")
	}

	match := managers.GetMatchByMatchId(matchId)

	json.NewEncoder(w).Encode(match)
}

func GetMatchesByWeekId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	weekId := vars["weekId"]
	if len(weekId) == 0 {
		panic("User did not provide both a teamId and a Season Id")
	}

	teamMatches := managers.GetMatchesByWeekId(weekId)

	json.NewEncoder(w).Encode(teamMatches)
}

func GetUpcomingMatchesByTeamIdAndSeasonId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	teamId := vars["teamId"]
	seasonId := vars["seasonId"]
	if len(teamId) == 0 || len(seasonId) == 0 {
		panic("User did not provide both a teamId and a Season Id")
	}

	upcomingMatches := managers.GetUpcomingMatchesByTeamIdAndSeasonId(teamId, seasonId)

	json.NewEncoder(w).Encode(upcomingMatches)
}
