package managers

import (
	"log"
	"sort"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
	"github.com/jinzhu/gorm"
)

func GetAllActiveCollegeTeams() []structs.Team {
	db := dbprovider.GetInstance().GetDB()

	var teams []structs.Team

	err := db.Where("is_active = ? and is_nba = ?", true, false).
		Find(&teams).Error
	if err != nil {
		log.Fatal(err)
	}
	return teams
}

func GetAllActiveCollegeTeamsWithSeasonStats(seasonID, weekID, viewType string) []structs.CollegeTeamResponse {
	db := dbprovider.GetInstance().GetDB()

	var teams []structs.Team

	if viewType == "SEASON" {
		err := db.Preload("TeamSeasonStats", "season_id = ?", seasonID).Where("is_active = ?", true).
			Find(&teams).Error
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := db.Preload("TeamStats", "season_id = ? AND week_id = ?", seasonID, weekID).
			Where("is_active = ?", true).
			Find(&teams).Error
		if err != nil {
			log.Fatal(err)
		}
	}

	var ctResponse []structs.CollegeTeamResponse

	for _, team := range teams {
		if len(team.TeamStats) == 0 && viewType == "WEEK" {
			continue
		}
		if team.TeamSeasonStats.ID == 0 && viewType == "SEASON" {
			continue
		}
		var teamStat structs.TeamStats
		if viewType == "WEEK" {
			teamStat = team.TeamStats[0]
		}
		var seasonsResponse structs.TeamSeasonStatsResponse
		if viewType == "SEASON" {
			seasonsResponse = structs.TeamSeasonStatsResponse{
				ID:                        team.TeamSeasonStats.ID,
				TeamID:                    team.ID,
				SeasonID:                  team.TeamSeasonStats.SeasonID,
				GamesPlayed:               team.TeamSeasonStats.GamesPlayed,
				Points:                    team.TeamSeasonStats.Points,
				PointsAgainst:             team.TeamSeasonStats.PointsAgainst,
				PPG:                       team.TeamSeasonStats.PPG,
				PAPG:                      team.TeamSeasonStats.PAPG,
				PointsDiff:                team.TeamSeasonStats.PPG - team.TeamSeasonStats.PAPG,
				Possessions:               team.TeamSeasonStats.Possessions,
				PossessionsPerGame:        team.TeamSeasonStats.PossessionsPerGame,
				FGM:                       team.TeamSeasonStats.FGM,
				FGA:                       team.TeamSeasonStats.FGA,
				FGPercent:                 team.TeamSeasonStats.FGPercent,
				FGMPG:                     team.TeamSeasonStats.FGMPG,
				FGAPG:                     team.TeamSeasonStats.FGAPG,
				FGMAgainst:                team.TeamSeasonStats.FGMAgainst,
				FGAAgainst:                team.TeamSeasonStats.FGAAgainst,
				FGPercentAgainst:          team.TeamSeasonStats.FGPercentAgainst,
				FGMAPG:                    team.TeamSeasonStats.FGMAPG,
				FGAAPG:                    team.TeamSeasonStats.FGAAPG,
				FGMDiff:                   team.TeamSeasonStats.FGMPG - team.TeamSeasonStats.FGMAPG,
				FGADiff:                   team.TeamSeasonStats.FGAPG - team.TeamSeasonStats.FGAAPG,
				FGPercentDiff:             team.TeamSeasonStats.FGPercent - team.TeamSeasonStats.FGPercentAgainst,
				ThreePointsMade:           team.TeamSeasonStats.ThreePointsMade,
				ThreePointAttempts:        team.TeamSeasonStats.ThreePointAttempts,
				ThreePointPercent:         team.TeamSeasonStats.ThreePointPercent,
				ThreePointsMadeAgainst:    team.TeamSeasonStats.ThreePointsMadeAgainst,
				ThreePointAttemptsAgainst: team.TeamSeasonStats.ThreePointAttemptsAgainst,
				ThreePointPercentAgainst:  team.TeamSeasonStats.ThreePointPercentAgainst,
				TPMPG:                     team.TeamSeasonStats.TPMPG,
				TPAPG:                     team.TeamSeasonStats.TPAPG,
				TPMAPG:                    team.TeamSeasonStats.TPMAPG,
				TPAAPG:                    team.TeamSeasonStats.TPAAPG,
				TPMDiff:                   team.TeamSeasonStats.TPMPG - team.TeamSeasonStats.TPMAPG,
				TPADiff:                   team.TeamSeasonStats.TPAPG - team.TeamSeasonStats.TPAAPG,
				TPPercentDiff:             team.TeamSeasonStats.ThreePointPercent - team.TeamSeasonStats.ThreePointPercentAgainst,
				FTM:                       team.TeamSeasonStats.FTM,
				FTA:                       team.TeamSeasonStats.FTA,
				FTPercent:                 team.TeamSeasonStats.FTPercent,
				FTMAgainst:                team.TeamSeasonStats.FTMAgainst,
				FTAAgainst:                team.TeamSeasonStats.FTAAgainst,
				FTPercentAgainst:          team.TeamSeasonStats.FTPercentAgainst,
				FTMPG:                     team.TeamSeasonStats.FTMPG,
				FTAPG:                     team.TeamSeasonStats.FTAPG,
				FTMAPG:                    team.TeamSeasonStats.FTMAPG,
				FTAAPG:                    team.TeamSeasonStats.FTAAPG,
				FTMDiff:                   team.TeamSeasonStats.FTMPG - team.TeamSeasonStats.FTMAPG,
				FTADiff:                   team.TeamSeasonStats.FTAPG - team.TeamSeasonStats.FTAAPG,
				FTPercentDiff:             team.TeamSeasonStats.FTPercent - team.TeamSeasonStats.FTPercentAgainst,
				Rebounds:                  team.TeamSeasonStats.Rebounds,
				OffRebounds:               team.TeamSeasonStats.OffRebounds,
				DefRebounds:               team.TeamSeasonStats.DefRebounds,
				ReboundsPerGame:           team.TeamSeasonStats.ReboundsPerGame,
				OffReboundsPerGame:        team.TeamSeasonStats.OffReboundsPerGame,
				DefReboundsPerGame:        team.TeamSeasonStats.DefReboundsPerGame,
				ReboundsAllowed:           team.TeamSeasonStats.ReboundsAllowed,
				ReboundsAllowedPerGame:    team.TeamSeasonStats.ReboundsAllowedPerGame,
				OffReboundsAllowed:        team.TeamSeasonStats.OffReboundsAllowed,
				OffReboundsAllowedPerGame: team.TeamSeasonStats.OffReboundsAllowedPerGame,
				DefReboundsAllowed:        team.TeamSeasonStats.DefReboundsAllowed,
				DefReboundsAllowedPerGame: team.TeamSeasonStats.DefReboundsAllowedPerGame,
				ReboundsDiff:              team.TeamSeasonStats.ReboundsPerGame - team.TeamSeasonStats.ReboundsAllowedPerGame,
				OReboundsDiff:             team.TeamSeasonStats.OffReboundsPerGame - team.TeamSeasonStats.OffReboundsAllowedPerGame,
				DReboundsDiff:             team.TeamSeasonStats.DefReboundsPerGame - team.TeamSeasonStats.DefReboundsAllowedPerGame,
				Assists:                   team.TeamSeasonStats.Assists,
				AssistsAllowed:            team.TeamSeasonStats.AssistsAllowed,
				AssistsPerGame:            team.TeamSeasonStats.AssistsPerGame,
				AssistsAllowedPerGame:     team.TeamSeasonStats.AssistsAllowedPerGame,
				AssistsDiff:               team.TeamSeasonStats.AssistsPerGame - team.TeamSeasonStats.AssistsAllowedPerGame,
				Steals:                    team.TeamSeasonStats.Steals,
				StealsAllowed:             team.TeamSeasonStats.StealsAllowed,
				StealsPerGame:             team.TeamSeasonStats.StealsPerGame,
				StealsAllowedPerGame:      team.TeamSeasonStats.StealsAllowedPerGame,
				StealsDiff:                team.TeamSeasonStats.StealsPerGame - team.TeamSeasonStats.StealsAllowedPerGame,
				Blocks:                    team.TeamSeasonStats.Blocks,
				BlocksAllowed:             team.TeamSeasonStats.BlocksAllowed,
				BlocksPerGame:             team.TeamSeasonStats.BlocksPerGame,
				BlocksAllowedPerGame:      team.TeamSeasonStats.BlocksAllowedPerGame,
				BlocksDiff:                team.TeamSeasonStats.BlocksPerGame - team.TeamSeasonStats.BlocksAllowedPerGame,
				TotalTurnovers:            team.TeamSeasonStats.TotalTurnovers,
				TurnoversAllowed:          team.TeamSeasonStats.TurnoversAllowed,
				TurnoversPerGame:          team.TeamSeasonStats.TurnoversPerGame,
				TurnoversAllowedPerGame:   team.TeamSeasonStats.TurnoversAllowedPerGame,
				TODiff:                    team.TeamSeasonStats.TurnoversPerGame - team.TeamSeasonStats.TurnoversAllowedPerGame,
				Fouls:                     team.TeamSeasonStats.Fouls,
				FoulsPerGame:              team.TeamSeasonStats.FoulsPerGame,
			}
		}

		var teamRes = structs.CollegeTeamResponse{
			ID:           team.ID,
			Team:         team.Team,
			Nickname:     team.Nickname,
			Abbr:         team.Abbr,
			ConferenceID: team.ConferenceID,
			Conference:   team.Conference,
			Coach:        team.Coach,
			OverallGrade: team.OverallGrade,
			OffenseGrade: team.OffenseGrade,
			DefenseGrade: team.DefenseGrade,
			IsNBA:        team.IsNBA,
			IsActive:     team.IsActive,
			SeasonStats:  seasonsResponse,
			Stats:        teamStat,
		}

		ctResponse = append(ctResponse, teamRes)
	}
	return ctResponse
}

