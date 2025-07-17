package managers

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/repository"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
	"gorm.io/gorm"
)

// UpdateGameplan -- Need to update
func UpdateGameplan(updateGameplanDto structs.UpdateGameplanDto) {
	db := dbprovider.GetInstance().GetDB()

	var teamId = strconv.Itoa(updateGameplanDto.TeamID)

	// Get Gameplans
	var gameplan = GetGameplansByTeam(teamId)

	ug := updateGameplanDto.Gameplan

	// If no changes made to gameplan

	// Otherwise, update the gameplan
	gameplan.UpdateGameplan(ug.Pace, ug.OffensiveFormation, ug.DefensiveFormation, ug.OffensiveStyle, ug.FocusPlayer)
	gameplan.UpdateToggles(ug.Toggle2pt, ug.Toggle3pt, ug.ToggleFN, ug.ToggleFT, ug.ToggleBW, ug.ToggleRB, ug.ToggleID, ug.TogglePD, ug.ToggleP2, ug.ToggleP3)
	fmt.Printf("%s", "Saving Gameplan for Team "+teamId+"\n")
	db.Save(&gameplan)

	// Get Players
	updatedPlayers := updateGameplanDto.CollegePlayers

	for _, player := range updatedPlayers {
		id := strconv.Itoa(int(player.PlayerID))
		record := GetCollegePlayerByPlayerID(id)
		record.UpdatePlayer(player.P1Minutes, player.P2Minutes, player.P3Minutes, player.PositionOne, player.PositionTwo, player.PositionThree, player.InsideProportion, player.MidRangeProportion, player.ThreePointProportion)
		db.Save(&record)
	}
}

func UpdateNBAGameplan(updateGameplanDto structs.UpdateGameplanDto) {
	db := dbprovider.GetInstance().GetDB()

	var teamId = strconv.Itoa(updateGameplanDto.TeamID)

	// Get Gameplans
	var gameplan = GetNBAGameplanByTeam(teamId)

	ug := updateGameplanDto.Gameplan

	// If no changes made to gameplan

	// Otherwise, update the gameplan
	gameplan.UpdateGameplan(ug.Pace, ug.OffensiveFormation, ug.DefensiveFormation, ug.OffensiveStyle, ug.FocusPlayer)
	gameplan.UpdateToggles(ug.Toggle2pt, ug.Toggle3pt, ug.ToggleFN, ug.ToggleFT, ug.ToggleBW, ug.ToggleRB, ug.ToggleID, ug.TogglePD, ug.ToggleP2, ug.ToggleP3)
	fmt.Printf("%s", "Saving Gameplan for Team "+teamId+"\n")
	db.Save(&gameplan)

	// Get Players
	updatedPlayers := updateGameplanDto.NBAPlayers

	for _, player := range updatedPlayers {
		id := strconv.Itoa(int(player.ID))
		record := GetNBAPlayerRecord(id)
		record.UpdatePlayer(player.P1Minutes, player.P2Minutes, player.P3Minutes, player.PositionOne, player.PositionTwo, player.PositionThree, player.InsideProportion, player.MidRangeProportion, player.ThreePointProportion)
		db.Save(&record)
	}
}

func GetAllCollegeGameplans() []structs.Gameplan {
	db := dbprovider.GetInstance().GetDB()

	var gameplans []structs.Gameplan

	db.Find(&gameplans)

	return gameplans
}

// GetGameplansByTeam
func GetGameplansByTeam(teamId string) structs.Gameplan {
	db := dbprovider.GetInstance().GetDB()

	var gameplans structs.Gameplan

	db.Where("team_id = ?", teamId).Order("game asc").Find(&gameplans)

	return gameplans
}

// GetGameplansByTeam
func GetNBAGameplanByTeam(teamId string) structs.NBAGameplan {
	db := dbprovider.GetInstance().GetDB()

	var gameplans structs.NBAGameplan

	db.Where("team_id = ?", teamId).Order("game asc").Find(&gameplans)

	return gameplans
}

