package managers

import (
	"log"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/repository"
	"github.com/CalebRose/SimNBA/structs"
	"gorm.io/gorm"
)

func GetConferenceStandingsByConferenceID(id string, seasonID string) []structs.CollegeStandings {
	db := dbprovider.GetInstance().GetDB()

	var standings []structs.CollegeStandings

	db.Where("conference_id = ? AND season_id = ?", id, seasonID).Order("conference_losses asc").Order("conference_wins desc").
		Order("total_losses asc").Order("total_wins desc").Find(&standings)

	return standings
}

func GetNBAConferenceStandingsByConferenceID(id string, seasonID string) []structs.NBAStandings {
	db := dbprovider.GetInstance().GetDB()

	var standings []structs.NBAStandings

	db.Where("conference_id = ? AND season_id = ?", id, seasonID).
		Order("total_losses asc").Order("total_wins desc").Find(&standings)

	return standings
}

func GetAllConferenceStandingsBySeasonID(seasonID string) []structs.CollegeStandings {
	db := dbprovider.GetInstance().GetDB()

	var standings []structs.CollegeStandings

	db.Where("season_id = ?", seasonID).Order("conference_id asc").Order("conference_losses asc").Order("conference_wins desc").
		Order("total_losses asc").Order("total_wins desc").Find(&standings)

	return standings
}

func GetAllNBAConferenceStandingsBySeasonID(seasonID string) []structs.NBAStandings {
	db := dbprovider.GetInstance().GetDB()

	var standings []structs.NBAStandings

	db.Where("season_id = ?", seasonID).Order("conference_id asc").
		Order("total_losses asc").Order("total_wins desc").Find(&standings)

	return standings
}

func GetNBAStandingsHistoryByTeamID(id string) []structs.NBAStandings {
	db := dbprovider.GetInstance().GetDB()

	var standings []structs.NBAStandings

	db.Where("team_id = ?", id).Find(&standings)

	return standings
}

func GetStandingsHistoryByTeamID(id string) []structs.CollegeStandings {
	db := dbprovider.GetInstance().GetDB()

	var standings []structs.CollegeStandings

	db.Where("team_id = ?", id).Find(&standings)

	return standings
}

func GetStandingsRecordByTeamID(id string, seasonID string) structs.CollegeStandings {
	db := dbprovider.GetInstance().GetDB()

	var standing structs.CollegeStandings

	db.Where("team_id = ? AND season_id = ?", id, seasonID).Find(&standing)

	return standing
}

func GetNBAStandingsRecordByTeamID(id string, seasonID string) structs.NBAStandings {
	db := dbprovider.GetInstance().GetDB()

	var standing structs.NBAStandings

	db.Where("team_id = ? AND season_id = ?", id, seasonID).Find(&standing)

	return standing
}

func GetNBAStandingsBySeasonID(seasonID string) []structs.NBAStandings {
	var standings []structs.NBAStandings
	db := dbprovider.GetInstance().GetDB()
	err := db.Where("season_id = ?", seasonID).Order("total_losses asc").Order("total_wins desc").
		Find(&standings).Error
	if err != nil {
		log.Fatal(err)
	}
	return standings
}

