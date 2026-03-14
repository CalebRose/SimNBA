package secrets

func GetPath() map[string]string {
	TeamsListPath := "./data/2026/2026_teams_import_map.csv"
	NBAStandingsPath := "./NBAStandings.csv"
	NBATeamsPath := "./data/InternationalSuperleague.csv"
	ArenaPath := "./data/Arenas.csv"
	DraftLotteryPath := "./data/2025/DraftStandings.csv"
	aiBehaviorPath := "./data/NewAIBehaviors.csv"
	extensionsPath := "./data/TempExtensions.csv"
	cbbMatchPath := "./data/2026/2026_cbb_game_data_test.csv"
	cbbConfTournamentPath := "./data/2025/2025_SimCBB_Conf_Tourneys.csv"
	cbbPostSeasonPath := "./data/2025/2025_SimCBB_Post_Season.csv"
	nbaMatchPath := "./data/2025/2025_SimNBA_Play_In.csv"
	nbaSeriesPath := "./data/2025/2025_SimNBA_Series.csv"
	draftPickPath := "./data/draft_picks.csv"
	collegePlayersPath := "./data/2026/Migration/2026_cbb_players_table.csv"
	historicCollegePlayersPath := "./data/2026/Migration/2026_historic_cbb_players_table.csv"
	nbaDrafteesPath := "./data/2026/Migration/2026_nba_draftees_table.csv"
	nbaPlayersPath := "./data/2026/Migration/2026_nba_players_table.csv"
	nbaRetiredPlayersPath := "./data/2026/Migration/2026_nba_retirees_table.csv"
	return map[string]string{
		"teams":                         TeamsListPath,
		"nbastandings":                  NBAStandingsPath,
		"nbateams":                      NBATeamsPath,
		"arenas":                        ArenaPath,
		"draftlottery":                  DraftLotteryPath,
		"ai":                            aiBehaviorPath,
		"extensions":                    extensionsPath,
		"cbbmatches":                    cbbMatchPath,
		"nbamatches":                    nbaMatchPath,
		"nbaseries":                     nbaSeriesPath,
		"draftpicks":                    draftPickPath,
		"cbbconftournamentmatches":      cbbConfTournamentPath,
		"cbbpostseasonmatches":          cbbPostSeasonPath,
		"college_players_2026":          collegePlayersPath,
		"historic_college_players_2026": historicCollegePlayersPath,
		"nba_players_2026":              nbaPlayersPath,
		"nba_draftees_2026":             nbaDrafteesPath,
		"nba_retired_players_2026":      nbaRetiredPlayersPath,
	}
}
