package util

// EventIDMap maps each play-by-play event ID to its enum name.
var EventIDMap = map[uint8]string{
	Tipoff:            "Tipoff",
	Ot_tipoff:         "Ot_tipoff",
	Steal:             "Steal",
	Turnover:          "Turnover",
	Move:              "Move",
	Pass_ball:         "Pass_ball",
	Heave:             "Heave",
	Free_throw:        "Free_throw",
	Shot_three:        "Shot_three",
	Shot_corner_three: "Shot_corner_three",
	Shot_inside:       "Shot_inside",
	Shot_paint:        "Shot_paint",
	Shot_midrange:     "Shot_midrange",
	QuarterOver:       "QuarterOver",
	HalfOver:          "HalfOver",
	GameOver:          "GameOver",
	OvertimeStart:     "OvertimeStart",
	OvertimeOver:      "OvertimeOver",
	Timeout:           "Timeout",
	Rebound:           "Rebound",
	Inbound:           "Inbound",
	Substitution:      "Substitution",
}

func ReturnStringFromEventIDPBPID(id uint8) string {
	if val, exists := EventIDMap[id]; exists {
		return val
	}
	return "Unknown Event ID"
}

// OutcomeIDMap maps each play-by-play outcome ID to its enum name.
var OutcomeIDMap = map[uint8]string{
	No_outcome:             "No_outcome",
	TipoffHomeWin:          "TipoffHomeWin",
	TipoffAwayWin:          "TipoffAwayWin",
	Heave_made:             "Heave_made",
	Heave_missed:           "Heave_missed",
	Shot_made:              "Shot_made",
	Shot_foul_made:         "Shot_foul_made",
	Shot_missed:            "Shot_missed",
	Shot_foul_missed:       "Shot_foul_missed",
	Shot_blocked:           "Shot_blocked",
	Shot_foul_blocked:      "Shot_foul_blocked",
	Shot_clock_violation:   "Shot_clock_violation",
	Move_success:           "Move_success",
	Move_foul:              "Move_foul",
	Pass_success:           "Pass_success",
	Pass_deflected:         "Pass_deflected",
	Pass_intercepted:       "Pass_intercepted",
	Pass_foul:              "Pass_foul",
	Move_cutoff:            "Move_cutoff",
	Offensive_charge:       "Offensive_charge",
	Move_trapped:           "Move_trapped",
	Steal_success:          "Steal_success",
	Out_of_bounds_turnover: "Out_of_bounds_turnover",
	No_passing_lane:        "No_passing_lane",
	Ft_made:                "Ft_made",
	Ft_missed:              "Ft_missed",
	Offensive_rebound:      "Offensive_rebound",
	Defensive_rebound:      "Defensive_rebound",
	Inbound_success:        "Inbound_success",
	MediaTimeout:           "MediaTimeout",
	TeamTimeout:            "TeamTimeout",
}

func ReturnStringFromOutcomeIDPBPID(id uint8) string {
	if val, exists := OutcomeIDMap[id]; exists {
		return val
	}
	return "Unknown Outcome ID"
}
