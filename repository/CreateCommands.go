package repository

import (
	"log"

	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
	"gorm.io/gorm"
)

func CreateDrafteeRecord(player structs.CollegePlayer, db *gorm.DB) {
	draftee := structs.NBADraftee{}
	draftee.Map(player)
	draftee.AssignPrimeAge(util.GenerateIntFromRange(24, 30))
	// Generate Draft Grades
	s2 := util.GenerateIntFromRange(int(draftee.MidRangeShooting)-3, int(draftee.MidRangeShooting)+3)
	s2Grade := util.GetDrafteeGrade(uint8(s2))
	s3 := util.GenerateIntFromRange(int(draftee.ThreePointShooting)-3, int(draftee.ThreePointShooting)+3)
	s3Grade := util.GetDrafteeGrade(uint8(s3))
	ft := util.GenerateIntFromRange(int(draftee.FreeThrow)-3, int(draftee.FreeThrow)+3)
	ftGrade := util.GetDrafteeGrade(uint8(ft))
	fn := util.GenerateIntFromRange(int(draftee.InsideShooting)-3, int(draftee.InsideShooting)+3)
	fnGrade := util.GetDrafteeGrade(uint8(fn))
	bw := util.GenerateIntFromRange(int(draftee.Ballwork)-3, int(draftee.Ballwork)+3)
	bwGrade := util.GetDrafteeGrade(uint8(bw))
	rb := util.GenerateIntFromRange(int(draftee.Rebounding)-3, int(draftee.Rebounding)+3)
	rbGrade := util.GetDrafteeGrade(uint8(rb))
	id := util.GenerateIntFromRange(int(draftee.InteriorDefense)-3, int(draftee.InteriorDefense)+3)
	idGrade := util.GetDrafteeGrade(uint8(id))
	pd := util.GenerateIntFromRange(int(draftee.PerimeterDefense)-3, int(draftee.PerimeterDefense)+3)
	pdGrade := util.GetDrafteeGrade(uint8(pd))
	ovrVal := ((s2 + s3 + ft) / 3) + fn + bw + rb + ((id + pd) / 2)
	ovr := util.GetOverallDraftGrade(ovrVal)
	draftee.ApplyGrades(s2Grade, s3Grade, ftGrade, fnGrade, bwGrade, rbGrade, idGrade, pdGrade, ovr)
	if draftee.ProPotentialGrade == 0 {
		pot := util.GeneratePotential()
		draftee.AssignProPotentialGrade(pot)
	}

	draftee.GetNBAPotentialGrade()

	err := db.Create(&draftee).Error
	if err != nil {
		log.Panicln("Could not save historic player record!")
	}
}

func CreateInternationalDrafteeRecord(player structs.NBAPlayer, db *gorm.DB) {
	draftee := structs.NBADraftee{}
	draftee.MapInternational(player)
	draftee.AssignPrimeAge(int(player.PrimeAge))
	// Generate Draft Grades
	s2 := util.GenerateIntFromRange(int(draftee.MidRangeShooting)-3, int(draftee.MidRangeShooting)+3)
	s2Grade := util.GetDrafteeGrade(uint8(s2))
	s3 := util.GenerateIntFromRange(int(draftee.ThreePointShooting)-3, int(draftee.ThreePointShooting)+3)
	s3Grade := util.GetDrafteeGrade(uint8(s3))
	ft := util.GenerateIntFromRange(int(draftee.FreeThrow)-3, int(draftee.FreeThrow)+3)
	ftGrade := util.GetDrafteeGrade(uint8(ft))
	fn := util.GenerateIntFromRange(int(draftee.InsideShooting)-3, int(draftee.InsideShooting)+3)
	fnGrade := util.GetDrafteeGrade(uint8(fn))
	bw := util.GenerateIntFromRange(int(draftee.Ballwork)-3, int(draftee.Ballwork)+3)
	bwGrade := util.GetDrafteeGrade(uint8(bw))
	rb := util.GenerateIntFromRange(int(draftee.Rebounding)-3, int(draftee.Rebounding)+3)
	rbGrade := util.GetDrafteeGrade(uint8(rb))
	id := util.GenerateIntFromRange(int(draftee.InteriorDefense)-3, int(draftee.InteriorDefense)+3)
	idGrade := util.GetDrafteeGrade(uint8(id))
	pd := util.GenerateIntFromRange(int(draftee.PerimeterDefense)-3, int(draftee.PerimeterDefense)+3)
	pdGrade := util.GetDrafteeGrade(uint8(pd))
	ovrVal := ((s2 + s3 + ft) / 3) + fn + bw + rb + ((id + pd) / 2)
	ovr := util.GetOverallDraftGrade(ovrVal)
	draftee.ApplyGrades(s2Grade, s3Grade, ftGrade, fnGrade, bwGrade, rbGrade, idGrade, pdGrade, ovr)
	if draftee.ProPotentialGrade == 0 {
		pot := util.GeneratePotential()
		draftee.AssignProPotentialGrade(pot)
	}

	draftee.GetNBAPotentialGrade()

	err := db.Create(&draftee).Error
	if err != nil {
		log.Panicln("Could not save draftee record!")
	}
}

