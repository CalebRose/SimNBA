package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/CalebRose/SimNBA/managers"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/gorilla/mux"
)

func AllPlayers(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	var players = managers.GetAllPlayers()

	json.NewEncoder(w).Encode(players)
}

func AllPlayersByTeamId(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)

	teamId := vars["teamId"]

	if len(teamId) == 0 {
		panic("User did not provide TeamID")
	}

	var players = managers.GetTeamRosterForRosterPage(teamId)

	json.NewEncoder(w).Encode(players)
}

func GetNBARosterByTeamID(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)

	teamId := vars["teamId"]

	if len(teamId) == 0 {
		panic("User did not provide TeamID")
	}

	var players = managers.GetAllNBAPlayersByTeamID(teamId)

	json.NewEncoder(w).Encode(players)
}

func AllCollegePlayers(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	var players = managers.GetAllCollegePlayers()

	json.NewEncoder(w).Encode(players)
}

func AllCollegeRecruits(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	var recruits = managers.GetAllCollegeRecruits()

	json.NewEncoder(w).Encode(recruits)
}

func AllJUCOCollegeRecruits(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	var recruits = managers.GetAllJUCOCollegeRecruits()

	json.NewEncoder(w).Encode(recruits)
}

func AllNBAPlayers(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	var players = managers.GetAllNBAPlayers()

	json.NewEncoder(w).Encode(players)
}

func AllNBAFreeAgents(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	players := managers.GetAllNBAPlayersByTeamID("0")

	json.NewEncoder(w).Encode(players)
}

func PlayerById(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)

	playerId := vars["playerId"]
	if len(playerId) == 0 {
		panic("User did not provide PlayerID")
	}

	player := managers.GetCollegePlayerByPlayerId(playerId)
	json.NewEncoder(w).Encode(player)
}

func AssignRedshirtForCollegePlayer(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	var redshirtDTO structs.RedshirtDTO
	err := json.NewDecoder(r.Body).Decode(&redshirtDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var player = managers.SetRedshirtStatusForPlayer(strconv.Itoa(redshirtDTO.PlayerID))

	json.NewEncoder(w).Encode(player)
}

// Old Method -- Not for official use
func NewPlayer(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)

	firstName := vars["firstname"]
	lastName := vars["lastname"]
	if len(firstName) == 0 || len(lastName) == 0 {
		log.Fatal("Need a first name and last name")
	}

	managers.CreateNewPlayer(firstName, lastName)

	fmt.Fprintf(w, "New Player Successfully Created")
}