func GetOpposingCollegiateTeamRoster(teamID string) []structs.CollegePlayer {
	ts := GetTimestamp()

	matches := GetCollegeTeamMatchesBySeasonId(strconv.Itoa(int(ts.SeasonID)), teamID)
	designatedMatch := structs.Match{}

	for _, m := range matches {
		if m.GameComplete {
			continue
		}
		designatedMatch = m
		break
	}

	opposingTeamID := ""
	if teamID == strconv.Itoa(int(designatedMatch.HomeTeamID)) {
		opposingTeamID = strconv.Itoa(int(designatedMatch.AwayTeamID))
	} else {
		opposingTeamID = strconv.Itoa(int(designatedMatch.HomeTeamID))
	}

	if opposingTeamID == "0" {
		return []structs.CollegePlayer{}
	}

	opposingRoster := GetCollegePlayersByTeamId(opposingTeamID)

	return opposingRoster
}

func GetOpposingNBATeamRoster(teamID string) []structs.NBAPlayer {
	ts := GetTimestamp()

	matches := GetNBATeamMatchesBySeasonId(strconv.Itoa(int(ts.SeasonID)), teamID)
	designatedMatch := structs.NBAMatch{}

	for _, m := range matches {
		if m.GameComplete {
			continue
		}
		designatedMatch = m
		break
	}

	opposingTeamID := ""
	if teamID == strconv.Itoa(int(designatedMatch.HomeTeamID)) {
		opposingTeamID = strconv.Itoa(int(designatedMatch.AwayTeamID))
	} else {
		opposingTeamID = strconv.Itoa(int(designatedMatch.HomeTeamID))
	}

	if opposingTeamID == "0" {
		return []structs.NBAPlayer{}
	}

	opposingRoster := GetAllNBAPlayersByTeamID(opposingTeamID)

	return opposingRoster
}

func CheckAllUserGameplans() {
	db := dbprovider.GetInstance().GetDB()
	teams := GetAllActiveCollegeTeams()
	collegeTeamMap := MakeCollegeTeamMap(teams)
	collegePlayers := GetAllCollegePlayers()
	collegePlayerMapByTeamID := MakeCollegePlayerMapByTeamID(collegePlayers, true)
	gameplans := GetAllCollegeGameplans()
	gameplanMap := MakeCollegeGameplanMap(gameplans)
	teamIDsToCheck := []uint{}
	teamIDStrings := []string{}
	for _, team := range teams {
		if !team.IsUserCoached {
			continue
		}
		totalMinutes := 0
		insideShotAllocation := 0.0
		midRangeShotAllocation := 0.0
		threePointShotAllocation := 0.0
		pgMinutes := 0
		sgMinutes := 0
		sfMinutes := 0
		pfMinutes := 0
		cMinutes := 0
		roster := collegePlayerMapByTeamID[team.ID]

		for _, player := range roster {
			if player.IsRedshirting || player.IsInjured {
				continue
			}
			positionOne := player.PositionOne
			positionTwo := player.PositionTwo
			positionThree := player.PositionThree
			if positionOne == "" && positionTwo == "" && positionThree == "" {
				continue
			}

			positions := []struct {
				pos     string
				minutes int
			}{
				{positionOne, player.P1Minutes},
				{positionTwo, player.P2Minutes},
				{positionThree, player.P3Minutes},
			}

			for _, p := range positions {
				switch p.pos {
				case "PG":
					pgMinutes += p.minutes
				case "SG":
					sgMinutes += p.minutes
				case "SF":
					sfMinutes += p.minutes
				case "PF":
					pfMinutes += p.minutes
				case "C":
					cMinutes += p.minutes
				}
			}

			insideShotAllocation += player.InsideProportion
			midRangeShotAllocation += player.MidRangeProportion
			threePointShotAllocation += player.ThreePointProportion
			totalMinutes += player.Minutes
		}
		totalShotAllocation := insideShotAllocation + midRangeShotAllocation + threePointShotAllocation

		if totalMinutes <= 199 || totalShotAllocation < 100 {
			message := fmt.Sprintf("Your team did not properly gameplan: %d\n of 200 Minutes were allocated across your roster", totalMinutes)
			CreateNotification("CBB", message, "Gameplan", team.ID)
			teamIDsToCheck = append(teamIDsToCheck, team.ID)
			teamIDStrings = append(teamIDStrings, strconv.Itoa(int(team.ID)))
		}
	}

	db.Model(&structs.CollegePlayer{}).Where("team_id in (?)", teamIDStrings).Updates(map[string]interface{}{
		"minutes":                0,
		"p1_minutes":             0,
		"p2_minutes":             0,
		"p3_minutes":             0,
		"position_one":           "",
		"position_two":           "",
		"position_three":         "",
		"inside_proportion":      0,
		"mid_range_proportion":   0,
		"three_point_proportion": 0,
	})

	for _, teamID := range teamIDsToCheck {
		team := collegeTeamMap[teamID]
		team.MarkTeamPenalty()
		repository.SaveCollegeTeamRecord(team, db)
		SetCollegeMinutesAndShotProportions(db, teamID, gameplanMap, collegePlayerMapByTeamID)
	}
}

