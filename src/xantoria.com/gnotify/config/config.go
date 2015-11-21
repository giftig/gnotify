package config

import (
	"io/ioutil"
	"log"
	"time"

	"gopkg.in/yaml.v1"
)

const (
	Master = "MASTER"
	Client = "CLIENT"
)

type loggingConfig struct {
	Type      string `yaml:"type"`
	File      string `yaml:"file"`
	Level     string `yaml:"level"`
	Formatter string `yaml:"formatter"`
}

type sourceConfig struct {
	Rest struct {
		Host      string
		Port      int
		PollFetch time.Duration `yaml:"poll_fetch"`
		Disabled  bool          `yaml:"disabled"`
	}
	Calendar struct {
		DatetimeFormat string
		DateFormat     string
		Disabled       bool `yaml:"disabled"`

		Polling struct {
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
		File     string
		Disabled bool `yaml:"disabled"`
	}
}

type nodeConfig struct {
	Type string `yaml:"type"`
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
type routingConfig struct {
	RecipientID string   `yaml:"recipient_id"`
	Groups      []string `yaml:"recipient_groups"`

	KnownRecipients []struct {
		ID     string   `yaml:"id"`
		Host   string   `yaml:"host"`
		Port   int      `yaml:"port"`
		Groups []string `yaml:"groups"`
	} `yaml:"known_recipients"`

	Master struct {
		Host string
		Port int
	}
}

type startupConfig struct {
	Delay time.Duration `yaml:"delay"`
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
var Node nodeConfig
var Logging loggingConfig
var Sources sourceConfig
var Notifications notificationConfig
var Static staticConfig
var EventTypes eventTypeConfig
var Routing routingConfig
var Persistence persistenceConfig
var Startup startupConfig

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
		Node          *nodeConfig
		Logging       *loggingConfig
		Sources       *sourceConfig
		Notifications *notificationConfig
		StaticConfig  *staticConfig
		EventTypes    *eventTypeConfig `yaml:"event_types"`
		Static        *staticConfig
		Routing       *routingConfig
		Persistence   *persistenceConfig
		Startup       *startupConfig
	}{
		Node:          &Node,
		Logging:       &Logging,
		Sources:       &Sources,
		Notifications: &Notifications,
		EventTypes:    &EventTypes,
		Static:        &Static,
		Routing:       &Routing,
		Persistence:   &Persistence,
		Startup:       &Startup,
	}

	if err = yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	clean()
}

// clean checks for any misconfigurations and fails fast. It also sets some sensible defaults
func clean() {
	// Default poll fetching from the master to 10 mins
	if Sources.Rest.PollFetch == 0 {
		Sources.Rest.PollFetch = 10 * time.Minute
	}

	// Default to server mode
	if Node.Type == "" {
		Node.Type = Master
	}

	if Node.Type != Master && Node.Type != Client {
		log.Fatalf("Node type must be %s or %s, not %s", Master, Client, Node.Type)
	}

	if Node.Type == Master && !Persistence.Persist {
		log.Fatalf("Cannot be a master node but not be db-backed!")
	}
	if Node.Type == Master && Routing.Master.Host != "" {
		log.Fatalf("We are acting as a master but also have a route to one?")
	}
}
