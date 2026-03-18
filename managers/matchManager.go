package managers

import (
	"fmt"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/repository"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
)

func GetTestMatches(request structs.TestRequest) structs.MatchStateResponse {
	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.SeasonID))

	arenaMap := GetArenaMap()

	collegeTeams := GetAllActiveCollegeTeams()
	collegeTeamAbbrMap := make(map[string]structs.Team)
	for _, team := range collegeTeams {
		collegeTeamAbbrMap[team.Abbr] = team
	}

	collegePlayers := GetAllCollegePlayers()
	collegePlayerMap := MakeCollegePlayerMapByTeamID(collegePlayers, true)
	collegeStandings := repository.FindAllCollegeStandingsRecords(repository.StandingsQuery{SeasonID: seasonID})
	collegeStandingsMap := MakeCollegeStandingsMap(collegeStandings)
	collegeLineups := repository.FindCollegeLineupRecords(repository.GameplanQuery{})
	collegeLineupMap := MakeCollegeLineupMapByTeamID(collegeLineups)
	collegeGameplans := GetAllCollegeGameplans()
	collegeGameplansMap := MakeCollegeGameplanMap(collegeGameplans)

	matchesList := make([]structs.MatchResponse, 0, len(request.TestMatches))

	for idx, m := range request.TestMatches {
		homeTeam := structs.MatchTeamResponse{}
		awayTeam := structs.MatchTeamResponse{}

		ht := collegeTeamAbbrMap[m.HomeTeam]
		at := collegeTeamAbbrMap[m.AwayTeam]
		homeTeam.Map(ht)
		awayTeam.Map(at)

		homeTeamLineup := collegeLineupMap[ht.ID]
		awayTeamLineup := collegeLineupMap[at.ID]

		homeGameplan := collegeGameplansMap[ht.ID]
		awayGameplan := collegeGameplansMap[at.ID]

		htGameLineUp := []structs.GameplanLineup{}
		for _, lp := range homeTeamLineup {
			htGameLineUp = append(htGameLineUp, lp.GameplanLineup)
		}

		atGameLineUp := []structs.GameplanLineup{}
		for _, lp := range awayTeamLineup {
			atGameLineUp = append(atGameLineUp, lp.GameplanLineup)
		}

		capacity := 0
		arena := arenaMap[ht.Arena]
		if arena.ID == 0 {
			capacity = 6000
		} else {
			capacity = int(arena.Capacity)
		}

		currentStandings := collegeStandingsMap[ht.ID]
		attendancePercent := getAttendancePercent(int(currentStandings.TotalWins), int(currentStandings.TotalLosses))
		if ts.CollegeWeek == 0 {
			attendancePercent = 1.0
		}
		fanCount := uint32(float64(capacity) * attendancePercent)
		hra := float64(fanCount) / float64(capacity)
		homeRoster := collegePlayerMap[ht.ID]
		awayRoster := collegePlayerMap[at.ID]

		homeGamePlayerRoster := []structs.GamePlayer{}
		awayGamePlayerRoster := []structs.GamePlayer{}

		for _, p := range homeRoster {
			gamePlayer := structs.GamePlayer{
				ID:         p.ID,
				BasePlayer: p.BasePlayer,
			}
			gamePlayer.CalculateModifiers(true, hra)
			homeGamePlayerRoster = append(homeGamePlayerRoster, gamePlayer)
		}

		for _, p := range awayRoster {
			gamePlayer := structs.GamePlayer{
				ID:         p.ID,
				BasePlayer: p.BasePlayer,
			}
			gamePlayer.CalculateModifiers(false, hra)
			awayGamePlayerRoster = append(awayGamePlayerRoster, gamePlayer)
		}

		livestreamChannel := 0
		if (ht.IsUserCoached) || (at.IsUserCoached) {
			if ht.ConferenceID%2 == 1 {
				livestreamChannel = 1
			} else {
				livestreamChannel = 2
			}
		} else {
			if ht.ConferenceID%2 == 1 {
				livestreamChannel = 3
			} else {
				livestreamChannel = 4
			}
		}

		match := structs.MatchResponse{
			MatchName:     "Test Match" + strconv.Itoa(idx+1),
			ID:            0,
			WeekID:        ts.CollegeWeekID,
			Week:          uint(ts.Season),
			SeasonID:      ts.SeasonID,
			HomeTeamID:    ht.ID,
			HomeTeam:      m.HomeTeam,
			AwayTeamID:    at.ID,
			AwayTeam:      m.AwayTeam,
			Arena:         arena.ArenaName,
			Capacity:      capacity,
			City:          ht.City,
			State:         ht.State,
			IsNeutralSite: m.IsNeutral,
			IsNBAMatch:    false,
			Channel:       uint(livestreamChannel),
			MatchData: structs.MatchDataResponse{
				HomeTeam:         homeTeam,
				HomeTeamRoster:   homeGamePlayerRoster,
				HomeTeamLineup:   htGameLineUp,
				HomeTeamGameplan: homeGameplan,
				AwayTeam:         awayTeam,
				AwayTeamRoster:   awayGamePlayerRoster,
				AwayTeamLineup:   atGameLineUp,
				AwayTeamGameplan: awayGameplan,
			},
		}
		matchesList = append(matchesList, match)
	}

	return structs.MatchStateResponse{
		Matches: matchesList,
		Week:    uint(ts.NBAWeek),
	}
}

