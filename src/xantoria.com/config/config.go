package config

import (
  "io/ioutil"
  "log"
  "os"
  "time"

  "launchpad.net/goyaml"
)

type YAMLConfig struct {
  Polling PollingConfig
  Logging LoggingConfig
}

type LoggingConfig struct {
  Type string
  File string
}

type PollingConfig struct {
  Notify, Sync *YAMLDuration
}
type YAMLDuration time.Duration

func (duration *YAMLDuration) SetYAML(tag string, value interface{}) bool {
  timeDuration, err := time.ParseDuration(value.(string))
  if err != nil { return false }

  *duration = YAMLDuration(timeDuration)
  return true
}


var Config YAMLConfig

/**
 * Load config from the given file and stick it into Config
 */
func LoadConfig(file string) error {
  data, err := ioutil.ReadFile(file)
  if err != nil { return err }

  return goyaml.Unmarshal(data, &Config)
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
