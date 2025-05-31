package managers

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	edmondsblossom "github.com/CalebRose/SimNBA/EdmondsBlossom"
	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/repository"
	"github.com/CalebRose/SimNBA/secrets"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
	"gorm.io/gorm"
)

func ImportMatchResultsToDB(Results structs.ImportMatchResultsDTO) {
	db := dbprovider.GetInstance().GetDB()
	fmt.Println("Starting import...")

	tsChn := make(chan structs.Timestamp)

	go func() {
		ts := GetTimestamp()
		tsChn <- ts
	}()

	timestamp := <-tsChn
	close(tsChn)

	matchType := ""

	if !timestamp.GamesARan {
		matchType = "A"
	} else if !timestamp.GamesBRan {
		matchType = "B"
	} else if !timestamp.GamesCRan {
		matchType = "C"
	} else if !timestamp.GamesDRan {
		matchType = "D"
	}

	var teamStats []structs.TeamStats
	var nbaTeamStats []structs.NBATeamStats

	// Import College Game Results
	for _, dto := range Results.CBBResults {
		record := make(chan structs.Match)
		go func() {
			asyncMatch := GetMatchByMatchId(dto.GameID)
			record <- asyncMatch
		}()

		gameRecord := <-record
		close(record)

		var playerStats []structs.CollegePlayerStats

		homeTeamChn := make(chan structs.Team)
		go func() {
			homeTeam := GetTeamByTeamID(strconv.Itoa(dto.TeamOne.ID))
			homeTeamChn <- homeTeam
		}()

		ht := <-homeTeamChn
		close(homeTeamChn)

		matchID, err := strconv.Atoi(dto.GameID)
		if err != nil {
			log.Fatalln("Could not convert string to int")
		}

		homeTeam := mapToCollegeTeamStatsObject(ht.ID, uint(matchID), timestamp.CollegeWeekID, uint(timestamp.NBAWeek), timestamp.SeasonID, matchType, dto.TeamOne, dto.TeamTwo)

		teamStats = append(teamStats, homeTeam)

		awayTeamChn := make(chan structs.Team)
		go func() {
			awayTeam := GetTeamByTeamID(strconv.Itoa(dto.TeamTwo.ID))
			awayTeamChn <- awayTeam
		}()

		at := <-awayTeamChn
		close(awayTeamChn)

		awayTeam := mapToCollegeTeamStatsObject(at.ID, uint(matchID), timestamp.CollegeWeekID, uint(timestamp.NBAWeek), timestamp.SeasonID, matchType, dto.TeamTwo, dto.TeamOne)

		teamStats = append(teamStats, awayTeam)

		for _, player := range dto.RosterOne {
			id := player.ID
			collegePlayerStats := mapToCBBPlayerStatsObject(player, id, matchID, timestamp.SeasonID, timestamp.CollegeWeekID, uint(timestamp.NBAWeek), matchType)
			playerStats = append(playerStats, collegePlayerStats)
		}

		for _, player := range dto.RosterTwo {
			id := player.ID
			collegePlayerStats := mapToCBBPlayerStatsObject(player, id, matchID, timestamp.SeasonID, timestamp.CollegeWeekID, uint(timestamp.NBAWeek), matchType)
			playerStats = append(playerStats, collegePlayerStats)
		}

		gameRecord.UpdateScore(dto.TeamOne.Stats.Points, dto.TeamTwo.Stats.Points)

		err = db.Save(&gameRecord).Error
		if err != nil {
			log.Panicln("Could not save Game " + strconv.Itoa(int(gameRecord.ID)) + "Between " + gameRecord.HomeTeam + " and " + gameRecord.AwayTeam)
		}

		err = db.CreateInBatches(&playerStats, len(playerStats)).Error
		if err != nil {
			log.Panicln("Could not save player stats from week " + strconv.Itoa(timestamp.CollegeWeek))
		}

		fmt.Println("Finished Game " + strconv.Itoa(int(gameRecord.ID)) + " Between " + gameRecord.HomeTeam + " and " + gameRecord.AwayTeam)
	}

	// Import NBA Game Results
	for _, dto := range Results.NBAResults {
		record := make(chan structs.NBAMatch)
		go func() {
			asyncMatch := GetNBAMatchByMatchId(dto.GameID)
			record <- asyncMatch
		}()

		gameRecord := <-record
		close(record)

		var playerStats []structs.NBAPlayerStats

		homeTeamChn := make(chan structs.NBATeam)
		go func() {
			homeTeam := GetNBATeamByTeamID(strconv.Itoa(dto.TeamOne.ID))
			homeTeamChn <- homeTeam
		}()

		ht := <-homeTeamChn
		close(homeTeamChn)

		matchID := util.ConvertStringToInt(dto.GameID)

		homeTeam := mapToNBATeamStatsObject(ht.ID, uint(matchID), timestamp.NBAWeekID, uint(timestamp.NBAWeek), timestamp.SeasonID, matchType, dto.TeamOne, dto.TeamTwo)

		nbaTeamStats = append(nbaTeamStats, homeTeam)

		awayTeamChn := make(chan structs.NBATeam)
		go func() {
			awayTeam := GetNBATeamByTeamID(strconv.Itoa(dto.TeamTwo.ID))
			awayTeamChn <- awayTeam
		}()

		at := <-awayTeamChn
		close(awayTeamChn)

		awayTeam := mapToNBATeamStatsObject(at.ID, uint(matchID), timestamp.NBAWeekID, uint(timestamp.NBAWeek), timestamp.SeasonID, matchType, dto.TeamTwo, dto.TeamOne)

		nbaTeamStats = append(nbaTeamStats, awayTeam)

		for _, player := range dto.RosterOne {
			id := player.ID
			nbaPlayerStats := mapToNBAPlayerStatsObject(player, id, matchID, timestamp.SeasonID, timestamp.NBAWeekID, uint(timestamp.NBAWeek), matchType)
			playerStats = append(playerStats, nbaPlayerStats)
		}

		for _, player := range dto.RosterTwo {
			id := player.ID
			nbaPlayerStats := mapToNBAPlayerStatsObject(player, id, matchID, timestamp.SeasonID, timestamp.NBAWeekID, uint(timestamp.NBAWeek), matchType)
			playerStats = append(playerStats, nbaPlayerStats)
		}

		gameRecord.UpdateScore(dto.TeamOne.Stats.Points, dto.TeamTwo.Stats.Points)

		err := db.Save(&gameRecord).Error
		if err != nil {
			log.Panicln("Could not save Game " + strconv.Itoa(int(gameRecord.ID)) + "Between " + gameRecord.HomeTeam + " and " + gameRecord.AwayTeam)
		}

		err = db.CreateInBatches(&playerStats, len(playerStats)).Error
		if err != nil {
			log.Panicln("Could not save player stats from week " + strconv.Itoa(timestamp.CollegeWeek))
		}

		fmt.Println("Finished Game " + strconv.Itoa(int(gameRecord.ID)) + " Between " + gameRecord.HomeTeam + " and " + gameRecord.AwayTeam)
	}

	// Import all college team stats
	for _, stats := range teamStats {
		err := db.Create(&stats).Error
		if err != nil {
			log.Panicln("Could not save team stats!")
		}
	}
	// Import All nba team stats
	for _, stats := range nbaTeamStats {
		err := db.Create(&stats).Error
		if err != nil {
			log.Panicln("Could not save team stats!")
		}
	}
	fmt.Println("Finished Import for all games")
}

