package notifier

import (
	"time"

	"xantoria.com/gnotify/log"
)

type Notification struct {
	Title, Message, Icon string
	Source, Id           string
	Priority             int
	Time                 time.Time
	Complete             bool
	Recipient            string
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
			log.Debug("Ignoring duplicate notification %s (%s)", notif.Id, notif.Title)
			continue
		}

		// Check if it's expired
		diff := notif.Time.Sub(time.Now())
		if diff <= 0 {
			log.Debug("Ignoring expired notification %s (%s)", notif.Id, notif.Title)
			continue
		}
		log.Info("Storing notification %s (%s)", notif.Id, notif.Title)

		notif.Deliver()
	}
}
