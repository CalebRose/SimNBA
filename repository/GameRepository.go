package repository

import (
	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
)

type GameQuery struct {
	TeamID                string
	SeasonID              string
	WeekID                string
	ExcludeCompletedGames bool
	IsPlayoffs            bool
	IsInternational       bool // NBA matches only
}

func FindCollegeMatchRecords(clauses GameQuery) []structs.Match {
	db := dbprovider.GetInstance().GetDB()

	var teamMatches []structs.Match

	query := db.Model(&structs.Match{})

	if clauses.SeasonID != "" {
		query = query.Where("season_id = ?", clauses.SeasonID)
	}
	if clauses.WeekID != "" {
		query = query.Where("week_id = ?", clauses.WeekID)
	}
	if clauses.TeamID != "" {
		query = query.Where("home_team_id = ? OR away_team_id = ?", clauses.TeamID, clauses.TeamID)
	}
	if clauses.ExcludeCompletedGames {
		query = query.Where("game_complete = ?", false)
	}
	if clauses.IsPlayoffs {
		query = query.Where("is_playoff_game = ?", true)
	}
	if clauses.IsInternational {
		query = query.Where("is_nba_match = ?", true)
	}

	query.Find(&teamMatches)

	return teamMatches
}
