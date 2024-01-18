package secrets

func GetPath() map[string]string {
	TeamsListPath := "./data/teamslist_v5.csv"
	NBAStandingsPath := "./NBAStandings.csv"
	NBATeamsPath := "./NBATeams.csv"
	ArenaPath := "./data/Arenas.csv"
	DraftLotteryPath := "./data/DraftStandings.csv"
	aiBehaviorPath := "./data/NewAIBehaviors.csv"
	extensionsPath := "./data/TempExtensions.csv"
	cbbMatchPath := "./data/2023_CBB_Conf_Tourneys.csv"
	nbaMatchPath := "./data/2023_SimNBA_Season.csv"
	draftPickPath := "./data/draft_picks.csv"
	return map[string]string{
		"teams":        TeamsListPath,
		"nbastandings": NBAStandingsPath,
		"nbateams":     NBATeamsPath,
		"arenas":       ArenaPath,
		"draftlottery": DraftLotteryPath,
		"ai":           aiBehaviorPath,
		"extensions":   extensionsPath,
		"cbbmatches":   cbbMatchPath,
		"nbamatches":   nbaMatchPath,
		"draftpicks":   draftPickPath,
	}
}