func GetMatchesForTimeslot() structs.MatchStateResponse {
	ts := GetTimestamp()
	if !ts.RunGames {
		return structs.MatchStateResponse{
			Matches:   []structs.MatchResponse{},
			MatchType: "",
			Week:      0,
		}
	}
	seasonID := strconv.Itoa(int(ts.SeasonID))
	weekID := strconv.Itoa(int(ts.CollegeWeekID))
	nbaWeekID := strconv.Itoa(int(ts.NBAWeekID))

	matchType := ""

	if !ts.GamesARan {
		matchType = "A"
	} else if !ts.GamesBRan {
		matchType = "B"
	} else if !ts.GamesCRan {
		matchType = "C"
	} else if !ts.GamesDRan {
		matchType = "D"
	}

	// Get College Matches
	collegeMatches := GetMatchesByWeekIdAndMatchType(weekID, seasonID, matchType)
	// Get Professional Matches
	nbaMatches := GetNBATeamMatchesByMatchType(nbaWeekID, seasonID, matchType)

	arenaMap := GetArenaMap()

	collegeTeams := GetAllActiveCollegeTeams()
	collegeTeamMap := MakeCollegeTeamMap(collegeTeams)
	nbaTeams := GetAllActiveNBATeams()
	nbaTeamMap := MakeNBATeamMap(nbaTeams)
	collegePlayers := GetAllCollegePlayers()
	collegePlayerMap := MakeCollegePlayerMapByTeamID(collegePlayers, true)
	nbaPlayers := GetAllNBAPlayers()
	nbaPlayerMap := MakeNBAPlayerMapByTeamID(nbaPlayers, true)
	collegeStandings := repository.FindAllCollegeStandingsRecords(repository.StandingsQuery{SeasonID: seasonID})
	collegeStandingsMap := MakeCollegeStandingsMap(collegeStandings)
	nbaStandings := repository.FindAllNBAStandingsRecords(repository.StandingsQuery{SeasonID: seasonID})
	nbaStandingsMap := MakeNBAStandingsMap(nbaStandings)
	collegeLineups := repository.FindCollegeLineupRecords(repository.GameplanQuery{})
	nbaLineups := repository.FindNBALineupRecords(repository.GameplanQuery{})
	collegeLineupMap := MakeCollegeLineupMapByTeamID(collegeLineups)
	nbaLineupMap := MakeNBALineupMapByTeamID(nbaLineups)
	collegeGameplans := GetAllCollegeGameplans()
	nbaGameplans := GetAllNBAGameplans()
	collegeGameplansMap := MakeCollegeGameplanMap(collegeGameplans)
	nbaGameplansMap := MakeNBAGameplanMap(nbaGameplans)

	matchesList := make([]structs.MatchResponse, 0, len(collegeMatches)+len(nbaMatches))

	if matchType == "" {
		return structs.MatchStateResponse{Matches: matchesList}
	}

	for _, c := range collegeMatches {
		if c.GameComplete {
			continue
		}

		ht := collegeTeamMap[c.HomeTeamID]
		at := collegeTeamMap[c.AwayTeamID]
		homeTeamLineup := collegeLineupMap[c.HomeTeamID]
		awayTeamLineup := collegeLineupMap[c.AwayTeamID]

		homeGameplan := collegeGameplansMap[c.HomeTeamID]
		awayGameplan := collegeGameplansMap[c.AwayTeamID]

		htGameLineUp := []structs.GameplanLineup{}
		for _, lp := range homeTeamLineup {
			htGameLineUp = append(htGameLineUp, lp.GameplanLineup)
		}

		atGameLineUp := []structs.GameplanLineup{}
		for _, lp := range awayTeamLineup {
			atGameLineUp = append(atGameLineUp, lp.GameplanLineup)
		}

		capacity := 0
		arena := arenaMap[ht.Arena]
		if arena.ID == 0 {
			capacity = 6000
		} else {
			capacity = int(arena.Capacity)
		}

		currentStandings := collegeStandingsMap[c.HomeTeamID]
		attendancePercent := getAttendancePercent(int(currentStandings.TotalWins), int(currentStandings.TotalLosses))
		if c.Week == 0 {
			attendancePercent = 1.0
		}
		fanCount := uint32(float64(capacity) * attendancePercent)
		hra := float64(fanCount) / float64(capacity)
		homeRoster := collegePlayerMap[c.HomeTeamID]
		awayRoster := collegePlayerMap[c.AwayTeamID]

		homeGamePlayerRoster := []structs.GamePlayer{}
		awayGamePlayerRoster := []structs.GamePlayer{}

		for _, p := range homeRoster {
			gamePlayer := structs.GamePlayer{
				ID:         p.ID,
				BasePlayer: p.BasePlayer,
			}
			gamePlayer.CalculateModifiers(true, hra)
			homeGamePlayerRoster = append(homeGamePlayerRoster, gamePlayer)
		}

		for _, p := range awayRoster {
			gamePlayer := structs.GamePlayer{
				ID:         p.ID,
				BasePlayer: p.BasePlayer,
			}
			gamePlayer.CalculateModifiers(false, hra)
			awayGamePlayerRoster = append(awayGamePlayerRoster, gamePlayer)
		}

		livestreamChannel := 0
		if (ht.IsUserCoached) || (at.IsUserCoached) {
			if ht.ConferenceID%2 == 1 {
				livestreamChannel = 1
			} else {
				livestreamChannel = 2
			}
		} else {
			if ht.ConferenceID%2 == 1 {
				livestreamChannel = 3
			} else {
				livestreamChannel = 4
			}
		}

		if c.IsPlayoffGame {
			livestreamChannel = 1
		} else if c.IsNITGame {
			livestreamChannel = 3
		} else if c.IsCBIGame {
			livestreamChannel = 4
		}

		match := structs.MatchResponse{
			MatchName:              c.MatchName,
			ID:                     c.ID,
			WeekID:                 c.WeekID,
			Week:                   c.Week,
			SeasonID:               c.SeasonID,
			HomeTeamID:             c.HomeTeamID,
			HomeTeam:               c.HomeTeam,
			AwayTeamID:             c.AwayTeamID,
			AwayTeam:               c.AwayTeam,
			MatchOfWeek:            c.MatchOfWeek,
			Arena:                  c.Arena,
			Capacity:               capacity,
			City:                   c.City,
			State:                  c.State,
			IsNeutralSite:          c.IsNeutralSite,
			IsNBAMatch:             false,
			IsConference:           c.IsConference,
			IsConferenceTournament: c.IsConferenceTournament,
			IsNITGame:              c.IsNITGame,
			IsPlayoffGame:          c.IsPlayoffGame,
			IsNationalChampionship: c.IsNationalChampionship,
			IsRivalryGame:          c.IsRivalryGame,
			IsInvitational:         c.IsInvitational,
			Channel:                uint(livestreamChannel),
			MatchData: structs.MatchDataResponse{
				HomeTeam: structs.MatchTeamResponse{
					ID:           ht.ID,
					TeamName:     ht.Team,
					Mascot:       ht.Nickname,
					Abbr:         ht.Abbr,
					Conference:   ht.Conference,
					Coach:        ht.Coach,
					ConferenceID: ht.ConferenceID,
					LeagueID:     1,
				},
				HomeTeamRoster:   homeGamePlayerRoster,
				HomeTeamLineup:   htGameLineUp,
				HomeTeamGameplan: homeGameplan,
				AwayTeam: structs.MatchTeamResponse{
					ID:           at.ID,
					TeamName:     at.Team,
					Mascot:       at.Nickname,
					Abbr:         at.Abbr,
					Conference:   at.Conference,
					Coach:        at.Coach,
					ConferenceID: at.ConferenceID,
					LeagueID:     1,
				},
				AwayTeamRoster:   awayGamePlayerRoster,
				AwayTeamLineup:   atGameLineUp,
				AwayTeamGameplan: awayGameplan,
			},
		}
		matchesList = append(matchesList, match)
	}

	// Iterate NBA Matches
	coinFlip := false
	for _, m := range nbaMatches {
		if m.GameComplete {
			continue
		}

		livestreamChannel := 0
		if coinFlip {
			livestreamChannel = 5
		} else {
			livestreamChannel = 6
		}
		coinFlip = !coinFlip
		if m.IsInternational {
			livestreamChannel = 7
		}

		ht := nbaTeamMap[m.HomeTeamID]
		at := nbaTeamMap[m.AwayTeamID]

		homeTeamLineUp := nbaLineupMap[ht.ID]
		awayTeamLineUp := nbaLineupMap[at.ID]

		homeGameplanLineup := []structs.GameplanLineup{}
		awayGameplanLineup := []structs.GameplanLineup{}

		for _, lp := range homeTeamLineUp {
			homeGameplanLineup = append(homeGameplanLineup, lp.GameplanLineup)
		}

		for _, lp := range awayTeamLineUp {
			awayGameplanLineup = append(awayGameplanLineup, lp.GameplanLineup)
		}

		hg := nbaGameplansMap[ht.ID]
		ag := nbaGameplansMap[at.ID]

		homeGameplan := structs.Gameplan{
			TeamID:             ht.ID,
			OffensiveFormation: hg.OffensiveFormation,
			DefensiveFormation: hg.DefensiveFormation,
			OffensiveStyle:     hg.OffensiveStyle,
			Pace:               hg.Pace,
			FocusPlayer:        hg.FocusPlayer,
		}

		awayGameplan := structs.Gameplan{
			TeamID:             at.ID,
			OffensiveFormation: ag.OffensiveFormation,
			DefensiveFormation: ag.DefensiveFormation,
			OffensiveStyle:     ag.OffensiveStyle,
			Pace:               ag.Pace,
			FocusPlayer:        ag.FocusPlayer,
		}

		capacity := 0
		arena := arenaMap[ht.Arena]
		if arena.ID == 0 {
			capacity = 6000
		} else {
			capacity = int(arena.Capacity)
		}
		currentStandings := nbaStandingsMap[m.HomeTeamID]
		attendancePercent := getAttendancePercent(int(currentStandings.TotalWins), int(currentStandings.TotalLosses))
		if m.Week == 0 {
			attendancePercent = 1.0
		}
		fanCount := uint32(float64(capacity) * attendancePercent)
		hra := float64(fanCount) / float64(capacity)
		homeRoster := nbaPlayerMap[m.HomeTeamID]
		awayRoster := nbaPlayerMap[m.AwayTeamID]

		homeGamePlayerRoster := []structs.GamePlayer{}
		awayGamePlayerRoster := []structs.GamePlayer{}

		for _, p := range homeRoster {
			gamePlayer := structs.GamePlayer{
				ID:         p.ID,
				BasePlayer: p.BasePlayer,
			}
			gamePlayer.CalculateModifiers(true, hra)
			homeGamePlayerRoster = append(homeGamePlayerRoster, gamePlayer)
		}

		for _, p := range awayRoster {
			gamePlayer := structs.GamePlayer{
				ID:         p.ID,
				BasePlayer: p.BasePlayer,
			}
			gamePlayer.CalculateModifiers(false, hra)
			awayGamePlayerRoster = append(awayGamePlayerRoster, gamePlayer)
		}

		match := structs.MatchResponse{
			MatchName:              m.MatchName,
			ID:                     m.ID,
			WeekID:                 m.WeekID,
			Week:                   m.Week,
			SeasonID:               m.SeasonID,
			HomeTeamID:             m.HomeTeamID,
			HomeTeam:               m.HomeTeam,
			AwayTeamID:             m.AwayTeamID,
			AwayTeam:               m.AwayTeam,
			MatchOfWeek:            m.MatchOfWeek,
			Arena:                  m.Arena,
			City:                   m.City,
			State:                  m.State,
			IsNeutralSite:          m.IsNeutralSite,
			IsNBAMatch:             true,
			IsConference:           m.IsConference,
			IsConferenceTournament: m.IsConferenceTournament,
			IsInternational:        m.IsInternational,
			IsPlayoffGame:          m.IsPlayoffGame,
			IsNationalChampionship: m.IsTheFinals,
			IsRivalryGame:          m.IsRivalryGame,
			Channel:                uint(livestreamChannel),
			MatchData: structs.MatchDataResponse{
				HomeTeam: structs.MatchTeamResponse{
					ID:           ht.ID,
					TeamName:     ht.Team,
					Mascot:       ht.Nickname,
					Abbr:         ht.Abbr,
					Conference:   ht.Conference,
					Coach:        ht.NBACoachName,
					ConferenceID: ht.ConferenceID,
					LeagueID:     1,
				},
				HomeTeamRoster:   homeGamePlayerRoster,
				HomeTeamLineup:   homeGameplanLineup,
				HomeTeamGameplan: homeGameplan,
				AwayTeam: structs.MatchTeamResponse{
					ID:           at.ID,
					TeamName:     at.Team,
					Mascot:       at.Nickname,
					Abbr:         at.Abbr,
					Conference:   at.Conference,
					Coach:        at.NBACoachName,
					ConferenceID: at.ConferenceID,
					LeagueID:     1,
				},
				AwayTeamRoster:   awayGamePlayerRoster,
				AwayTeamLineup:   awayGameplanLineup,
				AwayTeamGameplan: awayGameplan,
			},
		}
		matchesList = append(matchesList, match)
	}

	return structs.MatchStateResponse{
		Matches:   matchesList,
		MatchType: matchType,
		Week:      uint(ts.NBAWeek),
	}
}

