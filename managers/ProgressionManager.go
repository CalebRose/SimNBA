package managers

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
)

func ProgressionMain() {
	db := dbprovider.GetInstance().GetDB()
	fmt.Println(time.Now().UnixNano())
	rand.Seed(time.Now().UnixNano())

	collegeTeams := GetAllActiveCollegeTeams()

	for _, team := range collegeTeams {
		// var graduatingPlayers []structs.NBADraftee
		teamID := strconv.Itoa(int(team.ID))
		// roster := GetAllCollegePlayersWithStatsByTeamID(teamID, SeasonID)
		roster := GetCollegePlayersByTeamId(teamID)
		// croots := GetSignedRecruitsByTeamProfileID(teamID)

		for _, player := range roster {
			if player.HasProgressed {
				// player.FixAge()
				// err := db.Save(&player).Error
				// if err != nil {
				// 	log.Panicln("Could not save player record")
				// }
				continue
			}
			player = ProgressCollegePlayer(player)
			if player.IsRedshirting {
				player.SetRedshirtStatus()
			}

			player.SetExpectations(util.GetPlaytimeExpectations(player.Stars, player.Year))

			if (player.IsRedshirt && player.Year > 5) ||
				(!player.IsRedshirt && player.Year > 4) {
				player.GraduatePlayer()
				if player.Overall > 69 {
					draftee := structs.NBADraftee{}
					draftee.Map(player)
					draftee.AssignPrimeAge(util.GenerateIntFromRange(25, 30))

					err := db.Save(&draftee).Error
					if err != nil {
						log.Panicln("Could not save historic player record!")
					}
				}

				hcp := (structs.HistoricCollegePlayer)(player)

				err := db.Save(&hcp).Error
				if err != nil {
					log.Panicln("Could not save historic player record!")
				}
				// graduatingPlayers = append(graduatingPlayers, draftee)
				// CollegePlayer record will be deleted, but record will be mapped to a GraduatedCollegePlayer struct, and then saved in that table, along side with NFL Draftees table
				// GraduatedCollegePlayer will be a copy of the collegeplayers table, but only for historical players

				err = db.Delete(&player).Error
				if err != nil {
					log.Panicln("Could not delete old college player record.")
				}
			} else {
				err := db.Save(&player).Error
				if err != nil {
					log.Panicln("Could not save player record")
				}
			}

		}

		// for _, croot := range croots {
		// 	// Convert to College Player Record
		// 	cp := structs.CollegePlayer{}
		// 	cp.MapFromRecruit(croot, team)

		// 	// Save College Player Record
		// 	err := db.Save(&cp).Error
		// 	if err != nil {
		// 		log.Panicln("Could not save new college player record")
		// 	}

		// 	// Delete Recruit Record
		// }

		// Graduating players
		// err := db.CreateInBatches(&graduatingPlayers, len(graduatingPlayers)).Error
		// if err != nil {
		// 	log.Panicln("Could not save graduating players")
		// }
	}
}

func ProgressNBAPlayers() {
	db := dbprovider.GetInstance().GetDB()
	fmt.Println(time.Now().UnixNano())
	rand.Seed(time.Now().UnixNano())

	nbaTeams := GetAllActiveNBATeams()
	// Append empty team object to the end for Free Agents
	nbaTeams = append(nbaTeams, structs.NBATeam{})

	for _, team := range nbaTeams {
		teamID := strconv.Itoa(int(team.ID))

		roster := GetAllNBAPlayersByTeamID(teamID)

		for _, player := range roster {
			playerID := strconv.Itoa(int(player.ID))
			player = ProgressNBAPlayer(player)

			contract := GetNBAContractsByPlayerID(playerID)
			// Retiring Logic
			willPlayerRetire := util.WillPlayerRetire(player.Age, player.Overall)
			if willPlayerRetire {
				player.SetRetiringStatus()
				retiringPlayer := (structs.RetiredPlayer)(player)
				contract.RetireContract()
				db.Save(&contract)
				db.Create(&retiringPlayer)
				db.Delete(&player)
			} else {
				if player.IsMVP || player.IsDPOY || player.IsFirstTeamANBA {
					player.QualifyForSuperMax()
				} else if player.Overall > 100 {
					player.QualifiesForMax()
				}
				contract.ProgressContract()
				if contract.YearsRemaining == 0 && !contract.IsActive {
					player.BecomeFreeAgent()
				}

				// db.Save(&contract)
				db.Save(&player)
			}
		}

	}
}