func ImportNBATeamsAndArenas() {
	db := dbprovider.GetInstance().GetDB()
	path := secrets.GetPath()["nbateams"]
	nbaTeamsCSV := util.ReadCSV(path)

	for idx, row := range nbaTeamsCSV {
		if idx < 2 {
			continue
		}

		id := util.ConvertStringToInt(row[0])
		team := row[1]
		nickname := row[2]
		abbr := row[3]
		conferenceID := util.ConvertStringToInt(row[5])
		conference := row[6]
		city := row[8]
		state := row[9]
		country := row[10]
		arenaID := util.ConvertStringToInt(row[11])
		arena := row[12]
		// capacity := util.ConvertStringToInt(row[13])

		nbaTeam := structs.NBATeam{
			Team:         team,
			Nickname:     nickname,
			Abbr:         abbr,
			City:         city,
			State:        state,
			Country:      country,
			LeagueID:     2,
			League:       "International Super League",
			ConferenceID: uint(conferenceID),
			Conference:   conference,
			ArenaID:      uint(arenaID),
			Arena:        arena,
			IsActive:     true,
		}

		nbaTeam.AssignID(uint(id))

		db.Create(&nbaTeam)

		// teamArena := structs.Arena{
		// 	Model: gorm.Model{
		// 		ID: uint(arenaID),
		// 	},
		// 	ArenaName: arena,
		// 	City:      city,
		// 	State:     state,
		// 	Country:   country,
		// 	Capacity:  uint(capacity),
		// 	HomeTeam:  team,
		// }

		// teamArena.AssignID(uint(arenaID))

		// err := db.Create(&arena).Error
		// if err != nil {
		// 	log.Panicln(err)
		// }
	}

}

func SyncContractValues() {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	SalaryCap := ts.Y1Capspace

	players := GetAllNBAPlayers()

	for _, player := range players {
		val := 0.0

		// Check if Max or Supermax qualified
		if player.Overall > 109 && !player.MaxRequested {
			player.ToggleMaxRequested()
		}

		if (player.IsDPOY || player.IsMVP || player.IsFirstTeamANBA) && !player.IsSuperMaxQualified {
			// SUPER MAX
			player.ToggleSuperMax()
		}

		if player.IsSuperMaxQualified {
			if player.Year > 9 {
				val = SalaryCap * 0.35
			} else if player.Year > 6 {
				val = SalaryCap * 0.3
			} else {
				val = SalaryCap * 0.25
			}
		} else if player.MaxRequested {
			if player.Year > 9 {
				val = SalaryCap * 0.3
			} else if player.Year > 6 {
				val = SalaryCap * 0.25
			} else {
				val = SalaryCap * 0.2
			}
		} else {
			if player.Year > 9 {
				val = 2.5
			} else if player.Year > 8 {
				val = 2.4
			} else if player.Year > 7 {
				val = 2.3
			} else if player.Year > 6 {
				val = 2.2
			} else if player.Year > 5 {
				val = 2.0
			} else if player.Year > 4 {
				val = 1.9
			} else if player.Year > 3 {
				val = 1.8
			} else if player.Year > 2 {
				val = 1.7
			} else if player.Year > 1 {
				val = 1.6
			} else if player.Year > 0 {
				val = 1.5
			} else {
				val = 0.9
			}
		}
		player.AssignMinimumContractValue(val)

		db.Save(&player)
	}
}

func ImportFAPreferences() {
	fmt.Println(time.Now().UnixNano())
	db := dbprovider.GetInstance().GetDB()
	nbaPlayers := GetAllNBAPlayers()

	for _, p := range nbaPlayers {
		NegotiationRound := 0
		if p.Overall > 95 {
			NegotiationRound = util.GenerateIntFromRange(2, 4)
		} else {
			NegotiationRound = util.GenerateIntFromRange(3, 6)
		}

		SigningRound := NegotiationRound + 1

		p.AssignFAPreferences(uint(NegotiationRound), uint(SigningRound))

		repository.SaveProfessionalPlayerRecord(p, db)
	}
}

func ImportMinutesExpectations() {
	db := dbprovider.GetInstance().GetDB()
	nbaPlayers := GetAllNBAPlayers()

	for _, p := range nbaPlayers {
		if p.PlaytimeExpectations > 0 {
			continue
		}
		minutes := util.GetProfessionalPlaytimeExpectations(p.Age, int(p.PrimeAge), p.Overall)
		if minutes < 0 {
			minutes = util.GenerateNormalizedIntFromRange(0, 12)
		}
		p.SetMinutesExpectations(minutes)
		repository.SaveProfessionalPlayerRecord(p, db)
	}
}

func ImportPersonalities() {
	fmt.Println(time.Now().UnixNano())
	db := dbprovider.GetInstance().GetDB()
	recruits := GetAllRecruitRecords()

	for _, p := range recruits {
		if len(p.RecruitingBias) > 0 {
			continue
		}
		recruitingBias := util.GetRecruitingBias()

		p.SetRecruitingBias(recruitingBias)

		db.Save(&p)
	}

	collegePlayers := GetAllCollegePlayers()

	for _, p := range collegePlayers {
		if len(p.RecruitingBias) > 0 {
			continue
		}
		recruitingBias := util.GetRecruitingBias()

		p.SetRecruitingBias(recruitingBias)

		db.Save(&p)
	}
}

func ImportCBBGames() {
	db := dbprovider.GetInstance().GetDB()
	path := secrets.GetPath()["cbbmatches"]
	collegeMatches := util.ReadCSV(path)

	collegeTeams := GetAllActiveCollegeTeams()
	collegeMap := make(map[string]structs.Team)

	for _, t := range collegeTeams {
		collegeMap[t.Abbr] = t
	}

	for idx, row := range collegeMatches {
		if idx < 1 {
			continue
		}

		id := util.ConvertStringToInt(row[0])
		season := util.ConvertStringToInt(row[1])
		seasonID := season - 2020
		week := util.ConvertStringToInt(row[2])
		weekID := week + 20
		timeSlot := row[3]
		matchType := row[4]
		isConf := false
		if matchType == "CONF" {
			isConf = true
		}
		homeTeamAbbr := row[6]
		awayTeamAbbr := row[7]
		htRankStr := row[5]
		atRankStr := row[8]
		htRank := 0
		atRank := 0
		if htRankStr != "" {
			htRank = util.ConvertStringToInt(htRankStr)
		}
		if atRankStr != "" {
			atRank = util.ConvertStringToInt(atRankStr)
		}

		homeTeam := collegeMap[homeTeamAbbr]
		awayTeam := collegeMap[awayTeamAbbr]
		gameTitle := row[22]
		nextGameID := util.ConvertStringToInt(row[24])
		hoA := row[25]
		neutralSite := util.ConvertStringToBool(row[10])
		invitational := util.ConvertStringToBool(row[11])
		conferenceTournament := util.ConvertStringToBool(row[12])
		cbi := util.ConvertStringToBool(row[13])
		nit := util.ConvertStringToBool(row[14])
		tournament := util.ConvertStringToBool(row[15])
		nationalChamp := util.ConvertStringToBool(row[16])
		arena := row[19]
		city := row[20]
		state := row[21]
		homeCoach := homeTeam.Coach
		if homeCoach == "" {
			homeCoach = "AI"
		}
		awayCoach := awayTeam.Coach
		if awayCoach == "" {
			awayCoach = "AI"
		}

		match := structs.Match{
			Model:                  gorm.Model{ID: uint(id)},
			SeasonID:               uint(seasonID),
			WeekID:                 uint(weekID),
			Week:                   uint(week),
			MatchOfWeek:            timeSlot,
			IsConference:           isConf,
			HomeTeam:               homeTeamAbbr,
			HomeTeamID:             homeTeam.ID,
			AwayTeamID:             awayTeam.ID,
			HomeTeamCoach:          homeCoach,
			HomeTeamRank:           uint(htRank),
			AwayTeam:               awayTeamAbbr,
			AwayTeamCoach:          awayCoach,
			AwayTeamRank:           uint(atRank),
			MatchName:              gameTitle,
			NextGameID:             uint(nextGameID),
			NextGameHOA:            hoA,
			IsNeutralSite:          neutralSite,
			IsInvitational:         invitational,
			IsConferenceTournament: conferenceTournament,
			IsNITGame:              nit,
			IsCBIGame:              cbi,
			IsPlayoffGame:          tournament,
			IsNationalChampionship: nationalChamp,
			Arena:                  arena,
			City:                   city,
			State:                  state,
		}

		db.Create(&match)
	}
}

