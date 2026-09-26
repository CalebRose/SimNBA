package controller

import (
	"fmt"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/managers"
	"github.com/CalebRose/SimNBA/repository"
	"github.com/CalebRose/SimNBA/structs"
)

func FillAIBoardsViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && ts.CollegeWeek < 15 && ts.CollegeWeek > 0 {
		managers.FillAIRecruitingBoards()
	}

	if ts.RunCron && ts.CollegeSeasonOver {
		if ts.TransferPortalPhase == 2 {
			managers.AICoachPromisePhase()
		}
		if ts.TransferPortalPhase == 3 {
			managers.AICoachFillBoardsPhase()
		}
	}
}

func SyncAIBoardsViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && ts.CollegeWeek < 15 && ts.CollegeWeek > 0 {
		managers.ResetAIBoardsForCompletedTeams()
		managers.AllocatePointsToAIBoards()
	}

	if ts.RunCron && ts.CollegeSeasonOver {
		if ts.TransferPortalPhase == 3 {
			// Sync points and promise in the transfer portal
			managers.AICoachAllocateAndPromisePhase()
		}
	}
}

func SyncRecruitingViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && ts.CollegeWeek < 15 && ts.CollegeWeek > 0 {
		managers.SyncRecruiting()
	}
	if ts.RunCron && ts.CollegeSeasonOver {
		// Run First Phase of Transfer Portal
		if ts.TransferPortalPhase == 1 {
			managers.ProcessEarlyDeclareeAnnouncements()
			// managers.ProgressionMain()
			// managers.ProcessTransferIntention()
		} else if ts.TransferPortalPhase == 2 && !ts.ProgressedCollegePlayers {
			// Run Second Phase of Transfer Portal (Progressions & Move Players Over)
			// If CBB Progression wasn't ran
			managers.SyncPromises()
			managers.EnterTheTransferPortal()
		} else if ts.TransferPortalPhase == 3 {
			// Run Transfer Portal (Rounds 1-10)
			managers.SyncTransferPortal()
		}
	}

	if ts.RunCron && ts.IsOffSeason && !ts.CollegeSeasonOver && !ts.CrootsGenerated {
		// Reset Team Profiles, allocate bonus points where necessary
		managers.ProgressStandings()
		managers.RunDeclarationsAlgorithm()
		managers.DetermineRecruitingClassSize()
		managers.GenerateCollegeStandings()
		managers.GenerateNBAStandings()
		managers.GenerateCroots()
	}
}

func SyncToNextWeekViaCron() {
	ts := managers.GetTimestamp()
	if !ts.RunCron {
		return
	}

	if ts.IsNBAOffseason {
		// managers.SyncISLYouthDevelopment()
	}

	if (!ts.IsOffSeason || !ts.IsNBAOffseason) || (ts.CollegeSeasonOver && ts.NBASeasonOver && ts.FreeAgencyRound > 2) {
		managers.ProcessRecovery()
		managers.SyncToNextWeek()
	}

	if ts.IsNBAOffseason && ts.IsFreeAgencyLocked {
		// If NBA Progression Wasn't Ran, Run Progression
		if !ts.ProgressedProfessionalPlayers {
			// managers.ProgressNBAPlayers()
		}
	}
}

func CheckUserGameplansViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && ts.RunGames && (!ts.IsOffSeason || !ts.IsNBAOffseason) {
		managers.CheckAllUserGameplans()
	}
}

func ShowGamesViaCron() {
	ts := managers.GetTimestamp()
	if !ts.RunCron {
		return
	}
	if ts.Phase > 9 && ts.RunGames && (!ts.IsOffSeason || !ts.IsNBAOffseason) {
		managers.ShowGames()
	}
}

func RunAIGameplansViaCron() {
	ts := managers.GetTimestamp()
	if !ts.RunCron {
		return
	}
	if ts.Phase > 9 && (!ts.IsOffSeason || !ts.IsNBAOffseason) {
		managers.ProcessRecovery()
		val := managers.SetAIGameplans()
		if val {
			fmt.Println("AI Gameplans SET!")
		}
	}
}

func SyncFreeAgencyOffersViaCron() {
	ts := managers.GetTimestamp()
	if !ts.RunCron {
		return
	}

	if !ts.IsFreeAgencyLocked && !ts.IsDraftTime {
		managers.SyncAIOffers()
		managers.SyncFreeAgencyOffers()
	}
	if ts.NBASeasonOver {
		managers.RunExtensionsAlgorithm()
	}
	managers.AllocateCapsheets()
}

// Sync Phase Tuesday via Cron - Handling separate actions based on the current phase of the season.
func SyncPhaseTuesdayViaCron() {
	ts := managers.GetTimestamp()
	if !ts.RunCron {
		return
	}
	db := dbprovider.GetInstance().GetDB()
	if ts.Phase == 4 {
		managers.ImportNBAGames()
		managers.ImportISLGames()
	}
	if ts.Phase == 6 {
		managers.GenerateOOCSchedule()
	}

	// Calculate Player Minimum Values and Average Annual Value (AAV) for players
	if ts.Phase == 19 || ts.NBAWeek == 9 {
		// w := http.ResponseWriter(nil)
		// managers.CalculatePlayerMinimumAndAAVValues(w)
	}

	if (ts.CollegeSeasonOver && !ts.ProgressedCollegePlayers) || (ts.Phase == 31 && !ts.ProgressedCollegePlayers) {
		// Reset progression flags on all teams and players before running
		// db.Model(&structs.Team{}).Where("id > ?", 0).Updates(map[string]interface{}{"players_progressed": false, "recruits_added": false})
		db.Model(&structs.CollegePlayer{}).Where("id > ?", 0).Update("has_progressed", false)

		managers.ProgressionMain()
		ts.ToggleCollegeProgression()
		// managers.RecruitingAndTransferPortalCleanUp()
		repository.SaveTimeStamp(ts, db)
	}

	if (ts.NBAWeek == 28 || ts.Phase == 38) && ts.NBASeasonOver && !ts.ProgressedProfessionalPlayers {
		// Reset progression flags on all teams and players before running
		db.Model(&structs.NBAPlayer{}).Where("id > ?", 0).Update("has_progressed", false)
		managers.ProgressNBAPlayers()
		ts.ToggleProfessionalProgression()
		// managers.FreeAgencyCleanUp()
		repository.SaveTimeStamp(ts, db)
	}
}

// Sync Phase Wednesday via Cron - Handling separate actions based on the current phase of the season.
func SyncPhaseWednesdayViaCron() {
	ts := managers.GetTimestamp()
	if !ts.RunCron {
		return
	}

	// Generate Walkons in the middle of the night.
	if ts.CollegeWeek == 20 || ts.Phase == 30 {
		managers.GenerateCollegeWalkons()
		managers.AssignAllRecruitRanks()
	}
}

// Sync Phase Thursday via Cron - Handling separate actions based on the current phase of the season.
func SyncPhaseThursdayViaCron() {
	ts := managers.GetTimestamp()
	if !ts.RunCron {
		return
	}
}

// Sync Phase Friday via Cron - Handling separate actions based on the current phase of the season.
func SyncPhaseFridayViaCron() {
	ts := managers.GetTimestamp()
	if !ts.RunCron {
		return
	}
}
