package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/CalebRose/SimNBA/managers"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/gorilla/mux"
)

// GetRecruitingDataForOverviewPage - Returns all data needed for Recruiting Overview
func GetRecruitingDataForOverviewPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]

	if len(teamID) == 0 {
		panic("User did not provide teamID")
	}

	var dashboardResponse structs.DashboardTeamProfileResponse

	recruitingProfile := managers.GetRecruitingProfileForDashboardByTeamID(teamID)

	dashboardResponse.SetTeamProfile(recruitingProfile)

	json.NewEncoder(w).Encode(dashboardResponse)

}

// GetRecruitingDataForTeamBoardPage - Returns all data needed for team board
func GetRecruitingDataForTeamBoardPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]

	if len(teamID) == 0 {
		panic("User did not provide teamID")
	}

	var teamBoardResponse structs.TeamBoardTeamProfileResponse

	recruitingProfile := managers.GetRecruitingProfileForTeamBoardByTeamID(teamID)

	teamBoardResponse.SetTeamProfile(recruitingProfile)

	json.NewEncoder(w).Encode(teamBoardResponse)
}

// GetAllRecruitingProfiles
func GetAllRecruitingProfiles(w http.ResponseWriter, r *http.Request) {
	recruitingProfiles := managers.GetRecruitingProfileForRecruitSync()

	json.NewEncoder(w).Encode(recruitingProfiles)
}

func CreateRecruit(w http.ResponseWriter, r *http.Request) {
	var dto structs.CreateRecruitDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// validate info from DTO
	if len(dto.FirstName) == 0 || len(dto.LastName) == 0 || dto.Overall == 0 {
		log.Fatalln("ERROR: Did not provide all information for recruit.")
	}

	managers.CreateRecruit(dto)
	fmt.Println(w, "New Recruit Created")
}

func AddRecruitToBoard(w http.ResponseWriter, r *http.Request) {

	var recruitPointsDto structs.CreateRecruitProfileDto
	err := json.NewDecoder(r.Body).Decode(&recruitPointsDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	recruitingProfile := managers.AddRecruitToTeamBoard(recruitPointsDto)

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
	fmt.Printf("\nScholarship allocated to player " + strconv.Itoa(int(recruitingPointsProfile.RecruitID)) + ". Record saved")
	fmt.Printf("\nProfile: " + strconv.Itoa(int(recruitingProfile.TeamID)) + " Saved")
}

func RemoveRecruitFromBoard(w http.ResponseWriter, r *http.Request) {
	var updateRecruitPointsDto structs.UpdateRecruitPointsDto
	err := json.NewDecoder(r.Body).Decode(&updateRecruitPointsDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	recruitingPointsProfile := managers.RemoveRecruitFromBoard(updateRecruitPointsDto)

	fmt.Printf("\nPlayer " + strconv.Itoa(int(recruitingPointsProfile.RecruitID)) + " removed from board.")
}

func SaveRecruitingBoard(w http.ResponseWriter, r *http.Request) {
	var updateRecruitingBoardDto structs.UpdateRecruitingBoardDto
	err := json.NewDecoder(r.Body).Decode(&updateRecruitingBoardDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	recruitingProfile := managers.UpdateRecruitingProfile(updateRecruitingBoardDto)

	fmt.Println("Updated Recruiting Profile " + strconv.Itoa(int(recruitingProfile.TeamID)) + " and all associated players")
	w.WriteHeader(http.StatusOK)
}

func ExportCroots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")
	managers.ExportCroots(w)
}
