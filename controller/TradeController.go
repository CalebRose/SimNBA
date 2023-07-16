package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CalebRose/SimNBA/managers"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/gorilla/mux"
)

// Get Trade Block Data for Trade Block Page
func GetNBATradeBlockDataByTeamID(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	teamID := vars["teamID"]
	if len(teamID) == 0 {
		panic("User did not provide TeamID")
	}

	response := managers.GetTradeBlockDataByTeamID(teamID)

	json.NewEncoder(w).Encode(response)
}

// Get Trade Block Data for Trade Block Page
func GetAllAcceptedTrades(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	response := managers.GetAcceptedTradeProposals()
	json.NewEncoder(w).Encode(response)
}

// Get Trade Block Data for Trade Block Page
func GetAllRejectedTrades(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	response := managers.GetRejectedTradeProposals()

	json.NewEncoder(w).Encode(response)
}

// Place player on NBA Trade block
func PlaceNBAPlayerOnTradeBlock(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	playerID := vars["playerID"]
	if len(playerID) == 0 {
		panic("User did not provide playerID")
	}

	managers.PlaceNBAPlayerOnTradeBlock(playerID)

	json.NewEncoder(w).Encode("Player " + playerID + " placed on trade block.")
}

// Update Trade Preferences
func UpdateTradePreferences(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	var tradePreferenceDTO structs.NBATradePreferencesDTO
	err := json.NewDecoder(r.Body).Decode(&tradePreferenceDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.UpdateTradePreferences(tradePreferenceDTO)

	fmt.Fprintf(w, "Trade Preferences Updated")
}

// Create NBA Trade Proposal
func CreateNBATradeProposal(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	var tradeProposalDTO structs.NBATradeProposalDTO
	err := json.NewDecoder(r.Body).Decode(&tradeProposalDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.CreateTradeProposal(tradeProposalDTO)

	// recruitingProfile := managers.CreateRecruitingProfileForRecruit(tradeProposalDTO)
	fmt.Fprintf(w, "New Trade Proposal Created")
}

// Accept Trade Offer
func AcceptTradeOffer(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	proposalID := vars["proposalID"]
	if len(proposalID) == 0 {
		panic("User did not provide a proposalID")
	}

	managers.AcceptTradeProposal(proposalID)

	json.NewEncoder(w).Encode("Proposal " + proposalID + " has been accepted.")
}

// Reject Trade Offer
func RejectTradeOffer(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	proposalID := vars["proposalID"]
	if len(proposalID) == 0 {
		panic("User did not provide a proposalID")
	}

	managers.RejectTradeProposal(proposalID)

	json.NewEncoder(w).Encode("Proposal " + proposalID + " has been accepted.")
}

// Cancels Trade Offer
func CancelTradeOffer(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	proposalID := vars["proposalID"]
	if len(proposalID) == 0 {
		panic("User did not provide a proposalID")
	}

	managers.CancelTradeProposal(proposalID)

	json.NewEncoder(w).Encode("Proposal " + proposalID + " has been accepted.")
}

// SyncAcceptedTrade -- Admin approve a trade
func SyncAcceptedTrade(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	proposalID := vars["proposalID"]
	if len(proposalID) == 0 {
		panic("User did not provide a proposalID")
	}

	managers.SyncAcceptedTrade(proposalID)

	json.NewEncoder(w).Encode("Proposal " + proposalID + " has been accepted.")
}

// SyncAcceptedTrade -- Admin approve a trade
func VetoAcceptedTrade(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	proposalID := vars["proposalID"]
	if len(proposalID) == 0 {
		panic("User did not provide a proposalID")
	}

	managers.VetoTrade(proposalID)

	json.NewEncoder(w).Encode("Proposal " + proposalID + " has been accepted.")
}

// CleanUpRejectedTrades -- Remove all rejected trades from the DB
func CleanUpRejectedTrades(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	managers.RemoveRejectedTrades()

	json.NewEncoder(w).Encode("Removed all rejected trades from the interface.")
}
