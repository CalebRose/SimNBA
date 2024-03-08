package controller

import (
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
			managers.ProcessTransferIntention()
			managers.ProcessEarlyDeclareeAnnouncements()
		} else if ts.TransferPortalPhase == 2 && !ts.ProgressedCollegePlayers {
			// Run Second Phase of Transfer Portal (Progressions & Move Players Over)
			// If CBB Progression wasn't ran
			managers.ProgressionMain()
			managers.SyncPromises()
			managers.EnterTheTransferPortal()
		} else if ts.TransferPortalPhase == 3 {
			// Run Transfer Portal (Rounds 1-10)
			managers.SyncTransferPortal()
		}
	}

	if ts.RunCron && ts.IsOffSeason && !ts.CrootsGenerated {
		managers.GenerateCroots()
		// Reset Team Profiles, allocate bonus points where necessary
		// Generate Standings Records too
		managers.GenerateCollegeStandings()
		managers.GenerateNBAStandings()
	}
}

func SyncToNextWeekViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && ((!ts.IsOffSeason || !ts.IsNBAOffseason) || (ts.CollegeSeasonOver && ts.NBASeasonOver && ts.FreeAgencyRound > 2)) {
		managers.SyncToNextWeek()
	}

	if ts.RunCron && ts.IsNBAOffseason && ts.IsFreeAgencyLocked {
		// If NBA Progression Wasn't Ran, Run Progression
		if !ts.ProgressedProfessionalPlayers && ts.NBASeasonOver {
			managers.ProgressNBAPlayers()
		}
	}

}

func ShowGamesViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && (!ts.IsOffSeason || !ts.IsNBAOffseason) {
		managers.ShowGames()
	}
}

func SyncFreeAgencyOffersViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && !ts.IsFreeAgencyLocked {
		managers.SyncFreeAgencyOffers()
		if ts.NBASeasonOver {
			managers.RunExtensionsAlgorithm()
		}
		managers.MoveUpInOffseasonFreeAgency()
	}
}
