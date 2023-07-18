package managers

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"time"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
)

func SyncRecruiting(timestamp structs.Timestamp) {
	db := dbprovider.GetInstance().GetDB()
	fmt.Println(time.Now().UnixNano())
	//GetCurrentWeek

	if timestamp.RecruitingSynced {
		log.Fatalln("Recruiting already ran for this week. Please wait until next week to sync recruiting again.")
	}

	recruitProfilePointsMap := make(map[string]int)

	teamRecruitingProfiles := GetTeamRecruitingProfilesForRecruitSync()

	teamMap := make(map[string]*structs.TeamRecruitingProfile)

	for i := 0; i < len(teamRecruitingProfiles); i++ {
		teamMap[strconv.Itoa(int(teamRecruitingProfiles[i].ID))] = &teamRecruitingProfiles[i]
	}

	var modifier1 float64 = 75
	var modifierFor5Star float64 = 125
	weeksOfRecruiting := 15

	var recruitProfiles []structs.PlayerRecruitProfile

	var signeesLog []string

	// Get every recruit
	recruits := GetAllUnsignedRecruits()

	// Iterate through every recruit
	for _, recruit := range recruits {
		recruitProfiles = GetRecruitPlayerProfilesByRecruitId(strconv.Itoa(int(recruit.ID)))

		if len(recruitProfiles) == 0 {
			continue
		}

		var recruitProfilesWithScholarship []structs.PlayerRecruitProfile

		eligibleTeams := 0

		pointsPlaced := false

		var totalPointsOnRecruit float64 = 0

		var eligiblePointThreshold float64 = 0

		var signThreshold float64

		for i := 0; i < len(recruitProfiles); i++ {

			if recruitProfiles[i].CurrentWeeksPoints == 0 {
				continue
			} else {
				pointsPlaced = true
			}

			rpa := structs.RecruitPointAllocation{
				RecruitID:        recruitProfiles[i].RecruitID,
				TeamProfileID:    recruitProfiles[i].ProfileID,
				RecruitProfileID: recruitProfiles[i].ID,
				WeekID:           timestamp.CollegeWeekID,
			}

			var curr float64 = float64(recruitProfiles[i].CurrentWeeksPoints)

			// Region / State bonus
			if recruitProfiles[i].HasRegionBonus && recruit.Stars != 5 {
				curr = curr * 1.05
			} else if recruitProfiles[i].HasStateBonus && recruit.Stars != 5 {
				curr = curr * 1.1
			}
			// Bonus Points value when saving

			if recruitProfiles[i].CurrentWeeksPoints < 0 || recruitProfiles[i].CurrentWeeksPoints > 20 {
				curr = 0
				rpa.ApplyCaughtCheating()
			}

			rpa.UpdatePointsSpent(float64(recruitProfiles[i].CurrentWeeksPoints), curr)
			recruitProfiles[i].AllocateTotalPoints(curr)
			recruitProfilePointsMap[recruitProfiles[i].TeamAbbreviation] += recruitProfiles[i].CurrentWeeksPoints

			// Add RPA to point allocations list
			err := db.Create(&rpa).Error
			if err != nil {
				fmt.Println(err.Error())
				log.Fatalf("Could not save Point Allocation")
			}
		}

		if !pointsPlaced {
			fmt.Println("Skipping over " + recruit.FirstName + " " + recruit.LastName)
			continue
		}

		sort.Sort(structs.ByPoints(recruitProfiles))

		for i := 0; i < len(recruitProfiles) && pointsPlaced; i++ {
			recruitTeamProfile := teamMap[(strconv.Itoa(int(recruitProfiles[i].ProfileID)))]

			if recruitTeamProfile.TotalCommitments >= recruitTeamProfile.RecruitClassSize {
				continue
			}
			if eligiblePointThreshold == 0 && recruitProfiles[i].Scholarship {
				eligiblePointThreshold = float64(recruitProfiles[i].TotalPoints) * 0.66
			}

			if recruitProfiles[i].Scholarship && recruitProfiles[i].TotalPoints >= eligiblePointThreshold {
				totalPointsOnRecruit += recruitProfiles[i].TotalPoints
				eligibleTeams += 1
				recruitProfilesWithScholarship = append(recruitProfilesWithScholarship, recruitProfiles[i])
			}
		}

		var mod1 float64

		if recruit.Stars == 5 {
			mod1 = float64(modifierFor5Star)
		} else {
			mod1 = float64(modifier1)
		}

		// Change?
		// Assign point totals
		// If there are any modifiers
		// Evaluate
		firstMod := mod1 - float64(timestamp.CollegeWeek)
		secondMod := float64(eligibleTeams) / (float64(recruit.RecruitModifier) / 100)
		thirdMod := math.Log10(float64(weeksOfRecruiting - timestamp.CollegeWeek))
		signThreshold = firstMod * secondMod * thirdMod
		recruit.ApplySigningStatus(totalPointsOnRecruit, signThreshold)
		// Change logic to withold teams without available scholarships
		if totalPointsOnRecruit > signThreshold && eligibleTeams > 0 && pointsPlaced {
			var winningTeamID uint = 0
			var odds float64 = 0

			for winningTeamID == 0 {
				percentageOdds := 1 + rand.Float64()*(totalPointsOnRecruit-1)
				var currentProbability float64 = 0

				for i := 0; i < len(recruitProfilesWithScholarship); i++ {
					// If a team has no available scholarships or if a team has 25 commitments, continue
					currentProbability += recruitProfilesWithScholarship[i].TotalPoints
					if float64(percentageOdds) <= currentProbability {
						// WINNING TEAM
						winningTeamID = recruitProfilesWithScholarship[i].ProfileID
						odds = float64(recruitProfilesWithScholarship[i].TotalPoints) / float64(totalPointsOnRecruit) * 100
						break
					}
				}

				if winningTeamID > 0 {
					recruitTeamProfile := teamMap[(strconv.Itoa(int(winningTeamID)))]
					if recruitTeamProfile.TotalCommitments < recruitTeamProfile.RecruitClassSize {
						recruitTeamProfile.IncreaseCommitCount()
						teamAbbreviation := recruitTeamProfile.TeamAbbr
						recruit.AssignCollege(teamAbbreviation)

						newsLog := structs.NewsLog{
							WeekID:      timestamp.CollegeWeekID + 1,
							Week:        uint(timestamp.CollegeWeek),
							SeasonID:    timestamp.SeasonID,
							Season:      uint(timestamp.Season),
							MessageType: "Commitment",
							Message:     recruit.FirstName + " " + recruit.LastName + ", " + strconv.Itoa(recruit.Stars) + " star " + recruit.Position + " from " + recruit.State + ", " + recruit.Country + " has signed with " + recruit.TeamAbbr + " with " + strconv.Itoa(int(odds)) + " percent odds.",
						}

						db.Create(&newsLog)
						fmt.Println("Created new log!")

						db.Save(&recruitTeamProfile)
						fmt.Println("Saved " + recruitTeamProfile.TeamAbbr + " profile.")

						for i := 0; i < len(recruitProfiles); i++ {
							if recruitProfiles[i].ProfileID == winningTeamID {
								recruitProfiles[i].SignPlayer()
							} else {
								recruitProfiles[i].LockPlayer()
								if recruitProfiles[i].Scholarship {
									tp := teamMap[strconv.Itoa(int(recruitProfiles[i].ProfileID))]

									tp.ReallocateScholarship()
									err := db.Save(&tp).Error
									if err != nil {
										fmt.Println(err.Error())
										log.Fatalf("Could not sync recruiting profile.")
									}

									fmt.Println("Reallocated Scholarship to " + tp.TeamAbbr)
								}
							}
						}
					} else {
						recruitProfilesWithScholarship = util.FilterOutRecruitingProfile(recruitProfilesWithScholarship, int(winningTeamID))
						// If there are no longer any teams contending due to reaching the max class size, break the loop
						winningTeamID = 0
						if len(recruitProfilesWithScholarship) == 0 {
							break
						}

						totalPointsOnRecruit = 0
						for _, rp := range recruitProfilesWithScholarship {
							totalPointsOnRecruit += rp.TotalPoints
						}
					}
				}
			}
			recruit.UpdateTeamID(winningTeamID)
		}

		// Save Player Files towards Recruit
		for _, rp := range recruitProfiles {
			// Save Team Profile
			err := db.Save(&rp).Error
			if err != nil {
				fmt.Println(err.Error())
				log.Fatalf("Could not sync recruiting profile.")
			}

			fmt.Println("Save recruit profile from " + rp.TeamAbbreviation + " towards " + recruit.FirstName + " " + recruit.LastName)
		}

		// Save Recruit
		err := db.Save(&recruit).Error
		if err != nil {
			fmt.Println(err.Error())
			log.Fatalf("Could not sync recruit")
		}
		fmt.Println("Save Recruit " + recruit.FirstName + " " + recruit.LastName)
	}

	// Update rank system for all teams
	var maxESPNScore float64 = 0
	var minESPNScore float64 = 100000
	var maxRivalsScore float64 = 0
	var minRivalsScore float64 = 100000
	var max247Score float64 = 0
	var min247Score float64 = 100000

	for i := 0; i < len(teamRecruitingProfiles); i++ {

		signedRecruits := GetSignedRecruitsByTeamProfileID(strconv.Itoa(int(teamRecruitingProfiles[i].TeamID)))

		teamRecruitingProfiles[i].UpdateTotalSignedRecruits(len(signedRecruits))

		team247Rank := util.Get247TeamRanking(teamRecruitingProfiles[i], signedRecruits)
		teamESPNRank := util.GetESPNTeamRanking(teamRecruitingProfiles[i], signedRecruits)
		teamRivalsRank := util.GetRivalsTeamRanking(teamRecruitingProfiles[i], signedRecruits)

		teamRecruitingProfiles[i].Assign247Rank(team247Rank)
		teamRecruitingProfiles[i].AssignESPNRank(teamESPNRank)
		teamRecruitingProfiles[i].AssignRivalsRank(teamRivalsRank)
		if teamESPNRank > maxESPNScore {
			maxESPNScore = teamESPNRank
		}
		if teamESPNRank < minESPNScore {
			minESPNScore = teamESPNRank
		}
		if teamRivalsRank > maxRivalsScore {
			maxRivalsScore = teamRivalsRank
		}
		if teamRivalsRank < minRivalsScore {
			minRivalsScore = teamRivalsRank
		}
		if team247Rank > max247Score {
			max247Score = team247Rank
		}
		if team247Rank < min247Score {
			min247Score = team247Rank
		}
		fmt.Println("Setting Recruiting Ranks for " + teamRecruitingProfiles[i].TeamAbbr)

	}

	espnDivisor := (maxESPNScore - minESPNScore)
	divisor247 := (max247Score - min247Score)
	rivalsDivisor := (maxRivalsScore - minRivalsScore)

	for _, rp := range teamRecruitingProfiles {
		if recruitProfilePointsMap[rp.TeamAbbr] > rp.WeeklyPoints {
			rp.ApplyCaughtCheating()
		}

		var avg float64 = 0
		if espnDivisor > 0 && divisor247 > 0 && rivalsDivisor > 0 {
			distributionESPN := (rp.ESPNScore - minESPNScore) / espnDivisor
			distribution247 := (rp.Rank247Score - min247Score) / divisor247
			distributionRivals := (rp.RivalsScore - minRivalsScore) / rivalsDivisor

			avg = (distributionESPN + distribution247 + distributionRivals)

			rp.AssignCompositeRank(avg)
		}
		rp.ResetSpentPoints()

		// Save TEAM Recruiting Profile
		err := db.Save(&rp).Error
		if err != nil {
			fmt.Println(err.Error())
			log.Fatalf("Could not save Team Recruiting Profile")
		}
		fmt.Println("Saved Rank Scores for Team " + rp.TeamAbbr)
	}

	for _, log := range signeesLog {
		fmt.Println(log)
	}

	timestamp.ToggleRecruiting()

	err := db.Save(&timestamp).Error
	if err != nil {
		fmt.Println(err.Error())
		log.Fatalf("Could not save timestamp.")
	}
}