func ImportNBAGames() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	nbaTeams := GetOnlyNBATeams()
	schedule, err := GenerateNBASeasonSchedule(nbaTeams, ts.SeasonID, 5000, 68)
	if err != nil {
		log.Println("Generation Failed: ", err)
		return
	}
	finalUpload := []structs.NBAMatch{}
	for _, game := range schedule {
		seasonBase := ts.Season - 2000 // 2025
		game.AddWeekData((uint(seasonBase)*100 + uint(game.Round)), uint(game.Round), game.Slot)
		finalUpload = append(finalUpload, game.NBAMatch)
	}
	repository.CreateNBARecordsBatch(db, finalUpload, 100)
}

func ImportISLGames() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	nbaTeams := GetInternationalTeams()

	fullSchedule, err := GenerateISLSeasonSchedule(nbaTeams, ts.SeasonID, 15000)
	if err != nil {
		log.Println("Generation Failed: ", err)
		return
	}

	finalUpload := []structs.NBAMatch{}
	for _, game := range fullSchedule {
		seasonBase := ts.Season - 2000 // 2025
		game.AddWeekData((uint(seasonBase)*100 + uint(game.Round)), uint(game.Round), game.Slot)
		finalUpload = append(finalUpload, game.NBAMatch)
	}
	repository.CreateNBARecordsBatch(db, finalUpload, 100)
}

type ScheduledGame struct {
	structs.NBAMatch
	Round int
	Slot  string // "A", "B", "C", or "D"
}

func (sg *ScheduledGame) ApplyRoundAndSlot(round int, slot string) {
	sg.Round = round
	sg.TimeSlot = slot
}

func GenerateNBASeasonSchedule(
	teams []structs.NBATeam,
	seasonID uint,
	maxPartitionAttempts,
	totalRounds int,
) ([]ScheduledGame, error) {
	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)

	var master []ScheduledGame
	for i := 0; i < len(teams); i++ {
		for j := i + 1; j < len(teams); j++ {
			teamA := teams[i]
			teamB := teams[j]
			if teamA.ID == teamB.ID {
				continue
			}
			numGames := 0
			if teamA.DivisionID == teamB.DivisionID {
				numGames = 4
			} else {
				numGames = 2
			}
			half := numGames / 2
			extra := numGames % 2
			for x := 0; x < half+extra; x++ {
				g := generateProGameRecord(teamA, teamB, seasonID)
				master = append(master, ScheduledGame{NBAMatch: g})
			}
			// Team B hosts half
			for x := 0; x < half; x++ {
				g := generateProGameRecord(teamB, teamA, seasonID)
				master = append(master, ScheduledGame{NBAMatch: g})
			}
		}
	}
	counts := make(map[uint]int)
	for _, g := range master {
		counts[g.HomeTeamID]++
		counts[g.AwayTeamID]++
	}
	for _, t := range teams {
		if counts[t.ID] != totalRounds {
			log.Printf("⚠️ team %d has %d games (expected %d)\n",
				t.ID, counts[t.ID], totalRounds)
		}
	}

	rng.Shuffle(len(master), func(i, j int) {
		master[i], master[j] = master[j], master[i]
	})

	teamIDs := make([]uint, len(teams))
	for i, t := range teams {
		teamIDs[i] = t.ID
	}
	matchesPerRound := len(teams) / 2 // 32 teams → 16 games per round

	rounds, err := greedyPartition(rng, master, totalRounds, matchesPerRound, maxPartitionAttempts)
	if err != nil {
		return nil, err
	}
	var final []ScheduledGame
	for rIdx, roundGames := range rounds {
		roundNum := rIdx + 1
		// zero-based week index, then +1 to make it 1-based
		week := (roundNum-1)/4 + 1
		slot := []string{"A", "B", "C", "D"}[(roundNum-1)%4]

		for _, g := range roundGames {
			g.ApplyRoundAndSlot(week, slot)
			final = append(final, g)
		}
	}

	return final, nil
}

func GenerateISLSeasonSchedule(
	teams []structs.NBATeam,
	seasonID uint,
	maxPartitionAttempts int,
) ([]ScheduledGame, error) {
	// 1) Group teams by conference
	confGroups := make(map[uint][]structs.NBATeam)
	for _, t := range teams {
		confGroups[t.ConferenceID] = append(confGroups[t.ConferenceID], t)
	}
	// must be exactly 8 conferences of 20
	if len(confGroups) != 8 {
		return nil, fmt.Errorf("expected 8 conferences, got %d", len(confGroups))
	}
	for conf, grp := range confGroups {
		if len(grp) != 20 {
			return nil, fmt.Errorf("conf %d has %d teams (want 20)", conf, len(grp))
		}
	}

	// 2) Build each conference's 76 rounds of 10 games
	perConfRounds := make(map[uint][][]ScheduledGame, 8)
	totalRounds := 0
	for conf, grp := range confGroups {
		// leg1 = 19 rounds (single-leg RR)
		leg1 := makeRoundRobinRounds(grp, seasonID) // 19×10
		// leg2 = mirror of leg1
		leg2 := mirrorRounds(leg1) // 19×10
		// stack 4 legs => 76 rounds
		rounds := make([][]ScheduledGame, 0, 38)
		rounds = append(rounds, leg1...)
		rounds = append(rounds, leg2...)

		// sanity check
		if len(rounds) != 38 {
			return nil, fmt.Errorf("conf %d has %d rounds (want 76)", conf, len(rounds))
		}
		for i, rnd := range rounds {
			if len(rnd) != 10 {
				return nil, fmt.Errorf("conf %d round %d has %d games (want 10)", conf, i, len(rnd))
			}
		}

		perConfRounds[conf] = rounds
		totalRounds = len(rounds)
	}

	// 3) Zip into league rounds of 8×10=80 games
	confIDs := make([]uint, 0, len(perConfRounds))
	for c := range perConfRounds {
		confIDs = append(confIDs, c)
	}
	sort.Slice(confIDs, func(i, j int) bool { return confIDs[i] < confIDs[j] })

	leagueRounds := make([][]ScheduledGame, totalRounds)
	for r := 0; r < totalRounds; r++ {
		leagueRounds[r] = make([]ScheduledGame, 0, 80)
		for _, conf := range confIDs {
			leagueRounds[r] = append(leagueRounds[r], perConfRounds[conf][r]...)
		}
		if len(leagueRounds[r]) != 80 {
			return nil, fmt.Errorf("league round %d has %d games (want 80)", r, len(leagueRounds[r]))
		}
	}

	// 4) Slot into weeks A/B/C and flatten
	var final []ScheduledGame
	slots := []string{"A", "B", "C"}
	for r, rnd := range leagueRounds {
		week := (r / 3) + 1
		slot := slots[r%3]
		for _, g := range rnd {
			g.ApplyRoundAndSlot(week, slot)
			final = append(final, g)
		}
	}

	// 5) Final check: each team must appear exactly totalRounds times
	counts := make(map[uint]int, len(teams))
	for _, g := range final {
		counts[g.HomeTeamID]++
		counts[g.AwayTeamID]++
	}
	for _, t := range teams {
		if counts[t.ID] != totalRounds {
			return nil, fmt.Errorf("team %d has %d games (want %d)", t.ID, counts[t.ID], totalRounds)
		}
	}

	return final, nil
}

