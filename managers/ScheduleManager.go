package managers

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/CalebRose/SimNBA/dbprovider"
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

	// Write Results to CSV
	filename := fmt.Sprintf("./data/%d/%d_SimCBB_OOC_Games.csv", ts.Season, ts.Season)
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating CSV file: %v\n", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header row
	headerRow := []string{
		"MatchID", "Season", "SeasonID", "Week", "WeekID", "TimeSlot",
		"HomeTeamID", "HomeTeam", "HomeTeamCoach", "HomeTeamRank",
		"Arena", "City", "State",
		"AwayTeamID", "AwayTeam", "AwayTeamCoach", "AwayTeamRank",
		"IsNeutralSite", "IsInvitational", "IsConference", "IsConferenceTournament",
		"IsCBI", "IsNIT", "IsTournament", "IsNationalChampionship",
		"GameTitle", "NextGameID", "HomeOrAway",
	}

	err = writer.Write(headerRow)
	if err != nil {
		fmt.Printf("Error writing CSV header: %v\n", err)
		return
	}

	// Write match rows
	for idx, match := range matchesToUpload {
		matchID := idx + 1 // Generate sequential IDs
		isConferenceStr := "FALSE"
		if match.IsConference {
			isConferenceStr = "TRUE"
		}
		matchRow := []string{
			strconv.Itoa(matchID),
			strconv.Itoa(ts.Season),
			strconv.Itoa(int(match.SeasonID)),
			strconv.Itoa(int(match.Week)),
			strconv.Itoa(int(match.WeekID)),
			match.MatchOfWeek,
			strconv.Itoa(int(match.HomeTeamID)),
			match.HomeTeam,
			match.HomeTeamCoach,
			"", // HomeTeamRank (empty for OOC)
			match.Arena,
			match.City,
			match.State,
			strconv.Itoa(int(match.AwayTeamID)),
			match.AwayTeam,
			match.AwayTeamCoach,
			"",              // AwayTeamRank (empty for OOC)
			"FALSE",         // IsNeutralSite
			"FALSE",         // IsInvitational
			isConferenceStr, // IsConference
			"FALSE",         // IsConferenceTournament
			"FALSE",         // IsCBI
			"FALSE",         // IsNIT
			"FALSE",         // IsTournament
			"FALSE",         // IsNationalChampionship
			"",              // GameTitle
			"",              // NextGameID
			"",              // HomeOrAway
		}

		err = writer.Write(matchRow)
		if err != nil {
			fmt.Printf("Error writing match row %d: %v\n", matchID, err)
		}
	}

	fmt.Printf("Schedule generation complete. Total matches created: %d\n", len(matchesToUpload))
	fmt.Printf("CSV exported to: %s\n", filename)
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
	oocWeeks := []uint{1, 2, 3, 4} // Only weeks 1-4 for OOC games
	timeSlots := []string{"A", "B"}

	// Initialize tracking maps
	remainingGamesLeftByTeam := make(map[uint]int)
	timeSlotMap := make(map[uint]map[uint]map[string]bool)
	homeGamesByTeam := make(map[uint]int)
	awayGamesByTeam := make(map[uint]int)
	opponentsFacedMap := make(map[uint]map[uint]bool) // teamID -> opponentTeamID -> bool

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
	for _, match := range allActiveMatches {
		// Initialize nested maps if needed
		if timeSlotMap[match.HomeTeamID][match.Week] == nil {
			timeSlotMap[match.HomeTeamID][match.Week] = make(map[string]bool)
		}
		if timeSlotMap[match.AwayTeamID][match.Week] == nil {
			timeSlotMap[match.AwayTeamID][match.Week] = make(map[string]bool)
		}

		// Mark time slots as occupied
		timeSlotMap[match.HomeTeamID][match.Week][match.TimeSlot] = true
		timeSlotMap[match.AwayTeamID][match.Week][match.TimeSlot] = true
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
		}
	}

	// Validate the schedule - check if all teams that need OOC games got enough
	expectedGamesPerTeam := len(oocWeeks) * len(timeSlots) // 4 weeks * 2 slots = 8 games expected per team
	gamesScheduledPerTeam := make(map[uint]int)

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
	// Count the number of games per week and time slot to ensure we are filling them correctly
	weekTimeSlotCounts := make(map[uint]map[string]int)
	numberOfGamesExpectedPerSlot := len(collegeTeams) / 2 // Each game has 2 teams

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
