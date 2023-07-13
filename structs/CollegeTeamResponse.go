package structs

type CollegeTeamResponse struct {
	ID           uint
	Team         string
	Nickname     string
	Abbr         string
	ConferenceID uint
	Conference   string
	Coach        string
	OverallGrade string
	OffenseGrade string
	DefenseGrade string
	IsNBA        bool
	IsActive     bool
	SeasonStats  TeamSeasonStatsResponse
	Stats        TeamStats
}

type NBATeamResponse struct {
	ID              uint
	Team            string
	Nickname        string
	Abbr            string
	LeagueID        string
	League          string
	ConferenceID    uint
	Conference      string
	DivisionID      uint
	Division        string
	Coach           string
	OverallGrade    string
	OffenseGrade    string
	DefenseGrade    string
	IsActive        bool
	IsInternational bool
	SeasonStats     TeamSeasonStatsResponse
	Stats           NBATeamStats
}
