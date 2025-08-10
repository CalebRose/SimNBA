package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimNBA/managers"
)

func GenerateCoaches(w http.ResponseWriter, r *http.Request) {
	managers.GenerateCoachesForAITeams()
	json.NewEncoder(w).Encode("GENERATION COMPLETE")
}

func GenerateTestPlayers(w http.ResponseWriter, r *http.Request) {
	managers.GenerateTestPlayersForTP()
	json.NewEncoder(w).Encode("GENERATION COMPLETE")
}

func GeneratePlayers(w http.ResponseWriter, r *http.Request) {
	managers.GenerateNewTeams()
	json.NewEncoder(w).Encode("GENERATION COMPLETE")
}

func GenerateCroots(w http.ResponseWriter, r *http.Request) {
	// managers.ProgressStandings()
	// managers.RunDeclarationsAlgorithm()
	// managers.DetermineRecruitingClassSize()
	// managers.GenerateCollegeStandings()
	// managers.GenerateNBAStandings()
	managers.GenerateCroots()
	json.NewEncoder(w).Encode("GENERATION COMPLETE")
}

func GenerateCollegeWalkons(w http.ResponseWriter, r *http.Request) {
	managers.GenerateCollegeWalkons()
	json.NewEncoder(w).Encode("GENERATION COMPLETE")
}

func GenerateInternationalPlayers(w http.ResponseWriter, r *http.Request) {
	// managers.GenerateInternationalPlayers()
	managers.GenerateAdditionalWorldCupPlayers()
	json.NewEncoder(w).Encode("GENERATION COMPLETE")
}

func MoveISLPlayersToDraft(w http.ResponseWriter, r *http.Request) {
	managers.MoveISLPlayerToDraft()
	json.NewEncoder(w).Encode("GENERATION COMPLETE")
}

func GenerateInternationalRoster(w http.ResponseWriter, r *http.Request) {
	managers.FormISLRosters()
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

func GenerateNewAttributes(w http.ResponseWriter, r *http.Request) {
	managers.GenerateNewAttributes()
	json.NewEncoder(w).Encode("GENERATION COMPLETE")
}

func FormISLRosters(w http.ResponseWriter, r *http.Request) {
	managers.GenerateInternationalPlayers()
	json.NewEncoder(w).Encode("GENERATION COMPLETE")
}
