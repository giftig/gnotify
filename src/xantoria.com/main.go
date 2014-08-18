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
  syncOccasionally(syncTicker.C, notificationChannel)
}

func initNotifications(notifications <-chan *gcalendar.Notification) {
  for {
    notification := <-notifications
    log.Printf("Initialising notification %s", notification.Title)

    diff := notification.Time.Sub(time.Now())
    if diff > 0 {
      timer := time.NewTimer(diff)
      go func() {
        _ = <-timer.C
        log.Printf("Displaying notification %s", notification.Title)
      }()
    } else {
      notification.Complete = true
    }
  }
}

func syncOccasionally(
  ticks <-chan time.Time,
  notificationChannel chan *gcalendar.Notification,
) {
  for {
    _ = <-ticks
    log.Printf("Syncing with google calendar")

    // FIXME: For now I'm just going to create some notifications and send
    // them into the notifications channel
    notification := gcalendar.Notification{
      "Test notification",
      "This is my lovely test notification",
      "/home/giftiger_wunsch/Downloads/fire.png",
      time.Date(2014, 8, 18, 20, 48, 0, 0, time.FixedZone("UTC", 0)),
      false,
    }
    notificationChannel <- &notification
  }
}
