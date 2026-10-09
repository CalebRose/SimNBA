package managers

import (
	"testing"

	"github.com/CalebRose/SimNBA/structs"
)

func TestFillLineupSlotsSetsShotVolumeNormal(t *testing.T) {
	// Slots start with non-Normal values to prove the AI fill resets them explicitly.
	stale := structs.GameplanLineup{FSShotVolume: 1, SSShotVolume: 5, TSShotVolume: 2}
	collegeSlots := []*structs.CollegeLineup{{GameplanLineup: stale}, {GameplanLineup: stale}}
	collegeSlots[0].Position = "G"
	collegeSlots[1].Position = "C"
	fillCollegeLineupSlots(collegeSlots, []structs.CollegePlayer{})

	nbaSlots := []*structs.NBALineup{{GameplanLineup: stale}, {GameplanLineup: stale}}
	nbaSlots[0].Position = "F"
	nbaSlots[1].Position = "C"
	fillNBALineupSlots(nbaSlots, []structs.NBAPlayer{})

	for _, lineup := range []structs.GameplanLineup{collegeSlots[0].GameplanLineup, collegeSlots[1].GameplanLineup, nbaSlots[0].GameplanLineup, nbaSlots[1].GameplanLineup} {
		if lineup.FSShotVolume != structs.ShotVolumeNormal || lineup.SSShotVolume != structs.ShotVolumeNormal || lineup.TSShotVolume != structs.ShotVolumeNormal {
			t.Fatalf("expected AI-filled lineup to use Normal shot volume, got %d/%d/%d", lineup.FSShotVolume, lineup.SSShotVolume, lineup.TSShotVolume)
		}
	}
}