func GetMatchesByTeamIdAndSeasonId(teamId string, seasonId string) []structs.Match {
	db := dbprovider.GetInstance().GetDB()

	var teamMatches []structs.Match

	db.Where("(home_team_id = ? OR away_team_id = ?) AND season_id = ?", teamId, teamId, seasonId).Find(&teamMatches)

	return teamMatches
}

func GetCBBMatchesByTeamId(teamId string) []structs.Match {
	db := dbprovider.GetInstance().GetDB()

	var teamMatches []structs.Match

	db.Where("(home_team_id = ? OR away_team_id = ?)", teamId, teamId).Find(&teamMatches)

	return teamMatches
}

func GetNBAMatchesByTeamId(teamId string) []structs.NBAMatch {
	db := dbprovider.GetInstance().GetDB()

	var teamMatches []structs.NBAMatch

	db.Where("(home_team_id = ? OR away_team_id = ?)", teamId, teamId).Find(&teamMatches)

	return teamMatches
}

func GetProfessionalMatchesByTeamIdAndSeasonId(teamId string, seasonId string) []structs.NBAMatch {
	db := dbprovider.GetInstance().GetDB()

	var teamMatches []structs.NBAMatch

	db.Where("(home_team_id = ? OR away_team_id = ?) AND season_id = ?", teamId, teamId, seasonId).Find(&teamMatches)

	return teamMatches
}