func UpdateStandings(ts structs.Timestamp, MatchType string) {
	db := dbprovider.GetInstance().GetDB()

	if !ts.IsOffSeason {
		games := GetMatchesByWeekIdAndMatchType(strconv.Itoa(int(ts.CollegeWeekID)), strconv.Itoa(int(ts.SeasonID)), MatchType)
		teamMap := GetCollegeTeamMap()
		gameIDs := []uint{}
		for i := 0; i < len(games); i++ {
			game := games[i]
			if !game.GameComplete {
				continue
			}
			gameIDs = append(gameIDs, game.ID)
			HomeID := game.HomeTeamID
			AwayID := game.AwayTeamID
			homeID := strconv.Itoa(int(HomeID))
			awayID := strconv.Itoa(int(AwayID))
			seasonID := strconv.Itoa(int(ts.SeasonID))

			homeStandings := GetStandingsRecordByTeamID(homeID, seasonID)
			awayStandings := GetStandingsRecordByTeamID(awayID, seasonID)

			homeStandings.UpdateCollegeStandings(game)
			awayStandings.UpdateCollegeStandings(game)

			err := db.Save(&homeStandings).Error
			if err != nil {
				log.Panicln("Could not save standings for team " + homeID)
			}

			err = db.Save(&awayStandings).Error
			if err != nil {
				log.Panicln("Could not save standings for team " + awayID)
			}

			if game.NextGameID > 0 {

				nextGameID := strconv.Itoa(int(game.NextGameID))
				winningTeamID := 0
				winningTeam := ""
				winningCoach := ""
				winningTeamRank := 0
				arena := ""
				city := ""
				state := ""
				if game.HomeTeamWin {
					homeTeam := teamMap[HomeID]
					winningTeamID = int(game.HomeTeamID)
					winningTeam = game.HomeTeam
					winningTeamRank = int(game.HomeTeamRank)
					winningCoach = game.HomeTeamCoach
					arena = homeTeam.Arena
					city = homeTeam.City
					state = homeTeam.State
				} else {
					winningTeamID = int(game.AwayTeamID)
					winningTeam = game.AwayTeam
					winningTeamRank = int(game.AwayTeamRank)
					winningCoach = game.AwayTeamCoach
					awayTeam := teamMap[AwayID]
					arena = awayTeam.Arena
					city = awayTeam.City
					state = awayTeam.State
				}

				nextGame := GetMatchByMatchId(nextGameID)

				nextGame.AddTeam(game.NextGameHOA == "H", uint(winningTeamID), uint(winningTeamRank),
					winningTeam, winningCoach, arena, city, state)

				db.Save(&nextGame)
			}

			if game.IsNationalChampionship {
				ts.EndTheCollegeSeason()
				db.Save(&ts)
			}

			// if games[i].HomeTeamCoach != "AI" {
			// 	homeCoach := GetCollegeCoachByCoachName(games[i].HomeTeamCoach)
			// 	homeCoach.UpdateCoachRecord(games[i])

			// 	err = db.Save(&homeCoach).Error
			// 	if err != nil {
			// 		log.Panicln("Could not save coach record for team " + strconv.Itoa(HomeID))
			// 	}
			// }

			// if games[i].AwayTeamCoach != "AI" {
			// 	awayCoach := GetCollegeCoachByCoachName(games[i].AwayTeamCoach)
			// 	awayCoach.UpdateCoachRecord(games[i])
			// 	err = db.Save(&awayCoach).Error
			// 	if err != nil {
			// 		log.Panicln("Could not save coach record for team " + strconv.Itoa(AwayID))
			// 	}
			// }
		}

		db.Model(&structs.Match{}).Where("id IN ?", gameIDs).Update("is_revealed", true)
	}

	if !ts.IsNBAOffseason {
		nbaGames := GetNBATeamMatchesByMatchType(strconv.Itoa(int(ts.NBAWeekID)), strconv.Itoa(int(ts.SeasonID)), MatchType)
		nbaTeamMap := GetProfessionalTeamMap()
		gameIDs := []uint{}
		for _, game := range nbaGames {
			if !game.GameComplete || game.IsPreseason {
				continue
			}
			gameIDs = append(gameIDs, game.ID)
			HomeID := game.HomeTeamID
			AwayID := game.AwayTeamID

			if !game.IsPlayoffGame {
				homeStandings := GetNBAStandingsRecordByTeamID(strconv.Itoa(int(HomeID)), strconv.Itoa(int(ts.SeasonID)))
				awayStandings := GetNBAStandingsRecordByTeamID(strconv.Itoa(int(AwayID)), strconv.Itoa(int(ts.SeasonID)))

				homeStandings.UpdateNBAStandings(game)
				awayStandings.UpdateNBAStandings(game)

				err := db.Save(&homeStandings).Error
				if err != nil {
					log.Panicln("Could not save standings for team " + strconv.Itoa(int(HomeID)))
				}

				err = db.Save(&awayStandings).Error
				if err != nil {
					log.Panicln("Could not save standings for team " + strconv.Itoa(int(AwayID)))
				}
			}

			if game.IsPlayoffGame && game.SeriesID > 0 {
				seriesID := strconv.Itoa(int(game.SeriesID))
				series := GetNBASeriesBySeriesID(seriesID)
				winningID := 0
				if game.HomeTeamWin {
					winningID = int(game.HomeTeamID)
				} else {
					winningID = int(game.AwayTeamID)
				}
				series.UpdateWinCount(winningID)

				if series.GameCount <= 7 && (series.HomeTeamWins < 4 && series.AwayTeamWins < 4) {
					homeTeamID := 0
					nextHomeTeam := ""
					nextHomeTeamCoach := ""
					nextHomeRank := 0
					awayTeamID := 0
					nextAwayTeam := ""
					nextAwayTeamCoach := ""
					nextAwayRank := 0
					city := ""
					arena := ""
					state := ""
					country := ""
					switch series.GameCount {
					case 1, 2, 5, 7:
						homeTeam := nbaTeamMap[series.HomeTeamID]
						homeTeamID = int(series.HomeTeamID)
						nextHomeTeam = series.HomeTeam
						nextHomeTeamCoach = series.HomeTeamCoach
						nextHomeRank = int(series.HomeTeamRank)
						city = homeTeam.City
						arena = homeTeam.Arena
						state = homeTeam.State
						country = homeTeam.Country
						awayTeamID = int(series.AwayTeamID)
						nextAwayTeam = series.AwayTeam
						nextAwayTeamCoach = series.AwayTeamCoach
						nextAwayRank = int(series.AwayTeamRank)
					case 3, 4, 6:
						awayTeam := nbaTeamMap[series.AwayTeamID]
						homeTeamID = int(series.AwayTeamID)
						nextHomeTeam = series.AwayTeam
						nextHomeTeamCoach = series.AwayTeamCoach
						nextHomeRank = int(series.AwayTeamRank)
						city = awayTeam.City
						arena = awayTeam.Arena
						state = awayTeam.State
						country = awayTeam.Country
						awayTeamID = int(series.HomeTeamID)
						nextAwayTeam = series.HomeTeam
						nextAwayTeamCoach = series.HomeTeamCoach
						nextAwayRank = int(series.HomeTeamRank)
					}
					weekID := ts.NBAWeekID
					week := ts.NBAWeek
					matchOfWeek := "A"
					switch game.MatchOfWeek {
					case "A":
						matchOfWeek = "B"
					case "B":
						matchOfWeek = "C"
					case "C":
						matchOfWeek = "D"
					case "D":
						// Move game to next week
						weekID += 1
						week += 1
					}
					matchTitle := series.SeriesName + ": " + nextHomeTeam + " vs. " + nextAwayTeam
					nextGame := structs.NBAMatch{
						WeekID:          weekID,
						Week:            uint(week),
						SeasonID:        ts.SeasonID,
						SeriesID:        series.ID,
						MatchOfWeek:     matchOfWeek,
						MatchName:       matchTitle,
						HomeTeamID:      uint(homeTeamID),
						HomeTeam:        nextHomeTeam,
						HomeTeamCoach:   nextHomeTeamCoach,
						HomeTeamRank:    uint(nextHomeRank),
						AwayTeamID:      uint(awayTeamID),
						AwayTeam:        nextAwayTeam,
						AwayTeamCoach:   nextAwayTeamCoach,
						AwayTeamRank:    uint(nextAwayRank),
						City:            city,
						Arena:           arena,
						State:           state,
						Country:         country,
						IsPlayoffGame:   true,
						IsInternational: series.IsInternational,
					}

					db.Create(&nextGame)
				} else {
					if !series.IsTheFinals && series.NextSeriesID > 0 {
						// Promote Team to Next Series
						nextSeriesID := strconv.Itoa(int(series.NextSeriesID))
						nextSeriesHoa := series.NextSeriesHOA
						nextSeries := GetNBASeriesBySeriesID(nextSeriesID)
						var teamID uint = 0
						teamLabel := ""
						teamCoach := ""
						teamRank := 0
						if series.HomeTeamWin {
							teamID = series.HomeTeamID
							teamLabel = series.HomeTeam
							teamCoach = series.HomeTeamCoach
							teamRank = int(series.HomeTeamRank)
						} else {
							teamID = series.AwayTeamID
							teamLabel = series.AwayTeam
							teamCoach = series.AwayTeamCoach
							teamRank = int(series.AwayTeamRank)
						}
						nextSeries.AddTeam(nextSeriesHoa == "H", teamID, uint(teamRank), teamLabel, teamCoach)
						db.Save(&nextSeries)
					} else {
					}
				}
				db.Save(&series)
			}
		}

		db.Model(&structs.NBAMatch{}).Where("id IN ?", gameIDs).Update("is_revealed", true)
	}

}

func RegressStandings(ts structs.Timestamp, MatchType string) {
	db := dbprovider.GetInstance().GetDB()

	games := GetMatchesByWeekIdAndMatchType(strconv.Itoa(int(ts.CollegeWeekID)), strconv.Itoa(int(ts.SeasonID)), MatchType)

	for i := 0; i < len(games); i++ {
		HomeID := games[i].HomeTeamID
		AwayID := games[i].AwayTeamID

		homeStandings := GetStandingsRecordByTeamID(strconv.Itoa(int(HomeID)), strconv.Itoa(int(ts.SeasonID)))
		awayStandings := GetStandingsRecordByTeamID(strconv.Itoa(int(AwayID)), strconv.Itoa(int(ts.SeasonID)))

		homeStandings.RegressCollegeStandings(games[i])
		awayStandings.RegressCollegeStandings(games[i])

		err := db.Save(&homeStandings).Error
		if err != nil {
			log.Panicln("Could not save standings for team " + strconv.Itoa(int(HomeID)))
		}

		err = db.Save(&awayStandings).Error
		if err != nil {
			log.Panicln("Could not save standings for team " + strconv.Itoa(int(AwayID)))
		}
	}
}

