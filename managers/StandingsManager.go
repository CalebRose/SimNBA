package managers

import (
	"log"
	"strconv"
	"strings"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/repository"
	"github.com/CalebRose/SimNBA/structs"
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
		for i := 0; i < len(games); i++ {
			game := games[i]
			if !game.GameComplete {
				continue
			}
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
	}

	if !ts.IsNBAOffseason {
		nbaGames := GetNBATeamMatchesByMatchType(strconv.Itoa(int(ts.NBAWeekID)), strconv.Itoa(int(ts.SeasonID)), MatchType)
		nbaTeamMap := GetProfessionalTeamMap()
		for _, game := range nbaGames {
			if !game.GameComplete {
				continue
			}
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
					}
				}
				db.Save(&series)
			}
		}
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

	db.Model(&structs.CollegeStandings{}).Where("season_id = ?", seasonID).Updates(structs.CollegeStandings{Rank: 0})
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

		repository.SaveCBBTeamRecruitingProfile(*teamProfile, db)
	}
}

func GenerateCollegeStandings() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	teams := GetAllActiveCollegeTeams()

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

		db.Create(&standings)
	}
}

func GenerateNBAStandings() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	teams := GetAllActiveNBATeams()

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

		db.Create(&standings)
	}
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
