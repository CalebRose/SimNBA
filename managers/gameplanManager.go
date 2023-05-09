package managers

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
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
	fmt.Printf("Saving Gameplan for Team " + teamId + "\n")
	db.Save(&gameplan)

	// Get Players
	updatedPlayers := updateGameplanDto.CollegePlayers

	for _, player := range updatedPlayers {
		id := strconv.Itoa(int(player.PlayerID))
		record := GetCollegePlayerByPlayerID(id)
		record.UpdatePlayer(player.BasePlayer)
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
	fmt.Printf("Saving Gameplan for Team " + teamId + "\n")
	db.Save(&gameplan)

	// Get Players
	updatedPlayers := updateGameplanDto.NBAPlayers

	for _, player := range updatedPlayers {
		id := strconv.Itoa(int(player.ID))
		record := GetNBAPlayerRecord(id)
		record.UpdatePlayer(player.BasePlayer)
		db.Save(&record)
	}
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

	nextMatchType := ""
	if !ts.GamesARan {
		nextMatchType = "A"
	} else {
		nextMatchType = "B"
	}

	matches := GetTeamMatchesByWeekId(strconv.Itoa(int(ts.CollegeWeekID)), strconv.Itoa(int(ts.SeasonID)), nextMatchType, teamID)
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

	opposingRoster := GetCollegePlayersByTeamId(opposingTeamID)

	return opposingRoster
}

