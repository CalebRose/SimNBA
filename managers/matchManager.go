package managers

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/repository"
	"github.com/CalebRose/SimNBA/structs"
)

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

	// Wait Groups
	var collegeMatchesWg, nbaMatchesWg sync.WaitGroup

	// Mutex Lock
	var mutex sync.Mutex

	// Get College Matches
	collegeMatches := GetMatchesByWeekIdAndMatchType(weekID, seasonID, matchType)
	// Get Professional Matches
	nbaMatches := GetNBATeamMatchesByMatchType(nbaWeekID, seasonID, matchType)

	arenaMap := GetArenaMap()

	collegeMatchesWg.Add(len(collegeMatches))
	nbaMatchesWg.Add(len(nbaMatches))

	matchesList := make([]structs.MatchResponse, 0, len(collegeMatches)+len(nbaMatches))

	if matchType == "" {
		return structs.MatchStateResponse{Matches: matchesList}
	}
	sem := make(chan struct{}, 20) // Limit to 20 concurrent tasks

	for _, c := range collegeMatches {
		sem <- struct{}{}
		localC := c
		go func(c structs.Match) {
			defer func() { <-sem }()
			defer collegeMatchesWg.Done()
			if c.GameComplete {
				return
			}
			ht := GetTeamByTeamID(strconv.Itoa(int(c.HomeTeamID)))
			at := GetTeamByTeamID(strconv.Itoa(int(c.AwayTeamID)))

			capacity := 0

			mutex.Lock()
			arena := arenaMap[ht.Arena]
			if arena.ID == 0 {
				capacity = 6000
			} else {
				capacity = int(arena.Capacity)
			}
			mutex.Unlock()

			livestreamChannel := 0
			if (ht.IsUserCoached) || (at.IsUserCoached) {
				livestreamChannel = 1
			} else if ht.ConferenceID < 6 || ht.ConferenceID == 11 || at.ConferenceID < 6 || at.ConferenceID == 11 {
				livestreamChannel = 2
			} else if !(ht.ConferenceID < 6 || ht.ConferenceID == 11 || at.ConferenceID < 6 || at.ConferenceID == 11) && ((ht.ConferenceID > 5 && ht.ConferenceID < 14) || (at.ConferenceID > 5 && at.ConferenceID < 14)) {
				livestreamChannel = 3
			} else {
				livestreamChannel = 4
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
			}

			mutex.Lock()
			matchesList = append(matchesList, match)
			mutex.Unlock()
		}(localC)
	}

	// Iterate NBA Matches
	coinFlip := false
	for _, n := range nbaMatches {
		sem <- struct{}{}
		go func(m structs.NBAMatch) { // replace `YourNBAMatchType` with the actual type
			defer func() { <-sem }()
			defer nbaMatchesWg.Done()
			if m.GameComplete {
				return
			}
			livestreamChannel := 0
			if coinFlip {
				livestreamChannel = 5
				coinFlip = !coinFlip
			} else if !coinFlip {
				livestreamChannel = 6
				coinFlip = !coinFlip
			}
			if m.IsInternational {
				livestreamChannel = 7
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
			}
			mutex.Lock()
			matchesList = append(matchesList, match)
			mutex.Unlock()
		}(n)
	}
	collegeMatchesWg.Wait()
	nbaMatchesWg.Wait()

	for i := 0; i < cap(sem); i++ {
		sem <- struct{}{}
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

	db.Where("season_id = ?", seasonId).Find(&teamMatches)

	return teamMatches
}

func GetNBAMatchesBySeasonID(seasonId string) []structs.NBAMatch {
	db := dbprovider.GetInstance().GetDB()

	var teamMatches []structs.NBAMatch

	db.Where("season_id = ? AND is_international = false", seasonId).Find(&teamMatches)

	return teamMatches
}

func GetISLMatchesBySeasonID(seasonId string) []structs.NBAMatch {
	db := dbprovider.GetInstance().GetDB()

	var teamMatches []structs.NBAMatch

	db.Where("season_id = ? AND is_international = true", seasonId).Find(&teamMatches)

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

func GetAllActiveNBASeries() []structs.NBASeries {
	db := dbprovider.GetInstance().GetDB()

	var nbaSeries []structs.NBASeries

	db.Where("series_complete = ?", false).Find(&nbaSeries)

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