func FixPlayerStatsFromLastSeason() {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	lastSeasonID := ts.SeasonID - 1
	seasonIDSTR := strconv.Itoa(int(lastSeasonID))

	matches := GetCBBMatchesBySeasonID(seasonIDSTR)

	for _, m := range matches {
		id := strconv.Itoa(int(m.ID))
		stats := GetPlayerStatsByMatch(id)

		for _, stat := range stats {
			if stat.WeekID > 0 {
				continue
			}
			stat.MapNewProperties(m.WeekID, m.MatchOfWeek)

			db.Save(&stat)
		}
	}
}

func GetSchedulePageData(seasonId string) structs.MatchPageResponse {
	ts := GetTimestamp()
	collegeMatches := GetCBBMatchesBySeasonID(seasonId)
	nbaMatches := GetNBAMatchesBySeasonID(seasonId)
	islMatches := GetISLMatchesBySeasonID(seasonId)

	for _, m := range collegeMatches {
		showResults := m.Week <= uint(ts.CollegeWeek)

		if m.Week == uint(ts.CollegeWeek) && m.GameComplete &&
			((m.MatchOfWeek == "A" && !ts.GamesARan) || (m.MatchOfWeek == "B" && !ts.GamesBRan) ||
				(m.MatchOfWeek == "C" && !ts.GamesCRan) || (m.MatchOfWeek == "D" && !ts.GamesDRan)) {
			showResults = false
		}

		if !showResults {
			m.HideScore()
		}
	}

	for _, m := range nbaMatches {
		showResults := m.Week <= uint(ts.NBAWeek)

		if m.Week == uint(ts.NBAWeek) && m.GameComplete &&
			((m.MatchOfWeek == "A" && !ts.GamesARan) || (m.MatchOfWeek == "B" && !ts.GamesBRan) ||
				(m.MatchOfWeek == "C" && !ts.GamesCRan) || (m.MatchOfWeek == "D" && !ts.GamesDRan)) {
			showResults = false
		}

		if !showResults {
			m.HideScore()
		}
	}

	for _, m := range islMatches {
		showResults := m.Week <= uint(ts.NBAWeek)

		if m.Week == uint(ts.NBAWeek) && m.GameComplete &&
			((m.MatchOfWeek == "A" && !ts.GamesARan) || (m.MatchOfWeek == "B" && !ts.GamesBRan) ||
				(m.MatchOfWeek == "C" && !ts.GamesCRan) || (m.MatchOfWeek == "D" && !ts.GamesDRan)) {
			showResults = false
		}

		if !showResults {
			m.HideScore()
		}
	}

	return structs.MatchPageResponse{
		CBBGames: collegeMatches,
		NBAGames: nbaMatches,
		ISLGames: islMatches,
	}
}

