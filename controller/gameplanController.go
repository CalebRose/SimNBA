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

	vars := mux.Vars(r)

	teamId := vars["teamId"]

	if len(teamId) == 0 {
		panic("User did not provide TeamID")
	}

	var gameplans = managers.GetGameplansByTeam(teamId)

	json.NewEncoder(w).Encode(gameplans)
}

func UpdateGameplan(w http.ResponseWriter, r *http.Request) {

	var updateGameplanDto structs.UpdateGameplanDto

	err := json.NewDecoder(r.Body).Decode(&updateGameplanDto)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	managers.UpdateGameplan(updateGameplanDto)

	fmt.Println("Updated Gameplans and Players")
	w.WriteHeader(http.StatusOK)
}
