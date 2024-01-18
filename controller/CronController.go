package controller

import (
	"github.com/CalebRose/SimNBA/managers"
)

func FillAIBoardsViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && ts.CollegeWeek < 15 && ts.CollegeWeek > 0 {
		managers.FillAIRecruitingBoards()
	}
}

func SyncAIBoardsViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && ts.CollegeWeek < 15 && ts.CollegeWeek > 0 {
		managers.ResetAIBoardsForCompletedTeams()
		managers.AllocatePointsToAIBoards()
	}
}

func SyncRecruitingViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && ts.CollegeWeek < 15 && ts.CollegeWeek > 0 {
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
