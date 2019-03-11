package logger

import (
	"errors"
	"log/syslog"

	"github.com/sirupsen/logrus"
)

var levelToPriority = map[string]syslog.Priority{
	"emerg":   syslog.LOG_EMERG, /* system is unusable */
	"alert":   syslog.LOG_ALERT, /* action must be taken immediately */
	"crit":    syslog.LOG_CRIT,  /* critical conditions */
	"err":     syslog.LOG_ERR,   /* error conditions */
	"error":   syslog.LOG_ERR,
	"warning": syslog.LOG_WARNING, /* warning conditions */
	"warn":    syslog.LOG_WARNING, /* warning conditions */
	"notice":  syslog.LOG_NOTICE,  /* normal but significant condition */
	"info":    syslog.LOG_INFO,    /* informational */
	"debug":   syslog.LOG_DEBUG,   /* debug-level messages */
}

//Config is used to configure logger
type Config struct {
	AppName    string
	LogLevel   string
	LogConsole bool
	LogSyslog  *SyslogConfig
}

//SyslogConfig is a sub-config of LoggerConfig for syslog configuration
type SyslogConfig struct {
	RemoteIP    string
	LogPriority string
}

type loggerConfig struct {
	appName    string
	logLevel   logrus.Level
	logConsole bool
	logSyslog  *syslogConfig
}

type syslogConfig struct {
	remoteIP    string
	logPriority syslog.Priority
}

func validateConfig(input Config) (*loggerConfig, error) {
	if input.AppName == "" {
		return nil, errors.New("AppName is required")
	}

	if input.LogLevel == "" {
		return nil, errors.New("LogLevel is required")
	}

	logLvl, err := logrus.ParseLevel(input.LogLevel)
	if err != nil {
		return nil, errors.New("Invalid LogLevel")
	}
	var realLogSyslog *syslogConfig

	if input.LogSyslog != nil {
		if input.LogSyslog.RemoteIP == "" {
			return nil, errors.New("If LogSyslog is specified, LogSyslog.RemoteIP is required")
		}
		if input.LogSyslog.LogPriority == "" {
			return nil, errors.New("If LogSyslog is specified, LogSyslog.LogPriority is required")
		}
		logPrio, ok := levelToPriority[input.LogSyslog.LogPriority]
		if !ok {
			return nil, errors.New("Invalid LogPriority")
		}
		realLogSyslog = &syslogConfig{
			remoteIP:    input.LogSyslog.RemoteIP,
			logPriority: logPrio,
		}
	}

	if input.LogSyslog == nil && !input.LogConsole {
		return nil, errors.New("Must have at least either LogSyslog or LogConsole enabled")
	}
	c := loggerConfig{
		input.AppName,
		logLvl,
		input.LogConsole,
		realLogSyslog,
	}
	return &c, nil
}