func ResetStandings() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.SeasonID))
	standings := GetAllConferenceStandingsBySeasonID(seasonID)

	// College Standings
	for _, s := range standings {
		s.ResetStandings()
		teamID := strconv.Itoa(int(s.TeamID))
		matches := GetMatchesByTeamIdAndSeasonId(teamID, seasonID)

		for _, m := range matches {
			if !m.GameComplete {
				break
			}

			s.UpdateCollegeStandings(m)
		}

		repository.SaveCollegeStandingsRecord(s, db)
	}

	nbaStandings := GetAllNBAConferenceStandingsBySeasonID(seasonID)
	for _, s := range nbaStandings {
		s.ResetStandings()
		teamID := strconv.Itoa(int(s.TeamID))

		matches := GetNBATeamMatchesBySeasonId(seasonID, teamID)

		for _, m := range matches {
			if !m.GameComplete {
				break
			}
			s.UpdateNBAStandings(m)
		}

		repository.SaveNBAStandingsRecord(s, db)
	}
}

func ResetCollegeStandingsRanks() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.SeasonID))

	db.Model(&structs.CollegeStandings{}).Where("season_id = ?", seasonID).Update("rank", 0)
}

func GetCollegeStandingsMap(seasonID string) map[uint]structs.CollegeStandings {
	standingsMap := make(map[uint]structs.CollegeStandings)

	standings := GetAllConferenceStandingsBySeasonID(seasonID)
	for _, stat := range standings {
		standingsMap[stat.TeamID] = stat
	}

	return standingsMap
}

func ProgressStandings() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.SeasonID - 1))
	teams := GetAllActiveCollegeTeams()

	teamProfileMap := GetTeamProfileMap()
	standingsMap := GetCollegeStandingsMap(seasonID)
	// Update team profiles for bonus points
	for _, t := range teams {
		id := strconv.Itoa(int(t.ID))
		teamProfile := teamProfileMap[id]
		standings := standingsMap[t.ID]
		bonus := 0

		if standings.PostSeasonStatus == "Sweet 16" || standings.IsConferenceChampion {
			bonus = 1
		} else if standings.PostSeasonStatus == "Elite 8" {
			bonus = 2
		} else if standings.PostSeasonStatus == "Final Four" {
			bonus = 3
		} else if standings.PostSeasonStatus == "National Champion Runner-Up" {
			bonus = 4
		} else if standings.PostSeasonStatus == "National Champion" {
			bonus = 5
		}

		if bonus == 0 && standings.ConferenceLosses < 10 {
			conferenceID := strconv.Itoa(int(t.ConferenceID))
			confStandings := GetConferenceStandingsByConferenceID(conferenceID, seasonID)
			if confStandings[0].TeamID == t.ID {
				bonus = 1
			}
		}

		if standings.InvitationalChampion {
			bonus += 1
		}

		teamProfile.ResetSpentPoints()
		teamProfile.ResetScholarshipCount()
		teamProfile.AssignBonusPoints(bonus)

		repository.SaveTeamRecruitingProfileRecord(*teamProfile, db)
	}
}

func GenerateCollegeStandings() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	teams := GetAllActiveCollegeTeams()
	standingsBatch := []structs.CollegeStandings{}

	for _, t := range teams {
		if !t.IsActive {
			continue
		}

		standings := structs.CollegeStandings{
			TeamID:           t.ID,
			TeamName:         t.Team,
			TeamAbbr:         t.Abbr,
			SeasonID:         ts.SeasonID,
			Season:           ts.Season,
			ConferenceID:     t.ConferenceID,
			ConferenceName:   t.Conference,
			PostSeasonStatus: "None",
			BaseStandings: structs.BaseStandings{
				Coach: t.Coach,
			},
		}

		standingsBatch = append(standingsBatch, standings)
	}

	repository.CreateCollegeStandingsRecordsBatch(standingsBatch, db, 100)
}

func GenerateNBAStandings() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	teams := GetAllActiveNBATeams()

	standingsBatch := []structs.NBAStandings{}

	for _, t := range teams {
		if !t.IsActive {
			continue
		}
		coachName := t.NBACoachName
		if coachName == "AI" || len(coachName) == 0 {
			coachName = t.NBAOwnerName
		}
		label := t.Team + " " + t.Nickname
		strippedLabel := strings.TrimSpace(label)
		standings := structs.NBAStandings{
			TeamID:           t.ID,
			TeamName:         t.Team,
			TeamAbbr:         strippedLabel,
			SeasonID:         ts.SeasonID,
			Season:           ts.Season,
			ConferenceID:     t.ConferenceID,
			ConferenceName:   t.Conference,
			DivisionID:       t.DivisionID,
			DivisionName:     t.Division,
			LeagueID:         t.LeagueID,
			League:           t.League,
			PostSeasonStatus: "None",
			BaseStandings: structs.BaseStandings{
				Coach: coachName,
			},
		}

		standingsBatch = append(standingsBatch, standings)
	}

	repository.CreateNBAStandingsRecordsBatch(standingsBatch, db, 100)
}