func ProgressNBAPlayer(np structs.NBAPlayer) structs.NBAPlayer {
	stats := np.Stats
	totalMinutes := 0

	for _, stat := range stats {
		totalMinutes += stat.Minutes
	}

	var MinutesPerGame int = 0
	if len(stats) > 0 {
		MinutesPerGame = totalMinutes / len(stats)
	}
	age := np.Age + 1
	ageDifference := np.Age - int(np.PrimeAge)
	if ageDifference < 0 {
		ageDifference = 0
	}

	// Primary Progressions
	shooting2 := 0
	shooting3 := 0
	freeThrow := 0
	ballwork := 0
	rebounding := 0
	finishing := 0
	interiorDefense := 0
	perimeterDefense := 0

	attributeList := []string{}

	s2DiceRoll := util.GenerateIntFromRange(1, 20)
	s3DiceRoll := util.GenerateIntFromRange(1, 20)
	ftDiceRoll := util.GenerateIntFromRange(1, 20)
	fnDiceRoll := util.GenerateIntFromRange(1, 20)
	rbDiceRoll := util.GenerateIntFromRange(1, 20)
	bwDiceRoll := util.GenerateIntFromRange(1, 20)
	idDiceRoll := util.GenerateIntFromRange(1, 20)
	pdDiceRoll := util.GenerateIntFromRange(1, 20)

	potentialModifier := np.Potential / 20 // Guaranteed to be between 1-5

	if s2DiceRoll+potentialModifier > 15 || np.SpecShooting2 {
		attributeList = append(attributeList, "Shooting2")
	}

	if s3DiceRoll+potentialModifier > 15 || np.SpecShooting3 {
		attributeList = append(attributeList, "Shooting3")
	}
	if ftDiceRoll+potentialModifier > 15 || np.SpecFreeThrow {
		attributeList = append(attributeList, "FreeThrow")
	}
	if fnDiceRoll+potentialModifier > 15 || np.SpecFinishing {
		attributeList = append(attributeList, "Finishing")
	}
	if bwDiceRoll+potentialModifier > 15 || np.SpecBallwork {
		attributeList = append(attributeList, "Ballwork")
	}
	if rbDiceRoll+potentialModifier > 15 || np.SpecRebounding {
		attributeList = append(attributeList, "Rebounding")
	}
	if idDiceRoll+potentialModifier > 15 || np.SpecInteriorDefense {
		attributeList = append(attributeList, "InteriorDefense")
	}
	if pdDiceRoll+potentialModifier > 15 || np.SpecPerimeterDefense {
		attributeList = append(attributeList, "PerimeterDefense")
	}

	rand.Shuffle(len(attributeList), func(i, j int) {
		attributeList[i], attributeList[j] = attributeList[j], attributeList[i]
	})

	for _, attr := range attributeList {
		if attr == "Shooting2" {
			shooting2 = PlayerProgression(np.Potential, ageDifference, MinutesPerGame, np.PlaytimeExpectations, np.SpecShooting2, np.IsGLeague)
		} else if attr == "Shooting3" {
			shooting3 = PlayerProgression(np.Potential, ageDifference, MinutesPerGame, np.PlaytimeExpectations, np.SpecShooting3, np.IsGLeague)
		} else if attr == "FreeThrow" {
			freeThrow = PlayerProgression(np.Potential, ageDifference, MinutesPerGame, np.PlaytimeExpectations, np.SpecFreeThrow, np.IsGLeague)
		} else if attr == "Finishing" {
			finishing = PlayerProgression(np.Potential, ageDifference, MinutesPerGame, np.PlaytimeExpectations, np.SpecFinishing, np.IsGLeague)
		} else if attr == "Ballwork" {
			ballwork = PlayerProgression(np.Potential, ageDifference, MinutesPerGame, np.PlaytimeExpectations, np.SpecBallwork, np.IsGLeague)
		} else if attr == "Rebounding" {
			rebounding = PlayerProgression(np.Potential, ageDifference, MinutesPerGame, np.PlaytimeExpectations, np.SpecRebounding, np.IsGLeague)
		} else if attr == "InteriorDefense" {
			interiorDefense = PlayerProgression(np.Potential, ageDifference, MinutesPerGame, np.PlaytimeExpectations, np.SpecInteriorDefense, np.IsGLeague)
		} else if attr == "PerimeterDefense" {
			perimeterDefense = PlayerProgression(np.Potential, ageDifference, MinutesPerGame, np.PlaytimeExpectations, np.SpecPerimeterDefense, np.IsGLeague)
		}
	}

	stamina := ProgressStamina(np.Stamina, ageDifference)

	progressions := structs.NBAPlayerProgressions{
		Shooting2:        shooting2,
		Shooting3:        shooting3,
		Ballwork:         ballwork,
		Finishing:        finishing,
		Rebounding:       rebounding,
		InteriorDefense:  interiorDefense,
		PerimeterDefense: perimeterDefense,
		FreeThrow:        freeThrow,
		Age:              age,
		Stamina:          stamina,
	}

	np.Progress(progressions)

	return np
}