func FillAIRecruitingBoards() {
	db := dbprovider.GetInstance().GetDB()
	fmt.Println(time.Now().UnixNano())
	ts := GetTimestamp()

	AITeams := GetOnlyAITeamRecruitingProfiles()

	UnsignedRecruits := GetAllUnsignedRecruits()

	regionMap := util.GetRegionMap()

	boardCount := 30

	if ts.CollegeWeek > 3 {
		boardCount = 15
	}

	for _, team := range AITeams {
		count := 0
		if !team.IsAI || team.TotalCommitments >= team.RecruitClassSize {
			continue
		}
		id := strconv.Itoa(int(team.ID))
		existingBoard := GetAllRecruitsByProfileID(id)

		count = len(existingBoard)

		if count >= boardCount {
			continue
		}

		currentRoster := GetCollegePlayersByTeamId(id)
		teamNeedsMap := make(map[string]bool)
		positionCount := make(map[string]int)

		for _, r := range currentRoster {
			positionCount[r.Position] += 1
		}

		if positionCount["PG"] < 3 {
			teamNeedsMap["PG"] = true
		} else if positionCount["SG"] < 4 {
			teamNeedsMap["SG"] = true
		} else if positionCount["SF"] < 4 {
			teamNeedsMap["SF"] = true
		} else if positionCount["PF"] < 4 {
			teamNeedsMap["PF"] = true
		} else if positionCount["C"] < 3 {
			teamNeedsMap["C"] = true
		}

		for _, croot := range UnsignedRecruits {
			if count == boardCount {
				break
			}
			if croot.IsCustomCroot || (!teamNeedsMap[croot.Position] && ts.CollegeWeek < 10) ||
				(croot.Stars == 5 && team.AIQuality != "Blue Blood") {
				continue
			}

			crootProfile := GetPlayerRecruitProfileByPlayerId(strconv.Itoa(int(croot.ID)), strconv.Itoa(int(team.ID)))
			if crootProfile.ID > 0 || crootProfile.RemovedFromBoard || crootProfile.IsLocked {
				continue
			}

			crootProfiles := GetRecruitPlayerProfilesByRecruitId(strconv.Itoa(int(croot.ID)))

			leadingVal := util.IsAITeamContendingForCroot(crootProfiles)
			if leadingVal > 11 {
				continue
			}

			odds := 5

			if ts.CollegeWeek > 5 {
				odds = 15
			}

			if croot.Country == "USA" {
				if regionMap[croot.State] == team.Region {
					odds += 25
				}
				if croot.State == team.State {
					odds += 33
				}
			}
			/* Initial Base */
			if team.AIQuality == "Blue Blood" && croot.Stars == 5 {
				odds += 5
			} else if team.AIQuality == "Cinderella" && croot.Stars == 4 {
				odds += 15
			} else if team.AIQuality == "P6" && croot.Stars == 4 {
				odds += 8
			} else if team.AIQuality == "P6" && croot.Stars == 3 {
				odds += 20
			} else if team.AIQuality == "Mid-Major" && croot.Stars < 4 {
				odds += 7
			} else if team.AIQuality == "Mid-Major" && croot.Stars < 3 {
				odds += 25
			}

			if team.AIQuality == "Cinderella" && util.IsPlayerHighPotential(croot) {
				odds += 20
			}

			if team.AIValue == "Star" {
				odds += getOddsIncrementByStar(5, croot.Stars)
			} else if team.AIValue == "Potential" {
				odds += getOddsIncrementByPotential(5, croot.Potential)
			} else if team.AIValue == "Talent" {
				odds += getOddsIncrementByTalent(croot.Shooting2, croot.SpecShooting2, team.AIAttribute1 == "Shooting2" || team.AIAttribute2 == "Shooting2", team.AIQuality == "Mid-Major")
				odds += getOddsIncrementByTalent(croot.Shooting3, croot.SpecShooting3, team.AIAttribute1 == "Shooting3" || team.AIAttribute2 == "Shooting3", team.AIQuality == "Mid-Major")
				odds += getOddsIncrementByTalent(croot.Finishing, croot.SpecFinishing, team.AIAttribute1 == "Finishing" || team.AIAttribute2 == "Finishing", team.AIQuality == "Mid-Major")
				odds += getOddsIncrementByTalent(croot.FreeThrow, croot.SpecFreeThrow, team.AIAttribute1 == "FreeThrow" || team.AIAttribute2 == "FreeThrow", team.AIQuality == "Mid-Major")
				odds += getOddsIncrementByTalent(croot.Ballwork, croot.SpecBallwork, team.AIAttribute1 == "Ballwork" || team.AIAttribute2 == "Ballwork", team.AIQuality == "Mid-Major")
				odds += getOddsIncrementByTalent(croot.Rebounding, croot.SpecRebounding, team.AIAttribute1 == "Rebounding" || team.AIAttribute2 == "Rebounding", team.AIQuality == "Mid-Major")
				odds += getOddsIncrementByTalent(croot.InteriorDefense, croot.SpecInteriorDefense, team.AIAttribute1 == "InteriorDefense" || team.AIAttribute2 == "InteriorDefense", team.AIQuality == "Mid-Major")
				odds += getOddsIncrementByTalent(croot.PerimeterDefense, croot.SpecPerimeterDefense, team.AIAttribute1 == "PerimeterDefense" || team.AIAttribute2 == "PerimeterDefense", team.AIQuality == "Mid-Major")
			}

			chance := util.GenerateIntFromRange(1, 100)

			var teamsWithBoards []structs.PlayerRecruitProfile

			for _, team := range crootProfiles {
				if !team.RemovedFromBoard {
					teamsWithBoards = append(teamsWithBoards, team)
				}
			}

			if chance <= odds && len(teamsWithBoards) < 25 {
				playerProfile := structs.PlayerRecruitProfile{
					RecruitID:          croot.ID,
					ProfileID:          team.ID,
					SeasonID:           ts.SeasonID,
					TotalPoints:        0,
					CurrentWeeksPoints: 0,
					SpendingCount:      0,
					Scholarship:        false,
					ScholarshipRevoked: false,
					TeamAbbreviation:   team.TeamAbbr,
					RecruitModifier:    croot.RecruitModifier,
					IsSigned:           false,
					IsLocked:           false,
				}

				err := db.Save(&playerProfile).Error
				if err != nil {
					log.Fatalln("Could not add " + croot.FirstName + " " + croot.LastName + " to " + team.TeamAbbr + " Recruiting Board.")
				}

				count++
			}
		}
	}
}

