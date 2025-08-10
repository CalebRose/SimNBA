package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/CalebRose/SimNBA/managers"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/gorilla/mux"
)

// GetRecruitingDataForOverviewPage - Returns all data needed for Recruiting Overview
func GetRecruitingDataForOverviewPage(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
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
	EnableCors(&w)
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

// GetRecruitingDataForTeamBoardPage - Returns all data needed for team board
func GetRecruitingClassByTeamID(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	teamID := vars["teamID"]

	if len(teamID) == 0 {
		panic("User did not provide teamID")
	}

	class := managers.GetRecruitingClassByTeamID(teamID)

	json.NewEncoder(w).Encode(class)
}

// GetAllRecruitingProfiles
func GetAllRecruitingProfiles(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	recruitingProfiles := managers.GetRecruitingProfileForRecruitSync()

	json.NewEncoder(w).Encode(recruitingProfiles)
}

func CreateRecruit(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
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
	EnableCors(&w)
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

func AddRecruitToBoardV2(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	var recruitPointsDto structs.CreateRecruitProfileDtoV2
	err := json.NewDecoder(r.Body).Decode(&recruitPointsDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	recruitingProfile := managers.AddRecruitToTeamBoardV2(recruitPointsDto)

	json.NewEncoder(w).Encode(recruitingProfile)
}

func AllocateRecruitingPointsForRecruit(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	var updateRecruitPointsDto structs.UpdateRecruitPointsDto
	err := json.NewDecoder(r.Body).Decode(&updateRecruitPointsDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.AllocateRecruitingPointsForRecruit(updateRecruitPointsDto)
}

func SendScholarshipToRecruit(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	var updateRecruitPointsDto structs.UpdateRecruitPointsDto
	err := json.NewDecoder(r.Body).Decode(&updateRecruitPointsDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.SendScholarshipToRecruit(updateRecruitPointsDto)

}

func SendScholarshipToRecruitV2(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	var updateRecruitPointsDto structs.UpdateRecruitPointsDtoV2
	err := json.NewDecoder(r.Body).Decode(&updateRecruitPointsDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.SendScholarshipToRecruitV2(updateRecruitPointsDto)

}

func RemoveRecruitFromBoard(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	var updateRecruitPointsDto structs.UpdateRecruitPointsDto
	err := json.NewDecoder(r.Body).Decode(&updateRecruitPointsDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.RemoveRecruitFromBoard(updateRecruitPointsDto)
}

func SaveRecruitingBoard(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	var updateRecruitingBoardDto structs.UpdateRecruitingBoardDto
	err := json.NewDecoder(r.Body).Decode(&updateRecruitingBoardDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.UpdateRecruitingProfile(updateRecruitingBoardDto)

	w.WriteHeader(http.StatusOK)
}

func ExportCroots(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	w.Header().Set("Content-Type", "text/csv")
	managers.ExportCroots(w)
}

// Determine the size of a year's recruiting class based on current roster number, and players leaving/potentially leaving
func DetermineRecruitingClassSize(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	managers.DetermineRecruitingClassSize()
	json.NewEncoder(w).Encode("CLASS SIZE COMPLETE")
}