func ProgressCollegePlayer(cp structs.CollegePlayer) structs.CollegePlayer {
	stats := cp.Stats
	totalMinutes := 0

	for _, stat := range stats {
		totalMinutes += stat.Minutes
	}

	var MinutesPerGame int = 0
	if len(stats) > 0 {
		MinutesPerGame = totalMinutes / len(stats)
	}

	// Primary Progressions
	shooting2 := 0
	shooting3 := 0
	freeThrow := 0
	ballwork := 0
	rebounding := 0
	finishing := 0
	interiorDefense := 0
	perimeterDefense := 0

	attributeList := []string{}

	s2DiceRoll := util.GenerateIntFromRange(1, 20)
	s3DiceRoll := util.GenerateIntFromRange(1, 20)
	ftDiceRoll := util.GenerateIntFromRange(1, 20)
	fnDiceRoll := util.GenerateIntFromRange(1, 20)
	rbDiceRoll := util.GenerateIntFromRange(1, 20)
	bwDiceRoll := util.GenerateIntFromRange(1, 20)
	idDiceRoll := util.GenerateIntFromRange(1, 20)
	pdDiceRoll := util.GenerateIntFromRange(1, 20)

	potentialModifier := cp.Potential / 20 // Guaranteed to be between 1-5

	if s2DiceRoll+potentialModifier > 15 || cp.SpecShooting2 {
		attributeList = append(attributeList, "Shooting2")
	}

	if s3DiceRoll+potentialModifier > 15 || cp.SpecShooting3 {
		attributeList = append(attributeList, "Shooting3")
	}
	if ftDiceRoll+potentialModifier > 15 || cp.SpecFreeThrow {
		attributeList = append(attributeList, "FreeThrow")
	}
	if fnDiceRoll+potentialModifier > 15 || cp.SpecFinishing {
		attributeList = append(attributeList, "Finishing")
	}
	if bwDiceRoll+potentialModifier > 15 || cp.SpecBallwork {
		attributeList = append(attributeList, "Ballwork")
	}
	if rbDiceRoll+potentialModifier > 15 || cp.SpecRebounding {
		attributeList = append(attributeList, "Rebounding")
	}
	if idDiceRoll+potentialModifier > 15 || cp.SpecInteriorDefense {
		attributeList = append(attributeList, "InteriorDefense")
	}
	if pdDiceRoll+potentialModifier > 15 || cp.SpecPerimeterDefense {
		attributeList = append(attributeList, "PerimeterDefense")
	}

	rand.Shuffle(len(attributeList), func(i, j int) {
		attributeList[i], attributeList[j] = attributeList[j], attributeList[i]
	})

	for _, attr := range attributeList {
		if attr == "Shooting2" {
			shooting2 = CollegePlayerProgression(cp.Potential, MinutesPerGame, cp.PlaytimeExpectations, cp.SpecShooting2, cp.IsRedshirting)
		} else if attr == "Shooting3" {
			shooting3 = CollegePlayerProgression(cp.Potential, MinutesPerGame, cp.PlaytimeExpectations, cp.SpecShooting3, cp.IsRedshirting)
		} else if attr == "FreeThrow" {
			freeThrow = CollegePlayerProgression(cp.Potential, MinutesPerGame, cp.PlaytimeExpectations, cp.SpecFreeThrow, cp.IsRedshirting)
		} else if attr == "Finishing" {
			finishing = CollegePlayerProgression(cp.Potential, MinutesPerGame, cp.PlaytimeExpectations, cp.SpecFinishing, cp.IsRedshirting)
		} else if attr == "Ballwork" {
			ballwork = CollegePlayerProgression(cp.Potential, MinutesPerGame, cp.PlaytimeExpectations, cp.SpecBallwork, cp.IsRedshirting)
		} else if attr == "Rebounding" {
			rebounding = CollegePlayerProgression(cp.Potential, MinutesPerGame, cp.PlaytimeExpectations, cp.SpecRebounding, cp.IsRedshirting)
		} else if attr == "InteriorDefense" {
			interiorDefense = CollegePlayerProgression(cp.Potential, MinutesPerGame, cp.PlaytimeExpectations, cp.SpecInteriorDefense, cp.IsRedshirting)
		} else if attr == "PerimeterDefense" {
			perimeterDefense = CollegePlayerProgression(cp.Potential, MinutesPerGame, cp.PlaytimeExpectations, cp.SpecPerimeterDefense, cp.IsRedshirting)
		}
	}

	// Primary Progressions
	staminaCheck := ProgressStamina(cp.Stamina, 0)

	progressions := structs.CollegePlayerProgressions{
		Shooting2:        shooting2,
		Shooting3:        shooting3,
		Ballwork:         ballwork,
		FreeThrow:        freeThrow,
		Finishing:        finishing,
		Rebounding:       rebounding,
		InteriorDefense:  interiorDefense,
		PerimeterDefense: perimeterDefense,
		Stamina:          staminaCheck,
	}

	cp.Progress(progressions)

	return cp
}

