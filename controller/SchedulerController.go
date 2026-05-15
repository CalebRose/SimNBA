package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimNBA/managers"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/gorilla/mux"
)

// CreateCBBGameRequest accepts a CBBGameRequest body and persists it.
func CreateCBBGameRequest(w http.ResponseWriter, r *http.Request) {
	var request structs.CBBGameRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	managers.CreateCBBGameRequest(request)
	json.NewEncoder(w).Encode(true)
}

// AcceptCBBGameRequest marks the request as accepted.
func AcceptCBBGameRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := vars["requestID"]
	if len(requestID) == 0 {
		panic("User did not provide a requestID")
	}
	managers.AcceptCBBGameRequest(requestID)
	json.NewEncoder(w).Encode(true)
}

// RejectCBBGameRequest deletes the request.
func RejectCBBGameRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := vars["requestID"]
	if len(requestID) == 0 {
		panic("User did not provide a requestID")
	}
	managers.RejectCBBGameRequest(requestID)
	json.NewEncoder(w).Encode(true)
}

// ProcessCBBGameRequest converts an accepted request into a Match record.
func ProcessCBBGameRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := vars["requestID"]
	if len(requestID) == 0 {
		panic("User did not provide a requestID")
	}
	managers.ProcessCBBGameRequest(requestID)
	json.NewEncoder(w).Encode(true)
}

// VetoCBBGameRequest deletes the request via admin veto.
func VetoCBBGameRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := vars["requestID"]
	if len(requestID) == 0 {
		panic("User did not provide a requestID")
	}
	managers.VetoCBBGameRequest(requestID)
	json.NewEncoder(w).Encode(true)
}