func greedyPartition(
	rng *rand.Rand,
	allGames []ScheduledGame,
	totalRounds, matchesPerRound, perRoundMaxAttempts int,
) ([][]ScheduledGame, error) {
	// start with full pool
	remaining := append([]ScheduledGame(nil), allGames...)
	rounds := make([][]ScheduledGame, 0, totalRounds)

	for round := 1; round <= totalRounds; round++ {
		var thisRound []ScheduledGame
		success := false

		// keep an immutable snapshot for retries
		orig := append([]ScheduledGame(nil), remaining...)

		for attempt := 1; attempt <= perRoundMaxAttempts; attempt++ {
			// copy the snapshot
			pool := append([]ScheduledGame(nil), orig...)
			// shuffle the copy
			rng.Shuffle(len(pool), func(i, j int) {
				pool[i], pool[j] = pool[j], pool[i]
			})

			used := make(map[uint]bool, matchesPerRound*2)
			candidate := make([]ScheduledGame, 0, matchesPerRound)

			// greedily pull non‑conflicting games
			for i := 0; i < len(pool) && len(candidate) < matchesPerRound; {
				g := pool[i]
				if !used[g.HomeTeamID] && !used[g.AwayTeamID] {
					used[g.HomeTeamID] = true
					used[g.AwayTeamID] = true
					candidate = append(candidate, g)
					// remove from pool
					pool = append(pool[:i], pool[i+1:]...)
				} else {
					i++
				}
			}

			if len(candidate) == matchesPerRound {
				// success: lock in this round and remaining pool
				thisRound = candidate
				remaining = pool
				success = true
				break
			}
			// else: try again with fresh snapshot
		}

		if !success {
			return nil, fmt.Errorf(
				"round %d: failed to fill after %d attempts (got %d/%d games)",
				round, perRoundMaxAttempts, len(thisRound), matchesPerRound,
			)
		}

		rounds = append(rounds, thisRound)
	}

	if len(remaining) != 0 {
		return nil, fmt.Errorf(
			"leftover games after %d rounds: %d",
			totalRounds, len(remaining),
		)
	}
	return rounds, nil
}

func ImportNBAGamesOLD() {
	db := dbprovider.GetInstance().GetDB()
	path := secrets.GetPath()["nbamatches"]
	professionalMatches := util.ReadCSV(path)

	professionalTeams := GetAllActiveNBATeams()
	teamMap := make(map[string]structs.NBATeam)

	for _, t := range professionalTeams {
		teamMap[t.Team+" "+t.Nickname] = t
	}

	for idx, row := range professionalMatches {
		if idx < 1 {
			continue
		}

		id := util.ConvertStringToInt(row[0])
		season := util.ConvertStringToInt(row[1])
		seasonID := season - 2020
		week := util.ConvertStringToInt(row[2])
		weekID := week
		timeSlot := row[3]
		homeTeamStr := row[6]
		awayTeamStr := row[7]
		homeTeam := teamMap[homeTeamStr]
		awayTeam := teamMap[awayTeamStr]
		gameTitle := ""
		nextGameID := 0
		hoA := ""
		conference := util.ConvertStringToBool(row[12])
		divisional := util.ConvertStringToBool(row[13])
		international := util.ConvertStringToBool(row[14])
		arena := row[19]
		city := row[20]
		state := row[20]
		country := row[22]
		homeCoach := homeTeam.NBACoachName
		if homeCoach == "" {
			homeCoach = homeTeam.NBAOwnerName
			if homeCoach == "" {
				homeCoach = "AI"
			}
		}
		awayCoach := awayTeam.NBACoachName
		if awayCoach == "" {
			awayCoach = awayTeam.NBAOwnerName
			if awayCoach == "" {
				awayCoach = "AI"
			}
		}

		match := structs.NBAMatch{
			Model:           gorm.Model{ID: uint(id)},
			SeasonID:        uint(seasonID),
			WeekID:          uint(weekID),
			Week:            uint(week),
			MatchOfWeek:     timeSlot,
			IsConference:    conference,
			IsDivisional:    divisional,
			HomeTeam:        homeTeamStr,
			HomeTeamID:      homeTeam.ID,
			AwayTeamID:      awayTeam.ID,
			HomeTeamCoach:   homeCoach,
			AwayTeam:        awayTeamStr,
			AwayTeamCoach:   awayCoach,
			MatchName:       gameTitle,
			NextGameID:      uint(nextGameID),
			NextGameHOA:     hoA,
			IsNeutralSite:   conference,
			IsPlayoffGame:   false,
			IsTheFinals:     false,
			IsInternational: international,
			Arena:           arena,
			City:            city,
			State:           state,
			Country:         country,
		}

		db.Create(&match)
	}
}

