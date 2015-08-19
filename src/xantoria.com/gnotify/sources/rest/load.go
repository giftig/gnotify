package rest

import (
	"xantoria.com/gnotify/config"
	"xantoria.com/gnotify/log"
	"xantoria.com/gnotify/notifier"
)

func LoadEvents(notifC chan *notifier.Notification) {
	notifications := notifier.Fetch(config.Routing.RecipientID)

	log.Notice("Fetching stored notifications")
	for _, notif := range notifications {
		notifC <- &notif
	}
}
