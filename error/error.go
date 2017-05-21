package error

import (
	"github.com/getsentry/raven-go"
	"log"
	"github.com/fetzi/styx/config"
	"os"
	"strings"
	"fmt"
	"time"
)

type Severity string

var conf config.LoggingConfig

const (
	DEBUG = Severity("debug")
	INFO = Severity("info")
	WARNING = Severity("warning")
	ERROR = Severity("error")
	FATAL = Severity("fatal")
)

// Setup default config for logging
func init() {
	conf = config.LoggingConfig{
		Sentry: config.SentryConfig{
			Level:"error",
		},
		File: config.FileLoggingConfig{
			Level:"error",
		},
	}
}

// Register custom logging config
func Register(config config.LoggingConfig) {
	conf = config
	raven.SetDSN(config.Sentry.DSN)
}

// Log an error to multiple destinations
// error.LogError(errors.New("this is a warning"), error.WARNING, map[string]string{
// "foo": "bar",
// })
// error.LogError(errors.New("this is fatal"), error.FATAL, nil)
func LogError(err error, level Severity, tags map[string]string) {

	if (conf.Sentry.DSN != "" && shouldLog(Severity(conf.Sentry.Level), level)) {
		logToSentry(err, level, tags)
	}

	if (conf.File.Path != "" && shouldLog(Severity(conf.File.Level), level) ) {
		logToFile(err, level)
	}

	message := fmt.Sprintf("[%s]: %s", strings.ToUpper(fmt.Sprintf("%s", level)), err.Error())
	if (level == FATAL) {
		log.Fatal(message)
	} else {
		log.Println(message)
	}
}

// Log error to sentry
func logToSentry(err error, level Severity, tags map[string]string) {
	client := raven.DefaultClient

	packet := raven.NewPacket(err.Error())
	if (level == FATAL || level == ERROR) {
		packet = raven.NewPacket(err.Error(), raven.NewException(err, raven.NewStacktrace(1, 3, nil)))
	}

	packet.Level = raven.Severity(level)
	packet.AddTags(tags)
	client.Capture(packet, nil)
}

// Log error to filesystem
func logToFile(err error, level Severity) {
	f, e := os.OpenFile(conf.File.Path, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if e != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	message := fmt.Sprintf("%s [%s]: %s\r\n", time.Now(), strings.ToUpper(fmt.Sprintf("%s", level)), err)
	f.WriteString(message)
}

// Should a error be logged based on the max log level from config?
func shouldLog(logLevel Severity, level Severity) bool {
	return levelToInt(logLevel) <= levelToInt(level)
}

// Convert a log level to integer for comparison
func levelToInt(level Severity) int {
	switch level{
	case FATAL:
		return 5
	case ERROR:
		return 4
	case WARNING:
		return 3
	case INFO:
		return 2
	case DEBUG:
		return 1
	default:
		return 0
	}
}
