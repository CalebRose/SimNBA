package managers

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
)

func GetMatchesForTimeslot() structs.MatchStateResponse {
	ts := GetTimestamp()
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

	collegeMatchesWg.Add(len(collegeMatches))
	nbaMatchesWg.Add(len(nbaMatches))

	matchesList := make([]structs.MatchResponse, 0, len(collegeMatches)+len(nbaMatches))

	if matchType == "" {
		return structs.MatchStateResponse{Matches: matchesList}
	}
	sem := make(chan struct{}, 20) // Limit to 20 concurrent tasks

	for _, c := range collegeMatches {
		if c.GameComplete {
			continue
		}
		sem <- struct{}{}
		go func(c structs.Match) {
			defer func() { <-sem }()
			defer collegeMatchesWg.Done()

			ht := GetTeamByTeamID(strconv.Itoa(int(c.HomeTeamID)))
			at := GetTeamByTeamID(strconv.Itoa(int(c.AwayTeamID)))

			livestreamChannel := 0
			if (c.HomeTeamCoach != "AI" && c.HomeTeamCoach != "") || (c.AwayTeamCoach != "AI" && c.AwayTeamCoach != "") {
				livestreamChannel = 1
			} else if ht.ConferenceID < 6 || ht.ConferenceID == 11 || at.ConferenceID < 6 || at.ConferenceID == 11 {
				livestreamChannel = 2
			} else {
				livestreamChannel = 3
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
		}(c)
	}

	// Iterate NBA Matches
	coinFlip := false
	for _, n := range nbaMatches {
		if n.GameComplete {
			continue
		}
		sem <- struct{}{}
		go func(n structs.NBAMatch) { // replace `YourNBAMatchType` with the actual type
			defer func() { <-sem }()
			defer nbaMatchesWg.Done()
			livestreamChannel := 0
			if coinFlip {
				livestreamChannel = 4
				coinFlip = !coinFlip
			} else if !coinFlip {
				livestreamChannel = 5
				coinFlip = !coinFlip
			}
			if n.IsInternational {
				livestreamChannel = 6
			}

			match := structs.MatchResponse{
				MatchName:              n.MatchName,
				ID:                     n.ID,
				WeekID:                 n.WeekID,
				Week:                   n.Week,
				SeasonID:               n.SeasonID,
				HomeTeamID:             n.HomeTeamID,
				HomeTeam:               n.HomeTeam,
				AwayTeamID:             n.AwayTeamID,
				AwayTeam:               n.AwayTeam,
				MatchOfWeek:            n.MatchOfWeek,
				Arena:                  n.Arena,
				City:                   n.City,
				State:                  n.State,
				IsNeutralSite:          n.IsNeutralSite,
				IsNBAMatch:             true,
				IsConference:           n.IsConference,
				IsConferenceTournament: n.IsConferenceTournament,
				IsInternational:        n.IsInternational,
				IsPlayoffGame:          n.IsPlayoffGame,
				IsNationalChampionship: n.IsTheFinals,
				IsRivalryGame:          n.IsRivalryGame,
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
	collegeMatches := GetCBBMatchesBySeasonID(seasonId)
	nbaMatches := GetNBAMatchesBySeasonID(seasonId)
	islMatches := GetISLMatchesBySeasonID(seasonId)

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

	homePlayers := GetCollegePlayersWithMatchStatsByTeamId(strconv.Itoa(int(match.HomeTeamID)), matchId)
	awayPlayers := GetCollegePlayersWithMatchStatsByTeamId(strconv.Itoa(int(match.AwayTeamID)), matchId)
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

	homePlayers := GetNBAPlayersWithMatchStatsByTeamId(strconv.Itoa(int(match.HomeTeamID)), matchId)
	awayPlayers := GetNBAPlayersWithMatchStatsByTeamId(strconv.Itoa(int(match.AwayTeamID)), matchId)
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

// SAVE
