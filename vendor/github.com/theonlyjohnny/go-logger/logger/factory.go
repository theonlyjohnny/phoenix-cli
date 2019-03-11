package logger

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

type loggerFactory struct {
	loggers map[string]Logger
}

var factory = &loggerFactory{
	map[string]Logger{},
}

func (f *loggerFactory) getLogger(appName string) *Logger {
	logger, _ := f.loggers[appName]
	return &logger
}

func (f *loggerFactory) createLogger(appName string) (*Logger, error) {
	if f.isLoggerPresent(appName) {
		return nil, fmt.Errorf("Logger name %s already taken", appName)
	}

	logger := Logger{
		baseLogger: &logrus.Logger{
			Out:       os.Stderr,
			Formatter: new(CustomFormatter),
			Hooks:     make(logrus.LevelHooks),
			Level:     logrus.DebugLevel,
		},
		AppName: appName,
	}

	f.loggers[appName] = logger

	return f.getLogger(appName), nil
}

func (f *loggerFactory) isLoggerPresent(appName string) bool {
	_, ok := f.loggers[appName]
	return ok
}
