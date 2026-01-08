package managers

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/repository"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
	"gorm.io/gorm"
)

var upcomingTeam = "Prefers to play for an up-and-coming team"
var differentState = "Prefers to play in a different state"
var immediateStart = "Prefers to play for a team where he can start immediately"
var closeToHome = "Prefers to be close to home"
var nationalChampionshipContender = "Prefers to play for a national championship contender"
var specificCoach = "Prefers to play for a specific coach"
var legacy = "Legacy"
var richHistory = "Prefers to play for a team with a rich history"

func ProcessTransferIntention() {
	ts := GetTimestamp()
	db := dbprovider.GetInstance().GetDB()
	seasonID := strconv.Itoa(int(ts.SeasonID))
	allCollegePlayers := GetAllCollegePlayers()
	seasonStatMap := GetCollegePlayerSeasonStatMap(seasonID)
	fullRosterMap := GetFullTeamRosterWithCrootsMap()
	standingsMap := GetCollegeStandingsMap(seasonID)
	collegeTeamMap := GetCollegeTeamMap()
	// teamProfileMap := GetTeamProfileMap()
	transferCount := 0
	bigDrop := -25.0
	mediumDrop := -15.0
	smallDrop := -10.0
	giantDrop := -33.0
	// tinyDrop := -5.0
	// tinyGain := 5.0
	smallGain := 10.0
	mediumGain := 15.0
	bigGain := 25.0
	giantgain := 40.0
	for _, p := range allCollegePlayers {
		// Do not include redshirts and all graduating players
		if p.IsRedshirting || p.WillDeclare || p.TeamID == 0 || p.TransferStatus == 1 {
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
			ageMod = .001
		} else if p.Year == 2 && p.IsRedshirt {
			ageMod = .2
		} else if p.Year == 2 && !p.IsRedshirt {
			ageMod = .5
		} else if p.Year == 3 && p.IsRedshirt {
			ageMod = .65
		} else if p.Year == 3 && !p.IsRedshirt {
			ageMod = util.GenerateFloatFromRange(0.7, 1)
		} else if p.Year == 4 {
			ageMod = util.GenerateFloatFromRange(1, 1.25)
		} else if p.Year == 5 {
			ageMod = util.GenerateFloatFromRange(1.26, 1.45)
		}

		/// Higher star players are more likely to transfer
		if p.Stars == 0 {
			starMod = 1
		} else if p.Stars == 1 {
			starMod = .66
		} else if p.Stars == 2 {
			starMod = .75
		} else if p.Stars == 3 {
			starMod = util.GenerateFloatFromRange(0.9, 1.1)
		} else if p.Stars == 4 {
			starMod = util.GenerateFloatFromRange(1.11, 1.3)
		} else if p.Stars == 5 {
			starMod = util.GenerateFloatFromRange(1.31, 1.75)
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
		team := collegeTeamMap[p.TeamID]
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
			if team.State != p.State {
				biasMod = mediumGain
			} else {
				biasMod = mediumDrop
			}
		} else if p.RecruitingBias == differentState && p.Country == "USA" {
			if team.State != p.State {
				biasMod = mediumDrop
			} else {
				biasMod = mediumGain
			}
		} else if p.RecruitingBias == specificCoach {
			if team.Coach == p.RecruitingBiasValue {
				biasMod = mediumGain
			} else {
				biasMod = mediumDrop
			}
		} else if p.RecruitingBias == legacy {
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
		diceRoll := util.GenerateIntFromRange(1, 90)
		// NOT INTENDING TO TRANSFER
		transferInt := int(transferWeight)
		// Make it more likely for players on AI teams to transfer
		if !team.IsUserCoached {
			transferInt += util.GenerateIntFromRange(1, 20)
		}
		if diceRoll > transferInt {
			continue
		}

		status := getTransferStatus(int(transferWeight))

		// Is Intending to transfer
		p.DeclareTransferIntention(status)
		transferCount++
		if p.Stars > 2 {
			message := "Breaking News! " + strconv.Itoa(p.Stars) + " star " + p.Position + " " + p.FirstName + " " + p.LastName + " has announced their intention to transfer from " + p.TeamAbbr + "!"
			CreateNewsLog("CBB", message, "Transfer Portal", int(p.TeamID), ts)
		}
		repository.SaveCollegePlayerRecord(p, db)
		fmt.Println(strconv.Itoa(p.Year)+" YEAR "+p.TeamAbbr+" "+p.Position+" "+p.FirstName+" "+p.LastName+" HAS ANNOUNCED THEIR INTENTION TO TRANSFER | Weight: ", int(transferWeight))
	}
	transferPortalMessage := "Breaking News! About " + strconv.Itoa(transferCount) + " players intend to transfer from their current schools. Teams have one week to commit promises to retain players."
	CreateNewsLog("CBB", transferPortalMessage, "Transfer Portal", 0, ts)
	ts.EnactPromisePhase()
	repository.SaveTimeStamp(ts, db)
}

func AICoachPromisePhase() {
	db := dbprovider.GetInstance().GetDB()

	aiTeamProfiles := GetOnlyAITeamRecruitingProfiles()

	coachMap := GetActiveCollegeCoachMap()

	for _, team := range aiTeamProfiles {
		if !team.IsAI {
			continue
		}
		coach := coachMap[team.ID]
		if coach.ID == 0 {
			continue
		}
		teamID := strconv.Itoa(int(team.ID))
		roster := GetCollegePlayersByTeamId(teamID)
		// Coaches should not be putting out promises if they are over the max
		if len(roster) > 12 {
			continue
		}
		for _, p := range roster {
			if p.TransferStatus > 1 || p.TransferStatus == 0 {
				continue
			}
			collegePlayerID := strconv.Itoa(int(p.ID))
			promise := GetCollegePromiseByCollegePlayerID(collegePlayerID, teamID)
			if promise.ID != 0 {
				continue
			}

			promiseOdds := getBasePromiseOdds(coach.TeambuildingPreference, coach.PromiseTendency)
			diceRoll := util.GenerateIntFromRange(1, 100)

			if diceRoll < promiseOdds {
				// Commit Promise
				promiseLevel := getPromiseLevel(coach.PromiseTendency)
				promiseWeight := "Medium"
				promiseType := ""
				benchmarkStr := ""
				promiseBenchmark := 0

				bias := p.RecruitingBias
				if bias == closeToHome {
					promiseType = "Home State Game"
					benchmarkStr = p.State
				} else if bias == immediateStart && p.Overall > 48 {
					promiseType = "Minutes"
					promiseBenchmark = p.PlaytimeExpectations
					switch promiseLevel {
					case 1:
						promiseBenchmark += 5
						if promiseBenchmark > p.Stamina {
							promiseBenchmark = p.Stamina - 1
						}
					case -1:
						promiseBenchmark -= 1
					}

					promiseWeight = getPromiseWeightByMinutesOrWins(promiseType, promiseBenchmark)
				} else if bias == nationalChampionshipContender || bias == richHistory {
					// Promise based on wins
					promiseBenchmark = 20
					promiseType = "Wins"
					switch promiseLevel {
					case 1:
						promiseBenchmark += 5
					case -1:
						promiseBenchmark -= 5
					}
					promiseWeight = getPromiseWeightByMinutesOrWins(promiseType, promiseBenchmark)
				}

				if promiseType == "" {
					continue
				}

				collegePromise := structs.CollegePromise{
					TeamID:          team.ID,
					CollegePlayerID: p.ID,
					PromiseType:     promiseType,
					PromiseWeight:   promiseWeight,
					Benchmark:       promiseBenchmark,
					BenchmarkStr:    benchmarkStr,
					IsActive:        true,
				}

				db.Create(&collegePromise)
			}
		}
	}
}

func GetAllCollegePromises() []structs.CollegePromise {
	db := dbprovider.GetInstance().GetDB()

	p := []structs.CollegePromise{}

	err := db.Find(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []structs.CollegePromise{}
		} else {
			log.Fatal(err)
		}
	}
	return p
}

