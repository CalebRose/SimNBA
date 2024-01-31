package managers

import (
	"log"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
)

// Player Controls
func GetCollegePlayerByNameAndAbbr(firstName, lastName, abbr string) structs.CollegePlayerResponse {
	db := dbprovider.GetInstance().GetDB()

	var player structs.CollegePlayer

	err := db.Where("first_name = ? and last_name = ? and team_id = ?", firstName, lastName, abbr).Find(&player).Error
	if err != nil {
		log.Fatal(err)
	}
	shooting2Grade := util.GetAttributeGrade(player.Shooting2)
	shooting3Grade := util.GetAttributeGrade(player.Shooting3)
	freeThrowGrade := util.GetAttributeGrade(player.FreeThrow)
	finishingGrade := util.GetAttributeGrade(player.Finishing)
	reboundingGrade := util.GetAttributeGrade(player.Rebounding)
	ballworkGrade := util.GetAttributeGrade(player.Ballwork)
	interiorDefenseGrade := util.GetAttributeGrade(player.InteriorDefense)
	perimeterDefenseGrade := util.GetAttributeGrade(player.PerimeterDefense)
	potentialGrade := util.GetPotentialGrade(player.Potential)
	overallGrade := util.GetPlayerOverallGrade(player.Overall)

	return structs.CollegePlayerResponse{
		FirstName:             player.FirstName,
		LastName:              player.LastName,
		Position:              player.Position,
		Archetype:             player.Archetype,
		Age:                   player.Age,
		Year:                  player.Year,
		State:                 player.State,
		Country:               player.Country,
		Stars:                 player.Stars,
		Height:                player.Height,
		PotentialGrade:        potentialGrade,
		Shooting2Grade:        shooting2Grade,
		Shooting3Grade:        shooting3Grade,
		FreeThrowGrade:        freeThrowGrade,
		FinishingGrade:        finishingGrade,
		BallworkGrade:         ballworkGrade,
		ReboundingGrade:       reboundingGrade,
		InteriorDefenseGrade:  interiorDefenseGrade,
		PerimeterDefenseGrade: perimeterDefenseGrade,
		OverallGrade:          overallGrade,
		Stamina:               player.Stamina,
		PlaytimeExpectations:  player.PlaytimeExpectations,
		Minutes:               player.Minutes,
		Potential:             player.Potential,
		Personality:           player.Personality,
		RecruitingBias:        player.RecruitingBias,
		WorkEthic:             player.WorkEthic,
		AcademicBias:          player.AcademicBias,
		PlayerID:              player.PlayerID,
		TeamID:                player.TeamID,
		TeamAbbr:              player.TeamAbbr,
		IsRedshirting:         player.IsRedshirting,
		IsRedshirt:            player.IsRedshirt,
		PositionOne:           player.PositionOne,
		PositionTwo:           player.PositionTwo,
		PositionThree:         player.PositionThree,
		P1Minutes:             player.P1Minutes,
		P2Minutes:             player.P2Minutes,
		P3Minutes:             player.P3Minutes,
		InsideProportion:      player.InsideProportion,
		MidRangeProportion:    player.MidRangeProportion,
		ThreePointProportion:  player.ThreePointProportion,
	}
}

func GetCollegeRecruitByNameAndLocation(firstName, lastName string) []structs.Croot {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.Recruit
	var croots []structs.Croot

	err := db.Where("first_name = ? and last_name = ?", firstName, lastName).Find(&players).Error
	if err != nil {
		log.Fatal(err)
	}
	for _, recruit := range players {
		var croot structs.Croot
		croot.Map(recruit)

		overallGrade := util.GetOverallGrade(recruit.Overall)

		croot.SetOverallGrade(overallGrade)

		croots = append(croots, croot)
	}

	return croots
}

func GetNBAPlayerByNameAndAbbr(firstName, lastName, abbr string) structs.NBAPlayer {
	db := dbprovider.GetInstance().GetDB()

	var player structs.NBAPlayer

	err := db.Where("first_name = ? and last_name = ? and team_id = ?", firstName, lastName, abbr).Find(&player).Error
	if err != nil {
		log.Fatal(err)
	}

	return player
}

