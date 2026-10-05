package util

const (
	// EventIDs
	Tipoff            uint8 = 1
	Ot_tipoff         uint8 = 2
	Steal             uint8 = 3
	Turnover          uint8 = 4
	Move              uint8 = 5
	Pass_ball         uint8 = 6
	Heave             uint8 = 7
	Free_throw        uint8 = 8
	Shot_three        uint8 = 9
	Shot_corner_three uint8 = 10
	Shot_inside       uint8 = 11
	Shot_paint        uint8 = 12
	Shot_midrange     uint8 = 13
	QuarterOver       uint8 = 14
	HalfOver          uint8 = 15
	GameOver          uint8 = 16
	OvertimeStart     uint8 = 17
	OvertimeOver      uint8 = 18
	Timeout           uint8 = 19
	Rebound           uint8 = 20
	Inbound           uint8 = 21

	// OutcomeIDs
	No_outcome             uint8 = 0
	TipoffHomeWin          uint8 = 1
	TipoffAwayWin          uint8 = 2
	Heave_made             uint8 = 3
	Heave_missed           uint8 = 4
	Shot_made              uint8 = 5
	Shot_foul_made         uint8 = 6
	Shot_missed            uint8 = 7
	Shot_foul_missed       uint8 = 8
	Shot_blocked           uint8 = 9
	Shot_foul_blocked      uint8 = 10
	Shot_clock_violation   uint8 = 11
	Move_success           uint8 = 12
	Move_foul              uint8 = 13
	Pass_success           uint8 = 14
	Pass_deflected         uint8 = 15
	Pass_intercepted       uint8 = 16
	Pass_foul              uint8 = 17
	Move_cutoff            uint8 = 18
	Offensive_charge       uint8 = 19
	Move_trapped           uint8 = 20
	Steal_success          uint8 = 21
	Out_of_bounds_turnover uint8 = 22
	No_passing_lane        uint8 = 23
	Ft_made                uint8 = 24
	Ft_missed              uint8 = 25
	Offensive_rebound      uint8 = 26
	Defensive_rebound      uint8 = 27
	Inbound_success        uint8 = 28
)
