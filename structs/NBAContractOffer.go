package structs

import "github.com/jinzhu/gorm"

type NBAContractOfferDTO struct {
	ID             uint
	PlayerID       uint
	TeamID         uint
	TeamAbbr       string
	OriginalTeamID uint
	OriginalTeam   string
	PreviousTeamID uint
	PreviousTeam   string
	SeasonID       uint
	Team           string
	TotalYears     uint
	ContractType   string
	TotalCost      float64
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
	IsAccepted     bool
	IsRejected     bool
	// Do we want to kep track of the year?
}

type NBAContractOffer struct {
	gorm.Model
	PlayerID      uint
	TeamID        uint
	SeasonID      uint
	Team          string
	TotalYears    uint
	ContractType  string
	TotalCost     float64
	Year1Total    float64
	Year2Total    float64
	Year3Total    float64
	Year4Total    float64
	Year5Total    float64
	ContractValue float64
	Year1Opt      bool
	Year2Opt      bool
	Year3Opt      bool
	Year4Opt      bool
	Year5Opt      bool
	IsAccepted    bool
	IsRejected    bool
	IsActive      bool
	// Do we want to kep track of the year?
}

func (o *NBAContractOffer) AssignID(id uint) {
	o.ID = id
}

func (n *NBAContractOffer) CalculateOffer(offer NBAContractOfferDTO) {
	n.PlayerID = offer.PlayerID
	n.TeamID = offer.TeamID
	n.Team = offer.Team
	n.TotalYears = offer.TotalYears
	n.Year1Total = offer.Year1Total
	n.Year2Total = offer.Year2Total
	n.Year3Total = offer.Year3Total
	n.Year4Total = offer.Year4Total
	n.Year5Total = offer.Year5Total
	n.Year1Opt = offer.Year1Opt
	n.Year2Opt = offer.Year2Opt
	n.Year3Opt = offer.Year3Opt
	n.Year4Opt = offer.Year4Opt
	n.Year5Opt = offer.Year5Opt
	n.IsActive = true

	// Calculate Value
	y1BonusVal := n.Year1Total * 1
	y2BonusVal := n.Year2Total * 0.9
	y3BonusVal := n.Year3Total * 0.8
	y4BonusVal := n.Year4Total * 0.7
	y5BonusVal := n.Year5Total * 0.6

	n.ContractValue = y1BonusVal + y2BonusVal + y3BonusVal + y4BonusVal + y5BonusVal
}

func (o *NBAContractOffer) AcceptOffer() {
	o.IsAccepted = true
}

func (o *NBAContractOffer) RejectOffer() {
	o.IsRejected = false
}

func (o *NBAContractOffer) CancelOffer() {
	o.IsActive = false
}

// Sorting Funcs
type ByContractValue []NBAContractOffer

func (fo ByContractValue) Len() int      { return len(fo) }
func (fo ByContractValue) Swap(i, j int) { fo[i], fo[j] = fo[j], fo[i] }
func (fo ByContractValue) Less(i, j int) bool {
	return fo[i].ContractValue > fo[j].ContractValue
}