func GetAllCollegePromisesByTeamID(teamID string) []structs.CollegePromise {
	db := dbprovider.GetInstance().GetDB()

	p := []structs.CollegePromise{}

	err := db.Where("team_id = ?", teamID).Find(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []structs.CollegePromise{}
		} else {
			log.Fatal(err)
		}
	}
	return p
}

func GetCollegePromiseByID(id string) structs.CollegePromise {
	db := dbprovider.GetInstance().GetDB()

	p := structs.CollegePromise{}

	err := db.Where("id = ?", id).Find(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return structs.CollegePromise{}
		} else {
			log.Fatal(err)
		}
	}
	return p
}

func GetCollegePromiseByCollegePlayerID(id, teamID string) structs.CollegePromise {
	db := dbprovider.GetInstance().GetDB()

	p := structs.CollegePromise{}

	err := db.Where("college_player_id = ? AND team_id = ?", id, teamID).Find(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return structs.CollegePromise{}
		} else {
			log.Fatal(err)
		}
	}
	return p
}

func CreatePromise(promise structs.CollegePromise) structs.CollegePromise {
	db := dbprovider.GetInstance().GetDB()
	id := strconv.Itoa(int(promise.ID))
	collegePlayerID := strconv.Itoa(int(promise.CollegePlayerID))
	profileID := strconv.Itoa(int(promise.TeamID))

	existingPromise := GetCollegePromiseByID(id)
	if existingPromise.ID != 0 && existingPromise.ID > 0 {
		existingPromise.Reactivate(promise.PromiseType, promise.PromiseWeight, promise.Benchmark)
		db.Save(&existingPromise)
		assignPromiseToProfile(db, collegePlayerID, profileID, existingPromise.ID)
		return existingPromise
	}

	db.Create(&promise)

	assignPromiseToProfile(db, collegePlayerID, profileID, promise.ID)

	return promise
}

func assignPromiseToProfile(db *gorm.DB, collegePlayerID, profileID string, id uint) {
	tpProfile := GetOnlyTransferPortalProfileByPlayerID(collegePlayerID, profileID)
	if tpProfile.ID > 0 {
		tpProfile.AssignPromise(id)
		db.Save(&tpProfile)
	}
}

func UpdatePromise(promise structs.CollegePromise) {
	db := dbprovider.GetInstance().GetDB()
	id := strconv.Itoa(int(promise.ID))
	existingPromise := GetCollegePromiseByID(id)
	existingPromise.UpdatePromise(promise.PromiseType, promise.PromiseWeight, promise.Benchmark)
	db.Save(&existingPromise)
}

func CancelPromise(id string) {
	db := dbprovider.GetInstance().GetDB()
	promise := GetCollegePromiseByID(id)
	promise.Deactivate()
	db.Save(&promise)
}

func EnterTheTransferPortal() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	// Get All Teams
	teams := GetAllActiveCollegeTeams()

	for _, t := range teams {
		teamID := strconv.Itoa(int(t.ID))
		roster := GetCollegePlayersByTeamId(teamID)

		for _, p := range roster {
			if p.TransferStatus != 1 {
				continue
			}

			playerID := strconv.Itoa(int(p.ID))

			promise := GetCollegePromiseByCollegePlayerID(playerID, teamID)
			if promise.ID == 0 {
				p.WillTransfer()
				db.Save(&p)
				continue
			}
			// 1-100
			baseFloor := getTransferFloor(p.TransferLikeliness)
			// 10, 20, 40, 60, 70
			promiseModifier := getPromiseFloor(promise.PromiseWeight)
			difference := baseFloor - promiseModifier

			diceRoll := util.GenerateIntFromRange(1, 100)

			// Lets say the difference is 40. 60-20.
			if diceRoll < difference {
				// If the dice roll is within the 40%. They leave.
				// Okay this makes sense.

				p.WillTransfer()

				// Create News Log
				message := "Breaking News! " + p.TeamAbbr + " " + strconv.Itoa(p.Stars) + " Star " + p.Position + " " + p.FirstName + " " + p.LastName + " has officially entered the transfer portal!"
				CreateNewsLog("CBB", message, "Transfer Portal", int(p.PreviousTeamID), ts)

				repository.SaveCollegePlayerRecord(p, db)
				repository.DeleteCollegePromise(promise, db, false)
				continue
			}

			// Create News Log
			message := "Breaking News! " + p.TeamAbbr + " " + strconv.Itoa(p.Stars) + " Star " + p.Position + " " + p.FirstName + " " + p.LastName + " has withdrawn their name from the transfer portal!"
			CreateNewsLog("CBB", message, "Transfer Portal", int(p.PreviousTeamID), ts)

			promise.MakePromise()
			repository.SaveCollegePromiseRecord(promise, db)
			p.WillStay()
			repository.SaveCollegePlayerRecord(p, db)
		}
	}

	ts.EnactPortalPhase()
	repository.SaveTimeStamp(ts, db)
}

func AddTransferPlayerToBoard(transferPortalProfileDto structs.TransferPortalProfile) structs.TransferPortalProfile {
	db := dbprovider.GetInstance().GetDB()

	portalProfile := GetOnlyTransferPortalProfileByPlayerID(strconv.Itoa(int(transferPortalProfileDto.CollegePlayerID)), strconv.Itoa(int(transferPortalProfileDto.ProfileID)))

	// If Recruit Already Exists
	if portalProfile.CollegePlayerID != 0 && portalProfile.ProfileID != 0 {
		portalProfile.Reactivate()
		db.Save(&portalProfile)
		return portalProfile
	}

	newProfileForRecruit := structs.TransferPortalProfile{
		SeasonID:           uint(transferPortalProfileDto.SeasonID),
		CollegePlayerID:    uint(transferPortalProfileDto.CollegePlayerID),
		ProfileID:          uint(transferPortalProfileDto.ProfileID),
		TeamAbbreviation:   transferPortalProfileDto.TeamAbbreviation,
		TotalPoints:        0,
		CurrentWeeksPoints: 0,
		SpendingCount:      0,
		RemovedFromBoard:   false,
	}

	db.Create(&newProfileForRecruit)

	return newProfileForRecruit
}

