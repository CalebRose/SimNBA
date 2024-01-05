package structs

import (
	"github.com/jinzhu/gorm"
)

type CollegePlayer struct {
	gorm.Model
	BasePlayer
	PlayerID      uint
	TeamID        uint
	TeamAbbr      string
	IsRedshirt    bool
	IsRedshirting bool
	HasGraduated  bool
	HasProgressed bool
	WillDeclare   bool
	WillTransfer  bool
	LeavingStatus string
	Stats         []CollegePlayerStats     `gorm:"foreignKey:CollegePlayerID"`
	SeasonStats   CollegePlayerSeasonStats `gorm:"foreignKey:CollegePlayerID"`
}

func (c *CollegePlayer) SetRedshirtingStatus() {
	c.IsRedshirting = true
	if c.WillDeclare {
		c.WillDeclare = false
	}
}

func (c *CollegePlayer) UpdateMinutes(newMinutes int) {
	c.Minutes = newMinutes
}

func (c *CollegePlayer) SetID(id uint) {
	c.ID = id
}

func (cp *CollegePlayer) Progress(attr CollegePlayerProgressions) {
	cp.Age++
	cp.Year++
	cp.Ballwork += attr.Ballwork
	cp.Shooting2 += attr.Shooting2
	cp.Shooting3 += attr.Shooting3
	cp.FreeThrow += attr.FreeThrow
	cp.Finishing += attr.Finishing
	cp.InteriorDefense += attr.InteriorDefense
	cp.PerimeterDefense += attr.PerimeterDefense
	cp.Rebounding += attr.Rebounding
	cp.PotentialGrade = attr.PotentialGrade
	cp.Stamina = attr.Stamina
	cp.Overall = (int((cp.Shooting2 + cp.Shooting3 + cp.FreeThrow) / 3)) + cp.Finishing + cp.Ballwork + cp.Rebounding + int((cp.InteriorDefense+cp.PerimeterDefense)/2)
	cp.HasProgressed = true
}

func (cp *CollegePlayer) MapFromRecruit(r Recruit) {
	cp.ID = r.ID
	cp.TeamID = r.TeamID
	cp.TeamAbbr = r.TeamAbbr
	cp.PlayerID = r.PlayerID
	cp.State = r.State
	cp.Country = r.Country
	cp.Year = 1
	cp.IsRedshirt = false
	cp.IsRedshirting = false
	cp.HasGraduated = false
	cp.HasProgressed = true
	cp.Age = 19
	cp.FirstName = r.FirstName
	cp.LastName = r.LastName
	cp.Position = r.Position
	cp.Archetype = r.Archetype
	cp.Height = r.Height
	cp.Stars = r.Stars
	cp.Overall = r.Overall
	cp.Shooting2 = r.Shooting2
	cp.Shooting3 = r.Shooting3
	cp.FreeThrow = r.FreeThrow
	cp.Finishing = r.Finishing
	cp.Ballwork = r.Ballwork
	cp.Rebounding = r.Rebounding
	cp.InteriorDefense = r.InteriorDefense
	cp.PerimeterDefense = r.PerimeterDefense
	cp.Stamina = r.Stamina
	cp.Potential = r.Potential
	cp.ProPotentialGrade = r.ProPotentialGrade
	cp.PotentialGrade = r.PotentialGrade
	cp.FreeAgency = r.FreeAgency
	cp.Personality = r.Personality
	cp.RecruitingBias = r.RecruitingBias
	cp.WorkEthic = r.WorkEthic
	cp.AcademicBias = r.AcademicBias
	cp.SpecBallwork = r.SpecBallwork
	cp.SpecFinishing = r.SpecFinishing
	cp.SpecFreeThrow = r.SpecFreeThrow
	cp.SpecCount = r.SpecCount
	cp.SpecInteriorDefense = r.SpecInteriorDefense
	cp.SpecPerimeterDefense = r.SpecPerimeterDefense
	cp.SpecRebounding = r.SpecRebounding
	cp.SpecShooting2 = r.SpecShooting2
	cp.SpecShooting3 = r.SpecShooting3
}

func (cp *CollegePlayer) GraduatePlayer() {
	cp.HasGraduated = true
}

func (p *CollegePlayer) SetRedshirtStatus() {
	if p.IsRedshirting {
		p.IsRedshirting = false
		p.IsRedshirt = true
	}
}

func (p *CollegePlayer) SetProgressionStatus() {
	p.HasProgressed = true
}

func (p *CollegePlayer) FixAge() {
	p.Year++
	p.Age = p.Year + 18
	if p.IsRedshirt {
		p.Age++
	}
}

func (p *CollegePlayer) SetMinutes(val int) {
	p.Minutes = val
}

func (p *CollegePlayer) SetNewAttributes(ft int, id int, pd int) {
	p.FreeThrow = ft
	p.InteriorDefense = id
	p.PerimeterDefense = pd
}

func (b *CollegePlayer) SetNewPosition(pos string) {
	b.Position = pos
}

func (b *CollegePlayer) SetDeclarationStatus() {
	b.WillDeclare = true
}

func (b *CollegePlayer) StayHome() {
	b.WillDeclare = false
}

func (b *CollegePlayer) SetTransferStatus() {
	b.WillTransfer = true
}

func (b *CollegePlayer) TransferOut() {
	b.WillTransfer = false
	b.PreviousTeam = b.TeamAbbr
	b.PreviousTeamID = b.TeamID
	b.TeamAbbr = ""
	b.TeamID = 0
}

func (b *CollegePlayer) SetLeavingStatus(status string) {
	b.LeavingStatus = status
}

// Sorting Funcs
type ByPlayerOverall []CollegePlayer

func (cp ByPlayerOverall) Len() int      { return len(cp) }
func (cp ByPlayerOverall) Swap(i, j int) { cp[i], cp[j] = cp[j], cp[i] }
func (cp ByPlayerOverall) Less(i, j int) bool {
	return cp[i].Overall > cp[j].Overall
}
