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
	rand.Seed(time.Now().UnixNano())
	//GetCurrentWeek

	if timestamp.RecruitingSynced {
		log.Fatalln("Recruiting already ran for this week. Please wait until next week to sync recruiting again.")
	}

	recruitProfilePointsMap := make(map[string]int)

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

		var totalPointsOnRecruit float64 = 0

		var eligiblePointThreshold float64 = 0

		var signThreshold float64

		for i := 0; i < len(recruitProfiles); i++ {

			if recruitProfiles[i].CurrentWeeksPoints == 0 {
				continue
			}

			rpa := structs.RecruitPointAllocation{
				RecruitID:        recruitProfiles[i].RecruitID,
				TeamProfileID:    recruitProfiles[i].ProfileID,
				RecruitProfileID: recruitProfiles[i].ID,
				WeekID:           timestamp.CollegeWeekID,
			}

			var curr float64 = 0

			// Region / State bonus
			if recruitProfiles[i].HasRegionBonus {
				curr = float64(recruitProfiles[i].CurrentWeeksPoints) * 1.1
			} else if recruitProfiles[i].HasStateBonus {
				curr = float64(recruitProfiles[i].CurrentWeeksPoints) * 1.2
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

		sort.Sort(structs.ByPoints(recruitProfiles))

		for i := 0; i < len(recruitProfiles); i++ {
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
		if totalPointsOnRecruit > signThreshold && eligibleTeams > 0 {
			var winningTeamID uint = 0
			var odds float64 = 0

			for winningTeamID == 0 {
				percentageOdds := 1 + rand.Float64()*(totalPointsOnRecruit-1)
				var currentProbability float64 = 0

				for i := 0; i < len(recruitProfilesWithScholarship); i++ {
					// If a team has no available scholarships or if a team has 25 commitments, continue
					currentProbability += recruitProfilesWithScholarship[i].TotalPoints
					if currentProbability >= float64(percentageOdds) {
						// WINNING TEAM
						winningTeamID = recruitProfilesWithScholarship[i].ProfileID
						odds = float64(recruitProfilesWithScholarship[i].TotalPoints) / float64(totalPointsOnRecruit) * 100
						break
					}
				}

				if winningTeamID > 0 {
					recruitTeamProfile := GetOnlyTeamRecruitingProfileByTeamID(strconv.Itoa(int(winningTeamID)))
					if recruitTeamProfile.TotalCommitments < recruitTeamProfile.RecruitClassSize {
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

						for i := 0; i < len(recruitProfiles); i++ {
							if recruitProfiles[i].ProfileID == winningTeamID {
								recruitProfiles[i].SignPlayer()
							} else {
								recruitProfiles[i].LockPlayer()
								if recruitProfiles[i].Scholarship {
									tp := GetOnlyTeamRecruitingProfileByTeamID(strconv.Itoa(int(recruitProfiles[i].ProfileID)))

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
						if len(recruitProfilesWithScholarship) == 0 {
							break
						}
						winningTeamID = 0

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
	teamRecruitingProfiles := GetTeamRecruitingProfilesForRecruitSync()

	var totalESPNScore float64 = 0
	var total247Score float64 = 0
	var totalRivalsScore float64 = 0

	for i := 0; i < len(teamRecruitingProfiles); i++ {

		signedRecruits := GetSignedRecruitsByTeamProfileID(strconv.Itoa(int(teamRecruitingProfiles[i].TeamID)))

		teamRecruitingProfiles[i].UpdateTotalSignedRecruits(len(signedRecruits))

		team247Rank := util.Get247TeamRanking(teamRecruitingProfiles[i], signedRecruits)
		teamESPNRank := util.GetESPNTeamRanking(teamRecruitingProfiles[i], signedRecruits)
		teamRivalsRank := util.GetRivalsTeamRanking(teamRecruitingProfiles[i], signedRecruits)

		teamRecruitingProfiles[i].Assign247Rank(team247Rank)
		total247Score += team247Rank
		teamRecruitingProfiles[i].AssignESPNRank(teamESPNRank)
		totalESPNScore += teamESPNRank
		teamRecruitingProfiles[i].AssignRivalsRank(teamRivalsRank)
		totalRivalsScore += teamRivalsRank

		fmt.Println("Setting Recruiting Ranks for " + teamRecruitingProfiles[i].TeamAbbr)

	}

	averageESPNScore := totalESPNScore / 130
	average247score := total247Score / 130
	averageRivalScore := totalRivalsScore / 130

	for _, rp := range teamRecruitingProfiles {
		if recruitProfilePointsMap[rp.TeamAbbr] > rp.WeeklyPoints {
			rp.ApplyCaughtCheating()
		}

		var avg float64 = 0
		if averageESPNScore > 0 && average247score > 0 && averageRivalScore > 0 {
			distributionESPN := rp.ESPNScore / averageESPNScore
			distribution247 := rp.Rank247Score / average247score
			distributionRivals := rp.RivalsScore / averageRivalScore

			avg = (distributionESPN + distribution247 + distributionRivals) / 3

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
	rand.Seed(time.Now().UnixNano())
	ts := GetTimestamp()

	AITeams := GetOnlyAITeamRecruitingProfiles()

	UnsignedRecruits := GetAllUnsignedRecruits()

	regionMap := util.GetRegionMap()

	boardCount := 40

	if ts.CollegeWeek > 3 {
		boardCount = 15
	}

	for _, team := range AITeams {
		count := 0
		if !team.IsAI {
			continue
		}

		existingBoard := GetAllRecruitsByProfileID(strconv.Itoa(int(team.ID)))

		count = len(existingBoard)

		if count > boardCount {
			continue
		}

		for _, croot := range UnsignedRecruits {
			if count == boardCount {
				break
			}
			if croot.Stars == 5 || croot.IsCustomCroot ||
				(croot.Stars == 4 && team.AIQuality != "Blue Blood") ||
				(croot.Stars > 2 && team.AIQuality == "Bottom Feeder") {
				continue
			}

			crootProfile := GetPlayerRecruitProfileByPlayerId(strconv.Itoa(int(croot.ID)), strconv.Itoa(int(team.ID)))
			if crootProfile.RemovedFromBoard || crootProfile.IsLocked {
				continue
			}

			crootProfiles := GetRecruitPlayerProfilesByRecruitId(strconv.Itoa(int(croot.ID)))

			leadingVal := util.IsAITeamContendingForCroot(crootProfiles)
			if leadingVal > 11 {
				continue
			}

			odds := 5

			if croot.Country == "USA" {
				if regionMap[croot.State] == team.Region {
					odds = 25
				}
				if croot.State == team.State {
					odds = 33
				}
			}

			if team.AIQuality == "Offense" && util.IsPlayerOffensivelyStrong(croot) {
				odds += 50
			} else if team.AIQuality == "Defense" && util.IsPlayerDefensivelyStrong(croot) {
				odds += 50
			} else if team.AIQuality == "Cinderella" && util.IsPlayerHighPotential(croot) {
				odds += 50
			} else if team.AIQuality == "Blue Blood" && croot.Stars == 4 {
				odds += 50
			} else if team.AIQuality == "Bottom Feeder" && croot.Stars < 3 {
				odds += 40
			}

			chance := util.GenerateIntFromRange(1, 100)

			if chance <= odds {
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
	rand.Seed(time.Now().UnixNano())
	ts := GetTimestamp()

	AITeams := GetOnlyAITeamRecruitingProfiles()

	for _, team := range AITeams {
		if team.SpentPoints >= team.WeeklyPoints {
			continue
		}

		teamRecruits := GetAllRecruitsByProfileID(strconv.Itoa(int(team.ID)))

		for _, croot := range teamRecruits {
			// If a team has no more points to spend, break the loop
			if team.SpentPoints >= team.WeeklyPoints {
				break
			}
			// If a croot was signed OR has points already placed on the croot, move on to the next croot
			if croot.IsSigned || croot.CurrentWeeksPoints > 0 {
				continue
			}
			pointsRemaining := team.WeeklyPoints - team.SpentPoints
			removeCrootFromBoard := false
			num := 0
			// If a croot is locked and signed with a different team, remove from the team board and continue
			if croot.IsLocked && croot.TeamAbbreviation != team.TeamAbbr {
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
						if num >= pointsRemaining {
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
					max := 0

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
					if num >= pointsRemaining {
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
