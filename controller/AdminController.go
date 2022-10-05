package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimNBA/managers"
)

func GeneratePlayers(w http.ResponseWriter, r *http.Request) {
	players := managers.GenerateNewTeams()
	json.NewEncoder(w).Encode(players)
}