func GetHistoricalCBBRecordsByTeamID(TeamID string) structs.TeamRecordResponse {
	tsChn := make(chan structs.Timestamp)

	go func() {
		ts := GetTimestamp()
		tsChn <- ts
	}()

	timestamp := <-tsChn
	close(tsChn)
	historicGames := GetCBBMatchesByTeamId(TeamID)
	standings := GetStandingsHistoryByTeamID(TeamID)
	var ConferenceTournamentChampionships []string
	var sweetSixteens []string
	var eliteEights []string
	var finalFours []string
	var runnerUps []string
	var nationalChampionships []string
	overallWins := 0
	overallLosses := 0
	currentSeasonWins := 0
	currentSeasonLosses := 0
	conferenceTournamentWins := 0
	conferenceTournamentLosses := 0
	playoffWins := 0
	playoffLosses := 0
	nitWins := 0
	nitLosses := 0
	CBIWins := 0
	CBILosses := 0

	for _, s := range standings {
		if s.PostSeasonStatus == "Sweet Sixteen" {
			sweetSixteens = append(sweetSixteens, strconv.Itoa(s.Season))
		}

		if s.PostSeasonStatus == "Elite Eight" {
			eliteEights = append(eliteEights, strconv.Itoa(s.Season))
		}

		if s.PostSeasonStatus == "Final Four" {
			finalFours = append(finalFours, strconv.Itoa(s.Season))
		}

		if s.PostSeasonStatus == "National Champion Runner-Up" {
			runnerUps = append(runnerUps, strconv.Itoa(s.Season))
		}

		if s.PostSeasonStatus == "National Champion" {
			nationalChampionships = append(nationalChampionships, strconv.Itoa(s.Season))
		}

		if s.IsConferenceChampion {
			ConferenceTournamentChampionships = append(ConferenceTournamentChampionships, strconv.Itoa(s.Season))
		}

	}

	for _, game := range historicGames {
		if !game.GameComplete || (game.GameComplete && game.SeasonID == timestamp.SeasonID && game.WeekID == timestamp.CollegeWeekID) {
			continue
		}

		isAway := strconv.Itoa(int(game.AwayTeamID)) == TeamID

		if (isAway && game.AwayTeamWin) || (!isAway && game.HomeTeamWin) {
			overallWins++

			if game.SeasonID == timestamp.SeasonID {
				currentSeasonWins++
			}

			if game.IsConferenceTournament {
				conferenceTournamentWins++
			}

			if game.IsPlayoffGame {
				playoffWins++
			}

			if game.IsNITGame {
				nitWins++
			}

			if game.IsCBIGame {
				CBIWins++
			}

		} else {
			overallLosses++

			if game.SeasonID == timestamp.SeasonID {
				currentSeasonLosses++
			}

			if game.IsConferenceTournament {
				conferenceTournamentLosses++
			}

			if game.IsPlayoffGame {
				playoffLosses++
			}

			if game.IsNITGame {
				nitLosses++
			}

			if game.IsCBIGame {
				CBILosses++
			}
		}
	}

	response := structs.TeamRecordResponse{
		OverallWins:             overallWins,
		OverallLosses:           overallLosses,
		CurrentSeasonWins:       currentSeasonWins,
		CurrentSeasonLosses:     currentSeasonLosses,
		TournamentWins:          conferenceTournamentWins,
		TournamentLosses:        conferenceTournamentLosses,
		PlayoffWins:             playoffWins,
		PlayoffLosses:           playoffLosses,
		NITWins:                 nitWins,
		NITLosses:               nitLosses,
		CBIWins:                 CBIWins,
		CBILosses:               CBILosses,
		ConferenceChampionships: ConferenceTournamentChampionships,
		SweetSixteens:           sweetSixteens,
		EliteEights:             eliteEights,
		FinalFours:              finalFours,
		RunnerUps:               runnerUps,
		NationalChampionships:   nationalChampionships,
	}

	return response
}

func GetHistoricalNBARecordsByTeamID(teamID string) structs.TeamRecordResponse {
	tsChn := make(chan structs.Timestamp)

	go func() {
		ts := GetTimestamp()
		tsChn <- ts
	}()

	timestamp := <-tsChn
	close(tsChn)
	season := strconv.Itoa(int(timestamp.Season))
	historicGames := GetNBAMatchesByTeamId(teamID)
	nbaSeries := GetNBASeriesByTeamID(teamID)
	var ConferenceTournamentChampionships []string
	var firstRound []string
	var conferenceSemifinals []string
	var conferenceFinals []string
	var runnerUps []string
	var nationalChampionships []string
	overallWins := 0
	overallLosses := 0
	currentSeasonWins := 0
	currentSeasonLosses := 0
	playoffWins := 0
	playoffLosses := 0

	for _, s := range nbaSeries {
		homeTeamID := strconv.Itoa(int(s.HomeTeamID))
		if s.SeriesName == "First Round" && (s.HomeTeamWin && homeTeamID != teamID) {
			firstRound = append(firstRound, season)
		}

		if s.SeriesName == "Conference Semifinals" && (s.HomeTeamWin && homeTeamID != teamID) {
			conferenceSemifinals = append(conferenceSemifinals, season)
		}

		if s.SeriesName == "Conference Finals" && (s.HomeTeamWin && homeTeamID != teamID) {
			conferenceFinals = append(conferenceFinals, season)
		}

		if s.SeriesName == "The Finals" && (s.HomeTeamWin && homeTeamID != teamID) {
			runnerUps = append(runnerUps, season)
		} else if s.SeriesName == "The Finals" && ((s.HomeTeamWin && homeTeamID == teamID) || (s.AwayTeamWin && homeTeamID != teamID)) {
			nationalChampionships = append(nationalChampionships, season)
		}

		if s.SeriesName == "ISL Finals" && (s.HomeTeamWin && homeTeamID != teamID) {
			runnerUps = append(runnerUps, season)
		} else if s.SeriesName == "ISL Finals" && ((s.HomeTeamWin && homeTeamID == teamID) || (s.AwayTeamWin && homeTeamID != teamID)) {
			nationalChampionships = append(nationalChampionships, season)
		}
	}

	for _, game := range historicGames {
		if !game.GameComplete || (game.GameComplete && game.SeasonID == timestamp.SeasonID && game.WeekID == timestamp.CollegeWeekID) {
			continue
		}

		isAway := strconv.Itoa(int(game.AwayTeamID)) == teamID

		if (isAway && game.AwayTeamWin) || (!isAway && game.HomeTeamWin) {
			overallWins++

			if game.SeasonID == timestamp.SeasonID {
				currentSeasonWins++
			}

			if game.IsPlayoffGame {
				playoffWins++
			}

		} else {
			overallLosses++

			if game.SeasonID == timestamp.SeasonID {
				currentSeasonLosses++
			}

			if game.IsPlayoffGame {
				playoffLosses++
			}
		}
	}

	response := structs.TeamRecordResponse{
		OverallWins:             overallWins,
		OverallLosses:           overallLosses,
		CurrentSeasonWins:       currentSeasonWins,
		CurrentSeasonLosses:     currentSeasonLosses,
		PlayoffWins:             playoffWins,
		PlayoffLosses:           playoffLosses,
		ConferenceChampionships: ConferenceTournamentChampionships,
		RunnerUps:               runnerUps,
		NationalChampionships:   nationalChampionships,
	}

	return response
}

// UpdateCollegeRankings recalculates RPI, SOS, SOR, KenPom (from team season
// stats), quadrant records and the composite rank for the current season.
func UpdateCollegeRankings(ts structs.Timestamp) {
	GenerateCollegeRankings(strconv.Itoa(int(ts.SeasonID)))
}