func GetAllActiveNBATeamsWithSeasonStats(seasonID, weekID, viewType string) []structs.NBATeamResponse {
	db := dbprovider.GetInstance().GetDB()

	var teams []structs.NBATeam

	if viewType == "SEASON" {
		err := db.Preload("TeamSeasonStats", "season_id = ? AND games_played > 0", seasonID).Where("is_active = ?", true).
			Find(&teams).Error
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := db.Preload("TeamStats", "season_id = ? AND week_id = ?", seasonID, weekID).
			Where("is_active = ?", true).
			Find(&teams).Error
		if err != nil {
			log.Fatal(err)
		}
	}

	var ntResponse []structs.NBATeamResponse

	for _, team := range teams {
		if len(team.TeamStats) == 0 && viewType == "WEEK" {
			continue
		}
		var teamStat structs.NBATeamStats
		if viewType == "WEEK" {
			teamStat = team.TeamStats[0]
		}
		var seasonsResponse structs.TeamSeasonStatsResponse
		if viewType == "SEASON" {
			seasonsResponse = structs.TeamSeasonStatsResponse{
				ID:                        team.TeamSeasonStats.ID,
				TeamID:                    team.ID,
				SeasonID:                  team.TeamSeasonStats.SeasonID,
				GamesPlayed:               team.TeamSeasonStats.GamesPlayed,
				Points:                    team.TeamSeasonStats.Points,
				PointsAgainst:             team.TeamSeasonStats.PointsAgainst,
				PPG:                       team.TeamSeasonStats.PPG,
				PAPG:                      team.TeamSeasonStats.PAPG,
				PointsDiff:                team.TeamSeasonStats.PPG - team.TeamSeasonStats.PAPG,
				Possessions:               team.TeamSeasonStats.Possessions,
				PossessionsPerGame:        team.TeamSeasonStats.PossessionsPerGame,
				FGM:                       team.TeamSeasonStats.FGM,
				FGA:                       team.TeamSeasonStats.FGA,
				FGPercent:                 team.TeamSeasonStats.FGPercent,
				FGMPG:                     team.TeamSeasonStats.FGMPG,
				FGAPG:                     team.TeamSeasonStats.FGAPG,
				FGMAgainst:                team.TeamSeasonStats.FGMAgainst,
				FGAAgainst:                team.TeamSeasonStats.FGAAgainst,
				FGPercentAgainst:          team.TeamSeasonStats.FGPercentAgainst,
				FGMAPG:                    team.TeamSeasonStats.FGMAPG,
				FGAAPG:                    team.TeamSeasonStats.FGAAPG,
				FGMDiff:                   team.TeamSeasonStats.FGMPG - team.TeamSeasonStats.FGMAPG,
				FGADiff:                   team.TeamSeasonStats.FGAPG - team.TeamSeasonStats.FGAAPG,
				FGPercentDiff:             team.TeamSeasonStats.FGPercent - team.TeamSeasonStats.FGPercentAgainst,
				ThreePointsMade:           team.TeamSeasonStats.ThreePointsMade,
				ThreePointAttempts:        team.TeamSeasonStats.ThreePointAttempts,
				ThreePointPercent:         team.TeamSeasonStats.ThreePointPercent,
				ThreePointsMadeAgainst:    team.TeamSeasonStats.ThreePointsMadeAgainst,
				ThreePointAttemptsAgainst: team.TeamSeasonStats.ThreePointAttemptsAgainst,
				ThreePointPercentAgainst:  team.TeamSeasonStats.ThreePointPercentAgainst,
				TPMPG:                     team.TeamSeasonStats.TPMPG,
				TPAPG:                     team.TeamSeasonStats.TPAPG,
				TPMAPG:                    team.TeamSeasonStats.TPMAPG,
				TPAAPG:                    team.TeamSeasonStats.TPAAPG,
				TPMDiff:                   team.TeamSeasonStats.TPMPG - team.TeamSeasonStats.TPMAPG,
				TPADiff:                   team.TeamSeasonStats.TPAPG - team.TeamSeasonStats.TPAAPG,
				TPPercentDiff:             team.TeamSeasonStats.ThreePointPercent - team.TeamSeasonStats.ThreePointPercentAgainst,
				FTM:                       team.TeamSeasonStats.FTM,
				FTA:                       team.TeamSeasonStats.FTA,
				FTPercent:                 team.TeamSeasonStats.FTPercent,
				FTMAgainst:                team.TeamSeasonStats.FTMAgainst,
				FTAAgainst:                team.TeamSeasonStats.FTAAgainst,
				FTPercentAgainst:          team.TeamSeasonStats.FTPercentAgainst,
				FTMPG:                     team.TeamSeasonStats.FTMPG,
				FTAPG:                     team.TeamSeasonStats.FTAPG,
				FTMAPG:                    team.TeamSeasonStats.FTMAPG,
				FTAAPG:                    team.TeamSeasonStats.FTAAPG,
				FTMDiff:                   team.TeamSeasonStats.FTMPG - team.TeamSeasonStats.FTMAPG,
				FTADiff:                   team.TeamSeasonStats.FTAPG - team.TeamSeasonStats.FTAAPG,
				FTPercentDiff:             team.TeamSeasonStats.FTPercent - team.TeamSeasonStats.FTPercentAgainst,
				Rebounds:                  team.TeamSeasonStats.Rebounds,
				OffRebounds:               team.TeamSeasonStats.OffRebounds,
				DefRebounds:               team.TeamSeasonStats.DefRebounds,
				ReboundsPerGame:           team.TeamSeasonStats.ReboundsPerGame,
				OffReboundsPerGame:        team.TeamSeasonStats.OffReboundsPerGame,
				DefReboundsPerGame:        team.TeamSeasonStats.DefReboundsPerGame,
				ReboundsAllowed:           team.TeamSeasonStats.ReboundsAllowed,
				ReboundsAllowedPerGame:    team.TeamSeasonStats.ReboundsAllowedPerGame,
				OffReboundsAllowed:        team.TeamSeasonStats.OffReboundsAllowed,
				OffReboundsAllowedPerGame: team.TeamSeasonStats.OffReboundsAllowedPerGame,
				DefReboundsAllowed:        team.TeamSeasonStats.DefReboundsAllowed,
				DefReboundsAllowedPerGame: team.TeamSeasonStats.DefReboundsAllowedPerGame,
				ReboundsDiff:              team.TeamSeasonStats.ReboundsPerGame - team.TeamSeasonStats.ReboundsAllowedPerGame,
				OReboundsDiff:             team.TeamSeasonStats.OffReboundsPerGame - team.TeamSeasonStats.OffReboundsAllowedPerGame,
				DReboundsDiff:             team.TeamSeasonStats.DefReboundsPerGame - team.TeamSeasonStats.DefReboundsAllowedPerGame,
				Assists:                   team.TeamSeasonStats.Assists,
				AssistsAllowed:            team.TeamSeasonStats.AssistsAllowed,
				AssistsPerGame:            team.TeamSeasonStats.AssistsPerGame,
				AssistsAllowedPerGame:     team.TeamSeasonStats.AssistsAllowedPerGame,
				AssistsDiff:               team.TeamSeasonStats.AssistsPerGame - team.TeamSeasonStats.AssistsAllowedPerGame,
				Steals:                    team.TeamSeasonStats.Steals,
				StealsAllowed:             team.TeamSeasonStats.StealsAllowed,
				StealsPerGame:             team.TeamSeasonStats.StealsPerGame,
				StealsAllowedPerGame:      team.TeamSeasonStats.StealsAllowedPerGame,
				StealsDiff:                team.TeamSeasonStats.StealsPerGame - team.TeamSeasonStats.StealsAllowedPerGame,
				Blocks:                    team.TeamSeasonStats.Blocks,
				BlocksAllowed:             team.TeamSeasonStats.BlocksAllowed,
				BlocksPerGame:             team.TeamSeasonStats.BlocksPerGame,
				BlocksAllowedPerGame:      team.TeamSeasonStats.BlocksAllowedPerGame,
				BlocksDiff:                team.TeamSeasonStats.BlocksPerGame - team.TeamSeasonStats.BlocksAllowedPerGame,
				TotalTurnovers:            team.TeamSeasonStats.TotalTurnovers,
				TurnoversAllowed:          team.TeamSeasonStats.TurnoversAllowed,
				TurnoversPerGame:          team.TeamSeasonStats.TurnoversPerGame,
				TurnoversAllowedPerGame:   team.TeamSeasonStats.TurnoversAllowedPerGame,
				TODiff:                    team.TeamSeasonStats.TurnoversPerGame - team.TeamSeasonStats.TurnoversAllowedPerGame,
				Fouls:                     team.TeamSeasonStats.Fouls,
				FoulsPerGame:              team.TeamSeasonStats.FoulsPerGame,
			}
		}

		var teamRes = structs.NBATeamResponse{
			ID:              team.ID,
			Team:            team.Team,
			Nickname:        team.Nickname,
			Abbr:            team.Abbr,
			LeagueID:        team.League,
			League:          team.League,
			ConferenceID:    team.ConferenceID,
			Conference:      team.Conference,
			DivisionID:      team.DivisionID,
			Division:        team.Division,
			Coach:           team.NBACoachName,
			OverallGrade:    team.OverallGrade,
			OffenseGrade:    team.OffenseGrade,
			DefenseGrade:    team.DefenseGrade,
			IsActive:        team.IsActive,
			IsInternational: team.LeagueID > 1,
			SeasonStats:     seasonsResponse,
			Stats:           teamStat,
		}

		ntResponse = append(ntResponse, teamRes)
	}
	return ntResponse
}

