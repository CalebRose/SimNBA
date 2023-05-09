package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CalebRose/SimNBA/managers"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/gorilla/mux"
)

// GameplanController - For routes on Gameplans
func GetGameplansByTeamId(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)

	teamId := vars["teamId"]

	if len(teamId) == 0 {
		panic("User did not provide TeamID")
	}

	gp := managers.GetGameplansByTeam(teamId)
	opposingRoster := managers.GetOpposingCollegiateTeamRoster(teamId)

	res := structs.GameplanResponse{
		Gameplan:       gp,
		OpposingRoster: opposingRoster,
	}

	json.NewEncoder(w).Encode(res)
}

func UpdateGameplan(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	var updateGameplanDto structs.UpdateGameplanDto

	err := json.NewDecoder(r.Body).Decode(&updateGameplanDto)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	managers.UpdateGameplan(updateGameplanDto)

	fmt.Println("Updated Gameplans and Players")
	w.WriteHeader(http.StatusOK)
}

func GetNBAGameplanByTeamId(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)

	teamId := vars["teamId"]

	if len(teamId) == 0 {
		panic("User did not provide TeamID")
	}

	gp := managers.GetNBAGameplanByTeam(teamId)
	opposingRoster := managers.GetOpposingNBATeamRoster(teamId)

	res := structs.NBAGameplanResponse{
		Gameplan:       gp,
		OpposingRoster: opposingRoster,
	}

	json.NewEncoder(w).Encode(res)
}

func UpdateNBAGameplan(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	var updateGameplanDto structs.UpdateGameplanDto

	err := json.NewDecoder(r.Body).Decode(&updateGameplanDto)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	managers.UpdateNBAGameplan(updateGameplanDto)

	fmt.Println("Updated Gameplans and Players")
	w.WriteHeader(http.StatusOK)
}

func SetAIGameplans(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	ping := managers.SetAIGameplans()
	if ping {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("AI Gameplans Set.")
	} else {
		w.WriteHeader(http.StatusExpectationFailed)
		json.NewEncoder(w).Encode("AI Gameplans failed to set")
	}
}
