package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimNBA/managers"
	"github.com/gorilla/mux"
)

func GetMatchesByTeamIdAndSeasonId(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	teamId := vars["teamId"]
	seasonId := vars["seasonId"]
	if len(teamId) == 0 || len(seasonId) == 0 {
		panic("User did not provide both a teamId and a Season Id")
	}

	teamMatches := managers.GetMatchesByTeamIdAndSeasonId(teamId, seasonId)

	json.NewEncoder(w).Encode(teamMatches)
}

func GetNBAMatchesByTeamIdAndSeasonId(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	teamId := vars["teamId"]
	seasonId := vars["seasonId"]
	if len(teamId) == 0 || len(seasonId) == 0 {
		panic("User did not provide both a teamId and a Season Id")
	}

	teamMatches := managers.GetProfessionalMatchesByTeamIdAndSeasonId(teamId, seasonId)

	json.NewEncoder(w).Encode(teamMatches)
}

func GetMatchesBySeasonID(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	seasonId := vars["seasonID"]
	if len(seasonId) == 0 {
		panic("User did not provide both a teamId and a Season Id")
	}

	matchesObj := managers.GetSchedulePageData(seasonId)

	json.NewEncoder(w).Encode(matchesObj)
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

func GetMatchResultByMatchID(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	matchId := vars["matchId"]
	if len(matchId) == 0 {
		panic("User did not provide a matchId")
	}

	match := managers.GetMatchResultsByMatchID(matchId)

	json.NewEncoder(w).Encode(match)
}

func GetNBAMatchResultByMatchID(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	matchId := vars["matchId"]
	if len(matchId) == 0 {
		panic("User did not provide a matchId")
	}

	match := managers.GetNBAMatchResultsByMatchID(matchId)

	json.NewEncoder(w).Encode(match)
}

func GetMatchesByWeekId(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	weekId := vars["weekId"]
	if len(weekId) == 0 {
		panic("User did not provide both a teamId and a Season Id")
	}

	// teamMatches := managers.GetMatchesByWeekId(weekId)

	// json.NewEncoder(w).Encode(teamMatches)
}

func GetUpcomingMatchesByTeamIdAndSeasonId(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)

	teamId := vars["teamId"]
	seasonId := vars["seasonId"]
	if len(teamId) == 0 || len(seasonId) == 0 {
		panic("User did not provide both a teamId and a Season Id")
	}

	upcomingMatches := managers.GetUpcomingMatchesByTeamIdAndSeasonId(teamId, seasonId)

	json.NewEncoder(w).Encode(upcomingMatches)
}

func FixPlayerStatsFromLastSeason(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)

	managers.FixPlayerStatsFromLastSeason()

	json.NewEncoder(w).Encode("All done!")
}

func GetMatchesForSimulation(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)

	matches := managers.GetMatchesForTimeslot()

	json.NewEncoder(w).Encode(matches)
}

// Export Stats
func ExportMatchResults(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)

	seasonID := vars["seasonID"]
	weekID := vars["weekID"]
	nbaWeekID := vars["nbaWeekID"]
	matchType := vars["matchType"]
	w.Header().Set("Content-Type", "text/csv")
	managers.ExportMatchResults(w, seasonID, weekID, nbaWeekID, matchType)
}

// Export Stats
func SwapNBATeamsTEMP(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	managers.SwapNBATeamsTEMP()
}

func FixISLMatchData(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)

	managers.FixISLMatchData()

	json.NewEncoder(w).Encode("Done.")
}
