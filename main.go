package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/CalebRose/SimNBA/controller"
	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/nelkinda/health-go"
	"github.com/nelkinda/health-go/checks/sendgrid"
	"github.com/rs/cors"
)

func InitialMigration() {
	initiate := dbprovider.GetInstance().InitDatabase()
	if !initiate {
		log.Println("Initiate pool failure... Ending application")
		os.Exit(1)
	}
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", helloWorld).Methods("GET")
	// Health Controls
	HealthCheck := health.New(
		health.Health{
			Version:   "1",
			ReleaseID: "0.0.7-SNAPSHOT",
		},
		sendgrid.Health(),
	)
	myRouter.HandleFunc("/health", HealthCheck.Handler).Methods("GET")

	// Admin Controls
	// myRouter.HandleFunc("/admin/generateTeam", controller.GeneratePlayers).Methods("GET")
	// myRouter.HandleFunc("/admin/generateCroots", controller.GenerateCroots).Methods("GET")
	myRouter.HandleFunc("/admin/rankCroots", controller.RankCroots).Methods("GET")
	myRouter.HandleFunc("/admin/ai/fill/boards", controller.FillAIBoards).Methods("GET")
	myRouter.HandleFunc("/admin/ai/sync/boards", controller.SyncAIBoards).Methods("GET")
	myRouter.HandleFunc("/admin/recruiting/sync", controller.SyncRecruiting).Methods("GET")
	myRouter.HandleFunc("/admin/lock/recruiting", controller.LockRecruiting).Methods("GET")
	myRouter.HandleFunc("/admin/recruit/create", controller.CreateRecruit).Methods("POST")

	// myRouter.HandleFunc("/admin/generateGlobalPlayers", controller.GenerateGlobalPlayerRecords).Methods("GET")
	// myRouter.HandleFunc("/admin/migrate/data", controller.MigratePlayers).Methods("GET")
	// myRouter.HandleFunc("/admin/migrate/progress", controller.ProgressPlayers).Methods("GET")

	// Gameplan controls
	myRouter.HandleFunc("/gameplans/{teamId}", controller.GetGameplansByTeamId).Methods("GET")
	myRouter.HandleFunc("/gameplans/update", controller.UpdateGameplan).Methods("PUT")

	// Match Controls
	myRouter.HandleFunc("/match/{matchId}", controller.GetMatchByMatchId).Methods("GET")
	myRouter.HandleFunc("/match/team/{teamId}/season/{seasonId}", controller.GetMatchesByTeamIdAndSeasonId).Methods("GET")
	myRouter.HandleFunc("/match/week/{weekId}", controller.GetMatchesByTeamIdAndSeasonId).Methods("GET")
	myRouter.HandleFunc("/match/team/upcoming/{teamId}/season/{seasonId}", controller.GetUpcomingMatchesByTeamIdAndSeasonId).Methods("GET")
	// Player Controls
	myRouter.HandleFunc("/player/AllPlayers", controller.AllCollegePlayers).Methods("GET")
	myRouter.HandleFunc("/player/add/{firstname}/{lastname}", controller.NewPlayer).Methods("POST")
	myRouter.HandleFunc("/player/GetPlayer/{playerId}", controller.PlayerById).Methods("GET")
	myRouter.HandleFunc("/player/SetRedshirting/{playerId}", controller.SetRedshirtStatusByPlayerId).Methods("PUT")
	myRouter.HandleFunc("/players", controller.AllPlayers).Methods("GET")
	myRouter.HandleFunc("/players/{teamId}", controller.AllPlayersByTeamId).Methods("GET")
	myRouter.HandleFunc("/players/college", controller.AllCollegePlayers).Methods("GET")
	myRouter.HandleFunc("/players/college/recruits", controller.AllCollegeRecruits).Methods("GET")
	myRouter.HandleFunc("/players/nba", controller.AllNBAPlayers).Methods("GET")
	myRouter.HandleFunc("/players/nba/freeAgents", controller.AllNBAFreeAgents).Methods("GET")
	// Recruit Controls
	myRouter.HandleFunc("/recruiting/dashboard/{teamID}/", controller.GetRecruitingDataForOverviewPage).Methods("GET")
	myRouter.HandleFunc("/recruiting/profile/teamboard/{teamID}", controller.GetRecruitingDataForTeamBoardPage).Methods("GET")
	myRouter.HandleFunc("/recruiting/profile/all/", controller.GetAllRecruitingProfiles).Methods("GET")

	myRouter.HandleFunc("/recruit/createRecruitingPointsProfile", controller.AddRecruitToBoard).Methods("POST")
	myRouter.HandleFunc("/recruit/allocatePoints", controller.AllocateRecruitingPointsForRecruit).Methods("PUT")
	myRouter.HandleFunc("/recruit/toggleScholarship", controller.SendScholarshipToRecruit).Methods("POST")
	// myRouter.HandleFunc("/recruit/revokeScholarshipFromRecruit", controller.RevokeScholarshipFromRecruit).Methods("PUT")
	myRouter.HandleFunc("/recruit/removeRecruit", controller.RemoveRecruitFromBoard).Methods("POST")
	myRouter.HandleFunc("/recruit/saveRecruitingBoard", controller.SaveRecruitingBoard).Methods("POST")

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

	// StandingsControls
	myRouter.HandleFunc("/standings/college/conf/{conferenceId}/{seasonId}", controller.GetConferenceStandingsByConferenceID).Methods("GET")

	// Team Controls
	myRouter.HandleFunc("/team/{teamId}", controller.GetTeamByTeamID).Methods("GET")
	myRouter.HandleFunc("/team/removeUserFromTeam/{teamId}", controller.RemoveUserFromTeam).Methods("PUT")
	myRouter.HandleFunc("/teams", controller.AllTeams).Methods("GET")
	myRouter.HandleFunc("/teams/active", controller.AllActiveTeams).Methods("GET")
	myRouter.HandleFunc("/teams/active/college", controller.AllActiveCollegeTeams).Methods("GET")
	myRouter.HandleFunc("/teams/active/nba", controller.AllActiveNBATeams).Methods("GET")
	myRouter.HandleFunc("/teams/college/available", controller.AllAvailableTeams).Methods("GET")
	myRouter.HandleFunc("/teams/assign/ratings", controller.SyncTeamRatings).Methods("GET")
	myRouter.HandleFunc("/teams/coached", controller.AllCoachedTeams).Methods("GET")
	myRouter.HandleFunc("/teams/college", controller.AllCollegeTeams).Methods("GET")
	myRouter.HandleFunc("/teams/nba", controller.AllNBATeams).Methods("GET")
	// Timestamp Controls
	myRouter.HandleFunc("/simbba/get/timestamp", controller.GetCurrentTimestamp).Methods("GET")

	handler := cors.AllowAll().Handler(myRouter)

	log.Fatal(http.ListenAndServe(":8081", handler))
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World.")
}

func main() {
	InitialMigration()
	fmt.Println("Database initialized.")

	handleRequests()

	fmt.Println("Application Running")
}
