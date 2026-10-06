package managers

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/CalebRose/SimNBA/repository"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
)

func GetCBBPlayByPlayStreamData(streamType string) []structs.StreamResponse {
	ts := GetTimestamp()
	teamMap := GetCollegeTeamMap()
	collegePlayers := GetAllCollegePlayers()
	rosterMap := MakeCollegePlayerMapByTeamID(collegePlayers, true)
	weekID := strconv.Itoa(int(ts.CollegeWeekID))
	seasonID := strconv.Itoa(int(ts.SeasonID))
	gameDay := ts.GetGameDay()
	games := GetMatchesByWeekIdAndMatchType(weekID, seasonID, gameDay)
	_, gameType := ts.GetCurrentGameType(true)

	streams := []structs.StreamResponse{}

	for _, game := range games {
		if !game.GameComplete || game.IsRevealed {
			continue
		}
		homeTeam := teamMap[uint(game.HomeTeamID)]
		awayTeam := teamMap[uint(game.AwayTeamID)]

		if streamType == "tnt" {
			if !homeTeam.IsUserCoached && !awayTeam.IsUserCoached {
				continue
			}
			mod := game.ID % 4
			if mod == 1 {
				continue
			}
		}
		if streamType == "nbatv" {
			if !homeTeam.IsUserCoached && !awayTeam.IsUserCoached {
				continue
			}
			mod := game.ID % 4
			if mod == 0 {
				continue
			}
		}

		if streamType == "prime" {
			if !homeTeam.IsUserCoached && !awayTeam.IsUserCoached {
				continue
			}
			mod := game.ID % 4
			if mod == 2 {
				continue
			}
		}

		if streamType == "hbo" {
			if !homeTeam.IsUserCoached && !awayTeam.IsUserCoached {
				continue
			}
			mod := game.ID % 4
			if mod == 3 {
				continue
			}
		}

		gameID := strconv.Itoa(int(game.ID))
		var wg sync.WaitGroup
		var (
			playByPlays []structs.CollegePlayByPlay
			homePlayers []structs.CollegePlayer
			awayPlayers []structs.CollegePlayer
			homeStats   []structs.CollegePlayerStats
			awayStats   []structs.CollegePlayerStats
		)
		homePlayers = rosterMap[game.HomeTeamID]
		awayPlayers = rosterMap[game.AwayTeamID]
		wg.Add(2)

		go func() {
			defer wg.Done()
			stats := repository.FindCollegePlayerGameStatsRecords(seasonID, weekID, gameType, gameID)
			for _, s := range stats {
				if s.TeamID == game.HomeTeamID {
					homeStats = append(homeStats, s)
				} else {
					awayStats = append(awayStats, s)
				}
			}
		}()

		go func() {
			defer wg.Done()
			playByPlays = repository.FindCBBPlayByPlaysRecordsByGameID(gameID)
		}()

		wg.Wait()

		totalList := []structs.CollegePlayer{}
		totalList = append(totalList, homePlayers...)
		totalList = append(totalList, awayPlayers...)

		participantMap := MakeCollegePlayerMap(totalList)
		playbyPlayResponse := GenerateCBBPlayByPlayResponse(playByPlays, teamMap, participantMap, true, game.HomeTeamID, game.AwayTeamID)

		stream := structs.StreamResponse{
			GameID:            game.ID,
			HomeTeamID:        uint(game.HomeTeamID),
			HomeTeam:          game.HomeTeam,
			HomeTeamCoach:     homeTeam.Coach,
			HomeTeamRank:      game.HomeTeamRank,
			HomeLabel:         homeTeam.Team + " " + homeTeam.Nickname,
			HomeTeamDiscordID: homeTeam.DiscordID,
			AwayTeamID:        uint(game.AwayTeamID),
			AwayTeam:          game.AwayTeam,
			AwayTeamCoach:     awayTeam.Coach,
			AwayTeamRank:      game.AwayTeamRank,
			AwayTeamDiscordID: awayTeam.DiscordID,
			AwayLabel:         awayTeam.Team + " " + awayTeam.Nickname,
			Streams:           playbyPlayResponse,
			City:              game.City,
			State:             game.State,
			Country:           game.Country,
			ArenaID:           game.ArenaID,
			Arena:             game.Arena,
			Attendance:        uint(game.AttendanceCount),
		}

		streams = append(streams, stream)
	}

	return streams
}

