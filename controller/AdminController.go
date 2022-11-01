package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/CalebRose/SimNBA/managers"
)

func GeneratePlayers(w http.ResponseWriter, r *http.Request) {
	managers.GenerateNewTeams()
	json.NewEncoder(w).Encode("GENERATION COMPLETE")
}

func GenerateCroots(w http.ResponseWriter, r *http.Request) {
	managers.GenerateCroots()
	json.NewEncoder(w).Encode("GENERATION COMPLETE")
}

func RankCroots(w http.ResponseWriter, r *http.Request) {
	managers.AssignAllRecruitRanks()
	json.NewEncoder(w).Encode("Ranking COMPLETE")
}

func GenerateGlobalPlayerRecords(w http.ResponseWriter, r *http.Request) {
	managers.GenerateGlobalPlayerRecords()
	json.NewEncoder(w).Encode("GENERATION COMPLETE")
}

func MigratePlayers(w http.ResponseWriter, r *http.Request) {
	managers.MigrateOldPlayerDataToNewTables()
	json.NewEncoder(w).Encode("DONE!")
}

func ProgressPlayers(w http.ResponseWriter, r *http.Request) {
	managers.ProgressionMain()
	json.NewEncoder(w).Encode("DONE!")
}

func FillAIBoards(w http.ResponseWriter, r *http.Request) {
	managers.FillAIRecruitingBoards()
	json.NewEncoder(w).Encode("AI recruiting boards complete!")
}

func LockRecruiting(w http.ResponseWriter, r *http.Request) {
	managers.LockRecruiting()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Recruiting Locked")
}

func SyncAIBoards(w http.ResponseWriter, r *http.Request) {
	managers.AllocatePointsToAIBoards()
	json.NewEncoder(w).Encode("AI recruiting boards Synced!")
}

func SyncRecruiting(w http.ResponseWriter, r *http.Request) {
	ts := managers.GetTimestamp()
	managers.SyncRecruiting(ts)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Recruiting Synced for Week " + strconv.Itoa(ts.CollegeWeek))
}
