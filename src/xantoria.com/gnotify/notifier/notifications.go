package notifier

import (
	"fmt"
	"os/exec"
	"time"

	"xantoria.com/gnotify/config"
	"xantoria.com/gnotify/log"
)

type Notification struct {
	Title, Message, Icon string
	Source, Id           string
	Time                 time.Time
	Complete             bool
}

var notifications []Notification

// Display displays the notification to the user
// TODO: Support multiple, optional methods. For now, notify-send only
func (notification *Notification) Display() (err error) {
	log.Debug("Displaying notification %s (%s)", notification.Id, notification.Title)
	cmd := exec.Command(
		"/usr/bin/notify-send",
		"-i", notification.Icon,
		"-t", fmt.Sprintf("%d", config.Notifications.NotifySend.Duration/time.Millisecond),
		notification.Title,
		notification.Message,
	)

	err = cmd.Run()
	if err != nil {
		log.Critical("notify-send failed: `%s %s` (%v)", cmd.Path, cmd.Args, err)
	}
	return
}

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
		notification := <-notifications

		// Stick it in the registered notifications store, excluding duplicates
		inserted := AddNotification(*notification)
		if !inserted {
			log.Debug("Ignored duplicate notification %s (%s)", notification.Id, notification.Title)
			continue
		}

		log.Info("Storing notification %s (%s)", notification.Id, notification.Title)

		diff := notification.Time.Sub(time.Now())
		if diff > 0 {
			timer := time.NewTimer(diff)
			go func() {
				_ = <-timer.C
				err := notification.Display()
				if err == nil {
					// TODO: This doesn't actually do anything until they're actually
					// TODO: organised and properly synced with a local datastore.
					notification.Complete = true
				}
			}()
		} else {
			// This shouldn't really happen as we only ask google for future events
			// That's specific to google calendar notifications, of course.
			log.Warning("Notification %s (%s) has expired already", notification.Id, notification.Title)
		}
	}
}
