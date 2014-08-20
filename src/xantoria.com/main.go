package main

import (
  "log"
  "os"
  "time"

  "xantoria.com/config"
  "xantoria.com/gcalendar"
)

func main() {
  settingsFile := "settings.yaml"
  if len(os.Args) > 1 {
    settingsFile = os.Args[1]
  }

  if err := config.LoadConfig(settingsFile); err != nil {
    log.Fatal(err)
  }

  config.ConfigureLogger()

  syncTicker := time.NewTicker(config.Config.Polling.Sync)
  notificationChannel := make(chan *gcalendar.Notification)

  go initNotifications(notificationChannel)
  loadNotifications(syncTicker.C, notificationChannel)
}

func initNotifications(notifications <-chan *gcalendar.Notification) {
  for {
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
        log.Printf("NOTIFY: %s (%s)", notification.Id, notification.Title)
      }()
    } else {
      notification.Complete = true
    }
  }
}

func loadNotifications(
  ticks <-chan time.Time,
  notificationChannel chan *gcalendar.Notification,
) {
  for {
    log.Print("LOAD: Google calendar")
    gcalendar.GetCalendar(notificationChannel)
    _ = <-ticks
  }
}
