package structs

import "gorm.io/gorm"

type InboxResponse struct {
	CBBNotifications []Notification
	NBANotifications []Notification
}

type Notification struct {
	gorm.Model
	TeamID           uint
	League           string
	NotificationType string
	Message          string
	Subject          string
	IsRead           bool
}

func (n *Notification) ToggleIsRead() {
	n.IsRead = !n.IsRead
}
