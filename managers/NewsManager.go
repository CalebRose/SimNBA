package managers

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/repository"
	"github.com/CalebRose/SimNBA/structs"
)

func GetAllCBBNewsLogs() []structs.NewsLog {
	db := dbprovider.GetInstance().GetDB()

	var logs []structs.NewsLog

	err := db.Where("league = ?", "CBB").Find(&logs).Error
	if err != nil {
		fmt.Println(err)
	}

	return logs
}

func GetAllNBANewsLogs() []structs.NewsLog {
	db := dbprovider.GetInstance().GetDB()

	var logs []structs.NewsLog

	err := db.Where("league = ?", "NBA").Find(&logs).Error
	if err != nil {
		fmt.Println(err)
	}

	return logs
}

func CreateNewsLog(league, message, messageType string, teamID int, ts structs.Timestamp) {
	db := dbprovider.GetInstance().GetDB()

	seasonID := 0
	weekID := 0
	week := 0
	if league == "CBB" {
		seasonID = int(ts.SeasonID)
		weekID = int(ts.CollegeWeekID)
		week = ts.CollegeWeek
	} else {
		seasonID = int(ts.SeasonID)
		weekID = int(ts.NBAWeekID)
		week = ts.NBAWeek
	}

	news := structs.NewsLog{
		League:      league,
		Message:     message,
		MessageType: messageType,
		SeasonID:    uint(seasonID),
		WeekID:      uint(weekID),
		Week:        uint(week),
		TeamID:      uint(teamID),
	}

	db.Create(&news)
}

func GetNBARelatedNews(TeamID string) []structs.NewsLog {
	ts := GetTimestamp()

	newsLogs := GetAllNBANewsLogs()

	sort.Slice(newsLogs, func(i, j int) bool {
		return newsLogs[i].CreatedAt.Unix() > newsLogs[j].CreatedAt.Unix()
	})

	newsFeed := []structs.NewsLog{}

	recentEventsCount := 0
	personalizedNewsCount := 0
	for _, news := range newsLogs {
		if recentEventsCount == 5 && personalizedNewsCount == 5 {
			break
		}
		if news.SeasonID != ts.SeasonID && news.League != "NBA" {
			continue
		}
		if recentEventsCount < 5 {
			newsFeed = append(newsFeed, news)
			recentEventsCount += 1
		} else if news.TeamID > 0 && strconv.Itoa(int(news.TeamID)) == TeamID && personalizedNewsCount < 5 {
			newsFeed = append(newsFeed, news)
			personalizedNewsCount += 1
		}
	}

	return newsFeed
}

func GetCBBRelatedNews(TeamID string) []structs.NewsLog {
	ts := GetTimestamp()

	newsLogs := GetAllCBBNewsLogs()

	sort.Slice(newsLogs, func(i, j int) bool {
		return newsLogs[i].CreatedAt.Unix() > newsLogs[j].CreatedAt.Unix()
	})

	newsFeed := []structs.NewsLog{}

	recentEventsCount := 0
	personalizedNewsCount := 0
	for _, news := range newsLogs {
		if recentEventsCount == 5 && personalizedNewsCount == 5 {
			break
		}
		if news.SeasonID != ts.SeasonID && news.League != "CBB" {
			continue
		}
		if news.TeamID == 0 && recentEventsCount < 5 {
			newsFeed = append(newsFeed, news)
			recentEventsCount += 1
		} else if news.TeamID > 0 && strconv.Itoa(int(news.TeamID)) == TeamID && personalizedNewsCount < 5 {
			newsFeed = append(newsFeed, news)
			personalizedNewsCount += 1
		}
	}

	return newsFeed
}

func CreateNotification(league, message, messageType string, teamID uint) {
	db := dbprovider.GetInstance().GetDB()

	notification := structs.Notification{
		League:           league,
		Message:          message,
		NotificationType: messageType,
		TeamID:           teamID,
	}

	repository.CreateNotification(notification, db)
}

func GetBBAInbox(cbbID, nbaID string) structs.InboxResponse {
	cbbNoti := []structs.Notification{}
	nbaNoti := []structs.Notification{}

	if cbbID != "0" {
		cbbNoti = GetNotificationByTeamIDAndLeague("CBB", cbbID)
	}
	if nbaID != "0" {
		nbaNoti = GetNotificationByTeamIDAndLeague("NBA", nbaID)
	}

	return structs.InboxResponse{
		CBBNotifications: cbbNoti,
		NBANotifications: nbaNoti,
	}
}

func GetNotificationByTeamIDAndLeague(league, teamID string) []structs.Notification {
	db := dbprovider.GetInstance().GetDB()

	noti := []structs.Notification{}

	db.Where("league = ? and team_id = ?", league, teamID).Find(&noti)

	return noti
}

func GetNotificationById(id string) structs.Notification {
	db := dbprovider.GetInstance().GetDB()

	noti := structs.Notification{}
	db.Where("id = ?", id).Find(&noti)

	return noti
}

func ToggleNotification(id string) {
	db := dbprovider.GetInstance().GetDB()

	noti := GetNotificationById(id)

	if noti.ID == 0 {
		return
	}

	noti.ToggleIsRead()
	repository.SaveNotification(noti, db)
}

func DeleteNotification(id string) {
	db := dbprovider.GetInstance().GetDB()

	noti := GetNotificationById(id)

	if noti.ID == 0 {
		return
	}

	repository.DeleteNotificationRecord(noti, db)
}
