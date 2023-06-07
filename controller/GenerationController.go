package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimNBA/managers"
)

func GeneratePlayers(w http.ResponseWriter, r *http.Request) {
	managers.GenerateNewTeams()
	json.NewEncoder(w).Encode("GENERATION COMPLETE")
}

func GenerateCroots(w http.ResponseWriter, r *http.Request) {
	managers.GenerateCroots()
	json.NewEncoder(w).Encode("GENERATION COMPLETE")
}

func GenerateGlobalPlayerRecords(w http.ResponseWriter, r *http.Request) {
	managers.GenerateGlobalPlayerRecords()
	json.NewEncoder(w).Encode("GENERATION COMPLETE")
}

func GenerateAttributeSpecsForCollegeAndRecruits(w http.ResponseWriter, r *http.Request) {
	managers.GenerateAttributeSpecs()
	json.NewEncoder(w).Encode("GENERATION COMPLETE")
}

func GenerateGameplans(w http.ResponseWriter, r *http.Request) {
	managers.GenerateGameplans()
	json.NewEncoder(w).Encode("GENERATION COMPLETE")
}

func GenerateDraftWarRooms(w http.ResponseWriter, r *http.Request) {
	managers.GenerateDraftWarRooms()
	json.NewEncoder(w).Encode("GENERATION COMPLETE")
}

func GeneratePlaytimeExpectations(w http.ResponseWriter, r *http.Request) {
	managers.GeneratePlaytimeExpectations()
	json.NewEncoder(w).Encode("GENERATION COMPLETE")
}