func GetOpposingNBATeamRoster(teamID string) []structs.NBAPlayer {
	ts := GetTimestamp()

	nextMatchType := ""
	if !ts.GamesARan {
		nextMatchType = "A"
	} else if !ts.GamesBRan {
		nextMatchType = "B"
	} else {
		nextMatchType = "C"
	}

	matches := GetNBATeamMatchesByWeekId(strconv.Itoa(int(ts.CollegeWeekID)), strconv.Itoa(int(ts.SeasonID)), nextMatchType, teamID)
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

func SetAIGameplans() bool {
	db := dbprovider.GetInstance().GetDB()

	teams := GetAllActiveCollegeTeams()

	for _, team := range teams {
		if !team.IsActive {
			continue
		}

		if len(team.Coach) > 0 && team.Coach != "AI" {
			continue
		}

		pgCount := 0
		sgCount := 0
		sfCount := 0
		pfCount := 0
		cCount := 0
		pgMin := 0
		sgMin := 0
		sfMin := 0
		pfMin := 0
		cMin := 0

		pgList := []structs.CollegePlayer{}
		sgList := []structs.CollegePlayer{}
		sfList := []structs.CollegePlayer{}
		pfList := []structs.CollegePlayer{}
		cList := []structs.CollegePlayer{}

		gameplan := GetGameplansByTeam(strconv.Itoa(int(team.ID)))
		off := "Balanced"
		def := "Man-to-Man"
		ost := ""
		pace := "Balanced"

		roster := GetCollegePlayersByTeamId(strconv.Itoa(int(team.ID)))
		rMap := make(map[string]*structs.CollegePlayer)
		for i := 0; i < len(roster); i++ {
			id := strconv.Itoa(int(roster[i].ID))
			rMap[id] = &roster[i]
		}

		for _, c := range roster {
			if c.IsRedshirting {
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
			pgMin = 40
			sgMin = 40
			pfMin = 40
			sfMin = 40
			cMin = 40
		} else if ost == "Small Ball" {
			pgMin = 40
			sgMin = 80
			pfMin = 40
			sfMin = 40
			cMin = 0
		} else if ost == "Microball" {
			pgMin = 80
			sgMin = 80
			pfMin = 0
			sfMin = 40
			cMin = 0
		} else if ost == "Jumbo" {
			pgMin = 0
			sgMin = 40
			pfMin = 80
			sfMin = 40
			cMin = 40
		}
		sort.Slice(pgList, func(i, j int) bool {
			return pgList[i].Overall > pgList[j].Overall
		})

		sort.Slice(sgList, func(i, j int) bool {
			return sgList[i].Overall > sgList[j].Overall
		})

		sort.Slice(sfList, func(i, j int) bool {
			return sfList[i].Overall < sfList[j].Overall
		})

		sort.Slice(pfList, func(i, j int) bool {
			return pfList[i].Overall > pfList[j].Overall
		})

		sort.Slice(cList, func(i, j int) bool {
			return cList[i].Overall > cList[j].Overall
		})
		totalMinutes := 0
		if ost == "Traditional" {
			totalMinutes += setPositionMinutes(pgList, rMap, pgMin, "PG", ost)
			totalMinutes += setPositionMinutes(sgList, rMap, sgMin, "SG", ost)
			totalMinutes += setPositionMinutes(sfList, rMap, sfMin, "SF", ost)
			totalMinutes += setPositionMinutes(pfList, rMap, pfMin, "PF", ost)
			totalMinutes += setPositionMinutes(cList, rMap, cMin, "C", ost)
		} else if ost == "Jumbo" {
			totalMinutes += setPositionMinutes(cList, rMap, cMin, "C", ost)
			totalMinutes += setPositionMinutes(pfList, rMap, pfMin, "PF", ost)
			totalMinutes += setPositionMinutes(sfList, rMap, sfMin, "SF", ost)
			totalMinutes += setPositionMinutes(sgList, rMap, sgMin, "SG", ost)
		} else if ost == "Small Ball" {
			totalMinutes += setPositionMinutes(sgList, rMap, sgMin, "SG", ost)
			totalMinutes += setPositionMinutes(pgList, rMap, pgMin, "PG", ost)
			totalMinutes += setPositionMinutes(sfList, rMap, sfMin, "SF", ost)
			totalMinutes += setPositionMinutes(pfList, rMap, pfMin, "PF", ost)
		} else if ost == "Microball" {
			totalMinutes += setPositionMinutes(pgList, rMap, pgMin, "PG", ost)
			totalMinutes += setPositionMinutes(sgList, rMap, sgMin, "SG", ost)
			totalMinutes += setPositionMinutes(sfList, rMap, sfMin, "SF", ost)
		}

		// For testing purposes
		midProp := 0
		midLimit := 40
		insideProp := 0
		insideLimit := 40
		tpProp := 0
		tpLimit := 20
		sort.Slice(roster, func(i, j int) bool {
			return roster[i].Minutes > roster[j].Minutes && roster[i].Finishing > roster[j].Finishing
		})

		for i := 0; i < len(roster); i++ {
			if insideProp == insideLimit {
				break
			}
			if roster[i].Minutes == 0 || roster[i].IsRedshirting {
				continue
			}

			fn := 2
			if roster[i].Finishing > 15 {
				fn = 3
			}
			inside := (roster[i].Finishing / 4) * fn
			if insideProp < insideLimit {

				insideProp += inside
				if insideProp > insideLimit {
					diff := insideProp - insideLimit
					insideProp -= diff
					inside -= diff
				}
				roster[i].SetInsideProportion(inside)
			}
		}

		sort.Slice(roster, func(i, j int) bool {
			return roster[i].Minutes > roster[j].Minutes && roster[i].Shooting2 > roster[j].Shooting2
		})

		for i := 0; i < len(roster); i++ {
			if roster[i].Minutes == 0 || roster[i].IsRedshirting {
				continue
			}

			if insideProp < insideLimit {
				insdif := 0
				insdif = insideLimit - insideProp
				ins := int(roster[i].InsideProportion) + insdif
				insideProp += insdif
				roster[i].SetInsideProportion(ins)
			}

			if midProp == midLimit {
				continue
			}

			// Make shot proportion by formula
			s2 := 2
			if roster[i].Shooting2 > 15 {
				s2 = 3
			}
			mid := (roster[i].Shooting2 / 4) * s2
			if midProp < midLimit {
				midProp += mid
				if midProp > midLimit {
					diff := midProp - midLimit
					midProp -= diff
					mid -= diff
				}
				roster[i].SetMidShotProportion(mid)
			}
		}

		sort.Slice(roster, func(i, j int) bool {
			return roster[i].Minutes > roster[j].Minutes && roster[i].Shooting3 > roster[j].Shooting3
		})

		for i := 0; i < len(roster); i++ {
			if roster[i].Minutes == 0 || roster[i].IsRedshirting {
				continue
			}
			if midProp < midLimit {
				middif := 0
				middif = midLimit - midProp
				midProp += middif
				mid := int(roster[i].MidRangeProportion) + middif
				roster[i].SetMidShotProportion(mid)
			}

			if tpProp == tpLimit {
				continue
			}

			s3 := 1
			if team.FirstSeason == "2023" {
				s3 = util.GenerateIntFromRange(2, 3)
			}
			if roster[i].Shooting3 > 15 {
				s3 = 4
			} else if roster[i].Shooting3 > 10 {
				s3 = 3
			}
			tp := (roster[i].Shooting3 / 4) * s3
			tpProp += tp
			if tpProp > tpLimit {
				diff := tpProp - tpLimit
				tpProp -= diff
				tp -= diff
			}
			if tpProp >= 17 && tpProp < tpLimit {
				tpdif := tpLimit - tpProp
				tp += tpdif
				tpProp += tpdif
			}
			roster[i].SetThreePointProportion(tp)
		}

		if tpProp < tpLimit {
			for i := 0; i < len(roster); i++ {
				if roster[i].Minutes == 0 || roster[i].IsRedshirting {
					continue
				}
				if tpProp == tpLimit {
					break
				}

				s3 := 1
				if team.FirstSeason == "2023" {
					s3 = util.GenerateIntFromRange(2, 3)
				}
				if roster[i].Shooting3 > 15 {
					s3 = 4
				} else if roster[i].Shooting3 > 10 {
					s3 = 3
				}
				tp := (roster[i].Shooting3 / 4) * s3
				tpProp += tp
				if tpProp > tpLimit {
					diff := tpProp - tpLimit
					tpProp -= diff
					tp -= diff
				}
				roster[i].SetThreePointProportion(tp)
			}
		}

		for _, r := range roster {
			db.Save(&r)
		}

		gameplan.UpdateGameplan(pace, off, def, ost, "")

		db.Save(&gameplan)
	}

	return true
}

func setPositionMinutes(list []structs.CollegePlayer, rMap map[string]*structs.CollegePlayer, limit int, pos, ost string) int {
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
		diff = num - limit
		num -= diff
	}
	return num
}
