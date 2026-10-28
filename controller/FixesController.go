package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimNBA/managers"
)

func FixRecruitingProfileModifiers(w http.ResponseWriter, r *http.Request) {
	managers.FixExistingModifiersForRecruits()
}

func FixCBBRankings(w http.ResponseWriter, r *http.Request) {
	ts := managers.GetTimestamp()
	managers.UpdateCollegeRankings(ts)
	json.NewEncoder(w).Encode("DONE")
}

func FixTeamSeasonStats(w http.ResponseWriter, r *http.Request) {
	managers.FixTeamSeasonStats()
	json.NewEncoder(w).Encode("DONE")
}