func GetTeamByTeamID(teamId string) structs.Team {
	var team structs.Team
	db := dbprovider.GetInstance().GetDB()
	// Preload("RecruitingProfile").
	err := db.Where("id = ?", teamId).Find(&team).Error
	if err != nil {
		log.Fatal(err)
	}
	return team
}

func GetNBATeamByTeamID(teamId string) structs.NBATeam {
	var team structs.NBATeam
	db := dbprovider.GetInstance().GetDB()
	// Preload("RecruitingProfile").
	err := db.Preload("Capsheet").Where("id = ?", teamId).Find(&team).Error
	if err != nil {
		log.Fatal(err)
	}
	return team
}

func RemoveUserFromTeam(teamId string) structs.Team {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	team := GetTeamByTeamID(teamId)

	team.RemoveUser()

	standings := GetStandingsRecordByTeamID(teamId, strconv.Itoa(int(ts.SeasonID)))

	standings.UpdateCoach("AI")

	recruitingProfile := GetOnlyTeamRecruitingProfileByTeamID(teamId)

	recruitingProfile.ToggleAIBehavior(true)

	db.Save(&team)

	db.Save(&standings)

	db.Save(&recruitingProfile)

	return team
}

func GetTeamsInConference(db *gorm.DB, conference string) []structs.Team {
	var teams []structs.Team
	err := db.Where("conference = ?", conference).Find(&teams).Error
	if err != nil {
		log.Fatal(err)
	}

	return teams
}

