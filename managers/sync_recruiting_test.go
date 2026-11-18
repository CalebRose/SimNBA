package managers

import (
	"sort"
	"testing"

	"github.com/CalebRose/SimNBA/structs"
	"gorm.io/gorm"
)

func TestAllocateRecruitProfilePoints(t *testing.T) {
	profiles := []structs.RecruitPlayerProfile{
		{Model: gorm.Model{ID: 1}, RecruitID: 10, ProfileID: 1, TeamAbbreviation: "A", TotalPoints: 10, CurrentWeeksPoints: 10, Modifier: 2},
		{Model: gorm.Model{ID: 2}, RecruitID: 10, ProfileID: 2, TeamAbbreviation: "B", TotalPoints: 30, Modifier: 1},
		{Model: gorm.Model{ID: 3}, RecruitID: 10, ProfileID: 3, TeamAbbreviation: "C", TotalPoints: 60, CurrentWeeksPoints: 21, Modifier: 1},
	}
	teamPoints := make(map[string]float64)
	var allocations []structs.RecruitPointAllocation
	if !allocateRecruitProfilePoints(profiles, 7, teamPoints, func(allocation structs.RecruitPointAllocation) {
		allocations = append(allocations, allocation)
	}) {
		t.Fatal("expected points to be placed")
	}

	if profiles[0].TotalPoints != 50 || profiles[0].CurrentWeeksPoints != 0 || profiles[0].PreviouslySpentPoints != 10 || profiles[0].SpendingCount != 1 {
		t.Errorf("allocated profile was not updated: %+v", profiles[0])
	}
	if profiles[1].TotalPoints != 30 || profiles[1].CurrentWeeksPoints != 0 {
		t.Errorf("profile without weekly points changed: %+v", profiles[1])
	}
	if profiles[2].TotalPoints != 0 || profiles[2].CurrentWeeksPoints != 0 || profiles[2].PreviouslySpentPoints != 21 {
		t.Errorf("over-limit profile was not reset: %+v", profiles[2])
	}
	if teamPoints["A"] != 10 || teamPoints["C"] != 21 {
		t.Errorf("team points = %v, want A=10 and C=21", teamPoints)
	}
	if len(allocations) != 2 {
		t.Fatalf("got %d allocation records, want 2", len(allocations))
	}
	if allocations[0].RecruitProfileID != 1 || allocations[0].WeekID != 7 || allocations[0].Points != 10 || allocations[0].BonusAffectedPoints != 40 || allocations[0].CaughtCheating {
		t.Errorf("normal allocation = %+v", allocations[0])
	}
	if allocations[1].RecruitProfileID != 3 || allocations[1].Points != 21 || allocations[1].BonusAffectedPoints != 0 || !allocations[1].CaughtCheating {
		t.Errorf("over-limit allocation = %+v", allocations[1])
	}

	sort.Sort(structs.ByPoints(profiles))
	if profiles[0].ID != 1 || profiles[1].ID != 2 || profiles[2].ID != 3 {
		t.Errorf("profiles sorted using incorrect totals: %+v", profiles)
	}
}

func TestAllocateRecruitProfilePointsWithoutSpending(t *testing.T) {
	profiles := []structs.RecruitPlayerProfile{{TotalPoints: 20}}
	if allocateRecruitProfilePoints(profiles, 1, make(map[string]float64), func(structs.RecruitPointAllocation) {
		t.Error("unexpected allocation record")
	}) {
		t.Error("expected no points to be placed")
	}
	if profiles[0].TotalPoints != 20 {
		t.Errorf("total points = %v, want 20", profiles[0].TotalPoints)
	}
}