func SetAIGameplans() bool {
	db := dbprovider.GetInstance().GetDB()

	teams := GetAllActiveCollegeTeams()
	teamIDs := []string{}
	collegePlayers := GetAllCollegePlayers()
	collegePlayerMapByTeamID := MakeCollegePlayerMapByTeamID(collegePlayers, true)
	gameplans := GetAllCollegeGameplans()
	gameplanMap := MakeCollegeGameplanMap(gameplans)
	for _, team := range teams {
		if team.IsUserCoached || team.ID != 60 {
			continue
		}
		teamID := strconv.Itoa(int(team.ID))
		teamIDs = append(teamIDs, teamID)
	}

	// Update all players on AI team to reset minutes
	db.Model(&structs.CollegePlayer{}).Where("team_id in (?)", teamIDs).Updates(map[string]interface{}{
		"minutes":                0,
		"p1_minutes":             0,
		"p2_minutes":             0,
		"p3_minutes":             0,
		"position_one":           "",
		"position_two":           "",
		"position_three":         "",
		"inside_proportion":      0,
		"mid_range_proportion":   0,
		"three_point_proportion": 0,
	})

	// Update Minutes here

	for _, team := range teams {
		if team.IsUserCoached {
			continue
		}

		SetCollegeMinutesAndShotProportions(db, team.ID, gameplanMap, collegePlayerMapByTeamID)
	}

	islTeams := GetAllActiveNBATeams()
	islIDs := []string{}
	for _, team := range islTeams {
		if len(team.NBAOwnerName) > 0 {
			continue
		}
		teamID := strconv.Itoa(int(team.ID))
		islIDs = append(islIDs, teamID)
	}

	db.Model(&structs.NBAPlayer{}).Where("team_id in (?)", islIDs).Updates(map[string]interface{}{
		"minutes":                0,
		"p1_minutes":             0,
		"p2_minutes":             0,
		"p3_minutes":             0,
		"position_one":           "",
		"position_two":           "",
		"position_three":         "",
		"inside_proportion":      0,
		"mid_range_proportion":   0,
		"three_point_proportion": 0,
	})

	for _, team := range islTeams {
		if !team.IsActive {
			continue
		}

		if len(team.NBAOwnerName) > 0 {
			continue
		}

		pgCount := 0
		sgCount := 0
		sfCount := 0
		pfCount := 0
		cCount := 0
		pgMinutes := 0
		sgMinutes := 0
		sfMinutes := 0
		pfMinutes := 0
		cMinutes := 0

		pgList := []structs.NBAPlayer{}
		sgList := []structs.NBAPlayer{}
		sfList := []structs.NBAPlayer{}
		pfList := []structs.NBAPlayer{}
		cList := []structs.NBAPlayer{}

		gameplan := GetNBAGameplanByTeam(strconv.Itoa(int(team.ID)))
		off := "Balanced"
		def := "Man-to-Man"
		ost := ""
		pace := "Balanced"

		roster := GetOnlyNBAPlayersByTeamID(strconv.Itoa(int(team.ID)))
		rMap := make(map[string]*structs.NBAPlayer)
		for i := 0; i < len(roster); i++ {
			id := strconv.Itoa(int(roster[i].ID))
			rMap[id] = &roster[i]
		}

		for _, c := range roster {
			if c.IsGLeague {
				continue
			}

			if c.Position == "PG" {
				pgCount++
				pgList = append(pgList, c)
				sgList = append(sgList, c)
			} else if c.Position == "SG" {
				sgCount++
				sgList = append(sgList, c)
				pgList = append(pgList, c)
				sfList = append(sfList, c)
			} else if c.Position == "SF" {
				sfCount++
				sfList = append(sfList, c)
				sgList = append(sgList, c)
				pfList = append(pfList, c)
			} else if c.Position == "PF" {
				pfCount++
				pfList = append(pfList, c)
				sfList = append(sfList, c)
				cList = append(cList, c)
			} else if c.Position == "C" {
				cCount++
				cList = append(cList, c)
				pfList = append(pfList, c)
			}
		}

		if pgCount <= 2 && sgCount < 4 {
			ost = "Jumbo"
		} else if cCount <= 2 && pfCount < 4 {
			ost = util.PickFromStringList([]string{"Small Ball", "Microball"})
		} else {
			ost = "Traditional"
		}

		if ost == "Traditional" {
			pgMinutes = 48
			sgMinutes = 48
			pfMinutes = 48
			sfMinutes = 48
			cMinutes = 48
		} else if ost == "Small Ball" {
			pgMinutes = 48
			sgMinutes = 96
			pfMinutes = 48
			sfMinutes = 48
			cMinutes = 0
		} else if ost == "Microball" {
			pgMinutes = 96
			sgMinutes = 96
			pfMinutes = 0
			sfMinutes = 48
			cMinutes = 0
		} else if ost == "Jumbo" {
			pgMinutes = 0
			sgMinutes = 48
			pfMinutes = 96
			sfMinutes = 48
			cMinutes = 48
		}
		sort.Slice(pgList, func(i, j int) bool {
			return pgList[i].Overall > pgList[j].Overall
		})

		sort.Slice(sgList, func(i, j int) bool {
			return sgList[i].Overall > sgList[j].Overall
		})

		sort.Slice(sfList, func(i, j int) bool {
			return sfList[i].Overall > sfList[j].Overall
		})

		sort.Slice(pfList, func(i, j int) bool {
			return pfList[i].Overall > pfList[j].Overall
		})

		sort.Slice(cList, func(i, j int) bool {
			return cList[i].Overall > cList[j].Overall
		})
		totalMinutes := 0
		if ost == "Traditional" {
			totalMinutes += setNBAPositionMinutes(pgList, rMap, pgMinutes, "PG", ost)
			totalMinutes += setNBAPositionMinutes(sgList, rMap, sgMinutes, "SG", ost)
			totalMinutes += setNBAPositionMinutes(sfList, rMap, sfMinutes, "SF", ost)
			totalMinutes += setNBAPositionMinutes(pfList, rMap, pfMinutes, "PF", ost)
			totalMinutes += setNBAPositionMinutes(cList, rMap, cMinutes, "C", ost)
		} else if ost == "Jumbo" {
			totalMinutes += setNBAPositionMinutes(cList, rMap, cMinutes, "C", ost)
			totalMinutes += setNBAPositionMinutes(pfList, rMap, pfMinutes, "PF", ost)
			totalMinutes += setNBAPositionMinutes(sfList, rMap, sfMinutes, "SF", ost)
			totalMinutes += setNBAPositionMinutes(sgList, rMap, sgMinutes, "SG", ost)
		} else if ost == "Small Ball" {
			totalMinutes += setNBAPositionMinutes(sgList, rMap, sgMinutes, "SG", ost)
			totalMinutes += setNBAPositionMinutes(pgList, rMap, pgMinutes, "PG", ost)
			totalMinutes += setNBAPositionMinutes(sfList, rMap, sfMinutes, "SF", ost)
			totalMinutes += setNBAPositionMinutes(pfList, rMap, pfMinutes, "PF", ost)
		} else if ost == "Microball" {
			totalMinutes += setNBAPositionMinutes(pgList, rMap, pgMinutes, "PG", ost)
			totalMinutes += setNBAPositionMinutes(sgList, rMap, sgMinutes, "SG", ost)
			totalMinutes += setNBAPositionMinutes(sfList, rMap, sfMinutes, "SF", ost)
		}

		// For testing purposes
		teamMidRangeProportion := 0.0
		teamMidrangeLimit := 40.0
		teamInsideProportion := 0.0
		teamInsideLimit := 40.0
		teamThreePointProportion := 0.0
		teamThreePointLimit := 20.0

		sort.Slice(roster, func(i, j int) bool {
			return roster[i].Minutes > roster[j].Minutes && roster[i].Overall > roster[j].Overall
		})

		teamTotalSkill := 0
		for i := 0; i < len(roster); i++ {
			if roster[i].Minutes == 0 || roster[i].IsGLeague {
				continue
			}
			teamTotalSkill += roster[i].Shooting2 + roster[i].Shooting3 + roster[i].Finishing
		}

		// Loop for team shot proportions
		for i := 0; i < len(roster); i++ {
			if roster[i].Minutes == 0 || roster[i].IsGLeague {
				continue
			}
			totalSkill := roster[i].Shooting2 + roster[i].Shooting3 + roster[i].Finishing
			twoPointPercentage := float64(roster[i].Shooting2*100) / float64(totalSkill) * float64(roster[i].Minutes) / float64(roster[i].Stamina)
			threePointPercentage := float64(roster[i].Shooting3*100) / float64(totalSkill) * float64(roster[i].Minutes) / float64(roster[i].Stamina)
			insidePercentage := float64(roster[i].Finishing*100) / float64(totalSkill) * float64(roster[i].Minutes) / float64(roster[i].Stamina)
			teamInsideProportion += insidePercentage
			roster[i].SetInsideProportion(insidePercentage)
			teamMidRangeProportion += twoPointPercentage
			roster[i].SetMidShotProportion(twoPointPercentage)
			teamThreePointProportion += threePointPercentage
			roster[i].SetThreePointProportion(threePointPercentage)
		}

		insideProp := 0.0
		midProp := 0.0
		tpProp := 0.0

		// Motion
		if float64(teamThreePointProportion/teamMidRangeProportion) > 1.3 && float64(teamInsideProportion/teamMidRangeProportion) > 1.3 {
			off = "Motion"
			teamInsideLimit = 20
			teamMidrangeLimit = 10
			teamThreePointLimit = 70
			// Pick-And-Roll
		} else if float64(teamInsideProportion/teamMidRangeProportion) > 1.3 && float64(teamInsideProportion/teamThreePointProportion) > 1.3 {
			off = "Pick-and-Roll"
			teamInsideLimit = 40
			teamMidrangeLimit = 20
			teamThreePointLimit = 40
			// Post-Up
		} else if float64(teamInsideProportion/teamMidRangeProportion) > 1.5 && float64(teamInsideProportion/teamThreePointProportion) > 1.5 {
			off = "Post-Up"
			teamInsideLimit = 80
			teamMidrangeLimit = 15
			teamThreePointLimit = 5
			// Space-And-Post
		} else if float64(teamMidRangeProportion/teamInsideProportion) > 1.3 && float64(teamThreePointProportion/teamInsideProportion) > 1.3 {
			off = "Space-and-Post"
			teamInsideLimit = 20
			teamMidrangeLimit = 40
			teamThreePointLimit = 40
		}

		for i := 0; i < len(roster); i++ {
			if roster[i].Minutes == 0 || roster[i].IsGLeague {
				continue
			}
			normalizedInsideProportion := (roster[i].InsideProportion * float64(teamInsideLimit)) / teamInsideProportion
			insideProp += normalizedInsideProportion
			if insideProp > teamInsideLimit {
				diff := insideProp - teamInsideLimit
				insideProp -= diff
				normalizedInsideProportion -= diff
			}
			roster[i].SetInsideProportion(normalizedInsideProportion)

			normalizedMidrangeProportion := (roster[i].MidRangeProportion * float64(teamMidrangeLimit)) / teamMidRangeProportion
			midProp += normalizedMidrangeProportion
			if midProp > teamMidrangeLimit {
				diff := midProp - teamMidrangeLimit
				midProp -= diff
				normalizedMidrangeProportion -= diff
			}
			roster[i].SetMidShotProportion(normalizedMidrangeProportion)

			normalized3ptProportion := (roster[i].ThreePointProportion * float64(teamThreePointLimit)) / teamThreePointProportion
			tpProp += normalized3ptProportion
			if tpProp > teamThreePointLimit {
				diff := tpProp - teamThreePointLimit
				tpProp -= diff
				normalized3ptProportion -= diff
			}
			roster[i].SetThreePointProportion(normalized3ptProportion)
		}

		for _, r := range roster {
			db.Save(&r)
		}

		gameplan.UpdateGameplan(pace, off, def, ost, "")

		db.Save(&gameplan)
	}

	return true
}