func GenerateCollegeRankings(seasonID string) {
	db := dbprovider.GetInstance().GetDB()
	collegeTeams := GetAllActiveCollegeTeams()
	collegeTeamMap := MakeCollegeTeamMap(collegeTeams)
	collegeStandings := repository.FindAllCollegeStandingsRecords(repository.StandingsQuery{SeasonID: seasonID})
	collegeGames := repository.FindCollegeMatchRecords(repository.GameQuery{
		SeasonID: seasonID,
	})

	preSeasonCFBSeasonStats := repository.FindCollegeTeamSeasonStatsRecords(seasonID, "1")
	regularSeasonCFBSeasonStats := repository.FindCollegeTeamSeasonStatsRecords(seasonID, "2")
	postSeasonCFBSeasonStats := repository.FindCollegeTeamSeasonStatsRecords(seasonID, "3")

	_ = collegeTeamMap
	_ = preSeasonCFBSeasonStats
	_ = regularSeasonCFBSeasonStats
	_ = postSeasonCFBSeasonStats

	seasonStatsMap := make(map[uint][]structs.TeamSeasonStats)
	for _, stats := range preSeasonCFBSeasonStats {
		seasonStatsMap[stats.TeamID] = append(seasonStatsMap[stats.TeamID], stats)
	}
	for _, stats := range regularSeasonCFBSeasonStats {
		seasonStatsMap[stats.TeamID] = append(seasonStatsMap[stats.TeamID], stats)
	}
	for _, stats := range postSeasonCFBSeasonStats {
		seasonStatsMap[stats.TeamID] = append(seasonStatsMap[stats.TeamID], stats)
	}

	// Build lookup structures
	standingsMap := MakeCollegeStandingsMap(collegeStandings)
	gamesMap := MakeCollegeMatchMapByTeamID(collegeGames)

	// calcWinPct computes win percentage on the fly, since TotalWinPercentage
	// may not be populated in the DB yet.
	calcWinPct := func(wins, losses int) float64 {
		if wins+losses == 0 {
			return 0
		}
		return float64(wins) / float64(wins+losses)
	}

	// -------------------------------------------------------------------------
	// Pass 1: adjusted SOS, preliminary SOR, RPI, ConferenceStrengthAdj.
	//
	// Adjusted SOS: each opponent's win% is computed excluding the game they
	// played against the team being evaluated. This prevents a 12-0 team and
	// a 0-12 team from getting different SOS values just because of their own
	// result — if both faced the same 12 opponents, they should see the same SOS.
	// -------------------------------------------------------------------------
	for idx, standings := range collegeStandings {
		teamID := uint(standings.TeamID)
		teamGames := gamesMap[teamID]

		var adjOppWinPcts []float64
		var oppOppWinPcts []float64
		var adjConfOppWinPcts []float64

		sorNumerator := 0.0
		sorDenominator := 0.0

		for _, game := range teamGames {
			if !game.GameComplete {
				continue
			}

			isAway := game.AwayTeamID == teamID
			var oppID uint
			if isAway {
				oppID = game.HomeTeamID
			} else {
				oppID = game.AwayTeamID
			}

			oppStandings, ok := standingsMap[oppID]
			if !ok {
				continue
			}

			isWinner := (!isAway && game.HomeTeamWin) || (isAway && game.AwayTeamWin)

			// Adjusted opponent win%: remove the head-to-head game from O's record.
			// If O beat T (!isWinner), subtract one win; always subtract one game.
			adjOppWins := oppStandings.TotalWins
			if !isWinner {
				adjOppWins--
			}
			adjOppGames := oppStandings.TotalWins + oppStandings.TotalLosses - 1
			adjOppWinPct := 0.0
			if adjOppGames > 0 && adjOppWins >= 0 {
				adjOppWinPct = float64(adjOppWins) / float64(adjOppGames)
			}

			adjOppWinPcts = append(adjOppWinPcts, adjOppWinPct)

			if game.IsConference {
				adjConfOppWinPcts = append(adjConfOppWinPcts, adjOppWinPct)
			}

			// Opponents' opponents win% (RPI third component) — unadjusted
			oppGames := gamesMap[oppID]
			sumOppOpp := 0.0
			countOppOpp := 0
			for _, og := range oppGames {
				if !og.GameComplete {
					continue
				}
				isOppAway := og.AwayTeamID == oppID
				var oppOppID uint
				if isOppAway {
					oppOppID = og.HomeTeamID
				} else {
					oppOppID = og.AwayTeamID
				}
				if ooStandings, ok2 := standingsMap[oppOppID]; ok2 {
					sumOppOpp += calcWinPct(ooStandings.TotalWins, ooStandings.TotalLosses)
					countOppOpp++
				}
			}
			if countOppOpp > 0 {
				oppOppWinPcts = append(oppOppWinPcts, sumOppOpp/float64(countOppOpp))
			}

			// Preliminary SOR using adjusted opponent win% (refined in Pass 2)
			sorDenominator++
			if isWinner {
				sorNumerator += 0.5 + 0.5*adjOppWinPct
			}
		}

		// Adjusted SOS
		sos := float32(0)
		if len(adjOppWinPcts) > 0 {
			sum := 0.0
			for _, v := range adjOppWinPcts {
				sum += v
			}
			sos = float32(sum / float64(len(adjOppWinPcts)))
		}

		// Preliminary SOR
		sor := float32(0)
		if sorDenominator > 0 {
			sor = float32(sorNumerator / sorDenominator)
		}

		// RPI: 25% team win% + 50% adjusted avg opp win% (SOS) + 25% avg opp-opp win%
		avgOppOppWinPct := float32(0)
		if len(oppOppWinPcts) > 0 {
			sum := 0.0
			for _, v := range oppOppWinPcts {
				sum += v
			}
			avgOppOppWinPct = float32(sum / float64(len(oppOppWinPcts)))
		}
		teamWinPct := float32(calcWinPct(standings.TotalWins, standings.TotalLosses))
		rpi := 0.25*teamWinPct + 0.50*sos + 0.25*avgOppOppWinPct

		// ConferenceStrengthAdj — fall back to SOS for independent teams (conf IDs 13 and 22)
		// which never have IsConference = true games.
		confStrengthAdj := sos
		if len(adjConfOppWinPcts) > 0 {
			sum := 0.0
			for _, v := range adjConfOppWinPcts {
				sum += v
			}
			confStrengthAdj = float32(sum / float64(len(adjConfOppWinPcts)))
		}

		collegeStandings[idx].SOS = sos
		collegeStandings[idx].SOR = sor
		collegeStandings[idx].RPIRating = rpi
		collegeStandings[idx].ConferenceStrengthAdj = confStrengthAdj
	}

	// -------------------------------------------------------------------------
	// KenPom-style rating from team season stats (regular season + postseason;
	// preseason is excluded).
	//   raw offensive efficiency  = 100 * points / possessions
	//   raw defensive efficiency  = 100 * points allowed / possessions
	// (a team's opponents have ~the same number of possessions as the team).
	// Each team's efficiencies are then iteratively adjusted for the quality of
	// its opponents, exactly like KenPom's AdjO / AdjD:
	//   AdjO = rawO * leagueAvg / avg(opponent AdjD)
	//   AdjD = rawD * leagueAvg / avg(opponent AdjO)
	// KenPomRating = AdjO - AdjD (adjusted efficiency margin per 100 possessions).
	// -------------------------------------------------------------------------
	const homeCourtPoints = 3.5
	const ratingIterations = 100

	type efficiencyTotals struct {
		points, pointsAgainst, possessions float64
	}
	totalsMap := make(map[uint]*efficiencyTotals, len(collegeStandings))
	for _, statSet := range [][]structs.TeamSeasonStats{regularSeasonCFBSeasonStats, postSeasonCFBSeasonStats} {
		for _, st := range statSet {
			t, ok := totalsMap[st.TeamID]
			if !ok {
				t = &efficiencyTotals{}
				totalsMap[st.TeamID] = t
			}
			t.points += float64(st.Points)
			t.pointsAgainst += float64(st.PointsAgainst)
			t.possessions += float64(st.Possessions)
		}
	}

	rawO := make(map[uint]float64, len(totalsMap))
	rawD := make(map[uint]float64, len(totalsMap))
	totalPoints, totalPossessions := 0.0, 0.0
	for id, t := range totalsMap {
		if t.possessions <= 0 {
			continue
		}
		rawO[id] = 100 * t.points / t.possessions
		rawD[id] = 100 * t.pointsAgainst / t.possessions
		totalPoints += t.points
		totalPossessions += t.possessions
	}

	leagueAvgEff := 100.0
	if totalPossessions > 0 {
		leagueAvgEff = 100 * totalPoints / totalPossessions
	}

	// Average possessions per game, used to convert efficiency margin to points
	avgTempo := 70.0
	gamesPlayed := 0.0
	for _, statSet := range [][]structs.TeamSeasonStats{regularSeasonCFBSeasonStats, postSeasonCFBSeasonStats} {
		for _, st := range statSet {
			gamesPlayed += float64(st.GamesPlayed)
		}
	}
	if gamesPlayed > 0 && totalPossessions > 0 {
		avgTempo = totalPossessions / gamesPlayed
	}

	// Opponent lists come from completed games so schedule strength is respected
	opponents := make(map[uint][]uint, len(collegeStandings))
	for _, s := range collegeStandings {
		teamID := uint(s.TeamID)
		for _, game := range gamesMap[teamID] {
			if !game.GameComplete {
				continue
			}
			oppID := game.AwayTeamID
			if game.AwayTeamID == teamID {
				oppID = game.HomeTeamID
			}
			if _, ok := rawO[oppID]; ok {
				opponents[teamID] = append(opponents[teamID], oppID)
			}
		}
	}

	adjO := make(map[uint]float64, len(rawO))
	adjD := make(map[uint]float64, len(rawO))
	for id := range rawO {
		adjO[id] = rawO[id]
		adjD[id] = rawD[id]
	}

	for iter := 0; iter < ratingIterations; iter++ {
		nextO := make(map[uint]float64, len(adjO))
		nextD := make(map[uint]float64, len(adjD))
		for id := range rawO {
			opps := opponents[id]
			if len(opps) == 0 {
				nextO[id], nextD[id] = rawO[id], rawD[id]
				continue
			}
			sumOppD, sumOppO := 0.0, 0.0
			for _, oppID := range opps {
				sumOppD += adjD[oppID]
				sumOppO += adjO[oppID]
			}
			avgOppD := sumOppD / float64(len(opps))
			avgOppO := sumOppO / float64(len(opps))
			nextO[id] = rawO[id] * leagueAvgEff / avgOppD
			nextD[id] = rawD[id] * leagueAvgEff / avgOppO
		}
		// Re-center so league averages stay anchored at leagueAvgEff
		meanO, meanD := 0.0, 0.0
		for id := range nextO {
			meanO += nextO[id]
			meanD += nextD[id]
		}
		meanO /= float64(len(nextO))
		meanD /= float64(len(nextD))
		for id := range nextO {
			adjO[id] = nextO[id] * leagueAvgEff / meanO
			adjD[id] = nextD[id] * leagueAvgEff / meanD
		}
	}

	// Teams with no stats rate as average (0)
	ratings := make(map[uint]float64, len(collegeStandings))
	for _, s := range collegeStandings {
		teamID := uint(s.TeamID)
		if _, ok := adjO[teamID]; ok {
			ratings[teamID] = adjO[teamID] - adjD[teamID]
		}
	}

	totalTeams := len(collegeStandings)

	kenPomOrder := make([]int, totalTeams)
	for i := range kenPomOrder {
		kenPomOrder[i] = i
	}
	sort.SliceStable(kenPomOrder, func(i, j int) bool {
		return ratings[uint(collegeStandings[kenPomOrder[i]].TeamID)] > ratings[uint(collegeStandings[kenPomOrder[j]].TeamID)]
	})
	for pos, idx := range kenPomOrder {
		collegeStandings[idx].KenPomRank = uint16(pos + 1)
		collegeStandings[idx].KenPomRating = float32(ratings[uint(collegeStandings[idx].TeamID)])
	}

	rpiOrder := make([]int, totalTeams)
	for i := range rpiOrder {
		rpiOrder[i] = i
	}
	sort.SliceStable(rpiOrder, func(i, j int) bool {
		return collegeStandings[rpiOrder[i]].RPIRating > collegeStandings[rpiOrder[j]].RPIRating
	})
	for pos, idx := range rpiOrder {
		collegeStandings[idx].RPIRank = uint16(pos + 1)
	}

	// Quadrant rank (NET-like): average of KenPom and RPI ranks, KenPom breaks ties.
	quadOrder := make([]int, totalTeams)
	for i := range quadOrder {
		quadOrder[i] = i
	}
	sort.SliceStable(quadOrder, func(i, j int) bool {
		a, b := collegeStandings[quadOrder[i]], collegeStandings[quadOrder[j]]
		aAvg := int(a.KenPomRank) + int(a.RPIRank)
		bAvg := int(b.KenPomRank) + int(b.RPIRank)
		if aAvg != bAvg {
			return aAvg < bAvg
		}
		return a.KenPomRank < b.KenPomRank
	})
	quadRankMap := make(map[uint]int, totalTeams)
	for pos, idx := range quadOrder {
		quadRankMap[uint(collegeStandings[idx].TeamID)] = pos + 1
	}

	// Reference team for SOR: the rating of a typical top-25 team.
	referenceRating := 0.0
	if totalTeams > 0 {
		refPos := 24
		if refPos >= totalTeams {
			refPos = totalTeams - 1
		}
		referenceRating = float64(collegeStandings[kenPomOrder[refPos]].KenPomRating)
	}

	// -------------------------------------------------------------------------
	// Pass 2: SOR and quadrant records.
	//
	// SOR: sum of (result - P(reference top-25 team beats this opponent at this
	// location)) averaged per game. Beating tough opponents earns more; losing
	// to weak ones costs more. Win probability is logistic on KenPom rating diff.
	//
	// Quadrants (NCAA NET style, by opponent quadrant rank and location):
	//   Q1: home 1-30,   neutral 1-50,   away 1-75
	//   Q2: home 31-75,  neutral 51-100, away 76-135
	//   Q3: home 76-160, neutral 101-200, away 136-240
	//   Q4: everything below
	// -------------------------------------------------------------------------
	quadrantFor := func(oppRank int, isAway, isNeutral bool) int {
		var limits [3]int
		switch {
		case isNeutral:
			limits = [3]int{50, 100, 200}
		case isAway:
			limits = [3]int{75, 135, 240}
		default:
			limits = [3]int{30, 75, 160}
		}
		for q, limit := range limits {
			if oppRank <= limit {
				return q + 1
			}
		}
		return 4
	}

	quadrantWeights := [4]float64{1.0, 0.5, 0.25, 0.0}
	quadrantLossWeights := [4]float64{0.0, -0.1, -0.5, -1.0}

	for idx, standings := range collegeStandings {
		teamID := uint(standings.TeamID)

		var wins, losses [4]uint8
		sorSum := 0.0
		sorGames := 0
		quadRating := 0.0

		for _, game := range gamesMap[teamID] {
			if !game.GameComplete {
				continue
			}

			isAway := game.AwayTeamID == teamID
			oppID := game.AwayTeamID
			if isAway {
				oppID = game.HomeTeamID
			}
			oppRating, hasRating := ratings[oppID]
			oppQuadRank, hasQuadRank := quadRankMap[oppID]
			if !hasRating || !hasQuadRank {
				continue
			}

			isWinner := (!isAway && game.HomeTeamWin) || (isAway && game.AwayTeamWin)

			// Reference team's expected margin (points) against this opponent
			expectedMargin := (referenceRating - oppRating) * avgTempo / 100
			if !game.IsNeutralSite {
				if isAway {
					expectedMargin -= homeCourtPoints
				} else {
					expectedMargin += homeCourtPoints
				}
			}
			winProb := 1.0 / (1.0 + math.Exp(-expectedMargin/10.0))
			result := 0.0
			if isWinner {
				result = 1.0
			}
			sorSum += result - winProb
			sorGames++

			q := quadrantFor(oppQuadRank, isAway, game.IsNeutralSite) - 1
			if isWinner {
				wins[q]++
				quadRating += quadrantWeights[q]
			} else {
				losses[q]++
				quadRating += quadrantLossWeights[q]
			}
		}

		sor := float32(0)
		if sorGames > 0 {
			sor = float32(sorSum / float64(sorGames))
		}

		collegeStandings[idx].SOR = sor
		collegeStandings[idx].Q1Wins, collegeStandings[idx].Q2Wins = wins[0], wins[1]
		collegeStandings[idx].Q3Wins, collegeStandings[idx].Q4Wins = wins[2], wins[3]
		collegeStandings[idx].Q1Losses, collegeStandings[idx].Q2Losses = losses[0], losses[1]
		collegeStandings[idx].Q3Losses, collegeStandings[idx].Q4Losses = losses[2], losses[3]
		collegeStandings[idx].QuadrantRating = float32(quadRating)
	}

	// -------------------------------------------------------------------------
	// Pass 3: compute full composite ToucanScore and assign final ToucanRank.
	// -------------------------------------------------------------------------
	type teamScore struct {
		idx   int
		score float64
	}

	scores := make([]teamScore, len(collegeStandings))
	for idx, s := range collegeStandings {

		winPct := calcWinPct(s.TotalWins, s.TotalLosses)
		games := float64(s.TotalWins + s.TotalLosses)

		preseasonBonus := 0.0
		if s.PreseasonRank > 0 && s.PreseasonRank <= 25 {
			preseasonBonus = float64(26-int(s.PreseasonRank)) / float64(25) * 0.02
		}

		regularSeasonBonus := 0.0
		if s.Rank > 0 && s.Rank <= 25 {
			regularSeasonBonus = float64(26-int(s.Rank)) / float64(25) * 0.03
		}

		// Rank-based components scaled to 0..1 (1 = best)
		kenPomScore := 0.0
		rpiRankScore := 0.0
		if totalTeams > 1 {
			kenPomScore = 1.0 - float64(int(s.KenPomRank)-1)/float64(totalTeams-1)
			rpiRankScore = 1.0 - float64(int(s.RPIRank)-1)/float64(totalTeams-1)
		}

		quadPerGame := 0.0
		if games > 0 {
			quadPerGame = float64(s.QuadrantRating) / games
		}

		// Win%: 20%, KenPom: 20%, RPI rating: 10%, RPI rank: 5%, SOR: 15% (scaled
		// to roughly 0..1 around 0.5), SOS: 5%, quadrant rating per game: 15%,
		// conference strength: 5%, preseason: up to 2%, regular season rank: up to 3%
		score := winPct*0.20 +
			kenPomScore*0.20 +
			float64(s.RPIRating)*0.10 +
			rpiRankScore*0.05 +
			(float64(s.SOR)+0.5)*0.15 +
			float64(s.SOS)*0.05 +
			quadPerGame*0.15 +
			float64(s.ConferenceStrengthAdj)*0.05 +
			preseasonBonus +
			regularSeasonBonus

		scores[idx] = teamScore{idx: idx, score: score}
	}

	sort.SliceStable(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})

	rank := uint16(1)
	for _, sc := range scores {
		collegeStandings[sc.idx].ToucanRank = rank
		rank++
	}

	// Persist all updated standings
	for idx := range collegeStandings {
		collegeStandings[idx].CalculatePercentages()
		repository.SaveCollegeStandingsRecord(collegeStandings[idx], db)
	}
}