func GetTeamRatings(t structs.Team) {
	db := dbprovider.GetInstance().GetDB()
	teamIDINT := int(t.ID)

	players := GetCollegePlayersByTeamId(strconv.Itoa(teamIDINT))

	offenseRating := 0
	defenseRating := 0
	overallRating := 0

	offenseSum := 0
	defenseSum := 0

	for idx, player := range players {
		if idx > 9 {
			break
		}
		offenseSum += player.Shooting2 + player.Shooting3 + player.Finishing
		defenseSum += player.Ballwork + player.Rebounding + player.Defense
	}

	offenseRating = offenseSum / 9
	defenseRating = defenseSum / 9
	overallRating = (offenseRating + defenseRating) / 2

	offLetterGrade := util.GetOffenseGrade(offenseRating)
	defLetterGrade := util.GetDefenseGrade(defenseRating)
	ovrLetterGrade := util.GetOverallGrade(overallRating)

	t.AssignRatings(offLetterGrade, defLetterGrade, ovrLetterGrade)

	err := db.Save(&t).Error
	if err != nil {
		log.Fatalln("Could not save team rating for " + t.Abbr)
	}
}

func GetNBATeamRatings(t structs.NBATeam) {
	db := dbprovider.GetInstance().GetDB()
	teamIDINT := int(t.ID)

	players := GetNBAPlayersWithContractsByTeamID(strconv.Itoa(teamIDINT))

	offenseRating := 0
	defenseRating := 0
	overallRating := 0

	offenseSum := 0
	defenseSum := 0
	count := 0
	for _, player := range players {
		if player.IsGLeague {
			continue
		}
		if count > 9 {
			break
		}
		offenseSum += player.Shooting2 + player.Shooting3 + player.Finishing + player.FreeThrow
		defenseSum += player.Ballwork + player.Rebounding + player.InteriorDefense + player.PerimeterDefense
		count++
	}

	offenseRating = offenseSum / 9
	defenseRating = defenseSum / 9
	overallRating = (offenseRating + defenseRating) / 2

	offLetterGrade := util.GetNBATeamGrade(offenseRating)
	defLetterGrade := util.GetNBATeamGrade(defenseRating)
	ovrLetterGrade := util.GetNBATeamGrade(overallRating)

	t.AssignRatings(offLetterGrade, defLetterGrade, ovrLetterGrade)

	err := db.Save(&t).Error
	if err != nil {
		log.Fatalln("Could not save team rating for " + t.Abbr)
	}
}

