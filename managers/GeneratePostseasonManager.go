package managers

import (
	"sort"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
	"gorm.io/gorm"
)

// Enums for the amount of teams within a conference tournament, this will be used for determining how many teams to generate for each conference tournament and how to seed them within the tournament
const (
	CBB18Team = 1
	CBB16Team = 2
	CBB15Team = 3
	CBB14Team = 4
	CBB13Team = 5
	CBB12Team = 6
	CBB11Team = 7
	CBB10Team = 8
	CBB9Team  = 9
	CBB8Team  = 10
	BigTen    = 11
	ACC       = 12
	Big12     = 13
	SEC       = 14
	BigEast   = 15
	IvyLeague = 16
	WCC       = 17
)

// Pro Post Season Generation
func GenerateNBAPlayoffSeriesRecords() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.SeasonID))
	nbaStandings := GetNBAStandingsBySeasonID(seasonID)
	// Separate standings by conference
	nbaStandingsMap := make(map[uint][]structs.NBAStandings)
	nbaTeams := GetAllActiveNBATeams()
	nbaTeamMap := MakeNBATeamMap(nbaTeams)
	nbaSeries := []structs.NBASeries{}
	latestSeriesID := GetLatestNBASeriesID()

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

	// Filter the east & west standings using the below rules:
	// Need the division winners of each conference, followed by the top four remaining qualified teams
	filteredEastStandings := filterNBAStandingsForPlayoffs(nbaEastStandings)
	filteredWestStandings := filterNBAStandingsForPlayoffs(nbaWestStandings)

	nbaEastSeries := generateSeriesRecordsForNBAPlayoffs(ts, latestSeriesID, latestSeriesID+14, filteredEastStandings, "Eastern", nbaTeamMap)
	nbaWestSeries := generateSeriesRecordsForNBAPlayoffs(ts, latestSeriesID+7, latestSeriesID+14, filteredWestStandings, "Western", nbaTeamMap)

	// Generate Play-In Matches for Week 18A and 18B
	nbaSeries = append(nbaSeries, nbaEastSeries...)
	nbaSeries = append(nbaSeries, nbaWestSeries...)

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

	// Save playoff series to database
	for _, series := range nbaSeries {
		db.Create(&series)
	}
	for _, series := range internationalSeries {
		db.Create(&series)
	}
}

func filterNBAStandingsForPlayoffs(standings []structs.NBAStandings) []structs.NBAStandings {
	filteredStandings := []structs.NBAStandings{}

	// Sort by division and then total wins, we want division winners to be at the top of the standings followed by the rest of the teams sorted by total wins regardless of division
	sort.Slice(standings, func(i, j int) bool {
		return standings[i].DivisionID > standings[j].DivisionID && standings[i].TotalWins > standings[j].TotalWins
	})

	// Get division winners
	divisionWinners := make(map[uint]bool)
	qualifiedTeams := make(map[uint]structs.NBAStandings)

	// Get division winners first
	for _, standing := range standings {
		if standing.TeamID == 0 {
			continue
		}
		if divisionWinners[standing.DivisionID] {
			continue
		}
		divisionWinners[standing.DivisionID] = true
		qualifiedTeams[standing.TeamID] = standing
		filteredStandings = append(filteredStandings, standing)
	}

	// Get next best teams until we have 8 total teams for the playoffs
	for _, standing := range standings {
		if standing.TeamID == 0 {
			continue
		}
		if len(filteredStandings) >= 8 {
			break
		}
		if _, ok := qualifiedTeams[standing.TeamID]; ok {
			continue
		}
		qualifiedTeams[standing.TeamID] = standing
		filteredStandings = append(filteredStandings, standing)
	}

	return filteredStandings
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
	seed1Coach := getProCoachName(seed1Team)
	seed2Team := nbaTeamMap[seed2.TeamID]
	seed2Coach := getProCoachName(seed2Team)
	seed3Team := nbaTeamMap[seed3.TeamID]
	seed3Coach := getProCoachName(seed3Team)
	seed4Team := nbaTeamMap[seed4.TeamID]
	seed4Coach := getProCoachName(seed4Team)
	seed5Team := nbaTeamMap[seed5.TeamID]
	seed5Coach := getProCoachName(seed5Team)
	seed6Team := nbaTeamMap[seed6.TeamID]
	seed6Coach := getProCoachName(seed6Team)
	seed7Team := nbaTeamMap[seed7.TeamID]
	seed7Coach := getProCoachName(seed7Team)
	seed8Team := nbaTeamMap[seed8.TeamID]
	seed8Coach := getProCoachName(seed8Team)
	seed9Team := nbaTeamMap[seed9.TeamID]
	seed9Coach := getProCoachName(seed9Team)
	seed10Team := nbaTeamMap[seed10.TeamID]
	seed10Coach := getProCoachName(seed10Team)
	seed11Team := nbaTeamMap[seed11.TeamID]
	seed11Coach := getProCoachName(seed11Team)
	seed12Team := nbaTeamMap[seed12.TeamID]
	seed12Coach := getProCoachName(seed12Team)
	seed13Team := nbaTeamMap[seed13.TeamID]
	seed13Coach := getProCoachName(seed13Team)
	seed14Team := nbaTeamMap[seed14.TeamID]
	seed14Coach := getProCoachName(seed14Team)
	seed15Team := nbaTeamMap[seed15.TeamID]
	seed15Coach := getProCoachName(seed15Team)
	seed16Team := nbaTeamMap[seed16.TeamID]
	seed16Coach := getProCoachName(seed16Team)

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
	seed7 := standings[6] // Index 6 = 7th place
	seed8 := standings[7] // Index 7 = 8th place

	// Playoff series 1v8 (winner of play-in), 2v7 (winner of play-in), 3v6, 4v5
	seed1Team := nbaTeamMap[seed1.TeamID]
	seed1Coach := getProCoachName(seed1Team)
	seed2Team := nbaTeamMap[seed2.TeamID]
	seed2Coach := getProCoachName(seed2Team)
	seed3Team := nbaTeamMap[seed3.TeamID]
	seed3Coach := getProCoachName(seed3Team)
	seed4Team := nbaTeamMap[seed4.TeamID]
	seed4Coach := getProCoachName(seed4Team)
	seed5Team := nbaTeamMap[seed5.TeamID]
	seed5Coach := getProCoachName(seed5Team)
	seed6Team := nbaTeamMap[seed6.TeamID]
	seed6Coach := getProCoachName(seed6Team)
	seed7Team := nbaTeamMap[seed7.TeamID]
	seed7Coach := getProCoachName(seed7Team)
	seed8Team := nbaTeamMap[seed8.TeamID]
	seed8Coach := getProCoachName(seed8Team)

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
		AwayTeamID:    seed8.TeamID,
		AwayTeam:      seed8.TeamName,
		AwayTeamCoach: seed8Coach,
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
		AwayTeamID:    seed7.TeamID,
		AwayTeam:      seed7.TeamName,
		AwayTeamCoach: seed7Coach,
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

func getProCoachName(team structs.NBATeam) string {
	coachName := team.NBACoachName
	if coachName == "" {
		coachName = team.NBAOwnerName
		if coachName == "" {
			coachName = "AI"
		}
	}
	return coachName
}

// College Conference Tournament Generation
func GenerateConferenceTournamentRecords() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.SeasonID))
	collegeStandings := GetAllConferenceStandingsBySeasonID(seasonID)
	latestMatchID := GetLatestCollegeMatchID()
	collegeConferenceMap := make(map[uint][]structs.CollegeStandings)
	teamMap := GetCollegeTeamMap()
	for _, standing := range collegeStandings {
		if standing.TeamID == 0 {
			continue
		}
		collegeConferenceMap[standing.ConferenceID] = append(collegeConferenceMap[standing.ConferenceID], standing)
	}

	// Note - there are about 31 college conferences
	// Sort each conference by conference wins descending, we will use this order to determine seeding for the conference tournaments
	for conferenceID := range collegeConferenceMap {
		sort.Slice(collegeConferenceMap[conferenceID], func(i, j int) bool {
			return collegeConferenceMap[conferenceID][i].ConferenceWins > collegeConferenceMap[conferenceID][j].ConferenceWins
		})
	}

	postSeasonMatches := []structs.Match{}

	// Go through each conference & build up the conference tournament matches based on the standings, we want to give byes to the top 4 teams in each conference and then have the 5-12 seeds play in the first round, with the winners of those matches playing the top 4 seeds in the quarterfinals, and then semifinals and finals after that. We will also need to assign a game number to each match, so we will need to get the latest game number from the database and then increment from there for each match we create.
	for _, conference := range collegeConferenceMap {
		if len(conference) < 8 {
			continue
		}
		tournamentBracketEnum := getTournamentBracketEnumByConferenceID(uint(len(conference)), conference[0].ConferenceID)
		if tournamentBracketEnum == 0 {
			continue
		}
		conferenceTournamentMatches := generateConferenceTournamentMatchesBase(ts, latestMatchID+1, conference, uint(tournamentBracketEnum), teamMap)
		lastRecord := conferenceTournamentMatches[len(conferenceTournamentMatches)-1]
		latestMatchID = lastRecord.ID
		postSeasonMatches = append(postSeasonMatches, conferenceTournamentMatches...)
	}

	for _, match := range postSeasonMatches {
		db.Create(&match)
	}
}

func getTournamentBracketEnumByConferenceID(conferenceTeamNum, conferenceID uint) int {
	// Use IRL format for the following conferences
	switch conferenceID {
	case 1:
		return ACC
	case 2:
		return BigTen
	case 3:
		return Big12
	case 5:
		return SEC
	case 11:
		return BigEast
	case 26:
		return IvyLeague
	case 13:
		return WCC

	}

	// Use generic format for the following conferences, which is based on the number of teams in the conference tournament, so for example if there are 12 teams in the conference tournament then we would use the CBB12Team enum to determine how to seed the teams and how many rounds there are in the tournament
	switch conferenceTeamNum {
	case 18:
		return CBB18Team
	case 16:
		return CBB16Team
	case 15:
		return CBB15Team
	case 14:
		return CBB14Team
	case 13:
		return CBB13Team
	case 12:
		return CBB12Team
	case 11:
		return CBB11Team
	case 10:
		return CBB10Team
	case 9:
		return CBB9Team
	case 8:
		return CBB8Team
	default:
		return 0
	}
}

func generateConferenceTournamentMatchesBase(ts structs.Timestamp, latestMatchID uint, standings []structs.CollegeStandings, conferenceEnum uint, teamMap map[uint]structs.Team) []structs.Match {
	switch conferenceEnum {
	case CBB18Team:
		return []structs.Match{}
	case CBB16Team:
		return []structs.Match{}
	case CBB15Team:
		return generateCBB15TeamTournamentMatches(ts, latestMatchID, standings, teamMap)
	case CBB14Team:
		return generateCBB14TeamTournamentMatches(ts, latestMatchID, standings, teamMap)
	case CBB13Team:
		return generateCBB13TeamTournamentMatches(ts, latestMatchID, standings, teamMap)
	case CBB12Team:
		return generateCBB12TeamTournamentMatches(ts, latestMatchID, standings, teamMap)
	case CBB11Team:
		return generateCBB11TeamTournamentMatches(ts, latestMatchID, standings, teamMap)
	case CBB10Team:
		return generateCBB10TeamTournamentMatches(ts, latestMatchID, standings, teamMap)
	case CBB9Team:
		return generateCBB9TeamTournamentMatches(ts, latestMatchID, standings, teamMap)
	case CBB8Team:
		return generateCBB8TeamTournamentMatches(ts, latestMatchID, standings, teamMap)
	case BigTen:
		return generateBigTenTournamentMatches(ts, latestMatchID, standings, teamMap)
	case ACC:
		return generateACCtournamentMatches(ts, latestMatchID, standings, teamMap)
	case Big12:
		return generateBig12TournamentMatches(ts, latestMatchID, standings, teamMap)
	case SEC:
		return generateSECTournamentMatches(ts, latestMatchID, standings, teamMap)
	case BigEast:
		return generateBigEastTournamentMatches(ts, latestMatchID, standings, teamMap)
	case IvyLeague:
		return generateIvyLeagueTournamentMatches(ts, latestMatchID, standings, teamMap)
	case WCC:
		return generateWCCTournamentMatches(ts, latestMatchID, standings, teamMap)
	default:
		return []structs.Match{}
	}
}