func SetCollegeMinutesAndShotProportions(db *gorm.DB, teamID uint, gameplanMap map[uint]structs.Gameplan, collegePlayerMapByTeamID map[uint][]structs.CollegePlayer) {

	pgCount := 0
	sgCount := 0
	sfCount := 0
	pfCount := 0
	cCount := 0
	pgMinutes := 0
	sgMinutes := 0
	sfMinutes := 0
	pfMinutes := 0
	cMinutes := 0

	pgList := []structs.CollegePlayer{}
	sgList := []structs.CollegePlayer{}
	sfList := []structs.CollegePlayer{}
	pfList := []structs.CollegePlayer{}
	cList := []structs.CollegePlayer{}

	gameplan := gameplanMap[teamID]
	off := "Balanced"
	def := "Man-to-Man"
	ost := ""
	pace := "Balanced"

	roster := collegePlayerMapByTeamID[teamID]
	rMap := make(map[string]*structs.CollegePlayer)
	for i := 0; i < len(roster); i++ {
		id := strconv.Itoa(int(roster[i].ID))
		rMap[id] = &roster[i]
	}

	for _, c := range roster {
		if c.IsRedshirting || c.IsInjured {
			continue
		}

		if c.Position == "PG" {
			pgCount++
			pgList = append(pgList, c)
			sgList = append(sgList, c)
		} else if c.Position == "SG" {
			sgCount++
			sgList = append(sgList, c)
			pgList = append(pgList, c)
			sfList = append(sfList, c)
		} else if c.Position == "SF" {
			sfCount++
			sfList = append(sfList, c)
			sgList = append(sgList, c)
			pfList = append(pfList, c)
		} else if c.Position == "PF" {
			pfCount++
			pfList = append(pfList, c)
			sfList = append(sfList, c)
			cList = append(cList, c)
		} else if c.Position == "C" {
			cCount++
			cList = append(cList, c)
			pfList = append(pfList, c)
		}
	}

	if pgCount <= 2 && sgCount < 4 {
		ost = "Jumbo"
	} else if cCount <= 2 && pfCount < 4 {
		ost = util.PickFromStringList([]string{"Small Ball", "Microball"})
	} else {
		ost = "Traditional"
	}

	if ost == "Traditional" {
		pgMinutes = 40
		sgMinutes = 40
		pfMinutes = 40
		sfMinutes = 40
		cMinutes = 40
	} else if ost == "Small Ball" {
		pgMinutes = 40
		sgMinutes = 80
		pfMinutes = 40
		sfMinutes = 40
		cMinutes = 0
	} else if ost == "Microball" {
		pgMinutes = 80
		sgMinutes = 80
		pfMinutes = 00
		sfMinutes = 40
		cMinutes = 0
	} else if ost == "Jumbo" {
		pgMinutes = 0
		sgMinutes = 40
		pfMinutes = 80
		sfMinutes = 40
		cMinutes = 40
	}
	sort.Slice(pgList, func(i, j int) bool {
		return pgList[i].Overall > pgList[j].Overall
	})

	sort.Slice(sgList, func(i, j int) bool {
		return sgList[i].Overall > sgList[j].Overall
	})

	sort.Slice(sfList, func(i, j int) bool {
		return sfList[i].Overall > sfList[j].Overall
	})

	sort.Slice(pfList, func(i, j int) bool {
		return pfList[i].Overall > pfList[j].Overall
	})

	sort.Slice(cList, func(i, j int) bool {
		return cList[i].Overall > cList[j].Overall
	})
	totalMinutes := 0
	if ost == "Traditional" {
		totalMinutes += setPositionMinutes(pgList, rMap, pgMinutes, "PG", ost)
		totalMinutes += setPositionMinutes(sgList, rMap, sgMinutes, "SG", ost)
		totalMinutes += setPositionMinutes(sfList, rMap, sfMinutes, "SF", ost)
		totalMinutes += setPositionMinutes(pfList, rMap, pfMinutes, "PF", ost)
		totalMinutes += setPositionMinutes(cList, rMap, cMinutes, "C", ost)
	} else if ost == "Jumbo" {
		totalMinutes += setPositionMinutes(cList, rMap, cMinutes, "C", ost)
		totalMinutes += setPositionMinutes(pfList, rMap, pfMinutes, "PF", ost)
		totalMinutes += setPositionMinutes(sfList, rMap, sfMinutes, "SF", ost)
		totalMinutes += setPositionMinutes(sgList, rMap, sgMinutes, "SG", ost)
	} else if ost == "Small Ball" {
		totalMinutes += setPositionMinutes(sgList, rMap, sgMinutes, "SG", ost)
		totalMinutes += setPositionMinutes(pgList, rMap, pgMinutes, "PG", ost)
		totalMinutes += setPositionMinutes(sfList, rMap, sfMinutes, "SF", ost)
		totalMinutes += setPositionMinutes(pfList, rMap, pfMinutes, "PF", ost)
	} else if ost == "Microball" {
		totalMinutes += setPositionMinutes(pgList, rMap, pgMinutes, "PG", ost)
		totalMinutes += setPositionMinutes(sgList, rMap, sgMinutes, "SG", ost)
		totalMinutes += setPositionMinutes(sfList, rMap, sfMinutes, "SF", ost)
	}

	// For testing purposes
	teamMidRangeProportion := 0.0
	teamMidrangeLimit := 40.0
	teamInsideProportion := 0.0
	teamInsideLimit := 40.0
	teamThreePointProportion := 0.0
	teamThreePointLimit := 20.0

	sort.Slice(roster, func(i, j int) bool {
		return roster[i].Minutes > roster[j].Minutes && roster[i].Overall > roster[j].Overall
	})

	teamTotalSkill := 0
	for i := 0; i < len(roster); i++ {
		if roster[i].Minutes == 0 || roster[i].IsRedshirting || roster[i].IsInjured {
			continue
		}
		teamTotalSkill += roster[i].Shooting2 + roster[i].Shooting3 + roster[i].Finishing
	}

	// Loop for team shot proportions
	for i := 0; i < len(roster); i++ {
		if roster[i].Minutes == 0 || roster[i].IsRedshirting || roster[i].IsInjured {
			continue
		}
		totalSkill := roster[i].Shooting2 + roster[i].Shooting3 + roster[i].Finishing
		twoPointPercentage := float64(roster[i].Shooting2*100) / float64(totalSkill) * float64(roster[i].Minutes) / float64(roster[i].Stamina)
		threePointPercentage := float64(roster[i].Shooting3*100) / float64(totalSkill) * float64(roster[i].Minutes) / float64(roster[i].Stamina)
		insidePercentage := float64(roster[i].Finishing*100) / float64(totalSkill) * float64(roster[i].Minutes) / float64(roster[i].Stamina)
		teamInsideProportion += insidePercentage
		roster[i].SetInsideProportion(insidePercentage)
		teamMidRangeProportion += twoPointPercentage
		roster[i].SetMidShotProportion(twoPointPercentage)
		teamThreePointProportion += threePointPercentage
		roster[i].SetThreePointProportion(threePointPercentage)
	}

	insideProp := 0.0
	midProp := 0.0
	tpProp := 0.0

	// Motion
	if float64(teamThreePointProportion/teamMidRangeProportion) > 1.3 && float64(teamInsideProportion/teamMidRangeProportion) > 1.3 {
		off = "Motion"
		teamInsideLimit = 20
		teamMidrangeLimit = 10
		teamThreePointLimit = 70
		// Pick-And-Roll
	} else if float64(teamInsideProportion/teamMidRangeProportion) > 1.3 && float64(teamInsideProportion/teamThreePointProportion) > 1.3 {
		off = "Pick-and-Roll"
		teamInsideLimit = 40
		teamMidrangeLimit = 20
		teamThreePointLimit = 40
		// Post-Up
	} else if float64(teamInsideProportion/teamMidRangeProportion) > 1.5 && float64(teamInsideProportion/teamThreePointProportion) > 1.5 {
		off = "Post-Up"
		teamInsideLimit = 80
		teamMidrangeLimit = 15
		teamThreePointLimit = 5
		// Space-And-Post
	} else if float64(teamMidRangeProportion/teamInsideProportion) > 1.3 && float64(teamThreePointProportion/teamInsideProportion) > 1.3 {
		off = "Space-and-Post"
		teamInsideLimit = 20
		teamMidrangeLimit = 40
		teamThreePointLimit = 40
	}

	for i := 0; i < len(roster); i++ {
		if roster[i].Minutes == 0 || roster[i].IsRedshirting || roster[i].IsInjured {
			continue
		}
		normalizedInsideProportion := (roster[i].InsideProportion * float64(teamInsideLimit)) / teamInsideProportion
		insideProp += normalizedInsideProportion
		if insideProp > teamInsideLimit {
			diff := insideProp - teamInsideLimit
			insideProp -= diff
			normalizedInsideProportion -= diff
		}
		roster[i].SetInsideProportion(normalizedInsideProportion)

		normalizedMidrangeProportion := (roster[i].MidRangeProportion * float64(teamMidrangeLimit)) / teamMidRangeProportion
		midProp += normalizedMidrangeProportion
		if midProp > teamMidrangeLimit {
			diff := midProp - teamMidrangeLimit
			midProp -= diff
			normalizedMidrangeProportion -= diff
		}
		roster[i].SetMidShotProportion(normalizedMidrangeProportion)

		normalized3ptProportion := (roster[i].ThreePointProportion * float64(teamThreePointLimit)) / teamThreePointProportion
		tpProp += normalized3ptProportion
		if tpProp > teamThreePointLimit {
			diff := tpProp - teamThreePointLimit
			tpProp -= diff
			normalized3ptProportion -= diff
		}
		roster[i].SetThreePointProportion(normalized3ptProportion)
	}

	for _, r := range roster {
		repository.SaveCollegePlayerRecord(r, db)
	}

	gameplan.UpdateGameplan(pace, off, def, ost, "")

	db.Save(&gameplan)
}

