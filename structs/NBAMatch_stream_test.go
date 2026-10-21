package structs

import (
	"testing"

	"gorm.io/gorm"
)

func TestGetNBAStreamChannel(t *testing.T) {
	g := func(id uint, playoff bool) NBAMatch {
		return NBAMatch{Model: gorm.Model{ID: id}, IsPlayoffGame: playoff}
	}
	cases := []struct {
		name       string
		game       NBAMatch
		home, away uint
		want       string
	}{
		{"excluded id", g(25533, false), 1, 2, ""},
		{"excluded id playoff", g(25537, true), 40, 41, ""},
		{"both non-nba regular", g(100, false), 40, 41, ""},
		{"both non-nba playoff", g(100, true), 40, 41, "int"},
		{"one non-nba playoff", g(100, true), 5, 41, "int"},
		{"one non-nba regular uses mod", g(101, false), 5, 41, "nbatv"},
		{"mod 0", g(100, false), 1, 2, "tnt"},
		{"mod 1", g(101, false), 1, 2, "nbatv"},
		{"mod 2", g(102, false), 1, 2, "prime"},
		{"mod 3", g(103, false), 1, 2, "hbo"},
		{"nba playoff uses mod", g(102, true), 1, 2, "prime"},
	}
	for _, c := range cases {
		if got := GetNBAStreamChannel(c.game, c.home, c.away); got != c.want {
			t.Errorf("%s: got %q want %q", c.name, got, c.want)
		}
	}
}