func CreateCollegePlayerRecord(croot structs.Recruit, db *gorm.DB, fromProgression bool) {
	cp := structs.CollegePlayer{}
	cp.MapFromRecruit(croot)
	expectations := util.GetPlaytimeExpectations(int(cp.Stars), int(cp.Year), int(cp.Overall))
	cp.SetExpectations(uint8(expectations))
	// Save College Player Record
	err := db.Create(&cp).Error
	if err != nil {
		log.Panicln("Could not save new college player record")
	}

	if fromProgression {
		DeleteCollegeRecruitRecord(croot, db)
	}
}

func CreateHistoricPlayerRecord(player structs.CollegePlayer, db *gorm.DB) {
	hcp := structs.HistoricCollegePlayer{}
	hcp.Map(player)

	err := db.Save(&hcp).Error
	if err != nil {
		log.Panicln("Could not save historic player record!")
	}
}

func CreateRetireeRecord(retiree structs.RetiredPlayer, db *gorm.DB) {
	// Save College Player Record
	retiree.Offers = nil
	retiree.WaiverOffers = nil
	retiree.Extensions = nil
	retiree.Contract = structs.NBAContract{}
	retiree.Stats = nil
	retiree.SeasonStats = structs.NBAPlayerSeasonStats{}
	err := db.Create(&retiree).Error
	if err != nil {
		log.Panicln("Could not save new college recruit record")
	}
}

func CreatePlayerRecruitProfileRecord(cp structs.RecruitPlayerProfile, db *gorm.DB) {
	// Save College Player Record
	err := db.Create(&cp).Error
	if err != nil {
		log.Panicln("Could not save new college recruit record")
	}
}

func CreateProfessionalContractRecord(contract structs.NBAContract, db *gorm.DB) {
	// Save NBA Contract Record
	err := db.Create(&contract).Error
	if err != nil {
		log.Panicln("Could not create contract record")
	}
}

func CreateCollegePromiseRecord(promise structs.CollegePromise, db *gorm.DB) {
	// Save College Player Record
	err := db.Create(&promise).Error
	if err != nil {
		log.Panicln("Could not save new college recruit record")
	}
}

func CreateISLScoutingReportRecord(report structs.ISLScoutingReport, db *gorm.DB) {
	// Save ISL Scout Report Record
	err := db.Create(&report).Error
	if err != nil {
		log.Panicln("Could not create new scout report record")
	}
}

func CreateNotification(noti structs.Notification, db *gorm.DB) {
	err := db.Create(&noti).Error
	if err != nil {
		log.Panicln("Could not create notification record!")
	}
}

func CreateNBARecordsBatch(db *gorm.DB, fds []structs.NBAMatch, batchSize int) error {
	total := len(fds)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(fds[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreatePlayerRecruitProfileRecordsBatch(db *gorm.DB, cp []structs.RecruitPlayerProfile, batchSize int) error {
	total := len(cp)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(cp[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateProContractRecordsBatch(db *gorm.DB, cp []structs.NBAContract, batchSize int) error {
	total := len(cp)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(cp[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateCollegeMatchesRecordsBatch(db *gorm.DB, fds []structs.Match, batchSize int) error {
	total := len(fds)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(fds[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateCollegePlayersRecordBatch(db *gorm.DB, fds []structs.CollegePlayer, batchSize int) error {
	total := len(fds)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(fds[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}
