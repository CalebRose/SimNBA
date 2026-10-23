package controller

import (
	"net/http"

	"github.com/CalebRose/SimNBA/managers"
)

func FixRecruitingProfileModifiers(w http.ResponseWriter, r *http.Request) {
	managers.FixExistingModifiersForRecruits()
}
