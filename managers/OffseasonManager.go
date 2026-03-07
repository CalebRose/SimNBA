package managers

import (
	"math"
	"strconv"
	"strings"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/repository"
	"github.com/CalebRose/SimNBA/structs"
)

func UpdateTeamProfileAffinities() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	teamRecruitingProfiles := GetTeamRecruitingProfilesForRecruitSync()
	teamProfileMap := MakeTeamRecruitingProfileMapByTeamID(teamRecruitingProfiles)
	collegeTeams := GetAllActiveCollegeTeams()
	collegeGames := repository.FindCollegeMatchRecords(repository.GameQuery{ExcludeCompletedGames: false})
	collegeStandings := repository.FindAllCollegeStandingsRecords(repository.StandingsQuery{})
	collegeGamesByCoachMap := make(map[string][]structs.Match)
	collegeGamesByTeamIDMap := make(map[uint][]structs.Match)
	affinitiesMap := make(map[uint]structs.ProfileAttributes)

	for _, game := range collegeGames {
		if !game.GameComplete {
			continue
		}
		if game.HomeTeamCoach != "" && game.HomeTeamCoach != "AI" {
			collegeGamesByCoachMap[game.HomeTeamCoach] = append(collegeGamesByCoachMap[game.HomeTeamCoach], game)
		}
		if game.AwayTeamCoach != "" && game.AwayTeamCoach != "AI" {
			collegeGamesByCoachMap[game.AwayTeamCoach] = append(collegeGamesByCoachMap[game.AwayTeamCoach], game)
		}
		if game.HomeTeamID > 0 && game.AwayTeamID > 0 {
			collegeGamesByTeamIDMap[uint(game.HomeTeamID)] = append(collegeGamesByTeamIDMap[uint(game.HomeTeamID)], game)
			collegeGamesByTeamIDMap[uint(game.AwayTeamID)] = append(collegeGamesByTeamIDMap[uint(game.AwayTeamID)], game)
		}

	}

	seasonID := ts.SeasonID
	baseProgramDevSeasonID := seasonID - 4

	// Map all historical standings by team
	collegeStandingsMap := make(map[uint][]structs.CollegeStandings)
	for _, standing := range collegeStandings {
		collegeStandingsMap[uint(standing.TeamID)] = append(collegeStandingsMap[uint(standing.TeamID)], standing)
	}

	// Iterate by team
	// get all standings by team.ID
	// iterate by game
	for _, team := range collegeTeams {
		totalWins := 0
		totalLosses := 0
		totalCoachWins := 0
		totalCoachLosses := 0
		seasonMomentumWins := 0
		seasonMomentumLosses := 0
		homeWins := 0
		homeLosses := 0
		programPrestige := 5
		professionalPrestige := 5
		facilities := 5
		academics := 5
		campusLife := 5
		confPrestige := 5
		mediaSpotlight := 5
		teamStandings := collegeStandingsMap[uint(team.ID)]
		collegeGamesByTeam := collegeGamesByTeamIDMap[uint(team.ID)]
		collegeGamesByCoach := collegeGamesByCoachMap[team.Coach]

		// Iterate and track home wins and losses
		for _, game := range collegeGamesByTeam {
			if game.HomeTeamID == team.ID {
				if !game.IsNeutralSite && game.HomeTeamScore > game.AwayTeamScore {
					homeWins++
				} else {
					homeLosses++
				}
				if game.HomeTeamScore > game.AwayTeamScore {
					totalWins++
				} else {
					totalLosses++
				}
			} else if game.AwayTeamID == team.ID {
				if game.AwayTeamScore > game.HomeTeamScore {
					totalWins++
				} else {
					totalLosses++
				}
			}
		}

		tradPercentage := float64(totalWins) / float64(totalWins+totalLosses)
		atmospherePct := float64(homeWins) / float64(homeWins+homeLosses)

		// Iterate by historic standings for program development & seasonMomentum
		for _, standing := range teamStandings {
			// Season momentum
			if standing.SeasonID == ts.SeasonID-1 {
				seasonMomentumWins = standing.TotalWins
				seasonMomentumLosses = standing.TotalLosses
			}
			if standing.SeasonID < baseProgramDevSeasonID {
				continue
			}

			if standing.TotalWins == 0 {
				programPrestige -= 2
			}
			record := float64(standing.TotalWins) / float64(standing.TotalWins+standing.TotalLosses)
			if record < 0.25 {
				programPrestige -= 2
			} else if record < 0.5 {
				programPrestige -= 1
			} else if record >= 0.75 {
				programPrestige += 1
			}
			if standing.TotalLosses == 0 {
				programPrestige += 2
			}

			postSeason := standing.PostSeasonStatus
			postSeasonFlag := 0
			if postSeason == "NIT Champion" || postSeason == "NIT Championship" || postSeason == "CBI Championship" || postSeason == "Conference Champion" {
				postSeasonFlag = 1
			}
			if postSeason == "Sweet 16" || postSeason == "Elite 8" || postSeason == "Final Four" || postSeason == "National Champion Runner-Up" || postSeason == "National Champion" {
				postSeasonFlag = 1
			}
			if postSeason == "National Runner-Up" || postSeason == "National Champion" || postSeason == "Final Four" {
				postSeasonFlag = 2
			}
			if postSeason == "National Champion" {
				postSeasonFlag = 3
			}
			programPrestige += postSeasonFlag

			if programPrestige < 1 {
				programPrestige = 1
			}
			if programPrestige > 10 {
				programPrestige = 10
			}
		}

		seasonMomentumPct := float64(seasonMomentumWins) / float64(seasonMomentumWins+seasonMomentumLosses)

		// Iterate through games by Coach
		if team.Coach != "" && team.Coach != "AI" {
			for _, game := range collegeGamesByCoach {
				if game.HomeTeamCoach == team.Coach {
					if game.HomeTeamScore > game.AwayTeamScore {
						totalCoachWins += 1
					} else {
						totalCoachLosses += 1
					}
				}
				if game.AwayTeamCoach == team.Coach {
					if game.AwayTeamScore > game.HomeTeamScore {
						totalCoachWins += 1
					} else {
						totalCoachLosses += 1
					}
				}
			}
		}

		coachPct := float64(0.0)
		if (totalCoachWins + totalCoachLosses) > 0 {
			coachPct = float64(totalCoachWins) / float64(totalCoachWins+totalCoachLosses)
		}

		newCoachRating := uint8(coachPct * 10)
		if team.Coach == "" || team.Coach == "AI" {
			// Set default to 5
			newCoachRating = 5
		}

		affinitiesMap[uint(team.ID)] = structs.ProfileAttributes{
			ProgramPrestige:      uint8(programPrestige),
			ProfessionalPrestige: uint8(professionalPrestige),
			Traditions:           uint8(tradPercentage * 10),
			Facilities:           uint8(facilities),
			Atmosphere:           uint8(atmospherePct * 10),
			Academics:            uint8(academics),
			CampusLife:           uint8(campusLife),
			ConferencePrestige:   uint8(confPrestige),
			CoachRating:          newCoachRating,
			SeasonMomentum:       uint8(seasonMomentumPct * 10),
			MediaSpotlight:       uint8(mediaSpotlight),
		}
	}

	// Conference Prestige
	// Dynamically build seasonIDs up to the current season to avoid
	// processing unplayed seasons, which would incorrectly apply -1 to
	// every conference due to empty standings data.
	conferenceSeasonIDs := make([]uint, 0, seasonID)
	for i := uint(1); i <= uint(seasonID); i++ {
		conferenceSeasonIDs = append(conferenceSeasonIDs, i)
	}
	// Prefill Prestige Map to 5
	conferencePrestigeMap := make(map[uint]int)
	for _, team := range collegeTeams {
		conferencePrestigeMap[uint(team.ConferenceID)] = 5
	}

	for _, seasonID := range conferenceSeasonIDs {
		conferencePrestigeModMap := make(map[uint]uint)
		stadingsBySeason := repository.FindAllCollegeStandingsRecords(repository.StandingsQuery{SeasonID: strconv.Itoa(int(seasonID))})

		for _, standing := range stadingsBySeason {
			maxConferenceValue := 0
			postSeason := standing.PostSeasonStatus
			if strings.Contains(postSeason, "Sweet 16") || strings.Contains(postSeason, "NIT Champion") || strings.Contains(postSeason, "CBI Championship") || strings.Contains(postSeason, "Round of 32") {
				maxConferenceValue = 1
			} else if strings.Contains(postSeason, "Elite 8") {
				maxConferenceValue = 2
			} else if strings.Contains(postSeason, "Final Four") {
				maxConferenceValue = 3
			} else if strings.Contains(postSeason, "National Champion Runner-Up") || strings.Contains(postSeason, "National Champion") {
				maxConferenceValue = 4
			}
			// Make it the max of the current max value
			conferencePrestigeModMap[uint(standing.ConferenceID)] = uint(math.Max(float64(conferencePrestigeModMap[uint(standing.ConferenceID)]), float64(maxConferenceValue)))
		}

		conferenceIds := []uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 14, 15, 16, 17, 18, 19, 20, 21, 23, 24, 25, 26, 27}
		for _, id := range conferenceIds {
			currentSeasonPrestigeMod := conferencePrestigeModMap[id]
			conferenceMod := 0
			switch currentSeasonPrestigeMod {
			case 0:
				conferenceMod -= 1
			case 1:
				conferenceMod += 0
			case 2:
				conferenceMod += 1
			case 3:
				conferenceMod += 2
			case 4:
				conferenceMod += 3
			}
			newPrestige := conferencePrestigeMap[id] + conferenceMod
			conferencePrestigeMap[id] = newPrestige
		}
	}

	for _, team := range collegeTeams {
		teamProfile := teamProfileMap[team.ID]
		affinities := affinitiesMap[team.ID]
		conferencePrestige := conferencePrestigeMap[uint(team.ConferenceID)]
		if conferencePrestige < 1 {
			conferencePrestige = 1
		} else if conferencePrestige > 10 {
			conferencePrestige = 10
		}

		if team.ConferenceID > 13 && conferencePrestige < 1 {
			conferencePrestige = 5
		}
		affinities.ConferencePrestige = uint8(conferencePrestige)
		teamProfile.UpdateTeamProfileAffinities(affinities)
		repository.SaveCBBTeamRecruitingProfile(teamProfile, db)
	}
}
