package managers

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/repository"
	"github.com/CalebRose/SimNBA/structs"
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
		db.Save(&record)
	}
}

func UpdateNBAGameplan(updateGameplanDto structs.UpdateGameplanDto) {
	// Will need to redesign this function to account for new updates
}

func GetAllCollegeGameplans() []structs.Gameplan {
	db := dbprovider.GetInstance().GetDB()

	var gameplans []structs.Gameplan

	db.Find(&gameplans)

	return gameplans
}

func GetAllNBAGameplans() []structs.NBAGameplan {
	db := dbprovider.GetInstance().GetDB()

	var gameplans []structs.NBAGameplan

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
	// db := dbprovider.GetInstance().GetDB()
	// teams := GetAllActiveCollegeTeams()
	// // collegeTeamMap := MakeCollegeTeamMap(teams)
	// collegePlayers := GetAllCollegePlayers()
	// collegePlayerMapByTeamID := MakeCollegePlayerMapByTeamID(collegePlayers, true)
	// gameplans := GetAllCollegeGameplans()
	// gameplanMap := MakeCollegeGameplanMap(gameplans)
	// // NBA
	// nbaTeams := GetAllActiveNBATeams()
	// nbaTeamMap := MakeNBATeamMap(nbaTeams)
	// nbaPlayers := GetAllNBAPlayers()
	// nbaPlayerMapByTeamID := MakeNBAPlayerMapByTeamID(nbaPlayers, false)
	// nbaGameplans := GetAllNBAGameplans()
	// nbaGameplanMap := MakeNBAGameplanMap(nbaGameplans)
	// nbaTeamIDsToCheck := []uint{}
	// nbaTeamIDStrings := []string{}
}

func GetAllCollegeLineups() []structs.CollegeLineup {
	return repository.FindCollegeLineupRecords(repository.GameplanQuery{})
}

func GetAllNBALineups() []structs.NBALineup {
	return repository.FindNBALineupRecords(repository.GameplanQuery{})
}

// calcShotProportions derives per-player shot type splits from shooting attributes.
// Returns (inside, mid, threePoint) proportions that sum to 100.
func calcShotProportions(inside, mid, three uint8) (uint8, uint8, uint8) {
	total := uint16(inside) + uint16(mid) + uint16(three)
	if total == 0 {
		return 34, 33, 33
	}
	inProp := uint8(uint16(inside) * 100 / total)
	midProp := uint8(uint16(mid) * 100 / total)
	tpProp := uint8(100) - inProp - midProp
	return inProp, midProp, tpProp
}

// fillCollegeLineupSlots assigns players from a sorted bucket to lineup slots using a
// round-robin draft style so the best players rotate across slots. College games are
// 40 minutes per position slot.
func fillCollegeLineupSlots(slots []*structs.CollegeLineup, players []structs.CollegePlayer) {
	n := len(slots)
	if n == 0 {
		return
	}
	for i, slot := range slots {
		slot.FirstStringID = 0
		slot.FSMinutes = 0
		slot.FSInsideProportion = 0
		slot.FSMidProportion = 0
		slot.FSThreeProportion = 0
		slot.SecondStringID = 0
		slot.SSMinutes = 0
		slot.SSInsideProportion = 0
		slot.SSMidProportion = 0
		slot.SSThreeProportion = 0
		slot.ThirdStringID = 0
		slot.TSMinutes = 0
		slot.TSInsideProportion = 0
		slot.TSMidProportion = 0
		slot.TSThreeProportion = 0

		fsIdx := i
		ssIdx := i + n
		tsIdx := i + n*2

		hasFS := fsIdx < len(players)
		hasSS := ssIdx < len(players)
		hasTS := tsIdx < len(players)

		switch {
		case hasFS && hasSS && hasTS:
			slot.FSMinutes = 24
			slot.SSMinutes = 13
			slot.TSMinutes = 3
		case hasFS && hasSS:
			slot.FSMinutes = 27
			slot.SSMinutes = 13
		case hasFS:
			slot.FSMinutes = 40
		}

		if hasFS {
			p := players[fsIdx]
			in, mid, tp := calcShotProportions(p.InsideShooting, p.MidRangeShooting, p.ThreePointShooting)
			slot.FirstStringID = p.ID
			slot.FSInsideProportion = in
			slot.FSMidProportion = mid
			slot.FSThreeProportion = tp
		}
		if hasSS {
			p := players[ssIdx]
			in, mid, tp := calcShotProportions(p.InsideShooting, p.MidRangeShooting, p.ThreePointShooting)
			slot.SecondStringID = p.ID
			slot.SSInsideProportion = in
			slot.SSMidProportion = mid
			slot.SSThreeProportion = tp
		}
		if hasTS {
			p := players[tsIdx]
			in, mid, tp := calcShotProportions(p.InsideShooting, p.MidRangeShooting, p.ThreePointShooting)
			slot.ThirdStringID = p.ID
			slot.TSInsideProportion = in
			slot.TSMidProportion = mid
			slot.TSThreeProportion = tp
		}
	}
}