func RemovePlayerFromTransferPortalBoard(dto structs.TransferPortalProfile) {
	db := dbprovider.GetInstance().GetDB()

	playerID := strconv.Itoa(int(dto.CollegePlayerID))
	profileID := strconv.Itoa(int(dto.ProfileID))
	profile := GetOnlyTransferPortalProfileByPlayerID(playerID, profileID)

	profile.Deactivate()

	if profile.PromiseID.Int64 > 0 {
		promiseID := strconv.Itoa(int(profile.PromiseID.Int64))
		promise := GetCollegePromiseByID(promiseID)
		promise.Deactivate()
		profile.AssignPromise(0)
		repository.DeleteCollegePromise(promise, db, true)
	}

	repository.SaveTransferPortalProfile(profile, db)
}

func AllocatePointsToTransferPlayer(updateTransferPortalBoardDto structs.UpdateTransferPortalBoard) {
	db := dbprovider.GetInstance().GetDB()

	var teamId = strconv.Itoa(updateTransferPortalBoardDto.TeamID)
	var profile = GetOnlyTeamRecruitingProfileByTeamID(teamId)
	var portalProfiles = GetOnlyTransferPortalProfilesByTeamID(teamId)
	var updatedPlayers = updateTransferPortalBoardDto.Players

	currentPoints := 0

	for i := 0; i < len(portalProfiles); i++ {
		updatedRecruit := GetPlayerFromTransferPortalList(int(portalProfiles[i].CollegePlayerID), updatedPlayers)

		if portalProfiles[i].CurrentWeeksPoints != updatedRecruit.CurrentWeeksPoints {

			// Allocate Points to Profile
			currentPoints += updatedRecruit.CurrentWeeksPoints
			profile.AllocateSpentPoints(currentPoints)
			// If total not surpassed, allocate to the recruit and continue
			if profile.SpentPoints <= profile.WeeklyPoints {
				portalProfiles[i].AllocatePoints(updatedRecruit.CurrentWeeksPoints)
				fmt.Println("Saving recruit " + strconv.Itoa(int(portalProfiles[i].CollegePlayerID)))
			} else {
				panic("Error: Allocated more points for Profile " + strconv.Itoa(int(profile.TeamID)) + " than what is allowed.")
			}
			db.Save(&portalProfiles[i])
		} else {
			currentPoints += portalProfiles[i].CurrentWeeksPoints
			profile.AllocateSpentPoints(currentPoints)
		}
	}

	// Save profile
	db.Save(&profile)
}

func AICoachFillBoardsPhase() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.SeasonID))
	AITeams := GetOnlyAITeamRecruitingProfiles()
	// Shuffles the list of AI teams so that it's not always iterating from A-Z. Gives the teams at the lower end of the list a chance to recruit other croots
	rand.Shuffle(len(AITeams), func(i, j int) {
		AITeams[i], AITeams[j] = AITeams[j], AITeams[i]
	})
	transferPortalPlayers := GetTransferPortalPlayers()
	coachMap := GetActiveCollegeCoachMap()
	teamMap := GetCollegeTeamMap()
	regionMap := util.GetRegionMap()
	standingsMap := GetCollegeStandingsMap(seasonID)
	for _, teamProfile := range AITeams {
		if !teamProfile.IsAI {
			continue
		}
		if teamProfile.ID == 2 || teamProfile.ID == 8 || teamProfile.ID == 98 || teamProfile.ID == 106 {
			fmt.Println(teamProfile.TeamAbbr)
		}
		team := teamMap[teamProfile.ID]
		teamStandings := standingsMap[teamProfile.TeamID]
		teamID := strconv.Itoa(int(teamProfile.ID))
		coach := coachMap[teamProfile.ID]
		portalProfileMap := getTransferPortalProfileMapByTeamID(teamID)

		roster := GetCollegePlayersByTeamId(teamID)
		rosterSize := len(roster)
		// Roster sizes of 12 or higher should be ignored
		if rosterSize > 12 {
			continue
		}
		teamNeedsMap := make(map[string]bool)
		positionCount := make(map[string]int)

		if _, ok := teamNeedsMap["PG"]; !ok {
			teamNeedsMap["PG"] = true
		}
		if _, ok := teamNeedsMap["SG"]; !ok {
			teamNeedsMap["SG"] = true
		}
		if _, ok := teamNeedsMap["SF"]; !ok {
			teamNeedsMap["SF"] = true
		}
		if _, ok := teamNeedsMap["PF"]; !ok {
			teamNeedsMap["PF"] = true
		}
		if _, ok := teamNeedsMap["C"]; !ok {
			teamNeedsMap["C"] = true
		}

		if _, ok := positionCount["PG"]; !ok {
			positionCount["PG"] = 0
		}
		if _, ok := positionCount["SG"]; !ok {
			positionCount["SG"] = 0
		}
		if _, ok := positionCount["SF"]; !ok {
			positionCount["SF"] = 0
		}
		if _, ok := positionCount["PF"]; !ok {
			positionCount["PF"] = 0
		}
		if _, ok := positionCount["C"]; !ok {
			positionCount["C"] = 0
		}

		//
		for _, r := range roster {
			if r.WillDeclare {
				continue
			}
			positionCount[r.Position] += 1
		}

		if positionCount["PG"] >= 3 {
			teamNeedsMap["PG"] = false
		} else if positionCount["SG"] >= 4 {
			teamNeedsMap["SG"] = false
		} else if positionCount["SF"] >= 4 {
			teamNeedsMap["SF"] = false
		} else if positionCount["PF"] >= 4 {
			teamNeedsMap["PF"] = false
		} else if positionCount["C"] >= 3 {
			teamNeedsMap["C"] = false
		}
		for _, tp := range transferPortalPlayers {
			if !teamNeedsMap[tp.Position] || tp.PreviousTeamID == team.ID || portalProfileMap[tp.ID].CollegePlayerID == tp.ID || portalProfileMap[tp.ID].ID > 0 {
				continue
			}

			// Put together a player prestige rating to use as a qualifier on which teams will target specific players. Ideally more experienced coaches will be able to target higher rated players
			// playerPrestige := getPlayerPrestigeRating(tp.Stars, tp.Overall)
			// if coach.Prestige < playerPrestige {
			// 	continue
			// }
			bias := tp.RecruitingBias
			biasMod := 0
			postSeasonStatus := teamStandings.PostSeasonStatus
			if bias == richHistory {
				// Get multiple season standings
				teamHistory := GetStandingsHistoryByTeamID(teamID)
				averageWins := getAverageWins(teamHistory)
				biasMod += averageWins
				if teamProfile.AIQuality == "Blue Blood" {
					biasMod += 20
				}
			} else if bias == nationalChampionshipContender {
				switch postSeasonStatus {
				case "Sweet 16", "Elite 8":
					biasMod += 10
				case "Final Four":
					biasMod += 15
				case "National Championship Participant":
					biasMod += 20
				case "National Champions":
					biasMod += 25
				}
			} else if bias == upcomingTeam {
				biasMod += teamStandings.TotalWins
				if teamProfile.AIQuality == "Mid-Major" || teamProfile.AIQuality == "Cinderella" {
					biasMod += 15
				}
			} else if bias == differentState && tp.State != team.State {
				biasMod += 15
				if regionMap[tp.State] == team.State {
					biasMod += 5
				}
			} else if bias == specificCoach && tp.LegacyID == coach.ID {
				biasMod += 25
			} else if bias == legacy && tp.LegacyID == team.ID {
				biasMod += 25
			}

			diceRoll := util.GenerateIntFromRange(1, 50)
			if teamProfile.ID == 98 || teamProfile.ID == 106 {
				diceRoll -= 20
			}
			if diceRoll < biasMod {
				// Add Player to Board

				portalProfile := structs.TransferPortalProfile{
					ProfileID:        teamProfile.ID,
					CollegePlayerID:  tp.ID,
					SeasonID:         ts.SeasonID,
					TeamAbbreviation: teamProfile.TeamAbbr,
				}

				db.Create(&portalProfile)
			}
		}
	}
}

