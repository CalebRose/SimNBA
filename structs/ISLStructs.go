package structs

import "gorm.io/gorm"

type ISLScoutingDept struct {
	ID             uint
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

/*
	Resources for Identifying: 1 point == one player viewed.
	Resources for Scouting: 4 points == for attribute, 10 for potential
	Resources for signing: Must meet point requirement after scouting
*/

type ISLScoutingReport struct {
	gorm.Model
	PlayerID         uint
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
	PointRequirement uint // Requirement of points needed to sign on
	IsLocked         bool
	RemovedFromBoard bool
	IsSigned         bool
}

type ISLCountry struct {
	Name   string
	Weight int
}
