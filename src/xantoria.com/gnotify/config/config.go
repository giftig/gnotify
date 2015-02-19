package config

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v1"
)

type loggingConfig struct {
	Type string
	File string
}

type pollingConfig struct {
	Sync time.Duration
}

type notificationConfig struct {
	NotifySend notifySendConfig `yaml:"notify_send"`
}
type notifySendConfig struct {
	Duration time.Duration
}

type authConfig struct {
	Google googleAuthConfig
}
type googleAuthConfig struct {
	ClientID      string `yaml:"client_id"`
	Secret        string
	AuthEndpoint  string `yaml:"auth_endpoint"`
	TokenEndpoint string `yaml:"token_endpoint"`
	RedirectURI   string `yaml:"redirect_uri"`
	Scope         string
	Account       googleAccountConfig
}
type googleAccountConfig struct {
	Code       string
	CalendarID string `yaml:"calendar_id"`
}

type eventTypeConfig struct {
	Calendar calendarEventConfig
}
type calendarEventConfig struct {
	Icon, Label string
}

// Top-level config params
var Auth authConfig
var Polling pollingConfig
var Logging loggingConfig
var Notifications notificationConfig
var EventTypes eventTypeConfig
var DatetimeFormat string
var DateFormat string

/**
 * Load config from the given file and stick it into Config
 */
func LoadConfig(file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	// Unmarshal the config directly into package-level structs for each section
	cfg := struct {
		Auth           *authConfig
		Polling        *pollingConfig
		Logging        *loggingConfig
		Notifications  *notificationConfig
		EventTypes     *eventTypeConfig
		DatetimeFormat *string
		DateFormat     *string
	}{
		Auth:           &Auth,
		Polling:        &Polling,
		Logging:        &Logging,
		Notifications:  &Notifications,
		EventTypes:     &EventTypes,
		DatetimeFormat: &DatetimeFormat,
		DateFormat:     &DateFormat,
	}

	return yaml.Unmarshal(data, &cfg)
}

/**
 * Configure the Logger based on logging config
 */
func ConfigureLogger() {
	if Logging.Type == "" {
		Logging = loggingConfig{"console", ""}
	}

	switch Logging.Type {
	case "console":
		break // Console is the default logging configuration anyway
	case "file":
		f, err := os.OpenFile(
			Logging.File,
			os.O_RDWR|os.O_CREATE|os.O_APPEND,
			0644,
		)
		if err != nil {
			log.Fatalf("The logfile %s could not be accessed", Logging.File)
		}
		log.SetOutput(f)
	}
}