// fillNBALineupSlots assigns players from a sorted bucket to lineup slots using a
// round-robin draft style. NBA games are 48 minutes per position slot.
func fillNBALineupSlots(slots []*structs.NBALineup, players []structs.NBAPlayer) {
	n := len(slots)
	if n == 0 {
		return
	}
	for i, slot := range slots {
		slot.FirstStringID = 0
		slot.FSMinutes = 0
		slot.FSInsideProportion = 0
		slot.FSMidProportion = 0
		slot.FSThreeProportion = 0
		slot.SecondStringID = 0
		slot.SSMinutes = 0
		slot.SSInsideProportion = 0
		slot.SSMidProportion = 0
		slot.SSThreeProportion = 0
		slot.ThirdStringID = 0
		slot.TSMinutes = 0
		slot.TSInsideProportion = 0
		slot.TSMidProportion = 0
		slot.TSThreeProportion = 0

		fsIdx := i
		ssIdx := i + n
		tsIdx := i + n*2

		hasFS := fsIdx < len(players)
		hasSS := ssIdx < len(players)
		hasTS := tsIdx < len(players)

		switch {
		case hasFS && hasSS && hasTS:
			slot.FSMinutes = 30
			slot.SSMinutes = 13
			slot.TSMinutes = 5
		case hasFS && hasSS:
			slot.FSMinutes = 32
			slot.SSMinutes = 16
		case hasFS:
			slot.FSMinutes = 48
		}

		if hasFS {
			p := players[fsIdx]
			in, mid, tp := calcShotProportions(p.InsideShooting, p.MidRangeShooting, p.ThreePointShooting)
			slot.FirstStringID = p.ID
			slot.FSInsideProportion = in
			slot.FSMidProportion = mid
			slot.FSThreeProportion = tp
		}
		if hasSS {
			p := players[ssIdx]
			in, mid, tp := calcShotProportions(p.InsideShooting, p.MidRangeShooting, p.ThreePointShooting)
			slot.SecondStringID = p.ID
			slot.SSInsideProportion = in
			slot.SSMidProportion = mid
			slot.SSThreeProportion = tp
		}
		if hasTS {
			p := players[tsIdx]
			in, mid, tp := calcShotProportions(p.InsideShooting, p.MidRangeShooting, p.ThreePointShooting)
			slot.ThirdStringID = p.ID
			slot.TSInsideProportion = in
			slot.TSMidProportion = mid
			slot.TSThreeProportion = tp
		}
	}
}

func SetAIGameplans() bool {
	db := dbprovider.GetInstance().GetDB()

	teams := GetAllActiveCollegeTeams()
	collegePlayers := GetAllCollegePlayers()
	collegePlayerMapByTeamID := MakeCollegePlayerMapByTeamID(collegePlayers, true)
	collegeLineups := GetAllCollegeLineups()
	collegeLineupMap := MakeCollegeLineupMapByTeamID(collegeLineups)

	for _, team := range teams {
		// if team.IsUserCoached {
		// 	continue
		// }
		SetCollegeMinutesAndShotProportions(db, team.ID, collegeLineupMap, collegePlayerMapByTeamID)
	}

	islTeams := GetAllActiveNBATeams()
	nbaPlayers := GetAllNBAPlayers()
	nbaPlayerMapByTeamID := MakeNBAPlayerMapByTeamID(nbaPlayers, false)
	nbaLineups := GetAllNBALineups()
	nbaLineupMap := MakeNBALineupMapByTeamID(nbaLineups)

	for _, team := range islTeams {
		// if !team.IsActive {
		// 	continue
		// }

		// if len(team.NBAOwnerName) > 0 && team.NBAOwnerName != "AI" {
		// 	continue
		// }

		SetNBAMinutesAndShotProportions(db, team.ID, nbaLineupMap, nbaPlayerMapByTeamID)
	}

	return true
}