// Team Controls
func GetCollegeTeamDataByID(id string) structs.CollegeTeamResponseData {
	ts := GetTimestamp()
	seasonId := strconv.Itoa(int(ts.SeasonID))

	team := GetTeamByTeamID(id)
	standings := GetStandingsRecordByTeamID(id, seasonId)
	matches := GetMatchesByTeamIdAndSeasonId(id, seasonId)
	wins := 0
	losses := 0
	confWins := 0
	confLosses := 0
	matchList := []structs.Match{}
	for _, m := range matches {
		if m.Week > uint(ts.CollegeWeek) {
			break
		}
		gameNotRan := (m.MatchOfWeek == "A" && !ts.GamesARan) ||
			(m.MatchOfWeek == "B" && !ts.GamesBRan) ||
			(m.MatchOfWeek == "C" && !ts.GamesCRan) ||
			(m.MatchOfWeek == "D" && !ts.GamesDRan)

		if (m.Week < uint(ts.CollegeWeek)) || ((strconv.Itoa(int(m.HomeTeamID)) == id && m.HomeTeamWin) ||
			(strconv.Itoa(int(m.AwayTeamID)) == id && m.AwayTeamWin)) && !gameNotRan {
			wins += 1
			if m.IsConference {
				confWins += 1
			}
		} else if (m.Week < uint(ts.CollegeWeek)) || ((strconv.Itoa(int(m.HomeTeamID)) == id && m.AwayTeamWin) ||
			(strconv.Itoa(int(m.AwayTeamID)) == id && m.HomeTeamWin)) && !gameNotRan {
			losses += 1
			if m.IsConference {
				confLosses += 1
			}
		}
		if gameNotRan {
			m.HideScore()
		}
		if m.Week == uint(ts.CollegeWeek) {
			matchList = append(matchList, m)
		}
	}

	standings.MaskGames(wins, losses, confWins, confLosses)

	return structs.CollegeTeamResponseData{
		TeamData:        team,
		TeamStandings:   standings,
		UpcomingMatches: matchList,
	}
}

func GetNBATeamDataByID(id string) structs.NBATeamResponseData {
	ts := GetTimestamp()
	seasonId := strconv.Itoa(int(ts.SeasonID))

	team := GetNBATeamByTeamID(id)
	standings := GetNBAStandingsRecordByTeamID(id, seasonId)
	matches := GetProfessionalMatchesByTeamIdAndSeasonId(id, seasonId)
	matchList := []structs.NBAMatch{}
	for _, m := range matches {
		if m.Week < uint(ts.NBAWeek) {
			continue
		}
		if m.Week > uint(ts.NBAWeek) {
			break
		}
		matchList = append(matchList, m)
	}

	return structs.NBATeamResponseData{
		TeamData:        team,
		TeamStandings:   standings,
		UpcomingMatches: matchList,
	}
}

// Stats

// Standings
func GetCollegeConferenceStandingsByConference(conf string) []structs.CollegeStandings {
	ts := GetTimestamp()
	seasonId := strconv.Itoa(int(ts.SeasonID))
	standings := GetConferenceStandingsByConferenceID(conf, seasonId)

	return standings
}

func GetNBAConferenceStandingsByConference(conf string) []structs.NBAStandings {
	ts := GetTimestamp()
	seasonId := strconv.Itoa(int(ts.SeasonID))
	standings := GetNBAConferenceStandingsByConferenceID(conf, seasonId)

	return standings
}

// Matches
func GetCollegeMatchesByConfAndDay(conf, day string) []structs.Match {
	ts := GetTimestamp()
	seasonId := strconv.Itoa(int(ts.SeasonID))

	teamMap := make(map[string]bool)

	standings := GetConferenceStandingsByConferenceID(conf, seasonId)

	for _, s := range standings {
		teamMap[s.TeamAbbr] = true
	}

	matches := GetCBBMatchesBySeasonID(seasonId)
	matchList := []structs.Match{}

	for _, m := range matches {
		if m.Week < uint(ts.CollegeWeek) {
			continue
		}
		if (teamMap[m.HomeTeam] || teamMap[m.AwayTeam]) && m.MatchOfWeek == day {
			matchList = append(matchList, m)
		}
		if m.Week > uint(ts.CollegeWeek) {
			break
		}
	}

	return matchList
}
