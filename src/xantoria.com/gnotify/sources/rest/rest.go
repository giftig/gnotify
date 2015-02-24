package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"xantoria.com/gnotify/config"
	"xantoria.com/gnotify/log"
	"xantoria.com/gnotify/notifier"
)

var notificationC chan<- *notifier.Notification

// handleNotification decodes a JSON notification in the request body and sends it to the
// notification initialisation goroutine via the notification channel
// TODO: Respect the "propagate" boolean
func handleNotification(w http.ResponseWriter, r *http.Request, propagate bool) {
	var notif notifier.Notification

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&notif); err != nil {
		log.Debug("Decode error: ", err)
		http.Error(w, "", 400)
		return
	}
	notificationC <- &notif

	w.WriteHeader(200)
}

// triggerNotification triggers a "new" notification, i.e. it has not been passed on from another
// gnotify node, it's been triggered as a new event
func triggerNotification(w http.ResponseWriter, r *http.Request) {
	// FIXME: Should override the notification "source" here, too. = REST
	handleNotification(w, r, true)
}

// routeNotification triggers handling of a notification which has been routed from another node
func routeNotification(w http.ResponseWriter, r *http.Request) {
	handleNotification(w, r, false)
}

func Listen(notC chan<- *notifier.Notification) {
	notificationC = notC

	http.HandleFunc("/notify/trigger/", triggerNotification)
	http.HandleFunc("/notify/route/", routeNotification)

	serviceUrl := fmt.Sprintf("%s:%d", config.Sources.Rest.Host, config.Sources.Rest.Port)

	log.Notice("Starting REST interface on %s...", serviceUrl)
	log.Fatal(http.ListenAndServe(serviceUrl, nil))
}
