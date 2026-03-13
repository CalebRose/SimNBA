package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimNBA/managers"
)

func UpdateTeamProfileAffinities(w http.ResponseWriter, r *http.Request) {
	managers.UpdateTeamProfileAffinities()
	json.NewEncoder(w).Encode("Done!")
}
