package managers

import (
	"fmt"
	"log"
	"math"
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
				player.FixAge()
				err := db.Save(&player).Error
				if err != nil {
					log.Panicln("Could not save player record")
				}
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
				// draftee := structs.NBADraftee{}
				// draftee.Map(player)
				// draftee.AssignPrimeAge(util.GenerateIntFromRange(25, 30))

				// err := db.Save(&draftee).Error
				// if err != nil {
				// 	log.Panicln("Could not save historic player record!")
				// }

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
	nbaTeams = append(nbaTeams, structs.Team{})

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
	// stats := cp.Stats
	// totalMinutes := 0

	// for _, stat := range stats {
	// 	totalMinutes += stat.Minutes
	// }

	// var MinutesPerGame int = 0
	// if len(stats) > 0 {
	// 	MinutesPerGame = totalMinutes / len(stats)
	// }
	age := np.Age + 1
	ageDifference := np.Age - int(np.PrimeAge)
	if ageDifference < 0 {
		ageDifference = 0
	}

	shooting2 := PlayerProgression(np.Potential, np.Shooting2, ageDifference)
	shooting3 := PlayerProgression(np.Potential, np.Shooting3, ageDifference)
	freeThrow := PlayerProgression(np.Potential, np.FreeThrow, ageDifference)
	ballwork := PlayerProgression(np.Potential, np.Ballwork, ageDifference)
	rebounding := PlayerProgression(np.Potential, np.Rebounding, ageDifference)
	interiorDefense := PlayerProgression(np.Potential, np.InteriorDefense, ageDifference)
	perimeterDefense := PlayerProgression(np.Potential, np.PerimeterDefense, ageDifference)
	finishing := PlayerProgression(np.Potential, np.Finishing, ageDifference)
	stamina := ProgressStamina(np.Stamina, ageDifference)
	overall := int((shooting2+shooting3+freeThrow)/3) + ballwork + finishing + rebounding + int((perimeterDefense+interiorDefense)/2)

	progressions := structs.NBAPlayerProgressions{
		Shooting2:        shooting2,
		Shooting3:        shooting3,
		Ballwork:         ballwork,
		Finishing:        finishing,
		Rebounding:       rebounding,
		InteriorDefense:  interiorDefense,
		PerimeterDefense: perimeterDefense,
		FreeThrow:        freeThrow,
		Overall:          overall,
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

	shooting2 := 0
	shooting3 := 0
	finishing := 0
	freeThrow := 0
	ballwork := 0
	rebounding := 0
	interiorDefense := 0
	perimeterDefense := 0

	if cp.Position == "G" {
		// Primary Progressions
		shooting2 = PrimaryProgression(cp.Potential, cp.Shooting2, cp.Position, MinutesPerGame, "Shooting2", cp.IsRedshirting)
		shooting3 = PrimaryProgression(cp.Potential, cp.Shooting3, cp.Position, MinutesPerGame, "Shooting3", cp.IsRedshirting)
		ballwork = PrimaryProgression(cp.Potential, cp.Ballwork, cp.Position, MinutesPerGame, "Ballwork", cp.IsRedshirting)

		// Secondary
		freeThrow = SecondaryProgression(cp.Potential, cp.FreeThrow)
		rebounding = SecondaryProgression(cp.Potential, cp.Rebounding)
		interiorDefense = SecondaryProgression(cp.Potential, cp.InteriorDefense)
		perimeterDefense = SecondaryProgression(cp.Potential, cp.PerimeterDefense)
		finishing = SecondaryProgression(cp.Potential, cp.Finishing)

	} else if cp.Position == "F" {
		// Primary
		shooting2 = PrimaryProgression(cp.Potential, cp.Shooting2, cp.Position, MinutesPerGame, "Shooting2", cp.IsRedshirting)
		rebounding = PrimaryProgression(cp.Potential, cp.Rebounding, cp.Position, MinutesPerGame, "Rebounding", cp.IsRedshirting)
		finishing = PrimaryProgression(cp.Potential, cp.Finishing, cp.Position, MinutesPerGame, "Finishing", cp.IsRedshirting)
		freeThrow = PrimaryProgression(cp.Potential, cp.FreeThrow, cp.Position, MinutesPerGame, "FreeThrow", cp.IsRedshirting)
		perimeterDefense = PrimaryProgression(cp.Potential, cp.PerimeterDefense, cp.Position, MinutesPerGame, "Perimeter Defense", cp.IsRedshirting)
		// Secondary
		interiorDefense = SecondaryProgression(cp.Potential, cp.InteriorDefense)
		shooting3 = SecondaryProgression(cp.Potential, cp.Shooting3)
		ballwork = SecondaryProgression(cp.Potential, cp.Ballwork)

	} else if cp.Position == "C" {
		// Primary
		rebounding = PrimaryProgression(cp.Potential, cp.Rebounding, cp.Position, MinutesPerGame, "Rebounding", cp.IsRedshirting)
		interiorDefense = PrimaryProgression(cp.Potential, cp.InteriorDefense, cp.Position, MinutesPerGame, "Interior Defense", cp.IsRedshirting)
		finishing = PrimaryProgression(cp.Potential, cp.Finishing, cp.Position, MinutesPerGame, "Finishing", cp.IsRedshirting)

		// Secondary
		shooting2 = SecondaryProgression(cp.Potential, cp.Shooting2)
		shooting3 = SecondaryProgression(cp.Potential, cp.Shooting3)
		ballwork = SecondaryProgression(cp.Potential, cp.Ballwork)
		freeThrow = SecondaryProgression(cp.Potential, cp.FreeThrow)
	}

	overall := (int((shooting2 + shooting3 + freeThrow) / 3)) + finishing + ballwork + rebounding + int((interiorDefense+perimeterDefense)/2)

	progressions := structs.CollegePlayerProgressions{
		Shooting2:        shooting2,
		Shooting3:        shooting3,
		Ballwork:         ballwork,
		FreeThrow:        freeThrow,
		Finishing:        finishing,
		Rebounding:       rebounding,
		InteriorDefense:  interiorDefense,
		PerimeterDefense: perimeterDefense,
		Overall:          overall,
	}

	cp.Progress(progressions)

	return cp
}

func PlayerProgression(progression int, input int, ageDifference int) int {
	min := -1
	max := 1
	if progression > 74 {
		max = 4
	} else if progression > 56 {
		max = 3
	} else if progression > 38 {
		max = 2
	}

	regressionMax := 0
	if ageDifference > 0 && ageDifference < 4 {
		regressionMax = ageDifference
	} else if ageDifference > 3 {
		regressionMax = 4
	}
	max = max - regressionMax
	min = min - regressionMax

	return input + util.GenerateIntFromRange(min, max)
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

func PrimaryProgression(progression int, input int, position string, mpg int, attribute string, isRedshirting bool) int {
	if input == 0 {
		return 1
	}

	modifier := GetModifiers(position, mpg, attribute)

	var progress float64 = 0

	if !isRedshirting {
		progress = ((1 - math.Pow((float64(input)/99.0), 15)) * math.Log10(float64(input)) * (0.3 + modifier)) * (1 + (float64(progression) / 70))
	} else {
		progress = ((1 - math.Pow((float64(input)/99), 15)) * math.Log10(float64(input)) * 1.115 * (1 + (float64(progression / 60))))
	}

	if progress+float64(input) > 20 {
		progress = 20
	} else {
		progress = progress + float64(input)
	}

	return int(math.Round(progress))
}

func SecondaryProgression(progression int, input int) int {
	num := rand.Intn(99)

	if num < progression && input < 20 {
		input = input + util.GenerateIntFromRange(1, 4)
		return input
	} else {
		return input
	}
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
