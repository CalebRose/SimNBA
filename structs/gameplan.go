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

func (gl *GameplanLineup) MapLineupData(updated GameplanLineup) {
	gl.FirstStringID = updated.FirstStringID
	gl.FSMinutes = updated.FSMinutes
	gl.FSInsideProportion = updated.FSInsideProportion
	gl.FSMidProportion = updated.FSMidProportion
	gl.FSThreeProportion = updated.FSThreeProportion
	gl.SecondStringID = updated.SecondStringID
	gl.SSMinutes = updated.SSMinutes
	gl.SSInsideProportion = updated.SSInsideProportion
	gl.SSMidProportion = updated.SSMidProportion
	gl.SSThreeProportion = updated.SSThreeProportion
	gl.ThirdStringID = updated.ThirdStringID
	gl.TSMinutes = updated.TSMinutes
	gl.TSInsideProportion = updated.TSInsideProportion
	gl.TSMidProportion = updated.TSMidProportion
	gl.TSThreeProportion = updated.TSThreeProportion
}

// ClearPlayer removes a player from every string in this lineup and clears the
// associated minutes and shot allocations. It reports whether the lineup
// contained the player.
func (gl *GameplanLineup) ClearPlayer(playerID uint) bool {
	cleared := false
	if gl.FirstStringID == playerID {
		gl.FirstStringID = 0
		gl.FSMinutes = 0
		gl.FSInsideProportion = 0
		gl.FSMidProportion = 0
		gl.FSThreeProportion = 0
		cleared = true
	}
	if gl.SecondStringID == playerID {
		gl.SecondStringID = 0
		gl.SSMinutes = 0
		gl.SSInsideProportion = 0
		gl.SSMidProportion = 0
		gl.SSThreeProportion = 0
		cleared = true
	}
	if gl.ThirdStringID == playerID {
		gl.ThirdStringID = 0
		gl.TSMinutes = 0
		gl.TSInsideProportion = 0
		gl.TSMidProportion = 0
		gl.TSThreeProportion = 0
		cleared = true
	}
	return cleared
}

type CollegeLineup struct {
	GameplanLineup
}

type NBALineup struct {
	GameplanLineup
}
