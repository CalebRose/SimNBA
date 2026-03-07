package managers

import (
	"log"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/repository"
	"github.com/CalebRose/SimNBA/secrets"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
)

func CleanNBAPlayerTables() {
	db := dbprovider.GetInstance().GetDB()

	nbaPlayers := GetAllNBAPlayers()
	retiredPlayers := GetAllRetiredPlayers()

	for _, n := range nbaPlayers {
		id := strconv.Itoa(int(n.ID))
		contracts := GetContractsByPlayerID(id)
		hasActiveContract := false
		activeContractCount := 0
		var ac structs.NBAContract // Active Contract

		for _, c := range contracts {
			if c.IsActive {
				hasActiveContract = true
				activeContractCount++
				if ac.ID == 0 {
					ac = c
				}
			}
			if activeContractCount > 1 {
				c.RetireContract()
				db.Delete(&c)
			}
		}

		if !n.IsFreeAgent && hasActiveContract {
			continue
		}

		if n.IsFreeAgent && !hasActiveContract {
			continue
		}

		// If an nba player is not a free agent and they have no contracts
		if !n.IsFreeAgent && n.TeamID != 0 && (len(contracts) == 0 || !hasActiveContract) {
			n.BecomeFreeAgent()
			db.Save(&n)
			continue
		}
		if (n.IsFreeAgent || n.TeamID == 0 || n.Team == "FA") && hasActiveContract {
			n.SignWithTeam(ac.TeamID, ac.Team, false, 0)
			db.Save(&n)
			continue
		}
	}

	for _, r := range retiredPlayers {
		id := strconv.Itoa(int(r.ID))
		contracts := GetContractsByPlayerID(id)
		for _, c := range contracts {
			if !c.IsComplete || c.IsActive {
				c.RetireContract()
				db.Delete(&c)
			}
		}
	}
}

func MigrateRecruits() {
	db := dbprovider.GetInstance().GetDB()

	croots := GetAllRecruitRecords()

	for _, croot := range croots {
		// Convert to College Player Record
		cp := structs.CollegePlayer{}
		cp.MapFromRecruit(croot)

		// Save College Player Record
		err := db.Create(&cp).Error
		if err != nil {
			log.Panicln("Could not save new college player record")
		}

		// Delete Recruit Record
		db.Delete(&croot)
	}
}

func ProgressContractsByOneYear() {
	db := dbprovider.GetInstance().GetDB()

	nbaPlayers := GetAllNBAPlayers()

	for _, n := range nbaPlayers {
		if n.IsFreeAgent {
			continue
		}
		id := strconv.Itoa(int(n.ID))
		contract := GetContractByPlayerID(id)

		contract.ProgressContract()
		if !contract.IsActive || contract.IsComplete {
			n.BecomeFreeAgent()
			db.Save(&n)
		}
		db.Save(&contract)
	}
}

func MigrateNewAIRecruitingValues() {
	db := dbprovider.GetInstance().GetDB()

	path := secrets.GetPath()["ai"]
	teams := util.ReadCSV(path)

	for idx, row := range teams {
		if idx == 0 {
			continue
		}

		id := row[0]
		aiValue := row[7]
		attr1 := row[8]
		attr2 := row[9]

		attributeList := []string{"Finishing", "FreeThrow", "Shooting2", "Shooting3", "Ballwork", "Rebounding", "InteriorDefense", "PerimeterDefense"}

		for len(attr1) == 0 {
			attr1 = util.PickFromStringList(attributeList)
		}

		for len(attr2) == 0 || attr1 == attr2 {
			attr2 = util.PickFromStringList(attributeList)
		}

		teamProfile := GetOnlyTeamRecruitingProfileByTeamID(id)
		teamProfile.SetNewBehaviors(aiValue, attr1, attr2)

		db.Save(&teamProfile)
	}
}

func MigrateMissingRecruits() {
	db := dbprovider.GetInstance().GetDB()
	recruits := GetAllRecruitRecords()

	for _, croot := range recruits {
		// Check for ID conflicts and resolve them before creating college player records
		resolvedCroot := resolveRecruitIDConflict(croot, db)
		repository.CreateCollegePlayerRecord(resolvedCroot, db, true)
	}
}
