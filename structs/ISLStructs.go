package structs

import "gorm.io/gorm"

type ISLScoutingDept struct {
	ID             uint `gorm:"primaryKey"`
	TeamID         uint
	TeamLabel      string
	Prestige       uint8 // 1-5
	Resources      uint8 // 1-100. Resets at the end of each week
	IdentityPool   uint8 // Resources pooled into identify players
	ScoutingPool   uint8 // Resources spent towards scouting players in a week
	InvestingPool  uint8 // Resources spent towards investing in players
	ScoutingCount  uint8 // Total number of players successfully scouted
	IdentityBias   bool  // True will consider players in all regional countries, false will only consider their native country
	BehaviorBias   uint8 // 1 == Aggressive, 2 == normal, 3 == conservative
	Finishing      uint8 // These are modifiers that can lower the cost of scouting specific attributes, making it easier for teams to scout more efficiently.
	Shooting2      uint8 // The range of modifiers can go from 0-9.
	Shooting3      uint8 // At the end of each season, teams will be able to allocate to the modifiers. The more successful a team is, the more likely they will attain higher modifiers
	FreeThrow      uint8 // Every even season, the modifiers will randomly lower. This is to keep teams in check and keep youth development competitive
	Ballwork       uint8 // Each modifier costs 1 points to increment
	Rebounding     uint8 //
	IntDefense     uint8 //
	PerDefense     uint8 //
	Potential      uint8 //
	IdentityMod    uint8 // Can range from 1-3, lowers the cost of identifying a player. More costly (3 mod points)
	ModifierPoints uint8 // Total points available to spend on modifiers
}

func (d *ISLScoutingDept) ResetPoints() {
	d.Resources = 100
	d.InvestingPool = 0
	d.IdentityPool = 0
	d.ScoutingPool = 0
}

func (d *ISLScoutingDept) IncrementPool(pool, points uint8) {
	switch pool {
	case 1:
		d.IdentityPool += points
	case 2:
		d.ScoutingPool += points
	default:
		d.InvestingPool += points
	}
	d.Resources -= points
}

/*
	Resources for Identifying: 1 point == one player viewed.
	Resources for Scouting: 4 points == for attribute, 10 for potential
	Resources for signing: Must meet point requirement after scouting
*/

type ISLScoutingReport struct {
	gorm.Model
	PlayerID         uint
	TeamID           uint
	Finishing        bool
	Shooting2        bool
	Shooting3        bool
	FreeThrow        bool
	Rebounding       bool
	Ballwork         bool
	IntDefense       bool
	PerDefense       bool
	Potential        bool
	Overall          bool
	TotalPoints      uint8
	CurrentPoints    uint8
	PointRequirement uint // Requirement of points needed to sign on
	IsLocked         bool
	RemovedFromBoard bool
	IsSigned         bool
	Count            uint8
}

func (r *ISLScoutingReport) SetPointRequirement(points uint) {
	r.PointRequirement = points
}

func (r *ISLScoutingReport) AllocatePoints(points uint8) {
	r.CurrentPoints += uint8(points)
}

func (r *ISLScoutingReport) IncrementTotalPoints() {
	r.TotalPoints += r.CurrentPoints
	r.CurrentPoints = 0
}

func (r *ISLScoutingReport) RemovePlayerFromBoard() {
	r.RemovedFromBoard = true
	r.PointRequirement = 0
}

func (r *ISLScoutingReport) LockBoard(isSigned bool) {
	if isSigned {
		r.IsSigned = true
	}
	r.IsLocked = true
}

func (r *ISLScoutingReport) RevealAttribute(attr string) {
	switch attr {
	case "fn":
		r.Finishing = true
	case "sh2":
		r.Shooting2 = true
	case "sh3":
		r.Shooting3 = true
	case "ft":
		r.FreeThrow = true
	case "bw":
		r.Ballwork = true
	case "rb":
		r.Rebounding = true
	case "ind":
		r.IntDefense = true
	case "prd":
		r.PerDefense = true
	case "pot":
		r.Potential = true
	}
	r.Count += 1
	if r.Count > 3 {
		r.Overall = true
	}
}

type ISLCountry struct {
	Name   string
	Weight int
}
