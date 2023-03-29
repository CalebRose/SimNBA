package controller

import (
	"net/http"

	"github.com/CalebRose/SimNBA/managers"
)

func ImportNBATeamsAndArenas(w http.ResponseWriter, r *http.Request) {
	managers.ImportNBATeamsAndArenas()
}
