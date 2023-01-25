package structs

import "github.com/jinzhu/gorm"

type NBACapSheet struct {
	gorm.Model
	TeamID               uint
	CurrentSeason        uint
	Year1Total           float64
	Year1Cap             float64
	Year1CashTransferred float64
	Year1CashReceived    float64
	Year2Total           float64
	Year2Cap             float64
	Year2CashTransferred float64
	Year2CashReceived    float64
	Year3Total           float64
	Year3Cap             float64
	Year3CashTransferred float64
	Year3CashReceived    float64
	Year4Total           float64
	Year4Cap             float64
	Year4CashTransferred float64
	Year4CashReceived    float64
	Year5Total           float64
	Year5Cap             float64
	Year5CashTransferred float64
	Year5CashReceived    float64
	IsOverMax            bool
	// Contracts            []NBAContract
}

func (cs *NBACapSheet) SyncTotals(year1 float64, year2 float64, year3 float64, year4 float64, year5 float64) {
	cs.Year1Total = year1 + cs.Year1CashTransferred - cs.Year1CashReceived
	cs.Year2Total = year2 + cs.Year2CashTransferred - cs.Year2CashReceived
	cs.Year3Total = year3 + cs.Year3CashTransferred - cs.Year3CashReceived
	cs.Year4Total = year4 + cs.Year4CashTransferred - cs.Year4CashReceived
	cs.Year5Total = year5 + cs.Year5CashTransferred - cs.Year5CashReceived
}

func (cs *NBACapSheet) SyncByYear() {
	cs.Year1CashReceived = cs.Year2CashReceived
	cs.Year1CashTransferred = cs.Year2CashTransferred
	cs.Year2CashReceived = cs.Year3CashReceived
	cs.Year2CashTransferred = cs.Year3CashTransferred
	cs.Year3CashReceived = cs.Year4CashReceived
	cs.Year3CashTransferred = cs.Year4CashTransferred
	cs.Year4CashReceived = cs.Year5CashReceived
	cs.Year4CashTransferred = cs.Year5CashTransferred
	cs.Year5CashReceived = 0
	cs.Year5CashTransferred = 0
	cs.Year1Cap = cs.Year2Cap
	cs.Year2Cap = cs.Year3Cap
	cs.Year3Cap = cs.Year4Cap
	cs.Year4Cap = cs.Year5Cap
	cs.Year5Cap = cs.Year5Cap * 1.12
}
