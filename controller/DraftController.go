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
