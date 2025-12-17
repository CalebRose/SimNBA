package controller

import (
	"net/http"

	"github.com/CalebRose/SimNBA/managers"
)

func MigrateMissingRecruits(w http.ResponseWriter, r *http.Request) {
	managers.MigrateMissingRecruits()
}
