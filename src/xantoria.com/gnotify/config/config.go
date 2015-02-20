package config

import (
	"io/ioutil"
	"log"
	"time"

	"gopkg.in/yaml.v1"
)

type loggingConfig struct {
	Type      string `yaml:"type"`
	File      string `yaml:"file"`
	Level     string `yaml:"level"`
	Formatter string `yaml:"formatter"`
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
func LoadConfig(file string) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// Unmarshal the config directly into package-level structs for each section
	cfg := struct {
		Auth           *authConfig
		Polling        *pollingConfig
		Logging        *loggingConfig
		Notifications  *notificationConfig
		EventTypes     *eventTypeConfig `yaml:"event_types"`
		DatetimeFormat *string          `yaml:"datetime_format"`
		DateFormat     *string          `yaml:"date_format"`
	}{
		Auth:           &Auth,
		Polling:        &Polling,
		Logging:        &Logging,
		Notifications:  &Notifications,
		EventTypes:     &EventTypes,
		DatetimeFormat: &DatetimeFormat,
		DateFormat:     &DateFormat,
	}

	if err = yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}
}
