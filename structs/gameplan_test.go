package structs

import "testing"

func TestGameplanLineupClearPlayerClearsAllAssignments(t *testing.T) {
	lineup := GameplanLineup{
		FirstStringID:      7,
		FSMinutes:          10,
		FSInsideProportion: 25,
		FSMidProportion:    35,
		FSThreeProportion:  40,
		FSShotVolume:       5,
		SecondStringID:     7,
		SSMinutes:          5,
		SSInsideProportion: 20,
		SSMidProportion:    30,
		SSThreeProportion:  50,
		SSShotVolume:       1,
		TSShotVolume:       2,
		ThirdStringID:      8,
		TSMinutes:          1,
	}

	if !lineup.ClearPlayer(7) {
		t.Fatal("expected assigned player to be cleared")
	}

	if lineup.FirstStringID != 0 || lineup.FSMinutes != 0 ||
		lineup.FSInsideProportion != 0 || lineup.FSMidProportion != 0 || lineup.FSThreeProportion != 0 || lineup.FSShotVolume != 0 ||
		lineup.SecondStringID != 0 || lineup.SSMinutes != 0 ||
		lineup.SSInsideProportion != 0 || lineup.SSMidProportion != 0 || lineup.SSThreeProportion != 0 || lineup.SSShotVolume != 0 {
		t.Fatal("expected every cleared assignment setting to be reset")
	}
	if lineup.ThirdStringID != 8 || lineup.TSMinutes != 1 || lineup.TSShotVolume != 2 {
		t.Fatal("expected unrelated assignments to remain unchanged")
	}
}
