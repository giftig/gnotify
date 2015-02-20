package calendar

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"time"

	"code.google.com/p/goauth2/oauth"
	"github.com/skratchdot/open-golang/open"

	"xantoria.com/gnotify/config"
	"xantoria.com/gnotify/log"
	"xantoria.com/gnotify/notifier"
)

// These structs represent JSON received from the google calendar API
type CalendarEvents struct {
	Kind, Updated string
	Items         []CalendarEvent
}
type CalendarEvent struct {
	Id, Status, Summary string
	Start, End          CalendarDate
}
type CalendarDate struct {
	Date     string
	Datetime string `json:"dateTime"`
}

// authenticate auths the service with google via OAuth
func authenticate() (transport *oauth.Transport) {
	googleConfig := config.Auth.Google
	code := googleConfig.Account.Code

	// Configure and create the OAuth Transport
	oauthConfig := &oauth.Config{
		ClientId:     googleConfig.ClientID,
		ClientSecret: googleConfig.Secret,
		RedirectURL:  googleConfig.RedirectURI,
		Scope:        googleConfig.Scope,
		AuthURL:      googleConfig.AuthEndpoint,
		TokenURL:     googleConfig.TokenEndpoint,
		TokenCache:   oauth.CacheFile("_oauth_cache.json"),
	}
	transport = &oauth.Transport{Config: oauthConfig}

	token, err := oauthConfig.TokenCache.Token()

	// We don't have a cached token: we'll need to request one
	if err != nil {
		if code == "" {
			log.Warning("The account code needs to be set in the config.")
			open.Run(oauthConfig.AuthCodeURL(""))
			return
		}
		if token, err = transport.Exchange(code); err != nil {
			log.Fatal("OAuth token exchange failed", err)
		}
	}

	transport.Token = token
	return
}

// GetCalendar connects to google calendar, synchronises notifications based on the calendar
// contents, and pushes new notifications to the provided channel
func GetCalendar(notifs chan *notifier.Notification) {
	log.Notice(
		"Fetching calendar and caching notifications (calendar %s)",
		config.Auth.Google.Account.CalendarID,
	)

	transport := authenticate()
	now := url.QueryEscape(time.Now().Format(config.DatetimeFormat))

	// Get future events
	r, err := transport.Client().Get(fmt.Sprintf(
		"https://www.googleapis.com/calendar/v3/calendars/%s/events?"+
			"alwaysIncludeEmail=false&"+
			"maxAttendees=1&"+
			"timeMin=%s&"+
			"timeZone=UTC",
		config.Auth.Google.Account.CalendarID,
		now,
	))
	if err != nil {
		log.Fatal("Request failed: ", err)
	}
	defer r.Body.Close()

	responseText, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Error reading response body: %v", err)
	}
	log.Debug(
		"REQUEST: %s",
		fmt.Sprintf(
			"https://www.googleapis.com/calendar/v3/calendars/%s/events?"+
				"alwaysIncludeEmail=false&"+
				"maxAttendees=1&"+
				"timeMin=%s&"+
				"timeZone=UTC",
			config.Auth.Google.Account.CalendarID,
			now,
		),
	)

	var data CalendarEvents
	json.Unmarshal(responseText, &data)

	log.Info("Retrieved %d events", len(data.Items))

	for _, event := range data.Items {
		// Detect date or datetime fields for the event and pick the right format to parse
		var rawTime, timeFormat string
		if event.Start.Datetime != "" {
			rawTime = event.Start.Datetime
			timeFormat = config.DatetimeFormat
		} else {
			rawTime = event.Start.Date
			timeFormat = config.DateFormat
		}

		eventTime, err := time.Parse(timeFormat, rawTime)
		if err != nil {
			log.Error("Couldn't parse time '%s'! Skipping event %s", rawTime, event.Id)
			continue
		}

		log.Info("Created notification for event %s (%s)", event.Id, eventTime)
		notif := notifier.Notification{
			Title:    "Calendar event",
			Message:  event.Summary,
			Icon:     config.EventTypes.Calendar.Icon,
			Source:   config.EventTypes.Calendar.Label,
			Id:       event.Id,
			Time:     eventTime,
			Complete: event.Status == "complete",
		}
		notifs <- &notif
	}
}

func LoadEvents(
	ticks <-chan time.Time,
	notificationChannel chan *notifier.Notification,
) {
	for {
		GetCalendar(notificationChannel)
		_ = <-ticks
	}
}
