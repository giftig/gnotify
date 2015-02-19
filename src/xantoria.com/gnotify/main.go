package main

import (
	"os"
	"path/filepath"
	"time"

	"xantoria.com/gnotify/config"
	"xantoria.com/gnotify/log"
)

func main() {
	settingsFile, _ := filepath.Abs(filepath.Dir(os.Args[0]) + "/../etc/gnotify.conf")
	if len(os.Args) > 1 {
		settingsFile = os.Args[1]
	}

	// Load config and initialise log
	config.LoadConfig(settingsFile)
	log.Init()
	log.Info("Service starting...")

	syncTicker := time.NewTicker(config.Polling.Sync)
	notificationChannel := make(chan *Notification)

	go initNotifications(notificationChannel)
	loadNotifications(syncTicker.C, notificationChannel)
}

/**
 * Wait for notifications on the given channel and initialise them, adding
 * them to the stored notifications and starting a timer to trigger their
 * display.
 */
func initNotifications(notifications <-chan *Notification) {
	for {
		notification := <-notifications

		// Stick it in the registered notifications store, excluding duplicates
		inserted := AddNotification(*notification)
		if !inserted {
			log.Debug("DUP: %s (%s)", notification.Id, notification.Title)
			continue
		}

		log.Info("INIT: %s (%s)", notification.Id, notification.Title)

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
			log.Warning("OLD: %s (%s)", notification.Id, notification.Title)
		}
	}
}

func loadNotifications(
	ticks <-chan time.Time,
	notificationChannel chan *Notification,
) {
	for {
		log.Notice("LOAD: Google calendar")
		GetCalendar(notificationChannel)
		_ = <-ticks
	}
}
