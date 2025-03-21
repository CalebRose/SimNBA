package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimNBA/managers"
	"github.com/gorilla/mux"
)

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
