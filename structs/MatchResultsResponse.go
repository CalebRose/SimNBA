package structs

type MatchResultsResponse struct {
	HomePlayers []MatchResultsPlayer
	AwayPlayers []MatchResultsPlayer
	HomeStats   MatchResultsTeam
	AwayStats   MatchResultsTeam
}

type MatchResultsTeam struct {
	Team               string
	FirstHalfScore     int
	SecondQuarterScore int
	SecondHalfScore    int
	FourthQuarterScore int
	OvertimeScore      int
	Points             int
	Possessions        int
}

type MatchResultsPlayer struct {
	TeamID             uint
	FirstName          string
	LastName           string
	Position           string
	Archetype          string
	League             string
	Year               uint
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
}

type ByPlayedMinutes []MatchResultsPlayer

func (c ByPlayedMinutes) Len() int           { return len(c) }
func (c ByPlayedMinutes) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c ByPlayedMinutes) Less(i, j int) bool { return c[i].Minutes > c[j].Minutes }
