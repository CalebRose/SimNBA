package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/CalebRose/SimNBA/controller"
	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/ws"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
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

func monitorDBForUpdates() {
	var ts structs.Timestamp
	for {
		currentTS := controller.GetUpdatedTimestamp()
		if currentTS.UpdatedAt.After(ts.UpdatedAt) {
			ts = currentTS
			err := ws.BroadcastTSUpdate(ts)
			if err != nil {
				log.Printf("Error broadcasting timestamp: %v", err)
			}
		}

		time.Sleep(60 * time.Second)
	}
}

func handleRequests() http.Handler {
	myRouter := mux.NewRouter().StrictSlash(true)
	// Health Controls
	HealthCheck := health.New(
		health.Health{
			Version:   "1",
			ReleaseID: "0.0.7-SNAPSHOT",
		},
		sendgrid.Health(),
	)
	apiRouter := myRouter.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/health", HealthCheck.Handler).Methods("GET")

	// Admin Controls
	apiRouter.HandleFunc("/admin/generate/ts/models/", controller.CreateTSModelsFile).Methods("GET")
	apiRouter.HandleFunc("/admin/rankCroots", controller.RankCroots).Methods("GET")
	apiRouter.HandleFunc("/admin/ai/fill/boards", controller.FillAIBoards).Methods("GET")
	apiRouter.HandleFunc("/admin/ai/sync/boards", controller.SyncAIBoards).Methods("GET")
	apiRouter.HandleFunc("/admin/recruiting/sync", controller.SyncRecruiting).Methods("GET")
	// apiRouter.HandleFunc("/admin/transferportal/sync", controller.SyncTransferPortal).Methods("GET")
	apiRouter.HandleFunc("/admin/lock/recruiting", controller.LockRecruiting).Methods("GET")
	apiRouter.HandleFunc("/admin/recruit/create", controller.CreateRecruit).Methods("POST")
	apiRouter.HandleFunc("/admin/ai/gameplans/", controller.SetAIGameplans).Methods("GET")
	apiRouter.HandleFunc("/admin/results/import/", controller.ImportMatchResults).Methods("POST")
	apiRouter.HandleFunc("/admin/show/results", controller.ShowGames).Methods("GET")
	// apiRouter.HandleFunc("/admin/show/b", controller.ShowBGames).Methods("GET")
	// apiRouter.HandleFunc("/admin/regress/a", controller.RegressAGamesByOneWeek).Methods("GET")
	// apiRouter.HandleFunc("/admin/regress/b", controller.RegressBGamesByOneWeek).Methods("GET")
	apiRouter.HandleFunc("/admin/week/sync", controller.SyncToNextWeek).Methods("GET")
	apiRouter.HandleFunc("/admin/sync/contract/values", controller.SyncContractValues).Methods("GET")
	apiRouter.HandleFunc("/simbba/matches/simulation", controller.GetMatchesForSimulation).Methods("GET")
	apiRouter.HandleFunc("/simcbb/user/gameplans/fix", controller.CheckUserGameplans).Methods("GET")

	// apiRouter.HandleFunc("/admin/generateGlobalPlayers", controller.GenerateGlobalPlayerRecords).Methods("GET")
	// apiRouter.HandleFunc("/admin/generate/gameplans", controller.GenerateGameplans).Methods("GET")
	// apiRouter.HandleFunc("/admin/generate/warrooms", controller.GenerateDraftWarRooms).Methods("GET")
	// apiRouter.HandleFunc("/admin/generate/draft/grades", controller.GenerateDraftGrades).Methods("GET")
	// apiRouter.HandleFunc("/admin/generate/draft/rounds", controller.GeneratePredictionRound).Methods("GET")
	// apiRouter.HandleFunc("/admin/migrate/data", controller.MigratePlayers).Methods("GET")
	// apiRouter.HandleFunc("/admin/migrate/progress", controller.ProgressPlayers).Methods("GET")
	// apiRouter.HandleFunc("/admin/migrate/draftees", controller.ProgressPlayers).Methods("GET")
	// apiRouter.HandleFunc("/admin/migrate/new/teams", controller.ImportNewTeams).Methods("GET")
	// apiRouter.HandleFunc("/admin/migrate/nba/players", controller.MigrateNBAPlayersToTables).Methods("GET")
	// apiRouter.HandleFunc("/admin/progress/nba/players", controller.ProgressNBAPlayers).Methods("GET")
	// apiRouter.HandleFunc("/admin/clean/nba/tables", controller.CleanNBAPlayerTables).Methods("GET")

	// Bootstrap
	apiRouter.HandleFunc("/bootstrap/teams/", controller.BootstrapTeamData).Methods("GET")
	apiRouter.HandleFunc("/bootstrap/one/{collegeID}/{proID}", controller.BootstrapBasketballData).Methods("GET")
	apiRouter.HandleFunc("/bootstrap/two/{collegeID}/{proID}", controller.SecondBootstrapBasketballData).Methods("GET")
	apiRouter.HandleFunc("/bootstrap/three/{collegeID}/{proID}", controller.ThirdBootstrapBasketballData).Methods("GET")
	apiRouter.HandleFunc("/bootstrap/news/{collegeID}/{proID}", controller.BootstrapNewsData).Methods("GET")

	// Capsheet Controls
	apiRouter.HandleFunc("/nba/capsheet/generate", controller.GenerateCapsheets).Methods("GET")
	apiRouter.HandleFunc("/nba/contracts/get/value", controller.CalculateContracts).Methods("GET")

	// Draft Controls
	// apiRouter.HandleFunc("/nba/draft/conduct/lottery", controller.ConductDraftLottery).Methods("GET")
	apiRouter.HandleFunc("/nba/draft/export/picks", controller.ExportDraftedPicks).Methods("POST")
	// apiRouter.HandleFunc("/nba/draft/run/combine", controller.RunNBACombine).Methods("GET")
	apiRouter.HandleFunc("/nba/draft/page/{teamID}", controller.GetDraftPageData).Methods("GET")
	apiRouter.HandleFunc("/nba/draft/time/change", controller.ToggleDraftTime).Methods("GET")
	apiRouter.HandleFunc("/nba/draft/create/scoutprofile", controller.AddPlayerToScoutBoard).Methods("POST")
	apiRouter.HandleFunc("/nba/draft/reveal/attribute", controller.RevealScoutingAttribute).Methods("POST")
	// apiRouter.HandleFunc("/recruit/revokeScholarshipFromRecruit", controller.RevokeScholarshipFromRecruit).Methods("PUT")
	apiRouter.HandleFunc("/nba/draft/remove/{id}", controller.RemovePlayerFromScoutBoard).Methods("GET")
	apiRouter.HandleFunc("/nba/draft/scout/{id}", controller.GetScoutingDataByDraftee).Methods("GET")
	// apiRouter.HandleFunc("/nba/draft/player/", controller.SaveRecruitingBoard).Methods("POST")

	// Exports
	apiRouter.HandleFunc("/export/cbb/players/all", controller.ExportCollegePlayers).Methods("GET")
	apiRouter.HandleFunc("/export/cbb/preseason", controller.ExportCBBPreseasonRanks).Methods("GET")
	apiRouter.HandleFunc("/export/cbb/team/{teamID}", controller.ExportCBBRosterToCSV).Methods("GET")
	apiRouter.HandleFunc("/export/nba/team/{teamID}", controller.ExportNBARosterToCSV).Methods("GET")

	// Free Agency Controls
	apiRouter.HandleFunc("/nba/freeagency/available/{teamID}", controller.FreeAgencyAvailablePlayers).Methods("GET")
	apiRouter.HandleFunc("/nba/freeagency/create/offer", controller.CreateFreeAgencyOffer).Methods("POST")
	apiRouter.HandleFunc("/nba/freeagency/cancel/offer", controller.CancelFreeAgencyOffer).Methods("POST")
	apiRouter.HandleFunc("/nba/freeagency/create/waiver", controller.CreateWaiverOffer).Methods("POST")
	apiRouter.HandleFunc("/nba/freeagency/cancel/waiver", controller.CancelWaiverOffer).Methods("POST")
	apiRouter.HandleFunc("/nba/extension/create/offer", controller.CreateExtensionOffer).Methods("POST")
	apiRouter.HandleFunc("/nba/extension/cancel/offer", controller.CancelExtensionOffer).Methods("POST")
	apiRouter.HandleFunc("/nba/freeagency/extensions/temp", controller.ExtendPlayers).Methods("GET")
	apiRouter.HandleFunc("/nba/freeagency/sync/round", controller.SyncFreeAgencyOffers).Methods("GET")

	// Gameplan controls
	apiRouter.HandleFunc("/cbb/gameplans/{teamId}", controller.GetGameplansByTeamId).Methods("GET")
	apiRouter.HandleFunc("/nba/gameplans/{teamId}", controller.GetNBAGameplanByTeamId).Methods("GET")
	apiRouter.HandleFunc("/cbb/gameplans/update", controller.UpdateGameplan).Methods("POST")
	apiRouter.HandleFunc("/nba/gameplans/update", controller.UpdateNBAGameplan).Methods("POST")

	// Generation Controls
	// apiRouter.HandleFunc("/admin/generateCoaches", controller.GenerateCoaches).Methods("GET")
	// apiRouter.HandleFunc("/admin/generateTeam", controller.GeneratePlayers).Methods("GET")
	// apiRouter.HandleFunc("/admin/generateTestPlayers", controller.GenerateTestPlayers).Methods("GET")
	// apiRouter.HandleFunc("/admin/generateCroots", controller.GenerateCroots).Methods("GET")
	// apiRouter.HandleFunc("/admin/generate/walkons", controller.GenerateCollegeWalkons).Methods("GET")
	// apiRouter.HandleFunc("/admin/generate/international", controller.GenerateInternationalPlayers).Methods("GET")
	// apiRouter.HandleFunc("/admin/generate/worldcup/players", controller.GenerateInternationalPlayers).Methods("GET")
	// apiRouter.HandleFunc("/admin/migrate/international", controller.MoveISLPlayersToDraft).Methods("GET")
	// apiRouter.HandleFunc("/admin/allocate/international/rosters", controller.GenerateInternationalRoster).Methods("GET")
	// apiRouter.HandleFunc("/admin/fix/nba/records", controller.FixNBASeasonTables).Methods("GET")
	// apiRouter.HandleFunc("/generate/new/attributes", controller.GenerateNewAttributes).Methods("GET")
	// apiRouter.HandleFunc("/fix/nba/matches", controller.SwapNBATeamsTEMP).Methods("GET")

	// Import
	// apiRouter.HandleFunc("/import/nba", controller.ImportNBATeamsAndArenas).Methods("GET")
	// apiRouter.HandleFunc("/import/cbb/games", controller.ImportCBBMatches).Methods("GET")
	// apiRouter.HandleFunc("/import/nba/games", controller.ImportNBAMatches).Methods("GET")
	// apiRouter.HandleFunc("/import/nba/playin/games", controller.ImportNBAMatchesOLD).Methods("GET")
	// apiRouter.HandleFunc("/import/isl/games", controller.ImportISLMatches).Methods("GET")
	// apiRouter.HandleFunc("/import/nba/series", controller.ImportNBASeries).Methods("GET")
	// apiRouter.HandleFunc("/rollback/nba/season", controller.RollbackNBASeason).Methods("GET")
	// apiRouter.HandleFunc("/import/custom/croots", controller.ImportCustomCroots).Methods("GET")
	// apiRouter.HandleFunc("/fix/empty/country", controller.FixEmptyCountryValues).Methods("GET")

	// apiRouter.HandleFunc("/import/archetypes", controller.ImportArchetypes).Methods("GET")
	// apiRouter.HandleFunc("/import/fa/preferences", controller.ImportFAPreferences).Methods("GET")
	// apiRouter.HandleFunc("/import/minutes/expectations", controller.ImportPlaytimeExpectations).Methods("GET")
	// apiRouter.HandleFunc("/import/positions", controller.ImportNewPositions).Methods("GET")
	// apiRouter.HandleFunc("/import/ai/values", controller.MigrateNewAIRecruitingValues).Methods("GET")
	// apiRouter.HandleFunc("/import/new/personalities", controller.ImportPersonalities).Methods("GET")
	// apiRouter.HandleFunc("/import/nba/picks", controller.ImportDraftPicks).Methods("GET")
	// apiRouter.HandleFunc("/migrate/remaining/croots", controller.MigrateRecruits).Methods("GET")

	// International Super League
	// apiRouter.HandleFunc("/import/isl/scoutingdept", controller.ImportISLScouting).Methods("GET")
	// apiRouter.HandleFunc("/isl/identify/players", controller.ISLIdentifyYouthPlayers).Methods("GET")
	// apiRouter.HandleFunc("/isl/scout/players", controller.ISLScoutYouthPlayers).Methods("GET")
	// apiRouter.HandleFunc("/isl/invest/players", controller.ISLInvestYouthPlayers).Methods("GET")
	// apiRouter.HandleFunc("/isl/quick/sync/players", controller.ISLSyncYouthPlayers).Methods("GET")
	// apiRouter.HandleFunc("/isl/quick/draft/players", controller.ISLGenerateNewBatch).Methods("GET")

	// Match Controls
	apiRouter.HandleFunc("/match/{matchId}", controller.GetMatchByMatchId).Methods("GET")
	// apiRouter.HandleFunc("/match/fix/isl", controller.FixISLMatchData).Methods("GET")
	apiRouter.HandleFunc("/match/export/results/{seasonID}/{weekID}/{nbaWeekID}/{matchType}", controller.ExportMatchResults).Methods("GET")
	apiRouter.HandleFunc("/match/result/cbb/{matchId}", controller.GetMatchResultByMatchID).Methods("GET")
	apiRouter.HandleFunc("/match/result/nba/{matchId}", controller.GetNBAMatchResultByMatchID).Methods("GET")
	apiRouter.HandleFunc("/match/team/{teamId}/season/{seasonId}", controller.GetMatchesByTeamIdAndSeasonId).Methods("GET")
	apiRouter.HandleFunc("/match/week/{weekId}", controller.GetMatchesByTeamIdAndSeasonId).Methods("GET")
	apiRouter.HandleFunc("/match/season/{seasonID}", controller.GetMatchesBySeasonID).Methods("GET")
	apiRouter.HandleFunc("/match/team/upcoming/{teamId}/season/{seasonId}", controller.GetUpcomingMatchesByTeamIdAndSeasonId).Methods("GET")
	apiRouter.HandleFunc("/cbb/match/data/{homeTeamAbbr}/{awayTeamAbbr}", controller.GetCBBMatchData).Methods("GET")
	apiRouter.HandleFunc("/nba/match/data/{homeTeamID}/{awayTeamID}", controller.GetNBAMatchData).Methods("GET")
	apiRouter.HandleFunc("/nba/match/team/{teamId}/season/{seasonId}", controller.GetNBAMatchesByTeamIdAndSeasonId).Methods("GET")

	// Migrations
	// apiRouter.HandleFunc("/migrate/faces", controller.MigrateFaceData).Methods("GET")

	// News Controls
	apiRouter.HandleFunc("/cbb/news/all/", controller.GetAllCBBNewsInASeason).Methods("GET")
	apiRouter.HandleFunc("/nba/news/all/", controller.GetAllNBANewsInASeason).Methods("GET")
	apiRouter.HandleFunc("/news/feed/{league}/{teamID}/", controller.GetNewsFeed).Methods("GET")

	// Notification Controls
	apiRouter.HandleFunc("/bba/inbox/get/{cbbID}/{nbaID}/", controller.GetBBAInbox).Methods("GET")
	apiRouter.HandleFunc("/notification/toggle/{notiID}", controller.ToggleNotificationAsRead).Methods("GET")
	apiRouter.HandleFunc("/notification/delete/{notiID}", controller.DeleteNotification).Methods("GET")

	// Player Controls
	apiRouter.HandleFunc("/player/AllPlayers", controller.AllCollegePlayers).Methods("GET")
	// apiRouter.HandleFunc("/player/add/{firstname}/{lastname}", controller.NewPlayer).Methods("POST")
	apiRouter.HandleFunc("/cbb/player/assign/redshirt/", controller.AssignRedshirtForCollegePlayer).Methods("POST")
	apiRouter.HandleFunc("/cbb/players/redshirt/{playerId}/", controller.RedshirtCBBPlayer).Methods("GET")
	apiRouter.HandleFunc("/player/GetPlayer/{playerId}", controller.PlayerById).Methods("GET")
	// apiRouter.HandleFunc("/player/SetRedshirting/{playerId}", controller.SetRedshirtStatusByPlayerId).Methods("PUT")
	apiRouter.HandleFunc("/players", controller.AllPlayers).Methods("GET")
	apiRouter.HandleFunc("/players/{teamId}", controller.AllPlayersByTeamId).Methods("GET")
	apiRouter.HandleFunc("/players/college", controller.AllCollegePlayers).Methods("GET")
	apiRouter.HandleFunc("/players/college/recruits", controller.AllCollegeRecruits).Methods("GET")
	apiRouter.HandleFunc("/collegeplayers/check/declaration", controller.CheckDeclarationStatus).Methods("GET")
	apiRouter.HandleFunc("/players/nba", controller.AllNBAPlayers).Methods("GET")
	apiRouter.HandleFunc("/players/nba/freeAgents", controller.AllNBAFreeAgents).Methods("GET")
	apiRouter.HandleFunc("/nba/players/{teamId}", controller.GetNBARosterByTeamID).Methods("GET")
	apiRouter.HandleFunc("/cbb/players/cut/{playerID}", controller.CutPlayerFromCBBTeam).Methods("GET")
	apiRouter.HandleFunc("/nba/players/cut/{playerID}", controller.CutPlayerFromNBATeam).Methods("GET")
	apiRouter.HandleFunc("/nba/players/activate/option/{contractID}", controller.ActivateOption).Methods("GET")
	apiRouter.HandleFunc("/nba/players/place/gleague/{playerID}", controller.PlaceNBAPlayerInGLeague).Methods("GET")
	apiRouter.HandleFunc("/nba/players/place/twoway/{playerID}", controller.AssignNBAPlayerAsTwoWay).Methods("GET")

	// Poll Controls
	apiRouter.HandleFunc("/college/poll/create/", controller.CreatePollSubmission).Methods("POST")
	apiRouter.HandleFunc("/college/poll/sync", controller.SyncCollegePoll).Methods("GET")
	apiRouter.HandleFunc("/college/poll/official/season/{seasonID}", controller.GetOfficialPollsBySeasonID).Methods("GET")
	apiRouter.HandleFunc("/college/poll/submission/{username}", controller.GetPollSubmission).Methods("GET")

	// Recruit Controls
	apiRouter.HandleFunc("/recruiting/dashboard/{teamID}/", controller.GetRecruitingDataForOverviewPage).Methods("GET")
	apiRouter.HandleFunc("/recruiting/profile/teamboard/{teamID}", controller.GetRecruitingDataForTeamBoardPage).Methods("GET")
	apiRouter.HandleFunc("/recruiting/profile/all/", controller.GetAllRecruitingProfiles).Methods("GET")
	apiRouter.HandleFunc("/recruiting/profile/determine/size/", controller.DetermineRecruitingClassSize).Methods("GET")
	apiRouter.HandleFunc("/recruiting/class/{teamID}/", controller.GetRecruitingClassByTeamID).Methods("GET")
	apiRouter.HandleFunc("/recruiting/add/recruit/", controller.AddRecruitToBoardV2).Methods("POST")
	apiRouter.HandleFunc("/recruit/createRecruitingPointsProfile", controller.AddRecruitToBoard).Methods("POST")
	apiRouter.HandleFunc("/recruit/allocatePoints", controller.AllocateRecruitingPointsForRecruit).Methods("PUT")
	apiRouter.HandleFunc("/recruit/toggleScholarship", controller.SendScholarshipToRecruit).Methods("POST")
	apiRouter.HandleFunc("/recruit/toggle/Scholarship/v2", controller.SendScholarshipToRecruitV2).Methods("POST")
	// apiRouter.HandleFunc("/recruit/revokeScholarshipFromRecruit", controller.RevokeScholarshipFromRecruit).Methods("PUT")
	apiRouter.HandleFunc("/recruit/removeRecruit", controller.RemoveRecruitFromBoard).Methods("POST")
	apiRouter.HandleFunc("/recruit/remove/recruit/v2", controller.RemoveRecruitFromBoardV2).Methods("POST")
	apiRouter.HandleFunc("/recruit/saveRecruitingBoard", controller.SaveRecruitingBoard).Methods("POST")
	apiRouter.HandleFunc("/recruit/save/ai/settings", controller.SaveAIBehavior).Methods("POST")
	apiRouter.HandleFunc("/recruit/save/ai/toggle/{teamID}", controller.SaveAIBehavior).Methods("GET")
	apiRouter.HandleFunc("/croots/export/all", controller.ExportCroots).Methods("GET")

	// Request Controls
	apiRouter.HandleFunc("/requests/", controller.GetTeamRequests).Methods("GET")
	apiRouter.HandleFunc("/requests/createTeamRequest", controller.CreateTeamRequest).Methods("POST")
	apiRouter.HandleFunc("/requests/approveTeamRequest", controller.ApproveTeamRequest).Methods("PUT")
	apiRouter.HandleFunc("/requests/rejectTeamRequest", controller.RejectTeamRequest).Methods("DELETE")
	apiRouter.HandleFunc("/requests/view/cbb/{teamID}/", controller.ViewCBBTeamUponRequest).Methods("GET")
	apiRouter.HandleFunc("/requests/view/nba/{teamID}/", controller.ViewNBATeamUponRequest).Methods("GET")
	apiRouter.HandleFunc("/nba/requests/all/", controller.GetNBATeamRequests).Methods("GET")
	apiRouter.HandleFunc("/nba/requests/create/", controller.CreateNBATeamRequest).Methods("POST")
	apiRouter.HandleFunc("/nba/requests/approve/", controller.ApproveNBATeamRequest).Methods("POST")
	apiRouter.HandleFunc("/nba/requests/reject/", controller.RejectNBATeamRequest).Methods("POST")
	apiRouter.HandleFunc("/nba/requests/revoke/", controller.RemoveNBAUserFromNBATeam).Methods("POST")

	// Run Controls
	// apiRouter.HandleFunc("/run/promises", controller.RunPromises).Methods("GET")

	// Stats Controls
	apiRouter.HandleFunc("/stats/export/{league}/{seasonID}/{weekID}/{matchType}/{viewType}/{playerView}", controller.ExportStats).Methods("GET")
	apiRouter.HandleFunc("/stats/player/{playerId}", controller.GetPlayerStats).Methods("GET")
	apiRouter.HandleFunc("/stats/player/{playerId}/match/{matchId}", controller.GetPlayerStatsByMatch).Methods("GET")
	apiRouter.HandleFunc("/stats/player/{playerId}/season/{seasonId}", controller.GetPlayerStatsBySeason).Methods("GET")
	apiRouter.HandleFunc("/stats/team/{teamId}/season/{seasonId}", controller.GetTeamStatsBySeason).Methods("GET")
	apiRouter.HandleFunc("/stats/team/{teamId}/match/{matchId}", controller.GetCBBTeamStatsByMatch).Methods("GET")
	apiRouter.HandleFunc("/stats/cbb/fix/player/stats", controller.FixPlayerStatsFromLastSeason).Methods("GET")
	apiRouter.HandleFunc("/stats/cbb/{seasonID}/{weekID}/{matchType}/{viewType}", controller.GetCBBStatsPageData).Methods("GET")
	apiRouter.HandleFunc("/stats/nba/{seasonID}/{weekID}/{matchType}/{viewType}", controller.GetNBAStatsPageData).Methods("GET")
	apiRouter.HandleFunc("/stats/nba/team/{teamId}/match/{matchId}", controller.GetNBATeamStatsByMatch).Methods("GET")
	apiRouter.HandleFunc("/stats/nba/match/{matchId}", controller.GetPlayerStatsByMatch).Methods("GET")

	// StandingsControls
	apiRouter.HandleFunc("/standings/college/conf/{conferenceId}/{seasonId}", controller.GetConferenceStandingsByConferenceID).Methods("GET")
	apiRouter.HandleFunc("/standings/nba/conf/{conferenceId}/{seasonId}", controller.GetNBAConferenceStandingsByConferenceID).Methods("GET")
	apiRouter.HandleFunc("/standings/college/season/{seasonId}", controller.GetAllConferenceStandings).Methods("GET")
	apiRouter.HandleFunc("/standings/nba/season/{seasonId}", controller.GetAllNBAConferenceStandings).Methods("GET")
	apiRouter.HandleFunc("/standings/reset/all", controller.ResetSeasonStandings).Methods("GET")

	// Team Controls
	apiRouter.HandleFunc("/team/{teamId}", controller.GetTeamByTeamID).Methods("GET")
	apiRouter.HandleFunc("/team/nba/{teamId}", controller.GetNBATeamByTeamID).Methods("GET")
	apiRouter.HandleFunc("/team/removeUserFromTeam/{teamId}", controller.RemoveUserFromTeam).Methods("PUT")
	apiRouter.HandleFunc("/team/nba/removeUserFromTeam/{teamId}", controller.RemoveNBAUserFromNBATeam).Methods("POST")
	apiRouter.HandleFunc("/teams", controller.AllTeams).Methods("GET")
	apiRouter.HandleFunc("/teams/active", controller.AllActiveTeams).Methods("GET")
	apiRouter.HandleFunc("/teams/active/college", controller.AllActiveCollegeTeams).Methods("GET")
	apiRouter.HandleFunc("/teams/college/available", controller.AllAvailableTeams).Methods("GET")
	apiRouter.HandleFunc("/teams/assign/ratings", controller.SyncTeamRatings).Methods("GET")
	apiRouter.HandleFunc("/teams/nba/assign/ratings", controller.SyncNBATeamRatings).Methods("GET")
	apiRouter.HandleFunc("/teams/coached", controller.AllCoachedTeams).Methods("GET")
	apiRouter.HandleFunc("/teams/college", controller.AllCollegeTeams).Methods("GET")
	apiRouter.HandleFunc("/teams/nba", controller.AllNBATeams).Methods("GET")
	apiRouter.HandleFunc("/teams/isl", controller.AllISLTeams).Methods("GET")
	apiRouter.HandleFunc("/teams/pro", controller.AllProfessionalTeams).Methods("GET")
	apiRouter.HandleFunc("/teams/cbb/dashboard/{teamID}", controller.GetCBBDashboardByTeamID).Methods("GET")
	apiRouter.HandleFunc("/teams/nba/dashboard/{teamID}", controller.GetNBADashboardByTeamID).Methods("GET")

	// Trade Controls
	apiRouter.HandleFunc("/trades/nba/all/accepted", controller.GetAllAcceptedTrades).Methods("GET")
	apiRouter.HandleFunc("/trades/nba/all/rejected", controller.GetAllRejectedTrades).Methods("GET")
	apiRouter.HandleFunc("/trades/nba/block/{teamID}", controller.GetNBATradeBlockDataByTeamID).Methods("GET")
	apiRouter.HandleFunc("/trades/nba/place/block/{playerID}", controller.PlaceNBAPlayerOnTradeBlock).Methods("GET")
	apiRouter.HandleFunc("/trades/nba/preferences/update", controller.UpdateTradePreferences).Methods("POST")
	apiRouter.HandleFunc("/trades/nba/create/proposal", controller.CreateNBATradeProposal).Methods("POST")
	apiRouter.HandleFunc("/trades/nba/proposal/accept/{proposalID}", controller.AcceptTradeOffer).Methods("GET")
	apiRouter.HandleFunc("/trades/nba/proposal/reject/{proposalID}", controller.RejectTradeOffer).Methods("GET")
	apiRouter.HandleFunc("/trades/nba/proposal/cancel/{proposalID}", controller.CancelTradeOffer).Methods("GET")
	apiRouter.HandleFunc("/admin/trades/accept/sync/{proposalID}", controller.SyncAcceptedTrade).Methods("GET")
	apiRouter.HandleFunc("/admin/trades/veto/sync/{proposalID}", controller.VetoAcceptedTrade).Methods("GET")
	apiRouter.HandleFunc("/admin/trades/cleanup", controller.CleanUpRejectedTrades).Methods("GET")

	// Transfer Intentions
	apiRouter.HandleFunc("/portal/sync/promises", controller.SyncPromises).Methods("GET")
	apiRouter.HandleFunc("/portal/transfer/intention", controller.ProcessTransferIntention).Methods("GET")
	apiRouter.HandleFunc("/portal/transfer/enter", controller.EnterTheTransferPortal).Methods("GET")
	apiRouter.HandleFunc("/portal/transfer/sync", controller.SyncTransferPortal).Methods("GET")
	apiRouter.HandleFunc("/portal/ai/generate/profiles", controller.FillUpTransferBoardsAI).Methods("GET")
	apiRouter.HandleFunc("/portal/ai/allocate/profiles", controller.AllocateAndPromisePlayersAI).Methods("GET")
	apiRouter.HandleFunc("/portal/page/data/{teamID}", controller.GetTransferPortalPageData).Methods("GET")
	apiRouter.HandleFunc("/portal/profile/create", controller.AddTransferPlayerToBoard).Methods("POST")
	apiRouter.HandleFunc("/portal/profile/remove/{profileID}", controller.RemovePlayerFromTransferPortalBoard).Methods("GET")
	apiRouter.HandleFunc("/portal/saveboard", controller.SaveTransferBoard).Methods("POST")
	apiRouter.HandleFunc("/portal/promise/create", controller.CreatePromise).Methods("POST")
	apiRouter.HandleFunc("/portal/promise/cancel/{promiseID}", controller.CancelPromise).Methods("GET")
	apiRouter.HandleFunc("/portal/promise/player/{playerID}/{teamID}", controller.GetPromiseByPlayerID).Methods("GET")
	apiRouter.HandleFunc("/portal/player/scout/{id}", controller.GetScoutingDataByTransfer).Methods("GET")
	apiRouter.HandleFunc("/portal/export/players/", controller.ExportPortalPlayersToCSV).Methods("GET")

	// Timestamp Controls
	apiRouter.HandleFunc("/simbba/get/timestamp", controller.GetCurrentTimestamp).Methods("GET")
	apiRouter.HandleFunc("/cbb/easter/egg/collude/", controller.CollusionButton).Methods("POST")

	// Discord Controls
	apiRouter.HandleFunc("/ds/cbb/player/id/{id}", controller.CBBPlayerByID).Methods("GET")
	apiRouter.HandleFunc("/ds/nba/player/id/{id}", controller.NBAPlayerByID).Methods("GET")
	apiRouter.HandleFunc("/ds/cbb/player/name/{firstName}/{lastName}/{abbr}", controller.CBBPlayerByNameAndAbbr).Methods("GET")
	apiRouter.HandleFunc("/ds/nba/player/name/{firstName}/{lastName}/{abbr}", controller.NBAPlayerByNameAndAbbr).Methods("GET")
	apiRouter.HandleFunc("/ds/cbb/croots/class/{teamID}/", controller.GetRecruitingClassByTeamID).Methods("GET")
	apiRouter.HandleFunc("/ds/cbb/croot/name/{firstName}/{lastName}", controller.GetCrootsByName).Methods("GET")
	apiRouter.HandleFunc("/ds/cbb/croot/{id}", controller.GetRecruitViaDiscord).Methods("GET")
	apiRouter.HandleFunc("/ds/cbb/team/{teamId}", controller.GetCollegeTeamData).Methods("GET")
	apiRouter.HandleFunc("/ds/nba/team/{teamId}", controller.GetNBATeamDataByID).Methods("GET")
	apiRouter.HandleFunc("/ds/cbb/conf/standings/{conferenceID}", controller.CollegeConferenceStandings).Methods("GET")
	apiRouter.HandleFunc("/ds/nba/conf/standings/{conferenceID}", controller.NBAConferenceStandings).Methods("GET")
	apiRouter.HandleFunc("/ds/cbb/conf/matches/{conferenceID}/{day}", controller.CollegeMatchesByConference).Methods("GET")
	apiRouter.HandleFunc("/ds/cbb/flex/{teamOneID}/{teamTwoID}/", controller.CompareCFBTeams).Methods("GET")
	apiRouter.HandleFunc("/ds/nba/flex/{teamOneID}/{teamTwoID}/", controller.CompareNFLTeams).Methods("GET")
	apiRouter.HandleFunc("/ds/cbb/assign/discord/{teamID}/{discordID}", controller.AssignDiscordIDtoCollegeTeam).Methods("GET")
	apiRouter.HandleFunc("/ds/nba/assign/discord/{teamID}/{discordID}/{username}", controller.AssignDiscordIDtoNBATeam).Methods("GET")

	// Websocket
	myRouter.HandleFunc("/ws", ws.WebSocketHandler)

	// Handler
	handler := cors.AllowAll().Handler(myRouter)

	return handler
}

