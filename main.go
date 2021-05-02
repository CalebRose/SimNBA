package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/CalebRose/SimNBA/config"
	"github.com/CalebRose/SimNBA/controller"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/rs/cors"
)

var db *gorm.DB
var err error
var c = config.Config()

func InitialMigration() {
	// 23.252.52.222
	// 68.66.216.54
	db, err := gorm.Open(c["db"], c["cs"])
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}

	defer db.Close()

	db.AutoMigrate(&structs.CollegeWeek{})
	db.AutoMigrate(&structs.Gameplan{})
	db.AutoMigrate(&structs.Match{})
	db.AutoMigrate(&structs.NBAWeek{})
	db.AutoMigrate(&structs.Player{})
	db.AutoMigrate(&structs.PlayerStats{})
	db.AutoMigrate(&structs.RecruitingPoints{})
	db.AutoMigrate(&structs.RecruitingProfile{})
	db.AutoMigrate(&structs.Request{})
	db.AutoMigrate(&structs.Season{})
	db.AutoMigrate(&structs.Team{})
	db.AutoMigrate(&structs.TeamStats{})
	db.AutoMigrate(&structs.Timestamp{})
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", helloWorld).Methods("GET")
	// Match Controls
	myRouter.HandleFunc("/match/{matchId}", controller.GetMatchByMatchId).Methods("GET")
	myRouter.HandleFunc("/match/team/{teamId}/season/{seasonId}", controller.GetMatchesByTeamIdAndSeasonId).Methods("GET")
	myRouter.HandleFunc("/match/week/{weekId}", controller.GetMatchesByTeamIdAndSeasonId).Methods("GET")
	myRouter.HandleFunc("/match/team/upcoming/{teamId}/season/{seasonId}", controller.GetUpcomingMatchesByTeamIdAndSeasonId).Methods("GET")
	// Player Controls
	myRouter.HandleFunc("/player/add/{firstname}/{lastname}", controller.NewPlayer).Methods("POST")
	myRouter.HandleFunc("/player/GetPlayer/{playerId}", controller.PlayerById).Methods("GET")
	myRouter.HandleFunc("/player/SetRedshirting/{playerId}", controller.PlayerById).Methods("PUT")
	myRouter.HandleFunc("/players", controller.AllPlayers).Methods("GET")
	myRouter.HandleFunc("/players/{teamId}", controller.AllPlayersByTeamId).Methods("GET")
	myRouter.HandleFunc("/players/college", controller.AllCollegePlayers).Methods("GET")
	myRouter.HandleFunc("/players/college/recruits", controller.AllCollegeRecruits).Methods("GET")
	myRouter.HandleFunc("/players/nba", controller.AllNBAPlayers).Methods("GET")
	myRouter.HandleFunc("/players/nba/freeAgents", controller.AllNBAFreeAgents).Methods("GET")
	// Recruit Controls
	myRouter.HandleFunc("/recruit/croots/{profileId}", controller.AllRecruitsByProfileID).Methods("GET")
	myRouter.HandleFunc("/recruit/profile/{teamId}", controller.RecruitingProfileByTeamID).Methods("GET")
	myRouter.HandleFunc("/recruit/createRecruitingPointsProfile", controller.CreateRecruitingPointsProfileForRecruit).Methods("POST")
	myRouter.HandleFunc("/recruit/allocatePoints", controller.AllocateRecruitingPointsForRecruit).Methods("PUT")
	myRouter.HandleFunc("/recruit/sendScholarshipToRecruit", controller.SendScholarshipToRecruit).Methods("PUT")
	myRouter.HandleFunc("/recruit/revokeScholarshipFromRecruit", controller.RevokeScholarshipFromRecruit).Methods("PUT")
	// Request Controls
	myRouter.HandleFunc("/requests/", controller.GetTeamRequests).Methods("GET")
	myRouter.HandleFunc("/requests/createTeamRequest", controller.CreateTeamRequest).Methods("POST")
	myRouter.HandleFunc("/requests/approveTeamRequest", controller.ApproveTeamRequest).Methods("PUT")
	myRouter.HandleFunc("/requests/rejectTeamRequest", controller.RejectTeamRequest).Methods("DELETE")
	// Stats Controls
	myRouter.HandleFunc("/stats/player/{playerId}", controller.GetPlayerStats).Methods("GET")
	myRouter.HandleFunc("/stats/player/{playerId}/match/{matchId}", controller.GetPlayerStatsByMatch).Methods("GET")
	myRouter.HandleFunc("/stats/player/{playerId}/season/{seasonId}", controller.GetPlayerStatsBySeason).Methods("GET")
	myRouter.HandleFunc("/stats/team/{teamId}/season/{seasonId}", controller.GetTeamStatsBySeason).Methods("GET")
	myRouter.HandleFunc("/stats/team/{teamId}/match/{matchId}", controller.GetTeamStatsByMatch).Methods("GET")
	// Team Controls
	myRouter.HandleFunc("/team/{teamId}", controller.GetTeamByTeamID).Methods("GET")
	myRouter.HandleFunc("/team/removeUserFromTeam/{teamId}", controller.RemoveUserFromTeam).Methods("PUT")
	myRouter.HandleFunc("/teams", controller.AllTeams).Methods("GET")
	myRouter.HandleFunc("/teams/active", controller.AllActiveTeams).Methods("GET")
	myRouter.HandleFunc("/teams/available", controller.AllAvailableTeams).Methods("GET")
	myRouter.HandleFunc("/teams/coached", controller.AllCoachedTeams).Methods("GET")
	myRouter.HandleFunc("/teams/college", controller.AllCollegeTeams).Methods("GET")
	myRouter.HandleFunc("/teams/nba", controller.AllNBATeams).Methods("GET")
	// Timestamp Controls
	myRouter.HandleFunc("/timestamp", controller.GetCurrentTimestamp).Methods("GET")

	handler := cors.AllowAll().Handler(myRouter)

	log.Fatal(http.ListenAndServe(":8081", handler))
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World.")
}

func main() {
	fmt.Println("GORM initiation")

	InitialMigration()

	handleRequests()

	fmt.Println("Application Running")
}
