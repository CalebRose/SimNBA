package controller

import (
	"net/http"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/managers"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/gorilla/mux"
)

func GetMatchesByTeamIdAndSeasonId(w http.ResponseWriter, r *http.Request) {
	db := dbprovider.GetInstance().GetDB()

	vars := mux.Vars(r)
	teamId := vars["teamId"]
	seasonId := vars["seasonId"]
	if len(teamId) == 0 || len(seasonId) == 0 {
		panic("User did not provide both a teamId and a Season Id")
	}

	var teamMatches []structs.Match

	db.Where("team_id = ? AND season_id = ?", teamId, seasonId).Find(teamMatches)
}

func GetMatchByMatchId(w http.ResponseWriter, r *http.Request) {
	db := dbprovider.GetInstance().GetDB()

	vars := mux.Vars(r)
	matchId := vars["matchId"]
	if len(matchId) == 0 {
		panic("User did not provide a matchId")
	}

	var match structs.Match

	db.Where("id = ?", matchId).Find(match)
}

func GetMatchesByWeekId(w http.ResponseWriter, r *http.Request) {
	db := dbprovider.GetInstance().GetDB()

	vars := mux.Vars(r)
	weekId := vars["weekId"]
	if len(weekId) == 0 {
		panic("User did not provide both a teamId and a Season Id")
	}

	var teamMatches []structs.Match

	db.Where("week_id = ?", weekId).Find(teamMatches)
}

func GetUpcomingMatchesByTeamIdAndSeasonId(w http.ResponseWriter, r *http.Request) {
	db := dbprovider.GetInstance().GetDB()

	vars := mux.Vars(r)
	teamId := vars["teamId"]
	seasonId := vars["seasonId"]
	if len(teamId) == 0 || len(seasonId) == 0 {
		panic("User did not provide both a teamId and a Season Id")
	}

	timeStamp := managers.GetTimestamp(db)

	var teamMatches []structs.Match

	db.Where("team_id = ? AND season_id = ? AND week_id > ?", teamId, seasonId, timeStamp.CollegeWeekID).Find(teamMatches)
}
