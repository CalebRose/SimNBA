package managers

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/repository"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
	"gorm.io/gorm"
)

func SyncRecruiting() {
	db := dbprovider.GetInstance().GetDB()
	fmt.Println(time.Now().UnixNano())
	//GetCurrentWeek
	timestamp := GetTimestamp()

	if timestamp.RecruitingSynced {
		log.Fatalln("Recruiting already ran for this week. Please wait until next week to sync recruiting again.")
	}

	if !timestamp.IsRecruitingLocked {
		timestamp.ToggleLockRecruiting()
		db.Save(&timestamp)
	}

	var modifier1 float64 = 75
	var modifierFor5Star float64 = 125
	weeksOfRecruiting := 15
	eligibleThresholdPercentage := 0.66
	pointLimit := 20.0
	recruitProfilePointsMap := make(map[string]float64)
	teamRecruitingProfiles := GetTeamRecruitingProfilesForRecruitSync()
	teamMap := make(map[string]*structs.TeamRecruitingProfile)
	for i := 0; i < len(teamRecruitingProfiles); i++ {
		teamMap[strconv.Itoa(int(teamRecruitingProfiles[i].ID))] = &teamRecruitingProfiles[i]
		recruitProfilePointsMap[teamRecruitingProfiles[i].TeamAbbr] = 0.0
	}
	allRecruitProfiles := GetAllRecruitPlayerProfiles()
	recruitProfileMap := MakeRecruitProfileMapByRecruitID(allRecruitProfiles)

	// Get every recruit
	recruits := GetAllUnsignedRecruits()

	// Iterate through every recruit
	skipOver := []int{1, 2}

	sort.Ints(skipOver)

	for _, recruit := range recruits {
		id := recruit.ID
		i := sort.SearchInts(skipOver, int(id))
		if i < len(skipOver) && skipOver[i] == int(id) {
			continue
		}
		recruitProfiles := recruitProfileMap[id]

		if len(recruitProfiles) == 0 {
			continue
		}

		var recruitProfilesWithScholarship []structs.PlayerRecruitProfile
		eligibleTeams := 0
		pointsPlaced := false
		var totalPointsOnRecruit float64 = 0
		var eligiblePointThreshold float64 = 0
		var signThreshold float64

		allocatePointsToRecruit(recruit, &recruitProfiles, teamMap, pointLimit, &pointsPlaced, timestamp, &recruitProfilePointsMap, db)

		if !pointsPlaced {
			continue
		}

		sort.Sort(structs.ByPoints(recruitProfiles))

		for i := 0; i < len(recruitProfiles) && pointsPlaced; i++ {
			recruitTeamProfile := teamMap[(strconv.Itoa(int(recruitProfiles[i].ProfileID)))]

			if recruitTeamProfile.TotalCommitments >= recruitTeamProfile.RecruitClassSize {
				continue
			}
			if eligiblePointThreshold == 0 && recruitProfiles[i].Scholarship {
				eligiblePointThreshold = float64(recruitProfiles[i].TotalPoints) * eligibleThresholdPercentage
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
		// Change logic to withold teams without available scholarships
		passedTheSigningThreshold := totalPointsOnRecruit > signThreshold && eligibleTeams > 0 && pointsPlaced
		recruit.ApplySigningStatus(totalPointsOnRecruit, signThreshold, passedTheSigningThreshold)
		if passedTheSigningThreshold {
			var winningTeamID uint = 0
			var odds float64 = 0

			for winningTeamID == 0 {
				percentageOdds := rand.Float64() * (totalPointsOnRecruit)
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
						recruitTeamProfile.AddStarPlayer(recruit.Stars)
						teamAbbreviation := recruitTeamProfile.TeamAbbr
						recruit.AssignCollege(teamAbbreviation)
						message := recruit.FirstName + " " + recruit.LastName + ", " + strconv.Itoa(recruit.Stars) + " star " + recruit.Position + " from " + recruit.State + ", " + recruit.Country + " has signed with " + recruit.TeamAbbr + " with " + strconv.Itoa(int(odds)) + " percent odds."
						CreateNewsLog("CBB", message, "Commitment", int(winningTeamID), timestamp)
						fmt.Println("Created new log!")

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
			repository.SaveCBBRecruitProfile(rp, db)
			fmt.Println("Save recruit profile from " + rp.TeamAbbreviation + " towards " + recruit.FirstName + " " + recruit.LastName)
		}

		// Save Recruit
		repository.SaveCBBRecruit(recruit, db)
	}

	updateTeamRankings(teamRecruitingProfiles, teamMap, recruitProfilePointsMap, db)

	if timestamp.IsRecruitingLocked {
		timestamp.ToggleLockRecruiting()
	}

	repository.SaveTimeStamp(timestamp, db)
}

func FillAIRecruitingBoards() {
	db := dbprovider.GetInstance().GetDB()
	fmt.Println(time.Now().UnixNano())
	ts := GetTimestamp()

	AITeams := GetOnlyAITeamRecruitingProfiles()
	// Shuffles the list of AI teams so that it's not always iterating from A-Z. Gives the teams at the lower end of the list a chance to recruit other croots
	rand.Shuffle(len(AITeams), func(i, j int) {
		AITeams[i], AITeams[j] = AITeams[j], AITeams[i]
	})

	UnsignedRecruits := GetAllUnsignedRecruits()
	profiles := GetAllRecruitPlayerProfiles()
	recruitProfileMapByRecruit := MakeRecruitProfileMapByRecruitID(profiles)
	recruitProfileMapByProfile := MakeRecruitProfileMapByProfileID(profiles)

	regionMap := util.GetRegionMap()

	boardCount := 20

	profileBatch := []structs.PlayerRecruitProfile{}

	for _, team := range AITeams {
		count := 0
		lateSeasonCount := 0
		if !team.IsAI || team.TotalCommitments >= team.RecruitClassSize || team.ScholarshipsAvailable == 0 || team.SpentPoints == 50 {
			continue
		}
		id := strconv.Itoa(int(team.ID))
		existingBoard := recruitProfileMapByProfile[team.ID]

		count = len(existingBoard)
		boardCount = team.RecruitClassSize * 5
		if count >= boardCount {
			continue
		}
		currentRoster := GetCollegePlayersByTeamId(id)
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

		for _, r := range currentRoster {
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

		for _, croot := range UnsignedRecruits {
			if count == boardCount || (ts.CollegeWeek > 12 && lateSeasonCount > 12) {
				break
			}
			recruitingNeed := teamNeedsMap[croot.Position]
			willPursue := true
			if croot.IsCustomCroot || (!recruitingNeed && ts.CollegeWeek < 12) ||
				(croot.Stars > team.AIStarMax || croot.Stars < team.AIStarMin) {
				willPursue = false
			}

			crootProfiles := recruitProfileMapByRecruit[croot.ID]
			crootProfile := FindRecruitPlayerProfileByProfileID(crootProfiles, team.ID)

			leadingVal := util.IsAITeamContendingForCroot(crootProfiles)
			if leadingVal > 14 {
				willPursue = false
			}

			if ts.CollegeWeek > 12 {
				willPursue = true
			}

			if !willPursue || crootProfile.RemovedFromBoard || crootProfile.IsLocked || crootProfile.ID > 0 {
				continue
			}

			odds := 5

			if ts.CollegeWeek > 5 {
				odds = 15
			}

			if ts.CollegeWeek > 12 {
				odds = 35
			}

			if croot.Country == "USA" {
				if regionMap[croot.State] == team.Region {
					odds += 30
				}
				if croot.State == team.State {
					odds += 38
				}
				if regionMap[croot.State] != team.Region && croot.State != team.State && team.AIQuality == "Mid-Major" {
					odds -= 15
				}
			}

			if team.AIQuality == "Cinderella" && util.IsPlayerHighPotential(croot) {
				odds += 5
			}

			if team.AIValue == "Star" {
				odds += getOddsIncrementByStar(5, croot.Stars)
			} else if team.AIValue == "Potential" {
				odds += getOddsIncrementByPotential(5, croot.Potential, team.AIQuality == "Mid-Major")
			} else if team.AIValue == "Talent" {
				odds += getOddsIncrementByTalent(croot.Shooting2, croot.Stars, croot.SpecShooting2, team.AIAttribute1 == "Shooting2" || team.AIAttribute2 == "Shooting2", team.AIQuality == "Mid-Major")
				odds += getOddsIncrementByTalent(croot.Shooting3, croot.Stars, croot.SpecShooting3, team.AIAttribute1 == "Shooting3" || team.AIAttribute2 == "Shooting3", team.AIQuality == "Mid-Major")
				odds += getOddsIncrementByTalent(croot.Finishing, croot.Stars, croot.SpecFinishing, team.AIAttribute1 == "Finishing" || team.AIAttribute2 == "Finishing", team.AIQuality == "Mid-Major")
				odds += getOddsIncrementByTalent(croot.FreeThrow, croot.Stars, croot.SpecFreeThrow, team.AIAttribute1 == "FreeThrow" || team.AIAttribute2 == "FreeThrow", team.AIQuality == "Mid-Major")
				odds += getOddsIncrementByTalent(croot.Ballwork, croot.Stars, croot.SpecBallwork, team.AIAttribute1 == "Ballwork" || team.AIAttribute2 == "Ballwork", team.AIQuality == "Mid-Major")
				odds += getOddsIncrementByTalent(croot.Rebounding, croot.Stars, croot.SpecRebounding, team.AIAttribute1 == "Rebounding" || team.AIAttribute2 == "Rebounding", team.AIQuality == "Mid-Major")
				odds += getOddsIncrementByTalent(croot.InteriorDefense, croot.Stars, croot.SpecInteriorDefense, team.AIAttribute1 == "InteriorDefense" || team.AIAttribute2 == "InteriorDefense", team.AIQuality == "Mid-Major")
				odds += getOddsIncrementByTalent(croot.PerimeterDefense, croot.Stars, croot.SpecPerimeterDefense, team.AIAttribute1 == "PerimeterDefense" || team.AIAttribute2 == "PerimeterDefense", team.AIQuality == "Mid-Major")
			}

			chance := util.GenerateIntFromRange(1, 100)

			teamCount := 0

			for _, p := range crootProfiles {
				if p.RemovedFromBoard {
					continue
				}
				teamCount++
			}

			teamLimit := 10
			if ts.CollegeWeek > 12 {
				teamLimit = 6
			}

			if chance <= odds && teamCount < teamLimit {
				if ts.CollegeWeek > 12 {
					lateSeasonCount += 1
				}
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
				crootProfiles = append(crootProfiles, playerProfile)
				recruitProfileMapByRecruit[croot.ID] = crootProfiles
				profileBatch = append(profileBatch, playerProfile)
				count++
			}
		}
	}

	repository.CreatePlayerRecruitProfileRecordsBatch(db, profileBatch, 100)
}

func AllocatePointsToAIBoards() {
	db := dbprovider.GetInstance().GetDB()
	fmt.Println(time.Now().UnixNano())
	ts := GetTimestamp()

	AITeams := GetOnlyAITeamRecruitingProfiles()

	// Shuffles the list of AI teams so that it's not always iterating from A-Z. Gives the teams at the lower end of the list a chance to recruit other croots
	rand.Shuffle(len(AITeams), func(i, j int) {
		AITeams[i], AITeams[j] = AITeams[j], AITeams[i]
	})

	for _, team := range AITeams {
		if team.SpentPoints >= team.WeeklyPoints || team.TotalCommitments >= team.RecruitClassSize {
			continue
		}
		id := strconv.Itoa(int(team.ID))
		currentRoster := GetCollegePlayersByTeamId(id)
		signedCroots := GetSignedRecruitsByTeamProfileID(id)
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

		for _, r := range currentRoster {
			if r.WillDeclare {
				continue
			}
			positionCount[r.Position] += 1
		}

		for _, r := range signedCroots {
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

		teamRecruits := GetAllRecruitsByProfileID(strconv.Itoa(int(team.ID)))

		sort.Slice(teamRecruits, func(i, j int) bool {
			iCroot := teamRecruits[i].Recruit
			jCroot := teamRecruits[j].Recruit
			return iCroot.Stars > jCroot.Stars
		})

		for _, croot := range teamRecruits {
			// If a team has no more points to spend, break the loop
			pointsRemaining := team.WeeklyPoints - team.SpentPoints
			if team.SpentPoints >= team.WeeklyPoints || pointsRemaining <= 0 || (pointsRemaining < 1 && pointsRemaining > 0) {
				break
			}

			recruitingNeed := teamNeedsMap[croot.Recruit.Position]
			// If a croot was signed OR has points already placed on the croot, move on to the next croot
			if croot.IsSigned || croot.CurrentWeeksPoints > 0 || croot.ScholarshipRevoked || !recruitingNeed {
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
					if float64(croot.PreviouslySpentPoints)+croot.TotalPoints >= float64(leadingTeamVal)*0.66 || leadingTeamVal < 14 {
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

					min := team.AIMaxThreshold
					max := team.AIMinThreshold

					if team.AIBehavior == "Conservative" {
						max -= 2
					} else if team.AIBehavior == "Aggressive" {
						min += 2
						max += 2
					}
					if max > 20 {
						max = 20
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
					if leadingValPoints < 14 {
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
				if croot.IsLocked {
					fmt.Println("Because " + croot.Recruit.FirstName + " " + croot.Recruit.LastName + " signed with a different team, they are being removed from " + team.TeamAbbr + "'s Recruiting Board.")
				} else {
					fmt.Println("Because " + croot.Recruit.FirstName + " " + croot.Recruit.LastName + " is heavily considering other teams, they are being removed from " + team.TeamAbbr + "'s Recruiting Board.")
				}
				repository.SaveCBBRecruitProfile(croot, db)
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
			repository.SaveCBBRecruitProfile(croot, db)
			fmt.Println(team.TeamAbbr + " allocating " + strconv.Itoa(num) + " points to " + croot.Recruit.FirstName + " " + croot.Recruit.LastName)

			positionCount[croot.Recruit.Position] += 1
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
		}
		// Save Team Profile after iterating through recruits
		fmt.Println("Saved " + team.TeamAbbr + " Recruiting Board!")
		repository.SaveCBBTeamRecruitingProfile(team, db)
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
				repository.SaveCBBRecruitProfile(croot, db)
			}
			team.ResetSpentPoints()
			repository.SaveCBBTeamRecruitingProfile(team, db)
		}
	}
}

func getOddsIncrementByStar(init int, stars int) int {
	return init * stars
}

func getOddsIncrementByPotential(init int, potential int, isMidMajor bool) int {
	divisor := 10
	if isMidMajor {
		divisor = 20
	}
	potentialFloor := potential / divisor
	return init * potentialFloor
}

func getOddsIncrementByTalent(attr, stars int, attrspec, attrMatch bool, isMidMajor bool) int {
	attrRequirement := 14
	if attrMatch && (attrspec || attr > attrRequirement) {
		if stars > 3 && isMidMajor {
			return 10
		}
		return 25
	}
	return 0
}

func allocatePointsToRecruit(recruit structs.Recruit, recruitProfiles *[]structs.PlayerRecruitProfile, teamMap map[string]*structs.TeamRecruitingProfile, pointLimit float64, pointsPlaced *bool, timestamp structs.Timestamp, recruitProfilePointsMap *map[string]float64, db *gorm.DB) {
	// numWorkers := 3
	var mapMutex sync.Mutex
	numWorkers := runtime.NumCPU()
	if numWorkers > 3 {
		numWorkers = 3
	}
	jobs := make(chan int, len(*recruitProfiles))
	results := make(chan error, len(*recruitProfiles))

	// This starts up numWorkers number of workers, initially blocked because there are no jobs yet.
	for w := 1; w <= numWorkers; w++ {
		go func(jobs <-chan int, results chan<- error, w int) {
			for i := range jobs {
				if (*recruitProfiles)[i].CurrentWeeksPoints == 0 {
					results <- nil
					continue
				}
				abbr := (*recruitProfiles)[i].ProfileID
				mapMutex.Lock()
				bonus := teamMap[strconv.Itoa(int(abbr))].BonusPoints
				mapMutex.Unlock()
				err := processRecruitProfile(i, recruit, recruitProfiles, float64(bonus), pointLimit, pointsPlaced, timestamp, recruitProfilePointsMap, &mapMutex, db)
				results <- err
			}
		}(jobs, results, w)
	}

	// Here we send len(*recruitProfiles) jobs and then close the channel.
	for i := 0; i < len(*recruitProfiles); i++ {
		jobs <- i
	}
	close(jobs)

	// Finally, we collect all the results.
	// This ensures the function doesn't return until we've processed all recruit profiles.
	for i := 0; i < len(*recruitProfiles); i++ {
		err := <-results
		if err != nil {
			fmt.Println(err)
			log.Fatalf("Could not process recruit profile: %v", err)
		}
	}
}

func processRecruitProfile(i int, recruit structs.Recruit, recruitProfiles *[]structs.PlayerRecruitProfile, bonus, pointLimit float64, pointsPlaced *bool, timestamp structs.Timestamp, recruitProfilePointsMap *map[string]float64, m *sync.Mutex, db *gorm.DB) error {
	regionBonus := 1.05
	stateBonus := 1.1
	*pointsPlaced = true

	rpa := structs.RecruitPointAllocation{
		RecruitID:        (*recruitProfiles)[i].RecruitID,
		TeamProfileID:    (*recruitProfiles)[i].ProfileID,
		RecruitProfileID: (*recruitProfiles)[i].ID,
		WeekID:           timestamp.CollegeWeekID,
	}

	var curr float64 = float64((*recruitProfiles)[i].CurrentWeeksPoints)

	// Region / State bonus
	if (*recruitProfiles)[i].HasRegionBonus && recruit.Stars != 5 {
		curr = curr * regionBonus
	} else if (*recruitProfiles)[i].HasStateBonus && recruit.Stars != 5 {
		curr = curr * stateBonus
	}
	// Bonus Points value when saving
	if curr > 0 {
		curr += bonus
	}

	if (*recruitProfiles)[i].CurrentWeeksPoints < 0 || (*recruitProfiles)[i].CurrentWeeksPoints > 20 {
		curr = 0
		rpa.ApplyCaughtCheating()
	}

	rpa.UpdatePointsSpent(float64((*recruitProfiles)[i].CurrentWeeksPoints), curr)
	(*recruitProfiles)[i].AllocateTotalPoints(curr)

	m.Lock()
	(*recruitProfilePointsMap)[(*recruitProfiles)[i].TeamAbbreviation] += float64((*recruitProfiles)[i].CurrentWeeksPoints)
	m.Unlock()

	// Add RPA to point allocations list
	err := db.Create(&rpa).Error
	if err != nil {
		return fmt.Errorf("could not save point allocation: %v", err)
	}
	return nil
}

func updateTeamRankings(teamRecruitingProfiles []structs.TeamRecruitingProfile, teamMap map[string]*structs.TeamRecruitingProfile, recruitProfilePointsMap map[string]float64, db *gorm.DB) {
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

		teamRecruitingProfiles[i].Assign247Rank(team247Rank)
		teamRecruitingProfiles[i].AssignESPNRank(teamESPNRank)
		teamRecruitingProfiles[i].AssignRivalsRank(teamRivalsRank)
	}

	espnDivisor := (maxESPNScore - minESPNScore)
	divisor247 := (max247Score - min247Score)
	rivalsDivisor := (maxRivalsScore - minRivalsScore)
	for _, rp := range teamRecruitingProfiles {
		if recruitProfilePointsMap[rp.TeamAbbr] > float64(rp.WeeklyPoints) {
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
		repository.SaveCBBTeamRecruitingProfile(rp, db)
		fmt.Println("Saved Rank Scores for Team " + rp.TeamAbbr)
	}
}

func MakeRecruitProfileMapByRecruitID(profiles []structs.PlayerRecruitProfile) map[uint][]structs.PlayerRecruitProfile {
	recruitProfileMap := make(map[uint][]structs.PlayerRecruitProfile)

	for _, profile := range profiles {
		if profile.RemovedFromBoard {
			continue
		}
		if len(recruitProfileMap[profile.RecruitID]) == 0 {
			recruitProfileMap[profile.RecruitID] = []structs.PlayerRecruitProfile{profile}
		} else {
			recruitProfileMap[profile.RecruitID] = append(recruitProfileMap[profile.RecruitID], profile)
		}
	}

	return recruitProfileMap
}

func MakeRecruitProfileMapByProfileID(profiles []structs.PlayerRecruitProfile) map[uint][]structs.PlayerRecruitProfile {
	recruitProfileMap := make(map[uint][]structs.PlayerRecruitProfile)

	for _, profile := range profiles {
		if profile.RemovedFromBoard {
			continue
		}
		if len(recruitProfileMap[profile.ProfileID]) == 0 {
			recruitProfileMap[profile.ProfileID] = []structs.PlayerRecruitProfile{profile}
		} else {
			recruitProfileMap[profile.ProfileID] = append(recruitProfileMap[profile.ProfileID], profile)
		}
	}

	return recruitProfileMap
}
