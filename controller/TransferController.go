package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/CalebRose/SimNBA/managers"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/gorilla/mux"
)

func ProcessTransferIntention(w http.ResponseWriter, r *http.Request) {
	managers.ProcessTransferIntention()
}

func CreatePromise(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	var createPromiseDto structs.CollegePromise
	err := json.NewDecoder(r.Body).Decode(&createPromiseDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	promise := managers.CreatePromise(createPromiseDto)

	json.NewEncoder(w).Encode(promise)
}

func UpdatePromise(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	var createPromiseDto structs.CollegePromise
	err := json.NewDecoder(r.Body).Decode(&createPromiseDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.UpdatePromise(createPromiseDto)
}

func CancelPromise(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	promiseID := vars["promiseID"]

	if len(promiseID) == 0 {
		panic("User did not provide Promise ID")
	}

	managers.CancelPromise(promiseID)
}

func GetPromiseByPlayerID(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	id := vars["playerID"]
	teamID := vars["teamID"]
	if len(id) == 0 {
		panic("User did not provide proper IDs")
	}

	promise := managers.GetCollegePromiseByCollegePlayerID(id, teamID)

	encodedJson, err := json.Marshal(promise)
	if err != nil {
		log.Printf("Error encoding JSON: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(encodedJson)
}

func GetTransferPortalPageData(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	teamID := vars["teamID"]
	if len(teamID) == 0 {
		panic("User did not provide proper IDs")
	}

	data := managers.GetTransferPortalData(teamID)

	encodedJson, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error encoding JSON: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(encodedJson)
}

func AddTransferPlayerToBoard(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	var transferPortalProfile structs.TransferPortalProfile
	err := json.NewDecoder(r.Body).Decode(&transferPortalProfile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	profile := managers.AddTransferPlayerToBoard(transferPortalProfile)

	json.NewEncoder(w).Encode(profile)
}

func RemovePlayerFromTransferPortalBoard(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	var transferPortalProfile structs.TransferPortalProfile
	err := json.NewDecoder(r.Body).Decode(&transferPortalProfile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.RemovePlayerFromTransferPortalBoard(transferPortalProfile)
	json.NewEncoder(w).Encode(true)
}

func SaveTransferBoard(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	var transferPortalProfile structs.UpdateTransferPortalBoard
	err := json.NewDecoder(r.Body).Decode(&transferPortalProfile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.AllocatePointsToTransferPlayer(transferPortalProfile)
}

func GetScoutingDataByTransfer(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	id := vars["id"]
	if len(id) == 0 {
		panic("User did not provide scout profile id")
	}

	data := managers.GetTransferScoutingDataByPlayerID(id)

	json.NewEncoder(w).Encode(data)
}

func FillUpTransferBoardsAI(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	managers.AICoachFillBoardsPhase()

	json.NewEncoder(w).Encode("AI Boards for Transfer Portal Complete.")
}

func AllocateAndPromisePlayersAI(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	managers.AICoachAllocateAndPromisePhase()

	json.NewEncoder(w).Encode("Allocated and promised.")
}

func SyncPromises(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	managers.SyncPromises()

	json.NewEncoder(w).Encode("Allocated and promised.")
}

func ExportPortalPlayersToCSV(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")
	managers.ExportTransferPortalToCSV(w)
}
