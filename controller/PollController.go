package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/CalebRose/SimNBA/managers"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/gorilla/mux"
)

func CreatePollSubmission(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	var dto structs.CollegePollSubmission
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// validate info from DTO
	if len(dto.Username) == 0 {
		log.Fatalln("ERROR: Cannot submit poll.")
	}

	poll := managers.CreatePoll(dto)
	json.NewEncoder(w).Encode(poll)
}

func GetPollSubmission(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	username := vars["username"]

	poll := managers.GetPollSubmissionByUsernameWeekAndSeason(username)
	json.NewEncoder(w).Encode(poll)
}

func SyncCollegePoll(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	managers.SyncCollegePollSubmissionForCurrentWeek()
}

func GetOfficialPollByWeekIDAndSeasonID(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	weekID := vars["weekID"]
	seasonID := vars["seasonID"]
	if len(weekID) == 0 {
		panic("User did not provide teamID")
	}
	poll := managers.GetOfficialPollByWeekIDAndSeasonID(weekID, seasonID)

	json.NewEncoder(w).Encode(poll)
}
