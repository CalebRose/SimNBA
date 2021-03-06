package dbprovider

import (
	"fmt"
	"log"
	"sync"

	"github.com/CalebRose/SimNBA/config"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
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
	db, err = gorm.Open(c["db"], c["cs"])
	if err != nil {
		log.Fatal(err)
		return false
	}
	db.AutoMigrate(&structs.CollegeWeek{})
	db.AutoMigrate(&structs.Gameplan{})
	db.AutoMigrate(&structs.Match{})
	db.AutoMigrate(&structs.NBAWeek{})
	db.AutoMigrate(&structs.Player{})
	db.AutoMigrate(&structs.PlayerStats{})
	db.AutoMigrate(&structs.RecruitingPoints{})
	db.AutoMigrate(&structs.RecruitingProfile{})
	db.AutoMigrate(&structs.Request{})
	db.AutoMigrate(&structs.Season{})
	db.AutoMigrate(&structs.Team{})
	db.AutoMigrate(&structs.TeamStats{})
	db.AutoMigrate(&structs.Timestamp{})
	return true
}

func (p *Provider) GetDB() *gorm.DB {
	return db
}
