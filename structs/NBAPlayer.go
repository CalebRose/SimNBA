package structs

import "github.com/jinzhu/gorm"

type NBAPlayer struct {
	gorm.Model
	BasePlayer
	CollegeID           uint
	College             string
	DraftPickID         uint
	DraftedRound        uint
	DraftPick           uint
	DraftedTeamID       uint
	DraftedTeam         string
	IsNBA               bool
	MaxRequested        bool
	IsSuperMaxQualified bool
	IsFreeAgent         bool
	IsGLeague           bool
	IsTwoWay            bool
	IsWaived            bool
	IsOnTradeBlock      bool
	IsFirstTeamANBA     bool
	IsDPOY              bool
	IsMVP               bool
	IsInternational     bool
	IsIntGenerated      bool
	IsIntDeclared       bool
	IsRetiring          bool
	IsAcceptingOffers   bool
	IsNegotiating       bool
	MinimumValue        float64
	SigningRound        uint
	NegotiationRound    uint
	Rejections          int8
	HasProgressed       bool
	Offers              []NBAContractOffer   `gorm:"foreignKey:PlayerID"`
	WaiverOffers        []NBAWaiverOffer     `gorm:"foreignKey:PlayerID"`
	Extensions          []NBAExtensionOffer  `gorm:"foreignKey:NBAPlayerID"`
	Contract            NBAContract          `gorm:"foreignKey:PlayerID"`
	Stats               []NBAPlayerStats     `gorm:"foreignKey:NBAPlayerID"`
	SeasonStats         NBAPlayerSeasonStats `gorm:"foreignKey:NBAPlayerID"`
}

type ByTotalContract []NBAPlayer

func (rp ByTotalContract) Len() int      { return len(rp) }
func (rp ByTotalContract) Swap(i, j int) { rp[i], rp[j] = rp[j], rp[i] }
func (rp ByTotalContract) Less(i, j int) bool {
	p1 := rp[i].Contract
	p2 := rp[j].Contract
	p1Total := p1.Year1Total + p1.Year2Total + p1.Year3Total + p1.Year4Total + p1.Year5Total
	p2Total := p2.Year1Total + p2.Year2Total + p2.Year3Total + p2.Year4Total + p2.Year5Total
	return p1Total > p2Total
}

func (n *NBAPlayer) SetID(id uint) {
	n.ID = id
}

func (n *NBAPlayer) SetRetiringStatus() {
	n.IsRetiring = true
}

func (n *NBAPlayer) BecomeUDFA() {
	n.Team = "FA"
	n.TeamID = 0
	n.IsFreeAgent = true
	n.IsOnTradeBlock = false
	n.IsGLeague = false
	n.IsTwoWay = false
	n.IsAcceptingOffers = true
	n.IsNBA = true
	n.IsNegotiating = false
	n.IsAcceptingOffers = true
	n.MinimumValue = 0.7
}

func (n *NBAPlayer) BecomeFreeAgent() {
	n.Team = "FA"
	n.TeamID = 0
	n.IsFreeAgent = true
	n.IsOnTradeBlock = false
	n.IsGLeague = false
	n.IsTwoWay = false
	n.IsAcceptingOffers = true
	n.AssignMinimumContractValue(0)
}

func (n *NBAPlayer) BecomeInternationalDraftee() {
	n.CollegeID = n.TeamID
	n.College = n.Team
	n.Team = "DRAFT"
	n.TeamID = 0
	n.IsFreeAgent = false
	n.IsOnTradeBlock = false
	n.IsGLeague = false
	n.IsTwoWay = false
	n.IsAcceptingOffers = false
	n.IsIntDeclared = true
}

func (n *NBAPlayer) DraftInternationalPlayer(pickID, round, number, teamID uint, team string) {
	n.DraftPickID = pickID
	n.DraftedRound = round
	n.DraftPick = number
	n.DraftedTeam = team
	n.DraftedTeamID = teamID
	n.IsNBA = true
	n.Team = team
	n.TeamID = teamID
	n.IsFreeAgent = false
	n.IsWaived = false
	n.IsGLeague = false
	n.IsTwoWay = false
	n.IsAcceptingOffers = false
	n.IsNegotiating = false
}

func (n *NBAPlayer) SignWithTeam(teamID uint, team string, isFAorExt bool, minValue float64) {
	n.Team = team
	n.TeamID = teamID
	n.IsFreeAgent = false
	n.IsWaived = false
	n.IsGLeague = false
	n.IsTwoWay = false
	n.IsAcceptingOffers = false
	n.IsNegotiating = false
	n.IsInternational = teamID < 33
}

