package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimNBA/managers"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/gorilla/mux"
)

// FreeAgencyAvailablePlayers - Get All Available NFL Players for Free Agency Page
func FreeAgencyAvailablePlayers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamId := vars["teamID"]
	var players = managers.GetAllAvailableNBAPlayers(teamId)

	json.NewEncoder(w).Encode(players)
}

// FreeAgencyAvailablePlayers - Get All Available NFL Players for Free Agency Page
func CreateFreeAgencyOffer(w http.ResponseWriter, r *http.Request) {
	var freeAgencyOfferDTO structs.NBAContractOfferDTO
	err := json.NewDecoder(r.Body).Decode(&freeAgencyOfferDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var offer = managers.CreateFAOffer(freeAgencyOfferDTO)

	json.NewEncoder(w).Encode(offer)
}

// FreeAgencyAvailablePlayers - Get All Available NFL Players for Free Agency Page
func CancelFreeAgencyOffer(w http.ResponseWriter, r *http.Request) {
	var freeAgencyOfferDTO structs.NBAContractOfferDTO
	err := json.NewDecoder(r.Body).Decode(&freeAgencyOfferDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.CancelOffer(freeAgencyOfferDTO)

	json.NewEncoder(w).Encode(true)
}
