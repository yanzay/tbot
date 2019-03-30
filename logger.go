package tbot

import "log"

// Logger defines interface for any compatible logger
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Print(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

type nopLogger struct{}

func (nopLogger) Debugf(format string, args ...interface{}) {}
func (nopLogger) Infof(format string, args ...interface{})  {}
func (nopLogger) Printf(format string, args ...interface{}) {}
func (nopLogger) Warnf(format string, args ...interface{})  {}
func (nopLogger) Errorf(format string, args ...interface{}) {}
func (nopLogger) Debug(args ...interface{})                 {}
func (nopLogger) Info(args ...interface{})                  {}
func (nopLogger) Print(args ...interface{})                 {}
func (nopLogger) Warn(args ...interface{})                  {}
func (nopLogger) Error(args ...interface{})                 {}

type BasicLogger struct{}

func (BasicLogger) Debugf(format string, args ...interface{}) { log.Printf(format, args...) }
func (BasicLogger) Infof(format string, args ...interface{})  { log.Printf(format, args...) }
func (BasicLogger) Printf(format string, args ...interface{}) { log.Printf(format, args...) }
func (BasicLogger) Warnf(format string, args ...interface{})  { log.Printf(format, args...) }
func (BasicLogger) Errorf(format string, args ...interface{}) { log.Printf(format, args...) }
func (BasicLogger) Debug(args ...interface{})                 { log.Print(args...) }
func (BasicLogger) Info(args ...interface{})                  { log.Print(args...) }
func (BasicLogger) Print(args ...interface{})                 { log.Print(args...) }
func (BasicLogger) Warn(args ...interface{})                  { log.Print(args...) }
func (BasicLogger) Error(args ...interface{})                 { log.Print(args...) }
