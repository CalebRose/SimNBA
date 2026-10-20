package managers

import (
	"sort"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/repository"
	"github.com/CalebRose/SimNBA/structs"
)

func GenerateTeamRatings() {
	db := dbprovider.GetInstance().GetDB()
	collegeTeams := GetAllActiveCollegeTeams()
	nbaTeams := GetAllActiveNBATeams()

	collegePlayers := GetAllCollegePlayers()
	nbaPlayers := GetAllNBAPlayers()

	collegePlayerMap := MakeCollegePlayerMapByTeamID(collegePlayers, true)
	nbaPlayerMap := MakeNBAPlayerMapByTeamID(nbaPlayers, true)

	var collegeTeamGrades []structs.TeamGrade
	var nbaTeamGrades []structs.TeamGrade
	for _, team := range collegeTeams {
		players := collegePlayerMap[team.ID]

		basePlayers := MakeBasePlayerList(players, []structs.NBAPlayer{})

		offenseGrade := calculateOffenseGrade(basePlayers)
		defenseGrade := calculateDefenseGrade(basePlayers)
		overallGrade := calculateOverallGrade(offenseGrade, defenseGrade)

		tg := structs.TeamGrade{}
		tg.SetOffenseGradeNumber(offenseGrade)
		tg.SetDefenseGradeNumber(defenseGrade)
		tg.SetOverallGradeNumber(overallGrade)

		collegeTeamGrades = append(collegeTeamGrades, tg)
	}
	cbbOffenseMean, cbbOffenseStdDev := calculateMeanAndStdDev(collegeTeamGrades, "offense")
	cbbDefenseMean, cbbDefenseStdDev := calculateMeanAndStdDev(collegeTeamGrades, "defense")
	cbbOverallMean, cbbOverallStdDev := calculateMeanAndStdDev(collegeTeamGrades, "overall")

	for i, team := range collegeTeams {
		tg := collegeTeamGrades[i]
		offenseLetter := assignLetterGrade(tg.OffenseGradeNumber, cbbOffenseMean, cbbOffenseStdDev)
		defenseLetter := assignLetterGrade(tg.DefenseGradeNumber, cbbDefenseMean, cbbDefenseStdDev)
		overallLetter := assignLetterGrade(tg.OverallGradeNumber, cbbOverallMean, cbbOverallStdDev)
		teamToUpdate := team
		teamToUpdate.AssignRatings(offenseLetter, defenseLetter, overallLetter)
		repository.SaveCollegeTeamRecord(teamToUpdate, db)
	}

	for _, team := range nbaTeams {
		players := nbaPlayerMap[team.ID]

		basePlayers := MakeBasePlayerList([]structs.CollegePlayer{}, players)

		offenseGrade := calculateOffenseGrade(basePlayers)
		defenseGrade := calculateDefenseGrade(basePlayers)
		overallGrade := calculateOverallGrade(offenseGrade, defenseGrade)

		tg := structs.TeamGrade{}
		tg.SetOffenseGradeNumber(offenseGrade)
		tg.SetDefenseGradeNumber(defenseGrade)
		tg.SetOverallGradeNumber(overallGrade)

		nbaTeamGrades = append(nbaTeamGrades, tg)
	}
	nbaOffenseMean, nbaOffenseStdDev := calculateMeanAndStdDev(nbaTeamGrades, "offense")
	nbaDefenseMean, nbaDefenseStdDev := calculateMeanAndStdDev(nbaTeamGrades, "defense")
	nbaOverallMean, nbaOverallStdDev := calculateMeanAndStdDev(nbaTeamGrades, "overall")

	for i, team := range nbaTeams {
		tg := nbaTeamGrades[i]
		offenseLetter := assignLetterGrade(tg.OffenseGradeNumber, nbaOffenseMean, nbaOffenseStdDev)
		defenseLetter := assignLetterGrade(tg.DefenseGradeNumber, nbaDefenseMean, nbaDefenseStdDev)
		overallLetter := assignLetterGrade(tg.OverallGradeNumber, nbaOverallMean, nbaOverallStdDev)
		teamToUpdate := team
		teamToUpdate.AssignRatings(offenseLetter, defenseLetter, overallLetter)
		repository.SaveNBATeamRecord(teamToUpdate, db)
	}
}

// calculateOffenseGrade calculates the offense grade based on top 3 centers and top 6 forwards
func calculateOffenseGrade(players []structs.BasePlayer) float64 {
	var centers []structs.BasePlayer
	var forwards []structs.BasePlayer
	var guards []structs.BasePlayer

	// Separate players by position
	for _, player := range players {
		switch player.Position {
		case "C":
			centers = append(centers, player)
		case "F":
			forwards = append(forwards, player)
		case "G":
			guards = append(guards, player)
		}
	}

	// Sort by overall rating (descending)
	sort.Slice(centers, func(i, j int) bool {
		return centers[i].Overall > centers[j].Overall
	})

	sort.Slice(forwards, func(i, j int) bool {
		return forwards[i].Overall > forwards[j].Overall
	})

	sort.Slice(guards, func(i, j int) bool {
		return guards[i].Overall > guards[j].Overall
	})

	// Get top 3 centers, top 6 forwards, and top 5 guards
	topCenters := getTopPlayers(centers, 3)
	topForwards := getTopPlayers(forwards, 6)
	topGuards := getTopPlayers(guards, 5)

	// Calculate average offensive rating
	var totalOffenseRating float64
	var playerCount int

	for _, center := range topCenters {
		offenseRating := calculatePlayerOffenseRating(center)
		totalOffenseRating += offenseRating
		playerCount++
	}

	for _, forward := range topForwards {
		offenseRating := calculatePlayerOffenseRating(forward)
		totalOffenseRating += offenseRating
		playerCount++
	}
	for _, guard := range topGuards {
		offenseRating := calculatePlayerOffenseRating(guard)
		totalOffenseRating += offenseRating
		playerCount++
	}

	if playerCount == 0 {
		return 0
	}

	return totalOffenseRating / float64(playerCount)
}