func ImportNBASeries() {
	db := dbprovider.GetInstance().GetDB()
	path := secrets.GetPath()["nbaseries"]
	professionalMatches := util.ReadCSV(path)
	professionalTeams := GetAllActiveNBATeams()
	teamMap := make(map[string]structs.NBATeam)

	for _, t := range professionalTeams {
		teamStr := t.Team + " " + t.Nickname
		trimmedStr := strings.TrimSpace(teamStr)
		teamMap[trimmedStr] = t
	}

	for idx, row := range professionalMatches {
		if idx < 1 {
			continue
		}

		id := util.ConvertStringToInt(row[0])
		season := util.ConvertStringToInt(row[1])
		seasonID := season - 2020
		homeTeamStr := strings.TrimSpace(row[5])
		awayTeamStr := strings.TrimSpace(row[6])
		homeTeam := teamMap[homeTeamStr]
		homeTeamRank := util.ConvertStringToInt(row[4])
		awayTeamRank := util.ConvertStringToInt(row[7])
		awayTeam := teamMap[awayTeamStr]
		nextGameID := util.ConvertStringToInt(row[14])
		hoA := row[15]
		seriesTitle := row[13]
		international := util.ConvertStringToBool(row[9])
		playoff := util.ConvertStringToBool(row[8])
		finals := util.ConvertStringToBool(row[10])
		homeCoach := homeTeam.NBACoachName
		if homeCoach == "" {
			homeCoach = homeTeam.NBAOwnerName
			if homeCoach == "" {
				homeCoach = "AI"
			}
		}
		awayCoach := awayTeam.NBACoachName
		if awayCoach == "" {
			awayCoach = awayTeam.NBAOwnerName
			if awayCoach == "" {
				awayCoach = "AI"
			}
		}

		match := structs.NBASeries{
			Model:           gorm.Model{ID: uint(id)},
			SeriesName:      seriesTitle,
			SeasonID:        uint(seasonID),
			HomeTeam:        homeTeamStr,
			HomeTeamID:      homeTeam.ID,
			AwayTeamID:      awayTeam.ID,
			HomeTeamCoach:   homeCoach,
			HomeTeamWins:    0,
			HomeTeamWin:     false,
			HomeTeamRank:    uint(homeTeamRank),
			AwayTeam:        awayTeamStr,
			AwayTeamCoach:   awayCoach,
			AwayTeamWins:    0,
			AwayTeamWin:     false,
			AwayTeamRank:    uint(awayTeamRank),
			GameCount:       0,
			NextSeriesID:    uint(nextGameID),
			NextSeriesHOA:   hoA,
			IsPlayoffGame:   playoff,
			IsTheFinals:     finals,
			IsInternational: international,
			SeriesComplete:  false,
		}

		db.Create(&match)
	}
}

func RollbackNBAGames() {
	db := dbprovider.GetInstance().GetDB()
	path := secrets.GetPath()["nbamatches"]
	professionalMatches := util.ReadCSV(path)

	professionalTeams := GetAllActiveNBATeams()
	teamMap := make(map[string]structs.NBATeam)
	// teamStatMap := make(map[uint][]structs.NBATeamStats)
	// teamStats := GetNBATeamStatsBySeasonID("3")
	existingMatchMap := make(map[uint]structs.NBAMatch)

	nbaMatches := GetNBAMatchesBySeasonID("3")
	islMatches := GetISLMatchesBySeasonID("3")

	for _, match := range nbaMatches {
		existingMatchMap[match.ID] = match
	}

	for _, match := range islMatches {
		existingMatchMap[match.ID] = match
	}

	// Collect all team stats
	// for _, stat := range teamStats {
	// 	id := stat.MatchID
	// 	teamStatMap[id] = append(teamStatMap[id], stat)
	// }

	for _, t := range professionalTeams {
		teamStr := t.Team + " " + t.Nickname
		trimmedStr := strings.TrimSpace(teamStr)
		teamMap[trimmedStr] = t
	}

	for idx, row := range professionalMatches {
		if idx < 1313 {
			continue
		}

		id := util.ConvertStringToInt(row[0])
		match := existingMatchMap[uint(id)]
		// No need for stats updates considering that the scores are the same. Only team information + arena
		homeTeamStr := strings.TrimSpace(row[6])
		awayTeamStr := strings.TrimSpace(row[7])
		homeTeam := teamMap[homeTeamStr]
		awayTeam := teamMap[awayTeamStr]
		city := homeTeam.City
		arena := homeTeam.Arena
		state := homeTeam.State
		homeCoach := homeTeam.NBACoachName
		if homeCoach == "" {
			homeCoach = homeTeam.NBAOwnerName
			if homeCoach == "" {
				homeCoach = "AI"
			}
		}
		awayCoach := awayTeam.NBACoachName
		if awayCoach == "" {
			awayCoach = awayTeam.NBAOwnerName
			if awayCoach == "" {
				awayCoach = "AI"
			}
		}
		match.RollbackMatch(homeTeamStr, homeCoach, awayTeamStr, awayCoach, city, state, arena, homeTeam.ID, awayTeam.ID)
		db.Save(&match)
	}
}

func ImportNBAStandings() {
	db := dbprovider.GetInstance().GetDB()
	path := secrets.GetPath()["nbastandings"]
	nbaStandingsCSV := util.ReadCSV(path)

	for idx, row := range nbaStandingsCSV {
		if idx < 1 {
			continue
		}

		id := util.ConvertStringToInt(row[0])
		teamID := util.ConvertStringToInt(row[1])
		team := row[2]
		seasonID := util.ConvertStringToInt(row[3])
		season := util.ConvertStringToInt(row[4])
		leagueID := util.ConvertStringToInt(row[5])
		league := row[6]
		conferenceID := util.ConvertStringToInt(row[7])
		conference := row[8]
		divisionID := util.ConvertStringToInt(row[9])
		division := row[10]
		postSeasonStatus := row[11]
		isConferenceChampion := util.ConvertStringToBool(row[12])
		totalWins := util.ConvertStringToInt(row[13])
		totalLosses := util.ConvertStringToInt(row[14])

		standings := structs.NBAStandings{
			TeamID:               uint(teamID),
			TeamName:             team,
			SeasonID:             uint(seasonID),
			Season:               season,
			LeagueID:             uint(leagueID),
			League:               league,
			ConferenceID:         uint(conferenceID),
			ConferenceName:       conference,
			DivisionID:           uint(divisionID),
			DivisionName:         division,
			PostSeasonStatus:     postSeasonStatus,
			IsConferenceChampion: isConferenceChampion,
			BaseStandings: structs.BaseStandings{
				TotalWins:   totalWins,
				TotalLosses: totalLosses,
			},
		}

		standings.AssignID(uint(id))

		db.Create(&standings)
	}
}

func ImportNewPositions() {
	db := dbprovider.GetInstance().GetDB()

	collegePlayers := GetAllCollegePlayers()
	recruits := GetAllRecruitRecords()

	for _, c := range collegePlayers {
		if c.Position == "C" {
			continue
		}
		shooting := (c.Shooting2 + c.Shooting3) / 2
		if c.Position == "G" {
			if shooting > c.Ballwork || c.Archetype == "Floor General" {
				c.SetNewPosition("PG")
			} else {
				c.SetNewPosition("SG")
			}
		} else {
			if c.Rebounding > shooting || c.Archetype == "Point Forward" {
				c.SetNewPosition("PF")
			} else {
				c.SetNewPosition("SF")
			}
		}
		db.Save(&c)
	}

	for _, r := range recruits {
		if r.Position == "C" {
			continue
		}
		shooting := (r.Shooting2 + r.Shooting3) / 2
		if r.Position == "G" {
			if shooting > r.Ballwork || r.Archetype == "Floor General" {
				r.SetNewPosition("PG")
			} else {
				r.SetNewPosition("SG")
			}
		} else {
			if r.Rebounding > shooting || r.Archetype == "Point Forward" {
				r.SetNewPosition("PF")
			} else {
				r.SetNewPosition("SF")
			}
		}
		db.Save(&r)
	}
}

