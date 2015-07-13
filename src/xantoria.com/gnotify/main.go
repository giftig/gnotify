package main

import (
	"os"
	"path/filepath"

	"xantoria.com/gnotify/config"
	"xantoria.com/gnotify/log"
	"xantoria.com/gnotify/notifier"
	"xantoria.com/gnotify/sources/calendar"
	"xantoria.com/gnotify/sources/rest"
	"xantoria.com/gnotify/sources/todo"
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

	notificationC := make(chan *notifier.Notification)

	// Bring up a goroutine to set up notifications as it receives them
	go notifier.InitNotifications(notificationC)

	// Load events from the calendar whenever syncTicker ticks (configurable)
	go calendar.LoadEvents(notificationC)

	// Load events from the configured to-do list
	go todo.LoadEvents(notificationC)

	// Listen for routed or freshly-triggered events over REST
	rest.Listen(notificationC)
}
