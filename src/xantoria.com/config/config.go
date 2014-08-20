package config

import (
  "io/ioutil"
  "log"
  "os"
  "time"

  "gopkg.in/yaml.v1"
)

type YAMLConfig struct {
  Auth AuthConfig
  Polling PollingConfig
  Logging LoggingConfig
  DatetimeFormat string `yaml:"datetime_format"`
  DateFormat string `yaml:"date_format"`
}

type LoggingConfig struct {
  Type string
  File string
}

type PollingConfig struct {
  Sync time.Duration
}

type AuthConfig struct {
  Google GoogleAuthConfig
}
type GoogleAuthConfig struct {
  ClientID string `yaml:"client_id"`
  Secret string
  AuthEndpoint string `yaml:"auth_endpoint"`
  TokenEndpoint string `yaml:"token_endpoint"`
  RedirectURI string `yaml:"redirect_uri"`
  Scope string
  Account GoogleAccountConfig
}
type GoogleAccountConfig struct {
  Code string
  CalendarID string `yaml:"calendar_id"`
}

var Config YAMLConfig

/**
 * Load config from the given file and stick it into Config
 */
func LoadConfig(file string) error {
  data, err := ioutil.ReadFile(file)
  if err != nil { return err }

  return yaml.Unmarshal(data, &Config)
}

/**
 * Configure the Logger based on logging config
 */
func ConfigureLogger() {
  loggerConfig := Config.Logging
  if loggerConfig.Type == "" {
    loggerConfig = LoggingConfig{"console", ""}
  }

  switch loggerConfig.Type {
    case "console":
      break  // Console is the default logging configuration anyway
    case "file":
      f, err := os.OpenFile(
        loggerConfig.File,
        os.O_RDWR | os.O_CREATE | os.O_APPEND,
        0644,
      )
      if err != nil {
        log.Fatalf("The logfile %s could not be accessed", loggerConfig.File)
      }
      log.SetOutput(f)
  }
}
