package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimNBA/managers"
	"github.com/gorilla/mux"
)

func GetNewsFeed(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	league := vars["league"]
	teamID := vars["teamID"]

	if league == "CBB" {
		newsLogs := managers.GetCBBRelatedNews(teamID)
		json.NewEncoder(w).Encode(newsLogs)
	} else {
		newsLogs := managers.GetNBARelatedNews(teamID)
		json.NewEncoder(w).Encode(newsLogs)
	}
}
