package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/managers"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/gorilla/mux"
)

func AllTeams(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	db := dbprovider.GetInstance().GetDB()

	var teams []structs.Team
	db.Order("team asc").Find(&teams)
	json.NewEncoder(w).Encode(teams)
}

func AllActiveTeams(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	db := dbprovider.GetInstance().GetDB()

	var teams []structs.Team
	db.Where("first_season is not null").Order("team asc").Find(&teams)
	json.NewEncoder(w).Encode(teams)
}

func AllActiveCollegeTeams(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	db := dbprovider.GetInstance().GetDB()

	var teams []structs.Team
	db.Where("first_season is not null AND coach is not null and is_nba = ?", false).Order("team asc").Find(&teams)
	json.NewEncoder(w).Encode(teams)
}

func AllAvailableTeams(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	db := dbprovider.GetInstance().GetDB()

	var teams []structs.Team
	db.Where("(coach = ? OR coach = ?) AND is_active = true", "AI", "").Order("team asc").Find(&teams)

	json.NewEncoder(w).Encode(teams)
}

func AllCoachedTeams(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	db := dbprovider.GetInstance().GetDB()

	var teams []structs.Team
	db.Where("coach is not null AND coach NOT IN (?,?)", "", "AI").Order("team asc").Find(&teams)
	json.NewEncoder(w).Encode(teams)
}

func AllCollegeTeams(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	db := dbprovider.GetInstance().GetDB()

	var teams []structs.Team
	db.Where("is_nba = ?, is_active = ?", false, true).Order("team asc").Find(&teams)
	json.NewEncoder(w).Encode(teams)
}

func AllNBATeams(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	db := dbprovider.GetInstance().GetDB()

	var teams []structs.NBATeam
	db.Order("team asc").Where("league_id = ?", "1").Find(&teams)
	json.NewEncoder(w).Encode(teams)
}

func AllProfessionalTeams(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	db := dbprovider.GetInstance().GetDB()

	var teams []structs.NBATeam
	db.Order("team asc").Find(&teams)
	json.NewEncoder(w).Encode(teams)
}

func GetTeamByTeamID(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	teamId := vars["teamId"]
	if len(teamId) == 0 {
		panic("User did not provide TeamID")
	}
	team := managers.GetTeamByTeamID(teamId)
	json.NewEncoder(w).Encode(team)
}

func GetNBATeamByTeamID(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	teamId := vars["teamId"]
	if len(teamId) == 0 {
		panic("User did not provide TeamID")
	}
	team := managers.GetNBATeamByTeamID(teamId)
	json.NewEncoder(w).Encode(team)
}

func RemoveUserFromTeam(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	teamId := vars["teamId"]
	if len(teamId) == 0 {
		panic("User did not provide TeamID")
	}
	team := managers.RemoveUserFromTeam(teamId)
	json.NewEncoder(w).Encode(team)
}

func SyncTeamRatings(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	db := dbprovider.GetInstance().GetDB()

	var teams []structs.Team
	db.Where("is_nba = ? AND is_active = ?", false, true).Order("team asc").Find(&teams)

	for _, team := range teams {
		managers.GetTeamRatings(team)
	}

	json.NewEncoder(w).Encode("Team Ratings Sync Done!")
}

func SyncNBATeamRatings(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	db := dbprovider.GetInstance().GetDB()

	var teams []structs.NBATeam
	db.Order("team asc").Find(&teams)

	for _, team := range teams {
		managers.GetNBATeamRatings(team)
	}

	json.NewEncoder(w).Encode("Team Ratings Sync Done!")
}
