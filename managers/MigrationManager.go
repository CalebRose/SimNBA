package managers

import (
	"log"
	"strconv"
	"strings"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/repository"
	"github.com/CalebRose/SimNBA/secrets"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
	"gorm.io/gorm"
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

// parseHeightToInches converts height strings in either "6-5" or "6'9" format to total inches.
func parseHeightToInches(height string) uint8 {
	if height == "" || height == "NULL" {
		return 0
	}
	var feet, inches int
	if strings.ContainsRune(height, '\'') {
		parts := strings.SplitN(height, "'", 2)
		feet = util.ConvertStringToInt(parts[0])
		inchStr := strings.TrimRight(parts[1], "\"")
		if inchStr != "" {
			inches = util.ConvertStringToInt(inchStr)
		}
	} else {
		parts := strings.SplitN(height, "-", 2)
		if len(parts) == 2 {
			feet = util.ConvertStringToInt(parts[0])
			inches = util.ConvertStringToInt(parts[1])
		}
	}
	return uint8(feet*12 + inches)
}

// safeInt wraps ConvertStringToInt to gracefully handle empty and NULL CSV values.
func safeInt(s string) int {
	if s == "" || s == "NULL" {
		return 0
	}
	return util.ConvertStringToInt(s)
}

func Migration2026Main() {
	db := dbprovider.GetInstance().GetDB()
	hsBlob := getCrootLocations("HS")
	MigrateCollegePlayers2026(db, hsBlob)
	MigrateHistoricCollegePlayers2026(db)
	MigrateNBADraftees2026(db, hsBlob)
	MigrateNBAPlayers2026(db, hsBlob)
	MigrateNBARetirees2026(db)
}

func MigrateCollegePlayers2026(db *gorm.DB, hsBlob map[string][]structs.CrootLocation) {
	playersCSV := util.ReadCSV(secrets.GetPath()["college_players_2026"])
	collegePlayersUpload := []structs.CollegePlayer{}

	for idx, row := range playersCSV {
		if idx == 0 {
			continue
		}

		id := safeInt(row[0])
		if id == 0 {
			continue
		}
		// Column indices:
		// 0:id 1:player_id 2:first_name 3:last_name 4:position 5:archetype
		// 6:team_id 7:team_abbr 8:age 9:year 10:state 11:country 12:height
		// 13:stars 14:overall(unused) 15:finishing 16:shooting2 17:shooting3
		// 18:free_throw 19:ballwork 20:rebounding 21:interior_defense 22:perimeter_defense
		// 23:potential_grade 24:pro_potential_grade 25:stamina 26:playtime_expectations
		// 27:minutes 28:potential 29:discipline 30:injury_rating 31:is_injured
		// 32:injury_name 33:injury_type 34:weeks_of_recovery 35:injury_reserve
		// 36:personality 37:free_agency 38:recruiting_bias 39:work_ethic 40:academic_bias
		// 41:is_redshirt 42:is_redshirting 43:has_graduated 44:has_progressed
		// 45:spec_shooting2 46:spec_shooting3 47:spec_finishing 48:spec_free_throw
		// 49:spec_ballwork 50:spec_rebounding 51:spec_interior_defense 52:spec_perimeter_defense
		// 53:spec_count 63:will_declare 64:previous_team_id 65:previous_team
		// 68:transfer_status 69:transfer_likeliness 70:recruiting_bias_value
		// 71:legacy_id 72:relative_id 73:notes 74:relative_type
		cp := structs.CollegePlayer{}
		cp.ID = uint(id)
		cp.PlayerID = uint(id)
		cp.FirstName = row[2]
		cp.LastName = row[3]
		pos := row[4]
		archetype := row[5]
		switch pos {
		case "PG":
			cp.Position = "G"
			cp.Archetype = util.PickFromStringList([]string{"Point Guard", "Combo Guard"})

		case "SG":
			cp.Position = "G"
			cp.Archetype = util.PickFromStringList([]string{"Shooting Guard", "Combo Guard", "Slasher", "3-and-D"})

		case "SF":
			cp.Position = "F"
			cp.Archetype = util.PickFromStringList([]string{"Small Forward", "Two-Way", "Swingman", "Point Forward"})

		case "PF":
			cp.Position = "F"
			cp.Archetype = util.PickFromStringList([]string{"Power Forward", "Two-Way", "Point Forward"})

		case "C":
			cp.Position = "C"
			cp.Archetype = util.PickFromStringList([]string{"Stretch Center", "Post Scorer", "Rim Protector"})
		}

		if archetype == "All Around" || archetype == "All-Around" {
			cp.Archetype = "All-Around"
		}

		cp.TeamID = uint(safeInt(row[6]))
		cp.Team = row[7]
		cp.Age = uint8(safeInt(row[8]))
		year := safeInt(row[9])
		stars := safeInt(row[13])
		cp.Year = uint8(year)
		state := row[10]
		cp.Country = row[11]
		city := ""
		hs := ""
		if cp.Country == "USA" {
			if len(state) > 2 {
				state = util.GetStateKey(state)
			}
			city, hs = getCityAndHighSchool(hsBlob[state])
			cp.City = city
			cp.HighSchool = hs
		}
		cp.State = state
		cp.Height = parseHeightToInches(row[12])
		weight := 200
		if cp.Position == "G" {
			weight = util.GenerateNormalizedIntFromRange(170, 230)
		}
		if cp.Position == "F" {
			weight = util.GenerateNormalizedIntFromRange(190, 270)
		}
		if cp.Position == "C" {
			weight = util.GenerateNormalizedIntFromRange(200, 300)
		}
		cp.Weight = uint16(weight)
		cp.Stars = uint8(stars)
		// Individual attributes rescaled from old range to new 1-50 range
		cp.InsideShooting = getNewMigrationAttribute(row[15], stars, year, false)
		cp.MidRangeShooting = getNewMigrationAttribute(row[16], stars, year, false)
		cp.ThreePointShooting = getNewMigrationAttribute(row[17], stars, year, false)
		cp.FreeThrow = getNewMigrationAttribute(row[18], stars, year, false)
		cp.Ballwork = getNewMigrationAttribute(row[19], stars, year, false)
		cp.Rebounding = getNewMigrationAttribute(row[20], stars, year, false)
		cp.InteriorDefense = getNewMigrationAttribute(row[21], stars, year, false)
		cp.PerimeterDefense = getNewMigrationAttribute(row[22], stars, year, false)
		// Agility, Stealing, Blocking not present in CSV — left as zero
		agilityBase := (cp.InsideShooting + cp.MidRangeShooting + cp.ThreePointShooting + cp.Ballwork) / 4
		stealingBase := (cp.InteriorDefense + cp.PerimeterDefense + cp.Ballwork) / 3
		blockingBase := (cp.InteriorDefense + cp.PerimeterDefense + cp.Rebounding) / 3
		newAgi := util.GenerateNormalizedIntFromRange(int(agilityBase)-5, int(agilityBase)+5)
		if newAgi < 1 {
			newAgi = 1
		} else if newAgi > 50 {
			newAgi = 50
		}
		newStealing := util.GenerateNormalizedIntFromRange(int(stealingBase)-5, int(stealingBase)+5)
		if newStealing < 1 {
			newStealing = 1
		} else if newStealing > 50 {
			newStealing = 50
		}

		newBlocking := util.GenerateNormalizedIntFromRange(int(blockingBase)-5, int(blockingBase)+5)
		if newBlocking < 1 {
			newBlocking = 1
		} else if newBlocking > 50 {
			newBlocking = 50
		}

		cp.Agility = uint8(newAgi)
		cp.Stealing = uint8(newStealing)
		cp.Blocking = uint8(newBlocking)
		cp.PotentialGrade = row[23]
		cp.ProPotentialGrade = uint8(safeInt(row[24]))
		cp.Stamina = uint8(safeInt(row[25]))
		cp.PlaytimeExpectations = uint8(safeInt(row[26]))
		cp.Potential = uint8(safeInt(row[28]))
		cp.Discipline = util.RescaleDisciplineAndInjury(safeInt(row[29]))
		cp.InjuryRating = util.RescaleDisciplineAndInjury(safeInt(row[30]))
		cp.IsInjured = false
		cp.Personality = row[36]
		cp.FreeAgency = row[37]
		cp.RecruitingBias = row[38]
		cp.WorkEthic = row[39]
		cp.AcademicBias = row[40]
		cp.IsRedshirt = util.ConvertStringToBool(row[41])
		cp.IsRedshirting = util.ConvertStringToBool(row[42])
		cp.HasGraduated = util.ConvertStringToBool(row[43])
		cp.HasProgressed = util.ConvertStringToBool(row[44])
		cp.SpecMidRangeShooting = util.ConvertStringToBool(row[45])
		cp.SpecThreePointShooting = util.ConvertStringToBool(row[46])
		cp.SpecInsideShooting = util.ConvertStringToBool(row[47])
		cp.SpecFreeThrow = util.ConvertStringToBool(row[48])
		cp.SpecBallwork = util.ConvertStringToBool(row[49])
		cp.SpecRebounding = util.ConvertStringToBool(row[50])
		cp.SpecInteriorDefense = util.ConvertStringToBool(row[51])
		cp.SpecPerimeterDefense = util.ConvertStringToBool(row[52])
		// SpecAgility, SpecStealing, SpecBlocking not present in CSV — left as false
		cp.SpecAgility = util.GenerateSpecialty(cp.Position, cp.Archetype, "Agility")
		cp.SpecStealing = util.GenerateSpecialty(cp.Position, cp.Archetype, "Stealing")
		cp.SpecBlocking = util.GenerateSpecialty(cp.Position, cp.Archetype, "Blocking")
		cp.WillDeclare = util.ConvertStringToBool(row[63])
		cp.PreviousTeamID = uint(safeInt(row[64]))
		cp.PreviousTeam = row[65]
		cp.TransferStatus = uint8(safeInt(row[68]))
		cp.TransferLikeliness = row[69]
		cp.RecruitingBiasValue = row[70]
		cp.LegacyID = uint(safeInt(row[71]))
		cp.RelativeID = uint8(safeInt(row[72]))
		cp.Notes = row[73]
		cp.RelativeType = uint8(safeInt(row[74]))
		// Overall derived from rescaled attributes; will be 0 if Archetype isn't in archetypeWeights
		cp.GetOverall()
		cp.GetSpecCount()

		collegePlayersUpload = append(collegePlayersUpload, cp)
	}

	repository.CreateCollegePlayerRecordsBatch(db, collegePlayersUpload, 250)
}

func MigrateHistoricCollegePlayers2026(db *gorm.DB) {
	playersCSV := util.ReadCSV(secrets.GetPath()["historic_college_players_2026"])
	hcpUpload := []structs.HistoricCollegePlayer{}

	for idx, row := range playersCSV {
		if idx == 0 {
			continue
		}
		id := safeInt(row[0])
		if id == 0 {
			continue
		}
		// Column indices:
		// 0:id 1:first_name 2:last_name 3:position 4:age 5:year 6:team_id 7:team_abbr
		// 8:state 9:country 10:stars 11:height 12:shooting2 13:shooting3 14:finishing
		// 15:ballwork 16:rebounding 17:defense(skip) 18:potential 19:potential_grade
		// 20:pro_potential_grade 21:stamina 22:playtime_expectations 23:minutes(skip)
		// 24:overall(unused) 25:personality 26:free_agency 27:recruiting_bias 28:work_ethic
		// 29:academic_bias 30:player_id 31:is_redshirt 32:is_redshirting 33:has_graduated
		// 34:has_progressed 38:free_throw 39:interior_defense 40:perimeter_defense
		// 41:spec_shooting2 42:spec_shooting3 43:spec_finishing 44:spec_free_throw
		// 45:spec_ballwork 46:spec_rebounding 47:spec_interior_defense 48:spec_perimeter_defense
		// 49:spec_count 50:archetype 60:will_declare 61:previous_team_id 62:previous_team
		// 65:transfer_status 66:transfer_likeliness 67:recruiting_bias_value 68:legacy_id
		// 69:discipline 70:injury_rating 71:is_injured 72:injury_name 73:injury_type
		// 74:weeks_of_recovery 75:injury_reserve 76:relative_id 77:notes 78:relative_type
		hp := structs.HistoricCollegePlayer{}
		hp.ID = uint(safeInt(row[0]))
		hp.FirstName = row[1]
		hp.LastName = row[2]
		pos := row[3]
		archetype := row[50]
		switch pos {
		case "PG":
			hp.Position = "G"
			hp.Archetype = util.PickFromStringList([]string{"Point Guard", "Combo Guard"})
		case "SG":
			hp.Position = "G"
			hp.Archetype = util.PickFromStringList([]string{"Shooting Guard", "Combo Guard", "Slasher", "3-and-D"})
		case "SF":
			hp.Position = "F"
			hp.Archetype = util.PickFromStringList([]string{"Small Forward", "Two-Way", "Swingman", "Point Forward"})
		case "PF":
			hp.Position = "F"
			hp.Archetype = util.PickFromStringList([]string{"Power Forward", "Two-Way", "Point Forward"})
		case "C":
			hp.Position = "C"
			hp.Archetype = util.PickFromStringList([]string{"Stretch Center", "Post Scorer", "Rim Protector"})
		}
		if archetype == "All Around" || archetype == "All-Around" {
			hp.Archetype = "All-Around"
		}
		hp.Age = uint8(safeInt(row[4]))
		year := safeInt(row[5])
		stars := safeInt(row[10])
		hp.Year = uint8(year)
		hp.TeamID = uint(safeInt(row[6]))
		hp.Team = row[7]
		hp.State = row[8]
		hp.Country = row[9]
		hp.Stars = uint8(stars)
		hp.Height = parseHeightToInches(row[11])
		weight := 200
		switch hp.Position {
		case "G":
			weight = util.GenerateNormalizedIntFromRange(170, 230)
		case "F":
			weight = util.GenerateNormalizedIntFromRange(190, 270)
		case "C":
			weight = util.GenerateNormalizedIntFromRange(200, 300)
		}
		hp.Weight = uint16(weight)
		hp.MidRangeShooting = util.RescaleAttribute(safeInt(row[12]))
		hp.ThreePointShooting = util.RescaleAttribute(safeInt(row[13]))
		hp.InsideShooting = util.RescaleAttribute(safeInt(row[14]))
		hp.Ballwork = util.RescaleAttribute(safeInt(row[15]))
		hp.Rebounding = util.RescaleAttribute(safeInt(row[16]))
		// row[17] = defense (combined old stat, skipped)
		hp.FreeThrow = util.RescaleAttribute(safeInt(row[38]))
		hp.InteriorDefense = util.RescaleAttribute(safeInt(row[39]))
		hp.PerimeterDefense = util.RescaleAttribute(safeInt(row[40]))
		agilityBase := (hp.InsideShooting + hp.MidRangeShooting + hp.ThreePointShooting + hp.Ballwork) / 4
		stealingBase := (hp.InteriorDefense + hp.PerimeterDefense + hp.Ballwork) / 3
		blockingBase := (hp.InteriorDefense + hp.PerimeterDefense + hp.Rebounding) / 3
		newAgi := util.GenerateNormalizedIntFromRange(int(agilityBase)-5, int(agilityBase)+5)
		if newAgi < 1 {
			newAgi = 1
		} else if newAgi > 50 {
			newAgi = 50
		}
		newStealing := util.GenerateNormalizedIntFromRange(int(stealingBase)-5, int(stealingBase)+5)
		if newStealing < 1 {
			newStealing = 1
		} else if newStealing > 50 {
			newStealing = 50
		}

		newBlocking := util.GenerateNormalizedIntFromRange(int(blockingBase)-5, int(blockingBase)+5)
		if newBlocking < 1 {
			newBlocking = 1
		} else if newBlocking > 50 {
			newBlocking = 50
		}

		hp.Agility = uint8(newAgi)
		hp.Stealing = uint8(newStealing)
		hp.Blocking = uint8(newBlocking)
		hp.Potential = uint8(safeInt(row[18]))
		hp.PotentialGrade = row[19]
		hp.ProPotentialGrade = uint8(safeInt(row[20]))
		hp.Stamina = uint8(safeInt(row[21]))
		hp.PlaytimeExpectations = uint8(safeInt(row[22]))
		hp.Personality = row[25]
		hp.FreeAgency = row[26]
		hp.RecruitingBias = row[27]
		hp.WorkEthic = row[28]
		hp.AcademicBias = row[29]
		hp.PlayerID = uint(safeInt(row[30]))
		hp.IsRedshirt = util.ConvertStringToBool(row[31])
		hp.IsRedshirting = util.ConvertStringToBool(row[32])
		hp.HasGraduated = util.ConvertStringToBool(row[33])
		hp.HasProgressed = util.ConvertStringToBool(row[34])
		hp.SpecMidRangeShooting = util.ConvertStringToBool(row[41])
		hp.SpecThreePointShooting = util.ConvertStringToBool(row[42])
		hp.SpecInsideShooting = util.ConvertStringToBool(row[43])
		hp.SpecFreeThrow = util.ConvertStringToBool(row[44])
		hp.SpecBallwork = util.ConvertStringToBool(row[45])
		hp.SpecRebounding = util.ConvertStringToBool(row[46])
		hp.SpecInteriorDefense = util.ConvertStringToBool(row[47])
		hp.SpecPerimeterDefense = util.ConvertStringToBool(row[48])
		hp.SpecAgility = util.GenerateSpecialty(hp.Position, hp.Archetype, "Agility")
		hp.SpecStealing = util.GenerateSpecialty(hp.Position, hp.Archetype, "Stealing")
		hp.SpecBlocking = util.GenerateSpecialty(hp.Position, hp.Archetype, "Blocking")
		hp.WillDeclare = util.ConvertStringToBool(row[60])
		hp.PreviousTeamID = uint(safeInt(row[61]))
		hp.PreviousTeam = row[62]
		hp.TransferStatus = uint8(safeInt(row[65]))
		hp.TransferLikeliness = row[66]
		hp.RecruitingBiasValue = row[67]
		hp.LegacyID = uint(safeInt(row[68]))
		hp.Discipline = util.RescaleDisciplineAndInjury(safeInt(row[69]))
		hp.InjuryRating = util.RescaleDisciplineAndInjury(safeInt(row[70]))
		hp.IsInjured = false
		hp.RelativeID = uint8(safeInt(row[76]))
		hp.Notes = row[77]
		hp.RelativeType = uint8(safeInt(row[78]))
		hp.GetOverall()
		hp.GetSpecCount()

		hcpUpload = append(hcpUpload, hp)
	}

	repository.CreateHistoricCollegePlayerRecordsBatch(db, hcpUpload, 250)
}

func MigrateNBADraftees2026(db *gorm.DB, hsBlob map[string][]structs.CrootLocation) {
	playersCSV := util.ReadCSV(secrets.GetPath()["nba_draftees_2026"])
	nbaDUpload := []structs.NBADraftee{}

	for idx, row := range playersCSV {
		if idx == 0 {
			continue
		}
		id := safeInt(row[0])
		if id == 0 {
			continue
		}
		// Column indices:
		// 0:id 1:first_name 2:last_name 3:position 4:age 5:prime_age 6:year
		// 7:state 8:country 9:college_id 10:college 11:stars 12:height
		// 13:overall(unused) 14:finishing 15:shooting2 16:shooting3 17:free_throw
		// 18:ballwork 19:rebounding 20:interior_defense 21:perimeter_defense
		// 22:overall_grade 23:finishing_grade 24:shooting2_grade 25:shooting3_grade
		// 26:free_throw_grade 27:ballwork_grade 28:rebounding_grade
		// 29:interior_defense_grade 30:perimeter_defense_grade
		// 31:pro_potential_grade 32:potential_grade 33:draft_pick_id 34:draft_pick
		// 35:drafted_team_id 36:drafted_team_abbr 37:spec_shooting2 38:spec_shooting3
		// 39:spec_finishing 40:spec_free_throw 41:spec_ballwork 42:spec_rebounding
		// 43:spec_interior_defense 44:spec_perimeter_defense 45:spec_count 46:archetype
		// 47:prediction 48:recruiting_bias_value 49:stamina 50:playtime_expectations
		// 52:previous_team_id 53:previous_team 54:personality 55:free_agency
		// 56:recruiting_bias 57:work_ethic 58:academic_bias 59:player_id
		// 80:potential 81:discipline 82:injury_rating 83:is_injured 84:injury_name
		// 85:injury_type 86:weeks_of_recovery 87:injury_reserve
		// 88:relative_id 89:relative_type 90:notes 91:is_international
		nd := structs.NBADraftee{}
		nd.ID = uint(safeInt(row[0]))
		nd.FirstName = row[1]
		nd.LastName = row[2]
		nd.IsInternational = util.ConvertStringToBool(row[91])
		ndPos := row[3]
		ndArch := row[46]
		switch ndPos {
		case "PG":
			nd.Position = "G"
			nd.Archetype = util.PickFromStringList([]string{"Point Guard", "Combo Guard"})
		case "SG":
			nd.Position = "G"
			nd.Archetype = util.PickFromStringList([]string{"Shooting Guard", "Combo Guard", "Slasher", "3-and-D"})
		case "SF":
			nd.Position = "F"
			nd.Archetype = util.PickFromStringList([]string{"Small Forward", "Two-Way", "Swingman", "Point Forward"})
		case "PF":
			nd.Position = "F"
			nd.Archetype = util.PickFromStringList([]string{"Power Forward", "Two-Way", "Point Forward"})
		case "C":
			nd.Position = "C"
			nd.Archetype = util.PickFromStringList([]string{"Stretch Center", "Post Scorer", "Rim Protector"})
		default:
			nd.Position = ndPos
			nd.Archetype = ndArch
		}
		if ndArch == "All Around" || ndArch == "All-Around" {
			nd.Archetype = "All-Around"
		}
		nd.Age = uint8(safeInt(row[4]))
		nd.PrimeAge = safeInt(row[5])
		year := safeInt(row[6])
		stars := safeInt(row[11])
		starsForGen := stars
		if nd.IsInternational {
			starsForGen = 2
		}
		nd.Year = uint8(year)
		state := row[7]
		nd.Country = row[8]
		if nd.Country == "USA" {
			if len(state) > 2 {
				state = util.GetStateKey(state)
			}
			city, hs := getCityAndHighSchool(hsBlob[state])
			nd.City = city
			nd.HighSchool = hs
		}
		nd.State = state
		nd.CollegeID = uint(safeInt(row[9]))
		nd.College = row[10]
		nd.Stars = uint8(stars)
		nd.Height = parseHeightToInches(row[12])
		ndWeight := 200
		switch nd.Position {
		case "G":
			ndWeight = util.GenerateNormalizedIntFromRange(170, 230)
		case "F":
			ndWeight = util.GenerateNormalizedIntFromRange(190, 270)
		case "C":
			ndWeight = util.GenerateNormalizedIntFromRange(200, 300)
		}
		nd.Weight = uint16(ndWeight)
		nd.InsideShooting = getNewMigrationAttribute(row[14], starsForGen, year, true)
		nd.MidRangeShooting = getNewMigrationAttribute(row[15], starsForGen, year, true)
		nd.ThreePointShooting = getNewMigrationAttribute(row[16], starsForGen, year, true)
		nd.FreeThrow = getNewMigrationAttribute(row[17], starsForGen, year, true)
		nd.Ballwork = getNewMigrationAttribute(row[18], starsForGen, year, true)
		nd.Rebounding = getNewMigrationAttribute(row[19], starsForGen, year, true)
		nd.InteriorDefense = getNewMigrationAttribute(row[20], starsForGen, year, true)
		nd.PerimeterDefense = getNewMigrationAttribute(row[21], starsForGen, year, true)
		agilityBase := (nd.InsideShooting + nd.MidRangeShooting + nd.ThreePointShooting + nd.Ballwork) / 4
		stealingBase := (nd.InteriorDefense + nd.PerimeterDefense + nd.Ballwork) / 3
		blockingBase := (nd.InteriorDefense + nd.PerimeterDefense + nd.Rebounding) / 3
		newAgi := util.GenerateNormalizedIntFromRange(int(agilityBase)-5, int(agilityBase)+5)
		if newAgi < 1 {
			newAgi = 1
		} else if newAgi > 50 {
			newAgi = 50
		}
		newStealing := util.GenerateNormalizedIntFromRange(int(stealingBase)-5, int(stealingBase)+5)
		if newStealing < 1 {
			newStealing = 1
		} else if newStealing > 50 {
			newStealing = 50
		}

		newBlocking := util.GenerateNormalizedIntFromRange(int(blockingBase)-5, int(blockingBase)+5)
		if newBlocking < 1 {
			newBlocking = 1
		} else if newBlocking > 50 {
			newBlocking = 50
		}

		nd.Agility = uint8(newAgi)
		nd.Stealing = uint8(newStealing)
		nd.Blocking = uint8(newBlocking)
		insBaseGrade := util.GenerateNormalizedIntFromRange(int(nd.InsideShooting)-5, int(nd.InsideShooting)+5)
		nd.InsideShootingGrade = util.GetAttributeGrade(uint8(insBaseGrade), 5)
		midBaseGrade := util.GenerateNormalizedIntFromRange(int(nd.MidRangeShooting)-5, int(nd.MidRangeShooting)+5)
		nd.MidrangeShootingGrade = util.GetAttributeGrade(uint8(midBaseGrade), 5)
		threeBaseGrade := util.GenerateNormalizedIntFromRange(int(nd.ThreePointShooting)-5, int(nd.ThreePointShooting)+5)
		nd.ThreePointShootingGrade = util.GetAttributeGrade(uint8(threeBaseGrade), 5)
		ftBaseGrade := util.GenerateNormalizedIntFromRange(int(nd.FreeThrow)-5, int(nd.FreeThrow)+5)
		nd.FreeThrowGrade = util.GetAttributeGrade(uint8(ftBaseGrade), 5)
		bwBaseGrade := util.GenerateNormalizedIntFromRange(int(nd.Ballwork)-5, int(nd.Ballwork)+5)
		nd.BallworkGrade = util.GetAttributeGrade(uint8(bwBaseGrade), 5)
		rbBaseGrade := util.GenerateNormalizedIntFromRange(int(nd.Rebounding)-5, int(nd.Rebounding)+5)
		nd.ReboundingGrade = util.GetAttributeGrade(uint8(rbBaseGrade), 5)
		insDefBaseGrade := util.GenerateNormalizedIntFromRange(int(nd.InteriorDefense)-5, int(nd.InteriorDefense)+5)
		nd.InteriorDefenseGrade = util.GetAttributeGrade(uint8(insDefBaseGrade), 5)
		perDefBaseGrade := util.GenerateNormalizedIntFromRange(int(nd.PerimeterDefense)-5, int(nd.PerimeterDefense)+5)
		nd.PerimeterDefenseGrade = util.GetAttributeGrade(uint8(perDefBaseGrade), 5)
		agilityBaseGrade := util.GenerateNormalizedIntFromRange(int(nd.Agility)-5, int(nd.Agility)+5)
		nd.AgilityGrade = util.GetAttributeGrade(uint8(agilityBaseGrade), 5)
		stealingBaseGrade := util.GenerateNormalizedIntFromRange(int(nd.Stealing)-5, int(nd.Stealing)+5)
		nd.StealingGrade = util.GetAttributeGrade(uint8(stealingBaseGrade), 5)
		blockingBaseGrade := util.GenerateNormalizedIntFromRange(int(nd.Blocking)-5, int(nd.Blocking)+5)
		nd.BlockingGrade = util.GetAttributeGrade(uint8(blockingBaseGrade), 5)
		nd.ProPotentialGrade = uint8(safeInt(row[31]))
		nd.PotentialGrade = row[32]
		nd.DraftPickID = uint(safeInt(row[33]))
		nd.DraftPick = row[34] // string — formatted pick label
		nd.DraftedTeamID = uint(safeInt(row[35]))
		nd.DraftedTeam = row[36]
		nd.SpecMidRangeShooting = util.ConvertStringToBool(row[37])
		nd.SpecThreePointShooting = util.ConvertStringToBool(row[38])
		nd.SpecInsideShooting = util.ConvertStringToBool(row[39])
		nd.SpecFreeThrow = util.ConvertStringToBool(row[40])
		nd.SpecBallwork = util.ConvertStringToBool(row[41])
		nd.SpecRebounding = util.ConvertStringToBool(row[42])
		nd.SpecInteriorDefense = util.ConvertStringToBool(row[43])
		nd.SpecPerimeterDefense = util.ConvertStringToBool(row[44])
		nd.SpecAgility = util.GenerateSpecialty(nd.Position, nd.Archetype, "Agility")
		nd.SpecStealing = util.GenerateSpecialty(nd.Position, nd.Archetype, "Stealing")
		nd.SpecBlocking = util.GenerateSpecialty(nd.Position, nd.Archetype, "Blocking")
		nd.Prediction = safeInt(row[47])
		nd.RecruitingBiasValue = row[48]
		nd.Stamina = uint8(safeInt(row[49]))
		nd.PlaytimeExpectations = uint8(safeInt(row[50]))
		nd.PreviousTeamID = uint(safeInt(row[52]))
		nd.PreviousTeam = row[53]
		nd.Personality = row[54]
		nd.FreeAgency = row[55]
		nd.RecruitingBias = row[56]
		nd.WorkEthic = row[57]
		nd.AcademicBias = row[58]
		nd.PlayerID = uint(safeInt(row[59]))
		nd.Potential = uint8(safeInt(row[80]))
		nd.Discipline = util.RescaleDisciplineAndInjury(safeInt(row[81]))
		nd.InjuryRating = util.RescaleDisciplineAndInjury(safeInt(row[82]))
		nd.IsInjured = false
		nd.RelativeID = uint8(safeInt(row[88]))
		nd.RelativeType = uint8(safeInt(row[89]))
		nd.Notes = row[90]
		nd.GetOverall()
		ovrBase := util.GenerateNormalizedIntFromRange(int(nd.Overall)-5, int(nd.Overall)+5)
		nd.OverallGrade = util.GetAttributeGrade(uint8(ovrBase), 5)
		nd.GetSpecCount()

		nbaDUpload = append(nbaDUpload, nd)
	}

	repository.CreateNBADrafteesRecordsBatch(db, nbaDUpload, 250)
}

func MigrateNBAPlayers2026(db *gorm.DB, hsBlob map[string][]structs.CrootLocation) {
	playersCSV := util.ReadCSV(secrets.GetPath()["nba_players_2026"])
	nbaPUpload := []structs.NBAPlayer{}

	for idx, row := range playersCSV {
		if idx == 0 {
			continue
		}
		id := safeInt(row[0])
		if id == 0 {
			continue
		}
		// Column indices:
		// 0:id 1:player_id 2:team_id 3:team_abbr 4:first_name 5:last_name
		// 6:position 7:archetype 8:age 9:year 10:height 11:overall(unused)
		// 12:shooting2 13:shooting3 14:finishing 15:free_throw 16:ballwork
		// 17:rebounding 18:interior_defense 19:perimeter_defense
		// 20:potential_grade 21:pro_potential_grade 22:stamina 23:playtime_expectations
		// 24:personality 25:free_agency 26:recruiting_bias 27:work_ethic 28:academic_bias
		// 29:college_id 30:college 31:draft_pick_id 32:draft_pick 33:drafted_team_id
		// 34:drafted_team_abbr 35:prime_age 36:is_free_agent 37:is_nba 38:is_g_league
		// 39:is_two_way 40:is_waived 41:is_on_trade_block 42:is_first_team_anba
		// 43:is_dpoy 44:is_mvp 45:is_international 46:is_super_max_qualified
		// 47:max_requested 48:is_retiring 59:state 60:country 61:potential
		// 62:spec_shooting2 63:spec_shooting3 64:spec_finishing 65:spec_free_throw
		// 66:spec_ballwork 67:spec_rebounding 68:spec_interior_defense
		// 69:spec_perimeter_defense 70:spec_count 71:drafted_round
		// 72:previous_team_id 73:previous_team 74:is_accepting_offers 75:is_negotiating
		// 76:minimum_value 77:signing_round 78:negotiation_round 79:recruiting_bias_value
		// 80:has_progressed 81:rejections 82:discipline 83:injury_rating 84:is_injured
		// 85:injury_name 86:injury_type 87:weeks_of_recovery 88:injury_reserve
		// 89:relative_id 90:notes 91:relative_type 95:stars 97:is_int_generated 98:is_int_declared
		np := structs.NBAPlayer{}
		np.ID = uint(safeInt(row[0]))
		np.PlayerID = uint(safeInt(row[1]))
		np.TeamID = uint(safeInt(row[2]))
		np.Team = row[3]
		np.FirstName = row[4]
		np.LastName = row[5]
		npPos := row[6]
		npArch := row[7]
		switch npPos {
		case "PG":
			np.Position = "G"
			np.Archetype = util.PickFromStringList([]string{"Point Guard", "Combo Guard"})
		case "SG":
			np.Position = "G"
			np.Archetype = util.PickFromStringList([]string{"Shooting Guard", "Combo Guard", "Slasher", "3-and-D"})
		case "SF":
			np.Position = "F"
			np.Archetype = util.PickFromStringList([]string{"Small Forward", "Two-Way", "Swingman", "Point Forward"})
		case "PF":
			np.Position = "F"
			np.Archetype = util.PickFromStringList([]string{"Power Forward", "Two-Way", "Point Forward"})
		case "C":
			np.Position = "C"
			np.Archetype = util.PickFromStringList([]string{"Stretch Center", "Post Scorer", "Rim Protector"})
		default:
			np.Position = npPos
			np.Archetype = npArch
		}
		if npArch == "All Around" || npArch == "All-Around" {
			np.Archetype = "All-Around"
		}
		np.Age = uint8(safeInt(row[8]))
		np.PrimeAge = uint8(safeInt(row[35]))
		year := safeInt(row[9])
		stars := safeInt(row[95])
		np.Year = uint8(year)
		np.Height = parseHeightToInches(row[10])
		npWeight := 200
		switch np.Position {
		case "G":
			npWeight = util.GenerateNormalizedIntFromRange(170, 230)
		case "F":
			npWeight = util.GenerateNormalizedIntFromRange(190, 270)
		case "C":
			npWeight = util.GenerateNormalizedIntFromRange(200, 300)
		}
		np.Weight = uint16(npWeight)
		np.InsideShooting = getNewMigrationAttribute(row[14], stars, year, false)
		np.MidRangeShooting = getNewMigrationAttribute(row[12], stars, year, false)
		np.ThreePointShooting = getNewMigrationAttribute(row[13], stars, year, false)
		np.FreeThrow = getNewMigrationAttribute(row[15], stars, year, false)
		np.Ballwork = getNewMigrationAttribute(row[16], stars, year, false)
		np.Rebounding = getNewMigrationAttribute(row[17], stars, year, false)
		np.InteriorDefense = getNewMigrationAttribute(row[18], stars, year, false)
		np.PerimeterDefense = getNewMigrationAttribute(row[19], stars, year, false)
		agilityBase := (np.InsideShooting + np.MidRangeShooting + np.ThreePointShooting + np.Ballwork) / 4
		stealingBase := (np.InteriorDefense + np.PerimeterDefense + np.Ballwork) / 3
		blockingBase := (np.InteriorDefense + np.PerimeterDefense + np.Rebounding) / 3
		newAgi := util.GenerateNormalizedIntFromRange(int(agilityBase)-5, int(agilityBase)+5)
		if newAgi < 1 {
			newAgi = 1
		} else if newAgi > 50 {
			newAgi = 50
		}
		newStealing := util.GenerateNormalizedIntFromRange(int(stealingBase)-5, int(stealingBase)+5)
		if newStealing < 1 {
			newStealing = 1
		} else if newStealing > 50 {
			newStealing = 50
		}

		newBlocking := util.GenerateNormalizedIntFromRange(int(blockingBase)-5, int(blockingBase)+5)
		if newBlocking < 1 {
			newBlocking = 1
		} else if newBlocking > 50 {
			newBlocking = 50
		}

		np.Agility = uint8(newAgi)
		np.Stealing = uint8(newStealing)
		np.Blocking = uint8(newBlocking)
		np.PotentialGrade = row[20]
		np.ProPotentialGrade = uint8(safeInt(row[21]))
		np.Stamina = uint8(safeInt(row[22]))
		np.PlaytimeExpectations = uint8(safeInt(row[23]))
		np.Personality = row[24]
		np.FreeAgency = row[25]
		np.RecruitingBias = row[26]
		np.WorkEthic = row[27]
		np.AcademicBias = row[28]
		np.CollegeID = uint(safeInt(row[29]))
		np.College = row[30]
		np.DraftPickID = uint(safeInt(row[31]))
		np.DraftPick = uint(safeInt(row[32]))
		np.DraftedTeamID = uint(safeInt(row[33]))
		np.DraftedTeam = row[34]
		np.PrimeAge = uint8(safeInt(row[35]))
		np.IsFreeAgent = util.ConvertStringToBool(row[36])
		np.IsNBA = util.ConvertStringToBool(row[37])
		np.IsGLeague = false
		np.IsTwoWay = false
		np.IsWaived = util.ConvertStringToBool(row[40])
		np.IsOnTradeBlock = false
		np.IsFirstTeamANBA = util.ConvertStringToBool(row[42])
		np.IsDPOY = util.ConvertStringToBool(row[43])
		np.IsMVP = util.ConvertStringToBool(row[44])
		np.IsInternational = util.ConvertStringToBool(row[45])
		np.IsSuperMaxQualified = util.ConvertStringToBool(row[46])
		np.MaxRequested = util.ConvertStringToBool(row[47])
		np.IsRetiring = util.ConvertStringToBool(row[48])
		npState := row[59]
		np.Country = row[60]
		if np.Country == "USA" {
			if len(npState) > 2 {
				npState = util.GetStateKey(npState)
			}
			npCity, npHs := getCityAndHighSchool(hsBlob[npState])
			np.City = npCity
			np.HighSchool = npHs
		}
		np.State = npState
		np.Potential = uint8(safeInt(row[61]))
		np.SpecMidRangeShooting = util.ConvertStringToBool(row[62])
		np.SpecThreePointShooting = util.ConvertStringToBool(row[63])
		np.SpecInsideShooting = util.ConvertStringToBool(row[64])
		np.SpecFreeThrow = util.ConvertStringToBool(row[65])
		np.SpecBallwork = util.ConvertStringToBool(row[66])
		np.SpecRebounding = util.ConvertStringToBool(row[67])
		np.SpecInteriorDefense = util.ConvertStringToBool(row[68])
		np.SpecPerimeterDefense = util.ConvertStringToBool(row[69])
		np.SpecAgility = util.GenerateSpecialty(np.Position, np.Archetype, "Agility")
		np.SpecStealing = util.GenerateSpecialty(np.Position, np.Archetype, "Stealing")
		np.SpecBlocking = util.GenerateSpecialty(np.Position, np.Archetype, "Blocking")
		np.DraftedRound = uint(safeInt(row[71]))
		np.PreviousTeamID = uint(safeInt(row[72]))
		np.PreviousTeam = row[73]
		np.IsAcceptingOffers = util.ConvertStringToBool(row[74])
		np.IsNegotiating = util.ConvertStringToBool(row[75])
		np.MinimumValue = util.ConvertStringToFloat(row[76])
		np.SigningRound = uint(safeInt(row[77]))
		np.NegotiationRound = uint(safeInt(row[78]))
		np.RecruitingBiasValue = row[79]
		np.HasProgressed = false
		np.Rejections = int8(safeInt(row[81]))
		np.Discipline = util.RescaleDisciplineAndInjury(safeInt(row[82]))
		np.InjuryRating = util.RescaleDisciplineAndInjury(safeInt(row[83]))
		np.IsInjured = false
		np.RelativeID = uint8(safeInt(row[89]))
		np.Notes = row[90]
		np.RelativeType = uint8(safeInt(row[91]))
		np.Stars = uint8(safeInt(row[95]))
		np.IsIntGenerated = util.ConvertStringToBool(row[97])
		np.IsIntDeclared = util.ConvertStringToBool(row[98])
		np.GetOverall()
		np.GetSpecCount()

		nbaPUpload = append(nbaPUpload, np)
	}

	repository.CreateNBAPlayerRecordsBatch(db, nbaPUpload, 250)
}

func MigrateNBARetirees2026(db *gorm.DB) {
	// Note: path key is "nba_retired_players_2026" (matches secrets/paths.go)
	playersCSV := util.ReadCSV(secrets.GetPath()["nba_retired_players_2026"])
	nbaRetireesUpload := []structs.RetiredPlayer{}

	for idx, row := range playersCSV {
		if idx == 0 {
			continue
		}
		id := safeInt(row[0])
		if id == 0 {
			continue
		}
		// Column indices:
		// 0:id 1:first_name 2:last_name 3:position 4:age 5:year 6:state 7:country
		// 8:stars 9:height 10:shooting2 11:shooting3 12:finishing 13:free_throw
		// 14:ballwork 15:rebounding 16:defense(skip) 17:interior_defense 18:perimeter_defense
		// 19:potential 20:potential_grade 21:pro_potential_grade 22:stamina
		// 23:playtime_expectations 24:minutes(skip) 25:overall(unused)
		// 26:personality 27:free_agency 28:recruiting_bias 29:work_ethic 30:academic_bias
		// 31:player_id 32:team_id 33:team_abbr 34:college_id 35:college
		// 36:draft_pick_id 37:draft_pick 38:drafted_team_id 39:drafted_team_abbr
		// 40:prime_age 41:is_nba 42:max_requested 43:is_super_max_qualified
		// 44:is_free_agent 45:is_g_league 46:is_two_way 47:is_waived 48:is_on_trade_block
		// 49:is_first_team_anba 50:is_dpoy 51:is_mvp 52:is_international 53:is_retiring
		// 66:spec_shooting2 67:spec_shooting3 68:spec_finishing 69:spec_free_throw
		// 70:spec_ballwork 71:spec_rebounding 72:spec_interior_defense
		// 73:spec_perimeter_defense 74:spec_count 75:drafted_round
		// 76:previous_team_id 77:previous_team 78:is_accepting_offers 79:is_negotiating
		// 80:minimum_value 81:signing_round 82:negotiation_round 83:archetype
		// 90:recruiting_bias_value 91:has_progressed 92:rejections 93:discipline
		// 94:injury_rating 95:is_injured 96:injury_name 97:injury_type
		// 98:weeks_of_recovery 99:injury_reserve 100:relative_id 101:notes
		// 102:relative_type 103:is_int_generated 104:is_int_declared
		rp := structs.RetiredPlayer{}
		rp.ID = uint(id)
		rp.FirstName = row[1]
		rp.LastName = row[2]
		rpPos := row[3]
		rpArch := row[83]
		switch rpPos {
		case "PG":
			rp.Position = "G"
			rp.Archetype = util.PickFromStringList([]string{"Point Guard", "Combo Guard"})
		case "SG":
			rp.Position = "G"
			rp.Archetype = util.PickFromStringList([]string{"Shooting Guard", "Combo Guard", "Slasher", "3-and-D"})
		case "SF":
			rp.Position = "F"
			rp.Archetype = util.PickFromStringList([]string{"Small Forward", "Two-Way", "Swingman", "Point Forward"})
		case "PF":
			rp.Position = "F"
			rp.Archetype = util.PickFromStringList([]string{"Power Forward", "Two-Way", "Point Forward"})
		case "C":
			rp.Position = "C"
			rp.Archetype = util.PickFromStringList([]string{"Stretch Center", "Post Scorer", "Rim Protector"})
		default:
			rp.Position = rpPos
			rp.Archetype = rpArch
		}
		if rpArch == "All Around" || rpArch == "All-Around" {
			rp.Archetype = "All-Around"
		}
		rp.Age = uint8(safeInt(row[4]))
		rp.PrimeAge = uint8(safeInt(row[40]))
		rp.Year = uint8(safeInt(row[5]))
		rp.State = row[6]
		rp.Country = row[7]
		rp.Stars = uint8(safeInt(row[8]))
		rp.Height = parseHeightToInches(row[9])
		rpWeight := 200
		switch rp.Position {
		case "G":
			rpWeight = util.GenerateNormalizedIntFromRange(170, 230)
		case "F":
			rpWeight = util.GenerateNormalizedIntFromRange(190, 270)
		case "C":
			rpWeight = util.GenerateNormalizedIntFromRange(200, 300)
		}
		rp.Weight = uint16(rpWeight)
		rp.MidRangeShooting = util.RescaleAttribute(safeInt(row[10]))
		rp.ThreePointShooting = util.RescaleAttribute(safeInt(row[11]))
		rp.InsideShooting = util.RescaleAttribute(safeInt(row[12]))
		rp.FreeThrow = util.RescaleAttribute(safeInt(row[13]))
		rp.Ballwork = util.RescaleAttribute(safeInt(row[14]))
		rp.Rebounding = util.RescaleAttribute(safeInt(row[15]))
		// row[16] = defense (combined old stat, skipped)
		rp.InteriorDefense = util.RescaleAttribute(safeInt(row[17]))
		rp.PerimeterDefense = util.RescaleAttribute(safeInt(row[18]))
		agilityBase := (rp.InsideShooting + rp.MidRangeShooting + rp.ThreePointShooting + rp.Ballwork) / 4
		stealingBase := (rp.InteriorDefense + rp.PerimeterDefense + rp.Ballwork) / 3
		blockingBase := (rp.InteriorDefense + rp.PerimeterDefense + rp.Rebounding) / 3
		newAgi := util.GenerateNormalizedIntFromRange(int(agilityBase)-5, int(agilityBase)+5)
		if newAgi < 1 {
			newAgi = 1
		} else if newAgi > 50 {
			newAgi = 50
		}
		newStealing := util.GenerateNormalizedIntFromRange(int(stealingBase)-5, int(stealingBase)+5)
		if newStealing < 1 {
			newStealing = 1
		} else if newStealing > 50 {
			newStealing = 50
		}

		newBlocking := util.GenerateNormalizedIntFromRange(int(blockingBase)-5, int(blockingBase)+5)
		if newBlocking < 1 {
			newBlocking = 1
		} else if newBlocking > 50 {
			newBlocking = 50
		}

		rp.Agility = uint8(newAgi)
		rp.Stealing = uint8(newStealing)
		rp.Blocking = uint8(newBlocking)
		rp.Potential = uint8(safeInt(row[19]))
		rp.PotentialGrade = row[20]
		rp.ProPotentialGrade = uint8(safeInt(row[21]))
		rp.Stamina = uint8(safeInt(row[22]))
		rp.PlaytimeExpectations = uint8(safeInt(row[23]))
		rp.Personality = row[26]
		rp.FreeAgency = row[27]
		rp.RecruitingBias = row[28]
		rp.WorkEthic = row[29]
		rp.AcademicBias = row[30]
		rp.PlayerID = uint(safeInt(row[31]))
		rp.TeamID = uint(safeInt(row[32]))
		rp.Team = row[33]
		rp.CollegeID = uint(safeInt(row[34]))
		rp.College = row[35]
		rp.DraftPickID = uint(safeInt(row[36]))
		rp.DraftPick = uint(safeInt(row[37]))
		rp.DraftedTeamID = uint(safeInt(row[38]))
		rp.DraftedTeam = row[39]
		rp.PrimeAge = uint8(safeInt(row[40]))
		rp.IsNBA = util.ConvertStringToBool(row[41])
		rp.MaxRequested = util.ConvertStringToBool(row[42])
		rp.IsSuperMaxQualified = util.ConvertStringToBool(row[43])
		rp.IsFreeAgent = util.ConvertStringToBool(row[44])
		rp.IsGLeague = util.ConvertStringToBool(row[45])
		rp.IsTwoWay = util.ConvertStringToBool(row[46])
		rp.IsWaived = util.ConvertStringToBool(row[47])
		rp.IsOnTradeBlock = util.ConvertStringToBool(row[48])
		rp.IsFirstTeamANBA = util.ConvertStringToBool(row[49])
		rp.IsDPOY = util.ConvertStringToBool(row[50])
		rp.IsMVP = util.ConvertStringToBool(row[51])
		rp.IsInternational = util.ConvertStringToBool(row[52])
		rp.IsRetiring = util.ConvertStringToBool(row[53])
		rp.SpecMidRangeShooting = util.ConvertStringToBool(row[66])
		rp.SpecThreePointShooting = util.ConvertStringToBool(row[67])
		rp.SpecInsideShooting = util.ConvertStringToBool(row[68])
		rp.SpecFreeThrow = util.ConvertStringToBool(row[69])
		rp.SpecBallwork = util.ConvertStringToBool(row[70])
		rp.SpecRebounding = util.ConvertStringToBool(row[71])
		rp.SpecInteriorDefense = util.ConvertStringToBool(row[72])
		rp.SpecPerimeterDefense = util.ConvertStringToBool(row[73])
		rp.SpecAgility = util.GenerateSpecialty(rp.Position, rp.Archetype, "Agility")
		rp.SpecStealing = util.GenerateSpecialty(rp.Position, rp.Archetype, "Stealing")
		rp.SpecBlocking = util.GenerateSpecialty(rp.Position, rp.Archetype, "Blocking")
		rp.DraftedRound = uint(safeInt(row[75]))
		rp.PreviousTeamID = uint(safeInt(row[76]))
		rp.PreviousTeam = row[77]
		rp.IsAcceptingOffers = util.ConvertStringToBool(row[78])
		rp.IsNegotiating = util.ConvertStringToBool(row[79])
		rp.MinimumValue = util.ConvertStringToFloat(row[80])
		rp.SigningRound = uint(safeInt(row[81]))
		rp.NegotiationRound = uint(safeInt(row[82]))
		rp.Archetype = row[83]
		rp.RecruitingBiasValue = row[90]
		rp.HasProgressed = util.ConvertStringToBool(row[91])
		rp.Rejections = int8(safeInt(row[92]))
		rp.Discipline = util.RescaleDisciplineAndInjury(safeInt(row[93]))
		rp.InjuryRating = util.RescaleDisciplineAndInjury(safeInt(row[94]))
		rp.IsInjured = false
		rp.RelativeID = uint8(safeInt(row[100]))
		rp.Notes = row[101]
		rp.RelativeType = uint8(safeInt(row[102]))
		rp.IsIntGenerated = util.ConvertStringToBool(row[103])
		rp.IsIntDeclared = util.ConvertStringToBool(row[104])
		rp.GetOverall()
		rp.GetSpecCount()

		nbaRetireesUpload = append(nbaRetireesUpload, rp)
	}

	repository.CreateNBARetiredPlayerRecordsBatch(db, nbaRetireesUpload, 250)
}

func getNewMigrationAttribute(value string, stars, year int, isDraft bool) uint8 {
	rescaledValue := util.RescaleAttribute(safeInt(value))
	// AdditivBonus: +1 per star above 2, +1 per year above 1
	starsMod := int8(3)
	yearMod := int8(2)
	if isDraft {
		starsMod = 2
		yearMod = 1
	}
	bonus := int(max(0, int8(stars)-starsMod)) + int(max(0, int8(year)-yearMod))
	boosted := int(rescaledValue) + bonus
	if boosted > 50 {
		boosted = 50
	}
	floated := float64(boosted)
	buffed := floated * util.GenerateFloatFromRange(1.08, 1.28)
	if buffed > 50 {
		buffed = 50
	}
	return uint8(buffed)
}
