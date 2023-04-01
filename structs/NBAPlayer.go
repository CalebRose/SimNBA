package structs

import "github.com/jinzhu/gorm"

type NBAPlayer struct {
	gorm.Model
	BasePlayer
	PlayerID             uint
	TeamID               uint
	TeamAbbr             string
	CollegeID            uint
	College              string
	DraftPickID          uint
	DraftedRound         uint
	DraftPick            uint
	DraftedTeamID        uint
	DraftedTeamAbbr      string
	PreviousTeamID       uint
	PreviousTeam         string
	PrimeAge             uint
	IsNBA                bool
	MaxRequested         bool
	IsSuperMaxQualified  bool
	IsFreeAgent          bool
	IsGLeague            bool
	IsTwoWay             bool
	IsWaived             bool
	IsOnTradeBlock       bool
	IsFirstTeamANBA      bool
	IsDPOY               bool
	IsMVP                bool
	IsInternational      bool
	IsRetiring           bool
	IsAcceptingOffers    bool
	IsNegotiating        bool
	MinimumValue         float64
	SigningRound         uint
	NegotiationRound     uint
	PositionOne          string
	PositionTwo          string
	PositionThree        string
	Position1Minutes     uint
	Position2Minutes     uint
	Position3Minutes     uint
	InsidePercentage     uint
	MidPercentage        uint
	ThreePointPercentage uint
	Offers               []NBAContractOffer   `gorm:"foreignKey:PlayerID"`
	Contract             NBAContract          `gorm:"foreignKey:PlayerID"`
	Stats                []NBAPlayerStats     `gorm:"foreignKey:NBAPlayerID"`
	SeasonStats          NBAPlayerSeasonStats `gorm:"foreignKey:NBAPlayerID"`
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

func (n *NBAPlayer) BecomeFreeAgent() {
	n.TeamAbbr = "FA"
	n.TeamID = 0
	n.IsFreeAgent = true
	n.IsOnTradeBlock = false
}

func (n *NBAPlayer) SignWithTeam(teamID uint, team string) {
	n.TeamAbbr = team
	n.TeamID = teamID
	n.IsFreeAgent = false
}

func (n *NBAPlayer) Progress(p NBAPlayerProgressions) {
	n.Shooting2 = p.Shooting2
	n.Shooting3 = p.Shooting3
	n.FreeThrow = p.FreeThrow
	n.Ballwork = p.Ballwork
	n.Finishing = p.Finishing
	n.Rebounding = p.Rebounding
	n.InteriorDefense = p.InteriorDefense
	n.PerimeterDefense = p.PerimeterDefense
	n.Overall = p.Overall
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

func (n *NBAPlayer) ToggleIsNegotiating() {
	n.IsNegotiating = true
}

func (np *NBAPlayer) ToggleTradeBlock() {
	np.IsOnTradeBlock = !np.IsOnTradeBlock
}

func (np *NBAPlayer) RemoveFromTradeBlock() {
	np.IsOnTradeBlock = false
}

func (np *NBAPlayer) WaivePlayer() {
	np.TeamID = 0
	np.TeamAbbr = ""
	np.IsWaived = true
	np.IsOnTradeBlock = false
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
	np.PreviousTeam = np.TeamAbbr
	np.PreviousTeamID = uint(np.TeamID)
	np.TeamID = id
	np.TeamAbbr = team
	np.IsOnTradeBlock = false
}

func (np *NBAPlayer) AssignMinimumContractValue(val float64) {
	np.MinimumValue = val
}

func (np *NBAPlayer) ToggleSuperMax() {
	np.MaxRequested = true
	np.IsSuperMaxQualified = !np.IsSuperMaxQualified
}

func (np *NBAPlayer) ToggleMaxRequested() {
	np.MaxRequested = !np.MaxRequested
}
