package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/managers"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/gorilla/mux"
)

// GameplanController - For routes on Gameplans
func GetGameplansByTeamId(w http.ResponseWriter, r *http.Request) {
	db := dbprovider.GetInstance().GetDB()
	vars := mux.Vars(r)
	teamId := vars["teamId"]
	if len(teamId) == 0 {
		panic("User did not provide TeamID")
	}

	var gameplans = managers.GetGameplansByTeam(db, teamId)
	json.NewEncoder(w).Encode(gameplans)
}

func UpdateGameplan(w http.ResponseWriter, r *http.Request) {
	db := dbprovider.GetInstance().GetDB()

	var updateGameplanDto structs.UpdateGameplanDto
	err := json.NewDecoder(r.Body).Decode(&updateGameplanDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var teamId = strconv.Itoa(updateGameplanDto.TeamID)
	// Get Gameplans
	var gameplans = managers.GetGameplansByTeam(db, teamId)

	for i := 0; i < len(gameplans); i++ {
		updatedGameplan := updateGameplanDto.Gameplans[i]
		if gameplans[i].Pace == updatedGameplan.Pace &&
			gameplans[i].ThreePointProportion == updatedGameplan.ThreePointProportion &&
			gameplans[i].JumperProportion == updatedGameplan.JumperProportion &&
			gameplans[i].PaintProportion == updatedGameplan.PaintProportion {
			continue
		}
		gameplans[i].UpdatePace(updateGameplanDto.Gameplans[i].Pace)
		gameplans[i].Update3PtProportion(updatedGameplan.ThreePointProportion)
		gameplans[i].UpdateJumperProportion(updatedGameplan.JumperProportion)
		gameplans[i].UpdatePaintProportion(updatedGameplan.PaintProportion)
		fmt.Printf("Saving Gameplan for Team " + teamId + "\n")
		db.Save(&gameplans[i])
	}

	// Get Players
	var players = managers.GetPlayersByTeamId(db, teamId)

	for i := 0; i < len(players); i++ {
		updatedPlayer := updateGameplanDto.Players[i]
		if players[i].MinutesA == updatedPlayer.MinutesA &&
			players[i].MinutesB == updatedPlayer.MinutesB &&
			players[i].MinutesC == updatedPlayer.MinutesC {
			continue
		}
		players[i].UpdateMinutesA(updatedPlayer.MinutesA)
		players[i].UpdateMinutesB(updatedPlayer.MinutesB)

		// If player is an NBA player, update Minutes for C Game
		if players[i].IsNBA == true {
			players[i].UpdateMinutesC(updateGameplanDto.Players[i].MinutesC)
		}
		fmt.Printf("Saving Player " + players[i].FirstName + " " + players[i].LastName + "\n")
		db.Save(&players[i])
	}

	fmt.Println("Updated Gameplans and Players")
	w.WriteHeader(http.StatusOK)
}
