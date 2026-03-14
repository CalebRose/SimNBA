package managers

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/repository"
	"github.com/CalebRose/SimNBA/structs"
)

// ShuffleTeams randomizes the order of teams in a slice
func ShuffleTeams(teams []structs.Team) []structs.Team {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	shuffled := make([]structs.Team, len(teams))
	copy(shuffled, teams)
	rng.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return shuffled
}

func GenerateOOCScheduleToCSV() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.SeasonID))
	collegeTeams := GetAllActiveCollegeTeams()

	maxAttempts := 5000
	var matchesToUpload []structs.Match
	var err error

	// Try multiple times to generate a complete schedule
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		matchesToUpload, err = attemptGenerateOOCSchedule(ts, seasonID, collegeTeams)
		if err == nil {
			fmt.Printf("Successfully generated OOC schedule on attempt %d with %d matches.\n", attempt, len(matchesToUpload))
			break
		}
		if attempt < maxAttempts {
			fmt.Printf("Attempt %d failed: %v. Retrying...\n", attempt, err)
		} else {
			fmt.Printf("Failed to generate complete schedule after %d attempts. Last error: %v\n", maxAttempts, err)
			fmt.Printf("Uploading partial schedule with %d matches.\n", len(matchesToUpload))
		}
	}

	repository.CreateCollegeMatchesRecordsBatch(db, matchesToUpload, 100)
}

func GenerateOOCSchedule() {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.SeasonID))
	collegeTeams := GetAllActiveCollegeTeams()

	maxAttempts := 5000
	var matchesToUpload []structs.Match
	var err error

	// Try multiple times to generate a complete schedule
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		matchesToUpload, err = attemptGenerateOOCSchedule(ts, seasonID, collegeTeams)
		if err == nil {
			fmt.Printf("Successfully generated OOC schedule on attempt %d with %d matches.\n", attempt, len(matchesToUpload))
			break
		}
		if attempt < maxAttempts {
			fmt.Printf("Attempt %d failed: %v. Retrying...\n", attempt, err)
		} else {
			fmt.Printf("Failed to generate complete schedule after %d attempts. Last error: %v\n", maxAttempts, err)
			fmt.Printf("Uploading partial schedule with %d matches.\n", len(matchesToUpload))
		}
	}

	// Upload matches to database
	for _, match := range matchesToUpload {
		err := db.Create(&match).Error
		if err != nil {
			fmt.Printf("Error uploading match: %v\n", err)
		}
	}

	fmt.Printf("Schedule generation complete. Total matches created: %d\n", len(matchesToUpload))
}

