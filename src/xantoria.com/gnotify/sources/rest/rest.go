package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/nu7hatch/gouuid"

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

	// Override source if provided
	if source != "" {
		notif.Source = source
	}

	// Create an ID if one wasn't provided
	if notif.Id == "" {
		id, _ := uuid.NewV4()
		notif.Id = id.String()
	}

	notificationC <- &notif

	w.WriteHeader(202)
	io.WriteString(w, notif.Id)
}

// triggerNotification triggers a "new" notification, i.e. it has not been passed on from another
// gnotify node, it's been triggered as a new event
func triggerNotification(w http.ResponseWriter, r *http.Request) {
	handleNotification(w, r, config.EventTypes.Rest.Label)
}

// routeNotification triggers handling of a notification which has been routed from another node
func routeNotification(w http.ResponseWriter, r *http.Request) {
	handleNotification(w, r, "")
}

// fetchNotifications looks in the database for notifications addressed to the provided
// destination and returns them to the client. This is for internal use: generally a node will
// ask an authoritative server, backed by couchdb, if it's had any messages recently
func fetchNotifications(w http.ResponseWriter, r *http.Request) {
	route := config.Routing.Master
	persist := config.Persistence

	// If we're not the master, direct them to where we think the master is
	if route.Host != "" {
		http.Redirect(w, r, fmt.Sprintf("%s:%d/notify/fetch/", route.Host, route.Port), 301)
		return
	}

	dest := r.FormValue("dest")
	if dest == "" {
		http.Error(w, "", 400)
		return
	}
	log.Info("Getting notifications for %q", dest)

	// TODO: Pagination on this endpoint would be nice
	u := fmt.Sprintf(
		"http://%s:%d/%s/_design/notifications/_view/pending_by_recipient_and_time?%s",
		persist.Couch.Host,
		persist.Couch.Port,
		persist.Couch.Db,
		fmt.Sprintf(
			"reduce=false&include_docs=true&startkey=%s&endkey=%s",
			url.QueryEscape(fmt.Sprintf("[%q, null]", dest)),
			url.QueryEscape(fmt.Sprintf("[%q, {}]", dest)),
		),
	)
	log.Debug("Hitting URL %s", u)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		log.Error("A problem occurred getting notifications for %q: %q", dest, err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		log.Error("Couch returned an HTTP %d when getting notifications for %q", resp.StatusCode, dest)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	var data struct {
		Rows []struct {
			Doc notifier.Notification `json:"doc"`
		} `json:"rows"`
	}
	json.NewDecoder(resp.Body).Decode(&data)

	docs := []notifier.Notification{}
	for _, row := range data.Rows {
		docs = append(docs, row.Doc)
	}

	log.Info("Found %d pending notifications for %s", len(docs), dest)

	result, err := json.Marshal(docs)
	if err != nil {
		log.Error("Could not marshal list of notifications from couch into JSON: %q, %q", docs, err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, string(result))
}

// completeNotification marks a given notification as completed
func completeNotification(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	// If we're not the master, 404
	if config.Node.Type != config.Master {
		http.Error(w, "", 404)
		return
	}

	id := r.FormValue("id")
	src := r.FormValue("src")

	if id == "" || src == "" {
		http.Error(w, "", 400)
		return
	}
	docId := url.QueryEscape(fmt.Sprintf("%s:%s", id, src))
	log.Info("Notification %s acknowledged by %s", docId, r.RemoteAddr)
	err := notifier.MarkComplete(docId)

	if err != nil {
		if err.Error() == "404" {
			http.Error(w, "Not Found", 404)
		} else {
			http.Error(w, "Internal Server Error", 500)
		}
		return
	}
	w.WriteHeader(200)
}

func Listen(notC chan<- *notifier.Notification) {
	notificationC = notC

	http.HandleFunc("/notify/trigger/", triggerNotification)
	http.HandleFunc("/notify/route/", routeNotification)
	http.HandleFunc("/notify/fetch/", fetchNotifications)
	http.HandleFunc("/notify/complete/", completeNotification)

	serviceUrl := fmt.Sprintf("%s:%d", config.Sources.Rest.Host, config.Sources.Rest.Port)

	log.Notice("Starting REST interface on %s...", serviceUrl)
	log.Fatal(http.ListenAndServe(serviceUrl, nil))
}
