package structs

import "github.com/jinzhu/gorm"

// Contract - The contract of which the player is obligated to
type NBAContract struct {
	gorm.Model
	PlayerID       uint
	TeamID         uint
	Team           string
	OriginalTeamID uint
	OriginalTeam   string
	PreviousTeamID uint
	PreviousTeam   string
	YearsRemaining uint
	ContractType   string
	TotalRemaining float64
	Year1Total     float64
	Year2Total     float64
	Year3Total     float64
	Year4Total     float64
	Year5Total     float64
	ContractValue  float64
	Year1Opt       bool
	Year2Opt       bool
	Year3Opt       bool
	Year4Opt       bool
	Year5Opt       bool
	IsDeadCap      bool
	IsActive       bool
	IsComplete     bool
	IsExtended     bool
	// Do we want to keep track of the year?
}

func (n *NBAContract) ProgressContract() {
	if n.YearsRemaining > 0 {
		n.YearsRemaining -= 1
	}
	n.Year1Total = n.Year2Total
	n.Year2Total = n.Year3Total
	n.Year3Total = n.Year4Total
	n.Year4Total = n.Year5Total
	n.Year5Total = 0
	n.Year1Opt = n.Year2Opt
	n.Year2Opt = n.Year3Opt
	n.Year3Opt = n.Year4Opt
	n.Year4Opt = n.Year5Opt
	n.Year5Opt = false
	n.TotalRemaining = n.Year1Total + n.Year2Total + n.Year3Total + n.Year4Total + n.Year5Total
	if n.YearsRemaining == 0 {
		n.RetireContract()
	}
}

func (n *NBAContract) DeactivateContract() {
	n.IsActive = false
}

func (n *NBAContract) RetireContract() {
	n.IsActive = false
	n.IsComplete = true
}

func (n *NBAContract) MapFromOffer(o NBAContractOffer) {
	n.PlayerID = o.PlayerID
	n.TeamID = o.TeamID
	n.Team = o.Team
	n.OriginalTeam = o.Team
	n.OriginalTeamID = o.TeamID
	n.YearsRemaining = o.TotalYears
	n.ContractType = o.ContractType
	n.TotalRemaining = o.TotalCost
	n.Year1Opt = o.Year1Opt
	n.Year1Total = o.Year1Total
	n.Year2Opt = o.Year2Opt
	n.Year2Total = o.Year2Total
	n.Year3Opt = o.Year3Opt
	n.Year3Total = o.Year3Total
	n.Year4Opt = o.Year4Opt
	n.Year4Total = o.Year4Total
	n.Year5Opt = o.Year5Opt
	n.Year5Total = o.Year5Total
	n.IsActive = true
}

func (c *NBAContract) CalculateContract() {
	// Calculate Value
	y1BonusVal := c.Year1Total * 1
	y2BonusVal := c.Year2Total * 0.9
	y3BonusVal := c.Year3Total * 0.8
	y4BonusVal := c.Year4Total * 0.7
	y5BonusVal := c.Year5Total * 0.6
	c.ContractValue = y1BonusVal + y2BonusVal + y3BonusVal + y4BonusVal + y5BonusVal
}

func (c *NBAContract) TradePlayer(TeamID uint, Team string) {
	c.PreviousTeamID = c.TeamID
	c.PreviousTeam = c.Team
	c.TeamID = TeamID
	c.Team = Team
}
