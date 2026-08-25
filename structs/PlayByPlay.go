package structs

import "gorm.io/gorm"

type PlayByPlayDTO struct {
	GameID              uint
	Quarter             uint8
	TimeOnClock         uint16
	ShotClock           uint16
	SecondsConsumed     uint16
	HomeTeamScore       uint16
	AwayTeamScore       uint16
	TeamID              uint
	BallCarrierID       uint
	AssistingPlayerID   uint
	PassedPlayerID      uint
	DefenderID          uint
	BlockingPlayerID    uint
	StealingPlayerID    uint
	FoulingPlayerID     uint
	SubstitutePlayerID  uint
	InjuryID            uint8
	InjuryType          uint8
	InjuryDuration      uint8
	PenaltyID           uint8
	HomeOffensiveSystem string
	HomeDefensiveSystem string
	AwayOffensiveSystem string
	AwayDefensiveSystem string
	EventID             uint8
	OutcomeID           uint8
	XAxis               int8
	YAxis               int8
	NextXAxis           int8
	NextYAxis           int8
}

type BasePlayByPlay struct {
	GameID              uint
	Quarter             uint8
	TimeOnClock         uint16
	ShotClock           uint16
	SecondsConsumed     uint16
	HomeTeamScore       uint16
	AwayTeamScore       uint16
	TeamID              uint
	BallCarrierID       uint
	AssistingPlayerID   uint
	PassedPlayerID      uint
	DefenderID          uint
	BlockingPlayerID    uint
	StealingPlayerID    uint
	FoulingPlayerID     uint
	SubstitutePlayerID  uint
	InjuryID            uint8
	InjuryType          uint8
	InjuryDuration      uint8
	PenaltyID           uint8
	HomeOffensiveSystem uint8
	HomeDefensiveSystem uint8
	AwayOffensiveSystem uint8
	AwayDefensiveSystem uint8
	EventID             uint8
	OutcomeID           uint8
	XAxis               int8
	YAxis               int8
	NextXAxis           int8
	NextYAxis           int8
}

type CollegePlayByPlay struct {
	gorm.Model
	BasePlayByPlay
}

type NBAPlayByPlay struct {
	gorm.Model
	BasePlayByPlay
}

type PlayByPlayResponse struct {
	GameID              uint
	Quarter             uint8
	TimeOnClock         string
	ShotClock           uint16
	SecondsConsumed     uint16
	HomeTeamScore       uint16
	AwayTeamScore       uint16
	TeamID              uint
	BallCarrierID       uint
	AssistingPlayerID   uint
	PassedPlayerID      uint
	DefenderID          uint
	BlockingPlayerID    uint
	StealingPlayerID    uint
	FoulingPlayerID     uint
	SubstitutePlayerID  uint
	InjuryID            uint8
	InjuryType          uint8
	InjuryDuration      uint8
	PenaltyID           uint8
	HomeOffensiveSystem uint8
	HomeDefensiveSystem uint8
	AwayOffensiveSystem uint8
	AwayDefensiveSystem uint8
	Event               string
	Outcome             string
	XAxis               int8
	YAxis               int8
	NextXAxis           int8
	NextYAxis           int8
	Result              string
	StreamResult        []string
}

func (p *PlayByPlayResponse) AddPlayInformation(toc, event, outcome string, xAxis, yAxis, nextXAxis, nextYAxis int8) {
	p.TimeOnClock = toc
	p.Event = event
	p.XAxis = xAxis
	p.YAxis = yAxis
	p.NextXAxis = nextXAxis
	p.NextYAxis = nextYAxis
	p.Outcome = outcome
}

func (p *PlayByPlayResponse) AddResult(result []string, isStream bool) {
	if isStream {
		p.StreamResult = result
	} else {
		p.Result = result[0]
	}
}

type GameResultsPlayer struct {
	ID uint
	BasePlayer
}

type StreamResponse struct {
	GameID            uint
	HomeTeamID        uint
	GameLabel         string
	HomeLabel         string
	HomeTeam          string
	HomeTeamCoach     string
	HomeTeamDiscordID string
	HomeTeamRank      uint
	AwayTeamID        uint
	AwayLabel         string
	AwayTeam          string
	AwayTeamCoach     string
	AwayTeamDiscordID string
	AwayTeamRank      uint
	ArenaID           uint
	Arena             string
	Capacity          uint
	Attendance        uint
	City              string
	State             string
	Country           string
	HomeTeamWin       bool
	AwayTeamWin       bool
	HTScore           uint
	ATScore           uint
	Streams           []PlayByPlayResponse
}
