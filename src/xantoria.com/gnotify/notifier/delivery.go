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

// Deliver takes a notification and sets it up to be displayed at the right time if we're the
// correct recipient, and routes it to the right recipient if we are aware of one
func (notif *Notification) Deliver() {
	// TODO: Check here if we know anything about the intended recipients and pass it on

	// Figure out if we're the intended recipient of this notification
	shouldDisplay := notif.Recipient == "" || notif.Recipient == config.Routing.RecipientID
	if !shouldDisplay {
		for _, group := range config.Routing.Groups {
			if group == notif.Recipient {
				shouldDisplay = true
				break
			}
		}
	}

	// Create a timer which displays the notification at the correct time if not expired
	diff := notif.Time.Sub(time.Now())
	if diff <= 0 {
		log.Info("Ignoring expired notification \"%s\" (%s)", notif.Id, notif.Title)
	} else if shouldDisplay {
		timer := time.NewTimer(diff)
		go func() {
			_ = <-timer.C
			notif.Display()
		}()
	}
}

// Display displays the notification to the user
func (notif *Notification) Display() {
	cfg := config.Notifications
	displayed := false

	if cfg.NotifySend.Enabled {
		notifySend(notif)
	}
	// TODO: Add more display methods here

	if displayed {
		log.Debug("Displayed notification %s (%s)", notif.Id, notif.Title)
		notif.Complete = true
		notif.Save()
	}
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
