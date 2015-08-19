package notifier

import (
	"time"

	"xantoria.com/gnotify/log"
)

type Notification struct {
	Title     string `json:"title"`
	Message   string `json:"message"`
	Recipient string `json:"recipient"`

	Id     string `json:"id"`
	Source string `json:"source"`

	Priority int    `json:"priority"`
	Icon     string `json:"icon"`

	Time     time.Time `json:"time"`
	Complete bool      `json:"complete"`
}

// Priority levels
const (
	TRIVIAL   = 10
	NORMAL    = 50
	IMPORTANT = 75
	URGENT    = 100
)

var notifications []Notification

// AddNotification tries to add a notification to the slice of stored notifications; returns
// true if it was added (source and id must be unique together).
func AddNotification(notification Notification) bool {
	for _, n := range notifications {
		if n.Id == notification.Id && n.Source == notification.Source {
			return false
		}
	}
	notifications = append(notifications, notification)
	return true
}

// initNotifications waits for notifications on the given channel and initialises them, adding
// them to the stored notifications and starting a timer to trigger their display at the right time
func InitNotifications(notifications <-chan *Notification) {
	for {
		notif := <-notifications

		// Stick it in the registered notifications store, excluding duplicates
		inserted := AddNotification(*notif)
		if !inserted {
			continue
		}

		log.Info("Storing notification \"%s\" (%s)", notif.Id, notif.Title)
		notif.Save()

		notif.Deliver()
	}
}
