package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CalebRose/SimNBA/managers"
)

func MigrateFaceData(w http.ResponseWriter, r *http.Request) {
	// managers.MigrateFaceDataToRecruits()
	managers.MigrateFaceDataToCollegePlayers()
	// managers.MigrateFaceDataToHistoricCollegePlayers()
	managers.MigrateFaceDataToProPlayers()

	fmt.Println("All Faces have been generated")
	w.WriteHeader(http.StatusOK)
}

func GetAllFaces(w http.ResponseWriter, r *http.Request) {
	faceData := managers.GetAllFaces()

	fmt.Println("Face data retrieved")
	json.NewEncoder(w).Encode(faceData)
}
