package notifier

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"xantoria.com/gnotify/config"
	"xantoria.com/gnotify/log"
)

// Save persists the notification by saving it to CouchDB
func (notif *Notification) Save() {
	cfg := config.Persistence
	if !cfg.Persist {
		return
	}
	docId := fmt.Sprintf("%s:%s", notif.Id, notif.Source)
	log.Info("Saving notification %s to couch", docId)

	url := fmt.Sprintf(
		"http://%s:%d/%s/%s",
		cfg.Couch.Host,
		cfg.Couch.Port,
		cfg.Couch.Db,
		docId,
	)

	data, err := json.Marshal(notif)
	if err != nil {
		log.Error("Failed to marshal notification for doc %s", docId)
		return
	}

	client := &http.Client{}
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		log.Error("Failed to insert notification %s into couch", docId)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		// A 409 isn't too unexpected: it means we tried to save a notification which already exists
		if resp.StatusCode == 409 {
			log.Warning("Tried to overwrite notification %s", docId)
			return
		}
		log.Error(
			"Unexpected status code from couch when inserting notification %s: %s",
			docId, resp.Status,
		)
	}
}

// MarkComplete calls the relevant couch update handler to mark a notification as complete
// Takes the document ID in couch as the only argument. Returns an error if the request failed
func MarkComplete(docId string) (err error) {
	log.Info("Marking notification %q complete in couch", docId)

	cfg := config.Persistence
	u := fmt.Sprintf(
		"http://%s:%d/%s/_design/notifications/_update/mark_complete/%s",
		cfg.Couch.Host,
		cfg.Couch.Port,
		cfg.Couch.Db,
		docId,
	)
	log.Debug("Hitting URL %s", u)

	resp, err := http.Post(u, "text/plain", nil)

	if err != nil {
		log.Error("A problem occurred marking %q as complete: %q", docId, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		if resp.StatusCode != 404 {
			log.Error("Couch returned an HTTP %d when marking %q complete", resp.StatusCode, docId)
		}

		err = errors.New(strconv.Itoa(resp.StatusCode))
		return
	}
	return nil
}

func (notif *Notification) MarkComplete() {
	master := config.Routing.Master
	// We're the master, so just call the internal util
	if master.Host == "" {
		docId := fmt.Sprintf("%s:%s", notif.Id, notif.Source)
		MarkComplete(docId)
		return
	}

	// Inform the master that our notification has been completed
	u := fmt.Sprintf(
		"http://%s:%d/notify/complete/?%s",
		master.Host,
		master.Port,
		fmt.Sprintf("id=%s&src=%s", url.QueryEscape(notif.Id), url.QueryEscape(notif.Source)),
	)
	log.Info("Marking %s as complete with master (%s)", notif.Id, master.Host)
	log.Debug("Hitting URL %s", u)

	resp, err := http.Post(u, "text/plain", nil)

	if err != nil {
		log.Error("Error while telling the master node that %q is complete: %q", notif.Id, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		log.Error(
			"Bad status code from master while marking %q complete: %d", notif.Id, resp.StatusCode,
		)
		return
	}
}

// Fetch contacts the master to retrieve notifications for the given destination
func Fetch(dest string) (results []Notification) {
	master := config.Routing.Master
	// We're the master, so we don't need to load up notifications
	if master.Host == "" {
		return
	}

	u := fmt.Sprintf(
		"http://%s:%d/notify/fetch/?%s",
		master.Host,
		master.Port,
		fmt.Sprintf("dest=%s", url.QueryEscape(dest)),
	)
	log.Debug("Hitting URL %s", u)

	resp, err := http.Get(u)

	if err != nil {
		log.Error("Error while fetching notifications from master: %q", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		log.Error("Bad status code from master while fetching notifications: %d", resp.StatusCode)
		return
	}

	json.NewDecoder(resp.Body).Decode(&results)
	return
}
