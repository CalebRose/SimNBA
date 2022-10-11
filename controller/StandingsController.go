package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimNBA/managers"
	"github.com/gorilla/mux"
)

func GetConferenceStandingsByConferenceID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	conferenceID := vars["conferenceId"]
	seasonID := vars["seasonId"]
	if len(conferenceID) == 0 {
		panic("User did not provide TeamID")
	}

	conferenceStandings := managers.GetConferenceStandingsByConferenceID(conferenceID, seasonID)

	json.NewEncoder(w).Encode(conferenceStandings)
}