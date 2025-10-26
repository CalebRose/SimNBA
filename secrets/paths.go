package secrets

func GetPath() map[string]string {
	TeamsListPath := "./data/teamslist_v5.csv"
	NBAStandingsPath := "./NBAStandings.csv"
	NBATeamsPath := "./data/InternationalSuperleague.csv"
	ArenaPath := "./data/Arenas.csv"
	DraftLotteryPath := "./data/2025/DraftStandings.csv"
	aiBehaviorPath := "./data/NewAIBehaviors.csv"
	extensionsPath := "./data/TempExtensions.csv"
	cbbMatchPath := "./data/2025/2025_SimCBB_Regular_Season.csv"
	cbbConfTournamentPath := "./data/2025/2025_SimCBB_Conf_Tourneys.csv"
	cbbPostSeasonPath := "./data/2025/2025_SimCBB_Post_Season.csv"
	nbaMatchPath := "./data/2024/2024_SimNBA_Play_In.csv"
	nbaSeriesPath := "./data/2024/2024_SimNBA_Series.csv"
	draftPickPath := "./data/draft_picks.csv"
	return map[string]string{
		"teams":                    TeamsListPath,
		"nbastandings":             NBAStandingsPath,
		"nbateams":                 NBATeamsPath,
		"arenas":                   ArenaPath,
		"draftlottery":             DraftLotteryPath,
		"ai":                       aiBehaviorPath,
		"extensions":               extensionsPath,
		"cbbmatches":               cbbMatchPath,
		"nbamatches":               nbaMatchPath,
		"nbaseries":                nbaSeriesPath,
		"draftpicks":               draftPickPath,
		"cbbconftournamentmatches": cbbConfTournamentPath,
		"cbbpostseasonmatches":     cbbPostSeasonPath,
	}
}
