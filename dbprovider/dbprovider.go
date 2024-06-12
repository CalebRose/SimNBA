package dbprovider

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/CalebRose/SimNBA/config"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/crypto/ssh"
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
	sshConfig := config.GetSSHConfig()
	localPort, localErr := setupSSHTunnel(&sshConfig)
	if localErr != nil {
		log.Fatalf("Failed to establish SSH tunnel: %v", localErr)
	}

	var err error
	c := config.Config(localPort)
	db, err = gorm.Open(mysql.Open(c["cs"]), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return false
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get underlying sql.DB:", err)
		return false
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// db.AutoMigrate(&structs.GlobalPlayer{})
	// db.AutoMigrate(&structs.CollegePlayerStats{})
	// db.AutoMigrate(&structs.CollegePlayerSeasonStats{})
	// db.AutoMigrate(&structs.CollegePlayer{})
	// db.AutoMigrate(&structs.HistoricCollegePlayer{})
	// db.AutoMigrate(&structs.UnsignedPlayer{})
	// db.AutoMigrate(&structs.TransferPortalProfile{})
	// db.AutoMigrate(&structs.CollegeCoach{})
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
	// db.AutoMigrate(&structs.CollegePollOfficial{})
	// db.AutoMigrate(&structs.CollegePollSubmission{})
	// db.AutoMigrate(&structs.CollegePromise{})
	// db.AutoMigrate(&structs.DraftPick{})
	// db.AutoMigrate(&structs.ISLScoutingDept{})
	// db.AutoMigrate(&structs.ISLScoutingReport{})
	// db.AutoMigrate(&structs.NBACapsheet{})
	// db.AutoMigrate(&structs.NBAContract{})
	// db.AutoMigrate(&structs.NBAContractOffer{})
	// db.AutoMigrate(&structs.NBAExtensionOffer{})
	// db.AutoMigrate(&structs.NBAConference{})
	// db.AutoMigrate(&structs.NBADivision{})
	// db.AutoMigrate(&structs.NBAWarRoom{})
	// db.AutoMigrate(&structs.NBADraftee{})
	// db.AutoMigrate(&structs.ScoutingProfile{})
	// db.AutoMigrate(&structs.NBAGameplan{})
	// db.AutoMigrate(&structs.NBAMatch{})
	// db.AutoMigrate(&structs.NBASeries{})
	// db.AutoMigrate(&structs.NBAPlayer{})
	// db.AutoMigrate(&structs.NBAPlayerStats{})
	// db.AutoMigrate(&structs.NBAPlayerSeasonStats{})
	// db.AutoMigrate(&structs.RetiredPlayer{})
	// db.AutoMigrate(&structs.NBARequest{})
	// db.AutoMigrate(&structs.NBAStandings{})
	// db.AutoMigrate(&structs.NBATeam{})
	// db.AutoMigrate(&structs.NBATeamStats{})
	// db.AutoMigrate(&structs.NBATeamSeasonStats{})
	// db.AutoMigrate(&structs.NBATradePreferences{})
	// db.AutoMigrate(&structs.NBATradeProposal{})
	// db.AutoMigrate(&structs.NBATradeOption{})
	// db.AutoMigrate(&structs.NBAUser{})
	// db.AutoMigrate(&structs.NBAWaiverOffer{})
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

// setupSSHTunnel establishes an SSH tunnel and forwards a local port to the remote database port.
// Returns the local port and any error encountered.
func setupSSHTunnel(config *config.SshTunnelConfig) (string, error) {
	sshConfig := &ssh.ClientConfig{
		User: config.SshUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.SshPassword),
		},
		// CAUTION: In production, you should use a more secure HostKeyCallback.
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to SSH server
	sshClient, err := ssh.Dial("tcp", net.JoinHostPort(config.SshHost, config.SshPort), sshConfig)
	if err != nil {
		return "", err
	}

	// Setup local port forwarding
	localListener, err := net.Listen("tcp", "localhost:"+config.LocalPort)
	if err != nil {
		return "", err
	}

	go func() {
		defer localListener.Close()
		for {
			localConn, err := localListener.Accept()
			if err != nil {
				log.Printf("Failed to accept local connection: %s", err)
				continue
			}

			// Handle the connection in a new goroutine
			go func() {
				defer localConn.Close()

				// Connect to the remote database server through the SSH tunnel
				remoteConn, err := sshClient.Dial("tcp", net.JoinHostPort(config.DbHost, config.DbPort))
				if err != nil {
					log.Printf("Failed to dial remote server: %s", err)
					return
				}
				defer remoteConn.Close()

				// Copy data between the local connection and the remote connection
				copyConn(localConn, remoteConn)
			}()
		}
	}()

	return localListener.Addr().String(), nil
}

// copyConn copies data between two io.ReadWriteCloser objects (e.g., network connections)
func copyConn(localConn, remoteConn io.ReadWriteCloser) {
	// Start goroutine to copy data from local to remote
	go func() {
		_, err := io.Copy(remoteConn, localConn)
		if err != nil {
			log.Printf("Error copying from local to remote: %v", err)
		}
		localConn.Close()
		remoteConn.Close()
	}()

	// Copy data from remote to local in the main goroutine (or vice versa)
	_, err := io.Copy(localConn, remoteConn)
	if err != nil {
		log.Printf("Error copying from remote to local: %v", err)
	}
	// Ensure connections are closed when copying is done or an error occurs
	localConn.Close()
	remoteConn.Close()
}
