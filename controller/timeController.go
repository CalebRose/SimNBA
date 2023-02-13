package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimNBA/managers"
)

// GetCurrentTimestamp - Get the Current Global Timestamp
func GetCurrentTimestamp(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	timestamp := managers.GetTimestamp()

	json.NewEncoder(w).Encode(timestamp)
}
