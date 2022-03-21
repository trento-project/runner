package internal

import (
	log "github.com/sirupsen/logrus"
)

func SetLogLevel(level string) {
	switch level {
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.Warnln("Unrecognized minimum log level; using 'info' as default")
		log.SetLevel(log.InfoLevel)
	}
}

func SetLogFormatter(timestampFormat string) {
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = timestampFormat
	log.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true
}
