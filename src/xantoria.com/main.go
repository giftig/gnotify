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

func initNotifications(notifications <-chan *gnotify.Notification) {
  for {
    // TODO: Stick notifications somewhere and make sure that
    // TODO: id and source are unique together
    notification := <-notifications
    log.Printf(
      "INIT: %s (%s)",
      notification.Id,
      notification.Title,
    )

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
      notification.Complete = true
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
