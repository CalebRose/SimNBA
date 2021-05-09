package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CalebRose/SimNBA/managers"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func AllTeams(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(c["db"], c["cs"])
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}

	defer db.Close()

	var teams []structs.Team
	db.Order("team asc").Find(&teams)
	json.NewEncoder(w).Encode(teams)
}

func AllActiveTeams(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(c["db"], c["cs"])
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}

	defer db.Close()

	var teams []structs.Team
	db.Where("first_season is not null").Order("team asc").Find(&teams)
	json.NewEncoder(w).Encode(teams)
}

func AllAvailableTeams(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(c["db"], c["cs"])
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}

	defer db.Close()

	var teams []structs.Team
	db.Where("first_season is not null AND coach is null OR coach = ?", "AI").Order("team asc").Find(&teams)
	json.NewEncoder(w).Encode(teams)
}

func AllCoachedTeams(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(c["db"], c["cs"])
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}

	defer db.Close()

	var teams []structs.Team
	db.Where("coach is not null AND coach NOT IN (?,?)", "", "AI").Order("team asc").Find(&teams)
	json.NewEncoder(w).Encode(teams)
}

func AllCollegeTeams(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(c["db"], c["cs"])
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}

	defer db.Close()

	var teams []structs.Team
	db.Where("is_nba = ?", false).Order("team asc").Find(&teams)
	json.NewEncoder(w).Encode(teams)
}

func AllNBATeams(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(c["db"], c["cs"])
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}

	defer db.Close()

	var teams []structs.Team
	db.Where("is_nba = ?", true).Order("team asc").Find(&teams)
	json.NewEncoder(w).Encode(teams)
}

func GetTeamByTeamID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamId := vars["teamId"]
	if len(teamId) == 0 {
		panic("User did not provide TeamID")
	}
	team := managers.GetTeamByTeamID(teamId)
	json.NewEncoder(w).Encode(team)
}

func RemoveUserFromTeam(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamId := vars["teamId"]
	if len(teamId) == 0 {
		panic("User did not provide TeamID")
	}
	team := managers.GetTeamByTeamID(teamId)
	team.RemoveUser()
	managers.RemoveUserFromTeam(team)
	json.NewEncoder(w).Encode(team)
}
