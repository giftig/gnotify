package log

import (
	"log"
	"os"

	"github.com/op/go-logging"

	"xantoria.com/gnotify/config"
)

var Logger *logging.Logger = logging.MustGetLogger("xangames")

// Expose the logging functions at package level
var Debug = Logger.Debug
var Info = Logger.Info
var Notice = Logger.Notice
var Warning = Logger.Warning
var Error = Logger.Error
var Critical = Logger.Critical
var Fatal = Logger.Fatal
var Fatalf = Logger.Fatalf
var Panic = Logger.Panic
var Panicf = Logger.Panicf

func Init() {
	var logStream *os.File
	var level logging.Level
	var err error
	cfg := config.Logging
	Logger.ExtraCalldepth = 1

	switch cfg.Type {
	case "file":
		if logStream, err = os.OpenFile(cfg.File, os.O_RDWR|os.O_APPEND, 0660); err != nil {
			log.Fatal("Error setting up logger: ", err)
		}

	case "console":
		fallthrough
	case "stdout":
		logStream = os.Stdout

	case "stderr":
		logStream = os.Stderr

	default:
		log.Fatal("Bad logger type: '%s'", cfg.Type)
	}

	backend := logging.NewLogBackend(logStream, "", 0)
	formattedBackend := logging.NewBackendFormatter(
		backend,
		logging.MustStringFormatter(cfg.Formatter),
	)
	levelledBackend := logging.AddModuleLevel(formattedBackend)

	// Select the right logging level based on logging.level in config
	switch cfg.Level {
	case "CRITICAL":
		level = logging.CRITICAL
	case "ERROR":
		level = logging.ERROR
	case "WARNING":
		level = logging.WARNING
	case "NOTICE":
		level = logging.NOTICE
	case "INFO":
		level = logging.INFO
	case "DEBUG":
		level = logging.DEBUG

	default:
		log.Fatalf("Bad logging level: %s", cfg.Level)
	}

	levelledBackend.SetLevel(level, "xangames")
	logging.SetBackend(levelledBackend)
}