func GetCBBMatchesBySeasonID(seasonId string) []structs.Match {
	db := dbprovider.GetInstance().GetDB()

	var teamMatches []structs.Match

	db.Order("week_id asc").Order("match_of_week asc").Where("season_id = ?", seasonId).Find(&teamMatches)

	return teamMatches
}

func GetNBAMatchesBySeasonID(seasonId string) []structs.NBAMatch {
	db := dbprovider.GetInstance().GetDB()

	var teamMatches []structs.NBAMatch

	db.Order("week_id asc").Order("match_of_week asc").Where("season_id = ? AND is_international = false", seasonId).Find(&teamMatches)

	return teamMatches
}

func GetISLMatchesBySeasonID(seasonId string) []structs.NBAMatch {
	db := dbprovider.GetInstance().GetDB()

	var teamMatches []structs.NBAMatch

	db.Order("week_id asc").Order("match_of_week asc").Where("season_id = ? AND is_international = true", seasonId).Find(&teamMatches)

	return teamMatches
}

func GetMatchByMatchId(matchId string) structs.Match {
	db := dbprovider.GetInstance().GetDB()

	var match structs.Match

	err := db.Where("id = ?", matchId).Find(&match).Error
	if err != nil {
		fmt.Println(err.Error())
	}

	return match
}

