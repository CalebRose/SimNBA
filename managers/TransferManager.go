package managers

import (
	"sync"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
)

func AICoachPromisePhase() {

}

func CreatePromise() {

}

func UpdatePromise() {

}

func DeletePromise() {

}

func EnterTheTransferPortal() {

}

func AddTransferPlayerToBoard() {

}

func RemovePlayerFromTransferPortal() {

}

func AllocatePointsToTransferPlayer() {

}

func AICoachFillBoardsPhase() {

}

func AICoachAllocateAndPromisePhase() {

}

func SyncTransferPortal() {

}

// At end of season, sync through promises to confirm if promises were made
func SyncPromises() {

}

func GetPromisesByTeamID(teamID string) []structs.CollegePromise {
	db := dbprovider.GetInstance().GetDB()

	var promises []structs.CollegePromise

	db.Where("team_id = ?", teamID).Find(&promises)

	return promises
}

func GetPortalProfilesByTeamID(teamID string) []structs.TransferPortalProfile {
	db := dbprovider.GetInstance().GetDB()

	var profiles []structs.TransferPortalProfile

	db.Preload("CollegePlayer").Where("profile_id = ?", teamID).Find(&profiles)

	return profiles
}

func GetTransferPortalData(teamID string) structs.TransferPortalResponse {
	var waitgroup sync.WaitGroup
	waitgroup.Add(4)
	profileChan := make(chan structs.TeamRecruitingProfile)
	playersChan := make(chan []structs.CollegePlayer)
	boardChan := make(chan []structs.TransferPortalProfile)
	promiseChan := make(chan []structs.CollegePromise)

	go func() {
		waitgroup.Wait()
		close(profileChan)
		close(playersChan)
		close(boardChan)
		close(promiseChan)
	}()

	go func() {
		defer waitgroup.Done()
		profile := GetOnlyTeamRecruitingProfileByTeamID(teamID)
		profileChan <- profile
	}()
	go func() {
		defer waitgroup.Done()
		tpPlayers := GetTransferPortalPlayers()
		playersChan <- tpPlayers
	}()
	go func() {
		defer waitgroup.Done()
		tpProfiles := GetPortalProfilesByTeamID(teamID)
		boardChan <- tpProfiles
	}()
	go func() {
		defer waitgroup.Done()
		cPromises := GetPromisesByTeamID(teamID)
		promiseChan <- cPromises
	}()

	teamProfile := <-profileChan
	players := <-playersChan
	board := <-boardChan
	promises := <-promiseChan

	return structs.TransferPortalResponse{
		Team:         teamProfile,
		Players:      players,
		TeamBoard:    board,
		TeamPromises: promises,
	}
}
