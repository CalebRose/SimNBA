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
	ts := GetTimestamp()

	collegeTeams := GetAllActiveCollegeTeams()

	for _, team := range collegeTeams {
		// var graduatingPlayers []structs.NBADraftee
		teamID := strconv.Itoa(int(team.ID))
		// roster := GetAllCollegePlayersWithStatsByTeamID(teamID, SeasonID)
		roster := GetCollegePlayersByTeamIdForProgression(teamID)
		croots := GetSignedRecruitsByTeamProfileID(teamID)
		recruitingProfile := GetOnlyTeamRecruitingProfileByTeamID(teamID)

		for _, player := range roster {
			if player.HasProgressed {
				// player.FixAge()
				// err := db.Save(&player).Error
				// if err != nil {
				// 	log.Panicln("Could not save player record")
				// }
				continue
			}

			minutesPerGame := getMinutesPlayed(player)

			willTransfer := determineTransferStatus(minutesPerGame, player)

			status := determineStatusLevel(recruitingProfile, player)
			isSenior := (player.Year == 4 && !player.IsRedshirt) || (player.Year == 5 && player.IsRedshirt)
			player = ProgressCollegePlayer(player, minutesPerGame, false)
			if player.IsRedshirting {
				player.SetRedshirtStatus()
			}

			player.SetExpectations(util.GetPlaytimeExpectations(player.Stars, player.Year, player.Overall))

			if player.WillDeclare && isSenior {
				player.GraduatePlayer()

				message := player.Position + " " + player.FirstName + " " + player.LastName + " has graduated from " + player.TeamAbbr + "!"

				CreateNewsLog("CBB", message, "Graduation", int(player.TeamID), ts)

				// Make draftee record
				draftee := structs.NBADraftee{}
				draftee.Map(player)
				draftee.AssignPrimeAge(util.GenerateIntFromRange(25, 30))

				err := db.Save(&draftee).Error
				if err != nil {
					log.Panicln("Could not save historic player record!")
				}

				hcp := (structs.HistoricCollegePlayer)(player)

				err = db.Save(&hcp).Error
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
				message := ""
				if player.WillDeclare {
					message = player.Position + " " + player.FirstName + " " + player.LastName + " plans to declare early from " + player.TeamAbbr + " for the upcoming SimNBA Draft!"
					CreateNewsLog("CBB", message, "Graduation", int(player.TeamID), ts)
				} else if willTransfer {
					player.DeclareTransferIntention(status)
					message = player.Position + " " + player.FirstName + " " + player.LastName + " intends to transfer from " + player.TeamAbbr + "."
					CreateNewsLog("CBB", message, "Transfer Portal", int(player.TeamID), ts)
				}

				err := db.Save(&player).Error
				if err != nil {
					log.Panicln("Could not save player record")
				}
			}

		}

		for _, croot := range croots {
			// Convert to College Player Record
			cp := structs.CollegePlayer{}
			cp.MapFromRecruit(croot)

			// Save College Player Record
			err := db.Create(&cp).Error
			if err != nil {
				log.Panicln("Could not save new college player record")
			}

			// Delete Recruit Record
			db.Delete(&croot)
		}

	}
	ts.ToggleCollegeProgression()
	ts.NextTransferPortalPhase()
	db.Save(&ts)
}

func ProgressNBAPlayers() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	fmt.Println(time.Now().UnixNano())

	nbaTeams := GetAllActiveNBATeams()
	// Append empty team object to the end for Free Agents
	nbaTeams = append(nbaTeams, structs.NBATeam{})

	for _, team := range nbaTeams {
		teamID := strconv.Itoa(int(team.ID))

		roster := GetAllNBAPlayersByTeamID(teamID)

		for _, player := range roster {
			playerID := strconv.Itoa(int(player.ID))
			player = ProgressNBAPlayer(player, false)

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

				db.Save(&contract)
				db.Save(&player)
			}
		}
	}
	ts.ToggleProfessionalProgression()
	db.Save(&ts)
}