func AllocatePointsToAIBoards() {
	db := dbprovider.GetInstance().GetDB()
	fmt.Println(time.Now().UnixNano())
	ts := GetTimestamp()

	AITeams := GetOnlyAITeamRecruitingProfiles()

	for _, team := range AITeams {
		if team.SpentPoints >= team.WeeklyPoints || team.TotalCommitments >= team.RecruitClassSize {
			continue
		}

		teamRecruits := GetAllRecruitsByProfileID(strconv.Itoa(int(team.ID)))

		for _, croot := range teamRecruits {
			// If a team has no more points to spend, break the loop
			pointsRemaining := team.WeeklyPoints - team.SpentPoints
			if team.SpentPoints >= team.WeeklyPoints || pointsRemaining <= 0 || (pointsRemaining < 1 && pointsRemaining > 0) {
				break
			}

			// If a croot was signed OR has points already placed on the croot, move on to the next croot
			if croot.IsSigned || croot.CurrentWeeksPoints > 0 || croot.ScholarshipRevoked {
				continue
			}
			removeCrootFromBoard := false
			num := 0
			// If a croot is locked and signed with a different team, remove from the team board and continue
			if croot.IsLocked && croot.TeamAbbreviation != croot.Recruit.TeamAbbr {
				removeCrootFromBoard = true
			}

			if !removeCrootFromBoard {
				profiles := GetRecruitPlayerProfilesByRecruitId(strconv.Itoa(int(croot.RecruitID)))

				// If an AI team previously spent points on a croot, use the previous week allocation.
				if croot.PreviouslySpentPoints > 0 {
					leadingTeamVal := util.IsAITeamContendingForCroot(profiles)
					// If the allocation to be placed keeps the team in the lead, or if the lead is by 11 points or less
					if float64(croot.PreviouslySpentPoints)+croot.TotalPoints >= float64(leadingTeamVal)*0.66 || leadingTeamVal < 11 {
						num = croot.PreviouslySpentPoints
						if num > pointsRemaining {
							num = pointsRemaining
						}
					} else {
						removeCrootFromBoard = true
					}
				} else {
					// Flip a coin. If heads, the team spends points on the croot.
					// Else, move on to the next croot
					maxChance := 2
					if ts.CollegeWeek > 3 {
						maxChance = 4
					}
					chance := util.GenerateIntFromRange(1, maxChance)
					if (chance < 2 && ts.CollegeWeek <= 3) || (chance < 4 && ts.CollegeWeek > 3) {
						continue
					}

					min := 1
					max := 12

					if team.AIBehavior == "Conservative" {
						max = 10
					} else if team.AIBehavior == "Aggressive" {
						min = 10
						max = 15
					} else {
						min = 6
						max = 12
					}

					num = util.GenerateIntFromRange(min, max)
					if num > pointsRemaining {
						num = pointsRemaining
					}
					// Check to see if other teams are contending
					leadingValPoints := util.IsAITeamContendingForCroot(profiles)
					if float64(num)+croot.TotalPoints < float64(leadingValPoints)*0.66 {
						removeCrootFromBoard = true
					}
					if leadingValPoints < 11 {
						removeCrootFromBoard = false
					}
				}
			}

			// If the Croot needs to be removed from the board, remove it and move on.
			if removeCrootFromBoard || (team.ScholarshipsAvailable == 0 && !croot.Scholarship) {
				if croot.Scholarship {
					croot.RevokeScholarship()
					team.ReallocateScholarship()
				}
				croot.RemoveRecruitFromBoard()
				fmt.Println("Because " + croot.Recruit.FirstName + " " + croot.Recruit.LastName + " is heavily considering other teams, they are being removed from " + team.TeamAbbr + "'s Recruiting Board.")
				db.Save(&croot)
				continue
			}

			// If final week, do a spread of points
			if ts.CollegeWeek == 14 {
				num = 5
			}

			// Allocate points and save
			croot.AllocatePoints(num)
			if !croot.Scholarship && team.ScholarshipsAvailable > 0 {
				croot.ToggleScholarship(true, false)
				team.SubtractScholarshipsAvailable()
			}
			team.AIAllocateSpentPoints(num)
			// Save croot
			db.Save(&croot)
			fmt.Println(team.TeamAbbr + " allocating " + strconv.Itoa(num) + " points to " + croot.Recruit.FirstName + " " + croot.Recruit.LastName)
		}
		// Save Team Profile after iterating through recruits
		fmt.Println("Saved " + team.TeamAbbr + " Recruiting Board!")
		db.Save(&team)
	}
}

