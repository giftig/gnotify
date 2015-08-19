package rest

import (
	"time"

	"xantoria.com/gnotify/config"
	"xantoria.com/gnotify/log"
	"xantoria.com/gnotify/notifier"
)

func LoadEvents(notifC chan *notifier.Notification) {
	pollFetch := config.Sources.Rest.PollFetch

	log.Notice("Polling for calendar events every %s", pollFetch)
	ticker := time.NewTicker(pollFetch)
	ticks := ticker.C

	for {
		notifications := notifier.Fetch(config.Routing.RecipientID)
		for _, notif := range notifications {
			notifC <- &notif
		}

		_ = <-ticks
	}
}
