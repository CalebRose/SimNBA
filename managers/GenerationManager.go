package managers

import (
	"sort"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
	"gorm.io/gorm"
)

func GenerateNBAPlayoffSeriesRecords() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.SeasonID))
	nbaStandings := GetNBAStandingsBySeasonID(seasonID)
	// Separate standings by conference
	nbaStandingsMap := make(map[uint][]structs.NBAStandings)
	nbaTeams := GetAllActiveNBATeams()
	nbaTeamMap := MakeNBATeamMap(nbaTeams)
	playInMatches := []structs.NBAMatch{}
	nbaSeries := []structs.NBASeries{}
	latestSeriesID := GetLatestNBASeriesID()
	latestNBAMatchID := GetLatestNBAMatchID()

	for _, standing := range nbaStandings {
		if standing.TeamID == 0 {
			continue
		}
		nbaStandingsMap[standing.ConferenceID] = append(nbaStandingsMap[standing.ConferenceID], standing)
	}

	// There are 10 different conferences since there's also ISL included. Just sort each conference by total wins descending and then assign seeds based on that.
	for conferenceID := range nbaStandingsMap {
		sort.Slice(nbaStandingsMap[conferenceID], func(i, j int) bool {
			return nbaStandingsMap[conferenceID][i].TotalWins > nbaStandingsMap[conferenceID][j].TotalWins
		})
	}

	// NBA is conferences 1 and 2, this is where we setup the bubble matches for the 7/8 and 9/10 seeds
	/*
		Generate records for NBA playoff series, or do it via an automated fashion.
		Four games for 18A: Seeded between the 7 & 8 seed for both the East & West Conference and 9 & 10 seed for both the East & West Conference.
		The winner of the 7/8 seed game becomes the 7th seeded team in the playoffs.
		The winner of the 9/10 seed game plays the loser of the 7/8 seed game for the #8th seed spot.
		These two rounds take place on 18/A and 18/B, with the playoffs starting on 19A
	*/
	nbaEastStandings := nbaStandingsMap[1]
	nbaWestStandings := nbaStandingsMap[2]

	nbaEastSeries := generateSeriesRecordsForNBAPlayoffs(ts, latestSeriesID, latestSeriesID+14, nbaEastStandings, "Eastern", nbaTeamMap)
	nbaWestSeries := generateSeriesRecordsForNBAPlayoffs(ts, latestSeriesID+7, latestSeriesID+14, nbaWestStandings, "Western", nbaTeamMap)
	eastSeed8PlayinID := nbaEastSeries[0].ID
	westSeed8PlayinID := nbaWestSeries[0].ID
	eastSeed7PlayinID := nbaEastSeries[1].ID
	westSeed7PlayinID := nbaWestSeries[1].ID

	// Generate Play-In Matches for Week 18A and 18B
	eastPlayInMatches := generatePlayInTournament(ts, latestNBAMatchID+1, eastSeed8PlayinID, eastSeed7PlayinID, nbaEastStandings, "Eastern")
	westPlayInMatches := generatePlayInTournament(ts, latestNBAMatchID+3, westSeed8PlayinID, westSeed7PlayinID, nbaWestStandings, "Western")

	playInMatches = append(playInMatches, eastPlayInMatches...)
	playInMatches = append(playInMatches, westPlayInMatches...)
	nbaSeries = append(nbaSeries, nbaEastSeries...)
	nbaSeries = append(nbaSeries, nbaWestSeries...)

	// Save play-in matches to database
	for _, match := range playInMatches {
		db.Create(&match)
	}

	// Save playoff series to database
	for _, series := range nbaSeries {
		db.Create(&series)
	}

	intSeriesID := latestSeriesID + 14
	// Get top 2 teams from each international conference and create a playoff series record for them as well, with the winner of each series playing each other in an international final series
	internationalStandings := []structs.NBAStandings{}

	for cID, conference := range nbaStandingsMap {
		// Skip over NBA
		if cID == 1 || cID == 2 {
			continue
		}
		// Sort each conference by total wins
		sort.Slice(conference, func(i, j int) bool {
			return conference[i].TotalWins > conference[j].TotalWins
		})
		// Take top 2 teams from each conference and add to international standings
		if len(conference) > 1 {
			internationalStandings = append(internationalStandings, conference[0], conference[1])
		}
	}

	// Sort international standings by total wins, we should have about 16 teams in this list since there are 8 international conferences
	sort.Slice(internationalStandings, func(i, j int) bool {
		return internationalStandings[i].TotalWins > internationalStandings[j].TotalWins
	})

	internationalSeries := generateSeriesRecordsForISLPlayoffs(ts, intSeriesID, intSeriesID+14, internationalStandings, "ISL Playoffs", nbaTeamMap)

	for _, series := range internationalSeries {
		db.Create(&series)
	}
}

