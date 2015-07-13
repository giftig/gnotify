package todo

import (
	"io/ioutil"
	"os"
	"strings"
	"time"

	"xantoria.com/gnotify/config"
	"xantoria.com/gnotify/log"
	"xantoria.com/gnotify/notifier"
)

func LoadEvents(notificationC chan *notifier.Notification) {
	f := config.Sources.Todo.File
	log.Notice("Checking to do list at %s", f)

	data, err := ioutil.ReadFile(f)
	msg := strings.TrimSpace(string(data[:]))

	if err != nil {
		if os.IsNotExist(err) {
			log.Info("To do list is not defined; nothing to report")
		} else {
			log.Error("Couldn't read to do list:", err)
		}
		return
	}

	if msg == "" {
		log.Info("To do list is empty; nothing to report.")
		return
	}

	notif := notifier.Notification{
		Title:    "To do list!",
		Message:  msg,
		Icon:     config.EventTypes.Todo.Icon,
		Source:   config.EventTypes.Todo.Label,
		Id:       "todo-list",
		Time:     time.Now().Add(time.Duration(10) * time.Second),
		Complete: false,
	}
	notificationC <- &notif
}
