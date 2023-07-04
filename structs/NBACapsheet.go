package structs

import "github.com/jinzhu/gorm"

type NBACapsheet struct {
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

func (cs *NBACapsheet) AssignID(id uint) {
	cs.ID = id
	cs.TeamID = id
}

func (cs *NBACapsheet) ResetCapsheet() {
	cs.Year1Total = 0
	cs.Year2Total = 0
	cs.Year3Total = 0
	cs.Year4Total = 0
	cs.Year5Total = 0
}

func (cs *NBACapsheet) SyncTotals(year1 float64, year2 float64, year3 float64, year4 float64, year5 float64) {
	cs.Year1Total = year1 + cs.Year1CashTransferred - cs.Year1CashReceived
	cs.Year2Total = year2 + cs.Year2CashTransferred - cs.Year2CashReceived
	cs.Year3Total = year3 + cs.Year3CashTransferred - cs.Year3CashReceived
	cs.Year4Total = year4 + cs.Year4CashTransferred - cs.Year4CashReceived
	cs.Year5Total = year5 + cs.Year5CashTransferred - cs.Year5CashReceived
}

func (cs *NBACapsheet) SyncByYear() {
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

func (nc *NBACapsheet) SubtractFromCapsheetViaTrade(contract NBAContract) {
	nc.Year1Cap += contract.Year1Total
	nc.Year1Total -= contract.Year1Total
	nc.Year2Total -= contract.Year2Total
	nc.Year3Total -= contract.Year3Total
	nc.Year4Total -= contract.Year4Total
	nc.Year5Total -= contract.Year5Total
}

func (nc *NBACapsheet) CutPlayerFromCapsheet(contract NBAContract) {
	nc.Year1Cap += contract.Year1Total + contract.Year2Total + contract.Year3Total + contract.Year4Total + contract.Year5Total
	nc.Year1Total -= contract.Year1Total
	nc.Year2Total -= contract.Year2Total
	nc.Year3Total -= contract.Year3Total
	nc.Year4Total -= contract.Year4Total
	nc.Year5Total -= contract.Year5Total
}

func (nc *NBACapsheet) NegotiateSalaryDifference(CashTransferring float64) {
	nc.Year1Total -= CashTransferring
	nc.Year1CashTransferred += CashTransferring
}

func (nc *NBACapsheet) AddContractViaTrade(contract NBAContract, differenceValue float64) {
	// nc.Y1Bonus += contract.Y1Bonus
	nc.Year1CashReceived += differenceValue
	nc.Year1Total += differenceValue
	nc.Year2Total += contract.Year2Total
	nc.Year3Total += contract.Year3Total
	nc.Year4Total += contract.Year4Total
	nc.Year5Total += contract.Year5Total
}