func generateSeriesRecordsForISLPlayoffs(ts structs.Timestamp, latestID, finalsID uint, standings []structs.NBAStandings, conferenceName string, nbaTeamMap map[uint]structs.NBATeam) []structs.NBASeries {
	series := []structs.NBASeries{}

	// Get top 16 seeds from standings
	seed1 := standings[0]   // Index 0 = 1st place
	seed2 := standings[1]   // Index 1 = 2nd place
	seed3 := standings[2]   // Index 2 = 3rd place
	seed4 := standings[3]   // Index 3 = 4th place
	seed5 := standings[4]   // Index 4 = 5th place
	seed6 := standings[5]   // Index 5 = 6th place
	seed7 := standings[6]   // Index 6 = 7th place
	seed8 := standings[7]   // Index 7 = 8th place
	seed9 := standings[8]   // Index 8 = 9th place
	seed10 := standings[9]  // Index 9 = 10th place
	seed11 := standings[10] // Index 10 = 11th place
	seed12 := standings[11] // Index 11 = 12th place
	seed13 := standings[12] // Index 12 = 13th place
	seed14 := standings[13] // Index 13 = 14th place
	seed15 := standings[14] // Index 14 = 15th place
	seed16 := standings[15] // Index 15 = 16th place

	// Playoff series 1v8 (winner of play-in), 2v7 (winner of play-in), 3v6, 4v5
	seed1Team := nbaTeamMap[seed1.TeamID]
	seed1Coach := getCoachName(seed1Team)
	seed2Team := nbaTeamMap[seed2.TeamID]
	seed2Coach := getCoachName(seed2Team)
	seed3Team := nbaTeamMap[seed3.TeamID]
	seed3Coach := getCoachName(seed3Team)
	seed4Team := nbaTeamMap[seed4.TeamID]
	seed4Coach := getCoachName(seed4Team)
	seed5Team := nbaTeamMap[seed5.TeamID]
	seed5Coach := getCoachName(seed5Team)
	seed6Team := nbaTeamMap[seed6.TeamID]
	seed6Coach := getCoachName(seed6Team)
	seed7Team := nbaTeamMap[seed7.TeamID]
	seed7Coach := getCoachName(seed7Team)
	seed8Team := nbaTeamMap[seed8.TeamID]
	seed8Coach := getCoachName(seed8Team)
	seed9Team := nbaTeamMap[seed9.TeamID]
	seed9Coach := getCoachName(seed9Team)
	seed10Team := nbaTeamMap[seed10.TeamID]
	seed10Coach := getCoachName(seed10Team)
	seed11Team := nbaTeamMap[seed11.TeamID]
	seed11Coach := getCoachName(seed11Team)
	seed12Team := nbaTeamMap[seed12.TeamID]
	seed12Coach := getCoachName(seed12Team)
	seed13Team := nbaTeamMap[seed13.TeamID]
	seed13Coach := getCoachName(seed13Team)
	seed14Team := nbaTeamMap[seed14.TeamID]
	seed14Coach := getCoachName(seed14Team)
	seed15Team := nbaTeamMap[seed15.TeamID]
	seed15Coach := getCoachName(seed15Team)
	seed16Team := nbaTeamMap[seed16.TeamID]
	seed16Coach := getCoachName(seed16Team)

	currentID := latestID + 1
	confSemiFinalsID1 := currentID + 8
	confSemiFinalsID2 := currentID + 9
	confSemiFinalsID3 := currentID + 10
	confSemiFinalsID4 := currentID + 11
	confFinalsID1 := currentID + 12
	confFinalsID2 := currentID + 13

	series1 := structs.NBASeries{
		Model: gorm.Model{
			ID: currentID,
		},
		SeriesName:    conferenceName + " Conference First Round - 1st Seed",
		SeasonID:      ts.SeasonID,
		HomeTeamID:    seed1.TeamID,
		HomeTeam:      seed1.TeamName,
		HomeTeamCoach: seed1Coach,
		HomeTeamRank:  1,
		AwayTeamRank:  16,
		AwayTeamCoach: seed16Coach,
		AwayTeamID:    seed16.TeamID,
		AwayTeam:      seed16.TeamName,
		GameCount:     1,
		NextSeriesID:  confSemiFinalsID1,
		IsPlayoffGame: true,
		NextSeriesHOA: "H",
	}
	series2 := structs.NBASeries{
		Model: gorm.Model{
			ID: currentID + 1,
		},
		SeriesName:    conferenceName + " Conference First Round - 2nd Seed",
		SeasonID:      ts.SeasonID,
		HomeTeamID:    seed2.TeamID,
		HomeTeam:      seed2.TeamName,
		HomeTeamCoach: seed2Coach,
		HomeTeamRank:  2,
		AwayTeamRank:  15,
		AwayTeamCoach: seed15Coach,
		AwayTeamID:    seed15.TeamID,
		AwayTeam:      seed15.TeamName,
		GameCount:     1,
		NextSeriesID:  confSemiFinalsID2,
		IsPlayoffGame: true,
		NextSeriesHOA: "H",
	}

	series3 := structs.NBASeries{
		Model: gorm.Model{
			ID: currentID + 2,
		},
		SeriesName:    conferenceName + " Conference First Round - 3rd Seed",
		SeasonID:      ts.SeasonID,
		HomeTeamID:    seed3.TeamID,
		HomeTeam:      seed3.TeamName,
		HomeTeamCoach: seed3Coach,
		HomeTeamRank:  3,
		AwayTeamRank:  14,
		AwayTeamID:    seed14.TeamID,
		AwayTeam:      seed14.TeamName,
		AwayTeamCoach: seed14Coach,
		GameCount:     1,
		NextSeriesID:  confSemiFinalsID3,
		IsPlayoffGame: true,
		NextSeriesHOA: "H",
	}

	series4 := structs.NBASeries{
		Model: gorm.Model{
			ID: currentID + 3,
		},
		SeriesName:    conferenceName + " Conference First Round - 4th Seed",
		SeasonID:      ts.SeasonID,
		HomeTeamID:    seed4.TeamID,
		HomeTeam:      seed4.TeamName,
		HomeTeamCoach: seed4Coach,
		HomeTeamRank:  4,
		AwayTeamRank:  13,
		AwayTeamID:    seed13.TeamID,
		AwayTeam:      seed13.TeamName,
		AwayTeamCoach: seed13Coach,
		GameCount:     1,
		NextSeriesID:  confSemiFinalsID4,
		IsPlayoffGame: true,
		NextSeriesHOA: "H",
	}

	series5 := structs.NBASeries{
		Model: gorm.Model{
			ID: currentID + 4,
		},
		SeriesName:    conferenceName + " Conference First Round - 5th Seed",
		SeasonID:      ts.SeasonID,
		HomeTeamID:    seed5.TeamID,
		HomeTeam:      seed5.TeamName,
		HomeTeamCoach: seed5Coach,
		HomeTeamRank:  5,
		AwayTeamRank:  12,
		AwayTeamID:    seed12.TeamID,
		AwayTeam:      seed12.TeamName,
		AwayTeamCoach: seed12Coach,
		GameCount:     1,
		NextSeriesID:  confSemiFinalsID4,
		IsPlayoffGame: true,
		NextSeriesHOA: "A",
	}

	series6 := structs.NBASeries{
		Model: gorm.Model{
			ID: currentID + 5,
		},
		SeriesName:    conferenceName + " Conference First Round - 6th Seed",
		SeasonID:      ts.SeasonID,
		HomeTeamID:    seed6.TeamID,
		HomeTeam:      seed6.TeamName,
		HomeTeamCoach: seed6Coach,
		HomeTeamRank:  6,
		AwayTeamRank:  11,
		AwayTeamID:    seed11.TeamID,
		AwayTeam:      seed11.TeamName,
		AwayTeamCoach: seed11Coach,
		GameCount:     1,
		NextSeriesID:  confSemiFinalsID3,
		IsPlayoffGame: true,
		NextSeriesHOA: "A",
	}

	series7 := structs.NBASeries{
		Model: gorm.Model{
			ID: currentID + 6,
		},
		SeriesName:    conferenceName + " Conference First Round - 7th Seed",
		SeasonID:      ts.SeasonID,
		HomeTeamID:    seed7.TeamID,
		HomeTeam:      seed7.TeamName,
		HomeTeamCoach: seed7Coach,
		HomeTeamRank:  7,
		AwayTeamRank:  10,
		AwayTeamID:    seed10.TeamID,
		AwayTeam:      seed10.TeamName,
		AwayTeamCoach: seed10Coach,
		GameCount:     1,
		NextSeriesID:  confSemiFinalsID2,
		IsPlayoffGame: true,
		NextSeriesHOA: "A",
	}

	series8 := structs.NBASeries{
		Model: gorm.Model{
			ID: currentID + 7,
		},
		SeriesName:    conferenceName + " Conference First Round - 8th Seed",
		SeasonID:      ts.SeasonID,
		HomeTeamID:    seed8.TeamID,
		HomeTeam:      seed8.TeamName,
		HomeTeamCoach: seed8Coach,
		HomeTeamRank:  8,
		AwayTeamRank:  9,
		AwayTeamID:    seed9.TeamID,
		AwayTeam:      seed9.TeamName,
		AwayTeamCoach: seed9Coach,
		GameCount:     1,
		NextSeriesID:  confSemiFinalsID1,
		IsPlayoffGame: true,
		NextSeriesHOA: "A",
	}

	series = append(series, series1, series2, series3, series4, series5, series6, series7, series8)

	// Make Conference Semifinal Games
	confSFSeries1 := structs.NBASeries{
		Model: gorm.Model{
			ID: confSemiFinalsID1,
		},
		SeriesName:    conferenceName + " Conference Semifinal",
		SeasonID:      ts.SeasonID,
		GameCount:     1,
		NextSeriesID:  confFinalsID1,
		IsPlayoffGame: true,
		NextSeriesHOA: "H",
	}
	confSFSeries2 := structs.NBASeries{
		Model: gorm.Model{
			ID: confSemiFinalsID2,
		},
		SeriesName:    conferenceName + " Conference Semifinal",
		SeasonID:      ts.SeasonID,
		GameCount:     1,
		NextSeriesID:  confFinalsID2,
		IsPlayoffGame: true,
		NextSeriesHOA: "H",
	}
	confSFSeries3 := structs.NBASeries{
		Model: gorm.Model{
			ID: confSemiFinalsID3,
		},
		SeriesName:    conferenceName + " Conference Semifinal",
		SeasonID:      ts.SeasonID,
		GameCount:     1,
		NextSeriesID:  confFinalsID2,
		IsPlayoffGame: true,
		NextSeriesHOA: "A",
	}
	confSFSeries4 := structs.NBASeries{
		Model: gorm.Model{
			ID: confSemiFinalsID4,
		},
		SeriesName:    conferenceName + " Conference Semifinal",
		SeasonID:      ts.SeasonID,
		GameCount:     1,
		NextSeriesID:  confFinalsID1,
		IsPlayoffGame: true,
		NextSeriesHOA: "A",
	}
	// Make Conference Final Games
	confFinals := structs.NBASeries{
		Model: gorm.Model{
			ID: confFinalsID1,
		},
		SeriesName:    conferenceName + " Conference Final",
		SeasonID:      ts.SeasonID,
		GameCount:     1,
		NextSeriesID:  finalsID,
		IsPlayoffGame: true,
		NextSeriesHOA: "H",
	}
	confFinals2 := structs.NBASeries{
		Model: gorm.Model{
			ID: confFinalsID2,
		},
		SeriesName:    conferenceName + " Conference Final",
		SeasonID:      ts.SeasonID,
		GameCount:     1,
		NextSeriesID:  finalsID,
		IsPlayoffGame: true,
		NextSeriesHOA: "H",
	}

	// ISL Finals
	theFinalsSeries := structs.NBASeries{
		Model: gorm.Model{
			ID: finalsID,
		},
		SeriesName:    conferenceName + " Conference Final",
		SeasonID:      ts.SeasonID,
		GameCount:     1,
		NextSeriesID:  0,
		IsPlayoffGame: true,
		NextSeriesHOA: "",
	}
	//

	series = append(series, confSFSeries1, confSFSeries2, confSFSeries3, confSFSeries4, confFinals, confFinals2, theFinalsSeries)

	return series
}

