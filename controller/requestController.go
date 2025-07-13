package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimNBA/managers"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/gorilla/mux"
)

func GetTeamRequests(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	requests := managers.GetAllTeamRequests()

	json.NewEncoder(w).Encode(requests)
}

func GetNBATeamRequests(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	requests := managers.GetAllNBATeamRequests()

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

	json.NewEncoder(w).Encode(request)
}

func CreateNBATeamRequest(w http.ResponseWriter, r *http.Request) {
	var request structs.NBARequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.CreateNBATeamRequest(request)

	json.NewEncoder(w).Encode(request)
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

	json.NewEncoder(w).Encode(request)
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

	json.NewEncoder(w).Encode(request)
}

func ApproveNBATeamRequest(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	var request structs.NBARequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil || request.ID == 0 {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.ApproveNBATeamRequest(request)

	json.NewEncoder(w).Encode(request)
}

func RejectNBATeamRequest(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	var request structs.NBARequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.RejectNBATeamRequest(request)

	json.NewEncoder(w).Encode(request)
}

func RemoveNBAUserFromNBATeam(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	var request structs.NBARequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.RemoveUserFromNBATeam(request)

	json.NewEncoder(w).Encode(request)
}

func ViewCBBTeamUponRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]
	if len(teamID) == 0 {
		panic("User did not provide TeamID")
	}

	res := managers.GetCBBTeamForAvailableTeamsPage(teamID)

	json.NewEncoder(w).Encode(res)
}

func ViewNBATeamUponRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]
	if len(teamID) == 0 {
		panic("User did not provide TeamID")
	}

	res := managers.GetNBATeamForAvailableTeamsPage(teamID)

	json.NewEncoder(w).Encode(res)
}