func loadEnvs() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("CANNOT LOAD ENV VARIABLES")
	}
}
func handleCron() *cron.Cron {

	c := cron.New()
	runJobs := os.Getenv("RUN_JOBS")
	if runJobs != "FALSE" {
		c.AddFunc("0 16 * * 3", controller.SyncRecruitingViaCron)
		c.AddFunc("0 6 * * 4,6", controller.SyncAIBoardsViaCron)
		c.AddFunc("0 10 * * 1,3,5,6", controller.CheckUserGameplansViaCron)
		c.AddFunc("0 20 * * 1,3,5,6", controller.ShowGamesViaCron)
		c.AddFunc("0 22 * * 1,3,5,6", controller.RunAIGameplansViaCron)
		c.AddFunc("0 10 * * 4", controller.FillAIBoardsViaCron)
		c.AddFunc("0 12 * * 0", controller.SyncToNextWeekViaCron)
		c.AddFunc("0 16 * * *", controller.SyncFreeAgencyOffersViaCron)
	}

	c.Start()

	return c
}

func main() {
	loadEnvs()
	InitialMigration()
	fmt.Println("Database initialized.")

	fmt.Println("Setting up polling...")
	go monitorDBForUpdates()

	fmt.Println("Loading cron...")
	cronJobs := handleCron()
	fmt.Println("Loading Handler Requests.")
	fmt.Println("Basketball Server Initialized.")
	srv := &http.Server{
		Addr:    ":8081",
		Handler: handleRequests(),
	}

	go func() {
		fmt.Println("Application Running on port 8081")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", srv.Addr, err)
		}
	}()

	// Create a channel to listen for system interrupts (Ctrl+C, etc.)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Block until a signal is received
	<-quit
	fmt.Println("Shutting down server...")

	// Gracefully shutdown the server with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Stop cron jobs
	cronJobs.Stop()

	// Shutdown the server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server exiting")
}