func CreatePreseasonRanking() {
	db := dbprovider.GetInstance().GetDB()
	timestamp := GetTimestamp()
	currentSeasonID := strconv.Itoa(int(timestamp.SeasonID))
	previousSeasonID := strconv.Itoa(int(timestamp.SeasonID - 1))

	teamRecruitingProfiles := GetTeamRecruitingProfilesForRecruitSync()

	teamRecruitingProfileMap := MakeTeamRecruitingProfileMapByTeamID(teamRecruitingProfiles)

	// Get all teams for current season
	teams := GetAllActiveCollegeTeams()

	// preSeasonCFBSeasonStats := repository.FindCollegeTeamSeasonStatsRecords(seasonID, "1")
	// regularSeasonCFBSeasonStats := repository.FindCollegeTeamSeasonStatsRecords(seasonID, "2")
	// postSeasonCFBSeasonStats := repository.FindCollegeTeamSeasonStatsRecords(seasonID, "3")
	// Get previous season data for baseline
	previousStandings := repository.FindAllCollegeStandingsRecords(repository.StandingsQuery{SeasonID: previousSeasonID})

	// Get current season standings (should be empty/reset for preseason)
	currentStandings := repository.FindAllCollegeStandingsRecords(repository.StandingsQuery{SeasonID: currentSeasonID})

	// currentStandingsMap := MakeCollegeStandingsMap(currentStandings)

	previousStandingsMap := MakeCollegeStandingsMap(previousStandings)

	// TODO: Implement these data sources as needed:
	// 1. Get roster data with player ratings
	players := GetAllCollegePlayers()

	collegeRosterMap := MakeCollegePlayerMapByTeamID(players, true)

	preseasonScores := calculatePreseasonScores(teams, previousStandingsMap, collegeRosterMap, teamRecruitingProfileMap)

	assignPreseasonRanks(preseasonScores, &currentStandings, db)
}