func ResetAIBoardsForCompletedTeams() {
	db := dbprovider.GetInstance().GetDB()

	AITeams := GetTeamRecruitingProfilesForRecruitSync()

	for _, team := range AITeams {
		// If a team already has the maximum allowed for their recruiting class, take all Recruit Profiles for that team where the recruit hasn't signed, and reset their total points.
		// This is so that these unsigned recruits can be recruited for and will allow the AI to put points onto those recruits.

		if team.TotalCommitments >= team.RecruitClassSize {
			teamRecruits := GetAllRecruitsByProfileID(strconv.Itoa(int(team.ID)))

			for _, croot := range teamRecruits {
				if croot.IsSigned || croot.IsLocked || croot.TotalPoints == 0 {
					continue
				}
				croot.ResetTotalPoints()
				if team.IsAI {
					croot.ToggleTotalMax()
				}
				db.Save(&croot)
			}
			team.ResetSpentPoints()
			db.Save(&team)
		}
	}
}

func getOddsIncrementByStar(init int, stars int) int {
	return init * stars
}

func getOddsIncrementByPotential(init int, potential int) int {
	potentialFloor := potential / 10
	return init * potentialFloor
}

func getOddsIncrementByTalent(attr int, attrspec, attrMatch bool, isMidMajor bool) int {
	attrRequirement := 15
	if isMidMajor {
		attrRequirement = 10
	}
	if attrMatch && (attrspec || attr > attrRequirement) {
		return 25
	}
	return 0
}
