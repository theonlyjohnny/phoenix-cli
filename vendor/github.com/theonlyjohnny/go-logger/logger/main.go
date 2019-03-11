package logger

import (
	"fmt"
	"io/ioutil"
	"testing"

	logrus_syslog "github.com/sirupsen/logrus/hooks/syslog"
)

func testLog(t *testing.T, base string, args ...interface{}) {
	if t != nil {
		t.Logf(base, args...)
	}
}

//CreateLogger takes in a Config instance and constructs a Logger instance scoped to appName.
//This Logger is accessible through GetLoggerByAppName, as well as the return args of this function
func CreateLogger(config Config) (*Logger, error) {
	return realInit(config, nil)
}

//GetLoggerByAppName returns a reference to a previously defined logger
//with the same appName. Will return nil if the logger was not previously
//constructed (see CreateLogger)
func GetLoggerByAppName(appName string) *Logger {
	return factory.getLogger(appName)
}

func getTestLogger(config Config, t *testing.T) (*Logger, error) {
	return realInit(config, t)
}

func realInit(config Config, t *testing.T) (*Logger, error) {
	if ok := factory.isLoggerPresent(config.AppName); !ok {
		return createNewLogger(config, t)
	}
	logger := factory.getLogger(config.AppName)

	return logger, nil
}

func createNewLogger(config Config, t *testing.T) (*Logger, error) {
	testLog(t, "validating config for logger with name: %s", config.AppName)
	realConfig, err := validateConfig(config)
	if err != nil {
		return nil, err
	}
	logger, err := factory.createLogger(realConfig.appName)
	if err != nil {
		return nil, err
	}

	baseLogger := logger.baseLogger
	baseLogger.Level = realConfig.logLevel

	if realConfig.logSyslog != nil {
		testLog(t, "finalPrio: %d", realConfig.logSyslog.logPriority)
		hook, err := logrus_syslog.NewSyslogHook("udp", realConfig.logSyslog.remoteIP, realConfig.logSyslog.logPriority, realConfig.appName)

		if err == nil {
			testLog(t, "adding syslog %#v to %v -- writer: %#v", hook, baseLogger, hook.Writer)
			baseLogger.Hooks.Add(hook)

		} else {
			return nil, fmt.Errorf("Failed to add syslog %s", err.Error())
		}
	}

	if !realConfig.logConsole {
		baseLogger.Out = ioutil.Discard
	}

	return logger, nil
}
