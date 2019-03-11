package logger

import (
	"log/syslog"
	"testing"

	"github.com/sirupsen/logrus"
)

const (
	remote string = "localhost:514"
)

func assert(t *testing.T, prop string, expected, actual interface{}) {
	if expected != actual {
		t.Fatalf("%s mismatch. Expected: %#v, actual: %#v", prop, expected, actual)
	}
}

func TestGetAllLogLevels(t *testing.T) {
	appName := t.Name()
	levels := logrus.AllLevels
	for _, level := range levels {
		levelNameBytes, err := level.MarshalText()
		if err != nil {
			continue
		}

		levelName := string(levelNameBytes)

		config := Config{
			AppName:    appName,
			LogLevel:   levelName,
			LogConsole: true,
			LogSyslog:  nil,
		}

		realConfig, err := validateConfig(config)
		if err != nil {
			t.Fatal(err)
		}

		assert(t, appName+"["+levelName+"]", level, realConfig.logLevel)
	}
}

func TestGetValidConsoleConfig(t *testing.T) {
	appName := t.Name()
	config := Config{
		appName,
		"debug",
		true,
		nil,
	}

	realConfig, err := validateConfig(config)
	// config, err := NewConfig(appName, "debug", true, nil)
	t.Logf("TestGetValidConfig config: %#v, realConfig: %#v", config, realConfig)
	if err != nil {
		t.Fatal(err)
	}

	var logSyslog *syslogConfig
	assert(t, "config.appName", appName, realConfig.appName)
	assert(t, "config.logLevel", logrus.DebugLevel, realConfig.logLevel)
	assert(t, "config.logConsole", true, realConfig.logConsole)
	assert(t, "config.logSyslog", logSyslog, realConfig.logSyslog)
}

func TestGetInvalidConsoleConfig(t *testing.T) {
	appName := t.Name()
	config := Config{
		appName,
		"debug",
		false,
		nil,
	}
	realConfig, err := validateConfig(config)
	t.Logf("%s config: %#v, realConfig: %#v", appName, config, realConfig)
	if err == nil {
		t.Fatal("err shouldn't be nil")
	}
}

func TestGetValidSyslogConfig(t *testing.T) {
	appName := t.Name()
	syslogConfig := SyslogConfig{
		remote,
		"debug",
	}
	config := Config{
		appName,
		"debug",
		false,
		&syslogConfig,
	}
	realConfig, err := validateConfig(config)
	t.Logf("%s config: %#v, realConfig: %#v", appName, config, realConfig)
	if err != nil {
		t.Fatal(err)
	}
	assert(t, "realConfig.appName", appName, realConfig.appName)
	assert(t, "realConfig.logLevel", logrus.DebugLevel, realConfig.logLevel)
	assert(t, "realConfig.logConsole", false, realConfig.logConsole)
	assert(t, "realConfig.logSyslog.remoteIP", syslogConfig.RemoteIP, realConfig.logSyslog.remoteIP)
	assert(t, "realConfig.logSyslog.logPriority", syslog.LOG_DEBUG, realConfig.logSyslog.logPriority)
}

func TestGetInvalidSyslogConfig(t *testing.T) {
	appName := t.Name()
	syslogConfig := SyslogConfig{}
	config := Config{
		appName,
		"debug",
		false,
		&syslogConfig,
	}

	realConfig, err := validateConfig(config)
	t.Logf("%s config: %#v, realConfig: %#v", appName, config, realConfig)
	if err == nil {
		t.Fatal("err shouldn't be nil")
	}
}

func TestCreateLogger(t *testing.T) {
	appName := t.Name()
	config := Config{
		appName,
		"debug",
		true,
		nil,
	}
	logger, err := CreateLogger(config)
	if err != nil {
		t.Fatal(err)
	}
	logger.Info("info")
	logger.Infof("infof %s %v %#v %d %b", "string", "vString", "complexVString", 5, 5)
	logger.Warn("warn")
	logger.Warnf("warnf %s %v %#v %d %b", "string", "vString", "complexVString", 5, 5)
	logger.Error("error")
	logger.Errorf("errorf %s %v %#v %d %b", "string", "vString", "complexVString", 5, 5)
	logger.Debug("debug")
	logger.Debugf("debugf %s %v %#v %d %b", "string", "vString", "complexVString", 5, 5)
}

func TestGetValidLogger(t *testing.T) {
	appName := t.Name()
	TestCreateLogger(t)
	// var nilLogger *Logger
	logger := GetLoggerByAppName(appName)
	if logger == nil {
		t.Fatalf("no returned logger")
	}
}
