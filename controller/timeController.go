package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/managers"
)

// GetCurrentTimestamp - Get the Current Global Timestamp
func GetCurrentTimestamp(w http.ResponseWriter, r *http.Request) {
	db := dbprovider.GetInstance().GetDB()

	timestamp := managers.GetTimestamp(db)
	json.NewEncoder(w).Encode(timestamp)
}
