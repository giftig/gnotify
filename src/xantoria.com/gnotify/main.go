package main

import (
	"os"
	"path/filepath"
	"time"

	"xantoria.com/gnotify/config"
	"xantoria.com/gnotify/log"
	"xantoria.com/gnotify/notifier"
	"xantoria.com/gnotify/sources/calendar"
)

func main() {
	settingsFile, _ := filepath.Abs(filepath.Dir(os.Args[0]) + "/../etc/gnotify.conf")
	if len(os.Args) > 1 {
		settingsFile = os.Args[1]
	}

	// Load config and initialise log
	config.LoadConfig(settingsFile)
	log.Init()
	log.Notice("Service starting...")

	syncTicker := time.NewTicker(config.Polling.Sync)
	notificationC := make(chan *notifier.Notification)

	// Bring up a goroutine to set up notifications as it receives them
	go notifier.InitNotifications(notificationC)

	// Load events from the calendar whenever syncTicker ticks (configurable)
	// FIXME: This should be a goroutine too, but we have to do something in this thread or the
	//			  program will simply exit. This should be handled better.
	calendar.LoadEvents(syncTicker.C, notificationC)
}
