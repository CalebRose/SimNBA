package managers

import (
	"strings"
	"testing"

	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
)

func play(event, outcome uint8) structs.BasePlayByPlay {
	return structs.BasePlayByPlay{EventID: event, OutcomeID: outcome, BallCarrierID: 7}
}

func TestFreeThrowNote(t *testing.T) {
	ft := func() structs.BasePlayByPlay { return play(util.Free_throw, util.Ft_made) }
	sub := play(util.Substitution, util.No_outcome)

	tests := []struct {
		name  string
		plays []structs.BasePlayByPlay
		idx   int
		want  string
	}{
		{"not a free throw", []structs.BasePlayByPlay{play(util.Shot_three, util.Shot_made)}, 0, ""},
		{"three attempts, first", []structs.BasePlayByPlay{play(util.Shot_three, util.Shot_foul_missed), ft(), ft(), ft()}, 1, " for free throw 1 of 3"},
		{"three attempts, last", []structs.BasePlayByPlay{play(util.Shot_three, util.Shot_foul_missed), ft(), ft(), ft()}, 3, " for free throw 3 of 3"},
		{"two attempts", []structs.BasePlayByPlay{play(util.Shot_midrange, util.Shot_foul_missed), ft(), ft()}, 2, " for free throw 2 of 2"},
		{"substitution between attempts", []structs.BasePlayByPlay{play(util.Shot_inside, util.Shot_foul_missed), ft(), sub, ft()}, 3, " for free throw 2 of 2"},
		{"and-one", []structs.BasePlayByPlay{play(util.Shot_paint, util.Shot_foul_made), ft()}, 1, " for the and-one"},
		{"missed front end of one-and-one", []structs.BasePlayByPlay{play(util.Move, util.Move_foul), ft()}, 1, " for the front end of a one-and-one"},
		{"bonus trip of two", []structs.BasePlayByPlay{play(util.Move, util.Move_foul), ft(), ft()}, 2, " for free throw 2 of 2"},
	}
	for _, tc := range tests {
		if got := freeThrowNote(tc.plays, tc.idx); got != tc.want {
			t.Errorf("%s: got %q, want %q", tc.name, got, tc.want)
		}
	}
}

func TestMissedShootingFoulFreeThrowCount(t *testing.T) {
	labelFor := func(uint) string { return "" }
	for event, want := range map[uint8]string{
		util.Shot_three: "three free throws", util.Shot_corner_three: "three free throws",
		util.Shot_midrange: "two free throws", util.Shot_inside: "two free throws", util.Shot_paint: "two free throws",
	} {
		got := generatePlayStatement(play(event, util.Shot_foul_missed), "", labelFor, "Team", "Team")
		if !strings.Contains(got, want) {
			t.Errorf("event %d: %q does not contain %q", event, got, want)
		}
	}
}