func generateSeriesRecordsForNBAPlayoffs(ts structs.Timestamp, latestID, finalsID uint, standings []structs.NBAStandings, conferenceName string, nbaTeamMap map[uint]structs.NBATeam) []structs.NBASeries {
	series := []structs.NBASeries{}

	// Get top 6 seeds from standings
	seed1 := standings[0] // Index 0 = 1st place
	seed2 := standings[1] // Index 1 = 2nd place
	seed3 := standings[2] // Index 2 = 3rd place
	seed4 := standings[3] // Index 3 = 4th place
	seed5 := standings[4] // Index 4 = 5th place
	seed6 := standings[5] // Index 5 = 6th place

	// Playoff series 1v8 (winner of play-in), 2v7 (winner of play-in), 3v6, 4v5
	seed1Team := nbaTeamMap[seed1.TeamID]
	seed1Coach := getCoachName(seed1Team)
	seed2Team := nbaTeamMap[seed2.TeamID]
	seed2Coach := getCoachName(seed2Team)
	seed3Team := nbaTeamMap[seed3.TeamID]
	seed3Coach := getCoachName(seed3Team)
	seed4Team := nbaTeamMap[seed4.TeamID]
	seed4Coach := getCoachName(seed4Team)
	seed5Team := nbaTeamMap[seed5.TeamID]
	seed5Coach := getCoachName(seed5Team)
	seed6Team := nbaTeamMap[seed6.TeamID]
	seed6Coach := getCoachName(seed6Team)

	currentID := latestID + 1

	series1 := structs.NBASeries{
		Model: gorm.Model{
			ID: currentID,
		},
		SeriesName:    conferenceName + " Conference First Round - 1st Seed",
		SeasonID:      ts.SeasonID,
		HomeTeamID:    seed1.TeamID,
		HomeTeam:      seed1.TeamName,
		HomeTeamCoach: seed1Coach,
		HomeTeamRank:  1,
		AwayTeamRank:  8,
		GameCount:     1,
		NextSeriesID:  currentID + 4,
		IsPlayoffGame: true,
		NextSeriesHOA: "H",
	}
	series2 := structs.NBASeries{
		Model: gorm.Model{
			ID: currentID + 1,
		},
		SeriesName:    conferenceName + " Conference First Round - 2nd Seed",
		SeasonID:      ts.SeasonID,
		HomeTeamID:    seed2.TeamID,
		HomeTeam:      seed2.TeamName,
		HomeTeamCoach: seed2Coach,
		HomeTeamRank:  2,
		AwayTeamRank:  7,
		GameCount:     1,
		NextSeriesID:  currentID + 5,
		IsPlayoffGame: true,
		NextSeriesHOA: "H",
	}

	series3 := structs.NBASeries{
		Model: gorm.Model{
			ID: currentID + 2,
		},
		SeriesName:    conferenceName + " Conference First Round - 3rd Seed",
		SeasonID:      ts.SeasonID,
		HomeTeamID:    seed3.TeamID,
		HomeTeam:      seed3.TeamName,
		HomeTeamCoach: seed3Coach,
		HomeTeamRank:  3,
		AwayTeamRank:  6,
		AwayTeamID:    seed6.TeamID,
		AwayTeam:      seed6.TeamName,
		AwayTeamCoach: seed6Coach,
		GameCount:     1,
		NextSeriesID:  currentID + 5,
		IsPlayoffGame: true,
		NextSeriesHOA: "A",
	}

	series4 := structs.NBASeries{
		Model: gorm.Model{
			ID: currentID + 3,
		},
		SeriesName:    conferenceName + " Conference First Round - 4th Seed",
		SeasonID:      ts.SeasonID,
		HomeTeamID:    seed4.TeamID,
		HomeTeam:      seed4.TeamName,
		HomeTeamCoach: seed4Coach,
		HomeTeamRank:  4,
		AwayTeamRank:  5,
		AwayTeamID:    seed5.TeamID,
		AwayTeam:      seed5.TeamName,
		AwayTeamCoach: seed5Coach,
		GameCount:     1,
		NextSeriesID:  currentID + 4,
		IsPlayoffGame: true,
		NextSeriesHOA: "A",
	}

	series = append(series, series1, series2, series3, series4)

	// Make Conference Semifinal Games
	confSFSeries1 := structs.NBASeries{
		Model: gorm.Model{
			ID: currentID + 4,
		},
		SeriesName:    conferenceName + " Conference Semifinal",
		SeasonID:      ts.SeasonID,
		HomeTeamID:    seed1.TeamID,
		HomeTeam:      seed1.TeamName,
		HomeTeamCoach: seed1Coach,
		HomeTeamRank:  1,
		AwayTeamRank:  8,
		GameCount:     1,
		NextSeriesID:  currentID + 6,
		IsPlayoffGame: true,
		NextSeriesHOA: "H",
	}
	confSFSeries2 := structs.NBASeries{
		Model: gorm.Model{
			ID: currentID + 5,
		},
		SeriesName:    conferenceName + " Conference Semifinal",
		SeasonID:      ts.SeasonID,
		HomeTeamID:    seed2.TeamID,
		HomeTeam:      seed2.TeamName,
		HomeTeamCoach: seed2Coach,
		HomeTeamRank:  2,
		AwayTeamRank:  7,
		GameCount:     1,
		NextSeriesID:  currentID + 6,
		IsPlayoffGame: true,
		NextSeriesHOA: "H",
	}
	// Make Conference Final Games
	confFinals := structs.NBASeries{
		Model: gorm.Model{
			ID: currentID + 6,
		},
		SeriesName:    conferenceName + " Conference Final",
		SeasonID:      ts.SeasonID,
		HomeTeamID:    seed2.TeamID,
		HomeTeam:      seed2.TeamName,
		HomeTeamCoach: seed2Coach,
		HomeTeamRank:  2,
		AwayTeamRank:  7,
		GameCount:     1,
		NextSeriesID:  finalsID,
		IsPlayoffGame: true,
		NextSeriesHOA: "H",
	}

	series = append(series, confSFSeries1, confSFSeries2, confFinals)

	return series
}