func GetNBAMatchByMatchId(matchId string) structs.NBAMatch {
	db := dbprovider.GetInstance().GetDB()

	var match structs.NBAMatch

	err := db.Where("id = ?", matchId).Find(&match).Error
	if err != nil {
		fmt.Println(err.Error())
	}

	return match
}

func AddNBAMatches() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	// Get team information for arena and team names
	nbaTeamMap := GetProfessionalTeamMap()
	firstMatch := CreateNBAMatch(ts, nbaTeamMap, 30, 27, "C", 62)
	secondMatch := CreateNBAMatch(ts, nbaTeamMap, 128, 47, "C", 80)

	db.Create(&firstMatch)
	db.Create(&secondMatch)
}

// CreateNBAMatch creates a new NBAMatch struct with the provided team IDs and optional series ID
func CreateNBAMatch(ts structs.Timestamp, nbaTeamMap map[uint]structs.NBATeam, homeTeamID, awayTeamID uint, gameDay string, seriesID uint) structs.NBAMatch {

	// Convert team IDs to integers for the struct

	homeTeam := nbaTeamMap[uint(homeTeamID)]
	awayTeam := nbaTeamMap[uint(awayTeamID)]

	// Create the match struct
	match := structs.NBAMatch{
		WeekID:        ts.NBAWeekID,
		Week:          uint(ts.NBAWeek),
		SeasonID:      ts.SeasonID,
		HomeTeamID:    uint(homeTeamID),
		HomeTeam:      homeTeam.Team,
		AwayTeamID:    uint(awayTeamID),
		AwayTeam:      awayTeam.Team,
		HomeTeamCoach: homeTeam.NBAOwnerName,
		AwayTeamCoach: awayTeam.NBAOwnerName,
		Arena:         homeTeam.Arena,
		City:          homeTeam.City,
		State:         homeTeam.State,
		MatchOfWeek:   gameDay,
		GameComplete:  false,
		HomeTeamWin:   false,
		AwayTeamWin:   false,
	}

	// Set series ID if provided
	if seriesID > 0 {
		match.SeriesID = seriesID
		match.IsPlayoffGame = true
	}

	return match
}

func GetMatchResultsByMatchID(matchId string) structs.MatchResultsResponse {
	match := GetMatchByMatchId(matchId)
	involvedPlayers := GetCollegePlayersWithMatchStatsByTeamId(match.HomeTeamID, match.AwayTeamID, matchId)
	homePlayers := []structs.MatchResultsPlayer{}
	awayPlayers := []structs.MatchResultsPlayer{}
	for _, p := range involvedPlayers {
		if p.TeamID == match.HomeTeamID {
			homePlayers = append(homePlayers, p)
		} else {
			awayPlayers = append(awayPlayers, p)
		}
	}
	homeStats := GetCBBTeamResultsByMatch(strconv.Itoa(int(match.HomeTeamID)), matchId)
	awayStats := GetCBBTeamResultsByMatch(strconv.Itoa(int(match.AwayTeamID)), matchId)

	return structs.MatchResultsResponse{
		HomePlayers: homePlayers,
		AwayPlayers: awayPlayers,
		HomeStats:   homeStats,
		AwayStats:   awayStats,
	}
}

