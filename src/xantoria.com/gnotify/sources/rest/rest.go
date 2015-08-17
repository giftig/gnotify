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
func handleNotification(w http.ResponseWriter, r *http.Request, source string) {
	var notif notifier.Notification

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&notif); err != nil {
		log.Debug("Decode error: ", err)
		http.Error(w, "", 400)
		return
	}

	// Override source is provided
	if source != "" {
		notif.Source = source
	}
	notificationC <- &notif

	w.WriteHeader(202)
}

// triggerNotification triggers a "new" notification, i.e. it has not been passed on from another
// gnotify node, it's been triggered as a new event
func triggerNotification(w http.ResponseWriter, r *http.Request) {
	handleNotification(w, r, "API")
}

// routeNotification triggers handling of a notification which has been routed from another node
func routeNotification(w http.ResponseWriter, r *http.Request) {
	handleNotification(w, r, "")
}

// fetchNotifications looks in the database for notifications addressed to the provided
// destination and returns them to the client. This is for internal use: generally a node will
// ask an authoritative server, backed by couchdb, if it's had any messages recently
func fetchNotifications(w http.ResponseWriter, r *http.Request) {
	cfg := config.Routing.Master
	// If we're not the master, direct them to where we think the master is
	if cfg.Host != "" {
		http.Redirect(w, r, fmt.Sprintf("%s:%d/notify/fetch/", cfg.Host, cfg.Port), 301)
		return
	}

	dest := r.FormValue("dest")
	if dest == "" {
		http.Error(w, "", 400)
		return
	}

	// FIXME: Look up notifications for `dest` in couch
	http.Error(w, fmt.Sprintf("Not yet implemented: cannot retrieve for %s", dest), 500)
	return
}

func Listen(notC chan<- *notifier.Notification) {
	notificationC = notC

	http.HandleFunc("/notify/trigger/", triggerNotification)
	http.HandleFunc("/notify/route/", routeNotification)
	http.HandleFunc("/notify/fetch/", fetchNotifications)

	serviceUrl := fmt.Sprintf("%s:%d", config.Sources.Rest.Host, config.Sources.Rest.Port)

	log.Notice("Starting REST interface on %s...", serviceUrl)
	log.Fatal(http.ListenAndServe(serviceUrl, nil))
}
