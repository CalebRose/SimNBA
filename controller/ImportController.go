package controller

import (
	"net/http"

	"github.com/CalebRose/SimNBA/managers"
)

func ImportNBATeamsAndArenas(w http.ResponseWriter, r *http.Request) {
	managers.ImportNBATeamsAndArenas()
}

func ImportNewPositions(w http.ResponseWriter, r *http.Request) {
	managers.ImportNewPositions()
}

func ImportNBAStandings(w http.ResponseWriter, r *http.Request) {
	managers.ImportNBAStandings()
}

func MigrateRecruits(w http.ResponseWriter, r *http.Request) {
	managers.MigrateRecruits()
}

func MigrateNewAIRecruitingValues(w http.ResponseWriter, r *http.Request) {
	managers.MigrateNewAIRecruitingValues()
}

func ImportNBAPersonalities(w http.ResponseWriter, r *http.Request) {
	managers.ImportNBAPersonalities()
}

func ImportCBBMatches(w http.ResponseWriter, r *http.Request) {
	managers.ImportCBBGames()
}

func ImportNBAMatches(w http.ResponseWriter, r *http.Request) {
	managers.ImportNBAGames()
}

func ImportDraftPicks(w http.ResponseWriter, r *http.Request) {
	managers.ImportDraftPicks()
}