func GetNBAPlayByPlayStreamData(streamType string) []structs.StreamResponse {
	ts := GetTimestamp()
	weekID := strconv.Itoa(int(ts.NBAWeekID))
	seasonID := strconv.Itoa(int(ts.SeasonID))
	gameDay := ts.GetGameDay()
	games := GetNBAMatchesByWeekIdAndMatchType(weekID, seasonID, gameDay)
	nbaTeams := GetAllActiveNBATeams()
	teamMap := MakeNBATeamMap(nbaTeams)
	nbaPlayers := GetAllNBAPlayers()
	rosterMap := MakeNBAPlayerMapByTeamID(nbaPlayers, true)
	streams := []structs.StreamResponse{}

	nbaGameplans := GetAllNBAGameplans()
	gameplanMap := MakeNBAGameplanMap(nbaGameplans)

	for _, game := range games {
		if game.ID == 25533 {
			continue
		}
		if !game.GameComplete || game.IsRevealed {
			continue
		}

		homeTeam := teamMap[uint(game.HomeTeamID)]
		awayTeam := teamMap[uint(game.AwayTeamID)]
		if streamType == "tnt" {
			if homeTeam.ID > 32 && awayTeam.ID > 32 {
				continue
			}
			mod := game.ID % 4
			if mod != 0 {
				continue
			}
		}
		if streamType == "nbatv" {
			if homeTeam.ID > 32 && awayTeam.ID > 32 {
				continue
			}
			mod := game.ID % 4
			if mod != 1 {
				continue
			}
		}

		if streamType == "prime" {
			if homeTeam.ID > 32 && awayTeam.ID > 32 {
				continue
			}
			mod := game.ID % 4
			if mod != 2 {
				continue
			}
		}

		if streamType == "hbo" {
			if homeTeam.ID > 32 && awayTeam.ID > 32 {
				continue
			}
			mod := game.ID % 4
			if mod != 3 {
				continue
			}
		}

		if streamType == "int" {
			if homeTeam.ID <= 32 && awayTeam.ID <= 32 {
				continue
			}
			if !game.IsPlayoffGame {
				continue
			}
		}

		gameID := strconv.Itoa(int(game.ID))
		var wg sync.WaitGroup
		var (
			playByPlays []structs.NBAPlayByPlay
			homeStats   []structs.NBAPlayerStats
			awayStats   []structs.NBAPlayerStats
		)
		homePlayers := rosterMap[game.HomeTeamID]
		awayPlayers := rosterMap[game.AwayTeamID]
		homeGameplan := gameplanMap[uint(game.HomeTeamID)]
		awayGameplan := gameplanMap[uint(game.AwayTeamID)]

		wg.Add(2)

		go func() {
			defer wg.Done()
			stats := repository.FindProPlayerGameStatsRecords(seasonID, weekID, gameDay, gameID)
			for _, s := range stats {
				if s.TeamID == game.HomeTeamID {
					homeStats = append(homeStats, s)
				} else {
					awayStats = append(awayStats, s)
				}
			}
		}()

		go func() {
			defer wg.Done()
			playByPlays = repository.FindNBAPlayByPlaysRecordsByGameID(gameID)
		}()

		wg.Wait()

		totalList := []structs.NBAPlayer{}
		totalList = append(totalList, homePlayers...)
		totalList = append(totalList, awayPlayers...)

		participantMap := MakeNBAPlayerMap(totalList)
		playbyPlayResponse := GenerateNBAPlayByPlayResponse(playByPlays, teamMap, participantMap, true, game.HomeTeamID, game.AwayTeamID)

		stream := structs.StreamResponse{
			GameID:                 game.ID,
			HomeTeamID:             uint(game.HomeTeamID),
			HomeTeam:               game.HomeTeam,
			HomeTeamCoach:          game.HomeTeamCoach,
			HomeTeamDiscordID:      homeTeam.OwnerDiscordID,
			HomeLabel:              game.HomeTeam,
			HomePace:               homeGameplan.Pace,
			HomeOffensiveFormation: homeGameplan.OffensiveFormation,
			HomeDefensiveFormation: homeGameplan.DefensiveFormation,
			AwayTeamID:             uint(game.AwayTeamID),
			AwayTeam:               game.AwayTeam,
			AwayTeamCoach:          game.AwayTeamCoach,
			AwayLabel:              game.AwayTeam,
			AwayPace:               awayGameplan.Pace,
			AwayOffensiveFormation: awayGameplan.OffensiveFormation,
			AwayDefensiveFormation: awayGameplan.DefensiveFormation,
			AwayTeamDiscordID:      awayTeam.OwnerDiscordID,
			Streams:                playbyPlayResponse,
			City:                   game.City,
			State:                  game.State,
		}

		streams = append(streams, stream)
	}

	return streams
}