func AICoachAllocateAndPromisePhase() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	AITeams := GetOnlyAITeamRecruitingProfiles()
	transferPortalPlayerMap := GetCollegePlayerMap()
	coachMap := GetActiveCollegeCoachMap()
	regionMap := util.GetRegionMap()
	// Shuffles the list of AI teams so that it's not always iterating from A-Z. Gives the teams at the lower end of the list a chance to recruit other croots
	rand.Shuffle(len(AITeams), func(i, j int) {
		AITeams[i], AITeams[j] = AITeams[j], AITeams[i]
	})

	for _, teamProfile := range AITeams {
		if teamProfile.SpentPoints >= teamProfile.WeeklyPoints {
			continue
		}
		teamID := strconv.Itoa(int(teamProfile.ID))
		portalProfiles := GetTransferPortalProfilesByTeamID(teamID)
		for _, p := range portalProfiles {
			if p.LockProfile && p.CurrentWeeksPoints > 0 {
				points := p.CurrentWeeksPoints
				teamProfile.AIAllocateSpentPoints(points * -1)
				p.Deactivate()
				repository.SaveTransferPortalProfile(p, db)
			}
		}
	}

	for _, teamProfile := range AITeams {
		if teamProfile.SpentPoints >= teamProfile.WeeklyPoints {
			continue
		}

		teamID := strconv.Itoa(int(teamProfile.ID))
		roster := GetCollegePlayersByTeamId(teamID)
		if len(roster) > 12 {
			continue
		}

		teamNeedsMap := make(map[string]bool)
		positionCount := make(map[string]int)

		if _, ok := teamNeedsMap["PG"]; !ok {
			teamNeedsMap["PG"] = true
		}
		if _, ok := teamNeedsMap["SG"]; !ok {
			teamNeedsMap["SG"] = true
		}
		if _, ok := teamNeedsMap["SF"]; !ok {
			teamNeedsMap["SF"] = true
		}
		if _, ok := teamNeedsMap["PF"]; !ok {
			teamNeedsMap["PF"] = true
		}
		if _, ok := teamNeedsMap["C"]; !ok {
			teamNeedsMap["C"] = true
		}

		if _, ok := positionCount["PG"]; !ok {
			positionCount["PG"] = 0
		}
		if _, ok := positionCount["SG"]; !ok {
			positionCount["SG"] = 0
		}
		if _, ok := positionCount["SF"]; !ok {
			positionCount["SF"] = 0
		}
		if _, ok := positionCount["PF"]; !ok {
			positionCount["PF"] = 0
		}
		if _, ok := positionCount["C"]; !ok {
			positionCount["C"] = 0
		}

		for _, r := range roster {
			if r.WillDeclare {
				continue
			}
			positionCount[r.Position] += 1
		}

		if positionCount["PG"] >= 3 {
			teamNeedsMap["PG"] = false
		} else if positionCount["SG"] >= 4 {
			teamNeedsMap["SG"] = false
		} else if positionCount["SF"] >= 4 {
			teamNeedsMap["SF"] = false
		} else if positionCount["PF"] >= 4 {
			teamNeedsMap["PF"] = false
		} else if positionCount["C"] >= 3 {
			teamNeedsMap["C"] = false
		}

		portalProfiles := GetTransferPortalProfilesByTeamID(teamID)
		for _, profile := range portalProfiles {
			if profile.CurrentWeeksPoints > 0 || profile.RemovedFromBoard {
				continue
			}
			tp := transferPortalPlayerMap[profile.CollegePlayerID]
			// If player has already signed or if the position has been fulfilled
			if tp.TeamID > 0 || tp.TransferStatus == 0 || tp.ID == 0 || !teamNeedsMap[tp.Position] {
				points := profile.CurrentWeeksPoints
				teamProfile.AIAllocateSpentPoints(points * -1)
				profile.Deactivate()
				repository.SaveTransferPortalProfile(profile, db)
				continue
			}
			playerID := strconv.Itoa(int(profile.CollegePlayerID))
			pointsRemaining := teamProfile.WeeklyPoints - teamProfile.SpentPoints
			if teamProfile.SpentPoints >= teamProfile.WeeklyPoints || pointsRemaining <= 0 || (pointsRemaining < 1 && pointsRemaining > 0) {
				break
			}

			removePlayerFromBoard := false
			num := 0

			profiles := GetTransferPortalProfilesByPlayerID(playerID)
			leadingTeamVal := util.IsAITeamContendingForPortalPlayer(profiles)
			if profile.CurrentWeeksPoints > 0 && profile.TotalPoints+float64(profile.CurrentWeeksPoints) >= float64(leadingTeamVal)*0.66 {
				// continue, leave everything alone
				continue
			} else if profile.CurrentWeeksPoints > 0 && profile.TotalPoints+float64(profile.CurrentWeeksPoints) < float64(leadingTeamVal)*0.66 {
				profile.Deactivate()
				db.Save(&profile)
				continue
			}

			maxChance := 2
			if ts.CollegeWeek > 3 {
				maxChance = 4
			}
			chance := util.GenerateIntFromRange(1, maxChance)
			if (chance < 2 && ts.TransferPortalPhase <= 3) || (chance < 4 && ts.TransferPortalPhase > 3) {
				continue
			}
			coach := coachMap[teamProfile.TeamID]

			min := coach.PointMin
			max := coach.PointMax
			if max > 10 {
				max = 10
			}
			num = util.GenerateIntFromRange(min, max)
			if num > pointsRemaining {
				num = pointsRemaining
			}

			if float64(num)+profile.TotalPoints < float64(leadingTeamVal)*0.66 {
				removePlayerFromBoard = true
			}
			if leadingTeamVal < 14 {
				removePlayerFromBoard = false
			}

			if removePlayerFromBoard {
				points := profile.CurrentWeeksPoints
				teamProfile.AIAllocateSpentPoints(points * -1)
				profile.Deactivate()
				repository.SaveTransferPortalProfile(profile, db)
				continue
			}
			profile.AllocatePoints(num)

			// Generate Promise based on coach bias
			if profile.PromiseID.Int64 == 0 && !profile.RolledOnPromise {
				promiseOdds := getBasePromiseOdds(coach.TeambuildingPreference, coach.PromiseTendency)
				diceRoll := util.GenerateIntFromRange(1, 100)

				if diceRoll < promiseOdds {
					// Commit Promise
					promiseLevel := getPromiseLevel(coach.PromiseTendency)
					promiseWeight := "Medium"
					promiseType := ""
					benchmarkStr := ""
					promiseBenchmark := 0

					bias := tp.RecruitingBias
					if bias == closeToHome && (teamProfile.State == tp.State) || teamProfile.Region == regionMap[tp.State] {
						promiseType = "Home State Game"
						benchmarkStr = tp.State
					} else if bias == immediateStart && tp.Overall > 55 {
						promiseType = "Minutes"
						promiseBenchmark = tp.PlaytimeExpectations
						switch promiseLevel {
						case 1:
							promiseBenchmark += 5
							if promiseBenchmark > tp.Stamina {
								promiseBenchmark = tp.Stamina - 1
							}
						case -1:
							promiseBenchmark -= 1
						}

						promiseWeight = getPromiseWeightByMinutesOrWins(promiseType, promiseBenchmark)
					} else if bias == nationalChampionshipContender || bias == richHistory {
						// Promise based on wins
						promiseBenchmark = 20
						promiseType = "Wins"
						switch promiseLevel {
						case 1:
							promiseBenchmark += 5
						case -1:
							promiseBenchmark -= 5
						}
						promiseWeight = getPromiseWeightByMinutesOrWins(promiseType, promiseBenchmark)
					} else if bias == legacy && tp.LegacyID == teamProfile.TeamID {
						promiseType = "Legacy"
					} else if bias == specificCoach && tp.LegacyID == coach.ID {
						promiseType = "Specific Coach"
					} else if bias == differentState && teamProfile.State != tp.State {
						promiseType = "Different State"
						promiseWeight = "Low"
					}

					collegePromise := structs.CollegePromise{
						TeamID:          teamProfile.TeamID,
						CollegePlayerID: tp.ID,
						PromiseType:     promiseType,
						PromiseWeight:   promiseWeight,
						Benchmark:       promiseBenchmark,
						BenchmarkStr:    benchmarkStr,
						IsActive:        true,
					}

					repository.CreateCollegePromiseRecord(collegePromise, db)
				}

				profile.ToggleRolledOnPromise()
			}
			// Save Profile
			repository.SaveTransferPortalProfile(profile, db)
		}
	}
}

