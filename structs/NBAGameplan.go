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
	PreserveTimeouts     bool
	Trigger1Enabled      bool
	Trigger1Type         uint8 // 1 == Designated Player, 2 == Fouls Per Half
	Trigger1Value        uint  // Could either be player ID or number of fouls per half
	Trigger2Enabled      bool
	Trigger2Value        uint // Number of points the opponent is up by
	Trigger3Enabled      bool // Designate Player Exhaustion Trigger
	Trigger3Value        uint // PlayerID of the player to monitor for exhaustion
	Trigger3Exhaustion   uint // Exhaustion threshold for the designated player
	Trigger4Enabled      bool // On-floor average exhaustion
	Trigger4Value        uint // Average exhaustion of all players on the floor
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