func GenerateCBBPlayByPlayResponse(playByPlays []structs.CollegePlayByPlay, teamMap map[uint]structs.Team, playerMap map[uint]structs.CollegePlayer, isStream bool, ht, at uint) []structs.PlayByPlayResponse {
	results := []structs.PlayByPlayResponse{}
	for idx, play := range playByPlays {
		if play.OutcomeID == util.No_outcome {
			continue
		}
		timeOnClock := FormatTimeToClock(play.TimeOnClock)
		event := util.ReturnStringFromEventIDPBPID(play.EventID)
		outcome := util.ReturnStringFromOutcomeIDPBPID(play.OutcomeID)
		possessingTeam := teamMap[uint(play.TeamID)]
		result := generateCollegeResultsString(play.BasePlayByPlay, playerMap, possessingTeam)

		res := structs.PlayByPlayResponse{
			GameID:              play.GameID,
			PlayNumber:          uint(idx) + 1,
			HomeTeamID:          ht,
			HomeTeamScore:       play.HomeTeamScore,
			AwayTeamID:          at,
			AwayTeamScore:       play.AwayTeamScore,
			Quarter:             play.Quarter,
			TimeOnClock:         timeOnClock,
			SecondsConsumed:     play.SecondsConsumed,
			Event:               event,
			Outcome:             outcome,
			XAxis:               play.XAxis,
			YAxis:               play.YAxis,
			NextXAxis:           play.NextXAxis,
			NextYAxis:           play.NextYAxis,
			TeamID:              play.TeamID,
			BallCarrierID:       play.BallCarrierID,
			PassedPlayerID:      play.PassedPlayerID,
			AssistingPlayerID:   play.AssistingPlayerID,
			DefenderID:          play.DefenderID,
			BlockingPlayerID:    play.BlockingPlayerID,
			StealingPlayerID:    play.StealingPlayerID,
			FoulingPlayerID:     play.FoulingPlayerID,
			SubstitutePlayerID:  play.SubstitutePlayerID,
			InjuryID:            play.InjuryID,
			InjuryType:          play.InjuryType,
			InjuryDuration:      play.InjuryDuration,
			PenaltyID:           play.PenaltyID,
			HomeOffensiveSystem: play.HomeOffensiveSystem,
			HomeDefensiveSystem: play.HomeDefensiveSystem,
			AwayOffensiveSystem: play.AwayOffensiveSystem,
			AwayDefensiveSystem: play.AwayDefensiveSystem,
			Result:              result,
		}

		results = append(results, res)
	}
	return results
}

