package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"xantoria.com/gnotify/config"
	"xantoria.com/gnotify/log"
)

const (
	NotifySend = iota
	AudioAlert = iota
)

// Deliver takes a notification and sets it up to be displayed at the right time if we're the
// correct recipient, and routes it to the right recipient if we are aware of one
func (notif *Notification) Deliver() {
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
	if !shouldDisplay {
		notif.reroute()
		return
	}

	if notif.Complete {
		log.Info("Ignoring expired notification \"%s\" (%s)", notif.Id, notif.Title)
		return
	}

	// Create a timer which displays the notification at the correct time if not expired
	diff := notif.Time.Sub(time.Now())
	log.Debug("DIFF: %v", diff)
	if diff <= 0 {
		go func() {
			notif.Display()
		}()
	} else {
		timer := time.NewTimer(diff)
		go func() {
			_ = <-timer.C
			notif.Display()
		}()
	}
}

// reroute checks if we know how to contact the recipient of the given notification and passes it
// on to them if we do. It respects both recipient ID and groups
func (notif *Notification) reroute() {
	// Marshal the notification ready to be sent to recipients
	data, err := json.Marshal(notif)
	if err != nil {
		log.Error("Failed to marshal notification for notification %s", notif.Id)
		return
	}
	rawData := bytes.NewBuffer(data)

	// TODO: Probably better to make this all asynchronous, spinning up goroutines for each
	for _, recipient := range config.Routing.KnownRecipients {
		validRecipient := false
		if recipient.ID == notif.Recipient {
			validRecipient = true
		}

		// If recipient ID didn't match, check groups
		if !validRecipient {
			for _, group := range recipient.Groups {
				if group == notif.Recipient {
					validRecipient = true
					break
				}
			}
		}

		// This obviously isn't who the message is for, so try the next one
		if !validRecipient {
			continue
		}

		// This is the correct recipient, so let's pass the message on
		url := fmt.Sprintf("http://%s:%d/notify/route/", recipient.Host, recipient.Port)
		resp, err := http.Post(url, "application/json", rawData)
		log.Info("Forwarding notification %s to %s (%s)...", notif.Id, recipient.ID, recipient.Host)

		// The client may well not be online, in which case they'll ask for it when they come on
		if err != nil {
			log.Warning(
				"Failed to deliver notification %s to %s (%s): %v",
				notif.Id, recipient.ID, recipient.Host, err,
			)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			log.Error(
				"Unexpected HTTP %d from client %s when delivering notification %s",
				resp.StatusCode,
				recipient.ID,
				notif.Id,
			)
			continue
		}
	}
}

// Display displays the notification to the user
func (notif *Notification) Display() {
	cfg := config.Notifications
	displayed := false

	if cfg.NotifySend.Enabled {
		go func() {
			notifySend(notif)
		}()
		displayed = true
	}
	if cfg.AudioAlert.Enabled {
		go func() {
			audioAlert(notif)
		}()
		displayed = true
	}
	// TODO: Add more display methods here

	if displayed {
		log.Debug("Displayed notification %s (%s)", notif.Id, notif.Title)
	}
}

// notifySend uses the notify-send application to display a message to the user
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

	// Patch for a bug in notify-send which causes it to not show messages
	// See www.archivum.info/ubuntu-bugs: Bug 1424243
	msg := strings.Replace(notif.Message, "&", "and", -1)

	cmd := exec.Command(
		"/usr/bin/env",
		"notify-send",
		"-i", filepath.Join(config.Static.IconPath, icon),
		"-t", fmt.Sprintf("%d", config.Notifications.NotifySend.Duration/time.Millisecond),
		"-u", urgency,
		notif.Title,
		msg,
	)

	log.Debug("Command: %s %s", cmd.Path, cmd.Args)
	if err := cmd.Run(); err != nil {
		log.Critical("notify-send failed: `%s %s` (%v)", cmd.Path, cmd.Args, err)
		return
	}
	notif.MarkComplete()
}

func audioAlert(notif *Notification) {
	cfg := config.Notifications.AudioAlert

	sound := cfg.DefaultSound
	maxThresh := 0

	for thresh, snd := range cfg.Sounds {
		if notif.Priority >= thresh && thresh > maxThresh {
			maxThresh = thresh
			sound = snd
		}
	}

	// Drop out early if we don't have a sound
	if sound == "" {
		log.Warning("[%s, priority %d] No sound configured; aborting", notif.Id, notif.Priority)
		return
	}

	// TODO: Implement a couple more drivers for sounds
	switch cfg.Driver {
	case "mplayer":
		cmd := exec.Command(
			"/usr/bin/env",
			"mplayer",
			"-really-quiet", // I shit you not, that's the actual flag
			"-msglevel", "all=0",
			"-endpos", fmt.Sprintf("%d", cfg.CutOffLength),
			"-loop", fmt.Sprintf("%d", cfg.Repeats),
			sound,
		)

		log.Debug("Command: %s %s", cmd.Path, cmd.Args)
		if err := cmd.Run(); err != nil {
			log.Critical("mplayer failed: `%s %s` (%v)", cmd.Path, cmd.Args, err)
			return
		}

	default:
		log.Error(
			"[%s] Driver %s not supported; currently only mplayer is supported! Aborting.",
			notif.Id, cfg.Driver,
		)
	}
}