func (n *NBAPlayer) Progress(p NBAPlayerProgressions) {
	n.HasProgressed = true
	n.MidRangeShooting += p.Shooting2
	n.ThreePointShooting += p.Shooting3
	n.FreeThrow += p.FreeThrow
	n.Ballwork += p.Ballwork
	n.InsideShooting += p.Finishing
	n.Rebounding += p.Rebounding
	n.InteriorDefense += p.InteriorDefense
	n.PerimeterDefense += p.PerimeterDefense
	n.Overall = (uint8((n.MidRangeShooting + n.ThreePointShooting + n.FreeThrow) / 3)) + n.InsideShooting + n.Ballwork + n.Rebounding + uint8((n.InteriorDefense+n.PerimeterDefense)/2)
	n.Age = p.Age
	n.Stamina = p.Stamina
	if n.Stamina < 1 {
		n.Stamina = 1
	}
	n.Year++
}
func (n *NBAPlayer) QualifyForSuperMax() {
	n.IsSuperMaxQualified = true
}

func (n *NBAPlayer) QualifiesForMax() {
	n.MaxRequested = true
}

func (n *NBAPlayer) DoesNotQualify() {
	n.MaxRequested = false
	n.IsSuperMaxQualified = false
}

func (n *NBAPlayer) ToggleIsNegotiating() {
	n.IsNegotiating = true
	n.IsAcceptingOffers = false
}

func (n *NBAPlayer) WaitUntilStartOfSeason() {
	n.IsNegotiating = false
	n.IsAcceptingOffers = false
}

func (np *NBAPlayer) ToggleTradeBlock() {
	np.IsOnTradeBlock = !np.IsOnTradeBlock
}

func (np *NBAPlayer) ToggleGLeague() {
	np.IsGLeague = !np.IsGLeague
}

func (np *NBAPlayer) ToggleTwoWay() {
	np.IsTwoWay = !np.IsTwoWay
}

func (np *NBAPlayer) RemoveFromTradeBlock() {
	np.IsOnTradeBlock = false
}

func (np *NBAPlayer) WaivePlayer() {
	np.PreviousTeamID = np.TeamID
	np.PreviousTeam = np.Team
	np.TeamID = 0
	np.Team = ""
	np.IsWaived = true
	np.IsOnTradeBlock = false
	np.IsGLeague = false
	np.IsTwoWay = false
	np.IsAcceptingOffers = true
}

func (np *NBAPlayer) ConvertWaivedPlayerToFA() {
	np.IsWaived = false
	np.IsFreeAgent = true
	np.IsAcceptingOffers = true
}

func (np *NBAPlayer) AssignFAPreferences(negotiation uint, signing uint) {
	np.NegotiationRound = negotiation
	np.SigningRound = signing
}

func (np *NBAPlayer) TradePlayer(id uint, team string) {
	np.PreviousTeam = np.Team
	np.PreviousTeamID = uint(np.TeamID)
	np.TeamID = id
	np.Team = team
	np.IsOnTradeBlock = false
}

func (np *NBAPlayer) AssignMinimumContractValue(val float64) {
	if val > 0 {
		np.MinimumValue = val
	} else {
		if np.Overall >= 100 {
			np.MaxRequested = true
		} else {
			np.MaxRequested = false
		}
		if np.Overall >= 90 && np.Overall <= 99 {
			np.MinimumValue = 20
		} else if np.Overall >= 80 && np.Overall <= 89 {
			np.MinimumValue = 10
		} else {
			if np.Year > 9 {
				np.MinimumValue = 2.5
			} else if np.Year > 8 {
				np.MinimumValue = 2.4
			} else if np.Year > 7 {
				np.MinimumValue = 2.3
			} else if np.Year > 6 {
				np.MinimumValue = 2.2
			} else if np.Year > 5 {
				np.MinimumValue = 2.0
			} else if np.Year > 4 {
				np.MinimumValue = 1.9
			} else if np.Year > 3 {
				np.MinimumValue = 1.8
			} else if np.Year > 2 {
				np.MinimumValue = 1.7
			} else if np.Year > 1 {
				np.MinimumValue = 1.6
			} else if np.Year > 0 {
				np.MinimumValue = 1.5
			} else {
				np.MinimumValue = 0.7
			}
		}
	}
}

func (np *NBAPlayer) ToggleSuperMax() {
	np.MaxRequested = true
	np.IsSuperMaxQualified = !np.IsSuperMaxQualified
}

func (np *NBAPlayer) ToggleMaxRequested() {
	np.MaxRequested = !np.MaxRequested
}

func (f *NBAPlayer) DeclineOffer(week int) {
	f.Rejections += 1
	if week >= 30 {
		f.Rejections += 2
	}
}

func (cp *NBAPlayer) AddSeasonStats(seasonStats NBAPlayerSeasonStats) {
	cp.SeasonStats = seasonStats
}

func (np *NBAPlayer) ApplyCountry(country string) {
	np.Country = country
}
