package managers

import (
	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
)

func GetNBAUserByUsername(username string) structs.NBAUser {
	db := dbprovider.GetInstance().GetDB()

	var user structs.NBAUser

	err := db.Where("username = ?", username).Find(&user).Error
	if err != nil || user.ID == 0 {
		user = structs.NBAUser{
			Username:    username,
			TeamID:      0,
			TotalWins:   0,
			TotalLosses: 0,
			IsActive:    true,
		}
	}

	return user
}

func GetAllCollegeCoaches() []structs.CollegeCoach {
	db := dbprovider.GetInstance().GetDB()

	coaches := []structs.CollegeCoach{}

	db.Find(&coaches)

	return coaches
}

func GetAllActiveCollegeCoaches() []structs.CollegeCoach {
	db := dbprovider.GetInstance().GetDB()

	coaches := []structs.CollegeCoach{}

	db.Where("is_retired = ? and team_id > ?", false, "0").Find(&coaches)

	return coaches
}

func GetActiveCollegeCoachMap() map[uint]structs.CollegeCoach {
	coachMap := make(map[uint]structs.CollegeCoach)

	coaches := GetAllCollegeCoaches()

	for _, coach := range coaches {
		if coach.IsRetired || coach.TeamID == 0 {
			continue
		}
		coachMap[coach.TeamID] = coach
	}

	return coachMap
}
