package controller

import (
	"fmt"
	"net/http"

	"github.com/CalebRose/SimNBA/managers"
)

func MigrateFaceData(w http.ResponseWriter, r *http.Request) {
	managers.MigrateFaceDataToRecruits()
	managers.MigrateFaceDataToCollegePlayers()
	managers.MigrateFaceDataToProPlayers()

	fmt.Println("All Faces have been generated")
	w.WriteHeader(http.StatusOK)
}
