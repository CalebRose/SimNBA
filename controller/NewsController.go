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

func GetBBAInbox(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cfbID := vars["cbbID"]
	nflID := vars["nbaID"]

	inbox := managers.GetBBAInbox(cfbID, nflID)
	json.NewEncoder(w).Encode(inbox)
}

func ToggleNotificationAsRead(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	notiID := vars["notiID"]
	managers.ToggleNotification(notiID)
	json.NewEncoder(w).Encode("Toggled Notification")
}

func DeleteNotification(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	notiID := vars["notiID"]
	managers.DeleteNotification(notiID)
	json.NewEncoder(w).Encode("Toggled Notification")
}