func SetCollegeMinutesAndShotProportions(db *gorm.DB, teamID uint, lineupMap map[uint][]structs.CollegeLineup, collegePlayerMapByTeamID map[uint][]structs.CollegePlayer) {

	roster := collegePlayerMapByTeamID[teamID]
	lineups := lineupMap[teamID]
	if len(lineups) == 0 || len(roster) == 0 {
		return
	}

	gPlayers := []structs.CollegePlayer{}
	fPlayers := []structs.CollegePlayer{}
	cPlayers := []structs.CollegePlayer{}

	for _, p := range roster {
		if p.IsRedshirting || p.IsInjured {
			continue
		}
		switch p.Position {
		case "G":
			gPlayers = append(gPlayers, p)
		case "F":
			fPlayers = append(fPlayers, p)
		case "C":
			cPlayers = append(cPlayers, p)
		}
	}

	sort.Slice(gPlayers, func(i, j int) bool { return gPlayers[i].Overall > gPlayers[j].Overall })
	sort.Slice(fPlayers, func(i, j int) bool { return fPlayers[i].Overall > fPlayers[j].Overall })
	sort.Slice(cPlayers, func(i, j int) bool { return cPlayers[i].Overall > cPlayers[j].Overall })

	gSlots := []*structs.CollegeLineup{}
	fSlots := []*structs.CollegeLineup{}
	cSlots := []*structs.CollegeLineup{}

	for i := range lineups {
		switch lineups[i].Position {
		case "G":
			gSlots = append(gSlots, &lineups[i])
		case "F":
			fSlots = append(fSlots, &lineups[i])
		case "C":
			cSlots = append(cSlots, &lineups[i])
		}
	}

	fillCollegeLineupSlots(gSlots, gPlayers)
	fillCollegeLineupSlots(fSlots, fPlayers)
	fillCollegeLineupSlots(cSlots, cPlayers)

	for _, lineup := range lineups {
		repository.SaveCollegeLineupRecord(lineup, db)
	}
}

func SetNBAMinutesAndShotProportions(db *gorm.DB, teamID uint, lineupMap map[uint][]structs.NBALineup, nbaPlayerMapByTeamID map[uint][]structs.NBAPlayer) {
	roster := nbaPlayerMapByTeamID[teamID]
	lineups := lineupMap[teamID]
	if len(lineups) == 0 || len(roster) == 0 {
		return
	}

	gPlayers := []structs.NBAPlayer{}
	fPlayers := []structs.NBAPlayer{}
	cPlayers := []structs.NBAPlayer{}

	for _, p := range roster {
		if p.IsGLeague || p.IsInjured {
			continue
		}
		switch p.Position {
		case "G":
			gPlayers = append(gPlayers, p)
		case "F":
			fPlayers = append(fPlayers, p)
		case "C":
			cPlayers = append(cPlayers, p)
		}
	}

	sort.Slice(gPlayers, func(i, j int) bool { return gPlayers[i].Overall > gPlayers[j].Overall })
	sort.Slice(fPlayers, func(i, j int) bool { return fPlayers[i].Overall > fPlayers[j].Overall })
	sort.Slice(cPlayers, func(i, j int) bool { return cPlayers[i].Overall > cPlayers[j].Overall })

	gSlots := []*structs.NBALineup{}
	fSlots := []*structs.NBALineup{}
	cSlots := []*structs.NBALineup{}

	for i := range lineups {
		switch lineups[i].Position {
		case "G":
			gSlots = append(gSlots, &lineups[i])
		case "F":
			fSlots = append(fSlots, &lineups[i])
		case "C":
			cSlots = append(cSlots, &lineups[i])
		}
	}

	fillNBALineupSlots(gSlots, gPlayers)
	fillNBALineupSlots(fSlots, fPlayers)
	fillNBALineupSlots(cSlots, cPlayers)

	for _, lineup := range lineups {
		repository.SaveNBALineupRecord(lineup, db)
	}
}