func ImportNewTeams() {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()
	path := secrets.GetPath()["teams"]
	teams := util.ReadCSV(path)

	for idx, row := range teams {
		if idx == 0 {
			continue
		}
		teamID := util.ConvertStringToInt(row[0])
		team := row[1]
		nickname := row[2]
		abbr := row[3]
		city := row[4]
		state := row[5]
		country := row[6]
		conferenceID := util.ConvertStringToInt(row[7])
		conference := row[8]
		season := row[10]
		isActive := util.ConvertStringToBool(row[13])

		t := structs.Team{
			Team:         team,
			Nickname:     nickname,
			Abbr:         abbr,
			City:         city,
			State:        state,
			Country:      country,
			ConferenceID: uint(conferenceID),
			Conference:   conference,
			IsActive:     isActive,
			IsNBA:        false,
			FirstSeason:  season,
			Model:        gorm.Model{ID: uint(teamID)},
		}

		db.Create(&t)

		standings := structs.CollegeStandings{
			TeamID:         uint(teamID),
			TeamName:       team,
			TeamAbbr:       abbr,
			SeasonID:       ts.SeasonID,
			Season:         ts.Season,
			ConferenceID:   uint(conferenceID),
			ConferenceName: conference,
		}

		db.Create(&standings)
	}
}

func ImportDraftPicks() {
	db := dbprovider.GetInstance().GetDB()

	path := secrets.GetPath()["draftpicks"]
	picks := util.ReadCSV(path)

	for idx, row := range picks {
		if idx == 0 {
			continue
		}
		id := util.ConvertStringToInt(row[0])
		round := util.ConvertStringToInt(row[1])
		draftNumber := util.ConvertStringToInt(row[2])
		season_id := util.ConvertStringToInt(row[3])
		season := util.ConvertStringToInt(row[4])
		drafteeID := 0
		teamID := util.ConvertStringToInt(row[5])
		team := row[6]
		originalTeamID := util.ConvertStringToInt(row[7])
		originalTeam := row[8]
		notes := ""

		draftPick := structs.DraftPick{
			DraftRound:     uint(round),
			DraftNumber:    uint(draftNumber),
			DrafteeID:      uint(drafteeID),
			TeamID:         uint(teamID),
			Team:           team,
			OriginalTeamID: uint(originalTeamID),
			OriginalTeam:   originalTeam,
			Season:         uint(season),
			Notes:          notes,
			SeasonID:       uint(season_id),
			Model:          gorm.Model{ID: uint(id)},
		}

		db.Create(&draftPick)
	}
}

func ImportISLScoutingDepts() {
	db := dbprovider.GetInstance().GetDB()

	islTeams := GetAllActiveNBATeams()

	for _, t := range islTeams {
		// Iterate over all NBA Teams
		if t.LeagueID == 1 {
			continue
		}

		id := strconv.Itoa(int(t.ID))

		existingDept := GetScoutingDeptByTeamID(id)

		if existingDept.ID > 0 {
			continue
		}
		teamName := t.Team + " " + t.Nickname
		formattedName := strings.TrimSpace(teamName)
		prestige := util.GenerateNormalizedIntFromRange(1, 5)
		identityBias := prestige > 3
		behaviorBias := util.GenerateNormalizedIntFromRange(1, 3)
		bonusPoints := 5 + (2 * prestige)
		fn := 0
		sh2 := 0
		sh3 := 0
		ft := 0
		bw := 0
		rb := 0
		ind := 0
		prd := 0
		pot := 0
		idn := 0

		selectionList := []string{"fn", "sh2", "sh3", "ft", "bw", "rb", "ind", "prd", "pot", "idn"}
		rand.Shuffle(len(selectionList), func(i, j int) {
			selectionList[i], selectionList[j] = selectionList[j], selectionList[i]
		})

		for bonusPoints > 0 {
			for _, attr := range selectionList {
				num := util.GenerateIntFromRange(1, 5)
				if bonusPoints-num < 0 {
					num += (bonusPoints - num)
				}
				if attr == "fn" {
					fn += num
				} else if attr == "sh2" {
					sh2 += num
				} else if attr == "sh3" {
					sh3 += num
				} else if attr == "ft" {
					ft += num
				} else if attr == "bw" {
					bw += num
				} else if attr == "rb" {
					rb += num
				} else if attr == "ind" {
					ind += num
				} else if attr == "prd" {
					prd += num
				} else if attr == "pot" {
					pot += num
				} else if attr == "idn" {
					idn += num
				}

				bonusPoints -= num
			}
		}

		newDept := structs.ISLScoutingDept{
			TeamID:         t.ID,
			TeamLabel:      formattedName,
			Prestige:       uint8(prestige),
			Resources:      100,
			IdentityPool:   0,
			ScoutingPool:   0,
			InvestingPool:  0,
			ScoutingCount:  0,
			IdentityBias:   identityBias,
			BehaviorBias:   uint8(behaviorBias),
			Finishing:      uint8(fn),
			Shooting2:      uint8(sh2),
			Shooting3:      uint8(sh3),
			FreeThrow:      uint8(ft),
			Ballwork:       uint8(bw),
			Rebounding:     uint8(rb),
			IntDefense:     uint8(ind),
			PerDefense:     uint8(prd),
			Potential:      uint8(pot),
			IdentityMod:    uint8(idn),
			ModifierPoints: 0,
		}

		db.Create(&newDept)
	}
}

func filterLotteryPicks(list []structs.DraftLottery, id uint) []structs.DraftLottery {
	newList := []structs.DraftLottery{}
	for _, l := range list {
		if l.ID != id {
			newList = append(newList, l)
		}
	}
	return newList
}

func mapToCollegeTeamStatsObject(teamID, matchID, weekID, week, seasonID uint, matchType string, TeamOne, TeamTwo structs.TeamResultsDTO) structs.TeamStats {
	return structs.TeamStats{
		TeamID:                    teamID,
		MatchID:                   matchID,
		WeekID:                    weekID,
		Week:                      week,
		SeasonID:                  seasonID,
		MatchType:                 matchType,
		Points:                    TeamOne.Stats.Points,
		Possessions:               TeamOne.Stats.Possessions,
		FGM:                       TeamOne.Stats.FGM,
		FGA:                       TeamOne.Stats.FGA,
		FGPercent:                 TeamOne.Stats.FGPercent,
		ThreePointsMade:           TeamOne.Stats.ThreePointsMade,
		ThreePointAttempts:        TeamOne.Stats.ThreePointAttempts,
		ThreePointPercent:         TeamOne.Stats.ThreePointPercent,
		FTM:                       TeamOne.Stats.FTM,
		FTA:                       TeamOne.Stats.FTA,
		FTPercent:                 TeamOne.Stats.FTPercent,
		Rebounds:                  TeamOne.Stats.Rebounds,
		OffRebounds:               TeamOne.Stats.OffRebounds,
		DefRebounds:               TeamOne.Stats.DefRebounds,
		Assists:                   TeamOne.Stats.Assists,
		Steals:                    TeamOne.Stats.Steals,
		Blocks:                    TeamOne.Stats.Blocks,
		TotalTurnovers:            TeamOne.Stats.TotalTurnovers,
		LargestLead:               TeamOne.Stats.LargestLead,
		FirstHalfScore:            TeamOne.Stats.FirstHalfScore,
		SecondHalfScore:           TeamOne.Stats.SecondHalfScore,
		OvertimeScore:             TeamOne.Stats.OvertimeScore,
		Fouls:                     TeamOne.Stats.Fouls,
		PointsAgainst:             TeamTwo.Stats.Points,
		FGMAgainst:                TeamTwo.Stats.FGM,
		FGAAgainst:                TeamTwo.Stats.FGA,
		FGPercentAgainst:          TeamTwo.Stats.FGPercent,
		ThreePointsMadeAgainst:    TeamTwo.Stats.ThreePointsMade,
		ThreePointAttemptsAgainst: TeamTwo.Stats.ThreePointAttempts,
		ThreePointPercentAgainst:  TeamTwo.Stats.ThreePointPercent,
		FTMAgainst:                TeamTwo.Stats.FTM,
		FTAAgainst:                TeamTwo.Stats.FTA,
		FTPercentAgainst:          TeamTwo.Stats.FTPercent,
		ReboundsAllowed:           TeamTwo.Stats.Rebounds,
		OffReboundsAllowed:        TeamTwo.Stats.OffRebounds,
		DefReboundsAllowed:        TeamTwo.Stats.DefRebounds,
		AssistsAllowed:            TeamTwo.Stats.Assists,
		StealsAllowed:             TeamTwo.Stats.Steals,
		BlocksAllowed:             TeamTwo.Stats.Blocks,
		TurnoversAllowed:          TeamTwo.Stats.TotalTurnovers,
	}
}

