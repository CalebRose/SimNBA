package structs

import "github.com/jinzhu/gorm"

type NBAGameplanResponse struct {
	Gameplan       NBAGameplan
	OpposingRoster []NBAPlayer
}

// NBAGameplan - A team's strategy for their weekly gameplan
type NBAGameplan struct {
	gorm.Model
	TeamID               uint
	Game                 string
	Pace                 string
	ThreePointProportion int
	JumperProportion     int
	PaintProportion      int
	FocusPlayer          string
	OffensiveFormation   string
	DefensiveFormation   string
	OffensiveStyle       string
	Toggle2pt            bool
	Toggle3pt            bool
	ToggleFT             bool
	ToggleFN             bool
	ToggleBW             bool
	ToggleRB             bool
	ToggleID             bool
	TogglePD             bool
	ToggleP2             bool
	ToggleP3             bool
}

func (g *NBAGameplan) UpdateGameplan(pace, of, df, os, fp string) {
	g.Pace = pace
	g.OffensiveFormation = of
	g.DefensiveFormation = df
	g.OffensiveStyle = os
	g.FocusPlayer = fp
}

func (g *NBAGameplan) UpdateToggles(tp, thp, fn, ft, bw, rb, id, pd, p2, p3 bool) {
	g.Toggle2pt = tp
	g.Toggle3pt = thp
	g.ToggleFN = fn
	g.ToggleFT = ft
	g.ToggleBW = bw
	g.ToggleRB = rb
	g.ToggleID = id
	g.TogglePD = pd
	g.ToggleP2 = p2
	g.ToggleP3 = p3
}