func GenerateNBAPlayByPlayResponse(playByPlays []structs.NBAPlayByPlay, teamMap map[uint]structs.NBATeam, playerMap map[uint]structs.NBAPlayer, isStream bool, ht, at uint) []structs.PlayByPlayResponse {
	results := []structs.PlayByPlayResponse{}
	for idx, play := range playByPlays {
		if play.OutcomeID == util.No_outcome {
			continue
		}
		timeOnClock := FormatTimeToClock(play.TimeOnClock)
		event := util.ReturnStringFromEventIDPBPID(play.EventID)
		outcome := util.ReturnStringFromOutcomeIDPBPID(play.OutcomeID)
		possessingTeam := teamMap[uint(play.TeamID)]

		result := generateProResultsString(play.BasePlayByPlay, playerMap, possessingTeam)

		res := structs.PlayByPlayResponse{
			GameID:              play.GameID,
			PlayNumber:          uint(idx) + 1,
			HomeTeamID:          ht,
			HomeTeamScore:       play.HomeTeamScore,
			AwayTeamID:          at,
			AwayTeamScore:       play.AwayTeamScore,
			Quarter:             play.Quarter,
			TimeOnClock:         timeOnClock,
			SecondsConsumed:     play.SecondsConsumed,
			Event:               event,
			Outcome:             outcome,
			XAxis:               play.XAxis,
			YAxis:               play.YAxis,
			NextXAxis:           play.NextXAxis,
			NextYAxis:           play.NextYAxis,
			TeamID:              play.TeamID,
			BallCarrierID:       play.BallCarrierID,
			PassedPlayerID:      play.PassedPlayerID,
			AssistingPlayerID:   play.AssistingPlayerID,
			DefenderID:          play.DefenderID,
			BlockingPlayerID:    play.BlockingPlayerID,
			StealingPlayerID:    play.StealingPlayerID,
			FoulingPlayerID:     play.FoulingPlayerID,
			SubstitutePlayerID:  play.SubstitutePlayerID,
			InjuryID:            play.InjuryID,
			InjuryType:          play.InjuryType,
			InjuryDuration:      play.InjuryDuration,
			PenaltyID:           play.PenaltyID,
			HomeOffensiveSystem: play.HomeOffensiveSystem,
			HomeDefensiveSystem: play.HomeDefensiveSystem,
			AwayOffensiveSystem: play.AwayOffensiveSystem,
			AwayDefensiveSystem: play.AwayDefensiveSystem,
			Result:              result,
		}

		results = append(results, res)
	}
	return results
}

func FormatTimeToClock(timeInSeconds uint16) string {
	minutes := timeInSeconds / 60
	seconds := timeInSeconds % 60
	formatted := fmt.Sprintf("%02d:%02d", minutes, seconds)
	return formatted
}

func getPlayerLabel(player structs.BasePlayer) string {
	if len(player.FirstName) == 0 {
		return ""
	}
	return player.Team + " " + player.Position + " " + player.FirstName + " " + player.LastName
}

func generateCollegeResultsString(play structs.BasePlayByPlay, playerMap map[uint]structs.CollegePlayer, possessingTeam structs.Team) string {
	labelFor := func(id uint) string {
		player, ok := playerMap[id]
		if !ok {
			return ""
		}
		return getPlayerLabel(player.BasePlayer)
	}
	return generatePlayStatement(play, labelFor, possessingTeam.Team)
}

func generateProResultsString(play structs.BasePlayByPlay, playerMap map[uint]structs.NBAPlayer, possessingTeam structs.NBATeam) string {
	labelFor := func(id uint) string {
		player, ok := playerMap[id]
		if !ok {
			return ""
		}
		return getPlayerLabel(player.BasePlayer)
	}
	return generatePlayStatement(play, labelFor, possessingTeam.Team)
}

func labelOrDefault(label, fallback string) string {
	if label == "" {
		return fallback
	}
	return label
}

func shotDescription(eventID uint8) string {
	switch eventID {
	case util.Shot_three:
		return "a three-pointer"
	case util.Shot_corner_three:
		return "a corner three"
	case util.Shot_inside:
		return "an inside shot"
	case util.Shot_paint:
		return "a shot in the paint"
	case util.Shot_midrange:
		return "a mid-range jumper"
	}
	return "a shot"
}