func GetNBAMatchResultsByMatchID(matchId string) structs.MatchResultsResponse {
	match := GetNBAMatchByMatchId(matchId)
	involvedPlayers := GetNBAPlayersWithMatchStatsByTeamId(match.HomeTeamID, match.AwayTeamID, matchId)
	homePlayers := []structs.MatchResultsPlayer{}
	awayPlayers := []structs.MatchResultsPlayer{}
	for _, p := range involvedPlayers {
		if p.TeamID == match.HomeTeamID {
			homePlayers = append(homePlayers, p)
		} else {
			awayPlayers = append(awayPlayers, p)
		}
	}
	homeStats := GetNBATeamResultsByMatch(strconv.Itoa(int(match.HomeTeamID)), matchId)
	awayStats := GetNBATeamResultsByMatch(strconv.Itoa(int(match.AwayTeamID)), matchId)

	return structs.MatchResultsResponse{
		HomePlayers: homePlayers,
		AwayPlayers: awayPlayers,
		HomeStats:   homeStats,
		AwayStats:   awayStats,
	}
}

func GetTeamMatchesByWeekId(weekId, seasonID, matchType, teamID string) []structs.Match {
	db := dbprovider.GetInstance().GetDB()

	var teamMatches []structs.Match

	db.Where("week_id = ? AND season_id = ? AND match_of_week = ? AND (home_team_id = ? OR away_team_id = ?)", weekId, seasonID, matchType, teamID, teamID).Find(&teamMatches)

	return teamMatches
}

func GetCollegeTeamMatchesBySeasonId(seasonID, teamID string) []structs.Match {
	db := dbprovider.GetInstance().GetDB()

	var teamMatches []structs.Match

	db.Where("season_id = ?  AND (home_team_id = ? OR away_team_id = ?)", seasonID, teamID, teamID).Find(&teamMatches)

	return teamMatches
}

func GetNBATeamMatchesByMatchType(weekId, seasonID, matchType string) []structs.NBAMatch {
	db := dbprovider.GetInstance().GetDB()

	var teamMatches []structs.NBAMatch

	db.Where("week_id = ? AND season_id = ? AND match_of_week = ?", weekId, seasonID, matchType).Find(&teamMatches)

	return teamMatches
}

func GetNBATeamMatchesByWeekId(weekId, seasonID, matchType, teamID string) []structs.NBAMatch {
	db := dbprovider.GetInstance().GetDB()

	var teamMatches []structs.NBAMatch

	db.Where("week_id = ? AND season_id = ? AND match_of_week = ? AND (home_team_id = ? OR away_team_id = ?)", weekId, seasonID, matchType, teamID, teamID).Find(&teamMatches)

	return teamMatches
}

func GetNBATeamMatchesBySeasonId(seasonID, teamID string) []structs.NBAMatch {
	db := dbprovider.GetInstance().GetDB()

	var teamMatches []structs.NBAMatch

	db.Where("season_id = ? AND (home_team_id = ? OR away_team_id = ?)", seasonID, teamID, teamID).Find(&teamMatches)

	return teamMatches
}

func GetNBAMatchesByWeekIdAndMatchType(weekId string, seasonID string, matchType string) []structs.NBAMatch {
	db := dbprovider.GetInstance().GetDB()

	var teamMatches []structs.NBAMatch

	db.Where("week_id = ? AND season_id = ? AND match_of_week = ?", weekId, seasonID, matchType).Find(&teamMatches)

	return teamMatches
}

func GetMatchesByWeekIdAndMatchType(weekId string, seasonID string, matchType string) []structs.Match {
	db := dbprovider.GetInstance().GetDB()

	var teamMatches []structs.Match

	db.Where("week_id = ? AND season_id = ? AND match_of_week = ?", weekId, seasonID, matchType).Find(&teamMatches)

	return teamMatches
}

func GetMatchesByWeekId(weekId, seasonID string) []structs.Match {
	db := dbprovider.GetInstance().GetDB()

	var teamMatches []structs.Match

	db.Order("match_of_week asc").Where("week_id = ? AND season_id = ?", weekId, seasonID).Find(&teamMatches)

	return teamMatches
}

