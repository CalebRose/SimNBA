package structs

import "github.com/jinzhu/gorm"

type NBAContractOffer struct {
	gorm.Model
	PlayerID     uint
	TeamID       uint
	SeasonID     uint
	Team         string
	TotalYears   uint
	ContractType string
	TotalCost    float64
	Year1Total   float64
	Year2Total   float64
	Year3Total   float64
	Year4Total   float64
	Year5Total   float64
	Year1Opt     bool
	Year2Opt     bool
	Year3Opt     bool
	Year4Opt     bool
	Year5Opt     bool
	IsAccepted   bool
	IsRejected   bool
	// Do we want to kep track of the year?
}

func (o *NBAContractOffer) AcceptOffer() {
	o.IsAccepted = true
}

func (o *NBAContractOffer) RejectOffer() {
	o.IsRejected = false
}
