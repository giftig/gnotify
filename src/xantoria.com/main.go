package main

import (
  "log"
  "os"

  "xantoria.com/config"
)

func main() {
  settingsFile := "settings.yaml"
  if len(os.Args) > 1 {
    settingsFile = os.Args[1]
  }

  if err := config.LoadConfig(settingsFile); err != nil {
    log.Fatal(err)
  }

  config.ConfigureLogger()
}
