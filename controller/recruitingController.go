package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/CalebRose/SimNBA/managers"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// AllRecruitsByProfileID - Get all Recruits By A Team's Recruiting Profile
func AllRecruitsByProfileID(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(c["db"], c["cs"])
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}

	defer db.Close()

	vars := mux.Vars(r)
	profileID := vars["profileId"]
	if len(profileID) == 0 {
		panic("User did not provide a Recruiting Profile ID")
	}
	var players []structs.Player

	db.Where("profile_id = ?", profileID).Find(&players)
	json.NewEncoder(w).Encode(players)
}

// RecruitingProfileByTeamID - Get Recruiting Profile by TeamID
func RecruitingProfileByTeamID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamId := vars["teamId"]
	if len(teamId) == 0 {
		panic("User did not provide TeamID")
	}

	db, err := gorm.Open(c["db"], c["cs"])
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}

	defer db.Close()

	profile := managers.GetRecruitingProfileByTeamId(db, teamId)
	json.NewEncoder(w).Encode(profile)
}

func CreateRecruitingPointsProfileForRecruit(w http.ResponseWriter, r *http.Request) {

	db, err := gorm.Open(c["db"], c["cs"])
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}

	defer db.Close()

	var recruitPointsDto structs.CreateRecruitPointsDto
	err = json.NewDecoder(r.Body).Decode(&recruitPointsDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	recruitingPointProfile := structs.RecruitingPoints{
		SeasonID:               recruitPointsDto.SeasonId,
		PlayerID:               recruitPointsDto.PlayerId,
		ProfileID:              recruitPointsDto.ProfileId,
		TotalPointsSpent:       0,
		CurrentPointsSpent:     0,
		Scholarship:            false,
		InterestLevel:          "None",
		InterestLevelThreshold: 0,
		Signed:                 false,
	}

	db.Create(&recruitingPointProfile)

	fmt.Fprintf(w, "New Recruiting Profile Created")
}

func AllocateRecruitingPointsForRecruit(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(c["db"], c["cs"])
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}

	defer db.Close()

	var updateRecruitPointsDto structs.UpdateRecruitPointsDto
	err = json.NewDecoder(r.Body).Decode(&updateRecruitPointsDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	recruitingProfile := managers.GetRecruitingProfileByTeamId(db, updateRecruitPointsDto.ProfileId)

	recruitingPointsProfile := managers.GetRecruitingPointsProfileByPlayerId(db,
		updateRecruitPointsDto.PlayerId,
		updateRecruitPointsDto.ProfileId)

	recruitingProfile.AllocateSpentPoints(updateRecruitPointsDto.SpentPoints)
	if recruitingProfile.SpentPoints > recruitingProfile.WeeklyPoints {
		fmt.Printf("Recruiting Profile " + updateRecruitPointsDto.ProfileId + " cannot spend more points than weekly amount")
		return
	}

	recruitingPointsProfile.AllocatePoints(updateRecruitPointsDto.SpentPoints)

	db.Save(&recruitingPointsProfile)

	db.Save(&recruitingProfile)

	fmt.Printf("Updated Recruiting Points Profile")
}

func SendScholarshipToRecruit(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(c["db"], c["cs"])
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}

	defer db.Close()

	var updateRecruitPointsDto structs.UpdateRecruitPointsDto
	err = json.NewDecoder(r.Body).Decode(&updateRecruitPointsDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	recruitingProfile := managers.GetRecruitingProfileByTeamId(db, updateRecruitPointsDto.ProfileId)

	if recruitingProfile.ScholarshipsAvailable == 0 {
		fmt.Printf("\nTeamId: " + updateRecruitPointsDto.ProfileId + " does not have any availabe scholarships")
		return
	}

	recruitingPointsProfile := managers.GetRecruitingPointsProfileByPlayerId(db,
		updateRecruitPointsDto.PlayerId,
		updateRecruitPointsDto.ProfileId)

	if recruitingPointsProfile.Scholarship == true {
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
	db, err := gorm.Open(c["db"], c["cs"])
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}

	defer db.Close()

	var updateRecruitPointsDto structs.UpdateRecruitPointsDto
	err = json.NewDecoder(r.Body).Decode(&updateRecruitPointsDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	recruitingProfile := managers.GetRecruitingProfileByTeamId(db, updateRecruitPointsDto.ProfileId)

	recruitingPointsProfile := managers.GetRecruitingPointsProfileByPlayerId(db,
		updateRecruitPointsDto.PlayerId,
		updateRecruitPointsDto.ProfileId)

	if recruitingPointsProfile.Scholarship == false {
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
