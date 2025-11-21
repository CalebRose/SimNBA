package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/CalebRose/SimNBA/managers"
	"github.com/gorilla/mux"
)

func BootstrapTeamData(w http.ResponseWriter, r *http.Request) {
	data := managers.GetBootstrapTeams()
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func BootstrapBasketballData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collegeID := vars["collegeID"]
	proID := vars["proID"]
	data := managers.GetBootstrapData(collegeID, proID)
	json.NewEncoder(w).Encode(data)
}

func SecondBootstrapBasketballData(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	collegeID := vars["collegeID"]
	proID := vars["proID"]
	data := managers.GetSecondBootstrapData(collegeID, proID)
	json.NewEncoder(w).Encode(data)
}

func ThirdBootstrapBasketballData(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	collegeID := vars["collegeID"]
	proID := vars["proID"]
	data := managers.GetThirdBootstrapData(collegeID, proID)
	json.NewEncoder(w).Encode(data)
}

func BootstrapNewsData(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	collegeID := vars["collegeID"]
	proID := vars["proID"]
	data := managers.GetNewsBootstrap(collegeID, proID)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
