package structs

import "github.com/jinzhu/gorm"

// Gameplan - A team's strategy for their weekly gameplan

type GameplanResponse struct {
	Gameplan       Gameplan
	OpposingRoster []CollegePlayer
}

type Gameplan struct {
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

func (g *Gameplan) UpdateGameplan(pace, of, df, os, fp string) {
	g.Pace = pace
	g.OffensiveFormation = of
	g.DefensiveFormation = df
	g.OffensiveStyle = os
	g.FocusPlayer = fp
}

func (g *Gameplan) UpdateToggles(tp, thp, fn, ft, bw, rb, id, pd, p2, p3 bool) {
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

// UpdatePace - Update the Pace of the Gameplan
func (g *Gameplan) UpdatePace(pace string) {
	g.Pace = pace
}

// Update3PtProportion
func (g *Gameplan) Update3PtProportion(ratio int) {
	g.ThreePointProportion = ratio
}

func (g *Gameplan) UpdateJumperProportion(ratio int) {
	g.JumperProportion = ratio
}

func (g *Gameplan) UpdatePaintProportion(ratio int) {
	g.PaintProportion = ratio
}

type GameplanLineup struct {
	gorm.Model         // Just ignore this, it's for GORM (primary ID).
	TeamID             uint
	Position           string // G, F, or C
	FirstStringID      uint   // PlayerID at first string
	FSMinutes          uint8
	FSInsideProportion uint8 // Proportion towards shooting inside shots
	FSMidProportion    uint8 // Proportion towards shooting midrange shots
	FSThreeProportion  uint8 // Proportion towards shooting three point shots
	SecondStringID     uint  // PlayerID at second string
	SSMinutes          uint8
	SSInsideProportion uint8
	SSMidProportion    uint8
	SSThreeProportion  uint8
	ThirdStringID      uint // PlayerID at third string
	TSMinutes          uint8
	TSInsideProportion uint8
	TSMidProportion    uint8
	TSThreeProportion  uint8
}

type CollegeLineup struct {
	GameplanLineup
}

type NBALineup struct {
	GameplanLineup
}
