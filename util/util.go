package util

import "math/rand"

func GenerateIntFromRange(min int, max int) int {
	return rand.Intn(max-min+1) + min
}

func PickFromStringList(list []string) string {
	return list[rand.Intn(len(list))]
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
	if rating > 37 {
		return "A"
	}
	if rating > 35 {
		return "B"
	}
	if rating > 33 {
		return "C"
	}
	if rating > 30 {
		return "D"
	}
	return "F"
}

func GetDefenseGrade(rating int) string {
	if rating > 36 {
		return "A"
	}
	if rating > 33 {
		return "B"
	}
	if rating > 30 {
		return "C"
	}
	if rating > 27 {
		return "D"
	}
	return "F"
}

func GetOverallGrade(rating int) string {
	if rating > 36 {
		return "A"
	}
	if rating > 33 {
		return "B"
	}
	if rating > 31 {
		return "C"
	}
	if rating > 29 {
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
