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
	data := managers.GetBootstrapDataLanding(collegeID, proID)
	json.NewEncoder(w).Encode(data)
}

func BootstrapLandingData(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	collegeID := vars["collegeID"]
	proID := vars["proID"]
	data := managers.GetBootstrapDataLanding(collegeID, proID)
	json.NewEncoder(w).Encode(data)

}

func BootstrapTeamRosterData(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	collegeID := vars["collegeID"]
	proID := vars["proID"]
	data := managers.GetBootstrapDataTeamRoster(collegeID, proID)
	json.NewEncoder(w).Encode(data)

}

func BootstrapRecruitingData(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	collegeID := vars["collegeID"]
	data := managers.GetBootstrapDataRecruiting(collegeID)
	json.NewEncoder(w).Encode(data)

}

func BootstrapFreeAgencyData(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	proID := vars["proID"]
	data := managers.GetBootstrapDataFreeAgency(proID)
	json.NewEncoder(w).Encode(data)

}

func BootstrapSchedulingData(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	username := vars["username"]
	collegeID := vars["collegeID"]
	seasonID := vars["seasonID"]
	data := managers.GetBootstrapDataScheduling(username, collegeID, seasonID)
	json.NewEncoder(w).Encode(data)

}

func BootstrapDraftData(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	proID := vars["proID"]
	data := managers.GetBootstrapDataDraft(proID)
	json.NewEncoder(w).Encode(data)

}

func BootstrapPortalData(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	collegeID := vars["collegeID"]
	data := managers.GetBootstrapDataPortal(collegeID)
	json.NewEncoder(w).Encode(data)

}

func BootstrapGameplanData(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	collegeID := vars["collegeID"]
	proID := vars["proID"]
	data := managers.GetBootstrapDataGameplan(collegeID, proID)
	json.NewEncoder(w).Encode(data)
}

func BootstrapNewsData(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	collegeID := vars["collegeID"]
	proID := vars["proID"]
	data := managers.GetNewsBootstrap(collegeID, proID)
	json.NewEncoder(w).Encode(data)

}