func GetCBBTeamByAbbreviation(abbr string) structs.Team {
	var team structs.Team
	db := dbprovider.GetInstance().GetDB()
	err := db.Where("abbr = ?", abbr).Find(&team).Error
	if err != nil {
		log.Fatal(err)
	}
	return team
}

func GetOnlyNBATeams() []structs.NBATeam {
	db := dbprovider.GetInstance().GetDB()

	var teams []structs.NBATeam

	err := db.Where("league_id = 1").Find(&teams).Error
	if err != nil {
		log.Fatal(err)
	}
	return teams
}

func GetAllActiveNBATeams() []structs.NBATeam {
	db := dbprovider.GetInstance().GetDB()

	var teams []structs.NBATeam

	err := db.Find(&teams).Error
	if err != nil {
		log.Fatal(err)
	}
	return teams
}

// GetTeamByTeamID - straightforward
func GetNBATeamWithCapsheetByTeamID(teamId string) structs.NBATeam {
	var team structs.NBATeam
	db := dbprovider.GetInstance().GetDB()
	err := db.Preload("Capsheet").Where("id = ?", teamId).Find(&team).Error
	if err != nil {
		log.Fatal(err)
	}
	return team
}

func FormISLRosters() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	islTeams := GetAllActiveNBATeams()
	playerSignedMap := make(map[uint]bool)
	freeAgents := GetAllFreeAgents()
	sort.Slice(freeAgents, func(i, j int) bool {
		iVal := freeAgents[i].Overall
		jVal := freeAgents[j].Overall
		return iVal > jVal
	})
	maxRosterCount := 13
	currentCount := 2
	// Format Team Needs
	islTeamNeedsMap := make(map[uint]structs.ISLTeamNeeds)
	for _, t := range islTeams {
		if t.LeagueID == 1 {
			continue
		}
		teamID := strconv.Itoa(int(t.ID))
		teamNeedsMap := make(map[string]bool)
		positionCount := make(map[string]int)

		roster := GetAllNBAPlayersByTeamID(teamID)

		for _, r := range roster {
			positionCount[r.Position] += 1
		}

		if positionCount["PG"] < 3 {
			teamNeedsMap["PG"] = true
		} else if positionCount["SG"] < 4 {
			teamNeedsMap["SG"] = true
		} else if positionCount["SF"] < 4 {
			teamNeedsMap["SF"] = true
		} else if positionCount["PF"] < 4 {
			teamNeedsMap["PF"] = true
		} else if positionCount["C"] < 3 {
			teamNeedsMap["C"] = true
		}

		islTeamNeedsMap[t.ID] = structs.ISLTeamNeeds{
			TeamNeedsMap:  teamNeedsMap,
			PositionCount: positionCount,
			TotalCount:    len(roster),
		}
	}

	reverseOrder := islTeams
	sort.Slice(reverseOrder, func(i, j int) bool {
		iVal := reverseOrder[i].ID
		jVal := reverseOrder[j].ID
		return iVal > jVal
	})

	goReverse := false
	oneYear := false
	for currentCount < maxRosterCount {

		order := islTeams
		if goReverse {
			order = reverseOrder
		}

		for _, t := range order {
			teamNeeds := islTeamNeedsMap[t.ID]
			teamName := t.Team + " " + t.Nickname

			for _, fa := range freeAgents {
				if playerSignedMap[fa.ID] {
					continue
				}
				isSGSFPF := false
				if fa.Position == "SG" || fa.Position == "SF" || fa.Position == "PF" {
					isSGSFPF = true
				}
				if (teamNeeds.PositionCount[fa.Position] > 3 && isSGSFPF) || (teamNeeds.PositionCount[fa.Position] > 2 && !isSGSFPF) {
					continue
				}

				// Increase Position Count Limit
				teamNeeds.IncrementPositionCount(fa.Position)

				// Sign Player
				playerSignedMap[fa.ID] = true
				fa.SignWithTeam(t.ID, teamName)

				yearsOnContract := 1
				y1 := 0.7
				y2 := 0.0
				if !oneYear {
					yearsOnContract = 2
					y2 = 0.7
				}

				Contract := structs.NBAContract{
					PlayerID:       fa.PlayerID,
					TeamID:         t.ID,
					Team:           teamName,
					OriginalTeamID: t.ID,
					OriginalTeam:   teamName,
					YearsRemaining: uint(yearsOnContract),
					ContractType:   "Min",
					Year1Total:     y1,
					Year2Total:     y2,
					TotalRemaining: y1 + y2,
					IsActive:       true,
					IsComplete:     false,
					IsExtended:     false,
				}

				db.Create(&Contract)
				db.Save(&fa)

				// News Log
				message := "FA " + fa.Position + " " + fa.FirstName + " " + fa.LastName + " has signed with the ISL Team " + teamName + " with a contract worth approximately $" + strconv.Itoa(int(Contract.ContractValue)) + " Million Dollars."
				CreateNewsLog("NBA", message, "Free Agency", 0, ts)
			}

		}
		currentCount += 1
		oneYear = !oneYear
	}
}