func GetUpcomingMatchesByTeamIdAndSeasonId(teamId string, seasonId string) []structs.Match {
	db := dbprovider.GetInstance().GetDB()

	timeStamp := GetTimestamp()

	var teamMatches []structs.Match

	db.Where("team_id = ? AND season_id = ? AND week_id > ?", teamId, seasonId, timeStamp.CollegeWeekID).Find(teamMatches)

	return teamMatches
}

// GetNBASeriesBySeriesID -- Get an NBA Playoff Series Record
func GetNBASeriesBySeriesID(seriesID string) structs.NBASeries {
	db := dbprovider.GetInstance().GetDB()

	var nbaSeries structs.NBASeries

	db.Where("id = ?", seriesID).Find(&nbaSeries)

	return nbaSeries
}

func GetNBASeriesByTeamID(teamID string) []structs.NBASeries {
	db := dbprovider.GetInstance().GetDB()

	var nbaSeries []structs.NBASeries

	db.Where("home_team_id = ? OR away_team_id = ?", teamID, teamID).Find(&nbaSeries)

	return nbaSeries
}

func GetLatestNBASeriesID() uint {
	db := dbprovider.GetInstance().GetDB()

	var nbaSeries structs.NBASeries

	db.Order("id desc").First(&nbaSeries)

	return nbaSeries.ID
}

func GetLatestNBAMatchID() uint {
	db := dbprovider.GetInstance().GetDB()

	var nbaMatch structs.NBAMatch

	db.Order("id desc").First(&nbaMatch)

	return nbaMatch.ID
}

func GetLatestCollegeMatchID() uint {
	db := dbprovider.GetInstance().GetDB()

	var collegeMatch structs.Match

	db.Order("id desc").First(&collegeMatch)

	return collegeMatch.ID
}

func GetAllActiveNBASeries(ts structs.Timestamp) []structs.NBASeries {
	db := dbprovider.GetInstance().GetDB()

	var nbaSeries []structs.NBASeries

	db.Where("season_id = ?", ts.SeasonID).Find(&nbaSeries)

	return nbaSeries
}

func SwapNBATeamsTEMP() {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.SeasonID))
	nbaMatches := GetNBAMatchesBySeasonID(seasonID)
	// islMatches := GetISLMatchesBySeasonID(seasonID)
	nbaTeamMap := GetProfessionalTeamMap()

	skippingID := 2147
	matchDay := "A"
	for _, m := range nbaMatches {
		if m.ID > uint(skippingID) && m.MatchOfWeek == matchDay {
			m.ResetScore()
			m.SwapTeams()
			nbaTeam := nbaTeamMap[m.HomeTeamID]
			m.AssignArena(nbaTeam.Arena, nbaTeam.City, nbaTeam.State)
			repository.SaveProfessionalMatchRecord(m, db)
		}
	}

	// for _, m := range islMatches {
	// 	if m.ID > uint(skippingID) && m.MatchOfWeek == matchDay {
	// 		m.ResetScore()
	// 		m.SwapTeams()
	// 		nbaTeam := nbaTeamMap[m.HomeTeamID]
	// 		m.AssignArena(nbaTeam.Arena, nbaTeam.City, nbaTeam.State)
	// 		repository.SaveProfessionalMatchRecord(m, db)
	// 	}
	// }
}

func FixISLMatchData() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.SeasonID))
	nbaMap := GetProfessionalTeamMapBByLabel()

	matches := GetISLMatchesBySeasonID(seasonID)

	for _, m := range matches {
		if !m.IsInternational {
			continue
		}

		homeTeam := nbaMap[m.HomeTeam]
		awayTeam := nbaMap[m.AwayTeam]
		m.AddTeam(true, homeTeam.ID, 0, m.HomeTeam, homeTeam.NBAOwnerName, homeTeam.Arena, homeTeam.City, homeTeam.State)
		m.AddTeam(false, awayTeam.ID, 0, m.AwayTeam, awayTeam.NBAOwnerName, "", "", "")

		db.Save(&m)
	}
}

func getAttendancePercent(wins, losses int) float64 {
	totalGames := wins + losses
	if totalGames < 4 {
		return util.GenerateFloatFromRange(0.90, 1.00)
	}

	winRate := float64(wins) / float64(totalGames)

	switch {
	case winRate >= 0.75:
		return util.GenerateFloatFromRange(0.95, 1.05)
	case winRate >= 0.5:
		return util.GenerateFloatFromRange(0.85, 0.94)
	case winRate >= 0.35:
		return util.GenerateFloatFromRange(0.65, 0.84)
	default:
		return util.GenerateFloatFromRange(0.4, 0.64)
	}
}
