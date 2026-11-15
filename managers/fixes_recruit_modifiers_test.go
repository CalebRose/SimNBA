package managers

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/CalebRose/SimNBA/structs"
	"gorm.io/gorm"
)

func TestFixExistingModifiersForRecruits(t *testing.T) {
	attributes := structs.ProfileAttributes{
		ProgramPrestige: 5, ProfessionalPrestige: 5, Traditions: 5,
		Facilities: 5, Atmosphere: 5, Academics: 5, CampusLife: 5,
		ConferencePrestige: 5, CoachRating: 5, SeasonMomentum: 5, MediaSpotlight: 5,
	}
	recruits := []structs.Recruit{
		{Model: gorm.Model{ID: 11}, BasePlayer: structs.BasePlayer{State: "CA", Country: "USA"}},
		{Model: gorm.Model{ID: 22}, BasePlayer: structs.BasePlayer{
			State: "CA", Country: "Canada", PlayerPreferences: structs.PlayerPreferences{ProgramPref: 9},
		}},
	}
	teams := []structs.TeamRecruitingProfile{
		{Model: gorm.Model{ID: 501}, TeamID: 1, State: "California", Region: "Pacific", ProfileAttributes: attributes},
		{Model: gorm.Model{ID: 502}, TeamID: 2, State: "California", Region: "Midwest", ProfileAttributes: attributes},
		{Model: gorm.Model{ID: 503}, TeamID: 3, State: "Nevada", Region: "Pacific", ProfileAttributes: attributes},
	}
	teams[2].ProgramPrestige = 9

	profiles := []structs.RecruitPlayerProfile{
		{Model: gorm.Model{ID: 101}, RecruitID: 11, ProfileID: 1, Modifier: 9, TotalPoints: 999, PreviouslySpentPoints: 10, CurrentWeeksPoints: 3, Scholarship: true},
		{Model: gorm.Model{ID: 102}, RecruitID: 11, ProfileID: 2, Modifier: 9, TotalPoints: 999, PreviouslySpentPoints: 7},
		{Model: gorm.Model{ID: 103}, RecruitID: 22, ProfileID: 3, Modifier: 9, TotalPoints: 999, PreviouslySpentPoints: 13},
		{Model: gorm.Model{ID: 104}, RecruitID: 11, ProfileID: 3, Modifier: 9, TotalPoints: 999},
	}
	var saved []structs.RecruitPlayerProfile
	fixExistingModifiersForRecruits(recruits, profiles, teams, func(profile structs.RecruitPlayerProfile) {
		saved = append(saved, profile)
	})

	if len(saved) != len(profiles) {
		t.Fatalf("saved %d profiles, want %d", len(saved), len(profiles))
	}
	tests := []struct {
		name     string
		modifier float64
		points   float64
	}{
		{"region takes precedence over state", (11.0*14.0/9.0)/10.5 + 0.1, ((11.0*14.0/9.0)/10.5 + 0.1) * 10},
		{"state bonus without region", (11.0*14.0/9.0)/10.5 + 0.2, ((11.0*14.0/9.0)/10.5 + 0.2) * 7},
		{"international recruit and team attributes", (10.0*14.0/9.0 + 13.0/9.0) / 10.5, (10.0*14.0/9.0 + 13.0/9.0) / 10.5 * 13},
		{"zero previously spent points", (10.0*14.0/9.0+26.0/9.0)/10.5 + 0.1, 0},
	}
	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := saved[i]
			if math.Abs(float64(got.Modifier)-tc.modifier) > 0.00001 {
				t.Errorf("modifier = %v, want approximately %v", got.Modifier, tc.modifier)
			}
			if math.Abs(got.TotalPoints-tc.points) > 0.00001 {
				t.Errorf("total points = %v, want approximately %v", got.TotalPoints, tc.points)
			}
			want := profiles[i]
			want.Modifier = got.Modifier
			want.TotalPoints = got.TotalPoints
			if !reflect.DeepEqual(got, want) {
				t.Errorf("unrelated profile fields changed: got %+v, want %+v", got, want)
			}
		})
	}
}

func TestFixExistingModifiersForRecruitsNoProfiles(t *testing.T) {
	fixExistingModifiersForRecruits(nil, nil, nil, func(structs.RecruitPlayerProfile) {
		t.Error("unexpected save with no profiles")
	})
}

