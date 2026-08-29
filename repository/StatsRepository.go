package repository

import (
	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
	"gorm.io/gorm"
)

func CreateCollegeTeamStatsBatch(db *gorm.DB, fds []structs.TeamStats, batchSize int) error {
	total := len(fds)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(fds[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateNBATeamStatsBatch(db *gorm.DB, fds []structs.NBATeamStats, batchSize int) error {
	total := len(fds)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(fds[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateCollegePlayerStatsBatch(db *gorm.DB, fds []structs.CollegePlayerStats, batchSize int) error {
	total := len(fds)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(fds[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateNBAPlayerStatsBatch(db *gorm.DB, fds []structs.NBAPlayerStats, batchSize int) error {
	total := len(fds)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(fds[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateCollegePlayByPlayBatch(db *gorm.DB, fds []structs.CollegePlayByPlay, batchSize int) error {
	total := len(fds)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(fds[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateNBAPlayByPlayBatch(db *gorm.DB, fds []structs.NBAPlayByPlay, batchSize int) error {
	total := len(fds)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(fds[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func FindCollegePlayerSeasonStatRecord(playerID, SeasonID, gameType string) structs.CollegePlayerSeasonStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats structs.CollegePlayerSeasonStats

	db.Order("points desc").Where("player_id = ? AND season_id = ? AND game_type = ?", playerID, SeasonID, gameType).Find(&playerStats)

	return playerStats
}

func FindCollegePlayerSeasonStatsRecords(SeasonID, gameType string) []structs.CollegePlayerSeasonStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.CollegePlayerSeasonStats

	db.Order("points desc").Where("season_id = ? AND game_type = ?", SeasonID, gameType).Find(&playerStats)

	return playerStats
}

func FindProPlayerSeasonStatsRecords(SeasonID, gameType string) []structs.NBAPlayerSeasonStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.NBAPlayerSeasonStats

	db.Order("points desc").Where("season_id = ? AND game_type = ?", SeasonID, gameType).Find(&playerStats)

	return playerStats
}

func FindCollegePlayerGameStatsRecords(SeasonID, WeekID, GameType, GameID string) []structs.CollegePlayerStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.CollegePlayerStats

	query := db.Model(&playerStats)
	if len(SeasonID) > 0 {
		query = query.Where("season_id = ?", SeasonID)
	}

	if len(WeekID) > 0 {
		query = query.Where("week_id = ?", WeekID)
	}

	if len(GameType) > 0 {
		query = query.Where("game_type = ?", GameType)
	}

	if len(GameID) > 0 {
		query = query.Where("game_id = ?", GameID)
	}

	query.Order("points desc").Find(&playerStats)

	return playerStats
}

func FindProPlayerGameStatsRecords(SeasonID, WeekID, GameType, GameID string) []structs.NBAPlayerStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.NBAPlayerStats
	query := db.Model(&playerStats)
	if len(SeasonID) > 0 {
		query = query.Where("season_id = ?", SeasonID)
	}

	if len(WeekID) > 0 {
		query = query.Where("week_id = ?", WeekID)
	}

	if len(GameType) > 0 {
		query = query.Where("game_type = ?", GameType)
	}

	if len(GameID) > 0 {
		query = query.Where("game_id = ?", GameID)
	}

	query.Order("points desc").Find(&playerStats)

	return playerStats
}

func FindCollegeTeamSeasonStatsRecords(SeasonID, gameType string) []structs.TeamSeasonStats {
	db := dbprovider.GetInstance().GetDB()

	var teamStats []structs.TeamSeasonStats

	db.Order("points desc").Where("season_id = ? AND game_type = ?", SeasonID, gameType).Find(&teamStats)

	return teamStats
}

func FindProTeamSeasonStatsRecords(SeasonID, gameType string) []structs.NBATeamSeasonStats {
	db := dbprovider.GetInstance().GetDB()

	var teamStats []structs.NBATeamSeasonStats

	db.Order("points desc").Where("season_id = ? AND game_type = ?", SeasonID, gameType).Find(&teamStats)

	return teamStats
}

func FindCollegeTeamGameStatsRecords(SeasonID, WeekID, GameType, GameID string) []structs.TeamStats {
	db := dbprovider.GetInstance().GetDB()

	var teamStats []structs.TeamStats
	query := db.Model(&teamStats)
	if len(SeasonID) > 0 {
		query = query.Where("season_id = ?", SeasonID)
	}

	if len(WeekID) > 0 {
		query = query.Where("week_id = ?", WeekID)
	}

	if len(GameType) > 0 {
		query = query.Where("game_type = ?", GameType)
	}

	if len(GameID) > 0 {
		query = query.Where("game_id = ?", GameID)
	}

	query.Order("points desc").Find(&teamStats)

	return teamStats
}

func FindProTeamGameStatsRecords(SeasonID, WeekID, GameType, GameID string) []structs.NBATeamStats {
	db := dbprovider.GetInstance().GetDB()

	var teamStats []structs.NBATeamStats
	query := db.Model(&teamStats)
	if len(SeasonID) > 0 {
		query = query.Where("season_id = ?", SeasonID)
	}

	if len(WeekID) > 0 {
		query = query.Where("week_id = ?", WeekID)
	}

	if len(GameType) > 0 {
		query = query.Where("game_type = ?", GameType)
	}

	if len(GameID) > 0 {
		query = query.Where("game_id = ?", GameID)
	}

	query.Order("points desc").Find(&teamStats)

	return teamStats
}
