package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/managers"
	"github.com/CalebRose/SimNBA/structs"
)

func GetTeamRequests(w http.ResponseWriter, r *http.Request) {
	db := dbprovider.GetInstance().GetDB()

	var requests []structs.RequestDTO
	db.Raw("SELECT requests.id, requests.team_id, teams.team, teams.abbr, requests.username, teams.conference, teams.is_nba, requests.is_approved FROM simfbaah_simnba.requests INNER JOIN simfbaah_simnba.teams on teams.id = requests.team_id WHERE requests.deleted_at is null AND requests.is_approved = 0").
		Scan(&requests)
	// db.Where("deleted_date is null AND is_approved = 0").Find(&requests)
	json.NewEncoder(w).Encode(requests)
}

func CreateTeamRequest(w http.ResponseWriter, r *http.Request) {
	db := dbprovider.GetInstance().GetDB()

	var request structs.Request
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(request)

	db.Create(&request)

	fmt.Fprintf(w, "Request Successfully Created")
}

func ApproveTeamRequest(w http.ResponseWriter, r *http.Request) {
	db := dbprovider.GetInstance().GetDB()

	var request structs.Request
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil || request.ID == 0 {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	request.ApproveTeamRequest()

	fmt.Println("Team Approved...")

	db.Save(&request)

	fmt.Println("Assigning team...")

	// Assign Team
	team := managers.GetTeamByTeamID(strconv.Itoa(request.TeamID))

	team.AssignUserToTeam(request.Username)

	db.Save(&team)

	// db.Model(&team).Where("id = ?", request.TeamID).Update("coach", request.Username)

	fmt.Fprintf(w, "Request: %+v", request)
}

func RejectTeamRequest(w http.ResponseWriter, r *http.Request) {
	db := dbprovider.GetInstance().GetDB()

	var request structs.Request

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	request.RejectTeamRequest()

	db.Delete(&request)

	fmt.Fprintf(w, "Request: %+v", request)
}