func SyncTransferPortal() {
	db := dbprovider.GetInstance().GetDB()
	//GetCurrentWeek
	ts := GetTimestamp()
	// Use IsRecruitingLocked to lock the TP when not in use
	teamProfileMap := GetTeamProfileMap()
	transferPortalPlayers := GetTransferPortalPlayers()
	transferPortalProfileMap := MakeFullTransferPortalProfileMap(transferPortalPlayers)
	rosterMap := GetFullTeamRosterWithCrootsMap()

	if !ts.IsRecruitingLocked {
		ts.ToggleLockRecruiting()
		db.Save(&ts)
	}

	for _, portalPlayer := range transferPortalPlayers {

		// Skip over players that have already transferred
		if portalPlayer.TransferStatus != 2 || portalPlayer.TeamID > 0 {
			continue
		}

		portalProfiles := transferPortalProfileMap[portalPlayer.ID]
		if len(portalProfiles) == 0 && ts.TransferPortalRound < 10 {
			continue
		}

		// If no one has a profile on them during round 10
		if len(portalProfiles) == 0 && ts.TransferPortalRound == 10 {
			roster := rosterMap[portalPlayer.PreviousTeamID]
			if len(roster) > 15 {
				continue
			}
			rosterMap[portalPlayer.PreviousTeamID] = append(rosterMap[portalPlayer.PreviousTeamID], portalPlayer)
			portalPlayer.WillReturn()
			db.Save(&portalPlayer)
			continue
		}

		totalPointsOnPlayer := 0.0
		eligiblePointThreshold := 0.0
		readyToSign := false
		minSpendingCount := 100
		maxSpendingCount := 0
		signingMinimum := 0.66
		teamCount := 0
		eligibleTeams := []structs.TransferPortalProfile{}

		for i := range portalProfiles {
			if portalProfiles[i].UpdatedAt.String() >= "2026-01-07 00:00:00" {
				continue
			}
			promiseID := strconv.Itoa(int(portalProfiles[i].PromiseID.Int64))

			promise := GetCollegePromiseByID(promiseID)

			multiplier := getMultiplier(promise)
			portalProfiles[i].AddPointsToTotal(multiplier)
		}

		sort.Slice(portalProfiles, func(i, j int) bool {
			return portalProfiles[i].TotalPoints > portalProfiles[j].TotalPoints
		})

		for i := range portalProfiles {
			if portalProfiles[i].UpdatedAt.String() >= "2026-01-07 00:00:00" {
				continue
			}
			roster := rosterMap[portalProfiles[i].ProfileID]
			if len(roster) > 15 {
				continue
			}
			if eligiblePointThreshold == 0.0 {
				eligiblePointThreshold = portalProfiles[i].TotalPoints * signingMinimum
			}
			if portalProfiles[i].TotalPoints >= eligiblePointThreshold {
				if portalProfiles[i].SpendingCount < minSpendingCount {
					minSpendingCount = portalProfiles[i].SpendingCount
				}
				if portalProfiles[i].SpendingCount > maxSpendingCount {
					maxSpendingCount = portalProfiles[i].SpendingCount
				}
				eligibleTeams = append(eligibleTeams, portalProfiles[i])
				totalPointsOnPlayer += portalProfiles[i].TotalPoints
				teamCount += 1
			}
		}

		if (teamCount >= 1 && minSpendingCount >= 2) || (teamCount >= 1 && ts.TransferPortalRound == 10) {
			// threshold met
			readyToSign = true
		}
		var winningTeamID uint = 0
		if readyToSign {
			var odds float64 = 0

			for winningTeamID == 0 {
				percentageOdds := rand.Float64() * (totalPointsOnPlayer)
				currentProbability := 0.0
				for _, profile := range eligibleTeams {
					currentProbability += profile.TotalPoints
					if percentageOdds <= currentProbability {
						// WINNING TEAM
						winningTeamID = profile.ProfileID
						odds = profile.TotalPoints / totalPointsOnPlayer * 100
						break
					}
				}

				if winningTeamID > 0 {
					winningTeamIDSTR := strconv.Itoa(int(winningTeamID))
					promise := GetCollegePromiseByCollegePlayerID(strconv.Itoa(int(portalPlayer.ID)), winningTeamIDSTR)
					if promise.ID > 0 {
						promise.MakePromise()
						db.Save(&promise)
					}

					teamProfile := teamProfileMap[winningTeamIDSTR]
					currentRoster := rosterMap[teamProfile.ID]
					if len(currentRoster) < 15 {
						portalPlayer.SignWithNewTeam(teamProfile.ID, teamProfile.TeamAbbr)
						message := portalPlayer.FirstName + " " + portalPlayer.LastName + ", " + strconv.Itoa(portalPlayer.Stars) + " star " + portalPlayer.Position + " from " + portalPlayer.PreviousTeam + " has signed with " + portalPlayer.TeamAbbr + " with " + strconv.Itoa(int(odds)) + " percent odds."
						CreateNewsLog("CBB", message, "Transfer Portal", int(winningTeamID), ts)
						fmt.Println("Created new log!")
						// Add player to existing roster map
						rosterMap[teamProfile.ID] = append(rosterMap[teamProfile.ID], portalPlayer)
						for i := range portalProfiles {
							if portalProfiles[i].ProfileID == winningTeamID {
								portalProfiles[i].SignPlayer()
								promise := GetCollegePromiseByCollegePlayerID(strconv.Itoa(int(portalPlayer.ID)), strconv.Itoa(int(winningTeamID)))
								if promise.ID > 0 {
									promise.MakePromise()
									repository.SaveCollegePromiseRecord(promise, db)
								}
								break
							}
						}

					} else {
						// Filter out profile
						eligibleTeams = util.FilterOutPortalProfile(eligibleTeams, winningTeamID)
						winningTeamID = 0
						if len(eligibleTeams) == 0 {
							break
						}

						totalPointsOnPlayer = 0
						for _, p := range eligibleTeams {
							totalPointsOnPlayer += p.TotalPoints
						}
					}
				}
			}

		}
		for _, p := range portalProfiles {
			if winningTeamID > 0 && p.ID != winningTeamID {
				p.RemovePromise()
				p.Lock()
			}
			if winningTeamID > 0 || p.SpendingCount > 0 {
				repository.SaveTransferPortalProfile(p, db)
			}
			fmt.Println("Save transfer portal profile from " + portalPlayer.TeamAbbr + " towards " + portalPlayer.FirstName + " " + portalPlayer.LastName)
			if winningTeamID > 0 && p.ProfileID != winningTeamID {
				promise := GetCollegePromiseByCollegePlayerID(strconv.Itoa(int(portalPlayer.ID)), strconv.Itoa(int(p.ProfileID)))
				if promise.ID > 0 {
					repository.DeleteCollegePromise(promise, db, false)
				}
			}
		}
		// Save Recruit
		if portalPlayer.TeamID > 0 {
			repository.SaveCollegePlayerRecord(portalPlayer, db)
		}
	}

	ts.IncrementTransferPortalRound()
	repository.SaveTimeStamp(ts, db)
}

