package util

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/CalebRose/SimNBA/structs"
)

func ReadCSV(path string) [][]string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal("Unable to read input file "+path, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	rows, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+path, err)
	}

	return rows
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Swap(min, max int) (int, int) {
	return max, min
}

func GenerateIntFromRange(min int, max int) int {
	return rand.Intn(max-min+1) + min
}

func GenerateFloatFromRange(min float64, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func GenerateNormalizedIntFromRange(min int, max int) int {
	mean := float64(min+max) / 2.0
	stdDev := float64(max-min) / 6.0 // This approximates the 3-sigma rule

	for {
		// Generate a number using normal distribution around the mean
		num := rand.NormFloat64()*stdDev + mean
		// Round to nearest integer and convert to int type
		intNum := int(num + 0.5) // Adding 0.5 before truncating simulates rounding
		// Check if the generated number is within bounds
		if intNum >= min && intNum <= max {
			return intNum
		}
		// If not within bounds, loop again
	}
}

func PickPositionFromList() string {
	return PickFromStringList([]string{"G", "G", "G", "F", "F", "F", "C"})
}

func PickFromStringList(list []string) string {
	return list[rand.Intn(len(list))]
}

func GeneratePrimeAge() int {
	chance := GenerateIntFromRange(1, 100)

	if chance < 5 {
		return 22
	} else if chance < 10 {
		return 23
	} else if chance < 15 {
		return 24
	} else if chance < 35 {
		return 25
	} else if chance < 50 {
		return 26
	} else if chance < 55 {
		return 27
	} else if chance < 65 {
		return 28
	} else if chance < 75 {
		return 29
	} else if chance < 80 {
		return 30
	} else if chance < 85 {
		return 31
	} else if chance < 90 {
		return 32
	} else if chance < 95 {
		return 33
	}
	return 34
}

func FormatHeight(height string) string {
	// Split the input string by the dash
	parts := strings.Split(height, "-")

	// Check if the input string is in the correct format
	if len(parts) != 2 {
		return "Invalid format"
	}

	// Construct the formatted height string
	formattedHeight := fmt.Sprintf("%s'%s\"", parts[0], parts[1])

	return formattedHeight
}

func ConvertFloatToString(f float64) string {
	return strconv.FormatFloat(f, 'f', 3, 64)
}

func GenerateISLAge() int {
	mean := 17.0
	standardDeviation := 1.5 // Adjust this value to change the spread

	value := rand.NormFloat64()*standardDeviation + mean

	// Clamp the value to the range [16, 22]
	if value <= 15 {
		value = 15
	} else if value > 20 {
		value = 20
	}

	return int(value)
}

func GenerateStamina() int {
	mean := 32.0
	standardDeviation := 1.6 // Adjust this value to change the spread

	value := rand.NormFloat64()*standardDeviation + mean

	// Clamp the value to the range [25, 38]
	if value < 25 {
		value = 25
	} else if value > 38 {
		value = 38
	}

	return int(value)
}

func GeneratePotential() uint8 {
	num := GenerateIntFromRange(1, 100)

	if num < 10 {
		return uint8(GenerateIntFromRange(1, 20))
	} else if num < 20 {
		return uint8(GenerateIntFromRange(21, 40))
	} else if num < 80 {
		return uint8(GenerateIntFromRange(41, 55))
	} else if num < 85 {
		return uint8(GenerateIntFromRange(56, 65))
	} else if num < 90 {
		return uint8(GenerateIntFromRange(66, 75))
	} else if num < 95 {
		return uint8(GenerateIntFromRange(76, 85))
	} else {
		return uint8(GenerateIntFromRange(86, 99))
	}
}

func GetPersonality() string {
	chance := GenerateIntFromRange(1, 3)
	if chance < 3 {
		return "Average"
	}
	list := []string{"Reserved",
		"Eccentric",
		"Motivational",
		"Disloyal",
		"Cooperative",
		"Irrational",
		"Focused",
		"Book Worm",
		"Motivation",
		"Abrasive",
		"Absent Minded",
		"Uncooperative",
		"Introvert",
		"Disruptive",
		"Outgoing",
		"Tough",
		"Paranoid",
		"Stoic",
		"Dramatic",
		"Extroverted",
		"Selfish",
		"Impatient",
		"Reliable",
		"Frail",
		"Relaxed",
		"Average",
		"Flamboyant",
		"Perfectionist",
		"Popular",
		"Jokester",
		"Narcissist"}

	return PickFromStringList(list)
}

func GetAcademicBias() string {
	chance := GenerateIntFromRange(1, 3)
	if chance < 3 {
		return "Average"
	}
	list := []string{"Takes AP classes",
		"Sits at the front of the class",
		"Seeks out tutoring",
		"Tutor",
		"Wants to finish degree",
		"Teacher's Pet",
		"Sits at the back of the class",
		"Values academics",
		"Studious",
		"Frequent visits to the principal",
		"Class Clown",
		"More likely to get academic probation",
		"Has other priorities",
		"Distracted",
		"Loves Learning",
		"Studies hard",
		"Less likely to get academic probation",
		"Never Studies",
		"Average",
		"Naturally book smart",
		"Borderline failing",
		"Skips classes often",
		"Didn't come here to play school"}

	return PickFromStringList(list)
}

func GetRecruitingBias() string {
	chance := GenerateIntFromRange(1, 3)
	if chance < 3 {
		return "Average"
	}
	list := []string{"Prefers to play in a different state",
		"Prefers to play for an up-and-coming team",
		"Open-Minded",
		"Prefers to play for a team where he can start immediately",
		"Prefers to be close to home",
		"Prefers to play for a national championship contender",
		"Prefers to play for a specific coach",
		"Average",
		"Legacy",
		"Prefers to play for a team with a rich history"}

	return PickFromStringList(list)
}

func GetWorkEthic() string {
	chance := GenerateIntFromRange(1, 3)
	if chance < 3 {
		return "Average"
	}
	list := []string{"Persistant",
		"Lazy",
		"Footwork king",
		"Hard-working",
		"Complacent",
		"Skips Leg Day",
		"Working-Class mentality",
		"Film Room Genius",
		"Focuses on Max Weight",
		"Track Athlete",
		"Average",
		"Center of Attention",
		"Gym Rat",
		"Focuses on Max Reps",
		"Loud",
		"Quiet",
		"Streams too much",
		"Trolls on Discord"}
	return PickFromStringList(list)
}

func GetFreeAgencyBias(age, ovr int) string {
	chance := GenerateIntFromRange(1, 3)
	if chance < 3 {
		return "Average"
	}
	list := []string{
		"Wants extensions",
		"Drafted team discount",
		"Loyal",
		"Hometown hero",
		"Adversarial",
	}

	midAgeList := []string{
		"I'm the starter",
		"Market-driven",
		"Money motivated",
	}

	veteranList := []string{
		"Highest bidder",
		"Championship seeking",
	}

	if age > 30 || ovr > 95 {
		list = append(list, veteranList...)
	} else if age > 24 {
		list = append(list, midAgeList...)
	}

	return PickFromStringList(list)
}

func GetOffenseGrade(rating int) string {
	if rating > 45 {
		return "A+"
	}
	if rating > 42 {
		return "A"
	}
	if rating > 39 {
		return "A-"
	}
	if rating > 36 {
		return "B+"
	}
	if rating > 33 {
		return "B"
	}
	if rating > 30 {
		return "B-"
	}
	if rating > 27 {
		return "C+"
	}
	if rating > 24 {
		return "C"
	}
	if rating > 21 {
		return "C-"
	}
	if rating > 18 {
		return "D+"
	}
	if rating > 15 {
		return "D"
	}
	if rating > 12 {
		return "D-"
	}
	return "F"
}

func GetDefenseGrade(rating int) string {
	if rating > 45 {
		return "A+"
	}
	if rating > 42 {
		return "A"
	}
	if rating > 39 {
		return "A-"
	}
	if rating > 36 {
		return "B+"
	}
	if rating > 33 {
		return "B"
	}
	if rating > 30 {
		return "B-"
	}
	if rating > 27 {
		return "C+"
	}
	if rating > 24 {
		return "C"
	}
	if rating > 21 {
		return "C-"
	}
	if rating > 18 {
		return "D+"
	}
	if rating > 15 {
		return "D"
	}
	if rating > 12 {
		return "D-"
	}
	return "F"
}

func GetOverallGrade(rating uint8) string {
	if rating > 45 {
		return "A+"
	}
	if rating > 42 {
		return "A"
	}
	if rating > 39 {
		return "A-"
	}
	if rating > 36 {
		return "B+"
	}
	if rating > 33 {
		return "B"
	}
	if rating > 30 {
		return "B-"
	}
	if rating > 27 {
		return "C+"
	}
	if rating > 24 {
		return "C"
	}
	if rating > 21 {
		return "C-"
	}
	if rating > 18 {
		return "D+"
	}
	if rating > 15 {
		return "D"
	}
	if rating > 12 {
		return "D-"
	}
	return "F"
}

// FOR 2023 Season ONLY
func GetOverallDraftGrade(rating int) string {
	if rating > 90 {
		return "A+"
	}
	if rating > 87 {
		return "A"
	}
	if rating > 84 {
		return "A-"
	}
	if rating > 81 {
		return "B+"
	}
	if rating > 78 {
		return "B"
	}
	if rating > 75 {
		return "B-"
	}
	if rating > 72 {
		return "C+"
	}
	if rating > 69 {
		return "C"
	}
	if rating > 66 {
		return "C-"
	}
	if rating > 63 {
		return "D+"
	}
	if rating > 60 {
		return "D"
	}
	if rating > 57 {
		return "D-"
	}
	return "F"
}

func GetNumericalSortValueByLetterGrade(grade string) int {
	switch grade {
	case "A+":
		return 1
	case "A":
		return 2
	case "A-":
		return 3
	case "B+":
		return 4
	case "B":
		return 5
	case "B-":
		return 6
	case "C+":
		return 7
	case "C":
		return 8
	case "C-":
		return 9
	case "D+":
		return 10
	case "D":
		return 11
	case "D-":
		return 12
	}
	return 13
}

func GetNBATeamGrade(rating int) string {
	if rating > 89 {
		return "A+"
	}
	if rating > 84 {
		return "A"
	}
	if rating > 79 {
		return "A-"
	}
	if rating > 74 {
		return "B+"
	}
	if rating > 70 {
		return "B"
	}
	if rating > 65 {
		return "B-"
	}
	if rating > 60 {
		return "C+"
	}
	if rating > 55 {
		return "C"
	}
	if rating > 50 {
		return "C-"
	}
	if rating > 45 {
		return "D+"
	}
	if rating > 40 {
		return "D"
	}
	if rating > 35 {
		return "D-"
	}
	return "F"
}

func GetAttributeGrade(rating uint8) string {
	if rating > 16 {
		return "A"
	} else if rating > 13 {
		return "B"
	} else if rating > 10 {
		return "C"
	} else if rating > 7 {
		return "D"
	}
	return "F"
}

func GetDrafteeGrade(rating uint8) string {
	if rating > 24 {
		return "A+"
	}
	if rating > 22 {
		return "A"
	}
	if rating > 20 {
		return "A-"
	}
	if rating > 18 {
		return "B+"
	}
	if rating > 16 {
		return "B"
	}
	if rating > 14 {
		return "B-"
	}
	if rating > 12 {
		return "C+"
	}
	if rating > 10 {
		return "C"
	}
	if rating > 8 {
		return "C-"
	}
	if rating > 5 {
		return "D"
	}
	return "F"
}

func GetPlayerOverallGrade(rating uint8) string {
	if rating > 69 {
		return "A"
	}
	if rating > 56 {
		return "B"
	}
	if rating > 48 {
		return "C"
	}
	if rating > 36 {
		return "D"
	}
	return "F"
}

func GetWeightedPotentialGrade(rating uint8) string {
	weightedRating := GenerateIntFromRange(int(rating)-15, int(rating)+15)
	if weightedRating > 100 {
		weightedRating = 99
	} else if weightedRating < 0 {
		weightedRating = 0
	}

	if weightedRating > 88 {
		return "A+"
	}
	if weightedRating > 80 {
		return "A"
	}
	if weightedRating > 74 {
		return "A-"
	}
	if weightedRating > 68 {
		return "B+"
	}
	if weightedRating > 62 {
		return "B"
	}
	if weightedRating > 56 {
		return "B-"
	}
	if weightedRating > 50 {
		return "C+"
	}
	if weightedRating > 44 {
		return "C"
	}
	if weightedRating > 38 {
		return "C-"
	}
	if weightedRating > 32 {
		return "D+"
	}
	if weightedRating > 26 {
		return "D"
	}
	if weightedRating > 20 {
		return "D-"
	}
	return "F"
}

func GetNBAProgressionRatingFromGrade(grade string) int {
	switch grade {
	case "A+":
		return GenerateIntFromRange(88, 100)
	case "A":
		return GenerateIntFromRange(81, 88)
	case "A-":
		return GenerateIntFromRange(75, 80)
	case "B+":
		return GenerateIntFromRange(69, 74)
	case "B":
		return GenerateIntFromRange(63, 68)
	case "B-":
		return GenerateIntFromRange(57, 62)
	case "C+":
		return GenerateIntFromRange(51, 56)
	case "C":
		return GenerateIntFromRange(45, 50)
	case "C-":
		return GenerateIntFromRange(39, 44)
	case "D+":
		return GenerateIntFromRange(33, 38)
	case "D":
		return GenerateIntFromRange(27, 32)
	case "D-":
		return GenerateIntFromRange(21, 26)
	}
	return GenerateIntFromRange(1, 20)
}

func GetPotentialGrade(rating uint8) string {

	if rating > 88 {
		return "A+"
	}
	if rating > 80 {
		return "A"
	}
	if rating > 74 {
		return "A-"
	}
	if rating > 68 {
		return "B+"
	}
	if rating > 62 {
		return "B"
	}
	if rating > 56 {
		return "B-"
	}
	if rating > 50 {
		return "C+"
	}
	if rating > 44 {
		return "C"
	}
	if rating > 38 {
		return "C-"
	}
	if rating > 32 {
		return "D+"
	}
	if rating > 26 {
		return "D"
	}
	if rating > 20 {
		return "D-"
	}
	return "F"
}

func GetPlaytimeExpectations(stars int, year int, overall int) int {
	mod := 0
	if overall > 30 {
		mod = GenerateIntFromRange(1, 3)
	}
	switch stars {
	case 5:
		if year == 4 {
			return GenerateIntFromRange(15, 23) + mod
		}
		return GenerateIntFromRange(10, 22) + mod
	case 4:
		if year == 4 {
			return GenerateIntFromRange(8, 15) + mod
		}
		return GenerateIntFromRange(7, 15) + mod
	case 3:
		if year == 4 {
			return GenerateIntFromRange(7, 11) + mod
		}
		return GenerateIntFromRange(1, 10) + mod
	case 2:
		switch year {
		case 4:
			return GenerateIntFromRange(4, 8) + mod
		case 3:
			return GenerateIntFromRange(1, 6) + mod
		}
		return GenerateIntFromRange(1, 5) + mod
	default:
		return 1 + mod
	}
}

func GetProfessionalPlaytimeExpectations(age, primeage, overall uint8) int {
	mod := calculateOverallModifier(int(overall))
	if age < 23 {
		mod -= 5
	} else if age >= primeage {
		mod -= (int(age) - int(primeage))
	}
	if overall < 70 {
		return GenerateIntFromRange(0, 12) + mod
	} else if overall < 80 {
		return GenerateIntFromRange(4, 16) + mod
	} else if overall < 90 {
		return GenerateIntFromRange(8, 20) + mod
	}

	// Superstar Players
	return GenerateIntFromRange(12, 28) + mod
}

// calculateOverallModifier - Returns a modifier between 0 and 100 based on the overall of the player
func calculateOverallModifier(overall int) int {
	minOverall := 60
	maxOverall := 100 // Changed to match your specific game
	minModifier := 1
	maxModifier := 10

	// Interpolate between min and max values
	modifier := (overall-minOverall)*(maxModifier-minModifier)/(maxOverall-minOverall) + minModifier

	// Apply the cap
	if modifier > maxModifier {
		return maxModifier
	} else {
		return modifier
	}
}

func ConvertStringToBool(str string) bool {
	return str == "TRUE" || str == "1"
}

func ConvertStringToInt(num string) int {
	if num == "" {
		return 0
	}
	val, err := strconv.Atoi(num)
	if err != nil {
		log.Fatalln("Could not convert string to int")
	}

	return val
}

func ConvertStringToFloat(num string) float64 {
	floatNum, error := strconv.ParseFloat(num, 64)
	if error != nil {
		fmt.Println("Could not convert string to float 64, resetting as 0.")
		return 0
	}
	return floatNum
}

func WillPlayerRetire(age, overall uint8) bool {
	if age > 25 && overall < 60 {
		return true
	}
	if age > 29 && overall < 80 {
		odds := 5
		if age == 31 {
			odds = 15
		} else if age == 32 {
			odds = 25
		} else if age == 33 {
			odds = 35
		} else if age == 34 {
			odds = 45
		} else if age == 35 {
			odds = 55
		} else if age == 36 {
			odds = 65
		} else if age > 36 {
			odds = 75
		}
		chance := GenerateIntFromRange(1, 100)
		if chance < odds {
			return true
		}
	}
	return false
}

func GetRoundAbbreviation(str string) string {
	switch str {
	case "1":
		return "1st Round"
	case "2":
		return "2nd Round"
	case "3":
		return "3rd Round"
	case "4":
		return "4th Round"
	case "5":
		return "5th Round"
	case "6":
		return "6th Round"
	}
	return "7th Round"
}

func GetCollegePlayerIDsBySeasonStats(cps []structs.CollegePlayerSeasonStats) []string {
	var list []string

	for _, cp := range cps {
		list = append(list, strconv.Itoa(int(cp.CollegePlayerID)))
	}

	return list
}

func GetNBAPlayerIDsBySeasonStats(nps []structs.NBAPlayerSeasonStats) []string {
	var list []string

	for _, cp := range nps {
		list = append(list, strconv.Itoa(int(cp.NBAPlayerID)))
	}

	return list
}

func GetLotteryChances(idx int) []uint {
	chancesMap := map[int][]uint{
		1:  {1400, 1340, 1270, 1200},
		2:  {1400, 1340, 1270, 1200},
		3:  {1185, 1185, 1185, 1185},
		4:  {1185, 1185, 1185, 1185},
		5:  {1185, 1185, 1185, 1185},
		6:  {900, 920, 940, 960},
		7:  {750, 780, 810, 840},
		8:  {600, 630, 670, 710},
		9:  {425, 425, 425, 425},
		10: {425, 425, 425, 425},
		11: {200, 220, 240, 260},
		12: {150, 160, 180, 200},
		13: {80, 80, 80, 80},
		14: {80, 80, 80, 80},
		15: {26, 26, 26, 26},
		16: {26, 26, 26, 26},
	}
	return chancesMap[idx]
}

func GetAttributeNew(position, attribute string, spec, isWalkon bool) int {
	mod := 0
	if spec {
		mod = 4
	}
	if position == "PG" || position == "SG" {
		switch attribute {
		case "Shooting2", "Shooting3", "Ballwork":
			mod += GenerateIntFromRange(1, 2)
		case "Rebounding", "Interior Defense":
			mod -= GenerateIntFromRange(0, 1)
		}
	}
	if position == "SG" || position == "SF" {
		if attribute == "Perimeter Defense" {
			mod += GenerateIntFromRange(1, 2)
		}
	}
	if position == "PF" || position == "SF" {
		switch attribute {
		case "Finishing":
			mod += GenerateIntFromRange(1, 2)
		case "Shooting3":
			mod -= GenerateIntFromRange(0, 1)
		}
	}
	if position == "C" {
		switch attribute {
		case "Finishing", "Interior Defense", "Rebounding":
			mod += GenerateIntFromRange(1, 2)
		case "Shooting2", "Shooting3", "FreeThrow", "Ballwork":
			mod -= GenerateIntFromRange(0, 1)
		}
	}
	if isWalkon {
		mod = 0
	}
	return GenerateIntFromRange(3, 14) + mod
}

func GetWeekIDBySeasonAndWeek(season uint, week uint) uint {
	// Format should be SSWW where SS is the last two digits of the season and WW is the week number with leading zeros if necessary
	seasonPart := season % 100
	weekPart := week
	if week < 10 {
		weekPart = week + 100 // This will ensure that when we convert to string, it will have a leading zero
	}

	return uint(seasonPart*100 + weekPart)
}

func GetStarRating(isCustom, isInt bool) int {
	roll := GenerateIntFromRange(1, 1000)
	if isInt {
		if roll < 3 {
			return 5
		}
		if roll < 40 {
			return 4
		}
		if roll < 275 {
			return 3
		}
		if roll < 650 {
			return 2
		}
		return 1
	}
	if isCustom {
		roll -= 100
	}
	if roll < 0 {
		roll = 1
	}
	if roll < 3 {
		return 6
	}
	if roll < 42 {
		return 5
	}
	if roll < 122 {
		return 4
	}
	if roll < 352 || isCustom {
		return 3
	}
	if roll < 652 {
		return 2
	}
	return 1
}

type Locale struct {
	Name   string
	Weight int
}

// Pick a US state or Canadian province for which the player is from
func PickState() string {
	states := []Locale{
		{Name: "TX", Weight: 45}, // Collective weight for less prominent states
		{Name: "CA", Weight: 40}, // Collective weight for less prominent states
		{Name: "NY", Weight: 40},
		{Name: "NC", Weight: 40}, // Collective weight for less prominent states
		{Name: "IL", Weight: 40},
		{Name: "FL", Weight: 40}, // Collective weight for less prominent states
		{Name: "OH", Weight: 30},
		{Name: "PA", Weight: 30},
		{Name: "GA", Weight: 30}, // Collective weight for less prominent states
		{Name: "IN", Weight: 30}, // Collective weight for less prominent states
		{Name: "NJ", Weight: 25}, // Collective weight for less prominent states
		{Name: "VA", Weight: 20}, // Collective weight for less prominent states
		{Name: "AZ", Weight: 20}, // Collective weight for less prominent states
		{Name: "TN", Weight: 20}, // Collective weight for less prominent states
		{Name: "SC", Weight: 20}, // Collective weight for less prominent states
		{Name: "KS", Weight: 20}, // Collective weight for less prominent states
		{Name: "LA", Weight: 20}, // Collective weight for less prominent states
		{Name: "KY", Weight: 15}, // Collective weight for less prominent states
		{Name: "MO", Weight: 15}, // Collective weight for less prominent states
		{Name: "MA", Weight: 15},
		{Name: "AL", Weight: 15}, // Collective weight for less prominent states
		{Name: "UT", Weight: 10}, // Collective weight for less prominent states
		{Name: "MI", Weight: 10},
		{Name: "CO", Weight: 10},
		{Name: "CT", Weight: 10},
		{Name: "MS", Weight: 10}, // Collective weight for less prominent states
		{Name: "OK", Weight: 10}, // Collective weight for less prominent states
		{Name: "WA", Weight: 8},  // Collective weight for less prominent states
		{Name: "VT", Weight: 5},
		{Name: "OR", Weight: 5}, // Collective weight for less prominent states
		{Name: "AK", Weight: 5},
		{Name: "NH", Weight: 5}, // Collective weight for less prominent states
		{Name: "RI", Weight: 5}, // Collective weight for less prominent states
		{Name: "ME", Weight: 5}, // Collective weight for less prominent states
		{Name: "MN", Weight: 5},
		{Name: "WI", Weight: 5},
		{Name: "ND", Weight: 3},
		{Name: "DE", Weight: 3}, // Collective weight for less prominent states
		{Name: "MT", Weight: 3}, // Collective weight for less prominent states
		{Name: "NE", Weight: 2}, // Collective weight for less prominent states
		{Name: "MD", Weight: 1}, // Collective weight for less prominent states
		{Name: "IA", Weight: 1}, // Collective weight for less prominent states
		{Name: "NM", Weight: 1}, // Collective weight for less prominent states
		{Name: "SD", Weight: 1}, // Collective weight for less prominent states
		{Name: "WY", Weight: 1}, // Collective weight for less prominent states
		{Name: "ID", Weight: 1}, // Collective weight for less prominent states
		{Name: "NV", Weight: 1}, // Collective weight for less prominent states
		{Name: "WV", Weight: 1}, // Collective weight for less prominent states
		{Name: "AR", Weight: 1}, // Collective weight for less prominent states
		{Name: "HI", Weight: 1}, // Collective weight for less prominent states
		{Name: "GM", Weight: 1}, // Collective weight for less prominent states
		{Name: "AS", Weight: 1}, // Collective weight for less prominent states
	}

	totalWeight := 0
	for _, state := range states {
		totalWeight += state.Weight
	}

	randomWeight := GenerateIntFromRange(0, totalWeight)
	for _, state := range states {
		if randomWeight < state.Weight {
			return state.Name
		}
		randomWeight -= state.Weight
	}
	return PickFromStringList([]string{"MN", "MI", "NY", "MA"})
}

func GetArchetype(pos string) string {
	diceRoll := GenerateIntFromRange(1, 1000)
	if diceRoll > 998 {
		return "All-Around"
	}
	switch pos {
	case "C":
		return PickFromStringList([]string{"Rim Protector", "Post Scorer", "Stretch Center"})
	case "F":
		return PickFromStringList([]string{"Power Forward", "Small Forward", "Point Forward", "Swingman", "Two-Way"})
	case "G":
		return PickFromStringList([]string{"Point Guard", "Shooting Guard", "Combo Guard", "Slasher", "3-and-D"})
	}
	return ""
}
