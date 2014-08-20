package gcalendar

import "time"

type Notification struct {
  Title, Message, Icon string
  Source, Id string
  Time time.Time
  Complete bool
}
