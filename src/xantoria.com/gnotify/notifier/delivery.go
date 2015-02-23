package notifier

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	"xantoria.com/gnotify/config"
	"xantoria.com/gnotify/log"
)

const (
	NotifySend = iota
)

// Display displays the notification to the user
func (notif *Notification) Display(route int) {
	log.Debug("Displaying notification %s (%s) via method %d", notif.Id, notif.Title, route)

	switch route {
	case NotifySend:
		notifySend(notif)

	default:
		log.Error("Unknown route ID %d", route)
		return
	}

	// TODO: Doesn't do anything until it's properly synced with a local datastore
	notif.Complete = true
}

// notifySend users the notify-send application to display a message to the user
func notifySend(notif *Notification) {
	// Decide on the correct urgency string for notify-send
	urgency := "critical"
	if notif.Priority < NORMAL {
		urgency = "low"
	} else if notif.Priority < IMPORTANT {
		urgency = "normal"
	}

	// Override the notification icon with an urgency-based one if it's not set
	icon := notif.Icon
	if icon == "" {
		icon = fmt.Sprintf("%s.png", urgency)
	}

	cmd := exec.Command(
		"/usr/bin/notify-send",
		"-i", filepath.Join(config.Static.IconPath, icon),
		"-t", fmt.Sprintf("%d", config.Notifications.NotifySend.Duration/time.Millisecond),
		"-u", urgency,
		notif.Title,
		notif.Message,
	)

	if err := cmd.Run(); err != nil {
		log.Critical("notify-send failed: `%s %s` (%v)", cmd.Path, cmd.Args, err)
	}
}
