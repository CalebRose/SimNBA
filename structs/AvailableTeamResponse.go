package structs

type AvailableTeamResponse struct {
	TeamName     string
	TeamAbbr     string
	Conference   string
	OverallGrade string
	OffenseGrade string
	DefenseGrade string
}

type TeamRecordResponse struct {
	OverallWins                int
	OverallLosses              int
	CurrentSeasonWins          int
	CurrentSeasonLosses        int
	TournamentWins             int
	TournamentLosses           int
	PlayoffWins                int
	PlayoffLosses              int
	NITWins                    int
	NITLosses                  int
	CBIWins                    int
	CBILosses                  int
	RegularSeasonChampionships []string
	ConferenceChampionships    []string
	SweetSixteens              []string
	EliteEights                []string
	FinalFours                 []string
	RunnerUps                  []string
	NationalChampionships      []string
	TopPlayers                 []TopPlayer
}

func (t *TeamRecordResponse) AddTopPlayers(players []TopPlayer) {
	t.TopPlayers = players
}

type TopPlayer struct {
	PlayerID     uint
	FirstName    string
	LastName     string
	Position     string
	Archetype    string
	OverallGrade string
	Overall      int
}

func (t *TopPlayer) MapCollegePlayer(player CollegePlayer, grade string) {
	t.PlayerID = player.ID
	t.FirstName = player.FirstName
	t.LastName = player.LastName
	t.Position = player.Position
	t.Archetype = player.Archetype
	t.Overall = player.Overall
	t.OverallGrade = grade
}

func (t *TopPlayer) MapNBAPlayer(player NBAPlayer) {
	t.PlayerID = player.ID
	t.FirstName = player.FirstName
	t.LastName = player.LastName
	t.Position = player.Position
	t.Archetype = player.Archetype
	t.Overall = player.Overall
}