func attemptGenerateOOCSchedule(ts structs.Timestamp, seasonID string, collegeTeams []structs.Team) ([]structs.Match, error) {
	maxNumberOfGamesPerTeam := 30
	oocWeeks := []uint{15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1} // Schedule backwards: fill constrained weeks (7-15) first
	timeSlots := []string{"A", "B"}

	// Initialize tracking maps
	remainingGamesLeftByTeam := make(map[uint]int)
	timeSlotMap := make(map[uint]map[uint]map[string]bool)
	homeGamesByTeam := make(map[uint]int)
	awayGamesByTeam := make(map[uint]int)
	opponentsFacedMap := make(map[uint]map[uint]bool) // teamID -> opponentTeamID -> bool
	// Count the number of games per week and time slot to ensure we are filling them correctly
	weekTimeSlotCounts := make(map[uint]map[string]int)
	numberOfGamesExpectedPerSlot := len(collegeTeams) / 2 // Each game has 2 teams
	gamesScheduledPerTeam := make(map[uint]int)

	// Create a map to track user vs AI teams
	isUserTeam := make(map[uint]bool)
	for _, team := range collegeTeams {
		remainingGamesLeftByTeam[team.ID] = maxNumberOfGamesPerTeam
		timeSlotMap[team.ID] = make(map[uint]map[string]bool)
		homeGamesByTeam[team.ID] = 0
		awayGamesByTeam[team.ID] = 0
		opponentsFacedMap[team.ID] = make(map[uint]bool)
		// Assume teams with a coach are user teams
		isUserTeam[team.ID] = team.IsUserCoached
	}

	// Process existing matches to update counters
	allActiveMatches := GetCBBMatchesBySeasonID(seasonID)
	// Check all active games to ensure no team is scheduled twice in the same timeslot
	// teamSlotSeen: teamID -> "week:timeSlot" -> matchID
	teamSlotSeen := make(map[uint]map[string]uint)
	for _, match := range allActiveMatches {
		if match.AwayTeamID == 0 && match.HomeTeamID == 0 {
			continue
		}
		if weekTimeSlotCounts[match.Week] == nil {
			weekTimeSlotCounts[match.Week] = make(map[string]int)
		}
		weekTimeSlotCounts[match.Week][match.MatchOfWeek]++
		slotKey := fmt.Sprintf("%d:%s", match.Week, match.MatchOfWeek)
		for _, teamID := range []uint{match.HomeTeamID, match.AwayTeamID} {
			if teamSlotSeen[teamID] == nil {
				teamSlotSeen[teamID] = make(map[string]uint)
			}
			if firstMatchID, exists := teamSlotSeen[teamID][slotKey]; exists {
				fmt.Printf("DOUBLE-SCHEDULED: TeamID %d appears in Week %d TimeSlot %s in both MatchID %d and MatchID %d\n",
					teamID, match.Week, match.MatchOfWeek, firstMatchID, match.ID)
			} else {
				teamSlotSeen[teamID][slotKey] = match.ID
			}
		}
	}

	for _, match := range allActiveMatches {
		if match.AwayTeamID == 0 && match.HomeTeamID == 0 {
			continue
		}
		// Initialize nested maps if needed, guarding against team IDs not present in collegeTeams
		if timeSlotMap[match.HomeTeamID] == nil {
			timeSlotMap[match.HomeTeamID] = make(map[uint]map[string]bool)
		}
		if timeSlotMap[match.HomeTeamID][match.Week] == nil {
			timeSlotMap[match.HomeTeamID][match.Week] = make(map[string]bool)
		}
		if timeSlotMap[match.AwayTeamID] == nil {
			timeSlotMap[match.AwayTeamID] = make(map[uint]map[string]bool)
		}
		if timeSlotMap[match.AwayTeamID][match.Week] == nil {
			timeSlotMap[match.AwayTeamID][match.Week] = make(map[string]bool)
		}

		gamesScheduledPerTeam[match.HomeTeamID]++
		gamesScheduledPerTeam[match.AwayTeamID]++

		// Mark time slots as occupied
		timeSlotMap[match.HomeTeamID][match.Week][match.MatchOfWeek] = true
		timeSlotMap[match.AwayTeamID][match.Week][match.MatchOfWeek] = true
		if opponentsFacedMap[match.HomeTeamID] == nil {
			opponentsFacedMap[match.HomeTeamID] = make(map[uint]bool)
		}
		if opponentsFacedMap[match.AwayTeamID] == nil {
			opponentsFacedMap[match.AwayTeamID] = make(map[uint]bool)
		}
		opponentsFacedMap[match.HomeTeamID][match.AwayTeamID] = true
		opponentsFacedMap[match.AwayTeamID][match.HomeTeamID] = true

		// Update remaining games and home/away counts
		remainingGamesLeftByTeam[match.HomeTeamID]--
		remainingGamesLeftByTeam[match.AwayTeamID]--
		homeGamesByTeam[match.HomeTeamID]++
		awayGamesByTeam[match.AwayTeamID]++
	}

	matchesToUpload := []structs.Match{}

	// Iterate through each week and time slot
	for _, week := range oocWeeks {
		for _, timeSlot := range timeSlots {
			// Build list of available teams for this week/timeslot
			var availableTeams []structs.Team
			for _, team := range collegeTeams {
				// Check if team needs games and has the time slot available
				if remainingGamesLeftByTeam[team.ID] > 0 {
					if timeSlotMap[team.ID][week] == nil || !timeSlotMap[team.ID][week][timeSlot] {
						availableTeams = append(availableTeams, team)
					}
				}
			}

			// Shuffle available teams for randomness
			availableTeams = ShuffleTeams(availableTeams)

			// Track which teams have been matched in this slot
			matchedInSlot := make(map[uint]bool)

			// Try to pair teams
			for i := 0; i < len(availableTeams); i++ {
				team := availableTeams[i]

				// Skip if already matched in this slot
				if matchedInSlot[team.ID] {
					continue
				}

				// Find an opponent
				opponentIndex := -1
				for j := i + 1; j < len(availableTeams); j++ {
					opponent := availableTeams[j]

					// Check if valid opponent: different conference, not yet matched, has games remaining, and haven't faced each other
					if !matchedInSlot[opponent.ID] &&
						opponent.ConferenceID != team.ConferenceID &&
						remainingGamesLeftByTeam[opponent.ID] > 0 &&
						!opponentsFacedMap[team.ID][opponent.ID] {
						opponentIndex = j
						break
					}
				}

				// If no opponent found, this attempt fails
				if opponentIndex == -1 {
					continue // Try next team
				}

				opponent := availableTeams[opponentIndex]

				// Decide which team is home based on home/away balance
				var homeTeam, awayTeam structs.Team

				// Prefer to balance user teams first
				if isUserTeam[team.ID] && !isUserTeam[opponent.ID] {
					// User team: prioritize balance
					if homeGamesByTeam[team.ID] <= awayGamesByTeam[team.ID] {
						homeTeam, awayTeam = team, opponent
					} else {
						homeTeam, awayTeam = opponent, team
					}
				} else if !isUserTeam[team.ID] && isUserTeam[opponent.ID] {
					// Opponent is user team: prioritize their balance
					if homeGamesByTeam[opponent.ID] <= awayGamesByTeam[opponent.ID] {
						homeTeam, awayTeam = opponent, team
					} else {
						homeTeam, awayTeam = team, opponent
					}
				} else {
					// Both user or both AI: balance based on current counts
					teamHomeDeficit := homeGamesByTeam[team.ID] - awayGamesByTeam[team.ID]
					oppHomeDeficit := homeGamesByTeam[opponent.ID] - awayGamesByTeam[opponent.ID]

					if teamHomeDeficit < oppHomeDeficit {
						homeTeam, awayTeam = team, opponent
					} else if teamHomeDeficit > oppHomeDeficit {
						homeTeam, awayTeam = opponent, team
					} else {
						// Equal deficit, randomize
						if rand.Intn(2) == 0 {
							homeTeam, awayTeam = team, opponent
						} else {
							homeTeam, awayTeam = opponent, team
						}
					}
				}

				// Create match
				match := structs.Match{
					SeasonID:      ts.SeasonID,
					Week:          week,
					WeekID:        ts.CollegeWeekID + week,
					MatchOfWeek:   timeSlot,
					HomeTeamID:    homeTeam.ID,
					HomeTeam:      homeTeam.Abbr,
					HomeTeamCoach: homeTeam.Coach,
					Arena:         homeTeam.Arena,
					City:          homeTeam.City,
					State:         homeTeam.State,
					AwayTeamID:    awayTeam.ID,
					AwayTeam:      awayTeam.Abbr,
					AwayTeamCoach: awayTeam.Coach,
					IsConference:  homeTeam.ConferenceID == awayTeam.ConferenceID,
				}

				matchesToUpload = append(matchesToUpload, match)

				// Update tracking data
				matchedInSlot[team.ID] = true
				matchedInSlot[opponent.ID] = true

				if timeSlotMap[homeTeam.ID][week] == nil {
					timeSlotMap[homeTeam.ID][week] = make(map[string]bool)
				}
				if timeSlotMap[awayTeam.ID][week] == nil {
					timeSlotMap[awayTeam.ID][week] = make(map[string]bool)
				}

				timeSlotMap[homeTeam.ID][week][timeSlot] = true
				timeSlotMap[awayTeam.ID][week][timeSlot] = true

				remainingGamesLeftByTeam[homeTeam.ID]--
				remainingGamesLeftByTeam[awayTeam.ID]--

				homeGamesByTeam[homeTeam.ID]++
				awayGamesByTeam[awayTeam.ID]++

				// Mark that these teams have faced each other
				opponentsFacedMap[homeTeam.ID][awayTeam.ID] = true
				opponentsFacedMap[awayTeam.ID][homeTeam.ID] = true
			}

			// Log teams that were available this slot but could not be paired
			for _, team := range availableTeams {
				if !matchedInSlot[team.ID] && remainingGamesLeftByTeam[team.ID] > 0 {
					fmt.Printf("Week %d, Slot %s: Could not schedule TeamID %d (%s)\n", week, timeSlot, team.ID, team.Abbr)
				}
			}
		}
	}

	// Cleanup pass: schedule any remaining games with relaxed constraints.
	// The only requirement is that the two teams are from different conferences;
	// the opponentsFacedMap check is intentionally dropped to allow rematches.
	fmt.Println("Starting cleanup scheduling pass for remaining unscheduled games...")
	for _, week := range oocWeeks {
		for _, timeSlot := range timeSlots {
			var cleanupAvailable []structs.Team
			for _, team := range collegeTeams {
				if remainingGamesLeftByTeam[team.ID] > 0 {
					if timeSlotMap[team.ID][week] == nil || !timeSlotMap[team.ID][week][timeSlot] {
						cleanupAvailable = append(cleanupAvailable, team)
					}
				}
			}

			if len(cleanupAvailable) < 2 {
				continue
			}

			cleanupAvailable = ShuffleTeams(cleanupAvailable)
			cleanupMatchedInSlot := make(map[uint]bool)

			for i := 0; i < len(cleanupAvailable); i++ {
				team := cleanupAvailable[i]
				if cleanupMatchedInSlot[team.ID] {
					continue
				}

				opponentIndex := -1
				for j := i + 1; j < len(cleanupAvailable); j++ {
					opponent := cleanupAvailable[j]
					// Relaxed: only require different conference and no previous matchup
					if !cleanupMatchedInSlot[opponent.ID] &&
						opponent.ConferenceID != team.ConferenceID &&
						remainingGamesLeftByTeam[opponent.ID] > 0 &&
						!opponentsFacedMap[team.ID][opponent.ID] {
						opponentIndex = j
						break
					}
				}

				if opponentIndex == -1 {
					continue
				}

				opponent := cleanupAvailable[opponentIndex]

				var homeTeam, awayTeam structs.Team
				if isUserTeam[team.ID] && !isUserTeam[opponent.ID] {
					if homeGamesByTeam[team.ID] <= awayGamesByTeam[team.ID] {
						homeTeam, awayTeam = team, opponent
					} else {
						homeTeam, awayTeam = opponent, team
					}
				} else if !isUserTeam[team.ID] && isUserTeam[opponent.ID] {
					if homeGamesByTeam[opponent.ID] <= awayGamesByTeam[opponent.ID] {
						homeTeam, awayTeam = opponent, team
					} else {
						homeTeam, awayTeam = team, opponent
					}
				} else {
					teamHomeDeficit := homeGamesByTeam[team.ID] - awayGamesByTeam[team.ID]
					oppHomeDeficit := homeGamesByTeam[opponent.ID] - awayGamesByTeam[opponent.ID]
					if teamHomeDeficit < oppHomeDeficit {
						homeTeam, awayTeam = team, opponent
					} else if teamHomeDeficit > oppHomeDeficit {
						homeTeam, awayTeam = opponent, team
					} else {
						if rand.Intn(2) == 0 {
							homeTeam, awayTeam = team, opponent
						} else {
							homeTeam, awayTeam = opponent, team
						}
					}
				}

				match := structs.Match{
					SeasonID:      ts.SeasonID,
					Week:          week,
					WeekID:        ts.CollegeWeekID + week,
					MatchOfWeek:   timeSlot,
					HomeTeamID:    homeTeam.ID,
					HomeTeam:      homeTeam.Abbr,
					HomeTeamCoach: homeTeam.Coach,
					Arena:         homeTeam.Arena,
					City:          homeTeam.City,
					State:         homeTeam.State,
					AwayTeamID:    awayTeam.ID,
					AwayTeam:      awayTeam.Abbr,
					AwayTeamCoach: awayTeam.Coach,
					IsConference:  homeTeam.ConferenceID == awayTeam.ConferenceID,
				}

				matchesToUpload = append(matchesToUpload, match)

				cleanupMatchedInSlot[team.ID] = true
				cleanupMatchedInSlot[opponent.ID] = true

				if timeSlotMap[homeTeam.ID][week] == nil {
					timeSlotMap[homeTeam.ID][week] = make(map[string]bool)
				}
				if timeSlotMap[awayTeam.ID][week] == nil {
					timeSlotMap[awayTeam.ID][week] = make(map[string]bool)
				}

				timeSlotMap[homeTeam.ID][week][timeSlot] = true
				timeSlotMap[awayTeam.ID][week][timeSlot] = true

				remainingGamesLeftByTeam[homeTeam.ID]--
				remainingGamesLeftByTeam[awayTeam.ID]--

				homeGamesByTeam[homeTeam.ID]++
				awayGamesByTeam[awayTeam.ID]++

				opponentsFacedMap[homeTeam.ID][awayTeam.ID] = true
				opponentsFacedMap[awayTeam.ID][homeTeam.ID] = true

				fmt.Printf("Cleanup: Scheduled Week %d, Slot %s: TeamID %d (%s) vs TeamID %d (%s)\n",
					week, timeSlot, homeTeam.ID, homeTeam.Abbr, awayTeam.ID, awayTeam.Abbr)
			}
		}
	}

	// Validate the schedule - check if all teams that need OOC games got enough
	expectedGamesPerTeam := len(oocWeeks) * len(timeSlots) // 4 weeks * 2 slots = 8 games expected per team

	for _, match := range matchesToUpload {
		gamesScheduledPerTeam[match.HomeTeamID]++
		gamesScheduledPerTeam[match.AwayTeamID]++
	}

	// Count teams with insufficient games
	insufficientTeams := 0
	for _, team := range collegeTeams {
		if gamesScheduledPerTeam[team.ID] < expectedGamesPerTeam {
			insufficientTeams++
		}
	}

	for _, match := range matchesToUpload {
		if weekTimeSlotCounts[match.Week] == nil {
			weekTimeSlotCounts[match.Week] = make(map[string]int)
		}
		weekTimeSlotCounts[match.Week][match.MatchOfWeek]++
	}

	// Print out the distribution of games per week and time slot
	for week, slotCounts := range weekTimeSlotCounts {
		for slot, count := range slotCounts {
			fmt.Printf("Week %d, Time Slot %s: %d games scheduled (expected ~%d)\n", week, slot, count, numberOfGamesExpectedPerSlot)
		}
	}

	// If too many teams don't have enough games, consider this attempt failed
	// No, we need a perfect schedule generated.
	tolerance := 0
	if insufficientTeams > tolerance {
		return matchesToUpload, fmt.Errorf("incomplete schedule: %d teams have insufficient games", insufficientTeams)
	}

	return matchesToUpload, nil
}
