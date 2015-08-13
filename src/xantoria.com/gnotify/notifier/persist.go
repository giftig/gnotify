package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

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

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		log.Error(
			"Unexpected status code from couch when inserting notification %s: %s",
			docId, resp.Status,
		)
	}
	defer resp.Body.Close()
}