func mapToNBATeamStatsObject(teamID, matchID, weekID, week, seasonID uint, matchType string, TeamOne, TeamTwo structs.TeamResultsDTO) structs.NBATeamStats {
	return structs.NBATeamStats{
		TeamID:                    teamID,
		MatchID:                   matchID,
		WeekID:                    weekID,
		Week:                      week,
		SeasonID:                  seasonID,
		MatchType:                 matchType,
		Points:                    TeamOne.Stats.Points,
		Possessions:               TeamOne.Stats.Possessions,
		FGM:                       TeamOne.Stats.FGM,
		FGA:                       TeamOne.Stats.FGA,
		FGPercent:                 TeamOne.Stats.FGPercent,
		ThreePointsMade:           TeamOne.Stats.ThreePointsMade,
		ThreePointAttempts:        TeamOne.Stats.ThreePointAttempts,
		ThreePointPercent:         TeamOne.Stats.ThreePointPercent,
		FTM:                       TeamOne.Stats.FTM,
		FTA:                       TeamOne.Stats.FTA,
		FTPercent:                 TeamOne.Stats.FTPercent,
		Rebounds:                  TeamOne.Stats.Rebounds,
		OffRebounds:               TeamOne.Stats.OffRebounds,
		DefRebounds:               TeamOne.Stats.DefRebounds,
		Assists:                   TeamOne.Stats.Assists,
		Steals:                    TeamOne.Stats.Steals,
		Blocks:                    TeamOne.Stats.Blocks,
		TotalTurnovers:            TeamOne.Stats.TotalTurnovers,
		LargestLead:               TeamOne.Stats.LargestLead,
		FirstHalfScore:            TeamOne.Stats.FirstHalfScore,
		SecondQuarterScore:        TeamOne.Stats.SecondQuarterScore,
		SecondHalfScore:           TeamOne.Stats.SecondHalfScore,
		FourthQuarterScore:        TeamOne.Stats.FourthQuarterScore,
		OvertimeScore:             TeamOne.Stats.OvertimeScore,
		Fouls:                     TeamOne.Stats.Fouls,
		PointsAgainst:             TeamTwo.Stats.Points,
		FGMAgainst:                TeamTwo.Stats.FGM,
		FGAAgainst:                TeamTwo.Stats.FGA,
		FGPercentAgainst:          TeamTwo.Stats.FGPercent,
		ThreePointsMadeAgainst:    TeamTwo.Stats.ThreePointsMade,
		ThreePointAttemptsAgainst: TeamTwo.Stats.ThreePointAttempts,
		ThreePointPercentAgainst:  TeamTwo.Stats.ThreePointPercent,
		FTMAgainst:                TeamTwo.Stats.FTM,
		FTAAgainst:                TeamTwo.Stats.FTA,
		FTPercentAgainst:          TeamTwo.Stats.FTPercent,
		ReboundsAllowed:           TeamTwo.Stats.Rebounds,
		OffReboundsAllowed:        TeamTwo.Stats.OffRebounds,
		DefReboundsAllowed:        TeamTwo.Stats.DefRebounds,
		AssistsAllowed:            TeamTwo.Stats.Assists,
		StealsAllowed:             TeamTwo.Stats.Steals,
		BlocksAllowed:             TeamTwo.Stats.Blocks,
		TurnoversAllowed:          TeamTwo.Stats.TotalTurnovers,
	}
}

func mapToCBBPlayerStatsObject(player structs.PlayerDTO, id, matchID int, seasonID, weekID, week uint, matchType string) structs.CollegePlayerStats {
	return structs.CollegePlayerStats{
		TeamID:             uint(player.TeamID),
		CollegePlayerID:    uint(id),
		MatchID:            uint(matchID),
		SeasonID:           seasonID,
		MatchType:          matchType,
		WeekID:             weekID,
		Week:               week,
		Year:               uint(player.Stats.Year),
		Minutes:            player.Stats.Minutes,
		Possessions:        player.Stats.Possessions,
		FGM:                player.Stats.FGM,
		FGA:                player.Stats.FGA,
		FGPercent:          player.Stats.FGPercent,
		ThreePointsMade:    player.Stats.ThreePointsMade,
		ThreePointAttempts: player.Stats.ThreePointAttempts,
		ThreePointPercent:  player.Stats.ThreePointPercent,
		FTM:                player.Stats.FTM,
		FTA:                player.Stats.FTA,
		FTPercent:          player.Stats.FTPercent,
		Points:             player.Stats.Points,
		TotalRebounds:      player.Stats.TotalRebounds,
		OffRebounds:        player.Stats.OffRebounds,
		DefRebounds:        player.Stats.DefRebounds,
		Assists:            player.Stats.Assists,
		Steals:             player.Stats.Steals,
		Blocks:             player.Stats.Blocks,
		Turnovers:          player.Stats.Turnovers,
		Fouls:              player.Stats.Fouls,
		FouledOut:          player.Stats.FouledOut,
		IsInjured:          player.Stats.IsInjured,
		InjuryName:         player.Stats.InjuryName,
		InjuryType:         player.Stats.InjuryType,
		WeeksOfRecovery:    player.Stats.WeeksOfRecovery,
	}
}

func mapToNBAPlayerStatsObject(player structs.PlayerDTO, id, matchID int, seasonID, weekID, week uint, matchType string) structs.NBAPlayerStats {
	return structs.NBAPlayerStats{
		TeamID:             uint(player.TeamID),
		NBAPlayerID:        uint(id),
		MatchID:            uint(matchID),
		SeasonID:           seasonID,
		WeekID:             weekID,
		Week:               week,
		Year:               uint(player.Stats.Year),
		MatchType:          matchType,
		Minutes:            player.Stats.Minutes,
		Possessions:        player.Stats.Possessions,
		FGM:                player.Stats.FGM,
		FGA:                player.Stats.FGA,
		FGPercent:          player.Stats.FGPercent,
		ThreePointsMade:    player.Stats.ThreePointsMade,
		ThreePointAttempts: player.Stats.ThreePointAttempts,
		ThreePointPercent:  player.Stats.ThreePointPercent,
		FTM:                player.Stats.FTM,
		FTA:                player.Stats.FTA,
		FTPercent:          player.Stats.FTPercent,
		Points:             player.Stats.Points,
		TotalRebounds:      player.Stats.TotalRebounds,
		OffRebounds:        player.Stats.OffRebounds,
		DefRebounds:        player.Stats.DefRebounds,
		Assists:            player.Stats.Assists,
		Steals:             player.Stats.Steals,
		Blocks:             player.Stats.Blocks,
		Turnovers:          player.Stats.Turnovers,
		Fouls:              player.Stats.Fouls,
		FouledOut:          player.Stats.FouledOut,
		IsInjured:          player.Stats.IsInjured,
		InjuryName:         player.Stats.InjuryName,
		InjuryType:         player.Stats.InjuryType,
		WeeksOfRecovery:    player.Stats.WeeksOfRecovery,
	}
}