func PlayerProgression(progression int, ageDifference int, mpg int, mr int, spec bool, isGleague bool) int {
	min := 0
	max := 0

	progressionCheck := util.GenerateIntFromRange(1, 100)
	if progressionCheck < progression {
		max = 1
	}

	if spec || progressionCheck < progression-25 {
		max = 2
	}

	regressionMax := 0
	if ageDifference > 0 {
		if ageDifference < 4 {
			regressionMax = ageDifference
		} else if ageDifference > 3 {
			regressionMax = 4
		}
		max = max - regressionMax
		min = min - regressionMax
	}

	if mpg < mr && !isGleague {
		diff := mr - mpg
		if diff >= 10 {
			regressionMax += 3
		} else if diff > 5 {
			regressionMax += 2
		} else if diff > 1 {
			regressionMax += 1
		}
		if max > 0 {
			max = 0
		}
		min = min - regressionMax
	}

	if spec && max > 0 {
		min = 1
	}
	return util.GenerateIntFromRange(min, max)
}

func ProgressStamina(stamina int, ageDifference int) int {
	min := -1
	max := 2
	if ageDifference > 0 && ageDifference < 3 {
		min = -2
		max = 1
	} else if ageDifference > 2 && ageDifference < 7 {
		min = -3
		max = 0
	} else if ageDifference > 6 {
		min = -5
		max = 0
	}

	return stamina + util.GenerateIntFromRange(min, max)
}

func CollegePlayerProgression(progression int, mpg int, minutesRequired int, spec bool, isRedshirting bool) int {
	min := 0
	max := 0

	progressionCheck := util.GenerateIntFromRange(1, 100)
	if progressionCheck < progression {
		max = 1
	}

	if spec || progressionCheck < progression-25 {
		max = 2
	}

	if mpg < minutesRequired && !isRedshirting {
		diff := minutesRequired - mpg
		regressionMax := 0
		if diff >= 10 {
			regressionMax = 3
		} else if diff > 5 {
			regressionMax = 2
		} else if diff > 1 {
			regressionMax = 1
		}

		max = 0
		min = min - regressionMax
	}
	if spec && max > 0 {
		min = 1
	}

	return util.GenerateIntFromRange(min, max)
}

func GetModifiers(position string, mpg int, attrib string) float64 {
	var minuteMod float64 = 0.0
	if mpg > 30 {
		minuteMod = rand.Float64()*(1.25-1) + 1
	} else if mpg > 20 {
		minuteMod = rand.Float64()*(1.1-0.9) + 0.9
	} else if mpg > 10 {
		minuteMod = rand.Float64()*(1.0-0.75) + 0.75
	} else if mpg > 5 {
		minuteMod = rand.Float64()*(.9-0.6) + 0.6
	} else {
		minuteMod = rand.Float64()*(0.75-0.5) + 0.5
	}
	return minuteMod
}
