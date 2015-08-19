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
	Todo struct {
		File string
	}
}

type notificationConfig struct {
	NotifySend notifySendConfig `yaml:"notify_send"`
}
type notifySendConfig struct {
	Duration time.Duration `yaml:"duration"`
	Enabled  bool          `yaml:"enabled"`
}

type eventTypeConfig struct {
	Calendar basicEventConfig
	Todo     basicEventConfig
	Rest     basicEventConfig
}
type basicEventConfig struct {
	Icon, Label string
}

type staticConfig struct {
	IconPath string `yaml:"icon_path"`
}

// Identifies an endpoint on which to inform a named recipient of a message
// FIXME: Refactor to be inline for consistency
type recipientConfig struct {
	ID   string `yaml:"id"`
	Host string
	Port int
}
type routingConfig struct {
	RecipientID     string            `yaml:"recipient_id"`
	Groups          []string          `yaml:"recipient_groups"`
	KnownRecipients []recipientConfig `yaml:"known_recipients"`
	Master          struct {
		Host string
		Port int
	}
}

type persistenceConfig struct {
	Couch struct {
		Host string
		Port int
		Db   string
	}
	Persist bool
}

// Top-level config params
var Logging loggingConfig
var Sources sourceConfig
var Notifications notificationConfig
var Static staticConfig
var EventTypes eventTypeConfig
var Routing routingConfig
var Persistence persistenceConfig

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
		Persistence   *persistenceConfig
	}{
		Logging:       &Logging,
		Sources:       &Sources,
		Notifications: &Notifications,
		EventTypes:    &EventTypes,
		Static:        &Static,
		Routing:       &Routing,
		Persistence:   &Persistence,
	}

	if err = yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	if !Persistence.Persist && Routing.Master.Host == "" {
		log.Fatalf("Improperly configured: cannot be a master node but not be db-backed.")
	}
}