// calculateDefenseGrade calculates the defense grade based on top 4 defenders
func calculateDefenseGrade(players []structs.BasePlayer) float64 {
	var centers []structs.BasePlayer
	var forwards []structs.BasePlayer
	var guards []structs.BasePlayer

	// Separate players by position
	for _, player := range players {
		switch player.Position {
		case "C":
			centers = append(centers, player)
		case "F":
			forwards = append(forwards, player)
		case "G":
			guards = append(guards, player)
		}
	}

	// Sort by overall rating (descending)
	sort.Slice(centers, func(i, j int) bool {
		return centers[i].Overall > centers[j].Overall
	})

	sort.Slice(forwards, func(i, j int) bool {
		return forwards[i].Overall > forwards[j].Overall
	})

	sort.Slice(guards, func(i, j int) bool {
		return guards[i].Overall > guards[j].Overall
	})

	// Get top 3 centers, top 6 forwards, and top 5 guards
	topCenters := getTopPlayers(centers, 3)
	topForwards := getTopPlayers(forwards, 6)
	topGuards := getTopPlayers(guards, 5)

	// Calculate average defensive rating
	var totalDefenseRating float64
	playerCount := 0
	for _, center := range topCenters {
		defenseRating := calculatePlayerDefenseRating(center)
		totalDefenseRating += defenseRating
		playerCount++
	}

	for _, forward := range topForwards {
		defenseRating := calculatePlayerDefenseRating(forward)
		totalDefenseRating += defenseRating
		playerCount++
	}

	for _, guard := range topGuards {
		defenseRating := calculatePlayerDefenseRating(guard)
		totalDefenseRating += defenseRating
		playerCount++
	}

	if playerCount == 0 {
		return 0
	}

	return totalDefenseRating / float64(playerCount)
}

// calculateOverallGrade calculates overall grade with 50% offense and 50% defense
func calculateOverallGrade(offenseGrade, defenseGrade float64) float64 {
	return (offenseGrade * 0.5) + (defenseGrade * 0.5)
}

// calculatePlayerOffenseRating calculates a player's offensive contribution
func calculatePlayerOffenseRating(player structs.BasePlayer) float64 {
	// Weight offensive attributes: shooting, passing, puck handling, agility
	offenseRating := float64(player.InsideShooting)*0.20 +
		float64(player.MidRangeShooting)*0.20 +
		float64(player.ThreePointShooting)*0.20 +
		float64(player.FreeThrow)*0.5 +
		float64(player.Ballwork)*0.15 +
		float64(player.Agility)*0.10

	return offenseRating
}

// calculatePlayerDefenseRating calculates a player's defensive contribution
func calculatePlayerDefenseRating(player structs.BasePlayer) float64 {
	// Weight defensive attributes: checking, shot blocking, strength, agility
	return float64(player.InteriorDefense)*0.3 +
		float64(player.PerimeterDefense)*0.3 +
		float64(player.Rebounding)*0.20 +
		float64(player.Blocking)*0.10 +
		float64(player.Stealing)*0.10
}

// getTopPlayers returns the top N players from a sorted list
func getTopPlayers(players []structs.BasePlayer, count int) []structs.BasePlayer {
	if len(players) <= count {
		return players
	}
	return players[:count]
}

// calculateMeanAndStdDev calculates mean and standard deviation for a specific grade type
func calculateMeanAndStdDev(teamGrades []structs.TeamGrade, gradeType string) (float64, float64) {
	if len(teamGrades) == 0 {
		return 0, 0
	}

	var sum float64
	for _, tg := range teamGrades {
		switch gradeType {
		case "offense":
			sum += tg.OffenseGradeNumber
		case "defense":
			sum += tg.DefenseGradeNumber
		case "overall":
			sum += tg.OverallGradeNumber
		}
	}

	mean := sum / float64(len(teamGrades))

	// Calculate variance
	var variance float64
	for _, tg := range teamGrades {
		var value float64
		switch gradeType {
		case "offense":
			value = tg.OffenseGradeNumber
		case "defense":
			value = tg.DefenseGradeNumber
		case "overall":
			value = tg.OverallGradeNumber
		}
		diff := value - mean
		variance += diff * diff
	}

	variance /= float64(len(teamGrades))
	stdDev := 0.0
	if variance > 0 {
		// Simple square root approximation using Newton's method
		stdDev = variance / 2
		for i := 0; i < 10; i++ {
			stdDev = (stdDev + variance/stdDev) / 2
		}
	}

	return mean, stdDev
}

// assignLetterGrade assigns letter grade based on standard deviations from mean
func assignLetterGrade(value, mean, stdDev float64) string {
	if stdDev == 0 {
		return "C"
	}

	z := (value - mean) / stdDev

	switch {
	case z >= 2.0:
		return "A+"
	case z >= 1.75:
		return "A"
	case z >= 1.5:
		return "A-"
	case z >= 1.25:
		return "B+"
	case z >= 1.0:
		return "B"
	case z >= 0.75:
		return "B-"
	case z >= 0.5:
		return "C+"
	case z >= -0.5:
		return "C"
	case z >= -0.75:
		return "C-"
	case z >= -1.0:
		return "D+"
	case z >= -1.5:
		return "D"
	case z >= -2.0:
		return "D-"
	default:
		return "F"
	}
}