// generatePlayStatement builds the play-by-play sentence for a single basketball play.
func generatePlayStatement(play structs.BasePlayByPlay, labelFor func(uint) string, teamLabel string) string {
	carrier := labelOrDefault(labelFor(play.BallCarrierID), "The ball handler")
	defender := labelOrDefault(labelFor(play.DefenderID), "the defender")
	receiver := labelOrDefault(labelFor(play.PassedPlayerID), "a teammate")
	assister := labelFor(play.AssistingPlayerID)
	blocker := labelOrDefault(labelFor(play.BlockingPlayerID), defender)
	stealer := labelOrDefault(labelFor(play.StealingPlayerID), defender)
	fouler := labelOrDefault(labelFor(play.FoulingPlayerID), defender)
	team := labelOrDefault(teamLabel, "the offense")

	statement := ""
	switch play.EventID {
	case util.Tipoff, util.Ot_tipoff:
		prefix := "Tip-off! "
		if play.EventID == util.Ot_tipoff {
			prefix = "Overtime tip-off! "
		}
		switch play.OutcomeID {
		case util.TipoffHomeWin, util.TipoffAwayWin:
			statement = prefix + team + " wins the tip and takes possession."
		default:
			statement = prefix
		}

	case util.Shot_three, util.Shot_corner_three, util.Shot_inside, util.Shot_paint, util.Shot_midrange:
		statement = carrier + " attempts " + shotDescription(play.EventID) + "..."
		switch play.OutcomeID {
		case util.Shot_made:
			statement += " GOOD!"
			if assister != "" {
				statement += " Assisted by " + assister + "."
			}
		case util.Shot_foul_made:
			statement += " GOOD, and " + fouler + " is called for the foul! " + carrier + " will go to the line for a bonus free throw."
		case util.Shot_missed:
			statement += " MISSED!"
		case util.Shot_foul_missed:
			statement += " MISSED, but " + fouler + " is called for a shooting foul! " + carrier + " heads to the line for three free throws."
		case util.Shot_blocked:
			statement += " BLOCKED by " + blocker + "!"
		case util.Shot_foul_blocked:
			statement += " BLOCKED by " + blocker + ", but a shooting foul is called on " + fouler + "!"
		}

	case util.Heave:
		statement = carrier + " heaves a desperation shot from deep..."
		switch play.OutcomeID {
		case util.Heave_made:
			statement += " AND IT GOES IN! What a shot for " + team + "!"
		case util.Heave_missed:
			statement += " but it misses the mark."
		}

	case util.Free_throw:
		statement = carrier + " steps to the line..."
		switch play.OutcomeID {
		case util.Ft_made:
			statement += " and sinks the free throw!"
		case util.Ft_missed:
			statement += " and misses the free throw."
		}

	case util.Rebound:
		switch play.OutcomeID {
		case util.Offensive_rebound:
			statement = carrier + " grabs the offensive rebound for " + team + "!"
		case util.Defensive_rebound:
			statement = carrier + " secures the defensive rebound."
		}

	case util.Inbound:
		statement = carrier + " inbounds the ball."

	case util.Move:
		switch play.OutcomeID {
		case util.Move_success:
			statement = carrier + " drives and works past " + defender + "."
		case util.Move_cutoff:
			statement = carrier + " tries to drive, but " + defender + " cuts off the lane."
		case util.Move_foul:
			statement = fouler + " fouls " + carrier + " on the drive!"
		case util.Offensive_charge:
			statement = carrier + " charges into " + defender + "! Offensive foul, and the ball goes the other way."
		case util.Move_trapped:
			statement = carrier + " gets trapped by the defense and has nowhere to go."
		}

	case util.Pass_ball:
		switch play.OutcomeID {
		case util.Pass_success:
			statement = carrier + " finds " + receiver + " with the pass."
		case util.Pass_deflected:
			statement = "The pass from " + carrier + " is deflected by " + defender + " and recovered by " + team + "."
		case util.Pass_intercepted:
			statement = "The pass from " + carrier + " is intercepted by " + stealer + "!"
		case util.Pass_foul:
			statement = fouler + " fouls " + carrier + " on the pass."
		case util.No_passing_lane:
			statement = carrier + " looks for an open teammate but finds no passing lane."
		}

	case util.Steal:
		statement = stealer + " steals the ball from " + carrier + "!"

	case util.Turnover:
		switch play.OutcomeID {
		case util.Shot_clock_violation:
			statement = "Shot clock violation on " + team + "! " + carrier + " couldn't get a shot off in time."
		case util.Out_of_bounds_turnover:
			statement = carrier + " steps out of bounds. Turnover, " + team + "."
		default:
			statement = "Turnover, " + team + "."
		}

	case util.Timeout:
		statement = team + " calls a timeout."
	case util.QuarterOver:
		statement = "End of quarter " + strconv.Itoa(int(play.Quarter)) + "."
	case util.HalfOver:
		statement = "End of the half."
	case util.OvertimeStart:
		statement = "We're heading to overtime!"
	case util.OvertimeOver:
		statement = "End of overtime."
	case util.GameOver:
		statement = "That's the end of the game!"
	}

	return statement
}
