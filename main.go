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
	"github.com/robfig/cron/v3"
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
	myRouter.HandleFunc("/admin/rankCroots", controller.RankCroots).Methods("GET")
	myRouter.HandleFunc("/admin/ai/fill/boards", controller.FillAIBoards).Methods("GET")
	myRouter.HandleFunc("/admin/ai/sync/boards", controller.SyncAIBoards).Methods("GET")
	myRouter.HandleFunc("/admin/recruiting/sync", controller.SyncRecruiting).Methods("GET")
	myRouter.HandleFunc("/admin/lock/recruiting", controller.LockRecruiting).Methods("GET")
	myRouter.HandleFunc("/admin/recruit/create", controller.CreateRecruit).Methods("POST")
	myRouter.HandleFunc("/admin/ai/gameplans/", controller.SetAIGameplans).Methods("GET")
	myRouter.HandleFunc("/admin/results/import/", controller.ImportMatchResults).Methods("POST")
	myRouter.HandleFunc("/admin/show/results", controller.ShowGames).Methods("GET")
	// myRouter.HandleFunc("/admin/show/b", controller.ShowBGames).Methods("GET")
	// myRouter.HandleFunc("/admin/regress/a", controller.RegressAGamesByOneWeek).Methods("GET")
	// myRouter.HandleFunc("/admin/regress/b", controller.RegressBGamesByOneWeek).Methods("GET")
	myRouter.HandleFunc("/admin/week/sync", controller.SyncToNextWeek).Methods("GET")
	myRouter.HandleFunc("/admin/sync/contract/values", controller.SyncContractValues).Methods("GET")
	myRouter.HandleFunc("/simbba/matches/simulation", controller.GetMatchesForSimulation).Methods("GET")

	// myRouter.HandleFunc("/admin/generateGlobalPlayers", controller.GenerateGlobalPlayerRecords).Methods("GET")
	// myRouter.HandleFunc("/admin/generate/gameplans", controller.GenerateGameplans).Methods("GET")
	// myRouter.HandleFunc("/admin/generate/warrooms", controller.GenerateDraftWarRooms).Methods("GET")
	// myRouter.HandleFunc("/admin/generate/draft/grades", controller.GenerateDraftGrades).Methods("GET")
	// myRouter.HandleFunc("/admin/generate/draft/rounds", controller.GeneratePredictionRound).Methods("GET")
	// myRouter.HandleFunc("/admin/migrate/data", controller.MigratePlayers).Methods("GET")
	// myRouter.HandleFunc("/admin/migrate/progress", controller.ProgressPlayers).Methods("GET")
	// myRouter.HandleFunc("/admin/migrate/new/teams", controller.ImportNewTeams).Methods("GET")
	// myRouter.HandleFunc("/admin/migrate/nba/players", controller.MigrateNBAPlayersToTables).Methods("GET")
	// myRouter.HandleFunc("/admin/progress/nba/players", controller.ProgressNBAPlayers).Methods("GET")
	// myRouter.HandleFunc("/admin/clean/nba/tables", controller.CleanNBAPlayerTables).Methods("GET")

	// Capsheet Controls
	myRouter.HandleFunc("/nba/capsheet/generate", controller.GenerateCapsheets).Methods("GET")
	myRouter.HandleFunc("/nba/contracts/get/value", controller.CalculateContracts).Methods("GET")

	// Draft Controls
	// myRouter.HandleFunc("/nba/draft/conduct/lottery", controller.ConductDraftLottery).Methods("GET")
	myRouter.HandleFunc("/nba/draft/export/picks", controller.ExportDraftedPicks).Methods("POST")
	myRouter.HandleFunc("/nba/draft/page/{teamID}", controller.GetDraftPageData).Methods("GET")
	myRouter.HandleFunc("/nba/draft/time/change", controller.ToggleDraftTime).Methods("GET")
	myRouter.HandleFunc("/nba/draft/create/scoutprofile", controller.AddPlayerToScoutBoard).Methods("POST")
	myRouter.HandleFunc("/nba/draft/reveal/attribute", controller.RevealScoutingAttribute).Methods("POST")
	// myRouter.HandleFunc("/recruit/revokeScholarshipFromRecruit", controller.RevokeScholarshipFromRecruit).Methods("PUT")
	myRouter.HandleFunc("/nba/draft/remove/{id}", controller.RemovePlayerFromScoutBoard).Methods("GET")
	myRouter.HandleFunc("/nba/draft/scout/{id}", controller.GetScoutingDataByDraftee).Methods("GET")
	// myRouter.HandleFunc("/nba/draft/player/", controller.SaveRecruitingBoard).Methods("POST")

	// Exports
	myRouter.HandleFunc("/export/cbb/preseason", controller.ExportCBBPreseasonRanks).Methods("GET")
	myRouter.HandleFunc("/export/cbb/team/{teamID}", controller.ExportCBBRosterToCSV).Methods("GET")
	myRouter.HandleFunc("/export/nba/team/{teamID}", controller.ExportNBARosterToCSV).Methods("GET")

	// Free Agency Controls
	myRouter.HandleFunc("/nba/freeagency/available/{teamID}", controller.FreeAgencyAvailablePlayers).Methods("GET")
	myRouter.HandleFunc("/nba/freeagency/create/offer", controller.CreateFreeAgencyOffer).Methods("POST")
	myRouter.HandleFunc("/nba/freeagency/cancel/offer", controller.CancelFreeAgencyOffer).Methods("POST")
	myRouter.HandleFunc("/nba/freeagency/create/waiver", controller.CreateWaiverOffer).Methods("POST")
	myRouter.HandleFunc("/nba/freeagency/cancel/waiver", controller.CancelWaiverOffer).Methods("POST")
	myRouter.HandleFunc("/nba/freeagency/extensions/temp", controller.ExtendPlayers).Methods("GET")
	myRouter.HandleFunc("/nba/freeagency/sync/round", controller.SyncFreeAgencyOffers).Methods("GET")

	// Gameplan controls
	myRouter.HandleFunc("/cbb/gameplans/{teamId}", controller.GetGameplansByTeamId).Methods("GET")
	myRouter.HandleFunc("/nba/gameplans/{teamId}", controller.GetNBAGameplanByTeamId).Methods("GET")
	myRouter.HandleFunc("/cbb/gameplans/update", controller.UpdateGameplan).Methods("POST")
	myRouter.HandleFunc("/nba/gameplans/update", controller.UpdateNBAGameplan).Methods("POST")

	// Generation Controls
	// myRouter.HandleFunc("/admin/generateCoaches", controller.GenerateCoaches).Methods("GET")
	// myRouter.HandleFunc("/admin/generateTeam", controller.GeneratePlayers).Methods("GET")
	// myRouter.HandleFunc("/admin/generateCroots", controller.GenerateCroots).Methods("GET")
	// myRouter.HandleFunc("/admin/generate/international", controller.GenerateInternationalPlayers).Methods("GET")
	// myRouter.HandleFunc("/admin/allocate/international/rosters", controller.GenerateInternationalRoster).Methods("GET")
	// myRouter.HandleFunc("/admin/fix/height", controller.FixHeight).Methods("GET")
	// myRouter.HandleFunc("/generate/playtime/expectations", controller.GeneratePlaytimeExpectations).Methods("GET")

	// Import
	// myRouter.HandleFunc("/import/nba", controller.ImportNBAStandings).Methods("GET")
	// myRouter.HandleFunc("/import/cbb/games", controller.ImportCBBMatches).Methods("GET")
	// myRouter.HandleFunc("/import/nba/games", controller.ImportNBAMatches).Methods("GET")
	// myRouter.HandleFunc("/import/archetypes", controller.ImportArchetypes).Methods("GET")
	// myRouter.HandleFunc("/import/fa/preferences", controller.ImportFAPreferences).Methods("GET")
	// myRouter.HandleFunc("/import/positions", controller.ImportNewPositions).Methods("GET")
	// myRouter.HandleFunc("/import/ai/values", controller.MigrateNewAIRecruitingValues).Methods("GET")
	// myRouter.HandleFunc("/import/nba/personalities", controller.ImportNBAPersonalities).Methods("GET")
	// myRouter.HandleFunc("/import/nba/picks", controller.ImportDraftPicks).Methods("GET")
	// myRouter.HandleFunc("/migrate/remaining/croots", controller.MigrateRecruits).Methods("GET")

	// Match Controls
	myRouter.HandleFunc("/match/{matchId}", controller.GetMatchByMatchId).Methods("GET")
	myRouter.HandleFunc("/match/export/results/{seasonID}/{weekID}/{nbaWeekID}/{matchType}", controller.ExportMatchResults).Methods("GET")
	myRouter.HandleFunc("/match/result/cbb/{matchId}", controller.GetMatchResultByMatchID).Methods("GET")
	myRouter.HandleFunc("/match/result/nba/{matchId}", controller.GetNBAMatchResultByMatchID).Methods("GET")
	myRouter.HandleFunc("/match/team/{teamId}/season/{seasonId}", controller.GetMatchesByTeamIdAndSeasonId).Methods("GET")
	myRouter.HandleFunc("/match/week/{weekId}", controller.GetMatchesByTeamIdAndSeasonId).Methods("GET")
	myRouter.HandleFunc("/match/season/{seasonID}", controller.GetMatchesBySeasonID).Methods("GET")
	myRouter.HandleFunc("/match/team/upcoming/{teamId}/season/{seasonId}", controller.GetUpcomingMatchesByTeamIdAndSeasonId).Methods("GET")
	myRouter.HandleFunc("/cbb/match/data/{homeTeamAbbr}/{awayTeamAbbr}", controller.GetCBBMatchData).Methods("GET")
	myRouter.HandleFunc("/nba/match/data/{homeTeamID}/{awayTeamID}", controller.GetNBAMatchData).Methods("GET")
	myRouter.HandleFunc("/nba/match/team/{teamId}/season/{seasonId}", controller.GetNBAMatchesByTeamIdAndSeasonId).Methods("GET")

	// News Controls
	myRouter.HandleFunc("/cbb/news/all/", controller.GetAllCBBNewsInASeason).Methods("GET")
	myRouter.HandleFunc("/nba/news/all/", controller.GetAllNBANewsInASeason).Methods("GET")
	myRouter.HandleFunc("/news/feed/{league}/{teamID}/", controller.GetNewsFeed).Methods("GET")

	// Player Controls
	myRouter.HandleFunc("/player/AllPlayers", controller.AllCollegePlayers).Methods("GET")
	// myRouter.HandleFunc("/player/add/{firstname}/{lastname}", controller.NewPlayer).Methods("POST")
	myRouter.HandleFunc("/cbb/player/assign/redshirt/", controller.AssignRedshirtForCollegePlayer).Methods("POST")
	myRouter.HandleFunc("/player/GetPlayer/{playerId}", controller.PlayerById).Methods("GET")
	// myRouter.HandleFunc("/player/SetRedshirting/{playerId}", controller.SetRedshirtStatusByPlayerId).Methods("PUT")
	myRouter.HandleFunc("/players", controller.AllPlayers).Methods("GET")
	myRouter.HandleFunc("/players/{teamId}", controller.AllPlayersByTeamId).Methods("GET")
	myRouter.HandleFunc("/players/college", controller.AllCollegePlayers).Methods("GET")
	myRouter.HandleFunc("/players/college/recruits", controller.AllCollegeRecruits).Methods("GET")
	myRouter.HandleFunc("/collegeplayers/export/all", controller.ExportCollegePlayers).Methods("GET")
	myRouter.HandleFunc("/collegeplayers/check/declaration", controller.CheckDeclarationStatus).Methods("GET")
	myRouter.HandleFunc("/players/nba", controller.AllNBAPlayers).Methods("GET")
	myRouter.HandleFunc("/players/nba/freeAgents", controller.AllNBAFreeAgents).Methods("GET")
	myRouter.HandleFunc("/nba/players/{teamId}", controller.GetNBARosterByTeamID).Methods("GET")
	myRouter.HandleFunc("/nba/players/cut/{playerID}", controller.CutPlayerFromNBATeam).Methods("GET")
	myRouter.HandleFunc("/nba/players/place/gleague/{playerID}", controller.PlaceNBAPlayerInGLeague).Methods("GET")
	myRouter.HandleFunc("/nba/players/place/twoway/{playerID}", controller.AssignNBAPlayerAsTwoWay).Methods("GET")

	// Poll Controls
	myRouter.HandleFunc("/college/poll/create/", controller.CreatePollSubmission).Methods("POST")
	myRouter.HandleFunc("/college/poll/sync", controller.SyncCollegePoll).Methods("GET")
	myRouter.HandleFunc("/college/poll/official/season/{seasonID}", controller.GetOfficialPollsBySeasonID).Methods("GET")
	myRouter.HandleFunc("/college/poll/submission/{username}", controller.GetPollSubmission).Methods("GET")

	// Recruit Controls
	myRouter.HandleFunc("/recruiting/dashboard/{teamID}/", controller.GetRecruitingDataForOverviewPage).Methods("GET")
	myRouter.HandleFunc("/recruiting/profile/teamboard/{teamID}", controller.GetRecruitingDataForTeamBoardPage).Methods("GET")
	myRouter.HandleFunc("/recruiting/profile/all/", controller.GetAllRecruitingProfiles).Methods("GET")
	myRouter.HandleFunc("/recruiting/profile/determine/size/", controller.DetermineRecruitingClassSize).Methods("GET")
	myRouter.HandleFunc("/recruiting/class/{teamID}/", controller.GetRecruitingClassByTeamID).Methods("GET")
	myRouter.HandleFunc("/recruit/createRecruitingPointsProfile", controller.AddRecruitToBoard).Methods("POST")
	myRouter.HandleFunc("/recruit/allocatePoints", controller.AllocateRecruitingPointsForRecruit).Methods("PUT")
	myRouter.HandleFunc("/recruit/toggleScholarship", controller.SendScholarshipToRecruit).Methods("POST")
	// myRouter.HandleFunc("/recruit/revokeScholarshipFromRecruit", controller.RevokeScholarshipFromRecruit).Methods("PUT")
	myRouter.HandleFunc("/recruit/removeRecruit", controller.RemoveRecruitFromBoard).Methods("POST")
	myRouter.HandleFunc("/recruit/saveRecruitingBoard", controller.SaveRecruitingBoard).Methods("POST")
	myRouter.HandleFunc("/croots/export/all", controller.ExportCroots).Methods("GET")

	// Request Controls
	myRouter.HandleFunc("/requests/", controller.GetTeamRequests).Methods("GET")
	myRouter.HandleFunc("/requests/createTeamRequest", controller.CreateTeamRequest).Methods("POST")
	myRouter.HandleFunc("/requests/approveTeamRequest", controller.ApproveTeamRequest).Methods("PUT")
	myRouter.HandleFunc("/requests/rejectTeamRequest", controller.RejectTeamRequest).Methods("DELETE")
	myRouter.HandleFunc("/nba/requests/all/", controller.GetNBATeamRequests).Methods("GET")
	myRouter.HandleFunc("/nba/requests/create/", controller.CreateNBATeamRequest).Methods("POST")
	myRouter.HandleFunc("/nba/requests/approve/", controller.ApproveNBATeamRequest).Methods("POST")
	myRouter.HandleFunc("/nba/requests/reject/", controller.RejectNBATeamRequest).Methods("DELETE")
	myRouter.HandleFunc("/nba/requests/revoke/", controller.RemoveNBAUserFromNBATeam).Methods("POST")

	// Stats Controls
	myRouter.HandleFunc("/stats/export/{league}/{seasonID}/{weekID}/{matchType}/{viewType}/{playerView}", controller.ExportStats).Methods("GET")
	myRouter.HandleFunc("/stats/player/{playerId}", controller.GetPlayerStats).Methods("GET")
	myRouter.HandleFunc("/stats/player/{playerId}/match/{matchId}", controller.GetPlayerStatsByMatch).Methods("GET")
	myRouter.HandleFunc("/stats/player/{playerId}/season/{seasonId}", controller.GetPlayerStatsBySeason).Methods("GET")
	myRouter.HandleFunc("/stats/team/{teamId}/season/{seasonId}", controller.GetTeamStatsBySeason).Methods("GET")
	myRouter.HandleFunc("/stats/team/{teamId}/match/{matchId}", controller.GetCBBTeamStatsByMatch).Methods("GET")
	myRouter.HandleFunc("/stats/cbb/fix/player/stats", controller.FixPlayerStatsFromLastSeason).Methods("GET")
	myRouter.HandleFunc("/stats/cbb/{seasonID}/{weekID}/{matchType}/{viewType}", controller.GetCBBStatsPageData).Methods("GET")
	myRouter.HandleFunc("/stats/nba/{seasonID}/{weekID}/{matchType}/{viewType}", controller.GetNBAStatsPageData).Methods("GET")
	myRouter.HandleFunc("/stats/nba/team/{teamId}/match/{matchId}", controller.GetNBATeamStatsByMatch).Methods("GET")
	myRouter.HandleFunc("/stats/nba/match/{matchId}", controller.GetPlayerStatsByMatch).Methods("GET")

	// StandingsControls
	myRouter.HandleFunc("/standings/college/conf/{conferenceId}/{seasonId}", controller.GetConferenceStandingsByConferenceID).Methods("GET")
	myRouter.HandleFunc("/standings/nba/conf/{conferenceId}/{seasonId}", controller.GetNBAConferenceStandingsByConferenceID).Methods("GET")
	myRouter.HandleFunc("/standings/college/season/{seasonId}", controller.GetAllConferenceStandings).Methods("GET")
	myRouter.HandleFunc("/standings/nba/season/{seasonId}", controller.GetAllNBAConferenceStandings).Methods("GET")

	// Team Controls
	myRouter.HandleFunc("/team/{teamId}", controller.GetTeamByTeamID).Methods("GET")
	myRouter.HandleFunc("/team/nba/{teamId}", controller.GetNBATeamByTeamID).Methods("GET")
	myRouter.HandleFunc("/team/removeUserFromTeam/{teamId}", controller.RemoveUserFromTeam).Methods("PUT")
	myRouter.HandleFunc("/team/nba/removeUserFromTeam/{teamId}", controller.RemoveUserFromTeam).Methods("PUT")
	myRouter.HandleFunc("/teams", controller.AllTeams).Methods("GET")
	myRouter.HandleFunc("/teams/active", controller.AllActiveTeams).Methods("GET")
	myRouter.HandleFunc("/teams/active/college", controller.AllActiveCollegeTeams).Methods("GET")
	myRouter.HandleFunc("/teams/college/available", controller.AllAvailableTeams).Methods("GET")
	myRouter.HandleFunc("/teams/assign/ratings", controller.SyncTeamRatings).Methods("GET")
	myRouter.HandleFunc("/teams/nba/assign/ratings", controller.SyncNBATeamRatings).Methods("GET")
	myRouter.HandleFunc("/teams/coached", controller.AllCoachedTeams).Methods("GET")
	myRouter.HandleFunc("/teams/college", controller.AllCollegeTeams).Methods("GET")
	myRouter.HandleFunc("/teams/nba", controller.AllNBATeams).Methods("GET")
	myRouter.HandleFunc("/teams/isl", controller.AllISLTeams).Methods("GET")
	myRouter.HandleFunc("/teams/pro", controller.AllProfessionalTeams).Methods("GET")

	// Trade Controls
	myRouter.HandleFunc("/trades/nba/all/accepted", controller.GetAllAcceptedTrades).Methods("GET")
	myRouter.HandleFunc("/trades/nba/all/rejected", controller.GetAllRejectedTrades).Methods("GET")
	myRouter.HandleFunc("/trades/nba/block/{teamID}", controller.GetNBATradeBlockDataByTeamID).Methods("GET")
	myRouter.HandleFunc("/trades/nba/place/block/{playerID}", controller.PlaceNBAPlayerOnTradeBlock).Methods("GET")
	myRouter.HandleFunc("/trades/nba/preferences/update", controller.UpdateTradePreferences).Methods("POST")
	myRouter.HandleFunc("/trades/nba/create/proposal", controller.CreateNBATradeProposal).Methods("POST")
	myRouter.HandleFunc("/trades/nba/proposal/accept/{proposalID}", controller.AcceptTradeOffer).Methods("GET")
	myRouter.HandleFunc("/trades/nba/proposal/reject/{proposalID}", controller.RejectTradeOffer).Methods("GET")
	myRouter.HandleFunc("/trades/nba/proposal/cancel/{proposalID}", controller.CancelTradeOffer).Methods("GET")
	myRouter.HandleFunc("/admin/trades/accept/sync/{proposalID}", controller.SyncAcceptedTrade).Methods("GET")
	myRouter.HandleFunc("/admin/trades/veto/sync/{proposalID}", controller.VetoAcceptedTrade).Methods("GET")
	myRouter.HandleFunc("/admin/trades/cleanup", controller.CleanUpRejectedTrades).Methods("GET")

	// Timestamp Controls
	myRouter.HandleFunc("/simbba/get/timestamp", controller.GetCurrentTimestamp).Methods("GET")
	myRouter.HandleFunc("/cbb/easter/egg/collude/", controller.CollusionButton).Methods("POST")

	// Discord Controls
	myRouter.HandleFunc("/dis/cbb/player/{firstName}/{lastName}/{abbr}", controller.CBBPlayerByNameAndAbbr).Methods("GET")
	myRouter.HandleFunc("/dis/nba/player/{firstName}/{lastName}/{abbr}", controller.NBAPlayerByNameAndAbbr).Methods("GET")
	myRouter.HandleFunc("/dis/cbb/croot/{firstName}/{lastName}", controller.GetCrootsByName).Methods("GET")
	myRouter.HandleFunc("/dis/cbb/team/{teamId}", controller.GetCollegeTeamData).Methods("GET")
	myRouter.HandleFunc("/dis/nba/team/{teamId}", controller.GetNBATeamDataByID).Methods("GET")
	myRouter.HandleFunc("/dis/cbb/conf/standings/{conferenceID}", controller.CollegeConferenceStandings).Methods("GET")
	myRouter.HandleFunc("/dis/nba/conf/standings/{conferenceID}", controller.NBAConferenceStandings).Methods("GET")
	myRouter.HandleFunc("/dis/cbb/conf/matches/{conferenceID}/{day}", controller.CollegeMatchesByConference).Methods("GET")
	handler := cors.AllowAll().Handler(myRouter)

	log.Fatal(http.ListenAndServe(":8081", handler))
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World.")
}

func handleCron() {
	go func() {
		c := cron.New()
		c.AddFunc("0 16 * * 3", controller.SyncRecruitingViaCron)
		c.AddFunc("0 6 * * 4,6", controller.SyncAIBoardsViaCron)
		c.AddFunc("0 20 * * 1,3,5,6", controller.ShowGamesViaCron)
		c.AddFunc("0 10 * * 4", controller.FillAIBoardsViaCron)
		c.AddFunc("0 12 * * 0", controller.SyncToNextWeekViaCron)
		c.AddFunc("0 12 * * 2", controller.SyncFreeAgencyOffersViaCron)
		c.Start()
	}()
}

func main() {
	InitialMigration()
	fmt.Println("Database initialized.")

	fmt.Println("Loading cron...")
	handleCron()

	fmt.Println("Loading Requests...")
	handleRequests()
	fmt.Println("Application Running")
}
