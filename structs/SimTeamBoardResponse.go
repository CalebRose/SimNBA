package structs

type SimTeamBoardResponse struct {
	ID                    uint
	TeamID                uint
	Team                  string
	TeamAbbreviation      string
	State                 string
	Region                string
	ScholarshipsAvailable int
	WeeklyPoints          int
	SpentPoints           int
	BonusPoints           int
	TotalCommitments      int
	RecruitClassSize      int
	ESPNScore             float64
	RivalsScore           float64
	Rank247Score          float64
	CompositeScore        float64
	RecruitingClassRank   int
	Recruits              []CrootProfile
}

func (stbr *SimTeamBoardResponse) Map(rtp TeamRecruitingProfile, c []CrootProfile) {
	stbr.ID = rtp.ID
	stbr.TeamID = rtp.TeamID
	stbr.TeamAbbreviation = rtp.TeamAbbr
	stbr.State = rtp.State
	stbr.Region = rtp.Region
	stbr.ScholarshipsAvailable = rtp.ScholarshipsAvailable
	stbr.WeeklyPoints = rtp.WeeklyPoints
	stbr.SpentPoints = rtp.SpentPoints
	stbr.TotalCommitments = rtp.TotalCommitments
	stbr.ESPNScore = rtp.ESPNScore
	stbr.RivalsScore = rtp.RivalsScore
	stbr.Rank247Score = rtp.Rank247Score
	stbr.CompositeScore = rtp.CompositeScore
	stbr.Recruits = c
	stbr.RecruitClassSize = rtp.RecruitClassSize
}

type CBBRosterResponse struct {
	Players  []CollegePlayerResponse
	Promises []CollegePromise
}

type DashboardResponseData struct {
	CollegeStandings []CollegeStandings
	NewsLogs         []NewsLog
	CollegeGames     []Match
	NBAStandings     []NBAStandings
	NBAGames         []NBAMatch
	TopCBBPlayers    []CollegePlayer
	TopNBAPlayers    []NBAPlayer
	TopTenPoll       CollegePollOfficial
	CBBTeamStats     TeamSeasonStats
	NBATeamStats     NBATeamSeasonStats
}
