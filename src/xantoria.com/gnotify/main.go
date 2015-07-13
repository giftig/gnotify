package main

import (
	"os"
	"os/user"
	"path/filepath"

	"xantoria.com/gnotify/config"
	"xantoria.com/gnotify/log"
	"xantoria.com/gnotify/notifier"
	"xantoria.com/gnotify/sources/calendar"
	"xantoria.com/gnotify/sources/rest"
	"xantoria.com/gnotify/sources/todo"
)

func pickSettingsFile() string {
	// Did we provide an explicit config file? That takes precedence
	if len(os.Args) > 1 {
		f := os.Args[1]

		if _, err := os.Stat(f); err != nil {
			// If the specific config file provided is no good, fail immediately
			log.Fatalf("Error in provided config file %s! %v", f, err)
		}

		return f
	}

	// Let's see if we the user has their own config file
	usr, err := user.Current()
	if err == nil {
		f := usr.HomeDir + "/.gnotify/gnotify.conf"
		if _, err := os.Stat(f); err == nil {
			return f
		}
	}

	// No user-specific config; let's try ../etc (eg. /opt/etc if we're installed to /opt/bin)
	if f, err := filepath.Abs(filepath.Dir(os.Args[0]) + "/../etc/gnotify.conf"); err == nil {
		if _, err := os.Stat(f); err == nil {
			return f
		}
	}

	log.Fatal("Couldn't find a usable config file!")
	return ""
}

func main() {
	settingsFile := pickSettingsFile()

	// Load config and initialise log
	config.LoadConfig(settingsFile)
	log.Init()
	log.Notice("Service starting")
	log.Info("Using settings %s", settingsFile)

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
