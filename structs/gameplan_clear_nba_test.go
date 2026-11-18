package structs

import "testing"

func nbaLineup(position string, first, second, third uint) NBALineup {
	return NBALineup{GameplanLineup: GameplanLineup{
		Position:       position,
		FirstStringID:  first,
		FSMinutes:      20,
		SecondStringID: second,
		SSMinutes:      10,
		ThirdStringID:  third,
		TSMinutes:      3, TSInsideProportion: 40, TSMidProportion: 30, TSThreeProportion: 30,
	}}
}

func TestClearPlayerFromNBALineupsClearsOnlyThatPlayer(t *testing.T) {
	// The San Antonio case: player 3271 is the third-string C and nowhere else.
	lineups := []NBALineup{
		nbaLineup("PG", 1, 2, 3),
		nbaLineup("C", 10, 11, 3271),
		nbaLineup("F", 20, 21, 22),
	}

	changed := ClearPlayerFromNBALineups(lineups, 3271)

	if len(changed) != 1 || changed[0].Position != "C" {
		t.Fatalf("expected only the C lineup to change, got %+v", changed)
	}
	c := lineups[1]
	if c.ThirdStringID != 0 || c.TSMinutes != 0 || c.TSInsideProportion != 0 || c.TSMidProportion != 0 || c.TSThreeProportion != 0 {
		t.Fatalf("the cleared slot and its minutes/shot allocations should be reset, got %+v", c.GameplanLineup)
	}
	if c.FirstStringID != 10 || c.FSMinutes != 20 || c.SecondStringID != 11 || c.SSMinutes != 10 {
		t.Fatalf("the other slots in that lineup must be untouched, got %+v", c.GameplanLineup)
	}
	if lineups[0].ThirdStringID != 3 || lineups[2].ThirdStringID != 22 {
		t.Fatal("lineups that did not contain the player must be untouched")
	}
	if changed[0].ThirdStringID != 0 {
		t.Fatal("the returned lineup should reflect the cleared slot so it can be saved")
	}
}

func TestClearPlayerFromNBALineupsHandlesPlayerInSeveralLineups(t *testing.T) {
	lineups := []NBALineup{
		nbaLineup("F", 7, 8, 9),
		nbaLineup("C", 5, 7, 6),
		nbaLineup("G", 1, 2, 3),
	}
	changed := ClearPlayerFromNBALineups(lineups, 7)
	if len(changed) != 2 {
		t.Fatalf("expected both lineups that named player 7 to be returned, got %d", len(changed))
	}
	if lineups[0].FirstStringID != 0 || lineups[1].SecondStringID != 0 {
		t.Fatal("player 7 should be cleared from both lineups")
	}
	if lineups[0].SecondStringID != 8 || lineups[1].FirstStringID != 5 {
		t.Fatal("other players must stay")
	}
}

func TestClearPlayerFromNBALineupsIgnoresUnknownPlayerAndZeroID(t *testing.T) {
	// Empty slots are stored as 0; asking to clear player 0 must not "change" them.
	lineups := []NBALineup{nbaLineup("C", 10, 0, 0), nbaLineup("F", 20, 21, 22)}

	if changed := ClearPlayerFromNBALineups(lineups, 999); len(changed) != 0 {
		t.Fatalf("a player who is in no lineup should change nothing, got %+v", changed)
	}
	if changed := ClearPlayerFromNBALineups(lineups, 0); len(changed) != 0 {
		t.Fatalf("player ID 0 means an empty slot and must change nothing, got %+v", changed)
	}
	if lineups[0].FirstStringID != 10 || lineups[0].FSMinutes != 20 {
		t.Fatal("lineups must be unchanged")
	}
}

func TestClearPlayerFromNBALineupsWithNoLineups(t *testing.T) {
	if changed := ClearPlayerFromNBALineups(nil, 5); len(changed) != 0 {
		t.Fatalf("expected no changes for no lineups, got %+v", changed)
	}
}
