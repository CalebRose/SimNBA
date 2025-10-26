package managers

import (
	"log"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/repository"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
)

// Player Controls
func GetCollegePlayerByID(id string) structs.CollegePlayerResponse {
	db := dbprovider.GetInstance().GetDB()

	var player structs.CollegePlayer

	err := db.Where("id = ?", id).Find(&player).Error
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

func GetCollegeRecruitByID(id string) []structs.Croot {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.Recruit
	var croots []structs.Croot

	err := db.Where("id = ?", id).Find(&players).Error
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

func GetNBAPlayerByID(id string) structs.NBAPlayer {
	db := dbprovider.GetInstance().GetDB()

	var player structs.NBAPlayer

	err := db.Where("id = ?", id).Find(&player).Error
	if err != nil {
		log.Fatal(err)
	}

	return player
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

		earlierWeek := m.Week < uint(ts.CollegeWeek)

		if ((strconv.Itoa(int(m.HomeTeamID)) == id && m.HomeTeamWin) ||
			(strconv.Itoa(int(m.AwayTeamID)) == id && m.AwayTeamWin)) && (earlierWeek || !gameNotRan) {
			wins += 1
			if m.IsConference {
				confWins += 1
			}
		} else if ((strconv.Itoa(int(m.HomeTeamID)) == id && m.AwayTeamWin) ||
			(strconv.Itoa(int(m.AwayTeamID)) == id && m.HomeTeamWin)) && (earlierWeek || !gameNotRan) {
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
	wins := 0
	losses := 0
	confWins := 0
	confLosses := 0
	matchList := []structs.NBAMatch{}
	for _, m := range matches {
		if m.Week > uint(ts.NBAWeek) {
			break
		}
		gameNotRan := (m.MatchOfWeek == "A" && !ts.GamesARan) ||
			(m.MatchOfWeek == "B" && !ts.GamesBRan) ||
			(m.MatchOfWeek == "C" && !ts.GamesCRan) ||
			(m.MatchOfWeek == "D" && !ts.GamesDRan)

		earlierWeek := m.Week < uint(ts.CollegeWeek)

		if ((strconv.Itoa(int(m.HomeTeamID)) == id && m.HomeTeamWin) ||
			(strconv.Itoa(int(m.AwayTeamID)) == id && m.AwayTeamWin)) && (earlierWeek || !gameNotRan) {
			wins += 1
			if m.IsConference {
				confWins += 1
			}
		} else if ((strconv.Itoa(int(m.HomeTeamID)) == id && m.AwayTeamWin) ||
			(strconv.Itoa(int(m.AwayTeamID)) == id && m.HomeTeamWin)) && (earlierWeek || !gameNotRan) {
			losses += 1
			if m.IsConference {
				confLosses += 1
			}
		}
		if gameNotRan {
			m.HideScore()
		}
		matchList = append(matchList, m)
	}

	standings.MaskGames(wins, losses, confWins, confLosses)

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

func AssignDiscordIDToCollegeTeam(tID, dID string) {
	db := dbprovider.GetInstance().GetDB()

	team := GetTeamByTeamID(tID)

	team.AssignDiscordID(dID)

	repository.SaveCollegeTeamRecord(team, db)
}

func AssignDiscordIDToNBATeam(tID, dID, un string) {
	db := dbprovider.GetInstance().GetDB()

	team := GetNBATeamByTeamID(tID)

	team.AssignDiscordID(dID, un)

	repository.SaveNBATeamRecord(team, db)
}

func CompareTwoCBBTeams(t1ID, t2ID string) structs.FlexComparisonModel {
	ts := GetTimestamp()
	teamOneChan := make(chan structs.Team)
	teamTwoChan := make(chan structs.Team)

	go func() {
		t1 := GetTeamByTeamID(t1ID)
		teamOneChan <- t1
	}()

	teamOne := <-teamOneChan
	close(teamOneChan)

	go func() {
		t2 := GetTeamByTeamID(t2ID)
		teamTwoChan <- t2
	}()

	teamTwo := <-teamTwoChan
	close(teamTwoChan)

	allTeamOneGames := GetCBBMatchesByTeamId(t1ID)

	t1Wins := 0
	t1Losses := 0
	t1Streak := 0
	t1CurrentStreak := 0
	t1LargestMarginSeason := 0
	t1LargestMarginDiff := 0
	t1LargestMarginScore := ""
	t2Wins := 0
	t2Losses := 0
	t2Streak := 0
	t2CurrentStreak := 0
	latestWin := ""
	t2LargestMarginSeason := 0
	t2LargestMarginDiff := 0
	t2LargestMarginScore := ""

	for _, game := range allTeamOneGames {
		if !game.GameComplete ||
			(game.Week == uint(ts.NBAWeek) &&
				((game.TimeSlot == "A" && !ts.GamesARan) ||
					(game.TimeSlot == "B" && !ts.GamesBRan) ||
					(game.TimeSlot == "C" && !ts.GamesCRan) ||
					(game.TimeSlot == "D" && !ts.GamesDRan))) {
			continue
		}
		doComparison := (game.HomeTeamID == teamOne.ID && game.AwayTeamID == teamTwo.ID) ||
			(game.HomeTeamID == teamTwo.ID && game.AwayTeamID == teamOne.ID)
		if !doComparison {
			continue
		}
		homeTeamTeamOne := game.HomeTeamID == teamOne.ID
		if homeTeamTeamOne {
			if game.HomeTeamWin {
				t1Wins += 1
				t1CurrentStreak += 1
				latestWin = game.HomeTeam
				diff := game.HomeTeamScore - game.AwayTeamScore
				if diff > t1LargestMarginDiff {
					t1LargestMarginDiff = diff
					t1LargestMarginSeason = int(game.SeasonID) + 2020
					t1LargestMarginScore = "" + strconv.Itoa(game.HomeTeamScore) + "-" + strconv.Itoa(game.AwayTeamScore)
				}
			} else {
				t1Streak = t1CurrentStreak
				t1CurrentStreak = 0
				t1Losses += 1
			}
		} else {
			if game.HomeTeamWin {
				t2Wins += 1
				t2CurrentStreak += 1
				latestWin = game.HomeTeam
				diff := game.HomeTeamScore - game.AwayTeamScore
				if diff > t2LargestMarginDiff {
					t2LargestMarginDiff = diff
					t2LargestMarginSeason = int(game.SeasonID) + 2020
					t2LargestMarginScore = "" + strconv.Itoa(game.HomeTeamScore) + "-" + strconv.Itoa(game.AwayTeamScore)
				}
			} else {
				t2Streak = t2CurrentStreak
				t2CurrentStreak = 0
				t2Losses += 1
			}
		}

		awayTeamTeamOne := game.AwayTeamID == teamOne.ID
		if awayTeamTeamOne {
			if game.AwayTeamWin {
				t1Wins += 1
				t1CurrentStreak += 1
				latestWin = game.AwayTeam
				diff := game.AwayTeamScore - game.HomeTeamScore
				if diff > t1LargestMarginDiff {
					t1LargestMarginDiff = diff
					t1LargestMarginSeason = int(game.SeasonID) + 2020
					t1LargestMarginScore = "" + strconv.Itoa(game.AwayTeamScore) + "-" + strconv.Itoa(game.HomeTeamScore)
				}
			} else {
				t1Streak = t1CurrentStreak
				t1CurrentStreak = 0
				t1Losses += 1
			}
		} else {
			if game.AwayTeamWin {
				t2Wins += 1
				t2CurrentStreak += 1
				latestWin = game.AwayTeam
				diff := game.AwayTeamScore - game.HomeTeamScore
				if diff > t2LargestMarginDiff {
					t2LargestMarginDiff = diff
					t2LargestMarginSeason = int(game.SeasonID) + 2020
					t2LargestMarginScore = "" + strconv.Itoa(game.AwayTeamScore) + "-" + strconv.Itoa(game.HomeTeamScore)
				}
			} else {
				t2Streak = t2CurrentStreak
				t2CurrentStreak = 0
				t2Losses += 1
			}
		}
	}

	if t1CurrentStreak > 0 && t1CurrentStreak > t1Streak {
		t1Streak = t1CurrentStreak
	}
	if t2CurrentStreak > 0 && t2CurrentStreak > t2Streak {
		t2Streak = t2CurrentStreak
	}

	currentStreak := 0
	currentStreak = max(t1CurrentStreak, t2CurrentStreak)

	return structs.FlexComparisonModel{
		TeamOneID:      teamOne.ID,
		TeamOne:        teamOne.Abbr,
		TeamOneWins:    uint(t1Wins),
		TeamOneLosses:  uint(t1Losses),
		TeamOneStreak:  uint(t1Streak),
		TeamOneMSeason: t1LargestMarginSeason,
		TeamOneMScore:  t1LargestMarginScore,
		TeamTwoID:      teamTwo.ID,
		TeamTwo:        teamTwo.Abbr,
		TeamTwoWins:    uint(t2Wins),
		TeamTwoLosses:  uint(t2Losses),
		TeamTwoStreak:  uint(t2Streak),
		TeamTwoMSeason: t2LargestMarginSeason,
		TeamTwoMScore:  t2LargestMarginScore,
		CurrentStreak:  uint(currentStreak),
		LatestWin:      latestWin,
	}
}

func CompareTwoNBATeams(t1ID, t2ID string) structs.FlexComparisonModel {
	ts := GetTimestamp()
	teamOneChan := make(chan structs.NBATeam)
	teamTwoChan := make(chan structs.NBATeam)

	go func() {
		t1 := GetNBATeamByTeamID(t1ID)
		teamOneChan <- t1
	}()

	teamOne := <-teamOneChan
	close(teamOneChan)

	go func() {
		t2 := GetNBATeamByTeamID(t2ID)
		teamTwoChan <- t2
	}()

	teamTwo := <-teamTwoChan
	close(teamTwoChan)

	allTeamOneGames := GetNBAMatchesByTeamId(t1ID)

	t1Wins := 0
	t1Losses := 0
	t1Streak := 0
	t1CurrentStreak := 0
	t1LargestMarginSeason := 0
	t1LargestMarginDiff := 0
	t1LargestMarginScore := ""
	t2Wins := 0
	t2Losses := 0
	t2Streak := 0
	t2CurrentStreak := 0
	latestWin := ""
	t2LargestMarginSeason := 0
	t2LargestMarginDiff := 0
	t2LargestMarginScore := ""

	for _, game := range allTeamOneGames {
		if !game.GameComplete ||
			(game.Week == uint(ts.NBAWeek) &&
				((game.TimeSlot == "A" && !ts.GamesARan) ||
					(game.TimeSlot == "B" && !ts.GamesBRan) ||
					(game.TimeSlot == "C" && !ts.GamesCRan) ||
					(game.TimeSlot == "D" && !ts.GamesDRan))) {
			continue
		}
		doComparison := (game.HomeTeamID == teamOne.ID && game.AwayTeamID == teamTwo.ID) ||
			(game.HomeTeamID == teamTwo.ID && game.AwayTeamID == teamOne.ID)

		if !doComparison {
			continue
		}
		homeTeamTeamOne := game.HomeTeamID == teamOne.ID
		if homeTeamTeamOne {
			if game.HomeTeamWin {
				t1Wins += 1
				t1CurrentStreak += 1
				latestWin = game.HomeTeam
				diff := game.HomeTeamScore - game.AwayTeamScore
				if diff > t1LargestMarginDiff {
					t1LargestMarginDiff = diff
					t1LargestMarginSeason = int(game.SeasonID) + 2020
					t1LargestMarginScore = "" + strconv.Itoa(game.HomeTeamScore) + "-" + strconv.Itoa(game.AwayTeamScore)
				}
			} else {
				t1Streak = t1CurrentStreak
				t1CurrentStreak = 0
				t1Losses += 1
			}
		} else {
			if game.HomeTeamWin {
				t2Wins += 1
				t2CurrentStreak += 1
				latestWin = game.HomeTeam
				diff := game.HomeTeamScore - game.AwayTeamScore
				if diff > t2LargestMarginDiff {
					t2LargestMarginDiff = diff
					t2LargestMarginSeason = int(game.SeasonID) + 2020
					t2LargestMarginScore = "" + strconv.Itoa(game.HomeTeamScore) + "-" + strconv.Itoa(game.AwayTeamScore)
				}
			} else {
				t2Streak = t2CurrentStreak
				t2CurrentStreak = 0
				t2Losses += 1
			}
		}

		awayTeamTeamOne := game.AwayTeamID == teamOne.ID
		if awayTeamTeamOne {
			if game.AwayTeamWin {
				t1Wins += 1
				t1CurrentStreak += 1
				latestWin = game.AwayTeam
				diff := game.AwayTeamScore - game.HomeTeamScore
				if diff > t1LargestMarginDiff {
					t1LargestMarginDiff = diff
					t1LargestMarginSeason = int(game.SeasonID) + 2020
					t1LargestMarginScore = "" + strconv.Itoa(game.AwayTeamScore) + "-" + strconv.Itoa(game.HomeTeamScore)
				}
			} else {
				t1Streak = t1CurrentStreak
				t1CurrentStreak = 0
				t1Losses += 1
			}
		} else {
			if game.AwayTeamWin {
				t2Wins += 1
				t2CurrentStreak += 1
				latestWin = game.AwayTeam
				diff := game.AwayTeamScore - game.HomeTeamScore
				if diff > t2LargestMarginDiff {
					t2LargestMarginDiff = diff
					t2LargestMarginSeason = int(game.SeasonID) + 2020
					t2LargestMarginScore = "" + strconv.Itoa(game.AwayTeamScore) + "-" + strconv.Itoa(game.HomeTeamScore)
				}
			} else {
				t2Streak = t2CurrentStreak
				t2CurrentStreak = 0
				t2Losses += 1
			}
		}
	}

	if t1CurrentStreak > 0 && t1CurrentStreak > t1Streak {
		t1Streak = t1CurrentStreak
	}
	if t2CurrentStreak > 0 && t2CurrentStreak > t2Streak {
		t2Streak = t2CurrentStreak
	}

	currentStreak := 0
	currentStreak = max(t1CurrentStreak, t2CurrentStreak)

	return structs.FlexComparisonModel{
		TeamOneID:      teamOne.ID,
		TeamOne:        teamOne.Abbr,
		TeamOneWins:    uint(t1Wins),
		TeamOneLosses:  uint(t1Losses),
		TeamOneStreak:  uint(t1Streak),
		TeamOneMSeason: t1LargestMarginSeason,
		TeamOneMScore:  t1LargestMarginScore,
		TeamTwoID:      teamTwo.ID,
		TeamTwo:        teamTwo.Abbr,
		TeamTwoWins:    uint(t2Wins),
		TeamTwoLosses:  uint(t2Losses),
		TeamTwoStreak:  uint(t2Streak),
		TeamTwoMSeason: t2LargestMarginSeason,
		TeamTwoMScore:  t2LargestMarginScore,
		CurrentStreak:  uint(currentStreak),
		LatestWin:      latestWin,
	}
}