func ProgressNBAPlayer(np structs.NBAPlayer, isISLGen bool) structs.NBAPlayer {
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

	developingPlayer := np.IsGLeague || isISLGen

	for _, attr := range attributeList {
		if attr == "Shooting2" {
			shooting2 = PlayerProgression(np.Potential, ageDifference, MinutesPerGame, np.PlaytimeExpectations, np.SpecShooting2, developingPlayer)
		} else if attr == "Shooting3" {
			shooting3 = PlayerProgression(np.Potential, ageDifference, MinutesPerGame, np.PlaytimeExpectations, np.SpecShooting3, developingPlayer)
		} else if attr == "FreeThrow" {
			freeThrow = PlayerProgression(np.Potential, ageDifference, MinutesPerGame, np.PlaytimeExpectations, np.SpecFreeThrow, developingPlayer)
		} else if attr == "Finishing" {
			finishing = PlayerProgression(np.Potential, ageDifference, MinutesPerGame, np.PlaytimeExpectations, np.SpecFinishing, developingPlayer)
		} else if attr == "Ballwork" {
			ballwork = PlayerProgression(np.Potential, ageDifference, MinutesPerGame, np.PlaytimeExpectations, np.SpecBallwork, developingPlayer)
		} else if attr == "Rebounding" {
			rebounding = PlayerProgression(np.Potential, ageDifference, MinutesPerGame, np.PlaytimeExpectations, np.SpecRebounding, developingPlayer)
		} else if attr == "InteriorDefense" {
			interiorDefense = PlayerProgression(np.Potential, ageDifference, MinutesPerGame, np.PlaytimeExpectations, np.SpecInteriorDefense, developingPlayer)
		} else if attr == "PerimeterDefense" {
			perimeterDefense = PlayerProgression(np.Potential, ageDifference, MinutesPerGame, np.PlaytimeExpectations, np.SpecPerimeterDefense, developingPlayer)
		}
	}

	stamina := ProgressStamina(np.Stamina, ageDifference, true)

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

func ProgressCollegePlayer(cp structs.CollegePlayer, mpg int, isGeneration bool) structs.CollegePlayer {
	var MinutesPerGame int = mpg

	if isGeneration {
		MinutesPerGame = 100
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

	pointLimit := GetPointLimit(cp.Potential)
	count := 0

	s2DiceRoll := util.GenerateIntFromRange(1, 20)
	s3DiceRoll := util.GenerateIntFromRange(1, 20)
	ftDiceRoll := util.GenerateIntFromRange(1, 20)
	fnDiceRoll := util.GenerateIntFromRange(1, 20)
	rbDiceRoll := util.GenerateIntFromRange(1, 20)
	bwDiceRoll := util.GenerateIntFromRange(1, 20)
	idDiceRoll := util.GenerateIntFromRange(1, 20)
	pdDiceRoll := util.GenerateIntFromRange(1, 20)

	potentialModifier := cp.Potential / 20 // Guaranteed to be between 1-5

	if cp.SpecShooting2 {
		attributeList = append(attributeList, "Shooting2")
	}
	if cp.SpecShooting3 {
		attributeList = append(attributeList, "Shooting3")
	}
	if cp.SpecFreeThrow {
		attributeList = append(attributeList, "FreeThrow")
	}
	if cp.SpecFinishing {
		attributeList = append(attributeList, "Finishing")
	}
	if cp.SpecBallwork {
		attributeList = append(attributeList, "Ballwork")
	}
	if cp.SpecRebounding {
		attributeList = append(attributeList, "Rebounding")
	}
	if cp.SpecInteriorDefense {
		attributeList = append(attributeList, "InteriorDefense")
	}
	if cp.SpecPerimeterDefense {
		attributeList = append(attributeList, "PerimeterDefense")
	}

	if s2DiceRoll+potentialModifier >= 15 {
		attributeList = append(attributeList, "Shooting2")
	}

	if s3DiceRoll+potentialModifier >= 15 {
		attributeList = append(attributeList, "Shooting3")
	}
	if ftDiceRoll+potentialModifier >= 15 {
		attributeList = append(attributeList, "FreeThrow")
	}
	if fnDiceRoll+potentialModifier >= 15 {
		attributeList = append(attributeList, "Finishing")
	}
	if bwDiceRoll+potentialModifier >= 15 {
		attributeList = append(attributeList, "Ballwork")
	}
	if rbDiceRoll+potentialModifier >= 15 {
		attributeList = append(attributeList, "Rebounding")
	}
	if idDiceRoll+potentialModifier >= 15 {
		attributeList = append(attributeList, "InteriorDefense")
	}
	if pdDiceRoll+potentialModifier >= 15 {
		attributeList = append(attributeList, "PerimeterDefense")
	}

	rand.Shuffle(len(attributeList), func(i, j int) {
		attributeList[i], attributeList[j] = attributeList[j], attributeList[i]
	})

	for _, attr := range attributeList {
		if count >= pointLimit {
			break
		}
		allocation := 0
		if attr == "Shooting2" {
			allocation = CollegePlayerProgression(cp.Potential, MinutesPerGame, cp.PlaytimeExpectations, cp.SpecShooting2, cp.IsRedshirting)
			if count+allocation > pointLimit {
				allocation = pointLimit - count
			}
			shooting2 += allocation
		} else if attr == "Shooting3" {
			allocation = CollegePlayerProgression(cp.Potential, MinutesPerGame, cp.PlaytimeExpectations, cp.SpecShooting3, cp.IsRedshirting)
			if count+allocation > pointLimit {
				allocation = pointLimit - count
			}
			shooting3 += allocation
		} else if attr == "FreeThrow" {
			allocation = CollegePlayerProgression(cp.Potential, MinutesPerGame, cp.PlaytimeExpectations, cp.SpecFreeThrow, cp.IsRedshirting)
			if count+allocation > pointLimit {
				allocation = pointLimit - count
			}
			freeThrow += allocation
		} else if attr == "Finishing" {
			allocation = CollegePlayerProgression(cp.Potential, MinutesPerGame, cp.PlaytimeExpectations, cp.SpecFinishing, cp.IsRedshirting)
			if count+allocation > pointLimit {
				allocation = pointLimit - count
			}
			finishing += allocation
		} else if attr == "Ballwork" {
			allocation = CollegePlayerProgression(cp.Potential, MinutesPerGame, cp.PlaytimeExpectations, cp.SpecBallwork, cp.IsRedshirting)
			if count+allocation > pointLimit {
				allocation = pointLimit - count
			}
			ballwork += allocation
		} else if attr == "Rebounding" {
			allocation = CollegePlayerProgression(cp.Potential, MinutesPerGame, cp.PlaytimeExpectations, cp.SpecRebounding, cp.IsRedshirting)
			if count+allocation > pointLimit {
				allocation = pointLimit - count
			}
			rebounding += allocation
		} else if attr == "InteriorDefense" {
			allocation = CollegePlayerProgression(cp.Potential, MinutesPerGame, cp.PlaytimeExpectations, cp.SpecInteriorDefense, cp.IsRedshirting)
			if count+allocation > pointLimit {
				allocation = pointLimit - count
			}
			interiorDefense += allocation
		} else if attr == "PerimeterDefense" {
			allocation = CollegePlayerProgression(cp.Potential, MinutesPerGame, cp.PlaytimeExpectations, cp.SpecPerimeterDefense, cp.IsRedshirting)
			if count+allocation > pointLimit {
				allocation = pointLimit - count
			}
			perimeterDefense += allocation
		}
		count += allocation
	}

	// Primary Progressions
	staminaCheck := ProgressStamina(cp.Stamina, 0, false)

	potentialGrade := util.GetWeightedPotentialGrade(cp.Potential)

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
		PotentialGrade:   potentialGrade,
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

func ProgressStamina(stamina, ageDifference int, isNBA bool) int {
	min := -1
	if !isNBA {
		min = 1
	} else if isNBA && ageDifference == 0 {
		min = 0
	}
	max := 5
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

	newStamina := stamina + util.GenerateIntFromRange(min, max)
	if newStamina > 48 {
		newStamina = 48
	}
	return newStamina
}

func CollegePlayerProgression(progression int, mpg int, minutesRequired int, spec bool, isRedshirting bool) int {
	min := 0
	max := 0

	progressionCheck := util.GenerateIntFromRange(1, 100)
	if progressionCheck <= progression {
		max = 1
	}

	if spec || progressionCheck <= progression-25 {
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

func GetPointLimit(pot int) int {
	floater := float64(pot)
	floor := floater / 10
	roundUp := math.Ceil(floor)
	if roundUp < 1 {
		floor = 1
	}
	roof := int(roundUp) + util.GenerateIntFromRange(0, 1)
	if roof > 10 {
		roof = 10
	}
	return util.GenerateIntFromRange(int(roundUp), roof)
}

func getMinutesPlayed(cp structs.CollegePlayer) int {
	stats := cp.Stats
	totalMinutes := 0

	for _, stat := range stats {
		totalMinutes += stat.Minutes
	}

	var MinutesPerGame int = 0
	if len(stats) > 0 {
		MinutesPerGame = totalMinutes / len(stats)
	}

	return MinutesPerGame
}

func determineTransferStatus(mpg int, cp structs.CollegePlayer) bool {
	if cp.WillDeclare || cp.IsRedshirting {
		return false
	}

	return mpg < cp.PlaytimeExpectations
}

func determineStatusLevel(rtp structs.TeamRecruitingProfile, cp structs.CollegePlayer) string {
	if cp.WillDeclare {
		return "High"
	}

	ovrRate := cp.Overall / 10
	mod := 0
	if rtp.AIQuality == "Blue Blood" {
		mod += 3
	} else if rtp.AIQuality == "P6" {
		mod += 2
	} else {
		mod += 1
	}

	sum := ovrRate + mod

	if sum > 7 {
		return "High"
	} else if sum > 3 {
		return "Medium"
	}
	return "Low"
}