func TestCalculateModifierTowardsRecruitAttributesAndPreferences(t *testing.T) {
	tests := []struct {
		name          string
		setPreference func(*structs.PlayerPreferences, uint8)
		setAttribute  func(*structs.ProfileAttributes, uint8)
	}{
		{"program", func(p *structs.PlayerPreferences, v uint8) { p.ProgramPref = v }, func(a *structs.ProfileAttributes, v uint8) { a.ProgramPrestige = v }},
		{"professional development", func(p *structs.PlayerPreferences, v uint8) { p.ProfDevPref = v }, func(a *structs.ProfileAttributes, v uint8) { a.ProfessionalPrestige = v }},
		{"traditions", func(p *structs.PlayerPreferences, v uint8) { p.TraditionsPref = v }, func(a *structs.ProfileAttributes, v uint8) { a.Traditions = v }},
		{"facilities", func(p *structs.PlayerPreferences, v uint8) { p.FacilitiesPref = v }, func(a *structs.ProfileAttributes, v uint8) { a.Facilities = v }},
		{"atmosphere", func(p *structs.PlayerPreferences, v uint8) { p.AtmospherePref = v }, func(a *structs.ProfileAttributes, v uint8) { a.Atmosphere = v }},
		{"academics", func(p *structs.PlayerPreferences, v uint8) { p.AcademicsPref = v }, func(a *structs.ProfileAttributes, v uint8) { a.Academics = v }},
		{"conference", func(p *structs.PlayerPreferences, v uint8) { p.ConferencePref = v }, func(a *structs.ProfileAttributes, v uint8) { a.ConferencePrestige = v }},
		{"coach", func(p *structs.PlayerPreferences, v uint8) { p.CoachPref = v }, func(a *structs.ProfileAttributes, v uint8) { a.CoachRating = v }},
		{"season momentum", func(p *structs.PlayerPreferences, v uint8) { p.SeasonMomentumPref = v }, func(a *structs.ProfileAttributes, v uint8) { a.SeasonMomentum = v }},
		{"campus life", func(p *structs.PlayerPreferences, v uint8) { p.CampusLifePref = v }, func(a *structs.ProfileAttributes, v uint8) { a.CampusLife = v }},
		{"media", func(p *structs.PlayerPreferences, v uint8) { p.MediaPref = v }, func(a *structs.ProfileAttributes, v uint8) { a.MediaSpotlight = v }},
	}
	prefs := structs.PlayerPreferences{
		ProgramPref: 5, ProfDevPref: 5, TraditionsPref: 5, FacilitiesPref: 5,
		AtmospherePref: 5, AcademicsPref: 5, ConferencePref: 5, CoachPref: 5,
		SeasonMomentumPref: 5, CampusLifePref: 5, MediaPref: 5,
	}
	team := structs.TeamRecruitingProfile{ProfileAttributes: structs.ProfileAttributes{
		ProgramPrestige: 5, ProfessionalPrestige: 5, Traditions: 5, Facilities: 5,
		Atmosphere: 5, Academics: 5, ConferencePrestige: 5, CoachRating: 5,
		SeasonMomentum: 5, CampusLife: 5, MediaSpotlight: 5,
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for _, value := range []uint8{2, 8} {
				t.Run(fmt.Sprintf("preference_%d", value), func(t *testing.T) {
					changed := prefs
					tc.setPreference(&changed, value)
					got := CalculateModifierTowardsRecruit(changed, team, false, false)
					want := (11.0 + (5.0-float64(value))/9.0) / 10.5
					if math.Abs(float64(got)-want) > 0.00001 {
						t.Errorf("modifier = %v, want approximately %v", got, want)
					}
				})
				t.Run(fmt.Sprintf("team_attribute_%d", value), func(t *testing.T) {
					changed := team
					tc.setAttribute(&changed.ProfileAttributes, value)
					got := CalculateModifierTowardsRecruit(prefs, changed, false, false)
					factor := 1.0 + (float64(value)-5.0)/9.0
					want := (10.0 + factor*factor) / 10.5
					if math.Abs(float64(got)-want) > 0.00001 {
						t.Errorf("modifier = %v, want approximately %v", got, want)
					}
				})
			}
		})
	}
}
