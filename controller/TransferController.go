package controller

import (
	"net/http"

	"github.com/CalebRose/SimNBA/managers"
)

func ProcessTransferIntention(w http.ResponseWriter, r *http.Request) {
	managers.ProcessTransferIntention()
}
