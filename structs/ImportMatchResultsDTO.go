package structs

type ImportMatchResultsDTO struct {
	CBBResults []MatchResultsDTO
	NBAResults []MatchResultsDTO
}

type MatchResultsDTO struct {
	GameID    string
	TeamOne   TeamResultsDTO
	TeamTwo   TeamResultsDTO
	RosterOne []PlayerDTO
	RosterTwo []PlayerDTO
}

type TeamResultsDTO struct {
	TeamName   string
	Mascot     string
	Abbr       string
	Conference string
	Coach      string
	ID         int
	Stats      TeamStatsDTO
}

type TeamStatsDTO struct {
	Points             int
	Possessions        int
	FGM                int
	FGA                int
	FGPercent          float64
	ThreePointsMade    int
	ThreePointAttempts int
	ThreePointPercent  float64
	FTM                int
	FTA                int
	FTPercent          float64
	Rebounds           int
	OffRebounds        int
	DefRebounds        int
	Assists            int
	Steals             int
	Blocks             int
	TotalTurnovers     int
	LargestLead        int
	FirstHalfScore     int
	SecondQuarterScore int
	SecondHalfScore    int
	FourthQuarterScore int
	OvertimeScore      int
	Fouls              int
}

type PlayerDTO struct {
	ID            int
	FirstName     string
	LastName      string
	TeamID        int
	TeamAbbr      string
	IsRedshirt    bool
	IsRedshirting bool
	Position      string
	Age           int
	Stars         int
	Height        string
	Shooting2     int
	Shooting3     int
	Finishing     int
	Ballwork      int
	Rebounding    int
	Defense       int
	Stamina       int
	Minutes       int
	Overall       int
	Usage         float64
	AdjShooting   float64
	AdjFinishing  float64
	AdjBallwork   float64
	AdjRebounding float64
	AdjDefense    float64
	ReboundingPer float64
	DefensePer    float64
	AssistPer     float64
	Stats         PlayerStatsDTO
}

type PlayerStatsDTO struct {
	PlayerID           int
	Year               int
	MatchType          string
	Minutes            int
	Possessions        int
	FGM                int
	FGA                int
	FGPercent          float64
	ThreePointsMade    int
	ThreePointAttempts int
	ThreePointPercent  float64
	FTM                int
	FTA                int
	FTPercent          float64
	Points             int
	TotalRebounds      int
	OffRebounds        int
	DefRebounds        int
	Assists            int
	Steals             int
	Blocks             int
	Turnovers          int
	Fouls              int
	FouledOut          bool
	IsInjured          bool
	InjuryName         string
	InjuryType         string
	WeeksOfRecovery    uint8
}
