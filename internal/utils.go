package internal

import (
	"context"
	"fmt"
	"time"

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

// Repeat executes a function at a given interval.
// the first tick runs immediately
func Repeat(operation string, tick func(), interval time.Duration, ctx context.Context) {
	tick()

	ticker := time.NewTicker(interval)
	msg := fmt.Sprintf("Next execution for operation %s in %s", operation, interval)
	log.Debugf(msg)

	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			tick()
			log.Debugf(msg)
		case <-ctx.Done():
			return
		}
	}
}
