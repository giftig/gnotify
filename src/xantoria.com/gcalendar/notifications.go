package gcalendar

import (
  "fmt"
  "log"
  "os/exec"
  "time"

  "xantoria.com/config"
)

type Notification struct {
  Title, Message, Icon string
  Source, Id string
  Time time.Time
  Complete bool
}

/**
 * Display the notification
 *
 * TODO: Support multiple, optional methods. For now, notify-send only
 */
func (notification Notification) Display() (err error) {
  log.Printf("NOTIFY: %s (%s)", notification.Id, notification.Title)
  cmd := exec.Command(
    "/usr/bin/notify-send",
    "-i", notification.Icon,
    "-t", fmt.Sprintf("%d", config.Config.Notifications.NotifySend.Duration / time.Second),
    notification.Title,
    notification.Message,
  )

  err = cmd.Run()
  if err != nil {
    log.Printf("NOTIFY: ERROR: Command failed: `%s %s`", cmd.Path, cmd.Args)
  }
  return
}
