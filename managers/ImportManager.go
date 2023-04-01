package managers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
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

	var teamStats []structs.TeamStats

	for _, dto := range Results.Results {
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

		homeTeam := structs.TeamStats{
			TeamID:                    uint(ht.ID),
			MatchID:                   uint(matchID),
			WeekID:                    gameRecord.WeekID,
			SeasonID:                  uint(gameRecord.SeasonID),
			Points:                    dto.TeamOne.Stats.Points,
			Possessions:               dto.TeamOne.Stats.Possessions,
			FGM:                       dto.TeamOne.Stats.FGM,
			FGA:                       dto.TeamOne.Stats.FGA,
			FGPercent:                 dto.TeamOne.Stats.FGPercent,
			ThreePointsMade:           dto.TeamOne.Stats.ThreePointsMade,
			ThreePointAttempts:        dto.TeamOne.Stats.ThreePointAttempts,
			ThreePointPercent:         dto.TeamOne.Stats.ThreePointPercent,
			FTM:                       dto.TeamOne.Stats.FTM,
			FTA:                       dto.TeamOne.Stats.FTA,
			FTPercent:                 dto.TeamOne.Stats.FTPercent,
			Rebounds:                  dto.TeamOne.Stats.Rebounds,
			OffRebounds:               dto.TeamOne.Stats.OffRebounds,
			DefRebounds:               dto.TeamOne.Stats.DefRebounds,
			Assists:                   dto.TeamOne.Stats.Assists,
			Steals:                    dto.TeamOne.Stats.Steals,
			Blocks:                    dto.TeamOne.Stats.Blocks,
			TotalTurnovers:            dto.TeamOne.Stats.TotalTurnovers,
			LargestLead:               dto.TeamOne.Stats.LargestLead,
			FirstHalfScore:            dto.TeamOne.Stats.FirstHalfScore,
			SecondHalfScore:           dto.TeamOne.Stats.SecondHalfScore,
			OvertimeScore:             dto.TeamOne.Stats.OvertimeScore,
			Fouls:                     dto.TeamOne.Stats.Fouls,
			PointsAgainst:             dto.TeamTwo.Stats.Points,
			FGMAgainst:                dto.TeamTwo.Stats.FGM,
			FGAAgainst:                dto.TeamTwo.Stats.FGA,
			FGPercentAgainst:          dto.TeamTwo.Stats.FGPercent,
			ThreePointsMadeAgainst:    dto.TeamTwo.Stats.ThreePointsMade,
			ThreePointAttemptsAgainst: dto.TeamTwo.Stats.ThreePointAttempts,
			ThreePointPercentAgainst:  dto.TeamTwo.Stats.ThreePointPercent,
			FTMAgainst:                dto.TeamTwo.Stats.FTM,
			FTAAgainst:                dto.TeamTwo.Stats.FTA,
			FTPercentAgainst:          dto.TeamTwo.Stats.FTPercent,
			ReboundsAllowed:           dto.TeamTwo.Stats.Rebounds,
			OffReboundsAllowed:        dto.TeamTwo.Stats.OffRebounds,
			DefReboundsAllowed:        dto.TeamTwo.Stats.DefRebounds,
			AssistsAllowed:            dto.TeamTwo.Stats.Assists,
			StealsAllowed:             dto.TeamTwo.Stats.Steals,
			BlocksAllowed:             dto.TeamTwo.Stats.Blocks,
			TurnoversAllowed:          dto.TeamTwo.Stats.TotalTurnovers,
		}

		teamStats = append(teamStats, homeTeam)

		awayTeamChn := make(chan structs.Team)
		go func() {
			awayTeam := GetTeamByTeamID(strconv.Itoa(dto.TeamTwo.ID))
			awayTeamChn <- awayTeam
		}()

		at := <-awayTeamChn
		close(awayTeamChn)

		awayTeam := structs.TeamStats{
			TeamID:                    at.ID,
			MatchID:                   uint(matchID),
			WeekID:                    gameRecord.WeekID,
			SeasonID:                  gameRecord.SeasonID,
			Points:                    dto.TeamTwo.Stats.Points,
			Possessions:               dto.TeamTwo.Stats.Possessions,
			FGM:                       dto.TeamTwo.Stats.FGM,
			FGA:                       dto.TeamTwo.Stats.FGA,
			FGPercent:                 dto.TeamTwo.Stats.FGPercent,
			ThreePointsMade:           dto.TeamTwo.Stats.ThreePointsMade,
			ThreePointAttempts:        dto.TeamTwo.Stats.ThreePointAttempts,
			ThreePointPercent:         dto.TeamTwo.Stats.ThreePointPercent,
			FTM:                       dto.TeamTwo.Stats.FTM,
			FTA:                       dto.TeamTwo.Stats.FTA,
			FTPercent:                 dto.TeamTwo.Stats.FTPercent,
			Rebounds:                  dto.TeamTwo.Stats.Rebounds,
			OffRebounds:               dto.TeamTwo.Stats.OffRebounds,
			DefRebounds:               dto.TeamTwo.Stats.DefRebounds,
			Assists:                   dto.TeamTwo.Stats.Assists,
			Steals:                    dto.TeamTwo.Stats.Steals,
			Blocks:                    dto.TeamTwo.Stats.Blocks,
			TotalTurnovers:            dto.TeamTwo.Stats.TotalTurnovers,
			LargestLead:               dto.TeamTwo.Stats.LargestLead,
			FirstHalfScore:            dto.TeamTwo.Stats.FirstHalfScore,
			SecondHalfScore:           dto.TeamTwo.Stats.SecondHalfScore,
			OvertimeScore:             dto.TeamTwo.Stats.OvertimeScore,
			Fouls:                     dto.TeamTwo.Stats.Fouls,
			PointsAgainst:             dto.TeamOne.Stats.Points,
			FGMAgainst:                dto.TeamOne.Stats.FGM,
			FGAAgainst:                dto.TeamOne.Stats.FGA,
			FGPercentAgainst:          dto.TeamOne.Stats.FGPercent,
			ThreePointsMadeAgainst:    dto.TeamOne.Stats.ThreePointsMade,
			ThreePointAttemptsAgainst: dto.TeamOne.Stats.ThreePointAttempts,
			ThreePointPercentAgainst:  dto.TeamOne.Stats.ThreePointPercent,
			FTMAgainst:                dto.TeamOne.Stats.FTM,
			FTAAgainst:                dto.TeamOne.Stats.FTA,
			FTPercentAgainst:          dto.TeamOne.Stats.FTPercent,
			ReboundsAllowed:           dto.TeamOne.Stats.Rebounds,
			OffReboundsAllowed:        dto.TeamOne.Stats.OffRebounds,
			DefReboundsAllowed:        dto.TeamOne.Stats.DefRebounds,
			AssistsAllowed:            dto.TeamOne.Stats.Assists,
			StealsAllowed:             dto.TeamOne.Stats.Steals,
			BlocksAllowed:             dto.TeamOne.Stats.Blocks,
			TurnoversAllowed:          dto.TeamOne.Stats.TotalTurnovers,
		}

		teamStats = append(teamStats, awayTeam)

		for _, player := range dto.RosterOne {
			id := player.ID
			collegePlayerStats := structs.CollegePlayerStats{
				CollegePlayerID:    uint(id),
				MatchID:            uint(matchID),
				SeasonID:           timestamp.SeasonID,
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
			}
			playerStats = append(playerStats, collegePlayerStats)
		}

		for _, player := range dto.RosterTwo {
			id := player.ID
			collegePlayerStats := structs.CollegePlayerStats{
				CollegePlayerID:    uint(id),
				MatchID:            uint(matchID),
				SeasonID:           timestamp.SeasonID,
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
			}

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

	for _, stats := range teamStats {
		err := db.Create(&stats).Error
		if err != nil {
			log.Panicln("Could not save team stats!")
		}
	}
	fmt.Println("Finished Import for all games")
}

func ImportNBATeamsAndArenas() {
	db := dbprovider.GetInstance().GetDB()
	path := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimNBA\\data\\NBATeams.csv"
	arenapath := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimNBA\\data\\Arenas.csv"
	nbaTeamsCSV := util.ReadCSV(path)

	for idx, row := range nbaTeamsCSV {
		if idx < 2 {
			continue
		}

		id := util.ConvertStringToInt(row[0])
		team := row[1]
		nickname := row[2]
		abbr := row[3]
		city := row[4]
		state := row[5]
		country := row[6]
		conferenceID := util.ConvertStringToInt(row[7])
		conference := row[8]
		divisionID := util.ConvertStringToInt(row[9])
		division := row[10]
		arenaID := util.ConvertStringToInt(row[11])
		arena := row[12]

		nbaTeam := structs.NBATeam{
			Team:         team,
			Nickname:     nickname,
			Abbr:         abbr,
			City:         city,
			State:        state,
			Country:      country,
			ConferenceID: uint(conferenceID),
			Conference:   conference,
			DivisionID:   uint(divisionID),
			Division:     division,
			ArenaID:      uint(arenaID),
			Arena:        arena,
			IsActive:     true,
		}

		nbaTeam.AssignID(uint(id))

		db.Create(&nbaTeam)
	}

	arenasCSV := util.ReadCSV(arenapath)

	for idx, row := range arenasCSV {
		if idx < 1 {
			continue
		}

		id := util.ConvertStringToInt(row[0])
		name := row[1]
		city := row[2]
		state := row[3]
		country := row[4]
		capacity := util.ConvertStringToInt(row[5])
		hometeam := row[6]

		arena := structs.Arena{
			ArenaName: name,
			City:      city,
			State:     state,
			Country:   country,
			Capacity:  uint(capacity),
			HomeTeam:  hometeam,
		}

		arena.AssignID(uint(id))

		db.Create(&arena)
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