// At end of season, sync through promises to confirm if promises were made
func SyncPromises() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.SeasonID))
	teamProfileMap := GetTeamProfileMap()
	activePromises := GetAllCollegePromises()
	collegePlayerMap := GetCollegePlayerMap()
	historicPlayerMap := GetHistoricCollegePlayerMap()
	standingsMap := GetCollegeStandingsMap(seasonID)
	seasonStatsMap := GetCollegePlayerSeasonStatMap(seasonID)

	for _, promise := range activePromises {
		if !promise.IsActive || !promise.PromiseMade {
			continue
		}
		isHistoric := false
		benchMarkStr := ""
		result := ""
		player := collegePlayerMap[promise.CollegePlayerID]
		if player.ID == 0 {
			player = historicPlayerMap[promise.CollegePlayerID]
			if player.ID == 0 {
				continue
			}
			isHistoric = true
		}
		// If player is already going to portal, carry on!
		if player.TransferStatus == 2 {
			// Remove promise since there was likely a preceding promise
			repository.DeleteCollegePromise(promise, db, false)
			continue
		}
		teamID := strconv.Itoa(int(promise.TeamID))
		team := teamProfileMap[teamID]

		seasonStats := seasonStatsMap[promise.CollegePlayerID]
		if promise.PromiseType == "Wins" {
			benchMarkStr = strconv.Itoa(int(promise.Benchmark))
			standings := standingsMap[promise.TeamID]
			result = strconv.Itoa(int(standings.TotalWins))
			if standings.TotalWins >= promise.Benchmark {
				promise.FulfillPromise()
			}
		} else if promise.PromiseType == "Minutes" {
			benchMarkStr = strconv.Itoa(int(promise.Benchmark))
			result = util.ConvertFloatToString(seasonStats.MinutesPerGame)
			if seasonStats.MinutesPerGame >= float64(promise.Benchmark) {
				promise.FulfillPromise()
			}
		} else if promise.PromiseType == "Home State Game" || promise.PromiseType == "Different State" {
			// Loop through games
			benchMarkStr = promise.BenchmarkStr
			result = "Did not play game in requested state."
			games := GetMatchesByTeamIdAndSeasonId(teamID, seasonID)
			for _, game := range games {
				stateKey := util.GetStateKey(promise.BenchmarkStr)
				if game.State == stateKey || game.State == promise.BenchmarkStr {
					result = ""
					promise.FulfillPromise()
					break
				}
			}
		} else if promise.PromiseType == "No Redshirt" {
			result = "Was Redshirted"
			if !player.IsRedshirting {
				result = ""
				promise.FulfillPromise()
			}
		} else if promise.PromiseType == "National Championship" {
			result = "Did not win the Natty."
			standings := standingsMap[promise.TeamID]
			if standings.PostSeasonStatus == "National Champion" {
				result = ""
				promise.FulfillPromise()
			}
		} else if promise.PromiseType == "Conference Championship" {
			result = "Did not win Conference Championship"
			standings := standingsMap[promise.TeamID]
			if standings.IsConferenceChampion {
				result = ""
				promise.FulfillPromise()
			}
		} else if promise.PromiseType == "Specific Coach" {
			// Fulfill for now, will need to adjust value
			promise.FulfillPromise()
		} else if promise.PromiseType == "Playoffs" {
			standings := standingsMap[promise.TeamID]
			postSeasonStatus := standings.PostSeasonStatus
			// postSeasonStatus has substring "Round of" or postSeasonStatus == "Sweet 16" or "Elite 8" or "Final 4" or contains "National Champion", fullfill
			if strings.Contains(postSeasonStatus, "Round of") || postSeasonStatus == "Sweet 16" || postSeasonStatus == "Elite 8" || postSeasonStatus == "Final Four" || strings.Contains(postSeasonStatus, "National Champion") {
				result = ""
				promise.FulfillPromise()
			}
		} else if promise.PromiseType == "Elite 8" {
			standings := standingsMap[promise.TeamID]
			postSeasonStatus := standings.PostSeasonStatus
			// postSeasonStatus has substring "Round of" or postSeasonStatus == "Sweet 16" or "Elite 8" or "Final 4" or contains "National Champion", fullfill
			if postSeasonStatus == "Elite 8" || postSeasonStatus == "Final Four" || strings.Contains(postSeasonStatus, "National Champion") {
				result = ""
				promise.FulfillPromise()
			}
		} else if promise.PromiseType == "Sweet 16" {
			standings := standingsMap[promise.TeamID]
			postSeasonStatus := standings.PostSeasonStatus
			// postSeasonStatus has substring "Round of" or postSeasonStatus == "Sweet 16" or "Elite 8" or "Final 4" or contains "National Champion", fullfill
			if postSeasonStatus == "Sweet 16" || postSeasonStatus == "Elite 8" || postSeasonStatus == "Final Four" || strings.Contains(postSeasonStatus, "National Champion") {
				result = ""
				promise.FulfillPromise()
			}
		} else if promise.PromiseType == "Final Four" || promise.PromiseType == "Final 4" {
			standings := standingsMap[promise.TeamID]
			postSeasonStatus := standings.PostSeasonStatus
			// postSeasonStatus has substring "Round of" or postSeasonStatus == "Sweet 16" or "Elite 8" or "Final 4" or contains "National Champion", fullfill
			if postSeasonStatus == "Final Four" || strings.Contains(postSeasonStatus, "National Champion") {
				result = ""
				promise.FulfillPromise()
			}
		}
		weightValue := getPromiseWeightValue(!promise.IsFullfilled, promise.PromiseWeight)
		team.AdjustPortalReputation(weightValue)
		repository.SaveCBBTeamRecruitingProfile(*team, db)
		if !promise.IsFullfilled && !isHistoric {
			message := "Breaking News! " + player.TeamAbbr + " " + player.FirstName + " " + player.LastName + " will be re-entering the portal after a promise was broken! Promise: " + promise.PromiseType + " | Expected: " + benchMarkStr + " | Result: " + result
			player.WillTransfer()
			repository.SaveCollegePlayerRecord(player, db)
			CreateNewsLog("CBB", message, "Portal", int(team.TeamID), ts)
		}
		repository.DeleteCollegePromise(promise, db, false)
	}
}

