// Package logger provides functionality for logging
package loggers

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

type AppLogger struct {
	Hostname string
	*logrus.Logger
}

// Get returns logger instance. App will exit if an error occurred while getting logger
func Get() AppLogger {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	return AppLogger{hostname, logrus.StandardLogger()}
}

func (l AppLogger) Info(args ...interface{}) {
	l.addFields().Info(args...)
}

func (l AppLogger) Infof(format string, args ...interface{}) {
	l.addFields().Infof(format, args...)
}

func (l AppLogger) Debug(args ...interface{}) {
	l.addFields().Debug(args...)
}

func (l AppLogger) Debugf(format string, args ...interface{}) {
	l.addFields().Debugf(format, args...)
}

// DebugCtxf is error logging with context
func (l AppLogger) DebugCtxf(context string, format string, args ...interface{}) {
	l.addContextFields(context).Debugf(format, args...)
}

func (l AppLogger) Error(args ...interface{}) {
	l.addFields().Error(args...)
}

func (l AppLogger) Errorf(format string, args ...interface{}) {
	l.addFields().Errorf(format, args...)
}

// ErrorCtx is error logging with context
func (l AppLogger) ErrorCtx(context string, args ...interface{}) {
	l.addContextFields(context).Error(args...)
}

// ErrorCtxf is formatted error logging with context
func (l AppLogger) ErrorCtxf(context string, format string, args ...interface{}) {
	l.addContextFields(context).Errorf(format, args...)
}

// addFields add additional info to log
func (l AppLogger) addFields() *logrus.Entry {
	file, line := getCaller()
	return l.Logger.WithField("hostname", l.Hostname).
		WithField("source", fmt.Sprintf("%s:%d", file, line))
}

func (l AppLogger) addContextFields(context string) *logrus.Entry {
	file, line := getCaller()
	return l.Logger.WithField("hostname", l.Hostname).
		WithField("source", fmt.Sprintf("%s:%d", file, line)).
		WithField("context", context)
}

// getCaller returns where a code is executed
func getCaller() (string, int) {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		s := strings.LastIndex(file, "/")
		file = file[s+1:]
	}
	return file, line
}
