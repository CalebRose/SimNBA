package managers

import (
	"fmt"
	"sort"
	"strconv"
	"sync"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
)

func ProcessTransferIntention() {
	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.SeasonID))
	allCollegePlayers := GetAllCollegePlayers()
	seasonStatMap := GetCollegePlayerSeasonStatMap(seasonID)
	fullRosterMap := GetFullTeamRosterWithCrootsMap()
	standingsMap := GetCollegeStandingsMap(seasonID)
	collegeTeamMap := GetCollegeTeamMap()
	// teamProfileMap := GetTeamProfileMap()
	transferCount := 0
	freshmanCount := 0
	redshirtFreshmanCount := 0
	sophomoreCount := 0
	redshirtSophomoreCount := 0
	juniorCount := 0
	redshirtJuniorCount := 0
	seniorCount := 0
	redshirtSeniorCount := 0
	lowCount := 0
	mediumCount := 0
	highCount := 0

	upcomingTeam := "Prefers to play for an up-and-coming team"
	differentState := "Prefers to play in a different state"
	immediateStart := "Prefers to play for a team where he can start immediately"
	closeToHome := "Prefers to be close to home"
	nationalChampionshipContender := "Prefers to play for a national championship contender"
	specificCoach := "Prefers to play for a specific coach"
	legacy := "Legacy"
	richHistory := "Prefers to play for a team with a rich history"
	bigDrop := -25.0
	mediumDrop := -15.0
	smallDrop := -10.0
	giantDrop := -33.0
	// tinyDrop := -5.0
	// tinyGain := 5.0
	smallGain := 10.0
	mediumGain := 15.0
	bigGain := 25.0
	giantgain := 33.0
	for _, p := range allCollegePlayers {
		// Do not include redshirts and all graduating players
		if p.IsRedshirting || p.WillDeclare {
			continue
		}

		transferWeight := 0.0

		// Modifiers on reasons why they would transfer
		minutesMod := 0.0
		ageMod := 0.0
		starMod := 0.0
		depthChartCompetitionMod := 0.0
		// schemeMod := 0.0
		biasMod := 0.0

		// Check Minutes
		seasonStats := seasonStatMap[p.ID]
		minutesPerGame := seasonStats.MinutesPerGame

		if minutesPerGame < float64(p.PlaytimeExpectations) {
			minutesMod = giantgain
		} else {
			minutesMod = giantDrop
		}

		// Check Age
		// The more experienced the player is in the league,
		// the more likely they will transfer.
		/// Have this be a multiplicative factor to odds
		if p.Year == 1 {
			ageMod = .33
		} else if p.Year == 2 {
			ageMod = .5
		} else if p.Year == 3 {
			ageMod = .66
		} else if p.Year == 4 {
			ageMod = 1
		}

		/// Higher star players are more likely to transfer
		if p.Stars == 0 {
			starMod = 1
		} else if p.Stars == 1 {
			starMod = .66
		} else if p.Stars == 2 {
			starMod = .75
		} else if p.Stars == 3 {
			starMod = 1
		} else if p.Stars == 4 {
			starMod = 1.2
		} else if p.Stars == 5 {
			starMod = 1.5
		}

		teamRoster := fullRosterMap[uint(p.TeamID)]
		filteredRosterByPosition := filterRosterByPosition(teamRoster, p.Position)
		youngerPlayerAhead := false
		idFound := false
		for idx, pl := range filteredRosterByPosition {
			if pl.Age < p.Age && !idFound {
				youngerPlayerAhead = true
			}
			if pl.ID == p.ID && idx > 2 {
				idFound = true
				// Check the index of the player.
				// If they're at the top of the list, they're considered to be starting caliber.
				depthChartCompetitionMod += bigGain
			}
		}

		// If there's a modifier applied and there's a younger player ahead on the roster, double the amount on the modifier
		if depthChartCompetitionMod > 0 {
			if youngerPlayerAhead {
				depthChartCompetitionMod += bigGain
			} else {
				depthChartCompetitionMod = .5 * depthChartCompetitionMod
			}
		}

		// Bias Mod
		if p.RecruitingBias == upcomingTeam {
			standings := standingsMap[p.TeamID]
			if standings.PostSeasonStatus == "National Champion" || standings.PostSeasonStatus == "Conference Champion" ||
				standings.PostSeasonStatus == "Round of 32" || standings.PostSeasonStatus == "Sweet 16" || standings.PostSeasonStatus == "Elite 8" || standings.PostSeasonStatus == "Final Four" {
				biasMod = mediumDrop
			} else {
				biasMod = mediumGain
			}
		} else if p.RecruitingBias == nationalChampionshipContender {
			standings := standingsMap[p.TeamID]
			if standings.PostSeasonStatus == "National Champion" ||
				standings.PostSeasonStatus == "Final Four" {
				biasMod = bigDrop
			} else {
				biasMod = bigGain
			}
		} else if p.RecruitingBias == immediateStart && minutesMod > 0 {
			biasMod = mediumDrop
		} else if p.RecruitingBias == immediateStart && minutesMod <= 0 {
			biasMod = mediumGain
		} else if p.RecruitingBias == closeToHome && p.Country == "USA" {
			team := collegeTeamMap[p.TeamID]
			if team.State != p.State {
				biasMod = mediumGain
			} else {
				biasMod = mediumDrop
			}
		} else if p.RecruitingBias == differentState && p.Country == "USA" {
			team := collegeTeamMap[p.TeamID]
			if team.State != p.State {
				biasMod = mediumDrop
			} else {
				biasMod = mediumGain
			}
		} else if p.RecruitingBias == specificCoach {
			team := collegeTeamMap[p.TeamID]
			if team.Coach == p.RecruitingBiasValue {
				biasMod = mediumGain
			} else {
				biasMod = mediumDrop
			}
		} else if p.RecruitingBias == legacy {
			team := collegeTeamMap[p.TeamID]
			legacyID := util.ConvertStringToInt(p.RecruitingBiasValue)
			if uint(legacyID) > 0 && team.ID == uint(legacyID) {
				biasMod = smallGain
			} else {
				biasMod = smallDrop
			}
		} else if p.RecruitingBias == richHistory {
			biasMod = 0.0
		}

		/// Not playing = 25, low depth chart = 16 or 33, scheme = 10, if you're all 3, that's a ~60% chance of transferring pre- modifiers
		transferWeight = starMod * ageMod * (minutesMod + depthChartCompetitionMod + biasMod)
		diceRoll := util.GenerateIntFromRange(1, 100)
		// NOT INTENDING TO TRANSFER
		transferInt := int(transferWeight)
		if diceRoll > transferInt {
			continue
		}

		if p.Year == 1 {
			fmt.Println("Dice Roll: ", diceRoll)
		}

		// Is Intending to transfer
		p.DeclareTransferIntention(strconv.Itoa(int(transferWeight)))
		transferCount++
		if p.Year == 1 && !p.IsRedshirt {
			freshmanCount++
		} else if p.Year == 2 && p.IsRedshirt {
			redshirtFreshmanCount++
		} else if p.Year == 2 && !p.IsRedshirt {
			sophomoreCount++
		} else if p.Year == 3 && p.IsRedshirt {
			redshirtSophomoreCount++
		} else if p.Year == 3 && !p.IsRedshirt {
			juniorCount++
		} else if p.Year == 4 && p.IsRedshirt {
			redshirtJuniorCount++
		} else if p.Year == 4 && !p.IsRedshirt {
			seniorCount++
		} else if p.Year == 5 && p.IsRedshirt {
			redshirtSeniorCount++
		}

		if transferWeight < 25 {
			lowCount++
		} else if transferWeight < 45 {
			mediumCount++
		} else {
			highCount++
		}
		fmt.Println(strconv.Itoa(p.Year)+" YEAR "+p.TeamAbbr+" "+p.Position+" "+p.FirstName+" "+p.LastName+" HAS ANNOUNCED THEIR INTENTION TO TRANSFER | Weight: ", int(transferWeight))
	}
	fmt.Println("Total number of players entering the transfer portal: ", transferCount)
	fmt.Println("Total number of freshmen entering the transfer portal: ", freshmanCount)
	fmt.Println("Total number of redshirt freshmen entering the transfer portal: ", redshirtFreshmanCount)
	fmt.Println("Total number of sophomores entering the transfer portal: ", sophomoreCount)
	fmt.Println("Total number of redshirt sophomores entering the transfer portal: ", redshirtSophomoreCount)
	fmt.Println("Total number of juniors entering the transfer portal: ", juniorCount)
	fmt.Println("Total number of redshirt juniors entering the transfer portal: ", redshirtJuniorCount)
	fmt.Println("Total number of seniors entering the transfer portal: ", seniorCount)
	fmt.Println("Total number of redshirt seniors entering the transfer portal: ", redshirtSeniorCount)
	fmt.Println("Total number of players with low likeliness to enter transfer portal: ", lowCount)
	fmt.Println("Total number of players with medium likeliness to enter transfer portal: ", mediumCount)
	fmt.Println("Total number of players with high likeliness to enter transfer portal: ", highCount)
}

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

func filterRosterByPosition(roster []structs.CollegePlayer, pos string) []structs.CollegePlayer {
	estimatedSize := len(roster) / 5 // Adjust this based on your data
	filteredList := make([]structs.CollegePlayer, 0, estimatedSize)
	for _, p := range roster {
		if p.Position != pos || p.WillDeclare {
			continue
		}
		filteredList = append(filteredList, p)
	}
	sort.Slice(filteredList, func(i, j int) bool {
		return filteredList[i].Overall > filteredList[j].Overall
	})

	return filteredList
}
