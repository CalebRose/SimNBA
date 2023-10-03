package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/CalebRose/SimNBA/managers"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/gorilla/mux"
)

func GetCBBMatchData(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	homeTeamAbbr := vars["homeTeamAbbr"]
	awayTeamAbbr := vars["awayTeamAbbr"]

	var response structs.CBBMatchDataResponse

	var homeTeam structs.Team
	var awayTeam structs.Team
	hTeamChan := make(chan structs.Team)
	aTeamChan := make(chan structs.Team)

	go func() {
		hg := managers.GetCBBTeamByAbbreviation(homeTeamAbbr)
		hTeamChan <- hg
	}()

	go func() {
		ag := managers.GetCBBTeamByAbbreviation(awayTeamAbbr)
		aTeamChan <- ag
	}()

	homeTeam = <-hTeamChan
	close(hTeamChan)
	awayTeam = <-aTeamChan
	close(aTeamChan)

	homeTeamID := strconv.Itoa(int(homeTeam.ID))
	awayTeamID := strconv.Itoa(int(awayTeam.ID))

	var homeTeamResponse structs.MatchTeamResponse
	var awayTeamResponse structs.MatchTeamResponse

	homeTeamResponse.Map(homeTeam)
	awayTeamResponse.Map(awayTeam)

	var homeTeamRoster []structs.CollegePlayer
	var awayTeamRoster []structs.CollegePlayer
	hRosterChan := make(chan []structs.CollegePlayer)
	aRosterChan := make(chan []structs.CollegePlayer)

	go func() {
		hg := managers.GetCollegePlayersByTeamId(homeTeamID)
		hRosterChan <- hg
	}()

	go func() {
		ag := managers.GetCollegePlayersByTeamId(awayTeamID)
		aRosterChan <- ag
	}()

	homeTeamRoster = <-hRosterChan
	close(hRosterChan)
	awayTeamRoster = <-aRosterChan
	close(aRosterChan)

	var homeTeamGameplan structs.Gameplan
	var awayTeamGameplan structs.Gameplan
	hGameplanChan := make(chan structs.Gameplan)
	aGameplanChan := make(chan structs.Gameplan)

	go func() {
		hg := managers.GetGameplansByTeam(homeTeamID)
		hGameplanChan <- hg
	}()

	go func() {
		ag := managers.GetGameplansByTeam(awayTeamID)
		aGameplanChan <- ag
	}()

	homeTeamGameplan = <-hGameplanChan
	close(hGameplanChan)
	awayTeamGameplan = <-aGameplanChan
	close(aGameplanChan)

	response.AssignHomeTeam(homeTeamResponse, homeTeamRoster, homeTeamGameplan)
	response.AssignAwayTeam(awayTeamResponse, awayTeamRoster, awayTeamGameplan)

	json.NewEncoder(w).Encode(response)
}

func GetNBAMatchData(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	homeTeamID := vars["homeTeamID"]
	awayTeamID := vars["awayTeamID"]

	var response structs.NBAMatchDataResponse

	var homeTeam structs.NBATeam
	var awayTeam structs.NBATeam
	hTeamChan := make(chan structs.NBATeam)
	aTeamChan := make(chan structs.NBATeam)

	go func() {
		hg := managers.GetNBATeamByTeamID(homeTeamID)
		hTeamChan <- hg
	}()

	go func() {
		ag := managers.GetNBATeamByTeamID(awayTeamID)
		aTeamChan <- ag
	}()

	homeTeam = <-hTeamChan
	close(hTeamChan)
	awayTeam = <-aTeamChan
	close(aTeamChan)

	var homeTeamResponse structs.MatchTeamResponse
	var awayTeamResponse structs.MatchTeamResponse

	homeTeamResponse.MapNBATeam(homeTeam)
	awayTeamResponse.MapNBATeam(awayTeam)

	var homeTeamRoster []structs.NBAPlayer
	var awayTeamRoster []structs.NBAPlayer
	hRosterChan := make(chan []structs.NBAPlayer)
	aRosterChan := make(chan []structs.NBAPlayer)

	go func() {
		hg := managers.GetAllNBAPlayersByTeamID(homeTeamID)
		hRosterChan <- hg
	}()

	go func() {
		ag := managers.GetAllNBAPlayersByTeamID(awayTeamID)
		aRosterChan <- ag
	}()

	homeTeamRoster = <-hRosterChan
	close(hRosterChan)
	awayTeamRoster = <-aRosterChan
	close(aRosterChan)

	var homeTeamGameplan structs.NBAGameplan
	var awayTeamGameplan structs.NBAGameplan
	hGameplanChan := make(chan structs.NBAGameplan)
	aGameplanChan := make(chan structs.NBAGameplan)

	go func() {
		hg := managers.GetNBAGameplanByTeam(homeTeamID)
		hGameplanChan <- hg
	}()

	go func() {
		ag := managers.GetNBAGameplanByTeam(awayTeamID)
		aGameplanChan <- ag
	}()

	homeTeamGameplan = <-hGameplanChan
	close(hGameplanChan)
	awayTeamGameplan = <-aGameplanChan
	close(aGameplanChan)

	response.AssignHomeTeam(homeTeamResponse, homeTeamRoster, homeTeamGameplan)
	response.AssignAwayTeam(awayTeamResponse, awayTeamRoster, awayTeamGameplan)

	json.NewEncoder(w).Encode(response)
}