func getPromiseWeightValue(isPenalty bool, weight string) int {
	switch weight {
	case "Low":
		if isPenalty {
			return -5
		}
		return 3
	case "Very Low":
		if isPenalty {
			return -3
		}
		return 1
	case "Medium":
		if isPenalty {
			return -10
		}
		return 8
	case "High":
		if isPenalty {
			return -20
		}
		return 15
	case "Very High":
		if isPenalty {
			return -30
		}
		return 20
	}
	return 0
}

func GetPromisesByTeamID(teamID string) []structs.CollegePromise {
	db := dbprovider.GetInstance().GetDB()

	var promises []structs.CollegePromise

	db.Where("team_id = ?", teamID).Find(&promises)

	return promises
}

func GetOnlyTransferPortalProfilesByTeamID(teamID string) []structs.TransferPortalProfile {
	db := dbprovider.GetInstance().GetDB()

	var profiles []structs.TransferPortalProfile

	db.Where("profile_id = ?", teamID).Find(&profiles)

	return profiles
}

func GetTransferPortalProfilesByPlayerID(playerID string) []structs.TransferPortalProfile {
	db := dbprovider.GetInstance().GetDB()

	var profiles []structs.TransferPortalProfile

	db.Where("college_player_id = ? AND removed_from_board = ?", playerID, false).Find(&profiles)

	return profiles
}

func GetTransferPortalProfilesByPlayerIDs(playerID []string) []structs.TransferPortalProfile {
	db := dbprovider.GetInstance().GetDB()

	var profiles []structs.TransferPortalProfile

	db.Where("college_player_id in (?) AND removed_from_board = ?", playerID, false).Find(&profiles)

	return profiles
}

func GetTransferPortalProfilesForPage(teamID string) []structs.TransferPortalProfileResponse {
	db := dbprovider.GetInstance().GetDB()

	var profiles []structs.TransferPortalProfile
	var response []structs.TransferPortalProfileResponse
	err := db.Where("profile_id = ? AND removed_from_board = ?", teamID, false).Find(&profiles).Error
	collegePlayers := GetAllCollegePlayers()
	collegePlayerMap := MakeCollegePlayerMap(collegePlayers)
	if err != nil {
		log.Fatalln("Error!: ", err)
	}

	for _, p := range profiles {
		if p.RemovedFromBoard {
			continue
		}
		cpResponse := structs.TransferPlayerResponse{}
		cp := collegePlayerMap[p.CollegePlayerID]
		ovr := util.GetPlayerOverallGrade(cp.Overall)
		cpResponse.Map(cp, ovr)

		pResponse := structs.TransferPortalProfileResponse{
			ID:                    p.ID,
			SeasonID:              p.SeasonID,
			CollegePlayerID:       p.CollegePlayerID,
			ProfileID:             p.ProfileID,
			PromiseID:             uint(p.PromiseID.Int64),
			TeamAbbreviation:      p.TeamAbbreviation,
			TotalPoints:           p.TotalPoints,
			CurrentWeeksPoints:    p.CurrentWeeksPoints,
			PreviouslySpentPoints: p.PreviouslySpentPoints,
			SpendingCount:         p.SpendingCount,
			RemovedFromBoard:      p.RemovedFromBoard,
			RolledOnPromise:       p.RolledOnPromise,
			LockProfile:           p.LockProfile,
			IsSigned:              p.IsSigned,
			Recruiter:             p.Recruiter,
			CollegePlayer:         cpResponse,
		}

		response = append(response, pResponse)

	}

	return response
}

func GetActiveTransferPortalProfiles() []structs.TransferPortalProfile {
	db := dbprovider.GetInstance().GetDB()

	var profiles []structs.TransferPortalProfile

	db.Where("removed_from_board = ?", false).Find(&profiles)

	return profiles
}

func GetTransferPortalProfilesByTeamID(teamID string) []structs.TransferPortalProfile {
	db := dbprovider.GetInstance().GetDB()

	var profiles []structs.TransferPortalProfile

	db.Where("profile_id = ? AND removed_from_board = ?", teamID, false).Find(&profiles)

	return profiles
}

func GetOnlyTransferPortalProfileByID(tppID string) structs.TransferPortalProfile {
	db := dbprovider.GetInstance().GetDB()

	var profile structs.TransferPortalProfile

	db.Where("id = ?", tppID).Find(&profile)

	return profile
}

func GetOnlyTransferPortalProfileByPlayerID(playerId, teamId string) structs.TransferPortalProfile {
	db := dbprovider.GetInstance().GetDB()

	var profile structs.TransferPortalProfile

	db.Where("college_player_id = ? and profile_id = ?", playerId, teamId).Find(&profile)

	return profile
}

