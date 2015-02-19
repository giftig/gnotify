package main

import (
	"fmt"
	"os/exec"
	"time"

	"xantoria.com/gnotify/config"
	"xantoria.com/gnotify/log"
)

type Notification struct {
	Title, Message, Icon string
	Source, Id           string
	Time                 time.Time
	Complete             bool
}

var notifications []Notification

/**
 * Display the notification
 *
 * TODO: Support multiple, optional methods. For now, notify-send only
 */
func (notification *Notification) Display() (err error) {
	log.Debug("NOTIFY: %s (%s)", notification.Id, notification.Title)
	cmd := exec.Command(
		"/usr/bin/notify-send",
		"-i", notification.Icon,
		"-t", fmt.Sprintf("%d", config.Notifications.NotifySend.Duration/time.Millisecond),
		notification.Title,
		notification.Message,
	)

	err = cmd.Run()
	if err != nil {
		log.Critical("NOTIFY: Command failed: `%s %s`", cmd.Path, cmd.Args)
	}
	return
}

/**
 * Try to add a notification to the slice of stored notifications; return
 * true if it was added (source and id must be unique together).
 */
func AddNotification(notification Notification) bool {
	for _, n := range notifications {
		if n.Id == notification.Id && n.Source == notification.Source {
			return false
		}
	}
	notifications = append(notifications, notification)
	return true
}