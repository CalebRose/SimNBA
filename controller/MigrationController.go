package controller

import (
	"fmt"
	"net/http"

	"github.com/CalebRose/SimNBA/managers"
)

func MigrateMissingRecruits(w http.ResponseWriter, r *http.Request) {
	managers.MigrateMissingRecruits()
}

func Migrate2026Data(w http.ResponseWriter, r *http.Request) {
	managers.Migration2026Main()

	fmt.Println("Migration Complete.")
	w.WriteHeader(http.StatusOK)
}

func GenerateCollegeAndNBALineupStructs(w http.ResponseWriter, r *http.Request) {
	managers.GenerateCollegeAndNBALineupStructs()

	fmt.Println("Lineup Generation Complete.")
	w.WriteHeader(http.StatusOK)
}
