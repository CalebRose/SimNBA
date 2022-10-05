package managers

import (
	"math"
	"math/rand"

	"github.com/CalebRose/SimNBA/structs"
)

func ProgressPlayer(cp structs.CollegePlayer) structs.CollegePlayer {
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
	ballwork := 0
	rebounding := 0
	defense := 0

	if cp.Position == "G" {
		// Primary Progressions
		shooting2 = PrimaryProgression(cp.Potential, cp.Shooting2, cp.Position, MinutesPerGame, "Shooting2", cp.IsRedshirting)
		shooting3 = PrimaryProgression(cp.Potential, cp.Shooting3, cp.Position, MinutesPerGame, "Shooting3", cp.IsRedshirting)
		ballwork = PrimaryProgression(cp.Potential, cp.Ballwork, cp.Position, MinutesPerGame, "Ballwork", cp.IsRedshirting)

		// Secondary
		rebounding = SecondaryProgression(cp.Potential, cp.Rebounding)
		defense = SecondaryProgression(cp.Potential, cp.Defense)
		finishing = SecondaryProgression(cp.Potential, cp.Finishing)

	} else if cp.Position == "F" {
		// Primary
		shooting2 = PrimaryProgression(cp.Potential, cp.Shooting2, cp.Position, MinutesPerGame, "Shooting2", cp.IsRedshirting)
		rebounding = PrimaryProgression(cp.Potential, cp.Rebounding, cp.Position, MinutesPerGame, "Rebounding", cp.IsRedshirting)
		finishing = PrimaryProgression(cp.Potential, cp.Finishing, cp.Position, MinutesPerGame, "Finishing", cp.IsRedshirting)
		// Secondary
		defense = SecondaryProgression(cp.Potential, cp.Defense)
		shooting3 = SecondaryProgression(cp.Potential, cp.Shooting3)
		ballwork = SecondaryProgression(cp.Potential, cp.Ballwork)

	} else if cp.Position == "C" {
		// Primary
		rebounding = PrimaryProgression(cp.Potential, cp.Rebounding, cp.Position, MinutesPerGame, "Rebounding", cp.IsRedshirting)
		defense = PrimaryProgression(cp.Potential, cp.Defense, cp.Position, MinutesPerGame, "Defense", cp.IsRedshirting)
		finishing = PrimaryProgression(cp.Potential, cp.Finishing, cp.Position, MinutesPerGame, "Finishing", cp.IsRedshirting)

		// Secondary
		shooting2 = SecondaryProgression(cp.Potential, cp.Shooting2)
		shooting3 = SecondaryProgression(cp.Potential, cp.Shooting3)
		ballwork = SecondaryProgression(cp.Potential, cp.Ballwork)
	}

	overall := int((shooting2+shooting3)/2) + ballwork + finishing + rebounding + defense

	progressions := structs.CollegePlayerProgressions{
		Shooting2:  shooting2,
		Shooting3:  shooting3,
		Ballwork:   ballwork,
		Finishing:  finishing,
		Rebounding: rebounding,
		Defense:    defense,
		Overall:    overall,
	}

	cp.Progress(progressions)

	return cp
}

func PrimaryProgression(progression int, input int, position string, spg int, attribute string, isRedshirting bool) int {
	if input == 0 {
		return 1
	}

	modifier := GetModifiers(position, spg, attribute)

	var progress float64 = 0

	if !isRedshirting {
		progress = ((1 - math.Pow((float64(input)/99.0), 15)) * math.Log10(float64(input)) * (0.3 + modifier)) * (1 + (float64(progression) / 70))
	} else {
		progress = ((1 - math.Pow((float64(input)/99), 15)) * math.Log10(float64(input)) * 1.125 * (1 + (float64(progression / 60))))
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
		input = input + 1
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
