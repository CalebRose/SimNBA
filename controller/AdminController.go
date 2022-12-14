package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/managers"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/gorilla/mux"
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
	managers.ResetAIBoardsForCompletedTeams()
	managers.AllocatePointsToAIBoards()
	json.NewEncoder(w).Encode("AI recruiting boards Synced!")
}

func SyncRecruiting(w http.ResponseWriter, r *http.Request) {
	ts := managers.GetTimestamp()
	managers.SyncRecruiting(ts)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Recruiting Synced for Week " + strconv.Itoa(ts.CollegeWeek))
}

func ImportMatchResults(w http.ResponseWriter, r *http.Request) {
	var dto structs.ImportMatchResultsDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	managers.ImportMatchResultsToDB(dto)
	w.WriteHeader(http.StatusOK)
}

func SyncToNextWeek(w http.ResponseWriter, r *http.Request) {
	managers.SyncToNextWeek()
	w.WriteHeader(http.StatusOK)

}

func ShowAGames(w http.ResponseWriter, r *http.Request) {
	managers.ShowAGames()
	w.WriteHeader(http.StatusOK)

}

func ShowBGames(w http.ResponseWriter, r *http.Request) {
	managers.ShowBGames()
	w.WriteHeader(http.StatusOK)

}

func GetAllNewsInASeason(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	seasonID := vars["seasonID"]

	newsLogs := managers.GetAllNewsLogs(seasonID)

	json.NewEncoder(w).Encode(newsLogs)
}

// Collusion Button
func CollusionButton(w http.ResponseWriter, r *http.Request) {
	db := dbprovider.GetInstance().GetDB()
	var collusionButton structs.CollusionDto

	ts := managers.GetTimestamp()

	err := json.NewDecoder(r.Body).Decode(&collusionButton)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newsLog := structs.NewsLog{
		WeekID:      uint(collusionButton.WeekID),
		SeasonID:    uint(ts.SeasonID),
		Week:        uint(ts.CollegeWeek),
		MessageType: "Collusion",
		Message:     collusionButton.Message,
	}

	db.Create(&newsLog)
}
