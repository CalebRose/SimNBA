package managers

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/repository"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
	"gorm.io/gorm"
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
		roster := GetCollegePlayersByTeamIdForProgression(teamID, ts)
		croots := GetSignedRecruitsByTeamProfileID(teamID)

		for _, player := range roster {
			processCollegePlayer(player, ts, db)
		}

		for _, croot := range croots {
			// Convert to College Player Record
			repository.CreateCollegePlayerRecord(croot, db, true)
		}
	}

	// Unsigned Players
	forgottenPlayersID := "0"
	roster := GetCollegePlayersByTeamIdForProgression(forgottenPlayersID, ts)
	for _, player := range roster {
		if player.PreviousTeamID == 368 {
			continue
		}
		processCollegePlayer(player, ts, db)
	}

	croots := GetAllUnsignedRecruits()
	for _, croot := range croots {
		repository.CreateCollegePlayerRecord(croot, db, true)
	}
	ts.ToggleCollegeProgression()
	db.Save(&ts)
}

func ProgressNBAPlayers() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	if ts.ProgressedProfessionalPlayers {
		return
	}
	fmt.Println(time.Now().UnixNano())

	nbaTeams := GetAllActiveNBATeams()
	// Append empty team object to the end for Free Agents
	nbaTeams = append(nbaTeams, structs.NBATeam{})

	for _, team := range nbaTeams {
		teamID := strconv.Itoa(int(team.ID))

		roster := GetAllNBAPlayersByTeamID(teamID)

		for _, player := range roster {
			if player.HasProgressed {
				continue
			}
			playerID := strconv.Itoa(int(player.ID))
			player = ProgressNBAPlayer(player, false)
			contract := GetNBAContractsByPlayerID(playerID)
			// Retiring Logic
			willPlayerRetire := util.WillPlayerRetire(player.Age, player.Overall)
			isInternationalProspect := player.IsIntGenerated
			willDeclare := InternationalDeclaration(player, isInternationalProspect)
			if willDeclare {
				message := player.TeamAbbr + " " + player.Position + " " + player.FirstName + " " + player.LastName + " from " + player.Country + " has declared for the SimNBA Draft!"
				CreateNewsLog("NBA", message, "Draft", 0, ts)
				// Create NBA Draftee Record
				player.BecomeInternationalDraftee()
				// Void contract
				contract.DeactivateContract()
				repository.SaveProfessionalPlayerRecord(player, db)
				repository.SaveProfessionalContractRecord(contract, db)
				repository.CreateInternationalDrafteeRecord(player, db)
			} else if willPlayerRetire {
				player.SetRetiringStatus()
				message := player.TeamAbbr + " " + player.Position + " " + player.FirstName + " " + player.LastName + " has announced his retirement. He retires at " + strconv.Itoa(player.Age) + " years old and played professionally for " + strconv.Itoa(player.Year) + " years."
				CreateNewsLog("NBA", message, "Retirement", 0, ts)
				retiringPlayer := (structs.RetiredPlayer)(player)
				contract.RetireContract()
				repository.SaveProfessionalContractRecord(contract, db)
				repository.CreateRetireeRecord(retiringPlayer, db)
				repository.DeleteProfessionalPlayerRecord(player, db)
			} else {
				if (player.IsMVP || player.IsDPOY || player.IsFirstTeamANBA) && player.Overall > 90 {
					player.QualifyForSuperMax()
				} else if player.Overall > 94 {
					player.QualifiesForMax()
				} else {
					player.DoesNotQualify()
				}
				contract.ProgressContract()
				if contract.YearsRemaining == 0 && !contract.IsActive && contract.IsComplete {
					extensions := GetExtensionOffersByPlayerID(playerID)
					acceptedExtension := structs.NBAExtensionOffer{}
					for _, e := range extensions {
						if !e.IsAccepted {
							repository.DeleteExtension(e, db)
							continue
						}
						acceptedExtension = e
						break
					}
					if acceptedExtension.ID > 0 {
						contract.MapFromExtension(acceptedExtension)
						player.AssignMinimumContractValue(contract.ContractValue)
						message := "Breaking News: " + player.Position + " " + player.FirstName + " " + player.LastName + " has official signed his extended offer with " + player.TeamAbbr + " for $" + strconv.Itoa(int(contract.ContractValue)) + " Million Dollars!"
						CreateNewsLog("NBA", message, "Free Agency", int(player.TeamID), ts)
						repository.DeleteExtension(acceptedExtension, db)
					} else {
						player.BecomeFreeAgent()
					}
				}
				if contract.ID > 0 {
					repository.SaveProfessionalContractRecord(contract, db)
				}
				repository.SaveProfessionalPlayerRecord(player, db)
			}
		}
	}
	ts.ToggleProfessionalProgression()
	repository.SaveTimeStamp(ts, db)
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
	ageDifference := age - int(np.PrimeAge)
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
	pointLimit := GetPointLimit(np.Potential)
	count := 0

	s2DiceRoll := util.GenerateIntFromRange(1, 20)
	s3DiceRoll := util.GenerateIntFromRange(1, 20)
	ftDiceRoll := util.GenerateIntFromRange(1, 20)
	fnDiceRoll := util.GenerateIntFromRange(1, 20)
	rbDiceRoll := util.GenerateIntFromRange(1, 20)
	bwDiceRoll := util.GenerateIntFromRange(1, 20)
	idDiceRoll := util.GenerateIntFromRange(1, 20)
	pdDiceRoll := util.GenerateIntFromRange(1, 20)

	potentialModifier := np.Potential / 20 // Guaranteed to be between 1-5

	if np.SpecShooting2 {
		attributeList = append(attributeList, "Shooting2")
	}

	if np.SpecShooting3 {
		attributeList = append(attributeList, "Shooting3")
	}
	if np.SpecFreeThrow {
		attributeList = append(attributeList, "FreeThrow")
	}
	if np.SpecFinishing {
		attributeList = append(attributeList, "Finishing")
	}
	if np.SpecBallwork {
		attributeList = append(attributeList, "Ballwork")
	}
	if np.SpecRebounding {
		attributeList = append(attributeList, "Rebounding")
	}
	if np.SpecInteriorDefense {
		attributeList = append(attributeList, "InteriorDefense")
	}
	if np.SpecPerimeterDefense {
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

	developingPlayer := np.IsGLeague || isISLGen || (np.TeamID > 32 && np.Age <= 18 && MinutesPerGame == 0)

	for _, attr := range attributeList {
		if count >= pointLimit {
			break
		}
		allocation := 0
		if attr == "Shooting2" {
			allocation = NBAPlayerProgression(np.Potential, ageDifference, MinutesPerGame, np.PlaytimeExpectations, np.SpecShooting2, developingPlayer)
			if allocation > 0 && count+allocation > pointLimit {
				allocation = pointLimit - count
			}
			shooting2 += allocation
		} else if attr == "Shooting3" {
			allocation = NBAPlayerProgression(np.Potential, ageDifference, MinutesPerGame, np.PlaytimeExpectations, np.SpecShooting3, developingPlayer)
			if allocation > 0 && count+allocation > pointLimit {
				allocation = pointLimit - count
			}
			shooting3 += allocation
		} else if attr == "FreeThrow" {
			allocation = NBAPlayerProgression(np.Potential, ageDifference, MinutesPerGame, np.PlaytimeExpectations, np.SpecFreeThrow, developingPlayer)
			if allocation > 0 && count+allocation > pointLimit {
				allocation = pointLimit - count
			}
			freeThrow += allocation
		} else if attr == "Finishing" {
			allocation = NBAPlayerProgression(np.Potential, ageDifference, MinutesPerGame, np.PlaytimeExpectations, np.SpecFinishing, developingPlayer)
			if allocation > 0 && count+allocation > pointLimit {
				allocation = pointLimit - count
			}
			finishing += allocation
		} else if attr == "Ballwork" {
			allocation = NBAPlayerProgression(np.Potential, ageDifference, MinutesPerGame, np.PlaytimeExpectations, np.SpecBallwork, developingPlayer)
			if allocation > 0 && count+allocation > pointLimit {
				allocation = pointLimit - count
			}
			ballwork += allocation
		} else if attr == "Rebounding" {

			allocation = NBAPlayerProgression(np.Potential, ageDifference, MinutesPerGame, np.PlaytimeExpectations, np.SpecRebounding, developingPlayer)
			if allocation > 0 && count+allocation > pointLimit {
				allocation = pointLimit - count
			}
			rebounding += allocation
		} else if attr == "InteriorDefense" {
			allocation = NBAPlayerProgression(np.Potential, ageDifference, MinutesPerGame, np.PlaytimeExpectations, np.SpecInteriorDefense, developingPlayer)
			if allocation > 0 && count+allocation > pointLimit {
				allocation = pointLimit - count
			}
			interiorDefense += allocation
		} else if attr == "PerimeterDefense" {
			allocation = NBAPlayerProgression(np.Potential, ageDifference, MinutesPerGame, np.PlaytimeExpectations, np.SpecPerimeterDefense, developingPlayer)
			if allocation > 0 && count+allocation > pointLimit {
				allocation = pointLimit - count
			}
			perimeterDefense += allocation
		}
		count += allocation
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

// In the event newly generated players are too talented, move them to the next draft
func MoveISLPlayerToDraft() {
	db := dbprovider.GetInstance().GetDB()
	players := GetAllYouthDevelopmentPlayers()

	for _, p := range players {
		if p.ID != 9204 && p.ID != 9272 && p.ID != 8637 {
			continue
		}
		p.BecomeInternationalDraftee()
		repository.SaveProfessionalPlayerRecord(p, db)
		repository.CreateInternationalDrafteeRecord(p, db)
	}

}

func processCollegePlayer(player structs.CollegePlayer, ts structs.Timestamp, db *gorm.DB) {
	if player.HasProgressed {
		return
	}

	minutesPerGame := getMinutesPlayed(player)
	isSenior := (player.Year == 4 && !player.IsRedshirt) || (player.Year == 5 && player.IsRedshirt)

	player = ProgressCollegePlayer(player, minutesPerGame, false)

	if player.IsRedshirting {
		player.SetRedshirtStatus()
	}

	player.SetExpectations(util.GetPlaytimeExpectations(player.Stars, player.Year, player.Overall))

	if player.WillDeclare || isSenior {
		// Graduate Player
		handlePlayerGraduation(player, ts, db)
	} else {
		// Save Player
		if player.TeamID == 0 {
			player.WillTransfer()
		}
		repository.SaveCollegePlayerRecord(player, db)
	}
}

func handlePlayerGraduation(player structs.CollegePlayer, ts structs.Timestamp, db *gorm.DB) {
	// Graduate Player
	player.GraduatePlayer()

	message := player.Position + " " + player.FirstName + " " + player.LastName + " has graduated from " + player.TeamAbbr + "!"

	// Create News Log
	CreateNewsLog("CBB", message, "Graduation", int(player.TeamID), ts)

	// Make draftee record
	repository.CreateDrafteeRecord(player, db)
	repository.CreateHistoricPlayerRecord(player, db)
	repository.DeleteCollegePlayerRecord(player, db)
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

	threshold := 10
	if s2DiceRoll+potentialModifier >= threshold {
		attributeList = append(attributeList, "Shooting2")
	}

	if s3DiceRoll+potentialModifier >= threshold {
		attributeList = append(attributeList, "Shooting3")
	}
	if ftDiceRoll+potentialModifier >= threshold {
		attributeList = append(attributeList, "FreeThrow")
	}
	if fnDiceRoll+potentialModifier >= threshold {
		attributeList = append(attributeList, "Finishing")
	}
	if bwDiceRoll+potentialModifier >= threshold {
		attributeList = append(attributeList, "Ballwork")
	}
	if rbDiceRoll+potentialModifier >= threshold {
		attributeList = append(attributeList, "Rebounding")
	}
	if idDiceRoll+potentialModifier >= threshold {
		attributeList = append(attributeList, "InteriorDefense")
	}
	if pdDiceRoll+potentialModifier >= threshold {
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

func NBAPlayerProgression(progression int, ageDifference int, mpg int, mr int, spec bool, isGleague bool) int {
	progressionCheck := util.GenerateIntFromRange(1, 100)
	max := calculateMaxProgression(progression, progressionCheck, spec)
	if ageDifference > 0 {
		max = adjustForAge(ageDifference, max)
	}
	if mpg < mr && !isGleague {
		max = adjustForPlaytime(mpg, mr, max)
	}

	min := 0

	if spec && max > 0 {
		min = 1
	}
	if max < min {
		min, max = util.Swap(min, max)
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
	progressionCheck := util.GenerateIntFromRange(1, 100)
	max := calculateMaxProgression(progression, progressionCheck, spec)

	if mpg < minutesRequired && !isRedshirting {
		max = adjustForPlaytime(mpg, minutesRequired, max)
	}
	min := 0

	if spec && max > 0 {
		min = 1
	}
	if max < min {
		min, max = util.Swap(min, max)
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

func calculateMaxProgression(progression, progressionCheck int, spec bool) int {
	maxProgression := 4
	if spec {
		maxProgression += 1
	}
	if progressionCheck <= progression {
		roof := progression / 16
		min := 1
		newMax := util.Min(maxProgression, roof)
		if newMax < min {
			min, newMax = util.Swap(min, newMax)
		}
		return util.GenerateIntFromRange(min, newMax)
	} else if progressionCheck <= progression+25 {
		return 1
	}
	return 0
}

func adjustForAge(ageDifference, max int) int {
	regressionMax := util.Min(ageDifference, 5)
	regressionChange := util.GenerateIntFromRange(1, 10)
	if regressionChange <= ageDifference {
		return max - regressionMax
	}
	return max
}

func adjustForPlaytime(mpg, mr, max int) int {
	diff := mr - mpg
	regressionMax := 0
	if diff == 0 {
		regressionMax = 1
	} else if diff >= 10 {
		regressionMax = 3
	} else if diff > 5 {
		regressionMax = 2
	} else if diff > 1 {
		regressionMax = 1
	}

	regressionChance := util.GenerateIntFromRange(1, 5)
	if regressionChance <= 2 {
		// 4 could be adjusted or parameterized based on design choice
		return max - regressionMax
	}
	return max
}
