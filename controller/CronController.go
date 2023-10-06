package controller

import (
	"github.com/CalebRose/SimNBA/managers"
)

func FillAIBoardsViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron {
		managers.FillAIRecruitingBoards()
	}
}

func SyncAIBoardsViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron {
		managers.ResetAIBoardsForCompletedTeams()
		managers.AllocatePointsToAIBoards()
	}
}

func SyncRecruitingViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron {
		managers.SyncRecruiting()
	}
}

func SyncToNextWeekViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron {
		managers.ResetCollegeStandingsRanks()
		managers.SyncToNextWeek()
		managers.SyncCollegePollSubmissionForCurrentWeek()
	}
}

func ShowGamesViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron {
		managers.ShowGames()
	}
}

func SyncFreeAgencyOffersViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron {
		managers.SyncFreeAgencyOffers()
		managers.MoveUpInOffseasonFreeAgency()
	}
}
