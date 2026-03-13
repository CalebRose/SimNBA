package repository

import (
	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
)

type StandingsQuery struct {
	SeasonID string
	TeamID   string
}

func FindAllCollegeStandingsRecords(clauses StandingsQuery) []structs.CollegeStandings {
	db := dbprovider.GetInstance().GetDB()

	var standings []structs.CollegeStandings

	query := db.Model(&structs.CollegeStandings{})
	if clauses.SeasonID != "" {
		query = query.Where("season_id = ?", clauses.SeasonID)
	}
	if clauses.TeamID != "" {
		query = query.Where("team_id = ?", clauses.TeamID)
	}
	query.Find(&standings)

	return standings
}