// calculatePreseasonScores computes preseason ranking scores based on multiple factors
func calculatePreseasonScores(teams []structs.Team, previousStandings map[uint]structs.CollegeStandings, collegeRosterMap map[uint][]structs.CollegePlayer, teamRecruitingProfileMap map[uint]structs.TeamRecruitingProfile) map[uint]float32 {
	scores := make(map[uint]float32)

	for _, team := range teams {
		var score float32 = 0.0

		teamProfile := teamRecruitingProfileMap[team.ID]

		// Factor 1: Previous Season Performance (40% weight)
		if prevStanding, exists := previousStandings[team.ID]; exists {
			// Base score on previous win percentage and RPI
			winPct := prevStanding.GetWinPercentage()
			rpiScore := prevStanding.RPIRating

			prevSeasonScore := (winPct * 0.6) + (rpiScore * 0.4)
			score += prevSeasonScore * 0.40
		} else {
			// New program or missing data - use neutral score
			score += 0.5 * 0.40
		}

		// Factor 2: Program Prestige (20% weight)
		prestigeScore := float32(teamProfile.ProgramPrestige) / 10.0
		score += prestigeScore * 0.20

		// Factor 3: Roster Talent (30% weight)
		roster := collegeRosterMap[team.ID]
		rosterScore := calculateRosterTalent(roster)
		score += rosterScore * 0.30

		// Factor 4: Coaching Stability (10% weight)
		// CoachRating is a number between 1 and 10
		coachingScore := float32(teamProfile.CoachRating) / 10.0
		score += coachingScore * 0.10

		scores[team.ID] = score
	}

	return scores
}

