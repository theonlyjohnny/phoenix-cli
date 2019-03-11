package logger

import (
	"io"
	"net/http"
	"runtime"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

//Logger is the base struct that exports all the functions used
//to send messages
type Logger struct {
	baseLogger *logrus.Logger
	AppName    string
}

//GetHTTPHandler wraps a http.Handler to log every incoming request
func (logger Logger) GetHTTPHandler(handler http.Handler) http.Handler {
	return &logHandler{
		logger: logger,
		next:   handler,
	}
}

//GetWriter exposes a generic *io.PipeWriter into which logs can be written
func (logger Logger) GetWriter() *io.PipeWriter {
	return logger.baseLogger.Writer()
}

//Error prints an unformated error log
func (logger Logger) Error(args ...interface{}) {
	logger.log(logrus.ErrorLevel, args...)
}

//Info prints an unformated info log
func (logger Logger) Info(args ...interface{}) {
	logger.log(logrus.InfoLevel, args...)
}

//Debug prints an unformated debug log
func (logger Logger) Debug(args ...interface{}) {
	logger.log(logrus.DebugLevel, args...)
}

//Warn prints an unformated warn log
func (logger Logger) Warn(args ...interface{}) {
	logger.log(logrus.WarnLevel, args...)
}

//Errorf prints an formated error log
func (logger Logger) Errorf(msg string, args ...interface{}) {
	logger.logf(logrus.ErrorLevel, msg, args...)
}

//Infof prints an formated info log
func (logger Logger) Infof(msg string, args ...interface{}) {
	logger.logf(logrus.InfoLevel, msg, args...)
}

//Debugf prints an formated debug log
func (logger Logger) Debugf(msg string, args ...interface{}) {
	logger.logf(logrus.DebugLevel, msg, args...)
}

//Warnf prints an formated warn log
func (logger Logger) Warnf(msg string, args ...interface{}) {
	logger.logf(logrus.WarnLevel, msg, args...)
}

func (logger Logger) log(level logrus.Level, args ...interface{}) {
	baseLogger := logger.baseLogger
	if baseLogger.Level >= level {
		_, file, line, ok := runtime.Caller(2)
		if !ok {
			file = "???"
			line = 0
		}
		splitFile := strings.Split(file, "/")
		file = splitFile[len(splitFile)-1]

		l := baseLogger.WithFields(logrus.Fields{
			"fileline": file + ":" + strconv.Itoa(line),
		})

		switch level {
		case logrus.DebugLevel:
			l.Debug(args...)
		case logrus.InfoLevel:
			l.Info(args...)
		case logrus.ErrorLevel:
			l.Error(args...)
		case logrus.WarnLevel:
			l.Warn(args...)
		}
	}
}

func (logger Logger) logf(level logrus.Level, msg string, args ...interface{}) {
	baseLogger := logger.baseLogger
	if baseLogger.Level >= level {
		_, file, line, ok := runtime.Caller(2)
		if !ok {
			file = "???"
			line = 0
		}
		splitFile := strings.Split(file, "/")
		file = splitFile[len(splitFile)-1]

		l := baseLogger.WithFields(logrus.Fields{
			"fileline": file + ":" + strconv.Itoa(line),
		})

		switch level {
		case logrus.DebugLevel:
			l.Debugf(msg, args...)
		case logrus.InfoLevel:
			l.Infof(msg, args...)
		case logrus.ErrorLevel:
			l.Errorf(msg, args...)
		case logrus.WarnLevel:
			l.Warnf(msg, args...)
		}
	}
}
