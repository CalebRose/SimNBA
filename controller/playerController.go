package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/managers"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// var connectionString = "simfbaah_davidross10:bestpunterev3r!@tcp(68.66.216.54)/simfbaah_simnba?parseTime=true"

func AllPlayers(w http.ResponseWriter, r *http.Request) {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.Player
	db.Find(&players)
	json.NewEncoder(w).Encode(players)
}

func AllPlayersByTeamId(w http.ResponseWriter, r *http.Request) {
	db := dbprovider.GetInstance().GetDB()

	vars := mux.Vars(r)
	teamId := vars["teamId"]
	if len(teamId) == 0 {
		panic("User did not provide TeamID")
	}
	var players = managers.GetPlayersByTeamId(db, teamId)

	json.NewEncoder(w).Encode(players)
}

func AllCollegePlayers(w http.ResponseWriter, r *http.Request) {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.Player
	db.Where("is_nba = ?", false).Find(&players)
	json.NewEncoder(w).Encode(players)
}

func AllCollegeRecruits(w http.ResponseWriter, r *http.Request) {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.Player
	db.Preload("RecruitingPoints", "total_points_spent > ?", 0, func(db *gorm.DB) *gorm.DB {
		return db.Order("total_points_spent DESC")
	}).Where("is_nba = ? AND team_id = 0", false).Find(&players)
	json.NewEncoder(w).Encode(players)
}

func AllJUCOCollegeRecruits(w http.ResponseWriter, r *http.Request) {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.Player
	db.Where("is_nba = ? AND team_id = 0 AND year > 0", false).Find(&players)
	json.NewEncoder(w).Encode(players)
}

func AllNBAPlayers(w http.ResponseWriter, r *http.Request) {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.Player
	db.Where("is_nba = ?", true).Find(&players)
	json.NewEncoder(w).Encode(players)
}

func AllNBAFreeAgents(w http.ResponseWriter, r *http.Request) {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.Player
	db.Where("is_nba = ? AND team_id is null", true).Find(&players)
	json.NewEncoder(w).Encode(players)
}

func PlayerById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playerId := vars["playerId"]
	if len(playerId) == 0 {
		panic("User did not provide PlayerID")
	}

	player := managers.GetPlayerByPlayerId(playerId)
	json.NewEncoder(w).Encode(player)
}

func SetRedshirtStatusByPlayerId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playerId := vars["playerId"]
	if len(playerId) == 0 {
		panic("User did not provide PlayerID")
	}

	player := managers.GetPlayerByPlayerId(playerId)
	player.SetRedshirtingStatus()
	managers.UpdatePlayer(player)
	json.NewEncoder(w).Encode(player)
}

func NewPlayer(w http.ResponseWriter, r *http.Request) {
	db := dbprovider.GetInstance().GetDB()

	vars := mux.Vars(r)
	firstName := vars["firstname"]
	lastName := vars["lastname"]
	fmt.Println(firstName)
	fmt.Println(lastName)

	db.Create(&structs.Player{FirstName: firstName, LastName: lastName,
		Position: "C", Year: 4, State: "WA", Country: "USA",
		Stars: 3, Height: "7'0", TeamID: 10, Shooting: 14,
		Finishing: 20, Ballwork: 18, Rebounding: 20, Defense: 19,
		PotentialGrade: 20, Stamina: 36, PlaytimeExpectations: 25,
		MinutesA: 35, Overall: 20, IsNBA: false,
		IsRedshirt: false, IsRedshirting: false})

	fmt.Fprintf(w, "New Player Successfully Created")
}