// makeRoundRobin returns N−1 rounds of a single home+away visit per opponent
func makeRoundRobinRounds(grp []structs.NBATeam, seasonID uint) [][]ScheduledGame {
	N := len(grp)
	rounds := make([][]ScheduledGame, N-1)

	// build the rotating array (teams[1..])
	rotating := make([]structs.NBATeam, N-1)
	copy(rotating, grp[1:])

	for r := 0; r < N-1; r++ {
		var games []ScheduledGame
		// 1) anchor vs its rotating opponent
		games = append(games, ScheduledGame{NBAMatch: generateProGameRecord(grp[0], rotating[r], seasonID)})
		// 2) the remaining N/2−1 cross‐pairs
		for i := 1; i < N/2; i++ {
			a := rotating[(r+i)%(N-1)]
			b := rotating[(r-i+(N-1))%(N-1)]
			games = append(games, ScheduledGame{NBAMatch: generateProGameRecord(a, b, seasonID)})
		}
		rounds[r] = games
	}
	return rounds
}

// mirrorRounds returns a “swap home/away” copy of each round
func mirrorRounds(legs [][]ScheduledGame) [][]ScheduledGame {
	out := make([][]ScheduledGame, len(legs))
	for i, rnd := range legs {
		mirror := make([]ScheduledGame, len(rnd))
		for j, sg := range rnd {
			rev := sg
			rev.HomeTeamID, rev.AwayTeamID = sg.AwayTeamID, sg.HomeTeamID
			rev.HomeTeam, rev.AwayTeam = sg.AwayTeam, sg.HomeTeam
			rev.Arena, rev.City, rev.State, rev.Country =
				sg.Arena, sg.City, sg.State, sg.Country
			mirror[j] = rev
		}
		out[i] = mirror
	}
	return out
}

func generateProGameRecord(teamA structs.NBATeam, teamB structs.NBATeam, seasonID uint) structs.NBAMatch {
	return structs.NBAMatch{
		SeasonID:        seasonID,
		HomeTeamID:      teamA.ID,
		HomeTeam:        teamA.Team,
		AwayTeamID:      teamB.ID,
		AwayTeam:        teamB.Team,
		Arena:           teamA.Arena,
		City:            teamA.City,
		State:           teamA.State,
		Country:         teamA.Country,
		IsConference:    teamA.ConferenceID == teamB.ConferenceID,
		IsDivisional:    teamA.DivisionID == teamB.DivisionID,
		IsInternational: teamA.ID > 32,
	}
}

// pair represents an undirected matchup between two teams
type pair struct{ A, B int }

// pickInterConferencePairs returns exactly 6 non-conference opponents per team,
// by making 6 “stubs” per team and randomly pairing them.
// It will either succeed or return an error after N attempts.
func pickInterConferencePairs(
	teams []structs.NBATeam,
) ([]pair, error) {
	N := len(teams)
	// map team ID -> index 0..N-1
	idxOf := make(map[uint]int, N)
	for i, t := range teams {
		idxOf[t.ID] = i
	}
	// build the set of *allowed* inter-conference edges
	allowed := make(map[pair]bool)
	for i := 0; i < N; i++ {
		for j := i + 1; j < N; j++ {
			if teams[i].ConferenceID != teams[j].ConferenceID {
				allowed[pair{A: i, B: j}] = true
			}
		}
	}

	var allPairs []pair
	// we need exactly 6 perfect matchings
	for round := 0; round < 6; round++ {
		// build a new Graph
		g := edmondsblossom.NewGraph(N)
		// add only the allowed edges
		for e := range allowed {
			u, v := int(e.A), int(e.B)
			g.G[u][v] = 1
			g.G[v][u] = 1
		}

		// find one perfect matching on the inter-conference graph
		edmondsblossom.FindMaxMatching(g)
		// ensure this was a *perfect* matching
		matched := 0
		for u := 0; u < N; u++ {
			if g.M[u] >= 0 {
				matched++
			}
		}
		log.Printf("After full matching: %d/%d matched\n", matched, N)

		// extract matched pairs, and remove them from allowed
		for u, v := range g.M {
			if v > u {
				key := pair{A: u, B: v}
				if !allowed[key] {
					// this should never happen if Blossom found only allowed edges
					return nil, fmt.Errorf("unexpected edge %v in matching", key)
				}
				allPairs = append(allPairs, key)
				// remove so we never pick this same inter-conf matchup again
				delete(allowed, key)
			}
		}

		// sanity check
		if len(allPairs) != (round+1)*(N/2) {
			return nil, fmt.Errorf("round %d: only %d pairs found", round+1, len(allPairs))
		}
	}

	return allPairs, nil

}

func BlossomPartition(teams []structs.NBATeam, master []ScheduledGame, totalRounds int) ([][]ScheduledGame, error) {
	teamIDs := make([]uint, len(teams))
	for i, t := range teams {
		teamIDs[i] = t.ID
	}
	idxOf := map[uint]int{}
	for i, id := range teamIDs {
		idxOf[id] = i
	}
	pool := map[pair][]ScheduledGame{}
	for _, g := range master { // master from your GenerateISL code
		A, B := idxOf[g.HomeTeamID], idxOf[g.AwayTeamID]
		if A > B {
			A, B = B, A
		}
		pool[pair{A, B}] = append(pool[pair{A, B}], g)
	}

	rounds := [][]ScheduledGame{}

	for round := 1; round <= totalRounds; round++ {
		// a) build a Graph of size N
		g := edmondsblossom.NewGraph(len(teams))
		// b) add an edge between i and j iff pool[[i,j]] is non-empty
		for key, games := range pool {
			if len(games) > 0 {
				g.G[key.A][key.B] = 1
				g.G[key.B][key.A] = 1
			}
		}
		// c) run blossom matching
		edmondsblossom.FindMaxMatching(g)
		// d) extract all matched pairs, pick one ScheduledGame from pool
		thisRound := make([]ScheduledGame, 0, len(teams)/2)
		used := make(map[pair]bool)
		for u, v := range g.M {
			if v > u {
				key := pair{A: u, B: v}
				if !used[key] {
					used[key] = true
					game := pool[key][0] // pop the first game
					pool[key] = pool[key][1:]
					thisRound = append(thisRound, game)
				}
			}
		}
		if len(thisRound) != len(teams)/2 {
			return nil, fmt.Errorf("round %d: only %d games scheduled", round, len(thisRound))
		}
		rounds = append(rounds, thisRound)
	}

	return rounds, nil
}
