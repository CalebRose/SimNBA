package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/CalebRose/SimNBA/managers"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/gorilla/mux"
)

// AllRecruitsByProfileID - Get all Recruits By A Team's Recruiting Profile
func AllRecruitsByProfileID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	profileID := vars["profileId"]
	if len(profileID) == 0 {
		panic("User did not provide a Recruiting Profile ID")
	}

	var recruitPoints = managers.GetAllRecruitsByProfileID(profileID)

	json.NewEncoder(w).Encode(recruitPoints)
}

// RecruitingProfileByTeamID - Get Recruiting Profile by TeamID
func RecruitingProfileByTeamID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamId := vars["teamId"]
	if len(teamId) == 0 {
		panic("User did not provide TeamID")
	}

	profile := managers.GetRecruitingProfileByTeamId(teamId)
	json.NewEncoder(w).Encode(profile)
}

func CreateRecruitingPointsProfileForRecruit(w http.ResponseWriter, r *http.Request) {

	var recruitPointsDto structs.CreateRecruitPointsDto
	err := json.NewDecoder(r.Body).Decode(&recruitPointsDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	recruitingProfile := managers.CreateRecruitingPointsProfileForRecruit(recruitPointsDto)

	json.NewEncoder(w).Encode(recruitingProfile)

	fmt.Fprintf(w, "New Recruiting Profile Created")
}

func AllocateRecruitingPointsForRecruit(w http.ResponseWriter, r *http.Request) {
	var updateRecruitPointsDto structs.UpdateRecruitPointsDto
	err := json.NewDecoder(r.Body).Decode(&updateRecruitPointsDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.AllocateRecruitingPointsForRecruit(updateRecruitPointsDto)

	fmt.Printf("Updated Recruiting Points Profile")
}

func SendScholarshipToRecruit(w http.ResponseWriter, r *http.Request) {
	var updateRecruitPointsDto structs.UpdateRecruitPointsDto
	err := json.NewDecoder(r.Body).Decode(&updateRecruitPointsDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	recruitingPointsProfile, recruitingProfile := managers.SendScholarshipToRecruit(updateRecruitPointsDto)
	fmt.Printf("\nScholarship allocated to player " + strconv.Itoa(recruitingPointsProfile.PlayerID) + ". Record saved")
	fmt.Printf("\nProfile: " + strconv.Itoa(recruitingProfile.TeamID) + " Saved")
}

func RevokeScholarshipFromRecruit(w http.ResponseWriter, r *http.Request) {
	var updateRecruitPointsDto structs.UpdateRecruitPointsDto
	err := json.NewDecoder(r.Body).Decode(&updateRecruitPointsDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	recruitingPointsProfile, recruitingProfile := managers.RevokeScholarshipFromRecruit(updateRecruitPointsDto)

	fmt.Printf("\nScholarship revoked from player " + strconv.Itoa(recruitingPointsProfile.PlayerID) + ". Record saved")
	fmt.Printf("\nProfile: " + strconv.Itoa(recruitingProfile.TeamID) + " Saved")
}

func RemoveRecruitFromBoard(w http.ResponseWriter, r *http.Request) {
	var updateRecruitPointsDto structs.UpdateRecruitPointsDto
	err := json.NewDecoder(r.Body).Decode(&updateRecruitPointsDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	recruitingPointsProfile := managers.RemoveRecruitFromBoard(updateRecruitPointsDto)

	fmt.Printf("\nPlayer " + strconv.Itoa(recruitingPointsProfile.PlayerID) + " removed from board.")
}

func SaveRecruitingBoard(w http.ResponseWriter, r *http.Request) {
	var updateRecruitingBoardDto structs.UpdateRecruitingBoardDto
	err := json.NewDecoder(r.Body).Decode(&updateRecruitingBoardDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	recruitingProfile := managers.UpdateRecruitingProfile(updateRecruitingBoardDto)

	fmt.Println("Updated Recruiting Profile " + strconv.Itoa(recruitingProfile.TeamID) + " and all associated players")
	w.WriteHeader(http.StatusOK)
}
