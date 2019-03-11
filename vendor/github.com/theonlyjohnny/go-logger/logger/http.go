package logger

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const formatPattern = "%s %d %d %.4f - %s"

//LogRecord contains information about an incoming HTTP request and its corresponding response
type LogRecord struct {
	http.ResponseWriter
	logger                Logger
	status                int
	responseBytes         int64
	ip                    string
	method, uri, protocol string
	time                  time.Time
	elapsedTime           time.Duration
}

type logHandler struct {
	logger Logger
	next   http.Handler
	io     io.Writer
}

func (h *logHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	if colon := strings.LastIndex(clientIP, ":"); colon != -1 {
		clientIP = clientIP[:colon]
	}

	record := &LogRecord{
		ResponseWriter: w,
		logger:         h.logger,
		ip:             clientIP,
		time:           time.Time{},
		method:         r.Method,
		uri:            r.RequestURI,
		protocol:       r.Proto,
		status:         http.StatusOK,
		elapsedTime:    time.Duration(0),
	}

	startTime := time.Now()
	h.next.ServeHTTP(record, r)
	finishTime := time.Now()

	record.time = finishTime.UTC()
	record.elapsedTime = finishTime.Sub(startTime)

	record.log()
}

func (r *LogRecord) log() {
	request := fmt.Sprintf("%s %s %s", r.method, r.uri, r.protocol)
	line := fmt.Sprintf(formatPattern, request, r.status, r.responseBytes, r.elapsedTime.Seconds(), r.ip)
	if r.status < 300 {
		r.logger.Info(line)
	} else if r.status > 499 {
		r.logger.Error(line)
	} else {
		r.logger.Warn(line)
	}
}

// Write acts like a proxy passing the given bytes buffer to the ResponseWritter
// and additionally counting the passed amount of bytes for logging usage.
func (r *LogRecord) Write(p []byte) (int, error) {
	written, err := r.ResponseWriter.Write(p)
	r.responseBytes += int64(written)
	return written, err
}

// WriteHeader adds the request status onto the LogRecord
func (r *LogRecord) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}
