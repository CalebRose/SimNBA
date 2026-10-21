package structs

import (
	"testing"

	"gorm.io/gorm"
)

func TestGetCBBStreamChannel(t *testing.T) {
	g := func(id uint) Match { return Match{Model: gorm.Model{ID: id}} }
	cases := []struct {
		name       string
		game       Match
		home, away bool
		want       string
	}{
		{"no user-coached team", g(100), false, false, ""},
		{"home user-coached mod 0", g(100), true, false, "cbs"},
		{"away user-coached mod 1", g(101), false, true, "tbs"},
		{"both user-coached mod 2", g(102), true, true, "espn"},
		{"user-coached mod 3", g(103), true, false, "ytv"},
	}
	for _, c := range cases {
		if got := GetCBBStreamChannel(c.game, c.home, c.away); got != c.want {
			t.Errorf("%s: got %q want %q", c.name, got, c.want)
		}
	}
}
