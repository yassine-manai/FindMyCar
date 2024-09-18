package config

import (
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	//"gopkg.in/natefinch/lumberjack.v2"
)

// https://zerolog.io/#thread-safe-lock-free-non-blocking-writer
var mu sync.Mutex

// InitLogger initializes the logger with the specified configurations.
func InitLogger() {
	// Set up zerolog with the specified configurations
	zerolog.TimeFieldFormat = time.RFC3339Nano

	// Configure console and file loggers
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006-01-02 15:04:05.999"}

	//multiWriter := zerolog.MultiLevelWriter(consoleWriter, fileLogger)

	// Configure logger
	mu.Lock()
	defer mu.Unlock()
	log.Logger = zerolog.New(consoleWriter).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Logger()
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().Local()
	}
	// Customize log format
	log.Logger = log.Logger.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "2006-01-02 15:04:05.999",

		FormatCaller: func(i interface{}) string {
			caller := i.(string)
			parts := strings.Split(caller, "/")
			fileAndLine := parts[len(parts)-1]
			return "\033[34m" + fileAndLine + "\033[0m"
		},
		// FormatFieldName :func(i interface{}) string {
		// 		if i == "LaneID" {
		// 				return "\033[42;30m" + i.(string) + "\033[0m"+": "
		// 		}
		// 		return i.(string)+": "
		// },
	})
}
