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
	roll := GenerateIntFromRange(1, 10)
	if roll > 9 {
		return "C"
	}
	return PickFromStringList([]string{"PG", "SG", "PF", "SF"})
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
	mean := 18.0
	standardDeviation := 1.5 // Adjust this value to change the spread

	value := rand.NormFloat64()*standardDeviation + mean

	// Clamp the value to the range [16, 22]
	if value <= 16 {
		value = 16
	} else if value > 22 {
		value = 22
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

func GeneratePotential() int {
	num := GenerateIntFromRange(1, 100)

	if num < 10 {
		return GenerateIntFromRange(1, 20)
	} else if num < 20 {
		return GenerateIntFromRange(21, 40)
	} else if num < 80 {
		return GenerateIntFromRange(41, 55)
	} else if num < 85 {
		return GenerateIntFromRange(56, 65)
	} else if num < 90 {
		return GenerateIntFromRange(66, 75)
	} else if num < 95 {
		return GenerateIntFromRange(76, 85)
	} else {
		return GenerateIntFromRange(86, 99)
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

func GetSpecialties(pos string) []string {
	chance := GenerateIntFromRange(0, 9)
	if chance < 1 {
		return []string{}
	}

	list := []string{}
	mod := 0
	diceRoll := 0

	// S2
	diceRoll = GenerateIntFromRange(1, 20)
	if pos == "SG" || pos == "SF" || pos == "PG" {
		mod = 2
	}
	if diceRoll+mod > 13 {
		list = append(list, "SpecShooting2")
	}

	// S3
	diceRoll = GenerateIntFromRange(1, 20)
	if pos == "PG" || pos == "SG" {
		mod = 2
	}
	if diceRoll+mod > 13 {
		list = append(list, "SpecShooting3")
	}
	// FT
	diceRoll = GenerateIntFromRange(1, 20)
	if pos == "SG" || pos == "SF" || pos == "PF" {
		mod = 2
	} else {
		mod = -1
	}
	if diceRoll+mod > 13 {
		list = append(list, "SpecFreeThrow")
	}

	// FN
	diceRoll = GenerateIntFromRange(1, 20)
	if pos == "SF" || pos == "C" || pos == "SG" {
		mod = 2
	}
	if diceRoll > 13 {
		list = append(list, "SpecFinishing")
	}

	// BW
	diceRoll = GenerateIntFromRange(1, 20)
	if pos == "PG" || pos == "SG" {
		mod = 2
	}
	if diceRoll+mod > 13 {
		list = append(list, "SpecBallwork")
	}

	// RB
	diceRoll = GenerateIntFromRange(1, 20)
	if pos == "C" || pos == "PF" {
		mod = 2
	}
	if diceRoll+mod > 13 {
		list = append(list, "SpecRebounding")
	}

	// ID
	diceRoll = GenerateIntFromRange(1, 20)
	if pos == "C" || pos == "PF" {
		mod = 2
	}
	if diceRoll+mod > 13 {
		list = append(list, "SpecInteriorDefense")
	}

	// PD
	diceRoll = GenerateIntFromRange(1, 20)
	if pos == "SG" || pos == "SF" {
		mod = 2
	}
	if diceRoll+mod > 13 {
		list = append(list, "SpecPerimeterDefense")
	}
	return list
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

func GetOverallGrade(rating int) string {
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
	if grade == "A+" {
		return 1
	} else if grade == "A" {
		return 2
	} else if grade == "A-" {
		return 3
	} else if grade == "B+" {
		return 4
	} else if grade == "B" {
		return 5
	} else if grade == "B-" {
		return 6
	} else if grade == "C+" {
		return 7
	} else if grade == "C" {
		return 8
	} else if grade == "C-" {
		return 9
	} else if grade == "D+" {
		return 10
	} else if grade == "D" {
		return 11
	} else if grade == "D-" {
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

func GetAttributeGrade(rating int) string {
	if rating > 15 {
		return "A"
	}
	if rating > 10 {
		return "B"
	}
	if rating > 5 {
		return "C"
	}
	return "D"
}

func GetDrafteeGrade(rating int) string {
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

func GetPlayerOverallGrade(rating int) string {
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

func GetWeightedPotentialGrade(rating int) string {
	weightedRating := GenerateIntFromRange(rating-15, rating+15)
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
	if grade == "A+" {
		return GenerateIntFromRange(88, 100)
	} else if grade == "A" {
		return GenerateIntFromRange(81, 88)
	} else if grade == "A-" {
		return GenerateIntFromRange(75, 80)
	} else if grade == "B+" {
		return GenerateIntFromRange(69, 74)
	} else if grade == "B" {
		return GenerateIntFromRange(63, 68)
	} else if grade == "B-" {
		return GenerateIntFromRange(57, 62)
	} else if grade == "C+" {
		return GenerateIntFromRange(51, 56)
	} else if grade == "C" {
		return GenerateIntFromRange(45, 50)
	} else if grade == "C-" {
		return GenerateIntFromRange(39, 44)
	} else if grade == "D+" {
		return GenerateIntFromRange(33, 38)
	} else if grade == "D" {
		return GenerateIntFromRange(27, 32)
	} else if grade == "D-" {
		return GenerateIntFromRange(21, 26)
	}
	return GenerateIntFromRange(1, 20)
}

func GetPotentialGrade(rating int) string {

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
	if overall > 60 {
		mod = GenerateIntFromRange(1, 3)
	}
	if stars == 5 {
		if year == 4 {
			return GenerateIntFromRange(15, 23) + mod
		}
		return GenerateIntFromRange(10, 22) + mod
	} else if stars == 4 {
		if year == 4 {
			return GenerateIntFromRange(8, 15) + mod
		}
		return GenerateIntFromRange(7, 15) + mod
	} else if stars == 3 {
		if year == 4 {
			return GenerateIntFromRange(7, 11) + mod
		}
		return GenerateIntFromRange(1, 10) + mod
	} else if stars == 2 {
		if year == 4 {
			return GenerateIntFromRange(4, 8) + mod
		} else if year == 3 {
			return GenerateIntFromRange(1, 6) + mod
		}
		return GenerateIntFromRange(1, 5) + mod
	} else {
		return 1 + mod
	}
}

func GetProfessionalPlaytimeExpectations(age, primeage, overall int) int {
	mod := calculateOverallModifier(overall)
	if age < 23 {
		mod -= 5
	} else if age >= primeage {
		mod -= (age - primeage)
	}
	if overall < 80 {
		return GenerateIntFromRange(0, 12) + mod
	} else if overall < 90 {
		return GenerateIntFromRange(8, 18) + mod
	}

	// Superstar Players
	return GenerateIntFromRange(10, 25) + mod
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

func WillPlayerRetire(age int, overall int) bool {
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
	if str == "1" {
		return "1st Round"
	} else if str == "2" {
		return "2nd Round"
	} else if str == "3" {
		return "3rd Round"
	} else if str == "4" {
		return "4th Round"
	} else if str == "5" {
		return "5th Round"
	} else if str == "6" {
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

func GetLotteryChances(idx int) uint {
	if idx <= 4 {
		return 140
	}
	if idx == 5 {
		return 125
	}
	if idx == 6 {
		return 105
	}
	if idx == 7 {
		return 90
	}
	if idx == 8 {
		return 75
	}
	if idx == 9 {
		return 60
	}
	if idx == 10 {
		return 45
	}
	if idx == 11 {
		return 30
	}
	if idx == 12 {
		return 20
	}
	if idx == 13 {
		return 15
	}
	if idx == 14 {
		return 9
	}
	if idx == 15 {
		return 4
	}
	if idx == 16 {
		return 2
	}
	return 1
}

func GetAttributeNew(position, attribute string, spec bool) int {
	mod := 0
	if spec {
		mod = 4
	}
	if position == "PG" || position == "SG" {
		if attribute == "Shooting2" || attribute == "Shooting3" ||
			attribute == "Ballwork" {
			mod += GenerateIntFromRange(1, 2)
		} else if attribute == "Rebounding" || attribute == "Interior Defense" {
			mod -= GenerateIntFromRange(0, 1)
		}
	} else if position == "SG" || position == "SF" {
		if attribute == "Perimeter Defense" {
			mod += GenerateIntFromRange(1, 2)
		}
	} else if position == "PF" || position == "SF" {
		if attribute == "Finishing" {
			mod += GenerateIntFromRange(1, 2)
		} else if attribute == "Shooting3" {
			mod -= GenerateIntFromRange(0, 1)
		}
	} else if position == "C" {
		if attribute == "Finishing" || attribute == "Interior Defense" ||
			attribute == "Rebounding" {
			mod += GenerateIntFromRange(1, 2)
		} else if attribute == "Shooting2" || attribute == "Shooting3" ||
			attribute == "FreeThrow" || attribute == "Ballwork" {
			mod -= GenerateIntFromRange(0, 1)
		}
	}
	return GenerateIntFromRange(3, 14) + mod
}
