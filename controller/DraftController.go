package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CalebRose/SimNBA/managers"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/gorilla/mux"
)

func ConductDraftLottery(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	managers.ConductDraftLottery()
}

func GenerateDraftGrades(w http.ResponseWriter, r *http.Request) {
	managers.GenerateDraftLetterGrades()
	fmt.Println(w, "Congrats, you generated the Letter Grades!")
}

func GeneratePredictionRound(w http.ResponseWriter, r *http.Request) {
	managers.DraftPredictionRound()
	fmt.Println(w, "Congrats, you generated the Round Predictions!")
}

func CheckDeclarationStatus(w http.ResponseWriter, r *http.Request) {
	managers.RunDeclarationsAlgorithm()
	fmt.Println(w, "Congrats, you generated the Letter Grades!")
}

func GetDraftPageData(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	teamID := vars["teamID"]
	if len(teamID) == 0 {
		panic("User did not provide TeamID")
	}
	// Get War Room
	// Get Scouting Profiles?
	// Get full list of draftable players

	warRoom := managers.GetNBAWarRoomByTeamID(teamID)
	draftees := managers.GetNBADrafteesForDraftPage()
	allNBATeams := managers.GetOnlyNBATeams()
	draftPicks := managers.GetAllCurrentSeasonDraftPicks()

	res := structs.NBADraftPageResponse{
		WarRoom:          warRoom,
		DraftablePlayers: draftees,
		NBATeams:         allNBATeams,
		AllDraftPicks:    draftPicks,
	}

	json.NewEncoder(w).Encode(res)
}

func AddPlayerToScoutBoard(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	var scoutProfileDto structs.ScoutingProfileDTO
	err := json.NewDecoder(r.Body).Decode(&scoutProfileDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	scoutingProfile := managers.CreateScoutingProfile(scoutProfileDto)

	json.NewEncoder(w).Encode(scoutingProfile)
}

func RevealScoutingAttribute(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	var revealAttributeDTO structs.RevealAttributeDTO
	err := json.NewDecoder(r.Body).Decode(&revealAttributeDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	saveComplete := managers.RevealScoutingAttribute(revealAttributeDTO)

	json.NewEncoder(w).Encode(saveComplete)

	fmt.Fprintf(w, "New Scout Profile Created")
}

func RemovePlayerFromScoutBoard(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	id := vars["id"]
	if len(id) == 0 {
		panic("User did not provide scout profile id")
	}

	managers.RemovePlayerFromScoutBoard(id)

	json.NewEncoder(w).Encode("Removed Player From Scout Board")
}

func GetScoutingDataByDraftee(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	id := vars["id"]
	if len(id) == 0 {
		panic("User did not provide scout profile id")
	}

	data := managers.GetScoutingDataByPlayerID(id)

	json.NewEncoder(w).Encode(data)
}
