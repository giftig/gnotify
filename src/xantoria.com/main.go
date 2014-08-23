package main

import (
  "log"
  "os"
  "path/filepath"
  "time"

  "xantoria.com/config"
  "xantoria.com/gnotify"
)

func main() {
  settingsFile, _ := filepath.Abs(filepath.Dir(os.Args[0]) + "/../etc/gnotify.conf")
  if len(os.Args) > 1 {
    settingsFile = os.Args[1]
  }

  if err := config.LoadConfig(settingsFile); err != nil {
    log.Fatal(err)
  }

  config.ConfigureLogger()

  syncTicker := time.NewTicker(config.Config.Polling.Sync)
  notificationChannel := make(chan *gnotify.Notification)

  go initNotifications(notificationChannel)
  loadNotifications(syncTicker.C, notificationChannel)
}

/**
 * Wait for notifications on the given channel and initialise them, adding
 * them to the stored notifications and starting a timer to trigger their
 * display.
 */
func initNotifications(notifications <-chan *gnotify.Notification) {
  for {
    notification := <-notifications

    // Stick it in the registered notifications store, excluding duplicates
    inserted := gnotify.AddNotification(*notification)
    if !inserted {
      log.Printf("DUP: %s (%s)", notification.Id, notification.Title)
      continue
    }

    log.Printf("INIT: %s (%s)", notification.Id, notification.Title)

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
      log.Printf("OLD: %s (%s)", notification.Id, notification.Title)
    }
  }
}

func loadNotifications(
  ticks <-chan time.Time,
  notificationChannel chan *gnotify.Notification,
) {
  for {
    log.Print("LOAD: Google calendar")
    gnotify.GetCalendar(notificationChannel)
    _ = <-ticks
  }
}
