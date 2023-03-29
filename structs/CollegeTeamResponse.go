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
}

type NBATeamResponse struct {
	ID           uint
	Team         string
	Nickname     string
	Abbr         string
	ConferenceID uint
	Conference   string
	DivisionID   uint
	Division     string
	Coach        string
	OverallGrade string
	OffenseGrade string
	DefenseGrade string
	IsActive     bool
	SeasonStats  TeamSeasonStatsResponse
}
