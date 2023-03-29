package dbprovider

import (
	"fmt"
	"log"
	"sync"

	"github.com/CalebRose/SimNBA/config"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Provider struct {
}

var db *gorm.DB
var once sync.Once
var instance *Provider

func GetInstance() *Provider {
	once.Do(func() {
		instance = &Provider{}
	})
	return instance
}

func (p *Provider) InitDatabase() bool {
	fmt.Println("Database initializing...")
	var err error
	c := config.Config()
	db, err = gorm.Open(mysql.Open(c["cs"]), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return false
	}
	// db.AutoMigrate(&structs.GlobalPlayer{})
	// db.AutoMigrate(&structs.CollegePlayerStats{})
	// db.AutoMigrate(&structs.CollegePlayerSeasonStats{})
	// db.AutoMigrate(&structs.CollegePlayer{})
	// db.AutoMigrate(&structs.HistoricCollegePlayer{})
	// db.AutoMigrate(&structs.RecruitPointAllocation{})
	// db.AutoMigrate(&structs.Recruit{})
	// db.AutoMigrate(&structs.PlayerRecruitProfile{})
	// db.AutoMigrate(&structs.TeamRecruitingProfile{})
	// db.AutoMigrate(&structs.CollegeWeek{})
	// db.AutoMigrate(&structs.CollegeConference{})
	// db.AutoMigrate(&structs.CollegeStandings{})
	// db.AutoMigrate(&structs.TeamStats{})
	// db.AutoMigrate(&structs.TeamSeasonStats{})
	// db.AutoMigrate(&structs.Team{})

	// db.AutoMigrate(&structs.DraftPick{})
	// db.AutoMigrate(&structs.NBACapsheet{})
	// db.AutoMigrate(&structs.NBAContract{})
	// db.AutoMigrate(&structs.NBAContractOffer{})
	// db.AutoMigrate(&structs.NBAConference{})
	// db.AutoMigrate(&structs.NBADivision{})
	// db.AutoMigrate(&structs.NBADraftee{})
	// db.AutoMigrate(&structs.NBAGameplan{})
	// db.AutoMigrate(&structs.NBAMatch{})
	// db.AutoMigrate(&structs.NBAPlayer{})
	// db.AutoMigrate(&structs.NBAPlayerStats{})
	// db.AutoMigrate(&structs.NBAPlayerSeasonStats{})
	// db.AutoMigrate(&structs.RetiredPlayer{})
	// db.AutoMigrate(&structs.NBARequest{})
	// db.AutoMigrate(&structs.NBATeam{})
	// db.AutoMigrate(&structs.NBATeamStats{})
	// db.AutoMigrate(&structs.NBATeamSeasonStats{})
	// db.AutoMigrate(&structs.NBATradeProposal{})
	// db.AutoMigrate(&structs.NBATradeOption{})
	// db.AutoMigrate(&structs.NBAUser{})
	// db.AutoMigrate(&structs.Arena{})

	// db.AutoMigrate(&structs.Gameplan{})
	// db.AutoMigrate(&structs.Match{})
	// db.AutoMigrate(&structs.NBAWeek{})
	// db.AutoMigrate(&structs.Player{})
	// db.AutoMigrate(&structs.PlayerStats{})
	// db.AutoMigrate(&structs.NewsLog{})
	// db.AutoMigrate(&structs.Request{})
	// db.AutoMigrate(&structs.Season{})
	// db.AutoMigrate(&structs.Timestamp{})
	return true
}

func (p *Provider) GetDB() *gorm.DB {
	return db
}