func generateCBB15TeamTournamentMatches(ts structs.Timestamp, latestMatchID uint, standings []structs.CollegeStandings, teamMap map[uint]structs.Team) []structs.Match {
	if len(standings) < 15 {
		return []structs.Match{}
	}

	conferenceName := standings[0].ConferenceName

	seed1 := standings[0]
	seed2 := standings[1]
	seed3 := standings[2]
	seed4 := standings[3]
	seed5 := standings[4]
	seed6 := standings[5]
	seed7 := standings[6]
	seed8 := standings[7]
	seed9 := standings[8]
	seed10 := standings[9]
	seed11 := standings[10]
	seed12 := standings[11]
	seed13 := standings[12]
	seed14 := standings[13]
	seed15 := standings[14]

	playIn1ID := latestMatchID
	playIn2ID := latestMatchID + 1
	playIn3ID := latestMatchID + 2
	playIn4ID := latestMatchID + 3
	playIn5ID := latestMatchID + 4
	playIn6ID := latestMatchID + 5
	playIn7ID := latestMatchID + 6
	qf1ID := latestMatchID + 7
	qf2ID := latestMatchID + 8
	qf3ID := latestMatchID + 9
	qf4ID := latestMatchID + 10
	sf1ID := latestMatchID + 11
	sf2ID := latestMatchID + 12
	finalID := latestMatchID + 13

	buildMatch := func(id uint, homeSeed, awaySeed structs.CollegeStandings, homeRank, awayRank uint, name, slot string, nextID uint, nextHOA string) structs.Match {
		homeTeam := teamMap[homeSeed.TeamID]
		awayTeam := teamMap[awaySeed.TeamID]
		homeCoach := homeTeam.Coach
		if homeCoach == "" {
			homeCoach = "AI"
		}
		awayCoach := awayTeam.Coach
		if awayCoach == "" {
			awayCoach = "AI"
		}

		return structs.Match{
			Model:                  gorm.Model{ID: id},
			MatchName:              name,
			WeekID:                 ts.CollegeWeekID,
			Week:                   uint(ts.CollegeWeek),
			SeasonID:               ts.SeasonID,
			HomeTeamID:             homeSeed.TeamID,
			HomeTeam:               homeTeam.Abbr,
			HomeTeamCoach:          homeCoach,
			HomeTeamRank:           homeRank,
			AwayTeamID:             awaySeed.TeamID,
			AwayTeam:               awayTeam.Abbr,
			AwayTeamCoach:          awayCoach,
			AwayTeamRank:           awayRank,
			MatchOfWeek:            slot,
			NextGameID:             nextID,
			NextGameHOA:            nextHOA,
			IsConference:           true,
			IsConferenceTournament: true,
			IsPlayoffGame:          false,
			IsNeutralSite:          false,
			Arena:                  homeTeam.Arena,
			City:                   homeTeam.City,
			State:                  homeTeam.State,
		}
	}

	matches := []structs.Match{}

	// 7 play-in games: 2v15, 3v14, 4v13, 5v12, 6v11, 7v10, 8v9
	playIn1 := buildMatch(playIn1ID, seed2, seed15, 2, 15, conferenceName+" Tournament First Round", "A", qf4ID, "A")
	playIn2 := buildMatch(playIn2ID, seed3, seed14, 3, 14, conferenceName+" Tournament First Round", "A", qf3ID, "A")
	playIn3 := buildMatch(playIn3ID, seed4, seed13, 4, 13, conferenceName+" Tournament First Round", "A", qf2ID, "A")
	playIn4 := buildMatch(playIn4ID, seed5, seed12, 5, 12, conferenceName+" Tournament First Round", "A", qf2ID, "H")
	playIn5 := buildMatch(playIn5ID, seed6, seed11, 6, 11, conferenceName+" Tournament First Round", "A", qf3ID, "H")
	playIn6 := buildMatch(playIn6ID, seed7, seed10, 7, 10, conferenceName+" Tournament First Round", "A", qf4ID, "H")
	playIn7 := buildMatch(playIn7ID, seed8, seed9, 8, 9, conferenceName+" Tournament First Round", "A", qf1ID, "A")

	// Quarterfinals: 1 vs 8/9 winner, 4/13 vs 5/12 winner, 3/14 vs 6/11 winner, 2/15 vs 7/10 winner
	qf1 := buildMatch(qf1ID, seed1, seed8, 1, 8, conferenceName+" Tournament Quarterfinals", "B", sf1ID, "H")
	qf1.AwayTeamID = 0
	qf1.AwayTeam = ""
	qf1.AwayTeamCoach = ""
	qf1.AwayTeamRank = 0

	qf2 := structs.Match{
		Model:                  gorm.Model{ID: qf2ID},
		MatchName:              conferenceName + " Tournament Quarterfinals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "B",
		NextGameID:             sf1ID,
		NextGameHOA:            "A",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	qf3 := structs.Match{
		Model:                  gorm.Model{ID: qf3ID},
		MatchName:              conferenceName + " Tournament Quarterfinals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "B",
		NextGameID:             sf2ID,
		NextGameHOA:            "H",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	qf4 := structs.Match{
		Model:                  gorm.Model{ID: qf4ID},
		MatchName:              conferenceName + " Tournament Quarterfinals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "B",
		NextGameID:             sf2ID,
		NextGameHOA:            "A",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	sf1 := structs.Match{
		Model:                  gorm.Model{ID: sf1ID},
		MatchName:              conferenceName + " Tournament Semifinals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "C",
		NextGameID:             finalID,
		NextGameHOA:            "H",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	sf2 := structs.Match{
		Model:                  gorm.Model{ID: sf2ID},
		MatchName:              conferenceName + " Tournament Semifinals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "C",
		NextGameID:             finalID,
		NextGameHOA:            "A",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	finalMatch := structs.Match{
		Model:                  gorm.Model{ID: finalID},
		MatchName:              conferenceName + " Tournament Finals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "D",
		NextGameID:             0,
		NextGameHOA:            "",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	matches = append(matches, playIn1, playIn2, playIn3, playIn4, playIn5, playIn6, playIn7, qf1, qf2, qf3, qf4, sf1, sf2, finalMatch)

	return matches
}

func generateCBB14TeamTournamentMatches(ts structs.Timestamp, latestMatchID uint, standings []structs.CollegeStandings, teamMap map[uint]structs.Team) []structs.Match {
	if len(standings) < 14 {
		return []structs.Match{}
	}

	conferenceName := standings[0].ConferenceName

	seed1 := standings[0]
	seed2 := standings[1]
	seed3 := standings[2]
	seed4 := standings[3]
	seed5 := standings[4]
	seed6 := standings[5]
	seed7 := standings[6]
	seed8 := standings[7]
	seed9 := standings[8]
	seed10 := standings[9]
	seed11 := standings[10]
	seed12 := standings[11]
	seed13 := standings[12]
	seed14 := standings[13]

	playIn1ID := latestMatchID
	playIn2ID := latestMatchID + 1
	playIn3ID := latestMatchID + 2
	playIn4ID := latestMatchID + 3
	playIn5ID := latestMatchID + 4
	playIn6ID := latestMatchID + 5
	qf1ID := latestMatchID + 6
	qf2ID := latestMatchID + 7
	qf3ID := latestMatchID + 8
	qf4ID := latestMatchID + 9
	sf1ID := latestMatchID + 10
	sf2ID := latestMatchID + 11
	finalID := latestMatchID + 12

	buildMatch := func(id uint, homeSeed, awaySeed structs.CollegeStandings, homeRank, awayRank uint, name, slot string, nextID uint, nextHOA string) structs.Match {
		homeTeam := teamMap[homeSeed.TeamID]
		awayTeam := teamMap[awaySeed.TeamID]
		homeCoach := homeTeam.Coach
		if homeCoach == "" {
			homeCoach = "AI"
		}
		awayCoach := awayTeam.Coach
		if awayCoach == "" {
			awayCoach = "AI"
		}

		return structs.Match{
			Model:                  gorm.Model{ID: id},
			MatchName:              name,
			WeekID:                 ts.CollegeWeekID,
			Week:                   uint(ts.CollegeWeek),
			SeasonID:               ts.SeasonID,
			HomeTeamID:             homeSeed.TeamID,
			HomeTeam:               homeTeam.Abbr,
			HomeTeamCoach:          homeCoach,
			HomeTeamRank:           homeRank,
			AwayTeamID:             awaySeed.TeamID,
			AwayTeam:               awayTeam.Abbr,
			AwayTeamCoach:          awayCoach,
			AwayTeamRank:           awayRank,
			MatchOfWeek:            slot,
			NextGameID:             nextID,
			NextGameHOA:            nextHOA,
			IsConference:           true,
			IsConferenceTournament: true,
			IsPlayoffGame:          false,
			IsNeutralSite:          false,
			Arena:                  homeTeam.Arena,
			City:                   homeTeam.City,
			State:                  homeTeam.State,
		}
	}

	matches := []structs.Match{}

	// 6 play-in games: 3v14, 4v13, 5v12, 6v11, 7v10, 8v9
	playIn1 := buildMatch(playIn1ID, seed3, seed14, 3, 14, conferenceName+" Tournament First Round", "A", qf3ID, "A")
	playIn2 := buildMatch(playIn2ID, seed4, seed13, 4, 13, conferenceName+" Tournament First Round", "A", qf2ID, "A")
	playIn3 := buildMatch(playIn3ID, seed5, seed12, 5, 12, conferenceName+" Tournament First Round", "A", qf2ID, "H")
	playIn4 := buildMatch(playIn4ID, seed6, seed11, 6, 11, conferenceName+" Tournament First Round", "A", qf3ID, "H")
	playIn5 := buildMatch(playIn5ID, seed7, seed10, 7, 10, conferenceName+" Tournament First Round", "A", qf4ID, "H")
	playIn6 := buildMatch(playIn6ID, seed8, seed9, 8, 9, conferenceName+" Tournament First Round", "A", qf1ID, "A")

	// Quarterfinals: 1 vs 8/9 winner, 4/13 vs 5/12 winner, 3/14 vs 6/11 winner, 2 vs 7/10 winner
	qf1 := buildMatch(qf1ID, seed1, seed8, 1, 8, conferenceName+" Tournament Quarterfinals", "B", sf1ID, "H")
	qf1.AwayTeamID = 0
	qf1.AwayTeam = ""
	qf1.AwayTeamCoach = ""
	qf1.AwayTeamRank = 0

	qf2 := structs.Match{
		Model:                  gorm.Model{ID: qf2ID},
		MatchName:              conferenceName + " Tournament Quarterfinals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "B",
		NextGameID:             sf1ID,
		NextGameHOA:            "A",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	qf3 := structs.Match{
		Model:                  gorm.Model{ID: qf3ID},
		MatchName:              conferenceName + " Tournament Quarterfinals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "B",
		NextGameID:             sf2ID,
		NextGameHOA:            "H",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	qf4 := buildMatch(qf4ID, seed2, seed7, 2, 7, conferenceName+" Tournament Quarterfinals", "B", sf2ID, "A")
	qf4.AwayTeamID = 0
	qf4.AwayTeam = ""
	qf4.AwayTeamCoach = ""
	qf4.AwayTeamRank = 0

	sf1 := structs.Match{
		Model:                  gorm.Model{ID: sf1ID},
		MatchName:              conferenceName + " Tournament Semifinals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "C",
		NextGameID:             finalID,
		NextGameHOA:            "H",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	sf2 := structs.Match{
		Model:                  gorm.Model{ID: sf2ID},
		MatchName:              conferenceName + " Tournament Semifinals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "C",
		NextGameID:             finalID,
		NextGameHOA:            "A",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	finalMatch := structs.Match{
		Model:                  gorm.Model{ID: finalID},
		MatchName:              conferenceName + " Tournament Finals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "D",
		NextGameID:             0,
		NextGameHOA:            "",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	matches = append(matches, playIn1, playIn2, playIn3, playIn4, playIn5, playIn6, qf1, qf2, qf3, qf4, sf1, sf2, finalMatch)

	return matches
}

func generateCBB13TeamTournamentMatches(ts structs.Timestamp, latestMatchID uint, standings []structs.CollegeStandings, teamMap map[uint]structs.Team) []structs.Match {
	if len(standings) < 13 {
		return []structs.Match{}
	}

	conferenceName := standings[0].ConferenceName

	seed1 := standings[0]
	seed2 := standings[1]
	seed3 := standings[2]
	seed4 := standings[3]
	seed5 := standings[4]
	seed6 := standings[5]
	seed7 := standings[6]
	seed8 := standings[7]
	seed9 := standings[8]
	seed10 := standings[9]
	seed11 := standings[10]
	seed12 := standings[11]
	seed13 := standings[12]

	playIn1ID := latestMatchID
	playIn2ID := latestMatchID + 1
	playIn3ID := latestMatchID + 2
	playIn4ID := latestMatchID + 3
	playIn5ID := latestMatchID + 4
	qf1ID := latestMatchID + 5
	qf2ID := latestMatchID + 6
	qf3ID := latestMatchID + 7
	qf4ID := latestMatchID + 8
	sf1ID := latestMatchID + 9
	sf2ID := latestMatchID + 10
	finalID := latestMatchID + 11

	buildMatch := func(id uint, homeSeed, awaySeed structs.CollegeStandings, homeRank, awayRank uint, name, slot string, nextID uint, nextHOA string) structs.Match {
		homeTeam := teamMap[homeSeed.TeamID]
		awayTeam := teamMap[awaySeed.TeamID]
		homeCoach := homeTeam.Coach
		if homeCoach == "" {
			homeCoach = "AI"
		}
		awayCoach := awayTeam.Coach
		if awayCoach == "" {
			awayCoach = "AI"
		}

		return structs.Match{
			Model:                  gorm.Model{ID: id},
			MatchName:              name,
			WeekID:                 ts.CollegeWeekID,
			Week:                   uint(ts.CollegeWeek),
			SeasonID:               ts.SeasonID,
			HomeTeamID:             homeSeed.TeamID,
			HomeTeam:               homeTeam.Abbr,
			HomeTeamCoach:          homeCoach,
			HomeTeamRank:           homeRank,
			AwayTeamID:             awaySeed.TeamID,
			AwayTeam:               awayTeam.Abbr,
			AwayTeamCoach:          awayCoach,
			AwayTeamRank:           awayRank,
			MatchOfWeek:            slot,
			NextGameID:             nextID,
			NextGameHOA:            nextHOA,
			IsConference:           true,
			IsConferenceTournament: true,
			IsPlayoffGame:          false,
			IsNeutralSite:          false,
			Arena:                  homeTeam.Arena,
			City:                   homeTeam.City,
			State:                  homeTeam.State,
		}
	}

	matches := []structs.Match{}

	// 5 play-in games: 4v13, 5v12, 6v11, 7v10, 8v9
	playIn1 := buildMatch(playIn1ID, seed4, seed13, 4, 13, conferenceName+" Tournament First Round", "A", qf2ID, "A")
	playIn2 := buildMatch(playIn2ID, seed5, seed12, 5, 12, conferenceName+" Tournament First Round", "A", qf2ID, "H")
	playIn3 := buildMatch(playIn3ID, seed6, seed11, 6, 11, conferenceName+" Tournament First Round", "A", qf3ID, "H")
	playIn4 := buildMatch(playIn4ID, seed7, seed10, 7, 10, conferenceName+" Tournament First Round", "A", qf4ID, "H")
	playIn5 := buildMatch(playIn5ID, seed8, seed9, 8, 9, conferenceName+" Tournament First Round", "A", qf1ID, "A")

	// Quarterfinals: 1 vs 8/9 winner, 4/13 vs 5/12 winner, 3 vs 6/11 winner, 2 vs 7/10 winner
	qf1 := buildMatch(qf1ID, seed1, seed8, 1, 8, conferenceName+" Tournament Quarterfinals", "B", sf1ID, "H")
	qf1.AwayTeamID = 0
	qf1.AwayTeam = ""
	qf1.AwayTeamCoach = ""
	qf1.AwayTeamRank = 0

	qf2 := structs.Match{
		Model:                  gorm.Model{ID: qf2ID},
		MatchName:              conferenceName + " Tournament Quarterfinals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "B",
		NextGameID:             sf1ID,
		NextGameHOA:            "A",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	qf3 := buildMatch(qf3ID, seed3, seed6, 3, 6, conferenceName+" Tournament Quarterfinals", "B", sf2ID, "H")
	qf3.AwayTeamID = 0
	qf3.AwayTeam = ""
	qf3.AwayTeamCoach = ""
	qf3.AwayTeamRank = 0

	qf4 := buildMatch(qf4ID, seed2, seed7, 2, 7, conferenceName+" Tournament Quarterfinals", "B", sf2ID, "A")
	qf4.AwayTeamID = 0
	qf4.AwayTeam = ""
	qf4.AwayTeamCoach = ""
	qf4.AwayTeamRank = 0

	sf1 := structs.Match{
		Model:                  gorm.Model{ID: sf1ID},
		MatchName:              conferenceName + " Tournament Semifinals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "C",
		NextGameID:             finalID,
		NextGameHOA:            "H",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	sf2 := structs.Match{
		Model:                  gorm.Model{ID: sf2ID},
		MatchName:              conferenceName + " Tournament Semifinals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "C",
		NextGameID:             finalID,
		NextGameHOA:            "A",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	finalMatch := structs.Match{
		Model:                  gorm.Model{ID: finalID},
		MatchName:              conferenceName + " Tournament Finals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "D",
		NextGameID:             0,
		NextGameHOA:            "",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	matches = append(matches, playIn1, playIn2, playIn3, playIn4, playIn5, qf1, qf2, qf3, qf4, sf1, sf2, finalMatch)

	return matches
}

func generateCBB12TeamTournamentMatches(ts structs.Timestamp, latestMatchID uint, standings []structs.CollegeStandings, teamMap map[uint]structs.Team) []structs.Match {
	if len(standings) < 12 {
		return []structs.Match{}
	}

	conferenceName := standings[0].ConferenceName

	seed1 := standings[0]
	seed2 := standings[1]
	seed3 := standings[2]
	seed4 := standings[3]
	seed5 := standings[4]
	seed6 := standings[5]
	seed7 := standings[6]
	seed8 := standings[7]
	seed9 := standings[8]
	seed10 := standings[9]
	seed11 := standings[10]
	seed12 := standings[11]

	playIn1ID := latestMatchID
	playIn2ID := latestMatchID + 1
	playIn3ID := latestMatchID + 2
	playIn4ID := latestMatchID + 3
	qf1ID := latestMatchID + 4
	qf2ID := latestMatchID + 5
	qf3ID := latestMatchID + 6
	qf4ID := latestMatchID + 7
	sf1ID := latestMatchID + 8
	sf2ID := latestMatchID + 9
	finalID := latestMatchID + 10

	buildMatch := func(id uint, homeSeed, awaySeed structs.CollegeStandings, homeRank, awayRank uint, name, slot string, nextID uint, nextHOA string) structs.Match {
		homeTeam := teamMap[homeSeed.TeamID]
		awayTeam := teamMap[awaySeed.TeamID]
		homeCoach := homeTeam.Coach
		if homeCoach == "" {
			homeCoach = "AI"
		}
		awayCoach := awayTeam.Coach
		if awayCoach == "" {
			awayCoach = "AI"
		}

		return structs.Match{
			Model:                  gorm.Model{ID: id},
			MatchName:              name,
			WeekID:                 ts.CollegeWeekID,
			Week:                   uint(ts.CollegeWeek),
			SeasonID:               ts.SeasonID,
			HomeTeamID:             homeSeed.TeamID,
			HomeTeam:               homeTeam.Abbr,
			HomeTeamCoach:          homeCoach,
			HomeTeamRank:           homeRank,
			AwayTeamID:             awaySeed.TeamID,
			AwayTeam:               awayTeam.Abbr,
			AwayTeamCoach:          awayCoach,
			AwayTeamRank:           awayRank,
			MatchOfWeek:            slot,
			NextGameID:             nextID,
			NextGameHOA:            nextHOA,
			IsConference:           true,
			IsConferenceTournament: true,
			IsPlayoffGame:          false,
			IsNeutralSite:          false,
			Arena:                  homeTeam.Arena,
			City:                   homeTeam.City,
			State:                  homeTeam.State,
		}
	}

	matches := []structs.Match{}

	// 4 play-in games: 5v12, 6v11, 7v10, 8v9
	playIn1 := buildMatch(playIn1ID, seed5, seed12, 5, 12, conferenceName+" Tournament First Round", "A", qf2ID, "H")
	playIn2 := buildMatch(playIn2ID, seed6, seed11, 6, 11, conferenceName+" Tournament First Round", "A", qf3ID, "H")
	playIn3 := buildMatch(playIn3ID, seed7, seed10, 7, 10, conferenceName+" Tournament First Round", "A", qf4ID, "H")
	playIn4 := buildMatch(playIn4ID, seed8, seed9, 8, 9, conferenceName+" Tournament First Round", "A", qf1ID, "A")

	// Quarterfinals: 1 vs 8/9 winner, 4 vs 5/12 winner, 3 vs 6/11 winner, 2 vs 7/10 winner
	qf1 := buildMatch(qf1ID, seed1, seed8, 1, 8, conferenceName+" Tournament Quarterfinals", "B", sf1ID, "H")
	qf1.AwayTeamID = 0
	qf1.AwayTeam = ""
	qf1.AwayTeamCoach = ""
	qf1.AwayTeamRank = 0

	qf2 := buildMatch(qf2ID, seed4, seed5, 4, 5, conferenceName+" Tournament Quarterfinals", "B", sf1ID, "A")
	qf2.AwayTeamID = 0
	qf2.AwayTeam = ""
	qf2.AwayTeamCoach = ""
	qf2.AwayTeamRank = 0

	qf3 := buildMatch(qf3ID, seed3, seed6, 3, 6, conferenceName+" Tournament Quarterfinals", "B", sf2ID, "H")
	qf3.AwayTeamID = 0
	qf3.AwayTeam = ""
	qf3.AwayTeamCoach = ""
	qf3.AwayTeamRank = 0

	qf4 := buildMatch(qf4ID, seed2, seed7, 2, 7, conferenceName+" Tournament Quarterfinals", "B", sf2ID, "A")
	qf4.AwayTeamID = 0
	qf4.AwayTeam = ""
	qf4.AwayTeamCoach = ""
	qf4.AwayTeamRank = 0

	sf1 := structs.Match{
		Model:                  gorm.Model{ID: sf1ID},
		MatchName:              conferenceName + " Tournament Semifinals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "C",
		NextGameID:             finalID,
		NextGameHOA:            "H",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	sf2 := structs.Match{
		Model:                  gorm.Model{ID: sf2ID},
		MatchName:              conferenceName + " Tournament Semifinals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "C",
		NextGameID:             finalID,
		NextGameHOA:            "A",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	finalMatch := structs.Match{
		Model:                  gorm.Model{ID: finalID},
		MatchName:              conferenceName + " Tournament Finals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "D",
		NextGameID:             0,
		NextGameHOA:            "",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	matches = append(matches, playIn1, playIn2, playIn3, playIn4, qf1, qf2, qf3, qf4, sf1, sf2, finalMatch)

	return matches
}

func generateCBB11TeamTournamentMatches(ts structs.Timestamp, latestMatchID uint, standings []structs.CollegeStandings, teamMap map[uint]structs.Team) []structs.Match {
	if len(standings) < 11 {
		return []structs.Match{}
	}

	conferenceName := standings[0].ConferenceName

	seed1 := standings[0]
	seed2 := standings[1]
	seed3 := standings[2]
	seed4 := standings[3]
	seed5 := standings[4]
	seed6 := standings[5]
	seed7 := standings[6]
	seed8 := standings[7]
	seed9 := standings[8]
	seed10 := standings[9]
	seed11 := standings[10]

	playIn1ID := latestMatchID
	playIn2ID := latestMatchID + 1
	playIn3ID := latestMatchID + 2
	qf1ID := latestMatchID + 3
	qf2ID := latestMatchID + 4
	qf3ID := latestMatchID + 5
	qf4ID := latestMatchID + 6
	sf1ID := latestMatchID + 7
	sf2ID := latestMatchID + 8
	finalID := latestMatchID + 9

	buildMatch := func(id uint, homeSeed, awaySeed structs.CollegeStandings, homeRank, awayRank uint, name, slot string, nextID uint, nextHOA string) structs.Match {
		homeTeam := teamMap[homeSeed.TeamID]
		awayTeam := teamMap[awaySeed.TeamID]
		homeCoach := homeTeam.Coach
		if homeCoach == "" {
			homeCoach = "AI"
		}
		awayCoach := awayTeam.Coach
		if awayCoach == "" {
			awayCoach = "AI"
		}

		return structs.Match{
			Model:                  gorm.Model{ID: id},
			MatchName:              name,
			WeekID:                 ts.CollegeWeekID,
			Week:                   uint(ts.CollegeWeek),
			SeasonID:               ts.SeasonID,
			HomeTeamID:             homeSeed.TeamID,
			HomeTeam:               homeTeam.Abbr,
			HomeTeamCoach:          homeCoach,
			HomeTeamRank:           homeRank,
			AwayTeamID:             awaySeed.TeamID,
			AwayTeam:               awayTeam.Abbr,
			AwayTeamCoach:          awayCoach,
			AwayTeamRank:           awayRank,
			MatchOfWeek:            slot,
			NextGameID:             nextID,
			NextGameHOA:            nextHOA,
			IsConference:           true,
			IsConferenceTournament: true,
			IsPlayoffGame:          false,
			IsNeutralSite:          false,
			Arena:                  homeTeam.Arena,
			City:                   homeTeam.City,
			State:                  homeTeam.State,
		}
	}

	matches := []structs.Match{}

	playIn1 := buildMatch(playIn1ID, seed8, seed9, 8, 9, conferenceName+" Tournament First Round", "A", qf1ID, "A")
	playIn2 := buildMatch(playIn2ID, seed7, seed10, 7, 10, conferenceName+" Tournament First Round", "A", qf4ID, "A")
	playIn3 := buildMatch(playIn3ID, seed6, seed11, 6, 11, conferenceName+" Tournament First Round", "A", qf3ID, "A")

	qf1 := buildMatch(qf1ID, seed1, seed8, 1, 8, conferenceName+" Tournament Quarterfinals", "B", sf1ID, "H")
	qf1.AwayTeamID = 0
	qf1.AwayTeam = ""
	qf1.AwayTeamCoach = ""
	qf1.AwayTeamRank = 0

	qf2 := buildMatch(qf2ID, seed4, seed5, 4, 5, conferenceName+" Tournament Quarterfinals", "B", sf1ID, "A")

	qf3 := buildMatch(qf3ID, seed3, seed6, 3, 6, conferenceName+" Tournament Quarterfinals", "B", sf2ID, "H")
	qf3.AwayTeamID = 0
	qf3.AwayTeam = ""
	qf3.AwayTeamCoach = ""
	qf3.AwayTeamRank = 0

	qf4 := buildMatch(qf4ID, seed2, seed7, 2, 7, conferenceName+" Tournament Quarterfinals", "B", sf2ID, "A")
	qf4.AwayTeamID = 0
	qf4.AwayTeam = ""
	qf4.AwayTeamCoach = ""
	qf4.AwayTeamRank = 0

	sf1 := structs.Match{
		Model:                  gorm.Model{ID: sf1ID},
		MatchName:              conferenceName + " Tournament Semifinals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "C",
		NextGameID:             finalID,
		NextGameHOA:            "H",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	sf2 := structs.Match{
		Model:                  gorm.Model{ID: sf2ID},
		MatchName:              conferenceName + " Tournament Semifinals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "C",
		NextGameID:             finalID,
		NextGameHOA:            "A",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	finalMatch := structs.Match{
		Model:                  gorm.Model{ID: finalID},
		MatchName:              conferenceName + " Tournament Finals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "D",
		NextGameID:             0,
		NextGameHOA:            "",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	matches = append(matches, playIn1, playIn2, playIn3, qf1, qf2, qf3, qf4, sf1, sf2, finalMatch)

	return matches
}

func generateCBB10TeamTournamentMatches(ts structs.Timestamp, latestMatchID uint, standings []structs.CollegeStandings, teamMap map[uint]structs.Team) []structs.Match {
	if len(standings) < 10 {
		return []structs.Match{}
	}

	conferenceName := standings[0].ConferenceName

	seed1 := standings[0]
	seed2 := standings[1]
	seed3 := standings[2]
	seed4 := standings[3]
	seed5 := standings[4]
	seed6 := standings[5]
	seed7 := standings[6]
	seed8 := standings[7]
	seed9 := standings[8]
	seed10 := standings[9]

	playIn1ID := latestMatchID
	playIn2ID := latestMatchID + 1
	qf1ID := latestMatchID + 2
	qf2ID := latestMatchID + 3
	qf3ID := latestMatchID + 4
	qf4ID := latestMatchID + 5
	sf1ID := latestMatchID + 6
	sf2ID := latestMatchID + 7
	finalID := latestMatchID + 8

	buildMatch := func(id uint, homeSeed, awaySeed structs.CollegeStandings, homeRank, awayRank uint, name, slot string, nextID uint, nextHOA string) structs.Match {
		homeTeam := teamMap[homeSeed.TeamID]
		awayTeam := teamMap[awaySeed.TeamID]
		homeCoach := homeTeam.Coach
		if homeCoach == "" {
			homeCoach = "AI"
		}
		awayCoach := awayTeam.Coach
		if awayCoach == "" {
			awayCoach = "AI"
		}

		return structs.Match{
			Model:                  gorm.Model{ID: id},
			MatchName:              name,
			WeekID:                 ts.CollegeWeekID,
			Week:                   uint(ts.CollegeWeek),
			SeasonID:               ts.SeasonID,
			HomeTeamID:             homeSeed.TeamID,
			HomeTeam:               homeTeam.Abbr,
			HomeTeamCoach:          homeCoach,
			HomeTeamRank:           homeRank,
			AwayTeamID:             awaySeed.TeamID,
			AwayTeam:               awayTeam.Abbr,
			AwayTeamCoach:          awayCoach,
			AwayTeamRank:           awayRank,
			MatchOfWeek:            slot,
			NextGameID:             nextID,
			NextGameHOA:            nextHOA,
			IsConference:           true,
			IsConferenceTournament: true,
			IsPlayoffGame:          false,
			IsNeutralSite:          false,
			Arena:                  homeTeam.Arena,
			City:                   homeTeam.City,
			State:                  homeTeam.State,
		}
	}

	matches := []structs.Match{}

	playIn1 := buildMatch(playIn1ID, seed8, seed9, 8, 9, conferenceName+" Tournament First Round", "A", qf1ID, "A")
	playIn2 := buildMatch(playIn2ID, seed7, seed10, 7, 10, conferenceName+" Tournament First Round", "A", qf4ID, "A")

	qf1 := buildMatch(qf1ID, seed1, seed8, 1, 8, conferenceName+" Tournament Quarterfinals", "B", sf1ID, "H")
	qf1.AwayTeamID = 0
	qf1.AwayTeam = ""
	qf1.AwayTeamCoach = ""
	qf1.AwayTeamRank = 0

	qf2 := buildMatch(qf2ID, seed4, seed5, 4, 5, conferenceName+" Tournament Quarterfinals", "B", sf1ID, "A")
	qf3 := buildMatch(qf3ID, seed3, seed6, 3, 6, conferenceName+" Tournament Quarterfinals", "B", sf2ID, "H")

	qf4 := buildMatch(qf4ID, seed2, seed7, 2, 7, conferenceName+" Tournament Quarterfinals", "B", sf2ID, "A")
	qf4.AwayTeamID = 0
	qf4.AwayTeam = ""
	qf4.AwayTeamCoach = ""
	qf4.AwayTeamRank = 0

	sf1 := structs.Match{
		Model:                  gorm.Model{ID: sf1ID},
		MatchName:              conferenceName + " Tournament Semifinals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "C",
		NextGameID:             finalID,
		NextGameHOA:            "H",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	sf2 := structs.Match{
		Model:                  gorm.Model{ID: sf2ID},
		MatchName:              conferenceName + " Tournament Semifinals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "C",
		NextGameID:             finalID,
		NextGameHOA:            "A",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	finalMatch := structs.Match{
		Model:                  gorm.Model{ID: finalID},
		MatchName:              conferenceName + " Tournament Finals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "D",
		NextGameID:             0,
		NextGameHOA:            "",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	matches = append(matches, playIn1, playIn2, qf1, qf2, qf3, qf4, sf1, sf2, finalMatch)

	return matches
}

func generateCBB9TeamTournamentMatches(ts structs.Timestamp, latestMatchID uint, standings []structs.CollegeStandings, teamMap map[uint]structs.Team) []structs.Match {
	if len(standings) < 9 {
		return []structs.Match{}
	}

	conferenceName := standings[0].ConferenceName

	seed1 := standings[0]
	seed2 := standings[1]
	seed3 := standings[2]
	seed4 := standings[3]
	seed5 := standings[4]
	seed6 := standings[5]
	seed7 := standings[6]
	seed8 := standings[7]
	seed9 := standings[8]

	playInID := latestMatchID
	qf1ID := latestMatchID + 1
	qf2ID := latestMatchID + 2
	qf3ID := latestMatchID + 3
	qf4ID := latestMatchID + 4
	sf1ID := latestMatchID + 5
	sf2ID := latestMatchID + 6
	finalID := latestMatchID + 7

	buildMatch := func(id uint, homeSeed, awaySeed structs.CollegeStandings, homeRank, awayRank uint, name, slot string, nextID uint, nextHOA string) structs.Match {
		homeTeam := teamMap[homeSeed.TeamID]
		awayTeam := teamMap[awaySeed.TeamID]
		homeCoach := homeTeam.Coach
		if homeCoach == "" {
			homeCoach = "AI"
		}
		awayCoach := awayTeam.Coach
		if awayCoach == "" {
			awayCoach = "AI"
		}

		return structs.Match{
			Model:                  gorm.Model{ID: id},
			MatchName:              name,
			WeekID:                 ts.CollegeWeekID,
			Week:                   uint(ts.CollegeWeek),
			SeasonID:               ts.SeasonID,
			HomeTeamID:             homeSeed.TeamID,
			HomeTeam:               homeTeam.Abbr,
			HomeTeamCoach:          homeCoach,
			HomeTeamRank:           homeRank,
			AwayTeamID:             awaySeed.TeamID,
			AwayTeam:               awayTeam.Abbr,
			AwayTeamCoach:          awayCoach,
			AwayTeamRank:           awayRank,
			MatchOfWeek:            slot,
			NextGameID:             nextID,
			NextGameHOA:            nextHOA,
			IsConference:           true,
			IsConferenceTournament: true,
			IsPlayoffGame:          false,
			IsNeutralSite:          false,
			Arena:                  homeTeam.Arena,
			City:                   homeTeam.City,
			State:                  homeTeam.State,
		}
	}

	matches := []structs.Match{}

	playIn := buildMatch(playInID, seed8, seed9, 8, 9, conferenceName+" Tournament First Round", "A", qf1ID, "A")

	qf1 := buildMatch(qf1ID, seed1, seed8, 1, 8, conferenceName+" Tournament Quarterfinals", "B", sf1ID, "H")
	qf1.AwayTeamID = 0
	qf1.AwayTeam = ""
	qf1.AwayTeamCoach = ""
	qf1.AwayTeamRank = 0

	qf2 := buildMatch(qf2ID, seed4, seed5, 4, 5, conferenceName+" Tournament Quarterfinals", "B", sf1ID, "A")
	qf3 := buildMatch(qf3ID, seed3, seed6, 3, 6, conferenceName+" Tournament Quarterfinals", "B", sf2ID, "H")
	qf4 := buildMatch(qf4ID, seed2, seed7, 2, 7, conferenceName+" Tournament Quarterfinals", "B", sf2ID, "A")

	sf1 := structs.Match{
		Model:                  gorm.Model{ID: sf1ID},
		MatchName:              conferenceName + " Tournament Semifinals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "C",
		NextGameID:             finalID,
		NextGameHOA:            "H",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	sf2 := structs.Match{
		Model:                  gorm.Model{ID: sf2ID},
		MatchName:              conferenceName + " Tournament Semifinals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "C",
		NextGameID:             finalID,
		NextGameHOA:            "A",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	finalMatch := structs.Match{
		Model:                  gorm.Model{ID: finalID},
		MatchName:              conferenceName + " Tournament Finals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "D",
		NextGameID:             0,
		NextGameHOA:            "",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	matches = append(matches, playIn, qf1, qf2, qf3, qf4, sf1, sf2, finalMatch)

	return matches
}

func generateCBB8TeamTournamentMatches(ts structs.Timestamp, latestMatchID uint, standings []structs.CollegeStandings, teamMap map[uint]structs.Team) []structs.Match {
	if len(standings) < 8 {
		return []structs.Match{}
	}

	conferenceName := standings[0].ConferenceName

	seed1 := standings[0]
	seed2 := standings[1]
	seed3 := standings[2]
	seed4 := standings[3]
	seed5 := standings[4]
	seed6 := standings[5]
	seed7 := standings[6]
	seed8 := standings[7]

	qf1ID := latestMatchID
	qf2ID := latestMatchID + 1
	qf3ID := latestMatchID + 2
	qf4ID := latestMatchID + 3
	sf1ID := latestMatchID + 4
	sf2ID := latestMatchID + 5
	finalID := latestMatchID + 6

	buildMatch := func(id uint, homeSeed, awaySeed structs.CollegeStandings, homeRank, awayRank uint, name, slot string, nextID uint, nextHOA string) structs.Match {
		homeTeam := teamMap[homeSeed.TeamID]
		awayTeam := teamMap[awaySeed.TeamID]
		homeCoach := homeTeam.Coach
		if homeCoach == "" {
			homeCoach = "AI"
		}
		awayCoach := awayTeam.Coach
		if awayCoach == "" {
			awayCoach = "AI"
		}

		return structs.Match{
			Model:                  gorm.Model{ID: id},
			MatchName:              name,
			WeekID:                 ts.CollegeWeekID,
			Week:                   uint(ts.CollegeWeek),
			SeasonID:               ts.SeasonID,
			HomeTeamID:             homeSeed.TeamID,
			HomeTeam:               homeTeam.Abbr,
			HomeTeamCoach:          homeCoach,
			HomeTeamRank:           homeRank,
			AwayTeamID:             awaySeed.TeamID,
			AwayTeam:               awayTeam.Abbr,
			AwayTeamCoach:          awayCoach,
			AwayTeamRank:           awayRank,
			MatchOfWeek:            slot,
			NextGameID:             nextID,
			NextGameHOA:            nextHOA,
			IsConference:           true,
			IsConferenceTournament: true,
			IsPlayoffGame:          false,
			IsNeutralSite:          false,
			Arena:                  homeTeam.Arena,
			City:                   homeTeam.City,
			State:                  homeTeam.State,
		}
	}

	matches := []structs.Match{}

	qf1 := buildMatch(qf1ID, seed1, seed8, 1, 8, conferenceName+" Tournament Quarterfinals", "A", sf1ID, "H")
	qf2 := buildMatch(qf2ID, seed4, seed5, 4, 5, conferenceName+" Tournament Quarterfinals", "A", sf1ID, "A")
	qf3 := buildMatch(qf3ID, seed3, seed6, 3, 6, conferenceName+" Tournament Quarterfinals", "A", sf2ID, "A")
	qf4 := buildMatch(qf4ID, seed2, seed7, 2, 7, conferenceName+" Tournament Quarterfinals", "A", sf2ID, "H")

	sf1 := structs.Match{
		Model:                  gorm.Model{ID: sf1ID},
		MatchName:              conferenceName + " Tournament Semifinals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "B",
		NextGameID:             finalID,
		NextGameHOA:            "H",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	sf2 := structs.Match{
		Model:                  gorm.Model{ID: sf2ID},
		MatchName:              conferenceName + " Tournament Semifinals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "B",
		NextGameID:             finalID,
		NextGameHOA:            "A",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	finalMatch := structs.Match{
		Model:                  gorm.Model{ID: finalID},
		MatchName:              conferenceName + " Tournament Finals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "C",
		NextGameID:             0,
		NextGameHOA:            "",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	matches = append(matches, qf1, qf2, qf3, qf4, sf1, sf2, finalMatch)

	return matches
}

func generateBigTenTournamentMatches(ts structs.Timestamp, latestMatchID uint, standings []structs.CollegeStandings, teamMap map[uint]structs.Team) []structs.Match {
	if len(standings) < 18 {
		return []structs.Match{}
	}

	conferenceName := standings[0].ConferenceName

	seed1 := standings[0]
	seed2 := standings[1]
	seed3 := standings[2]
	seed4 := standings[3]
	seed5 := standings[4]
	seed6 := standings[5]
	seed7 := standings[6]
	seed8 := standings[7]
	seed9 := standings[8]
	seed10 := standings[9]
	seed11 := standings[10]
	seed12 := standings[11]
	seed13 := standings[12]
	seed14 := standings[13]
	seed15 := standings[14]
	seed16 := standings[15]
	seed17 := standings[16]
	seed18 := standings[17]

	// First Round (4 games) - Slot A
	r1g1ID := latestMatchID
	r1g2ID := latestMatchID + 1
	r1g3ID := latestMatchID + 2
	r1g4ID := latestMatchID + 3

	// Second Round (4 games) - Slot B
	r2g1ID := latestMatchID + 4
	r2g2ID := latestMatchID + 5
	r2g3ID := latestMatchID + 6
	r2g4ID := latestMatchID + 7

	// Third Round (4 games) - Slot C
	r3g1ID := latestMatchID + 8
	r3g2ID := latestMatchID + 9
	r3g3ID := latestMatchID + 10
	r3g4ID := latestMatchID + 11

	// Quarterfinals (4 games) - Slot D
	qf1ID := latestMatchID + 12
	qf2ID := latestMatchID + 13
	qf3ID := latestMatchID + 14
	qf4ID := latestMatchID + 15

	// Semifinals (2 games) - Next week, Slot A
	sf1ID := latestMatchID + 16
	sf2ID := latestMatchID + 17

	// Finals (1 game) - Next week, Slot B
	finalID := latestMatchID + 18

	buildMatch := func(id uint, homeSeed, awaySeed structs.CollegeStandings, homeRank, awayRank uint, name, slot string, nextID uint, nextHOA string) structs.Match {
		homeTeam := teamMap[homeSeed.TeamID]
		awayTeam := teamMap[awaySeed.TeamID]
		homeCoach := homeTeam.Coach
		if homeCoach == "" {
			homeCoach = "AI"
		}
		awayCoach := awayTeam.Coach
		if awayCoach == "" {
			awayCoach = "AI"
		}

		return structs.Match{
			Model:                  gorm.Model{ID: id},
			MatchName:              name,
			WeekID:                 ts.CollegeWeekID,
			Week:                   uint(ts.CollegeWeek),
			SeasonID:               ts.SeasonID,
			HomeTeamID:             homeSeed.TeamID,
			HomeTeam:               homeTeam.Abbr,
			HomeTeamCoach:          homeCoach,
			HomeTeamRank:           homeRank,
			AwayTeamID:             awaySeed.TeamID,
			AwayTeam:               awayTeam.Abbr,
			AwayTeamCoach:          awayCoach,
			AwayTeamRank:           awayRank,
			MatchOfWeek:            slot,
			NextGameID:             nextID,
			NextGameHOA:            nextHOA,
			IsConference:           true,
			IsConferenceTournament: true,
			IsPlayoffGame:          false,
			IsNeutralSite:          false,
			Arena:                  homeTeam.Arena,
			City:                   homeTeam.City,
			State:                  homeTeam.State,
		}
	}

	matches := []structs.Match{}

	// First Round - 4 games (Slot A): 16v17, 15v18, 12v13, 11v14
	r1g1 := buildMatch(r1g1ID, seed16, seed17, 16, 17, conferenceName+" Tournament First Round", "A", r2g1ID, "A")
	r1g2 := buildMatch(r1g2ID, seed15, seed18, 15, 18, conferenceName+" Tournament First Round", "A", r2g2ID, "A")
	r1g3 := buildMatch(r1g3ID, seed12, seed13, 12, 13, conferenceName+" Tournament First Round", "A", r2g3ID, "A")
	r1g4 := buildMatch(r1g4ID, seed11, seed14, 11, 14, conferenceName+" Tournament First Round", "A", r2g4ID, "A")

	// Second Round - 4 games (Slot B): 9 vs 16/17 winner, 10 vs 15/18 winner, 5 vs 12/13 winner, 6 vs 11/14 winner
	r2g1 := buildMatch(r2g1ID, seed9, seed16, 9, 16, conferenceName+" Tournament Second Round", "B", r3g1ID, "A")
	r2g1.AwayTeamID = 0
	r2g1.AwayTeam = ""
	r2g1.AwayTeamCoach = ""
	r2g1.AwayTeamRank = 0

	r2g2 := buildMatch(r2g2ID, seed10, seed15, 10, 15, conferenceName+" Tournament Second Round", "B", r3g2ID, "A")
	r2g2.AwayTeamID = 0
	r2g2.AwayTeam = ""
	r2g2.AwayTeamCoach = ""
	r2g2.AwayTeamRank = 0

	r2g3 := buildMatch(r2g3ID, seed5, seed12, 5, 12, conferenceName+" Tournament Second Round", "B", r3g3ID, "A")
	r2g3.AwayTeamID = 0
	r2g3.AwayTeam = ""
	r2g3.AwayTeamCoach = ""
	r2g3.AwayTeamRank = 0

	r2g4 := buildMatch(r2g4ID, seed6, seed11, 6, 11, conferenceName+" Tournament Second Round", "B", r3g4ID, "A")
	r2g4.AwayTeamID = 0
	r2g4.AwayTeam = ""
	r2g4.AwayTeamCoach = ""
	r2g4.AwayTeamRank = 0

	// Third Round - 4 games (Slot C): 8 vs 9 winner, 7 vs 10 winner, 4 vs 5 winner, 3 vs 6 winner
	r3g1 := buildMatch(r3g1ID, seed8, seed9, 8, 9, conferenceName+" Tournament Third Round", "C", qf1ID, "A")
	r3g1.AwayTeamID = 0
	r3g1.AwayTeam = ""
	r3g1.AwayTeamCoach = ""
	r3g1.AwayTeamRank = 0

	r3g2 := buildMatch(r3g2ID, seed7, seed10, 7, 10, conferenceName+" Tournament Third Round", "C", qf3ID, "A")
	r3g2.AwayTeamID = 0
	r3g2.AwayTeam = ""
	r3g2.AwayTeamCoach = ""
	r3g2.AwayTeamRank = 0

	r3g3 := buildMatch(r3g3ID, seed4, seed5, 4, 5, conferenceName+" Tournament Third Round", "C", qf2ID, "A")
	r3g3.AwayTeamID = 0
	r3g3.AwayTeam = ""
	r3g3.AwayTeamCoach = ""
	r3g3.AwayTeamRank = 0

	r3g4 := buildMatch(r3g4ID, seed3, seed6, 3, 6, conferenceName+" Tournament Third Round", "C", qf4ID, "A")
	r3g4.AwayTeamID = 0
	r3g4.AwayTeam = ""
	r3g4.AwayTeamCoach = ""
	r3g4.AwayTeamRank = 0

	// Quarterfinals - 4 games (Slot D): Seeds 1-2 vs third round winners
	qf1 := buildMatch(qf1ID, seed1, seed8, 1, 8, conferenceName+" Tournament Quarterfinals", "D", sf1ID, "H")
	qf1.AwayTeamID = 0
	qf1.AwayTeam = ""
	qf1.AwayTeamCoach = ""
	qf1.AwayTeamRank = 0

	qf2 := buildMatch(qf2ID, seed4, seed5, 4, 5, conferenceName+" Tournament Quarterfinals", "D", sf1ID, "A")
	qf2.HomeTeamID = 0
	qf2.HomeTeam = ""
	qf2.HomeTeamCoach = ""
	qf2.HomeTeamRank = 0
	qf2.AwayTeamID = 0
	qf2.AwayTeam = ""
	qf2.AwayTeamCoach = ""
	qf2.AwayTeamRank = 0

	qf3 := buildMatch(qf3ID, seed2, seed7, 2, 7, conferenceName+" Tournament Quarterfinals", "D", sf2ID, "H")
	qf3.AwayTeamID = 0
	qf3.AwayTeam = ""
	qf3.AwayTeamCoach = ""
	qf3.AwayTeamRank = 0

	qf4 := buildMatch(qf4ID, seed3, seed6, 3, 6, conferenceName+" Tournament Quarterfinals", "D", sf2ID, "A")
	qf4.HomeTeamID = 0
	qf4.HomeTeam = ""
	qf4.HomeTeamCoach = ""
	qf4.HomeTeamRank = 0
	qf4.AwayTeamID = 0
	qf4.AwayTeam = ""
	qf4.AwayTeamCoach = ""
	qf4.AwayTeamRank = 0

	// Semifinals - 2 games (Next week, Slot A)
	sf1 := structs.Match{
		Model:                  gorm.Model{ID: sf1ID},
		MatchName:              conferenceName + " Tournament Semifinals",
		WeekID:                 ts.CollegeWeekID + 1,
		Week:                   uint(ts.CollegeWeek + 1),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "A",
		NextGameID:             finalID,
		NextGameHOA:            "H",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	sf2 := structs.Match{
		Model:                  gorm.Model{ID: sf2ID},
		MatchName:              conferenceName + " Tournament Semifinals",
		WeekID:                 ts.CollegeWeekID + 1,
		Week:                   uint(ts.CollegeWeek + 1),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "A",
		NextGameID:             finalID,
		NextGameHOA:            "A",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	// Finals - 1 game (Next week, Slot B)
	finalMatch := structs.Match{
		Model:                  gorm.Model{ID: finalID},
		MatchName:              conferenceName + " Tournament Finals",
		WeekID:                 ts.CollegeWeekID + 1,
		Week:                   uint(ts.CollegeWeek + 1),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "B",
		NextGameID:             0,
		NextGameHOA:            "",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	matches = append(matches, r1g1, r1g2, r1g3, r1g4, r2g1, r2g2, r2g3, r2g4, r3g1, r3g2, r3g3, r3g4, qf1, qf2, qf3, qf4, sf1, sf2, finalMatch)

	return matches
}

func generateACCtournamentMatches(ts structs.Timestamp, latestMatchID uint, standings []structs.CollegeStandings, teamMap map[uint]structs.Team) []structs.Match {
	if len(standings) < 15 {
		return []structs.Match{}
	}

	conferenceName := standings[0].ConferenceName

	seed1 := standings[0]
	seed2 := standings[1]
	seed3 := standings[2]
	seed4 := standings[3]
	seed5 := standings[4]
	seed6 := standings[5]
	seed7 := standings[6]
	seed8 := standings[7]
	seed9 := standings[8]
	seed10 := standings[9]
	seed11 := standings[10]
	seed12 := standings[11]
	seed13 := standings[12]
	seed14 := standings[13]
	seed15 := standings[14]

	// First Round (3 games) - Slot A
	r1g1ID := latestMatchID
	r1g2ID := latestMatchID + 1
	r1g3ID := latestMatchID + 2

	// Second Round (4 games) - Slot B
	r2g1ID := latestMatchID + 3
	r2g2ID := latestMatchID + 4
	r2g3ID := latestMatchID + 5
	r2g4ID := latestMatchID + 6

	// Quarterfinals (4 games) - Slot C
	qf1ID := latestMatchID + 7
	qf2ID := latestMatchID + 8
	qf3ID := latestMatchID + 9
	qf4ID := latestMatchID + 10

	// Semifinals (2 games) - Slot D
	sf1ID := latestMatchID + 11
	sf2ID := latestMatchID + 12

	// Finals (1 game) - Next week, Slot A
	finalID := latestMatchID + 13

	buildMatch := func(id uint, homeSeed, awaySeed structs.CollegeStandings, homeRank, awayRank uint, name, slot string, nextID uint, nextHOA string) structs.Match {
		homeTeam := teamMap[homeSeed.TeamID]
		awayTeam := teamMap[awaySeed.TeamID]
		homeCoach := homeTeam.Coach
		if homeCoach == "" {
			homeCoach = "AI"
		}
		awayCoach := awayTeam.Coach
		if awayCoach == "" {
			awayCoach = "AI"
		}

		return structs.Match{
			Model:                  gorm.Model{ID: id},
			MatchName:              name,
			WeekID:                 ts.CollegeWeekID,
			Week:                   uint(ts.CollegeWeek),
			SeasonID:               ts.SeasonID,
			HomeTeamID:             homeSeed.TeamID,
			HomeTeam:               homeTeam.Abbr,
			HomeTeamCoach:          homeCoach,
			HomeTeamRank:           homeRank,
			AwayTeamID:             awaySeed.TeamID,
			AwayTeam:               awayTeam.Abbr,
			AwayTeamCoach:          awayCoach,
			AwayTeamRank:           awayRank,
			MatchOfWeek:            slot,
			NextGameID:             nextID,
			NextGameHOA:            nextHOA,
			IsConference:           true,
			IsConferenceTournament: true,
			IsPlayoffGame:          false,
			IsNeutralSite:          false,
			Arena:                  homeTeam.Arena,
			City:                   homeTeam.City,
			State:                  homeTeam.State,
		}
	}

	matches := []structs.Match{}

	// First Round - 3 games (Slot A): 12v13, 11v14, 10v15
	r1g1 := buildMatch(r1g1ID, seed12, seed13, 12, 13, conferenceName+" Tournament First Round", "A", r2g1ID, "A")
	r1g2 := buildMatch(r1g2ID, seed11, seed14, 11, 14, conferenceName+" Tournament First Round", "A", r2g2ID, "A")
	r1g3 := buildMatch(r1g3ID, seed10, seed15, 10, 15, conferenceName+" Tournament First Round", "A", r2g3ID, "A")

	// Second Round - 4 games (Slot B): 5 vs 12/13 winner, 6 vs 11/14 winner, 7 vs 10/15 winner, 8v9
	r2g1 := buildMatch(r2g1ID, seed5, seed12, 5, 12, conferenceName+" Tournament Second Round", "B", qf4ID, "A")
	r2g1.AwayTeamID = 0
	r2g1.AwayTeam = ""
	r2g1.AwayTeamCoach = ""
	r2g1.AwayTeamRank = 0

	r2g2 := buildMatch(r2g2ID, seed6, seed11, 6, 11, conferenceName+" Tournament Second Round", "B", qf3ID, "A")
	r2g2.AwayTeamID = 0
	r2g2.AwayTeam = ""
	r2g2.AwayTeamCoach = ""
	r2g2.AwayTeamRank = 0

	r2g3 := buildMatch(r2g3ID, seed7, seed10, 7, 10, conferenceName+" Tournament Second Round", "B", qf2ID, "A")
	r2g3.AwayTeamID = 0
	r2g3.AwayTeam = ""
	r2g3.AwayTeamCoach = ""
	r2g3.AwayTeamRank = 0

	r2g4 := buildMatch(r2g4ID, seed8, seed9, 8, 9, conferenceName+" Tournament Second Round", "B", qf1ID, "A")

	// Quarterfinals - 4 games (Slot C): Seeds 1-4 vs second round winners
	qf1 := buildMatch(qf1ID, seed1, seed8, 1, 8, conferenceName+" Tournament Quarterfinals", "C", sf1ID, "H")
	qf1.AwayTeamID = 0
	qf1.AwayTeam = ""
	qf1.AwayTeamCoach = ""
	qf1.AwayTeamRank = 0

	qf2 := buildMatch(qf2ID, seed2, seed7, 2, 7, conferenceName+" Tournament Quarterfinals", "C", sf1ID, "A")
	qf2.AwayTeamID = 0
	qf2.AwayTeam = ""
	qf2.AwayTeamCoach = ""
	qf2.AwayTeamRank = 0

	qf3 := buildMatch(qf3ID, seed3, seed6, 3, 6, conferenceName+" Tournament Quarterfinals", "C", sf2ID, "H")
	qf3.AwayTeamID = 0
	qf3.AwayTeam = ""
	qf3.AwayTeamCoach = ""
	qf3.AwayTeamRank = 0

	qf4 := buildMatch(qf4ID, seed4, seed5, 4, 5, conferenceName+" Tournament Quarterfinals", "C", sf2ID, "A")
	qf4.AwayTeamID = 0
	qf4.AwayTeam = ""
	qf4.AwayTeamCoach = ""
	qf4.AwayTeamRank = 0

	// Semifinals - 2 games (Slot D)
	sf1 := structs.Match{
		Model:                  gorm.Model{ID: sf1ID},
		MatchName:              conferenceName + " Tournament Semifinals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "D",
		NextGameID:             finalID,
		NextGameHOA:            "H",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	sf2 := structs.Match{
		Model:                  gorm.Model{ID: sf2ID},
		MatchName:              conferenceName + " Tournament Semifinals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "D",
		NextGameID:             finalID,
		NextGameHOA:            "A",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	// Finals - 1 game (Next week, Slot A)
	finalMatch := structs.Match{
		Model:                  gorm.Model{ID: finalID},
		MatchName:              conferenceName + " Tournament Finals",
		WeekID:                 ts.CollegeWeekID + 1,
		Week:                   uint(ts.CollegeWeek + 1),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "A",
		NextGameID:             0,
		NextGameHOA:            "",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	matches = append(matches, r1g1, r1g2, r1g3, r2g1, r2g2, r2g3, r2g4, qf1, qf2, qf3, qf4, sf1, sf2, finalMatch)

	return matches
}

func generateBig12TournamentMatches(ts structs.Timestamp, latestMatchID uint, standings []structs.CollegeStandings, teamMap map[uint]structs.Team) []structs.Match {
	if len(standings) < 16 {
		return []structs.Match{}
	}

	conferenceName := standings[0].ConferenceName

	seed1 := standings[0]
	seed2 := standings[1]
	seed3 := standings[2]
	seed4 := standings[3]
	seed5 := standings[4]
	seed6 := standings[5]
	seed7 := standings[6]
	seed8 := standings[7]
	seed9 := standings[8]
	seed10 := standings[9]
	seed11 := standings[10]
	seed12 := standings[11]
	seed13 := standings[12]
	seed14 := standings[13]
	seed15 := standings[14]
	seed16 := standings[15]

	// First Round (4 games) - Slot A
	r1g1ID := latestMatchID
	r1g2ID := latestMatchID + 1
	r1g3ID := latestMatchID + 2
	r1g4ID := latestMatchID + 3

	// Second Round (4 games) - Slot B
	r2g1ID := latestMatchID + 4
	r2g2ID := latestMatchID + 5
	r2g3ID := latestMatchID + 6
	r2g4ID := latestMatchID + 7

	// Quarterfinals (4 games) - Slot C
	qf1ID := latestMatchID + 8
	qf2ID := latestMatchID + 9
	qf3ID := latestMatchID + 10
	qf4ID := latestMatchID + 11

	// Semifinals (2 games) - Slot D
	sf1ID := latestMatchID + 12
	sf2ID := latestMatchID + 13

	// Finals (1 game) - Next week, Slot A
	finalID := latestMatchID + 14

	buildMatch := func(id uint, homeSeed, awaySeed structs.CollegeStandings, homeRank, awayRank uint, name, slot string, nextID uint, nextHOA string) structs.Match {
		homeTeam := teamMap[homeSeed.TeamID]
		awayTeam := teamMap[awaySeed.TeamID]
		homeCoach := homeTeam.Coach
		if homeCoach == "" {
			homeCoach = "AI"
		}
		awayCoach := awayTeam.Coach
		if awayCoach == "" {
			awayCoach = "AI"
		}

		return structs.Match{
			Model:                  gorm.Model{ID: id},
			MatchName:              name,
			WeekID:                 ts.CollegeWeekID,
			Week:                   uint(ts.CollegeWeek),
			SeasonID:               ts.SeasonID,
			HomeTeamID:             homeSeed.TeamID,
			HomeTeam:               homeTeam.Abbr,
			HomeTeamCoach:          homeCoach,
			HomeTeamRank:           homeRank,
			AwayTeamID:             awaySeed.TeamID,
			AwayTeam:               awayTeam.Abbr,
			AwayTeamCoach:          awayCoach,
			AwayTeamRank:           awayRank,
			MatchOfWeek:            slot,
			NextGameID:             nextID,
			NextGameHOA:            nextHOA,
			IsConference:           true,
			IsConferenceTournament: true,
			IsPlayoffGame:          false,
			IsNeutralSite:          false,
			Arena:                  homeTeam.Arena,
			City:                   homeTeam.City,
			State:                  homeTeam.State,
		}
	}

	matches := []structs.Match{}

	// First Round - 4 games (Slot A)
	r1g1 := buildMatch(r1g1ID, seed12, seed13, 12, 13, conferenceName+" Tournament First Round", "A", r2g1ID, "A")
	r1g2 := buildMatch(r1g2ID, seed11, seed14, 11, 14, conferenceName+" Tournament First Round", "A", r2g2ID, "A")
	r1g3 := buildMatch(r1g3ID, seed10, seed15, 10, 15, conferenceName+" Tournament First Round", "A", r2g3ID, "A")
	r1g4 := buildMatch(r1g4ID, seed9, seed16, 9, 16, conferenceName+" Tournament First Round", "A", r2g4ID, "A")

	// Second Round - 4 games (Slot B): Seeds 5-8 vs first round winners
	r2g1 := buildMatch(r2g1ID, seed5, seed12, 5, 12, conferenceName+" Tournament Second Round", "B", qf1ID, "A")
	r2g1.AwayTeamID = 0
	r2g1.AwayTeam = ""
	r2g1.AwayTeamCoach = ""
	r2g1.AwayTeamRank = 0

	r2g2 := buildMatch(r2g2ID, seed6, seed11, 6, 11, conferenceName+" Tournament Second Round", "B", qf2ID, "A")
	r2g2.AwayTeamID = 0
	r2g2.AwayTeam = ""
	r2g2.AwayTeamCoach = ""
	r2g2.AwayTeamRank = 0

	r2g3 := buildMatch(r2g3ID, seed7, seed10, 7, 10, conferenceName+" Tournament Second Round", "B", qf3ID, "A")
	r2g3.AwayTeamID = 0
	r2g3.AwayTeam = ""
	r2g3.AwayTeamCoach = ""
	r2g3.AwayTeamRank = 0

	r2g4 := buildMatch(r2g4ID, seed8, seed9, 8, 9, conferenceName+" Tournament Second Round", "B", qf4ID, "A")
	r2g4.AwayTeamID = 0
	r2g4.AwayTeam = ""
	r2g4.AwayTeamCoach = ""
	r2g4.AwayTeamRank = 0

	// Quarterfinals - 4 games (Slot C): Seeds 1-4 vs second round winners
	qf1 := buildMatch(qf1ID, seed4, seed5, 4, 5, conferenceName+" Tournament Quarterfinals", "C", sf1ID, "A")
	qf1.AwayTeamID = 0
	qf1.AwayTeam = ""
	qf1.AwayTeamCoach = ""
	qf1.AwayTeamRank = 0

	qf2 := buildMatch(qf2ID, seed3, seed6, 3, 6, conferenceName+" Tournament Quarterfinals", "C", sf1ID, "H")
	qf2.AwayTeamID = 0
	qf2.AwayTeam = ""
	qf2.AwayTeamCoach = ""
	qf2.AwayTeamRank = 0

	qf3 := buildMatch(qf3ID, seed2, seed7, 2, 7, conferenceName+" Tournament Quarterfinals", "C", sf2ID, "A")
	qf3.AwayTeamID = 0
	qf3.AwayTeam = ""
	qf3.AwayTeamCoach = ""
	qf3.AwayTeamRank = 0

	qf4 := buildMatch(qf4ID, seed1, seed8, 1, 8, conferenceName+" Tournament Quarterfinals", "C", sf2ID, "H")
	qf4.AwayTeamID = 0
	qf4.AwayTeam = ""
	qf4.AwayTeamCoach = ""
	qf4.AwayTeamRank = 0

	// Semifinals - 2 games (Slot D)
	sf1 := structs.Match{
		Model:                  gorm.Model{ID: sf1ID},
		MatchName:              conferenceName + " Tournament Semifinals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "D",
		NextGameID:             finalID,
		NextGameHOA:            "H",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	sf2 := structs.Match{
		Model:                  gorm.Model{ID: sf2ID},
		MatchName:              conferenceName + " Tournament Semifinals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "D",
		NextGameID:             finalID,
		NextGameHOA:            "A",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	// Finals - 1 game (Next week, Slot A)
	finalMatch := structs.Match{
		Model:                  gorm.Model{ID: finalID},
		MatchName:              conferenceName + " Tournament Finals",
		WeekID:                 ts.CollegeWeekID + 1,
		Week:                   uint(ts.CollegeWeek + 1),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "A",
		NextGameID:             0,
		NextGameHOA:            "",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	matches = append(matches, r1g1, r1g2, r1g3, r1g4, r2g1, r2g2, r2g3, r2g4, qf1, qf2, qf3, qf4, sf1, sf2, finalMatch)

	return matches
}

func generateSECTournamentMatches(ts structs.Timestamp, latestMatchID uint, standings []structs.CollegeStandings, teamMap map[uint]structs.Team) []structs.Match {
	if len(standings) < 16 {
		return []structs.Match{}
	}

	conferenceName := standings[0].ConferenceName

	seed1 := standings[0]
	seed2 := standings[1]
	seed3 := standings[2]
	seed4 := standings[3]
	seed5 := standings[4]
	seed6 := standings[5]
	seed7 := standings[6]
	seed8 := standings[7]
	seed9 := standings[8]
	seed10 := standings[9]
	seed11 := standings[10]
	seed12 := standings[11]
	seed13 := standings[12]
	seed14 := standings[13]
	seed15 := standings[14]
	seed16 := standings[15]

	// First Round (4 games) - Slot A
	r1g1ID := latestMatchID
	r1g2ID := latestMatchID + 1
	r1g3ID := latestMatchID + 2
	r1g4ID := latestMatchID + 3

	// Second Round (4 games) - Slot B
	r2g1ID := latestMatchID + 4
	r2g2ID := latestMatchID + 5
	r2g3ID := latestMatchID + 6
	r2g4ID := latestMatchID + 7

	// Quarterfinals (4 games) - Slot C
	qf1ID := latestMatchID + 8
	qf2ID := latestMatchID + 9
	qf3ID := latestMatchID + 10
	qf4ID := latestMatchID + 11

	// Semifinals (2 games) - Slot D
	sf1ID := latestMatchID + 12
	sf2ID := latestMatchID + 13

	// Finals (1 game) - Next week, Slot A
	finalID := latestMatchID + 14

	buildMatch := func(id uint, homeSeed, awaySeed structs.CollegeStandings, homeRank, awayRank uint, name, slot string, nextID uint, nextHOA string) structs.Match {
		homeTeam := teamMap[homeSeed.TeamID]
		awayTeam := teamMap[awaySeed.TeamID]
		homeCoach := homeTeam.Coach
		if homeCoach == "" {
			homeCoach = "AI"
		}
		awayCoach := awayTeam.Coach
		if awayCoach == "" {
			awayCoach = "AI"
		}

		return structs.Match{
			Model:                  gorm.Model{ID: id},
			MatchName:              name,
			WeekID:                 ts.CollegeWeekID,
			Week:                   uint(ts.CollegeWeek),
			SeasonID:               ts.SeasonID,
			HomeTeamID:             homeSeed.TeamID,
			HomeTeam:               homeTeam.Abbr,
			HomeTeamCoach:          homeCoach,
			HomeTeamRank:           homeRank,
			AwayTeamID:             awaySeed.TeamID,
			AwayTeam:               awayTeam.Abbr,
			AwayTeamCoach:          awayCoach,
			AwayTeamRank:           awayRank,
			MatchOfWeek:            slot,
			NextGameID:             nextID,
			NextGameHOA:            nextHOA,
			IsConference:           true,
			IsConferenceTournament: true,
			IsPlayoffGame:          false,
			IsNeutralSite:          false,
			Arena:                  homeTeam.Arena,
			City:                   homeTeam.City,
			State:                  homeTeam.State,
		}
	}

	matches := []structs.Match{}

	// First Round - 4 games (Slot A): 16v9, 13v12, 15v10, 14v11
	r1g1 := buildMatch(r1g1ID, seed16, seed9, 16, 9, conferenceName+" Tournament First Round", "A", r2g1ID, "A")
	r1g2 := buildMatch(r1g2ID, seed13, seed12, 13, 12, conferenceName+" Tournament First Round", "A", r2g2ID, "A")
	r1g3 := buildMatch(r1g3ID, seed15, seed10, 15, 10, conferenceName+" Tournament First Round", "A", r2g3ID, "A")
	r1g4 := buildMatch(r1g4ID, seed14, seed11, 14, 11, conferenceName+" Tournament First Round", "A", r2g4ID, "A")

	// Second Round - 4 games (Slot B): 8v(16/9), 5v(13/12), 7v(15/10), 6v(14/11)
	r2g1 := buildMatch(r2g1ID, seed8, seed16, 8, 16, conferenceName+" Tournament Second Round", "B", qf1ID, "A")
	r2g1.AwayTeamID = 0
	r2g1.AwayTeam = ""
	r2g1.AwayTeamCoach = ""
	r2g1.AwayTeamRank = 0

	r2g2 := buildMatch(r2g2ID, seed5, seed13, 5, 13, conferenceName+" Tournament Second Round", "B", qf2ID, "A")
	r2g2.AwayTeamID = 0
	r2g2.AwayTeam = ""
	r2g2.AwayTeamCoach = ""
	r2g2.AwayTeamRank = 0

	r2g3 := buildMatch(r2g3ID, seed7, seed15, 7, 15, conferenceName+" Tournament Second Round", "B", qf3ID, "A")
	r2g3.AwayTeamID = 0
	r2g3.AwayTeam = ""
	r2g3.AwayTeamCoach = ""
	r2g3.AwayTeamRank = 0

	r2g4 := buildMatch(r2g4ID, seed6, seed14, 6, 14, conferenceName+" Tournament Second Round", "B", qf4ID, "A")
	r2g4.AwayTeamID = 0
	r2g4.AwayTeam = ""
	r2g4.AwayTeamCoach = ""
	r2g4.AwayTeamRank = 0

	// Quarterfinals - 4 games (Slot C): 1v(r2g1), 4v(r2g2), 2v(r2g3), 3v(r2g4)
	qf1 := buildMatch(qf1ID, seed1, seed8, 1, 8, conferenceName+" Tournament Quarterfinals", "C", sf1ID, "H")
	qf1.AwayTeamID = 0
	qf1.AwayTeam = ""
	qf1.AwayTeamCoach = ""
	qf1.AwayTeamRank = 0

	qf2 := buildMatch(qf2ID, seed4, seed5, 4, 5, conferenceName+" Tournament Quarterfinals", "C", sf1ID, "A")
	qf2.AwayTeamID = 0
	qf2.AwayTeam = ""
	qf2.AwayTeamCoach = ""
	qf2.AwayTeamRank = 0

	qf3 := buildMatch(qf3ID, seed2, seed7, 2, 7, conferenceName+" Tournament Quarterfinals", "C", sf2ID, "H")
	qf3.AwayTeamID = 0
	qf3.AwayTeam = ""
	qf3.AwayTeamCoach = ""
	qf3.AwayTeamRank = 0

	qf4 := buildMatch(qf4ID, seed3, seed6, 3, 6, conferenceName+" Tournament Quarterfinals", "C", sf2ID, "A")
	qf4.AwayTeamID = 0
	qf4.AwayTeam = ""
	qf4.AwayTeamCoach = ""
	qf4.AwayTeamRank = 0

	// Semifinals - 2 games (Slot D)
	sf1 := structs.Match{
		Model:                  gorm.Model{ID: sf1ID},
		MatchName:              conferenceName + " Tournament Semifinals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "D",
		NextGameID:             finalID,
		NextGameHOA:            "H",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	sf2 := structs.Match{
		Model:                  gorm.Model{ID: sf2ID},
		MatchName:              conferenceName + " Tournament Semifinals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "D",
		NextGameID:             finalID,
		NextGameHOA:            "A",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	// Finals - 1 game (Next week, Slot A)
	finalMatch := structs.Match{
		Model:                  gorm.Model{ID: finalID},
		MatchName:              conferenceName + " Tournament Finals",
		WeekID:                 ts.CollegeWeekID + 1,
		Week:                   uint(ts.CollegeWeek + 1),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "A",
		NextGameID:             0,
		NextGameHOA:            "",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	matches = append(matches, r1g1, r1g2, r1g3, r1g4, r2g1, r2g2, r2g3, r2g4, qf1, qf2, qf3, qf4, sf1, sf2, finalMatch)

	return matches
}

func generateBigEastTournamentMatches(ts structs.Timestamp, latestMatchID uint, standings []structs.CollegeStandings, teamMap map[uint]structs.Team) []structs.Match {
	if len(standings) < 11 {
		return []structs.Match{}
	}

	conferenceName := standings[0].ConferenceName

	seed1 := standings[0]
	seed2 := standings[1]
	seed3 := standings[2]
	seed4 := standings[3]
	seed5 := standings[4]
	seed6 := standings[5]
	seed7 := standings[6]
	seed8 := standings[7]
	seed9 := standings[8]
	seed10 := standings[9]
	seed11 := standings[10]

	// First Round (3 games) - Slot A
	playIn1ID := latestMatchID
	playIn2ID := latestMatchID + 1
	playIn3ID := latestMatchID + 2

	// Quarterfinals (4 games) - Slot B
	qf1ID := latestMatchID + 3
	qf2ID := latestMatchID + 4
	qf3ID := latestMatchID + 5
	qf4ID := latestMatchID + 6

	// Semifinals (2 games) - Slot C
	sf1ID := latestMatchID + 7
	sf2ID := latestMatchID + 8

	// Finals (1 game) - Slot D
	finalID := latestMatchID + 9

	buildMatch := func(id uint, homeSeed, awaySeed structs.CollegeStandings, homeRank, awayRank uint, name, slot string, nextID uint, nextHOA string) structs.Match {
		homeTeam := teamMap[homeSeed.TeamID]
		awayTeam := teamMap[awaySeed.TeamID]
		homeCoach := homeTeam.Coach
		if homeCoach == "" {
			homeCoach = "AI"
		}
		awayCoach := awayTeam.Coach
		if awayCoach == "" {
			awayCoach = "AI"
		}

		return structs.Match{
			Model:                  gorm.Model{ID: id},
			MatchName:              name,
			WeekID:                 ts.CollegeWeekID,
			Week:                   uint(ts.CollegeWeek),
			SeasonID:               ts.SeasonID,
			HomeTeamID:             homeSeed.TeamID,
			HomeTeam:               homeTeam.Abbr,
			HomeTeamCoach:          homeCoach,
			HomeTeamRank:           homeRank,
			AwayTeamID:             awaySeed.TeamID,
			AwayTeam:               awayTeam.Abbr,
			AwayTeamCoach:          awayCoach,
			AwayTeamRank:           awayRank,
			MatchOfWeek:            slot,
			NextGameID:             nextID,
			NextGameHOA:            nextHOA,
			IsConference:           true,
			IsConferenceTournament: true,
			IsPlayoffGame:          false,
			IsNeutralSite:          false,
			Arena:                  homeTeam.Arena,
			City:                   homeTeam.City,
			State:                  homeTeam.State,
		}
	}

	matches := []structs.Match{}

	// First Round (Slot A): 3 games - 8v9, 7v10, 6v11
	playIn1 := buildMatch(playIn1ID, seed8, seed9, 8, 9, conferenceName+" Tournament First Round", "A", qf2ID, "A")
	playIn2 := buildMatch(playIn2ID, seed7, seed10, 7, 10, conferenceName+" Tournament First Round", "A", qf3ID, "A")
	playIn3 := buildMatch(playIn3ID, seed6, seed11, 6, 11, conferenceName+" Tournament First Round", "A", qf4ID, "A")

	// Quarterfinals (Slot B): 4 games
	// 4v5 (pre-seeded)
	qf1 := buildMatch(qf1ID, seed4, seed5, 4, 5, conferenceName+" Tournament Quarterfinals", "B", sf1ID, "A")

	// 1 vs 8/9 winner
	qf2 := buildMatch(qf2ID, seed1, seed8, 1, 8, conferenceName+" Tournament Quarterfinals", "B", sf1ID, "H")
	qf2.AwayTeamID = 0
	qf2.AwayTeam = ""
	qf2.AwayTeamCoach = ""
	qf2.AwayTeamRank = 0

	// 2 vs 7/10 winner
	qf3 := buildMatch(qf3ID, seed2, seed7, 2, 7, conferenceName+" Tournament Quarterfinals", "B", sf2ID, "H")
	qf3.AwayTeamID = 0
	qf3.AwayTeam = ""
	qf3.AwayTeamCoach = ""
	qf3.AwayTeamRank = 0

	// 3 vs 6/11 winner
	qf4 := buildMatch(qf4ID, seed3, seed6, 3, 6, conferenceName+" Tournament Quarterfinals", "B", sf2ID, "A")
	qf4.AwayTeamID = 0
	qf4.AwayTeam = ""
	qf4.AwayTeamCoach = ""
	qf4.AwayTeamRank = 0

	// Semifinals (Slot C): 2 games
	sf1 := structs.Match{
		Model:                  gorm.Model{ID: sf1ID},
		MatchName:              conferenceName + " Tournament Semifinals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "C",
		NextGameID:             finalID,
		NextGameHOA:            "H",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	sf2 := structs.Match{
		Model:                  gorm.Model{ID: sf2ID},
		MatchName:              conferenceName + " Tournament Semifinals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "C",
		NextGameID:             finalID,
		NextGameHOA:            "A",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	// Finals (Slot D): 1 game
	finalMatch := structs.Match{
		Model:                  gorm.Model{ID: finalID},
		MatchName:              conferenceName + " Tournament Finals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "D",
		NextGameID:             0,
		NextGameHOA:            "",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	matches = append(matches, playIn1, playIn2, playIn3, qf1, qf2, qf3, qf4, sf1, sf2, finalMatch)

	return matches
}

func generateIvyLeagueTournamentMatches(ts structs.Timestamp, latestMatchID uint, standings []structs.CollegeStandings, teamMap map[uint]structs.Team) []structs.Match {
	if len(standings) < 8 {
		return []structs.Match{}
	}

	conferenceName := standings[0].ConferenceName

	seed1 := standings[0]
	seed2 := standings[1]
	seed3 := standings[2]
	seed4 := standings[3]

	sf1ID := latestMatchID + 1
	sf2ID := latestMatchID + 2
	finalID := latestMatchID + 3

	buildMatch := func(id uint, homeSeed, awaySeed structs.CollegeStandings, homeRank, awayRank uint, name, slot string, nextID uint, nextHOA string) structs.Match {
		homeTeam := teamMap[homeSeed.TeamID]
		awayTeam := teamMap[awaySeed.TeamID]
		homeCoach := homeTeam.Coach
		if homeCoach == "" {
			homeCoach = "AI"
		}
		awayCoach := awayTeam.Coach
		if awayCoach == "" {
			awayCoach = "AI"
		}

		return structs.Match{
			Model:                  gorm.Model{ID: id},
			MatchName:              name,
			WeekID:                 ts.CollegeWeekID,
			Week:                   uint(ts.CollegeWeek),
			SeasonID:               ts.SeasonID,
			HomeTeamID:             homeSeed.TeamID,
			HomeTeam:               homeTeam.Abbr,
			HomeTeamCoach:          homeCoach,
			HomeTeamRank:           homeRank,
			AwayTeamID:             awaySeed.TeamID,
			AwayTeam:               awayTeam.Abbr,
			AwayTeamCoach:          awayCoach,
			AwayTeamRank:           awayRank,
			MatchOfWeek:            slot,
			NextGameID:             nextID,
			NextGameHOA:            nextHOA,
			IsConference:           true,
			IsConferenceTournament: true,
			IsPlayoffGame:          false,
			IsNeutralSite:          false,
			Arena:                  homeTeam.Arena,
			City:                   homeTeam.City,
			State:                  homeTeam.State,
		}
	}

	matches := []structs.Match{}

	sf1 := buildMatch(sf1ID, seed1, seed4, 1, 4, conferenceName+" Tournament Semifinals", "A", finalID, "H")
	sf2 := buildMatch(sf2ID, seed2, seed3, 2, 3, conferenceName+" Tournament Semifinals", "A", finalID, "A")

	finalMatch := structs.Match{
		Model:                  gorm.Model{ID: finalID},
		MatchName:              conferenceName + " Tournament Finals",
		WeekID:                 ts.CollegeWeekID,
		Week:                   uint(ts.CollegeWeek),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "C",
		NextGameID:             0,
		NextGameHOA:            "",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	matches = append(matches, sf1, sf2, finalMatch)

	return matches
}

func generateWCCTournamentMatches(ts structs.Timestamp, latestMatchID uint, standings []structs.CollegeStandings, teamMap map[uint]structs.Team) []structs.Match {
	if len(standings) < 11 {
		return []structs.Match{}
	}

	conferenceName := standings[0].ConferenceName

	seed1 := standings[0]
	seed2 := standings[1]
	seed3 := standings[2]
	seed4 := standings[3]
	seed5 := standings[4]
	seed6 := standings[5]
	seed7 := standings[6]
	seed8 := standings[7]
	seed9 := standings[8]
	seed10 := standings[9]
	seed11 := standings[10]

	// 6-round progressive structure
	firstRoundID := latestMatchID
	r2g1ID := latestMatchID + 1
	r2g2ID := latestMatchID + 2
	r3g1ID := latestMatchID + 3
	r3g2ID := latestMatchID + 4
	qf1ID := latestMatchID + 5
	qf2ID := latestMatchID + 6
	sf1ID := latestMatchID + 7
	sf2ID := latestMatchID + 8
	finalID := latestMatchID + 9

	buildMatch := func(id uint, homeSeed, awaySeed structs.CollegeStandings, homeRank, awayRank uint, name, slot string, nextID uint, nextHOA string) structs.Match {
		homeTeam := teamMap[homeSeed.TeamID]
		awayTeam := teamMap[awaySeed.TeamID]
		homeCoach := homeTeam.Coach
		if homeCoach == "" {
			homeCoach = "AI"
		}
		awayCoach := awayTeam.Coach
		if awayCoach == "" {
			awayCoach = "AI"
		}

		return structs.Match{
			Model:                  gorm.Model{ID: id},
			MatchName:              name,
			WeekID:                 ts.CollegeWeekID,
			Week:                   uint(ts.CollegeWeek),
			SeasonID:               ts.SeasonID,
			HomeTeamID:             homeSeed.TeamID,
			HomeTeam:               homeTeam.Abbr,
			HomeTeamCoach:          homeCoach,
			HomeTeamRank:           homeRank,
			AwayTeamID:             awaySeed.TeamID,
			AwayTeam:               awayTeam.Abbr,
			AwayTeamCoach:          awayCoach,
			AwayTeamRank:           awayRank,
			MatchOfWeek:            slot,
			NextGameID:             nextID,
			NextGameHOA:            nextHOA,
			IsConference:           true,
			IsConferenceTournament: true,
			IsPlayoffGame:          false,
			IsNeutralSite:          false,
			Arena:                  homeTeam.Arena,
			City:                   homeTeam.City,
			State:                  homeTeam.State,
		}
	}

	matches := []structs.Match{}

	// First Round (Slot A): 1 game - 10v11
	firstRound := buildMatch(firstRoundID, seed10, seed11, 10, 11, conferenceName+" Tournament First Round", "A", r2g1ID, "A")

	// Second Round (Slot B): 2 games - 8v9, 7 vs 10/11 winner
	r2g1 := buildMatch(r2g1ID, seed7, seed10, 7, 10, conferenceName+" Tournament Second Round", "B", r3g1ID, "A")
	r2g1.AwayTeamID = 0
	r2g1.AwayTeam = ""
	r2g1.AwayTeamCoach = ""
	r2g1.AwayTeamRank = 0

	r2g2 := buildMatch(r2g2ID, seed8, seed9, 8, 9, conferenceName+" Tournament Second Round", "B", r3g2ID, "A")

	// Third Round (Slot C): 2 games - 6 vs 7 winner, 5 vs 8/9 winner
	r3g1 := buildMatch(r3g1ID, seed6, seed7, 6, 7, conferenceName+" Tournament Third Round", "C", qf1ID, "A")
	r3g1.AwayTeamID = 0
	r3g1.AwayTeam = ""
	r3g1.AwayTeamCoach = ""
	r3g1.AwayTeamRank = 0

	r3g2 := buildMatch(r3g2ID, seed5, seed8, 5, 8, conferenceName+" Tournament Third Round", "C", qf2ID, "A")
	r3g2.AwayTeamID = 0
	r3g2.AwayTeam = ""
	r3g2.AwayTeamRank = 0

	// Quarterfinals (Slot D): 2 games - 4 vs 6 winner, 3 vs 5 winner
	qf1 := buildMatch(qf1ID, seed4, seed6, 4, 6, conferenceName+" Tournament Quarterfinals", "D", sf1ID, "A")
	qf1.AwayTeamID = 0
	qf1.AwayTeam = ""
	qf1.AwayTeamCoach = ""
	qf1.AwayTeamRank = 0

	qf2 := buildMatch(qf2ID, seed3, seed5, 3, 5, conferenceName+" Tournament Quarterfinals", "D", sf2ID, "A")
	qf2.AwayTeamID = 0
	qf2.AwayTeam = ""
	qf2.AwayTeamCoach = ""
	qf2.AwayTeamRank = 0

	// Semifinals (Next week, Slot A): 2 games - 2 vs 4 winner, 1 vs 3 winner
	sf1 := buildMatch(sf1ID, seed2, seed4, 2, 4, conferenceName+" Tournament Semifinals", "A", finalID, "H")
	sf1.WeekID = ts.CollegeWeekID + 1
	sf1.Week = uint(ts.CollegeWeek + 1)
	sf1.AwayTeamID = 0
	sf1.AwayTeam = ""
	sf1.AwayTeamCoach = ""
	sf1.AwayTeamRank = 0

	sf2 := buildMatch(sf2ID, seed1, seed3, 1, 3, conferenceName+" Tournament Semifinals", "A", finalID, "A")
	sf2.WeekID = ts.CollegeWeekID + 1
	sf2.Week = uint(ts.CollegeWeek + 1)
	sf2.AwayTeamID = 0
	sf2.AwayTeam = ""
	sf2.AwayTeamCoach = ""
	sf2.AwayTeamRank = 0

	// Championship (Next week, Slot B): 1 game
	finalMatch := structs.Match{
		Model:                  gorm.Model{ID: finalID},
		MatchName:              conferenceName + " Tournament Championship",
		WeekID:                 ts.CollegeWeekID + 1,
		Week:                   uint(ts.CollegeWeek + 1),
		SeasonID:               ts.SeasonID,
		MatchOfWeek:            "B",
		NextGameID:             0,
		NextGameHOA:            "",
		IsConference:           true,
		IsConferenceTournament: true,
		IsPlayoffGame:          false,
		IsNeutralSite:          false,
	}

	matches = append(matches, firstRound, r2g1, r2g2, r3g1, r3g2, qf1, qf2, sf1, sf2, finalMatch)

	return matches
}
