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
	ToGLeague      bool
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
	Syncs         uint8
	ToGLeague     bool
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
	n.TotalCost = offer.Year1Total + offer.Year2Total + offer.Year3Total + offer.Year4Total + offer.Year5Total
	n.IsActive = true

	// Calculate Value
	y1BonusVal := n.Year1Total * 1
	y2BonusVal := n.Year2Total * 0.9
	y3BonusVal := n.Year3Total * 0.8
	y4BonusVal := n.Year4Total * 0.7
	y5BonusVal := n.Year5Total * 0.6

	n.ContractValue = y1BonusVal + y2BonusVal + y3BonusVal + y4BonusVal + y5BonusVal
}

func (o *NBAContractOffer) IncrementSyncs() {
	o.Syncs++
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

// Table for storing Extensions for contracted players
type NBAExtensionOffer struct {
	gorm.Model
	NBAPlayerID   uint
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
	Rejections    int
	IsAccepted    bool
	IsActive      bool
	IsRejected    bool
}

func (f *NBAExtensionOffer) AssignID(id uint) {
	f.ID = id
}

func (f *NBAExtensionOffer) CalculateOffer(offer NBAContractOfferDTO) {
	f.NBAPlayerID = offer.PlayerID
	f.TeamID = offer.TeamID
	f.Team = offer.Team
	f.TotalYears = offer.TotalYears
	f.Year1Total = offer.Year1Total
	f.Year1Opt = offer.Year1Opt
	f.Year2Total = offer.Year2Total
	f.Year2Opt = offer.Year2Opt
	f.Year3Total = offer.Year3Total
	f.Year3Opt = offer.Year3Opt
	f.Year4Total = offer.Year4Total
	f.Year4Opt = offer.Year4Opt
	f.Year5Total = offer.Year5Total
	f.Year5Opt = offer.Year5Opt
	f.IsActive = true
	f.TotalCost = offer.Year1Total + offer.Year2Total + offer.Year3Total + offer.Year4Total + offer.Year5Total

	// Calculate Value
	y1BonusVal := f.Year1Total * 1
	y2BonusVal := f.Year2Total * 0.9
	y3BonusVal := f.Year3Total * 0.8
	y4BonusVal := f.Year4Total * 0.7
	y5BonusVal := f.Year5Total * 0.6

	f.ContractValue = y1BonusVal + y2BonusVal + y3BonusVal + y4BonusVal + y5BonusVal
}

func (f *NBAExtensionOffer) AcceptOffer() {
	f.IsAccepted = true
	f.CancelOffer()
}

func (f *NBAExtensionOffer) DeclineOffer(week int) {
	f.Rejections += 1
	if f.Rejections > 2 || week >= 30 {
		f.IsRejected = true
	}
}

func (f *NBAExtensionOffer) CancelOffer() {
	f.IsActive = false
}
