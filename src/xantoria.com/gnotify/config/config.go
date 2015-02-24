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

type sourceConfig struct {
	Rest struct {
		Host string
		Port int
	}
	Calendar struct {
		DatetimeFormat string
		DateFormat     string
		Polling        struct {
			Sync time.Duration
		}
		Auth struct {
			ClientID      string `yaml:"client_id"`
			Secret        string
			AuthEndpoint  string `yaml:"auth_endpoint"`
			TokenEndpoint string `yaml:"token_endpoint"`
			RedirectURI   string `yaml:"redirect_uri"`
			Scope         string
			Account       struct {
				Code       string
				CalendarID string `yaml:"calendar_id"`
			}
		}
	}
}

type notificationConfig struct {
	NotifySend notifySendConfig `yaml:"notify_send"`
}
type notifySendConfig struct {
	Duration time.Duration
}

type eventTypeConfig struct {
	Calendar calendarEventConfig
}
type calendarEventConfig struct {
	Icon, Label string
}

type staticConfig struct {
	IconPath string `yaml:"icon_path"`
}

// Identifies an endpoint on which to inform a named recipient of a message
type recipientConfig struct {
	ID       string `yaml:"id"`
	Endpoint string `yaml:"endpoint"`
}
type routingConfig struct {
	RecipientID     string            `yaml:"recipient_id"`
	Groups          []string          `yaml:"recipient_groups"`
	KnownRecipients []recipientConfig `yaml:"known_recipients"`
}

// Top-level config params
var Logging loggingConfig
var Sources sourceConfig
var Notifications notificationConfig
var Static staticConfig
var EventTypes eventTypeConfig
var Routing routingConfig

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
		Logging       *loggingConfig
		Sources       *sourceConfig
		Notifications *notificationConfig
		StaticConfig  *staticConfig
		EventTypes    *eventTypeConfig `yaml:"event_types"`
		Static        *staticConfig
		Routing       *routingConfig
	}{
		Logging:       &Logging,
		Sources:       &Sources,
		Notifications: &Notifications,
		EventTypes:    &EventTypes,
		Static:        &Static,
		Routing:       &Routing,
	}

	if err = yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}
}
