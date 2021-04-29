package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CalebRose/SimNBA/managers"
	"github.com/jinzhu/gorm"
)

// GetCurrentTimestamp - Get the Current Global Timestamp
func GetCurrentTimestamp(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(c["db"], c["cs"])
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}

	defer db.Close()

	timestamp := managers.GetTimestamp(db)
	json.NewEncoder(w).Encode(timestamp)
}
