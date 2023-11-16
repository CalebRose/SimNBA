package util

import (
	"math"

	"github.com/CalebRose/SimNBA/structs"
)

func FilterOutRecruitingProfile(profiles []structs.PlayerRecruitProfile, ID int) []structs.PlayerRecruitProfile {
	var rp []structs.PlayerRecruitProfile

	for _, profile := range profiles {
		if int(profile.ID) != ID {
			rp = append(rp, profile)
		}
	}

	return rp
}

func Get247TeamRanking(rp structs.TeamRecruitingProfile, signedCroots []structs.Recruit) float64 {
	stddev := 10

	var Rank247 float64 = 0

	for idx, croot := range signedCroots {

		rank := float64((idx - 1) / stddev)

		expo := (-0.5 * (math.Pow(rank, 2)))

		weightedScore := (croot.Rank247 - 20) * math.Pow(math.E, expo)

		Rank247 += (weightedScore)
	}

	return Rank247
}

func GetESPNTeamRanking(rp structs.TeamRecruitingProfile, signedCroots []structs.Recruit) float64 {

	var espnRank float64 = 0

	for _, croot := range signedCroots {
		espnRank += croot.ESPNRank
	}

	return espnRank
}

func GetRivalsTeamRanking(rp structs.TeamRecruitingProfile, signedCroots []structs.Recruit) float64 {

	var rivalsRank float64 = 0

	for _, croot := range signedCroots {
		rivalsRank += croot.RivalsRank
	}

	return rivalsRank
}

func GetRegionMap() map[string]string {
	return map[string]string{
		"Alaska":         "Pacific",
		"California":     "Pacific",
		"Hawai'i":        "Pacific",
		"Idaho":          "Pacific",
		"Nevada":         "Pacific",
		"Oregon":         "Pacific",
		"Utah":           "Pacific",
		"Washington":     "Pacific",
		"Arizona":        "Southwest",
		"Arkansas":       "Southwest",
		"New Mexico":     "Southwest",
		"Oklahoma":       "Southwest",
		"Texas":          "Southwest",
		"Colorado":       "Plains",
		"Kansas":         "Plains",
		"Montana":        "Plains",
		"Nebraska":       "Plains",
		"North Dakota":   "Plains",
		"South Dakota":   "Plains",
		"Wyoming":        "Plains",
		"Illinois":       "Midwest",
		"Indiana":        "Midwest",
		"Iowa":           "Midwest",
		"Kentucky":       "Midwest",
		"Michigan":       "Midwest",
		"Minnesota":      "Midwest",
		"Missouri":       "Midwest",
		"Ohio":           "Midwest",
		"Wisconsin":      "Midwest",
		"Alabama":        "Southeast",
		"Florida":        "Southeast",
		"Georgia":        "Southeast",
		"Louisiana":      "Southeast",
		"Mississippi":    "Southeast",
		"North Carolina": "Southeast",
		"South Carolina": "Southeast",
		"Tennessee":      "Southeast",
		"Delaware":       "Mid-Atlantic",
		"Maryland":       "Mid-Atlantic",
		"New Jersey":     "Mid-Atlantic",
		"New York":       "Mid-Atlantic",
		"Pennsylvania":   "Mid-Atlantic",
		"Virginia":       "Mid-Atlantic",
		"West Virginia":  "Mid-Atlantic",
		"DC":             "Mid-Atlantic",
		"Connecticut":    "Northeast",
		"Maine":          "Northeast",
		"Massachusetts":  "Northeast",
		"New Hampshire":  "Northeast",
		"Rhode Island":   "Northeast",
		"Vermont":        "Northeast",
	}
}

func IsPlayerOffensivelyStrong(r structs.Recruit) bool {
	if r.Stars == 3 && (r.Shooting2 > 12 || r.Shooting3 > 12 || r.Finishing > 12) {
		return true
	} else if r.Stars == 2 && (r.Shooting2 > 10 || r.Shooting3 > 10 || r.Finishing > 10) {
		return true
	} else if r.Stars == 1 && (r.Shooting2 > 8 || r.Shooting3 > 8 || r.Finishing > 8) {
		return true
	}
	return false
}

func IsPlayerDefensivelyStrong(r structs.Recruit) bool {
	if r.Stars == 3 && (r.Rebounding > 12 || r.Defense > 12) {
		return true
	} else if r.Stars == 2 && (r.Rebounding > 10 || r.Defense > 10) {
		return true
	} else if r.Stars == 1 && (r.Rebounding > 8 || r.Defense > 8) {
		return true
	}
	return false
}

func IsPlayerHighPotential(r structs.Recruit) bool {
	return r.Potential > 70
}

func IsAITeamContendingForCroot(profiles []structs.PlayerRecruitProfile) int {
	if len(profiles) == 0 {
		return 0
	}
	leadingVal := 0
	for _, profile := range profiles {
		if profile.TotalPoints != 0 && profile.TotalPoints > float64(leadingVal) {
			leadingVal = int(profile.TotalPoints)
		}
	}

	return leadingVal
}

func RivalsModifiers() []int {
	return []int{100, 83, 82, 81, 80,
		76, 75, 74, 73, 72,
		69, 68, 67, 66, 65,
		64, 63, 62, 61, 60,
		59, 58, 57, 56, 55,
		53, 53, 53, 53, 53,
		51, 51, 51, 51, 51,
		49, 49, 49, 49, 49,
		47, 47, 47, 47, 47,
		45, 45, 45, 45, 45,
		43, 43, 43, 43, 43,
		41, 41, 41, 41, 41,
		40, 40, 40, 40, 40,
		39, 39, 39, 39, 39,
		38, 38, 38, 38, 38,
		37, 37, 37, 37, 37,
		36, 36, 36, 36, 36,
		35, 35, 35, 35, 35,
		34, 34, 34, 34, 34,
		33, 33, 33, 33, 33,
		32, 32, 32, 32, 32,
		31, 31, 31, 31, 31,
		30, 30, 30, 30, 30,
		29, 29, 29, 29, 29,
		28, 28, 28, 28, 28,
		27, 27, 27, 27, 27,
		26, 26, 26, 26, 26,
		25, 25, 25, 25, 25,
		24, 24, 24, 24, 24,
		23, 23, 23, 23, 23,
		22, 22, 22, 22, 22,
		21, 21, 21, 21, 21,
		20, 20, 20, 20, 20,
		19, 19, 19, 19, 19,
		18, 18, 18, 18, 18,
		17, 17, 17, 17, 17,
		16, 16, 16, 16, 16,
		15, 15, 15, 15, 15,
		14, 14, 14, 14, 14,
		13, 13, 13, 13, 13,
		12, 12, 12, 12, 12,
		11, 11, 11, 11, 11,
		10, 10, 10, 10, 10,
		9, 9, 9, 9, 9,
		8, 8, 8, 8, 8,
		7, 7, 7, 7, 7,
		6, 6, 6, 6, 6,
		5, 5, 5, 5, 5,
		4, 4, 4, 4, 4,
		3, 3, 3, 3, 3,
	}
}

func ESPNModifiers() map[string]map[string]string {
	return map[string]map[string]string{
		"C": {
			"Height": "7-0",
		},
		"PG": {
			"Height": "6-3",
		},
		"SG": {
			"Height": "6-5",
		},
		"PF": {
			"Height": "6-8",
		},
		"SF": {
			"Height": "6-7",
		},
	}
}
