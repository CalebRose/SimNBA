package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimNBA/managers"
	"github.com/gorilla/mux"
)

// CBBPlayerByNameAndAbbr - Get a college player record and share in the discord
func CBBPlayerByNameAndAbbr(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)

	firstName := vars["firstName"]
	lastName := vars["lastName"]
	abbr := vars["abbr"]
	if len(firstName) == 0 || len(lastName) == 0 || len(abbr) == 0 {
		panic("User did not provide enough information for call")
	}

	player := managers.GetCollegePlayerByNameAndAbbr(firstName, lastName, abbr)
	json.NewEncoder(w).Encode(player)
}

// NBAPlayerByNameAndAbbr - Get an NBA player record and share in the discord
func NBAPlayerByNameAndAbbr(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)

	firstName := vars["firstName"]
	lastName := vars["lastName"]
	abbr := vars["abbr"]
	if len(firstName) == 0 || len(lastName) == 0 || len(abbr) == 0 {
		panic("User did not provide enough information for call")
	}

	player := managers.GetNBAPlayerByNameAndAbbr(firstName, lastName, abbr)
	json.NewEncoder(w).Encode(player)
}

// GetCrootsByName - Get college croots with matching name
func GetCrootsByName(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)

	firstName := vars["firstName"]
	lastName := vars["lastName"]
	if len(firstName) == 0 || len(lastName) == 0 {
		panic("User did not provide enough information for call")
	}

	player := managers.GetCollegeRecruitByNameAndLocation(firstName, lastName)
	json.NewEncoder(w).Encode(player)
}

// GetCollegeTeamData - Get all season related data for a college team
func GetCollegeTeamData(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)

	teamId := vars["teamId"]
	if len(teamId) == 0 {
		panic("User did not provide enough information for call")
	}

	team := managers.GetCollegeTeamDataByID(teamId)
	json.NewEncoder(w).Encode(team)
}

// NBA Team Data
func GetNBATeamDataByID(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)

	teamId := vars["teamId"]
	if len(teamId) == 0 {
		panic("User did not provide enough information for call")
	}

	player := managers.GetNBATeamDataByID(teamId)
	json.NewEncoder(w).Encode(player)
}

// CollegeConferenceStandings
func CollegeConferenceStandings(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)

	conference := vars["conferenceID"]
	if len(conference) == 0 {
		panic("User did not provide enough information for call")
	}

	data := managers.GetCollegeConferenceStandingsByConference(conference)
	json.NewEncoder(w).Encode(data)
}

// Get NBAConferenceStandings
func NBAConferenceStandings(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)

	conf := vars["conferenceID"]
	if len(conf) == 0 {
		panic("User did not provide enough information for call")
	}

	player := managers.GetNBAConferenceStandingsByConference(conf)
	json.NewEncoder(w).Encode(player)
}

// CollegeMatchesByConference
func CollegeMatchesByConference(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)

	conf := vars["conferenceID"]
	day := vars["day"]
	if len(conf) == 0 {
		panic("User did not provide enough information for call")
	}

	player := managers.GetCollegeMatchesByConfAndDay(conf, day)
	json.NewEncoder(w).Encode(player)
}
