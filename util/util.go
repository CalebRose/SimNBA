package util

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
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

func GenerateIntFromRange(min int, max int) int {
	return rand.Intn(max-min+1) + min
}

func PickFromStringList(list []string) string {
	return list[rand.Intn(len(list))]
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
		"Prefers to play with former teammates",
		"Fanboy",
		"Multi-Sport",
		"Prefers to play for a team where he can start immediately",
		"Is going to school mainly for academics",
		"Prefers to be close to home",
		"Prefers to play for a national championship contender",
		"Prefers to play for a specific coach",
		"Average",
		"Legacy",
		"Prefers to play for a team with a rich history",
		"Wants that NIL money"}

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

func GetSpecialties(pos string) []string {
	chance := GenerateIntFromRange(0, 9)
	if chance < 2 {
		return []string{}
	}

	list := []string{}
	mod := 0
	diceRoll := 0

	// S2
	diceRoll = GenerateIntFromRange(1, 20)
	if pos == "G" || pos == "F" {
		mod = 3
	}
	if diceRoll+mod > 15 {
		list = append(list, "SpecShooting2")
	}

	// S3
	diceRoll = GenerateIntFromRange(1, 20)
	if pos == "G" {
		mod = 3
	}
	if diceRoll+mod > 15 {
		list = append(list, "SpecShooting3")
	}
	// FT
	diceRoll = GenerateIntFromRange(1, 20)
	if pos == "G" || pos == "F" {
		mod = 3
	} else {
		mod = -1
	}
	if diceRoll+mod > 15 {
		list = append(list, "SpecFreeThrow")
	}

	// FN
	diceRoll = GenerateIntFromRange(1, 20)
	if pos == "F" || pos == "C" {
		mod = 3
	}
	if diceRoll > 15 {
		list = append(list, "SpecFinishing")
	}

	// BW
	diceRoll = GenerateIntFromRange(1, 20)
	if pos == "G" {
		mod = 4
	}
	if diceRoll+mod > 15 {
		list = append(list, "SpecBallwork")
	}

	// RB
	diceRoll = GenerateIntFromRange(1, 20)
	if pos == "C" || pos == "F" {
		mod = 4
	}
	if diceRoll+mod > 15 {
		list = append(list, "SpecRebounding")
	}

	// ID
	diceRoll = GenerateIntFromRange(1, 20)
	if pos == "C" || pos == "F" {
		mod = 3
	}
	if diceRoll+mod > 15 {
		list = append(list, "SpecInteriorDefense")
	}

	// PD
	diceRoll = GenerateIntFromRange(1, 20)
	if pos == "G" || pos == "F" {
		mod = 3
	}
	if diceRoll+mod > 15 {
		list = append(list, "SpecPerimeterDefense")
	}
	return list
}

func GetFreeAgencyBias() string {
	chance := GenerateIntFromRange(1, 3)
	if chance < 3 {
		return "Average"
	}
	list := []string{"I'm the starter",
		"Market-driven",
		"Wants extensions",
		"Drafted team discount",
		"Highest bidder",
		"Championship seeking",
		"Loyal",
		"Average",
		"Hometown hero",
		"Money motivated",
		"Hates Tags",
		"Adversarial",
		"Hates Cleveland",
		"Will eventually play for LeBron"}

	return PickFromStringList(list)
}

func GetOffenseGrade(rating int) string {
	if rating > 40 {
		return "A"
	}
	if rating > 35 {
		return "B"
	}
	if rating > 32 {
		return "C"
	}
	if rating > 30 {
		return "D"
	}
	return "F"
}

func GetDefenseGrade(rating int) string {
	if rating > 40 {
		return "A"
	}
	if rating > 35 {
		return "B"
	}
	if rating > 32 {
		return "C"
	}
	if rating > 27 {
		return "D"
	}
	return "F"
}

func GetOverallGrade(rating int) string {
	if rating > 40 {
		return "A"
	}
	if rating > 35 {
		return "B"
	}
	if rating > 32 {
		return "C"
	}
	if rating > 30 {
		return "D"
	}
	return "F"
}

func GetNBATeamGrade(rating int) string {
	if rating > 74 {
		return "A+"
	}
	if rating > 70 {
		return "A"
	}
	if rating > 65 {
		return "A-"
	}
	if rating > 60 {
		return "B+"
	}
	if rating > 55 {
		return "B"
	}
	if rating > 50 {
		return "B-"
	}
	if rating > 45 {
		return "C+"
	}
	if rating > 40 {
		return "C"
	}
	if rating > 35 {
		return "C-"
	}
	if rating > 30 {
		return "D"
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

func GetPlaytimeExpectations(stars int, year int) int {
	if stars == 5 {
		if year == 4 {
			return GenerateIntFromRange(15, 29)
		} else if year == 3 {
			return GenerateIntFromRange(10, 25)
		} else if year == 2 {
			return GenerateIntFromRange(10, 20)
		}
		return GenerateIntFromRange(10, 20)
	} else if stars == 4 {
		if year == 4 {
			return GenerateIntFromRange(15, 25)
		} else if year == 3 {
			return GenerateIntFromRange(9, 20)
		} else if year == 2 {
			return GenerateIntFromRange(5, 17)
		}
		return GenerateIntFromRange(5, 15)
	} else if stars == 3 {
		if year == 4 {
			return GenerateIntFromRange(7, 21)
		} else if year == 3 {
			return GenerateIntFromRange(3, 17)
		} else if year == 2 {
			return GenerateIntFromRange(2, 13)
		}
		return GenerateIntFromRange(0, 10)
	} else if stars == 2 {
		if year == 4 {
			return GenerateIntFromRange(0, 13)
		} else if year == 3 {
			return GenerateIntFromRange(0, 13)
		} else if year == 2 {
			return GenerateIntFromRange(0, 9)
		}
		return GenerateIntFromRange(0, 6)
	} else {
		return GenerateIntFromRange(0, 5)
	}
}

func ConvertStringToBool(str string) bool {
	return str == "TRUE"
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
