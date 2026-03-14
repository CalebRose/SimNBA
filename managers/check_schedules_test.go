package managers

import (
	"fmt"
	"testing"
)

type slotKey struct {
	week int
	slot string
}

func checkScheduleEntries(t *testing.T, name string, numTeams int, entries []scheduleEntry) {
	seen := make(map[slotKey]map[int]string)
	counts := make(map[slotKey]int)
	expected := numTeams / 2

	for _, e := range entries {
		k := slotKey{int(e.week), e.slot}
		counts[k]++
		if seen[k] == nil {
			seen[k] = make(map[int]string)
		}
		label := fmt.Sprintf("(h=%d,a=%d)", e.homeIdx, e.awayIdx)
		for _, idx := range []int{e.homeIdx, e.awayIdx} {
			if prev, conflict := seen[k][idx]; conflict {
				t.Errorf("[%s] DUPLICATE: Week %d Slot %s — team %d in %s AND %s",
					name, e.week, e.slot, idx, prev, label)
			} else {
				seen[k][idx] = label
			}
		}
	}

	for k, c := range counts {
		if c != expected {
			t.Errorf("[%s] BAD GAME COUNT: Week %d Slot %s has %d games (expected %d)",
				name, k.week, k.slot, c, expected)
		}
	}
}

func TestElevenTeamSchedules(t *testing.T) {
	t.Run("18Game", func(t *testing.T) {
		checkScheduleEntries(t, "18G", 11, elevenTeam18GameSchedule)
	})
	t.Run("20Game", func(t *testing.T) {
		checkScheduleEntries(t, "20G", 11, elevenTeam20GameSchedule)
	})
}
