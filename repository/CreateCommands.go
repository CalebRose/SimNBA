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
	s2 := util.GenerateIntFromRange(draftee.Shooting2-3, draftee.Shooting2+3)
	s2Grade := util.GetDrafteeGrade(s2)
	s3 := util.GenerateIntFromRange(draftee.Shooting3-3, draftee.Shooting3+3)
	s3Grade := util.GetDrafteeGrade(s3)
	ft := util.GenerateIntFromRange(draftee.FreeThrow-3, draftee.FreeThrow+3)
	ftGrade := util.GetDrafteeGrade(ft)
	fn := util.GenerateIntFromRange(draftee.Finishing-3, draftee.Finishing+3)
	fnGrade := util.GetDrafteeGrade(fn)
	bw := util.GenerateIntFromRange(draftee.Ballwork-3, draftee.Ballwork+3)
	bwGrade := util.GetDrafteeGrade(bw)
	rb := util.GenerateIntFromRange(draftee.Rebounding-3, draftee.Rebounding+3)
	rbGrade := util.GetDrafteeGrade(rb)
	id := util.GenerateIntFromRange(draftee.InteriorDefense-3, draftee.InteriorDefense+3)
	idGrade := util.GetDrafteeGrade(id)
	pd := util.GenerateIntFromRange(draftee.PerimeterDefense-3, draftee.PerimeterDefense+3)
	pdGrade := util.GetDrafteeGrade(pd)
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

func CreateCollegePlayerRecord(croot structs.Recruit, db *gorm.DB, fromProgression bool) {
	cp := structs.CollegePlayer{}
	cp.MapFromRecruit(croot)
	cp.SetExpectations(util.GetPlaytimeExpectations(cp.Stars, cp.Year, cp.Overall))
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

func CreateRecruitRecord(croot structs.Recruit, db *gorm.DB) {
	// Save College Player Record
	err := db.Create(&croot).Error
	if err != nil {
		log.Panicln("Could not save new college recruit record")
	}
}

func CreatePlayerRecruitProfileRecord(cp structs.PlayerRecruitProfile, db *gorm.DB) {
	// Save College Player Record
	err := db.Create(&cp).Error
	if err != nil {
		log.Panicln("Could not save new college recruit record")
	}
}

func CreateGlobalPlayerRecord(player structs.GlobalPlayer, db *gorm.DB) {
	// Save College Player Record
	err := db.Create(&player).Error
	if err != nil {
		log.Panicln("Could not save new college recruit record")
	}
}

func CreateCollegePromiseRecord(promise structs.CollegePromise, db *gorm.DB) {
	// Save College Player Record
	err := db.Create(&promise).Error
	if err != nil {
		log.Panicln("Could not save new college recruit record")
	}
}
