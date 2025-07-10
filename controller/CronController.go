package controller

import (
	"fmt"

	"github.com/CalebRose/SimNBA/managers"
)

func FillAIBoardsViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && ts.CollegeWeek < 15 && ts.CollegeWeek > 0 {
		managers.FillAIRecruitingBoards()
	}

	if ts.RunCron && ts.IsOffSeason && ts.CollegeSeasonOver {
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

	if ts.RunCron && ts.IsOffSeason && ts.CollegeSeasonOver {
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
	if ts.RunCron && ts.IsOffSeason && ts.CollegeSeasonOver {
		// Run First Phase of Transfer Portal
		if ts.TransferPortalPhase == 1 {
			managers.ProcessEarlyDeclareeAnnouncements()
			managers.ProgressionMain()
			managers.ProcessTransferIntention()
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

	if ts.RunCron && ts.IsNBAOffseason {
		// managers.SyncISLYouthDevelopment()
	}

	if ts.RunCron && ((!ts.IsOffSeason || !ts.IsNBAOffseason) || (ts.CollegeSeasonOver && ts.NBASeasonOver && ts.FreeAgencyRound > 2)) {
		managers.ProcessRecovery()
		managers.SyncToNextWeek()
	}

	if ts.RunCron && ts.IsNBAOffseason && ts.IsFreeAgencyLocked {
		// If NBA Progression Wasn't Ran, Run Progression
		if !ts.ProgressedProfessionalPlayers {
			// managers.ProgressNBAPlayers()
		}
	}
}

func ShowGamesViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && ts.RunGames && (!ts.IsOffSeason || !ts.IsNBAOffseason) {
		managers.ShowGames()
	}
}

func RunAIGameplansViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && (!ts.IsOffSeason || !ts.IsNBAOffseason) {
		managers.ProcessRecovery()
		val := managers.SetAIGameplans()
		if val {
			fmt.Println("AI Gameplans SET!")
		}
	}
}

func SyncFreeAgencyOffersViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && !ts.IsFreeAgencyLocked && !ts.IsDraftTime {
		managers.SyncAIOffers()
		managers.SyncFreeAgencyOffers()
		managers.AllocateCapsheets()
	}
	if ts.RunCron && ts.NBASeasonOver {
		managers.RunExtensionsAlgorithm()
	}
}
