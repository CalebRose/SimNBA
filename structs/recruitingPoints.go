package structs

import "github.com/jinzhu/gorm"

// RecruitingPoints - The points allocated to one player
type RecruitingPoints struct {
	gorm.Model
	Year          int
	PlayerID      int
	Points        int
	PointsSpent   int
	Scholarship   bool
	InterestLevel string
	Signed        bool
}
