package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

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

	var wg sync.WaitGroup
	wg.Add(5)
	var (
		warRoom         structs.NBAWarRoom
		draftees        []structs.NBADraftee
		allNBATeams     []structs.NBATeam
		draftPicks      [2][]structs.DraftPick
		allCollegeTeams []structs.Team
	)

	go func() {
		defer wg.Done()
		warRoom = managers.GetNBAWarRoomByTeamID(teamID)
	}()

	go func() {
		defer wg.Done()
		draftees = managers.GetNBADrafteesForDraftPage()
	}()

	// GetAllNFLTeams
	go func() {
		defer wg.Done()
		allNBATeams = managers.GetOnlyNBATeams()
	}()

	// GetAllCurrentSeasonDraftPicksForDraftRoom
	go func() {
		defer wg.Done()
		draftPicks = managers.GetAllCurrentSeasonDraftPicks()
	}()

	// GetAllCollegeTeams
	go func() {
		defer wg.Done()
		allCollegeTeams = managers.GetAllActiveCollegeTeams()
	}()

	// Wait for all goroutines to complete
	wg.Wait()

	res := structs.NBADraftPageResponse{
		WarRoom:          warRoom,
		DraftablePlayers: draftees,
		NBATeams:         allNBATeams,
		AllDraftPicks:    draftPicks,
		CollegeTeams:     allCollegeTeams,
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

func ExportDraftedPicks(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	var draftPickDTO structs.ExportDraftPicksDTO
	err := json.NewDecoder(r.Body).Decode(&draftPickDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	saveComplete := managers.ExportDraftedPlayers(draftPickDTO.DraftPicks)

	json.NewEncoder(w).Encode(saveComplete)

	fmt.Fprintf(w, "Exported Players to new tables")
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

func ToggleDraftTime(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	managers.ToggleDraftTime()

	json.NewEncoder(w).Encode("Draft Time Changed")
}

func RunNBACombine(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	managers.NBACombineForDraft()

	json.NewEncoder(w).Encode("Draft Time Changed")
}
