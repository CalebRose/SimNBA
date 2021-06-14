package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/managers"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/gorilla/mux"
)

// AllRecruitsByProfileID - Get all Recruits By A Team's Recruiting Profile
func AllRecruitsByProfileID(w http.ResponseWriter, r *http.Request) {
	db := dbprovider.GetInstance().GetDB()

	vars := mux.Vars(r)
	profileID := vars["profileId"]
	if len(profileID) == 0 {
		panic("User did not provide a Recruiting Profile ID")
	}
	var recruitPoints []structs.RecruitingPoints

	db.Preload("Recruit").Where("profile_id = ?", profileID).Find(&recruitPoints)
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
	db := dbprovider.GetInstance().GetDB()

	var recruitPointsDto structs.CreateRecruitPointsDto
	err := json.NewDecoder(r.Body).Decode(&recruitPointsDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Check for Existing Points profile
	existingRecruitingPointsProfile := managers.
		GetRecruitingPointsProfileByPlayerId(strconv.Itoa(recruitPointsDto.PlayerId), strconv.Itoa(recruitPointsDto.ProfileId))

	// If Recruit Already Exists
	if existingRecruitingPointsProfile.PlayerID != 0 && existingRecruitingPointsProfile.ProfileID != 0 {
		existingRecruitingPointsProfile.ReplaceRecruitToBoard()
		db.Save(&existingRecruitingPointsProfile)
		json.NewEncoder(w).Encode(existingRecruitingPointsProfile)
		return
	}

	recruitingPointProfile := structs.RecruitingPoints{
		SeasonID:               recruitPointsDto.SeasonId,
		PlayerID:               recruitPointsDto.PlayerId,
		ProfileID:              recruitPointsDto.ProfileId,
		Team:                   recruitPointsDto.Team,
		TotalPointsSpent:       0,
		CurrentPointsSpent:     0,
		Scholarship:            false,
		InterestLevel:          "None",
		InterestLevelThreshold: 0,
		Signed:                 false,
		RemovedFromBoard:       false,
	}

	db.Create(&recruitingPointProfile)

	json.NewEncoder(w).Encode(recruitingPointProfile)

	fmt.Fprintf(w, "New Recruiting Profile Created")
}

func AllocateRecruitingPointsForRecruit(w http.ResponseWriter, r *http.Request) {
	db := dbprovider.GetInstance().GetDB()

	var updateRecruitPointsDto structs.UpdateRecruitPointsDto
	err := json.NewDecoder(r.Body).Decode(&updateRecruitPointsDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	recruitingProfile := managers.GetRecruitingProfileByTeamId(strconv.Itoa(updateRecruitPointsDto.ProfileId))

	recruitingPointsProfile := managers.GetRecruitingPointsProfileByPlayerId(
		strconv.Itoa(updateRecruitPointsDto.PlayerId),
		strconv.Itoa(updateRecruitPointsDto.ProfileId),
	)

	recruitingProfile.AllocateSpentPoints(updateRecruitPointsDto.SpentPoints)
	if recruitingProfile.SpentPoints > recruitingProfile.WeeklyPoints {
		fmt.Printf("Recruiting Profile " + strconv.Itoa(updateRecruitPointsDto.ProfileId) + " cannot spend more points than weekly amount")
		return
	}

	recruitingPointsProfile.AllocatePoints(updateRecruitPointsDto.SpentPoints)

	db.Save(&recruitingPointsProfile)

	db.Save(&recruitingProfile)

	fmt.Printf("Updated Recruiting Points Profile")
}

func SendScholarshipToRecruit(w http.ResponseWriter, r *http.Request) {
	db := dbprovider.GetInstance().GetDB()

	var updateRecruitPointsDto structs.UpdateRecruitPointsDto
	err := json.NewDecoder(r.Body).Decode(&updateRecruitPointsDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	recruitingProfile := managers.GetRecruitingProfileByTeamId(strconv.Itoa(updateRecruitPointsDto.ProfileId))

	if recruitingProfile.ScholarshipsAvailable == 0 {
		fmt.Printf("\nTeamId: " + strconv.Itoa(updateRecruitPointsDto.ProfileId) + " does not have any availabe scholarships")
		return
	}

	recruitingPointsProfile := managers.GetRecruitingPointsProfileByPlayerId(
		strconv.Itoa(updateRecruitPointsDto.PlayerId),
		strconv.Itoa(updateRecruitPointsDto.ProfileId),
	)

	if recruitingPointsProfile.Scholarship {
		fmt.Printf("\nRecruit " + strconv.Itoa(recruitingPointsProfile.PlayerID) + "already has a scholarship")
		return
	}

	recruitingPointsProfile.AllocateScholarship()
	recruitingProfile.SubtractScholarshipsAvailable()

	db.Save(&recruitingPointsProfile)
	fmt.Printf("\nScholarship allocated to player " + strconv.Itoa(recruitingPointsProfile.PlayerID) + ". Record saved")
	db.Save(&recruitingProfile)
	fmt.Printf("\nProfile: " + strconv.Itoa(recruitingProfile.TeamID) + " Saved")
}

func RevokeScholarshipFromRecruit(w http.ResponseWriter, r *http.Request) {
	db := dbprovider.GetInstance().GetDB()

	var updateRecruitPointsDto structs.UpdateRecruitPointsDto
	err := json.NewDecoder(r.Body).Decode(&updateRecruitPointsDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	recruitingProfile := managers.GetRecruitingProfileByTeamId(strconv.Itoa(updateRecruitPointsDto.ProfileId))

	recruitingPointsProfile := managers.GetRecruitingPointsProfileByPlayerId(
		strconv.Itoa(updateRecruitPointsDto.PlayerId),
		strconv.Itoa(updateRecruitPointsDto.ProfileId),
	)

	if recruitingPointsProfile.Scholarship {
		fmt.Printf("\nCannot revoke an inexistant scholarship from Recruit " + strconv.Itoa(recruitingPointsProfile.PlayerID))
		return
	}

	recruitingPointsProfile.RevokeScholarship()
	recruitingProfile.ReallocateScholarship()

	db.Save(&recruitingPointsProfile)
	fmt.Printf("\nScholarship revoked from player " + strconv.Itoa(recruitingPointsProfile.PlayerID) + ". Record saved")
	db.Save(&recruitingProfile)
	fmt.Printf("\nProfile: " + strconv.Itoa(recruitingProfile.TeamID) + " Saved")
}

func RemoveRecruitFromBoard(w http.ResponseWriter, r *http.Request) {
	db := dbprovider.GetInstance().GetDB()

	var updateRecruitPointsDto structs.UpdateRecruitPointsDto
	err := json.NewDecoder(r.Body).Decode(&updateRecruitPointsDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	recruitingPointsProfile := managers.GetRecruitingPointsProfileByPlayerId(
		strconv.Itoa(updateRecruitPointsDto.PlayerId),
		strconv.Itoa(updateRecruitPointsDto.ProfileId),
	)

	if recruitingPointsProfile.RemovedFromBoard {
		panic("Recruit already removed from board")
	}

	recruitingPointsProfile.RemoveRecruitFromBoard()
	db.Save(&recruitingPointsProfile)
}

func SaveRecruitingBoard(w http.ResponseWriter, r *http.Request) {
	db := dbprovider.GetInstance().GetDB()

	var updateRecruitingBoardDto structs.UpdateRecruitingBoardDto
	err := json.NewDecoder(r.Body).Decode(&updateRecruitingBoardDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var teamId = strconv.Itoa(updateRecruitingBoardDto.TeamID)

	var profile = managers.GetOnlyRecruitingProfileByTeamId(teamId)

	var recruitingPoints = managers.GetRecruitingPointsByTeamId(teamId)

	var updatedRecruits = updateRecruitingBoardDto.Recruits

	currentPoints := 0

	for i := 0; i < len(recruitingPoints); i++ {
		updatedRecruit := managers.GetRecruitFromRecruitsList(recruitingPoints[i].PlayerID, updatedRecruits)

		if updatedRecruit.CurrentPointsSpent > 0 &&
			recruitingPoints[i].CurrentPointsSpent != updatedRecruit.CurrentPointsSpent {

			// Allocate Points to Profile
			currentPoints += updatedRecruit.CurrentPointsSpent
			profile.AllocateSpentPoints(currentPoints)
			// If total not surpassed, allocate to the recruit and continue
			if profile.SpentPoints <= profile.WeeklyPoints {
				recruitingPoints[i].AllocatePoints(updatedRecruit.CurrentPointsSpent)
				fmt.Println("Saving recruit " + strconv.Itoa(recruitingPoints[i].PlayerID))
				db.Save(&recruitingPoints[i])
			} else {
				panic("Error: Allocated more points for Profile " + strconv.Itoa(profile.TeamID) + " than what is allowed.")
			}
		}
	}

	// Save profile
	db.Save(&profile)

	fmt.Println("Updated Recruiting Profile and Players")
	w.WriteHeader(http.StatusOK)
}
