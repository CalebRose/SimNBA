package controller

import (
	"encoding/json"
	"fmt"
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

	fmt.Fprintf(w, "New Promise Created")
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

	fmt.Fprintf(w, "Promise Updated")
}

func CancelPromise(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	promiseID := vars["promiseID"]

	if len(promiseID) == 0 {
		panic("User did not provide Promise ID")
	}

	managers.CancelPromise(promiseID)

	fmt.Fprintf(w, "Promise Cancelled.")
}