// generatePlayInTournament creates the play-in matches and series for a conference
func generatePlayInTournament(ts structs.Timestamp, nextMatchID, series1ID, series2ID uint, standings []structs.NBAStandings, conferenceName string) []structs.NBAMatch {
	matches := []structs.NBAMatch{}

	// Get seeds 7, 8, 9, 10 from standings
	seed7 := standings[6]  // Index 6 = 7th place
	seed8 := standings[7]  // Index 7 = 8th place
	seed9 := standings[8]  // Index 8 = 9th place
	seed10 := standings[9] // Index 9 = 10th place

	// Week 18A - First Play-In Round
	// Match 1: 7 seed vs 8 seed (higher seed at home)
	match7v8 := GenerateNBAMatch(ts, seed7, seed8, "PlayIn-7v8", 18)
	match7v8.ID = nextMatchID
	match7v8.MatchName = conferenceName + " Play-In: 7 vs 8"
	match7v8.MatchOfWeek = "A"
	match7v8.IsPlayInGame = true
	match7v8.NextSeriesID = series2ID
	match7v8.NextSeriesHOA = "A"          // 7/8 winner is always home team in 8th seed series
	match7v8.NextGameID = nextMatchID + 2 // This will be the 18B match between 7/8 loser and 9/10 winner
	matches = append(matches, match7v8)

	// Match 2: 9 seed vs 10 seed (higher seed at home)
	match9v10 := GenerateNBAMatch(ts, seed9, seed10, "PlayIn-9v10", 18)
	match9v10.ID = nextMatchID + 1
	match9v10.MatchName = conferenceName + " Play-In: 9 vs 10"
	match9v10.MatchOfWeek = "A"
	match9v10.IsPlayInGame = true
	match9v10.NextGameID = nextMatchID + 2 // This will be the 18B match between 7/8 loser and 9/10 winner
	matches = append(matches, match9v10)

	// Week 18B - Second Play-In Round
	// This match will be between loser of 7v8 and winner of 9v10
	// We create a placeholder match with NextGameID pointing to the 8th seed series
	matchFinal := structs.NBAMatch{
		Model: gorm.Model{
			ID: nextMatchID + 2,
		},
		WeekID:        ts.NBAWeekID, // Next week
		Week:          uint(ts.NBAWeek),
		SeasonID:      ts.SeasonID,
		MatchName:     conferenceName + " Play-In Final: 8th Seed",
		MatchOfWeek:   "B",
		IsPlayoffGame: false, // This is still a play-in game
		IsPlayInGame:  true,
		HomeTeamRank:  0, // To be determined from 7v8 loser
		AwayTeamRank:  0, // To be determined from 9v10 winner
		NextSeriesID:  series1ID,
		NextSeriesHOA: "A", // 8th seed series is always away team
	}
	matches = append(matches, matchFinal)

	return matches
}

