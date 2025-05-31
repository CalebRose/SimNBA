package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimNBA/managers"
)

func ISLIdentifyYouthPlayers(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	managers.ISLIdentityPhase()

	json.NewEncoder(w).Encode("Identified players.")
}

func ISLScoutYouthPlayers(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	managers.ISLScoutingPhase()

	json.NewEncoder(w).Encode("Scouted players.")
}

func ISLInvestYouthPlayers(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	managers.ISLInvestingPhase()

	json.NewEncoder(w).Encode("Invested in players.")
}

func ISLSyncYouthPlayers(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	managers.PickUpISLPlayers()

	json.NewEncoder(w).Encode("Synced players.")
}

func ISLGenerateNewBatch(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	managers.GenerateInternationalPlayersByTeam()

	json.NewEncoder(w).Encode("Synced players.")
}
