package managers

import (
	"log"
	"math"
	"math/rand"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
)

func AssignAllRecruitRanks() {
	db := dbprovider.GetInstance().GetDB()

	var recruits []structs.Recruit

	// var recruitsToSync []structs.Recruit

	db.Order("overall desc").Find(&recruits)

	rivalsModifiers := util.RivalsModifiers()

	for idx, croot := range recruits {
		// 247 Rankings
		rank247 := Get247Ranking(croot)
		// ESPN Rankings
		espnRank := GetESPNRanking(croot)

		// Rivals Ranking
		var rivalsRank float64 = 0
		if idx <= 249 {
			rivalsBonus := rivalsModifiers[idx]

			rivalsRank = GetRivalsRanking(croot.Stars, rivalsBonus)
		}

		var r float64 = croot.TopRankModifier

		if croot.TopRankModifier == 0 || croot.TopRankModifier < 0.95 || croot.TopRankModifier > 1.05 {
			r = 0.95 + rand.Float64()*(1.05-0.95)
		}

		croot.AssignRankValues(rank247, espnRank, rivalsRank, r)

		db.Save(&croot)

		// recruitsToSync = append(recruitsToSync, croot)
	}

}

func Get247Ranking(r structs.Recruit) float64 {
	ovr := r.Overall

	potentialGrade := Get247PotentialModifier(r.PotentialGrade)

	specGrade := float64(r.SpecCount) * 0.15

	return float64(ovr) + (potentialGrade * 2) + specGrade
}

func GetESPNRanking(r structs.Recruit) float64 {
	// ESPN Ranking = Star Rank + Archetype Modifier + weight difference + height difference
	// + potential val, and then round.

	starRank := GetESPNStarRank(r.Stars)
	potentialMod := GetESPNPotentialModifier(r.PotentialGrade)

	espnPositionMap := util.ESPNModifiers()
	espnHeight := getInches(espnPositionMap[r.Position]["Height"])
	playerHeight := getInches(r.Height)
	var heightMod float64 = float64(playerHeight / espnHeight)
	espnRanking := math.Round(float64(starRank) + potentialMod + heightMod)

	return espnRanking
}

func getInches(height string) int {
	feet := 0
	inches := 0
	pastDash := false
	for idx, char := range height {
		if string(char) != "-" {
			if !pastDash {
				str := string(char)
				ft, err := strconv.Atoi(str)
				if err != nil {
					log.Panic("Could not convert height to inches")
				}
				feet = ft
			} else {
				str := height[idx:]
				inc, err := strconv.Atoi(str)
				if err != nil {
					log.Panic("Could not convert height to inches")
				}
				inches = inc
			}
		} else {
			pastDash = true
		}
	}
	return (feet * 12) + inches
}

func GetRivalsRanking(stars int, bonus int) float64 {
	return GetRivalsStarModifier(stars) + float64(bonus)
}

func GetESPNStarRank(star int) int {
	switch star {
	case 5:
		return 95
	case 4:
		return 85
	case 3:
		return 75
	case 2:
		return 65
	}
	return 55
}

func GetArchetypeModifier(arch string) int {
	switch arch {
	case "Coverage", "Run Stopper", "Ball Hawk", "Man Coverage", "Pass Rusher", "Rushing":
		return 1
	case "Possession", "Field General", "Nose Tackle", "Blocking", "Line Captain":
		return -1
	case "Speed Rusher", "Pass Rush", "Scrambler", "Vertical Threat", "Speed":
		return 2
	}
	return 0
}

func Get247PotentialModifier(pg string) float64 {
	switch pg {
	case "A+":
		return 5.83
	case "A":
		return 5.06
	case "A-":
		return 4.77
	case "B+":
		return 4.33
	case "B":
		return 4.04
	case "B-":
		return 3.87
	case "C+":
		return 3.58
	case "C":
		return 3.43
	case "C-":
		return 3.31
	case "D+":
		return 3.03
	case "D":
		return 2.77
	case "D-":
		return 2.67
	}
	return 2.3
}

func GetESPNPotentialModifier(pg string) float64 {
	switch pg {
	case "A+":
		return 1
	case "A":
		return 0.9
	case "A-":
		return 0.8
	case "B+":
		return 0.6
	case "B":
		return 0.4
	case "B-":
		return 0.2
	case "C+":
		return 0
	case "C":
		return -0.15
	case "C-":
		return -0.3
	case "D+":
		return -0.6
	case "D":
		return -0.75
	case "D-":
		return -0.9
	}
	return -1
}

func GetPredictiveOverall(r structs.Recruit) int {
	currentOverall := r.Overall

	var potentialProg int

	switch r.PotentialGrade {
	case "B+", "A-", "A", "A+":
		potentialProg = 7
	case "B", "B-", "C+":
		potentialProg = 5
	default:
		potentialProg = 4
	}

	return currentOverall + (potentialProg * 3)
}

func GetRivalsStarModifier(stars int) float64 {
	switch stars {
	case 5:
		return 6.1
	case 4:
		return RoundToFixedDecimalPlace(rand.Float64()*((6.0-5.8)+5.8), 1)
	case 3:
		return RoundToFixedDecimalPlace(rand.Float64()*((5.7-5.5)+5.5), 1)
	case 2:
		return RoundToFixedDecimalPlace(rand.Float64()*((5.4-5.2)+5.2), 1)
	default:
		return 5
	}
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func RoundToFixedDecimalPlace(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