func GenerateNBAMatch(ts structs.Timestamp, nbaTeamA structs.NBAStandings, nbaTeamB structs.NBAStandings, matchType string, week int) structs.NBAMatch {
	// Get full team data
	teamA := GetNBATeamByTeamID(strconv.Itoa(int(nbaTeamA.TeamID)))
	teamB := GetNBATeamByTeamID(strconv.Itoa(int(nbaTeamB.TeamID)))

	// Higher seed gets home court (teamA is higher seed by convention)
	coachA := teamA.NBACoachName
	if coachA == "" {
		coachA = teamA.NBAOwnerName
		if coachA == "" {
			coachA = "AI"
		}
	}

	coachB := teamB.NBACoachName
	if coachB == "" {
		coachB = teamB.NBAOwnerName
		if coachB == "" {
			coachB = "AI"
		}
	}

	return structs.NBAMatch{
		WeekID:        ts.NBAWeekID + uint(week-int(ts.NBAWeek)),
		Week:          uint(week),
		SeasonID:      ts.SeasonID,
		HomeTeamID:    teamA.ID,
		HomeTeam:      teamA.Team,
		HomeTeamCoach: coachA,
		HomeTeamRank:  uint(getSeedFromStanding(nbaTeamA)),
		AwayTeamID:    teamB.ID,
		AwayTeam:      teamB.Team,
		AwayTeamCoach: coachB,
		AwayTeamRank:  uint(getSeedFromStanding(nbaTeamB)),
		Arena:         teamA.Arena,
		City:          teamA.City,
		State:         teamA.State,
		Country:       teamA.Country,
		IsConference:  teamA.ConferenceID == teamB.ConferenceID,
		IsPlayoffGame: false, // Play-in games are not playoff games yet
		GameComplete:  false,
	}
}

// Helper function to determine seed from standing position
func getSeedFromStanding(standing structs.NBAStandings) int {
	// In a real implementation, you'd calculate this based on conference position
	// For now, return a placeholder
	return 0
}

func getCoachName(team structs.NBATeam) string {
	coachName := team.NBACoachName
	if coachName == "" {
		coachName = team.NBAOwnerName
		if coachName == "" {
			coachName = "AI"
		}
	}
	return coachName
}
