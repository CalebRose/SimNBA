package structs

import "github.com/jinzhu/gorm"

type PlayerStats struct {
	gorm.Model
	PlayerID           int
	MatchID            int
	SeasonID           int
	Minutes            int
	Possessions        int
	FGM                int
	FGA                int
	FGPercent          float32
	ThreePointsMade    int
	ThreePointAttempts int
	ThreePointPercent  float32
	FTM                int
	FTA                int
	FTPercent          float32
	Points             int
	TotalRebounds      int
	OffRebounds        int
	DefRebounds        int
	Assists            int
	Steals             int
	Blocks             int
	Turnovers          int
}
