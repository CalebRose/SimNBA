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

	homeTeam := managers.GetCBBTeamByAbbreviation(homeTeamAbbr)
	awayTeam := managers.GetCBBTeamByAbbreviation(awayTeamAbbr)

	homeTeamID := strconv.Itoa(int(homeTeam.ID))
	awayTeamID := strconv.Itoa(int(awayTeam.ID))

	var homeTeamResponse structs.MatchTeamResponse
	var awayTeamResponse structs.MatchTeamResponse

	homeTeamResponse.Map(homeTeam)
	awayTeamResponse.Map(awayTeam)

	homeTeamRoster := managers.GetCollegePlayersByTeamId(homeTeamID)
	awayTeamRoster := managers.GetCollegePlayersByTeamId(awayTeamID)

	homeTeamGameplan := managers.GetGameplansByTeam(homeTeamID)
	awayTeamGameplan := managers.GetGameplansByTeam(awayTeamID)
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

	homeTeam := managers.GetNBATeamByTeamID(homeTeamID)
	awayTeam := managers.GetNBATeamByTeamID(awayTeamID)

	var homeTeamResponse structs.MatchTeamResponse
	var awayTeamResponse structs.MatchTeamResponse

	homeTeamResponse.MapNBATeam(homeTeam)
	awayTeamResponse.MapNBATeam(awayTeam)

	homeTeamRoster := managers.GetAllNBAPlayersByTeamID(homeTeamID)
	awayTeamRoster := managers.GetAllNBAPlayersByTeamID(awayTeamID)

	homeTeamGameplan := managers.GetNBAGameplanByTeam(homeTeamID)
	awayTeamGameplan := managers.GetNBAGameplanByTeam(awayTeamID)
	response.AssignHomeTeam(homeTeamResponse, homeTeamRoster, homeTeamGameplan)
	response.AssignAwayTeam(awayTeamResponse, awayTeamRoster, awayTeamGameplan)

	json.NewEncoder(w).Encode(response)
}
