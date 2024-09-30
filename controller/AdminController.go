package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimNBA/managers"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/gorilla/mux"
)

func EnableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
}

// RankCroots -- Assigns recruiting rankings for all current CBB recruits
func RankCroots(w http.ResponseWriter, r *http.Request) {
	managers.AssignAllRecruitRanks()
	json.NewEncoder(w).Encode("Ranking COMPLETE")
}

// func MigratePlayers(w http.ResponseWriter, r *http.Request) {
// 	managers.MigrateOldPlayerDataToNewTables()
// 	json.NewEncoder(w).Encode("DONE!")
// }

func ProgressPlayers(w http.ResponseWriter, r *http.Request) {
	// managers.ProgressionMain()
	// Sync Promises
	// managers.SyncPromises()
	// Enter Transfer Portal
	// managers.EnterTheTransferPortal()
	json.NewEncoder(w).Encode("DONE!")
}

func ImportNewTeams(w http.ResponseWriter, r *http.Request) {
	managers.ImportNewTeams()
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
	managers.ResetAIBoardsForCompletedTeams()
	managers.AllocatePointsToAIBoards()
	json.NewEncoder(w).Encode("AI recruiting boards Synced!")
}

func SyncRecruiting(w http.ResponseWriter, r *http.Request) {
	managers.SyncRecruiting()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Recruiting Synced for Week.")
}

func SyncTransferPortal(w http.ResponseWriter, r *http.Request) {
	managers.SyncTransferPortal()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("TP Synced for Week.")
}

func ImportMatchResults(w http.ResponseWriter, r *http.Request) {
	var dto structs.ImportMatchResultsDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	managers.ImportMatchResultsToDB(dto)
}

func SyncToNextWeek(w http.ResponseWriter, r *http.Request) {
	managers.ResetCollegeStandingsRanks()
	managers.SyncToNextWeek()
	ts := managers.GetTimestamp()
	managers.SyncCollegePollSubmissionForCurrentWeek(uint(ts.CollegeWeek), ts.CollegeWeekID, ts.SeasonID)
	w.WriteHeader(http.StatusOK)
}

func ShowGames(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	managers.ShowGames()
	w.WriteHeader(http.StatusOK)

}

func RegressBGamesByOneWeek(w http.ResponseWriter, r *http.Request) {
	managers.RegressGames("B")
	w.WriteHeader(http.StatusOK)
}

func RegressAGamesByOneWeek(w http.ResponseWriter, r *http.Request) {
	managers.RegressGames("D")
	w.WriteHeader(http.StatusOK)
}

func GetAllCBBNewsInASeason(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	newsLogs := managers.GetAllCBBNewsLogs()

	json.NewEncoder(w).Encode(newsLogs)
}

func GetAllNBANewsInASeason(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	newsLogs := managers.GetAllNBANewsLogs()

	json.NewEncoder(w).Encode(newsLogs)
}

// Collusion Button
func CollusionButton(w http.ResponseWriter, r *http.Request) {
	var collusionButton structs.CollusionDto

	ts := managers.GetTimestamp()

	err := json.NewDecoder(r.Body).Decode(&collusionButton)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.CreateNewsLog("CBB", collusionButton.Message, "Collusion", 0, ts)
}

func ProgressNBAPlayers(w http.ResponseWriter, r *http.Request) {
	managers.ProgressNBAPlayers()
	json.NewEncoder(w).Encode("Migration Complete")
}

func SyncContractValues(w http.ResponseWriter, r *http.Request) {
	managers.SyncContractValues()
	w.WriteHeader(http.StatusOK)
}

func CleanNBAPlayerTables(w http.ResponseWriter, r *http.Request) {
	managers.ProgressContractsByOneYear()
	w.WriteHeader(http.StatusOK)
}

func ExportCBBPreseasonRanks(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	w.Header().Set("Content-Type", "text/csv")
	managers.ExportCBBPreseasonRanks(w)
}

func ExportCBBRosterToCSV(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	w.Header().Set("Content-Type", "text/csv")

	vars := mux.Vars(r)
	teamId := vars["teamID"]

	if len(teamId) == 0 {
		panic("User did not provide TeamID")
	}

	managers.ExportCBBRosterToCSV(teamId, w)
}

func ExportNBARosterToCSV(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	w.Header().Set("Content-Type", "text/csv")

	vars := mux.Vars(r)
	teamId := vars["teamID"]

	if len(teamId) == 0 {
		panic("User did not provide TeamID")
	}

	managers.ExportNBARosterToCSV(teamId, w)
}
