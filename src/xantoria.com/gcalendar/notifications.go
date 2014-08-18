package gcalendar

import "time"

type Notification struct {
  Title, Message, Icon string
  Time time.Time
  Complete bool
}