func setPositionMinutes(list []structs.CollegePlayer, rMap map[string]*structs.CollegePlayer, limit int, pos, ost string) int {
	curr := 0
	for curr < limit {
		for _, c := range list {
			if curr >= limit {
				break
			}
			id := strconv.Itoa(int(c.ID))
			p := rMap[id]
			if p.Minutes == p.Stamina {
				continue
			}

			min := p.Minutes
			diff := p.Stamina - min

			if diff > 30 {
				diff = util.GenerateIntFromRange(25, 30)
				// If we have a negative number, reset to 0
			} else if diff < 0 {
				diff = 0
			}

			// If the difference is greater than the allowed limit - the current allocation, set difference to the number remaining
			if diff > limit-curr {
				diff = limit - curr
			}

			if p.P1Minutes == 0 {
				p.SetP1Minutes(diff, pos)
			} else if p.P2Minutes == 0 && p.PositionOne != pos {
				p.SetP2Minutes(diff, pos)
			} else if p.P3Minutes == 0 && p.PositionOne != pos && p.PositionTwo != pos {
				p.SetP3Minutes(diff, pos)
			}
			curr = addCurrentMinutes(curr, diff, limit)
		}
	}

	return curr
}

func setNBAPositionMinutes(list []structs.NBAPlayer, rMap map[string]*structs.NBAPlayer, limit int, pos, ost string) int {
	curr := 0
	for _, c := range list {
		if curr >= limit {
			break
		}
		id := strconv.Itoa(int(c.ID))
		p := rMap[id]
		if p.Minutes == p.Stamina {
			continue
		}

		min := p.Minutes
		diff := p.Stamina - min

		if diff > 30 {
			diff = 30
			// If we have a negative number, reset to 0
		} else if diff < 0 {
			diff = 0
		}

		if diff > limit-curr {
			diff = limit - curr
		}

		if p.P1Minutes == 0 {
			p.SetP1Minutes(diff, pos)
			curr = addCurrentMinutes(curr, diff, limit)
		} else if p.P2Minutes == 0 && p.PositionOne != pos {
			p.SetP2Minutes(diff, pos)
			curr = addCurrentMinutes(curr, diff, limit)
		} else if p.P3Minutes == 0 && p.PositionOne != pos && p.PositionTwo != pos {
			p.SetP3Minutes(diff, pos)
			curr = addCurrentMinutes(curr, diff, limit)
		}
	}
	return curr
}

func addCurrentMinutes(curr, diff, limit int) int {
	num := curr
	num += diff
	if num > limit {
		newDiff := num - limit
		num -= newDiff
	}
	return num
}
