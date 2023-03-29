package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CalebRose/SimNBA/managers"
	"github.com/CalebRose/SimNBA/structs"
)

func GetTeamRequests(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	requests := managers.GetAllTeamRequests()

	json.NewEncoder(w).Encode(requests)
}

func GetNBATeamRequests(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	requests := managers.GetAllTeamRequests()

	json.NewEncoder(w).Encode(requests)
}

func CreateTeamRequest(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	var request structs.Request
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.CreateTeamRequest(request)

	fmt.Fprintf(w, "Request Successfully Created")
}

func CreateNBATeamRequest(w http.ResponseWriter, r *http.Request) {
	var request structs.NBARequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.CreateNBATeamRequest(request)

	fmt.Fprintf(w, "Request Successfully Created")
}

func ApproveTeamRequest(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	var request structs.Request
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil || request.ID == 0 {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.ApproveTeamRequest(request)

	fmt.Fprintf(w, "Request: %+v", request)
}

func RejectTeamRequest(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	var request structs.Request

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.RejectTeamRequest(request)

	fmt.Fprintf(w, "Request: %+v", request)
}

func ApproveNBATeamRequest(w http.ResponseWriter, r *http.Request) {
	var request structs.NBARequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil || request.ID == 0 {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.ApproveNBATeamRequest(request)

	fmt.Fprintf(w, "Request: %+v", request)
}

func RejectNBATeamRequest(w http.ResponseWriter, r *http.Request) {
	var request structs.NBARequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.RejectNBATeamRequest(request)

	fmt.Fprintf(w, "Request: %+v", request)
}

func RemoveNBAUserFromNBATeam(w http.ResponseWriter, r *http.Request) {
	var request structs.NBARequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.RemoveUserFromNBATeam(request)

	// json.NewEncoder(w).Encode(team)
}