// calculateRosterTalent evaluates the overall talent level of a team's roster
func calculateRosterTalent(roster []structs.CollegePlayer) float32 {
	if len(roster) == 0 {
		return 0.3 // Below average for teams with no roster data
	}

	var totalTalent float32
	var playerCount int
	var starterBonus float32

	// Position weights for roster balance evaluation
	positionCounts := make(map[string]int)

	for _, player := range roster {
		// Use player's Overall rating as primary talent metric
		// Assuming Overall is 1-50, normalize to 0-1
		playerTalent := float32(player.Overall) / 50.0

		// Weight by player year (experience matters)
		yearMultiplier := float32(1.0)
		switch player.Year {
		case 1: // Freshman
			yearMultiplier = 0.8
		case 2: // Sophomore
			yearMultiplier = 0.9
		case 3: // Junior
			yearMultiplier = 1.0
		case 4: // Senior
			yearMultiplier = 1.1
		case 5: // Graduate/5th year
			yearMultiplier = 1.15
		}

		// Apply experience multiplier
		adjustedTalent := playerTalent * yearMultiplier

		// Bonus for star players (45+ overall)
		if player.Overall >= 45 {
			starterBonus += 0.05 // Each elite player adds 5% bonus
		}

		totalTalent += adjustedTalent
		playerCount++

		// Track position balance
		positionCounts[player.Position]++
	}

	// Calculate average roster talent
	avgTalent := totalTalent / float32(playerCount)

	// Roster depth bonus/penalty
	depthModifier := float32(1.0)
	if playerCount < 8 { // Thin roster penalty
		depthModifier = 0.85
	} else if playerCount > 11 { // Good depth bonus
		depthModifier = 1.1
	}

	// Position balance bonus (ensure reasonable distribution)
	balanceModifier := calculatePositionBalance(positionCounts)

	// Final roster score with all modifiers
	finalScore := (avgTalent + starterBonus) * depthModifier * balanceModifier

	// Ensure score stays within reasonable bounds (0.0 to 1.0)
	if finalScore > 1.0 {
		finalScore = 1.0
	}
	if finalScore < 0.0 {
		finalScore = 0.0
	}

	return finalScore
}

// calculatePositionBalance evaluates roster balance across positions
func calculatePositionBalance(positionCounts map[string]int) float32 {
	// Expected minimum players by position for hockey
	expectedMins := map[string]int{
		"C": 2, // Centers
		"F": 2, // Forwards
		"G": 2, // Guards
	}

	balanceScore := float32(1.0)

	for position, expectedMin := range expectedMins {
		actual := positionCounts[position]
		if actual < expectedMin {
			// Penalty for being short at a position
			shortfall := float32(expectedMin - actual)
			balanceScore -= shortfall * 0.05 // 5% penalty per missing player
		}
	}

	// Ensure balance score doesn't go below 0.7 (max 30% penalty)
	if balanceScore < 0.7 {
		balanceScore = 0.7
	}

	return balanceScore
}

// assignPreseasonRanks assigns preseason rankings based on calculated scores
func assignPreseasonRanks(preseasonScores map[uint]float32, standings *[]structs.CollegeStandings, db *gorm.DB) {
	type teamScore struct {
		teamID uint
		score  float32
	}

	var sortedScores []teamScore
	for teamID, score := range preseasonScores {
		sortedScores = append(sortedScores, teamScore{teamID, score})
	}

	// Sort by score descending
	for i := 0; i < len(sortedScores)-1; i++ {
		for j := i + 1; j < len(sortedScores); j++ {
			if sortedScores[i].score < sortedScores[j].score {
				sortedScores[i], sortedScores[j] = sortedScores[j], sortedScores[i]
			}
		}
	}

	// Assign preseason ranks
	for i := range *standings {
		standing := &(*standings)[i]
		for rank, teamScore := range sortedScores {
			if teamScore.teamID == standing.TeamID {
				// Set preseason ranking - you might want a separate PreseasonRank field
				standing.PreseasonRank = uint16(rank) + 1

				// Initialize RPI with preseason score for display
				standing.RPIRating = teamScore.score

				repository.SaveCollegeStandingsRecord(*standing, db)
				break
			}
		}
	}
}