func GetTransferPortalData(teamID string) structs.TransferPortalResponse {
	var waitgroup sync.WaitGroup
	waitgroup.Add(5)
	profileChan := make(chan structs.TeamRecruitingProfile)
	playersChan := make(chan []structs.TransferPlayerResponse)
	boardChan := make(chan []structs.TransferPortalProfileResponse)
	promiseChan := make(chan []structs.CollegePromise)
	teamsChan := make(chan []structs.Team)

	go func() {
		waitgroup.Wait()
		close(profileChan)
		close(playersChan)
		close(boardChan)
		close(promiseChan)
		close(teamsChan)
	}()

	go func() {
		defer waitgroup.Done()
		profile := GetOnlyTeamRecruitingProfileByTeamID(teamID)
		profileChan <- profile
	}()
	go func() {
		defer waitgroup.Done()
		tpPlayers := GetTransferPortalPlayersForPage()
		playersChan <- tpPlayers
	}()
	go func() {
		defer waitgroup.Done()
		tpProfiles := GetTransferPortalProfilesForPage(teamID)
		boardChan <- tpProfiles
	}()
	go func() {
		defer waitgroup.Done()
		cPromises := GetPromisesByTeamID(teamID)
		promiseChan <- cPromises
	}()
	go func() {
		defer waitgroup.Done()
		cTeams := GetAllActiveCollegeTeams()
		teamsChan <- cTeams
	}()

	teamProfile := <-profileChan
	players := <-playersChan
	board := <-boardChan
	promises := <-promiseChan
	teams := <-teamsChan

	return structs.TransferPortalResponse{
		Team:         teamProfile,
		Players:      players,
		TeamBoard:    board,
		TeamPromises: promises,
		TeamList:     teams,
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

func getTransferPortalProfileMap(players []structs.CollegePlayer) map[uint][]structs.TransferPortalProfile {
	portalMap := make(map[uint][]structs.TransferPortalProfile)
	var mu sync.Mutex     // to safely update the map
	var wg sync.WaitGroup // to wait for all goroutines to finish
	semaphore := make(chan struct{}, 10)
	for _, p := range players {
		semaphore <- struct{}{}
		wg.Add(1)
		go func(c structs.CollegePlayer) {
			defer wg.Done()
			playerID := strconv.Itoa(int(c.ID))
			portalProfiles := GetTransferPortalProfilesByPlayerID(playerID)
			mu.Lock()
			portalMap[c.ID] = portalProfiles
			mu.Unlock()

			<-semaphore
		}(p)
	}
	wg.Wait()
	close(semaphore)
	return portalMap
}

// GetTransferFloor -- Get the Base Floor to determine if a player will transfer or not based on a promise
func getTransferFloor(likeliness string) int {
	min := 25
	max := 100
	switch likeliness {
	case "Low":
		max = 40
	case "Medium":
		min = 45
		max = 70
	default:
		min = 75
	}

	return util.GenerateIntFromRange(min, max)
}

// getPromiseFloor -- Get the modifier towards the floor value above
func getPromiseFloor(weight string) int {
	if weight == "Very Low" {
		return 10
	}
	if weight == "Low" {
		return 20
	}
	if weight == "Medium" {
		return 40
	}
	if weight == "High" {
		return 60
	}
	return util.GenerateIntFromRange(70, 80)
}

func getPromiseWeightByMinutesOrWins(category string, benchmark int) string {
	weight := "Medium"
	if category == "Minutes" {
		if benchmark <= 40 {
			weight = "Very High"
		}
		if benchmark <= 25 {
			weight = "High"
		}
		if benchmark <= 20 {
			weight = "Medium"
		}
		if benchmark <= 10 {
			weight = "Low"
		}
		if benchmark <= 5 {
			weight = "Very Low"
		}
	} else {
		// Wins
		if benchmark <= 40 {
			weight = "Extremely High"
		}
		if benchmark <= 30 {
			weight = "Very High"
		}
		if benchmark <= 25 {
			weight = "High"
		}
		if benchmark <= 20 {
			weight = "Medium"
		}
		if benchmark <= 15 {
			weight = "Low"
		}
		if benchmark <= 10 {
			weight = "Very Low"
		}
		if benchmark <= 5 {
			weight = "Extremely Low"
		}
	}

	return weight
}

func getPlayerPrestigeRating(stars, overall int) int {
	prestige := 1

	starMod := stars / 2
	if starMod <= 0 {
		starMod = 1
	}

	overallMod := overall / 10
	if overallMod <= 0 {
		overallMod = 1
	}

	return prestige + starMod + overallMod
}

func getAverageWins(standings []structs.CollegeStandings) int {
	wins := 0
	for _, s := range standings {
		wins += s.TotalWins
	}

	totalStandings := len(standings)
	if totalStandings > 0 {
		wins = wins / len(standings)
	}

	return wins
}

func getBasePromiseOdds(tbPref, ptTendency string) int {
	promiseOdds := 50
	if tbPref == "Recruiting" {
		promiseOdds += 20
	} else if ptTendency == "Transfer" {
		promiseOdds -= 20
	}

	return promiseOdds
}

func getTransferStatus(weight int) string {
	if weight < 20 {
		return "Low"
	}
	if weight < 40 {
		return "Medium"
	}
	return "High"
}

func getPromiseLevel(pt string) int {
	promiseLevel := 0
	switch pt {
	case "Over-Promise":
		promiseLevel = 1
	case "Under-Promise":
		promiseLevel = -1
	}
	return promiseLevel
}

func getMultiplier(pr structs.CollegePromise) float64 {
	if pr.ID == 0 {
		return 1
	}
	weight := pr.PromiseWeight
	switch weight {
	case "Why even try?":
		return 0.5
	case "Extremely Low":
		return 1.01
	case "Very Low":
		return 1.05
	case "Low":
		return 1.1
	case "Medium":
		return 1.3
	case "High":
		return 1.5
	case "Very High":
		return 1.75
	case "Extremely High":
		return 2
	case "If you make this promise then you better win it!":
		return 2.25
	}
	// Default
	return 1
}

func GetPlayerFromTransferPortalList(id int, profiles []structs.TransferPortalProfile) structs.TransferPortalProfile {
	var profile structs.TransferPortalProfile

	for i := 0; i < len(profiles); i++ {
		if profiles[i].CollegePlayerID == uint(id) {
			profile = profiles[i]
			break
		}
	}

	return profile
}

func GetTransferScoutingDataByPlayerID(id string) structs.ScoutingDataResponse {
	ts := GetTimestamp()

	seasonID := ts.SeasonID
	seasonIDSTR := strconv.Itoa(int(seasonID))

	draftee := GetCollegePlayerByPlayerID(id)

	seasonStats := GetPlayerSeasonStatsByPlayerID(id, seasonIDSTR)
	teamID := strconv.Itoa(int(draftee.PreviousTeamID))
	collegeStandings := GetStandingsRecordByTeamID(teamID, seasonIDSTR)

	return structs.ScoutingDataResponse{
		DrafteeSeasonStats: seasonStats,
		TeamStandings:      collegeStandings,
	}
}

func getTransferPortalProfileMapByTeamID(id string) map[uint]structs.TransferPortalProfile {
	profiles := GetOnlyTransferPortalProfilesByTeamID(id)

	profileMap := make(map[uint]structs.TransferPortalProfile)

	for _, profile := range profiles {
		profileMap[profile.CollegePlayerID] = profile
	}

	return profileMap
}
