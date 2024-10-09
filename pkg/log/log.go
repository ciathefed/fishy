package log

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
)

var logger = log.NewWithOptions(os.Stderr, log.Options{
	// ReportCaller:    true,
	ReportTimestamp: true,
	TimeFormat:      time.Kitchen,
	Level:           log.DebugLevel,
})

func Info(msg interface{}, keyVals ...interface{}) {
	logger.Info(msg, keyVals...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Debug(msg interface{}, keyVals ...interface{}) {
	logger.Debug(msg, keyVals...)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Warn(msg interface{}, keyVals ...interface{}) {
	logger.Warn(msg, keyVals...)
}

func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func Error(msg interface{}, keyVals ...interface{}) {
	logger.Error(msg, keyVals...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Fatal(msg interface{}, keyVals ...interface{}) {
	logger.Fatal(msg, keyVals...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}
