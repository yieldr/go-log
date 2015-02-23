// Copyright 2013 CoreOS, Inc.
// Copyright 2014 Yieldr
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package log

import (
	"fmt"
	"os"
	"time"
)

type Logger interface {
	Emergency(...interface{})
	Emergencyf(string, ...interface{})
	Alert(...interface{})
	Alertf(string, ...interface{})
	Critical(...interface{})
	Criticalf(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
	Warning(...interface{})
	Warningf(string, ...interface{})
	Notice(...interface{})
	Noticef(string, ...interface{})
	Info(...interface{})
	Infof(string, ...interface{})
	Debug(...interface{})
	Debugf(string, ...interface{})
}

var (
	BasicFormat  = "%s [%s] %s\n"
	BasicFields  = []string{"time", "priority", "message"}
	RichFormat   = "%s [%9s] %d %s - %s\n"
	RichFields   = []string{"time", "priority", "seq", "prefix", "message"}
	SyslogFormat = "[%s] %s\n"
	SyslogFields = []string{"priority", "message"}
)

// logger represents an active logging object that forwards log messages to its
// underlying sinks.
type logger struct {
	sinks   []Sink    // the sinks this logger will log to
	prefix  string    // static field available to all log sinks under this logger
	created time.Time // time when this logger was created
	seq     uint64    // sequential number of log message, starting at 1
}

// New creates a new logger with the supplied options.
func New(sinks ...Sink) Logger {
	logger := &logger{
		created: time.Now(),
		seq:     0,
		sinks:   sinks,
	}
	return logger
}

func (logger *logger) Log(p Priority, v ...interface{}) {
	fields := Fields{
		"priority":     func() interface{} { return p },
		"message":      func() interface{} { return fmt.Sprint(v...) },
		"prefix":       logger.prefixFn,  // static field available to all sinks
		"time":         logger.timeFn,    // formatted time of log entry
		"start_time":   logger.createdFn, // start time of the logger
		"elapsed_time": logger.elapsedFn, // relative time of log entry since started
		"seq":          logger.seqFn,     // auto-incrementing sequence number
		"pid":          logger.pidFn,     // process id
	}
	for _, sink := range logger.sinks {
		sink.Log(fields)
	}
}

func (logger *logger) Logf(priority Priority, format string, v ...interface{}) {
	logger.Log(priority, fmt.Sprintf(format, v...))
}

func (logger *logger) Emergency(v ...interface{}) {
	logger.Log(EMERGENCY, v...)
}

func (logger *logger) Emergencyf(format string, v ...interface{}) {
	logger.Log(EMERGENCY, fmt.Sprintf(format, v...))
}

func (logger *logger) Alert(v ...interface{}) {
	logger.Log(ALERT, v...)
}

func (logger *logger) Alertf(format string, v ...interface{}) {
	logger.Log(ALERT, fmt.Sprintf(format, v...))
}

func (logger *logger) Critical(v ...interface{}) {
	logger.Log(CRITICAL, v...)
}

func (logger *logger) Criticalf(format string, v ...interface{}) {
	logger.Log(CRITICAL, fmt.Sprintf(format, v...))
}

func (logger *logger) Error(v ...interface{}) {
	logger.Log(ERROR, v...)
}

func (logger *logger) Errorf(format string, v ...interface{}) {
	logger.Log(ERROR, fmt.Sprintf(format, v...))
}

func (logger *logger) Warning(v ...interface{}) {
	logger.Log(WARNING, v...)
}

func (logger *logger) Warningf(format string, v ...interface{}) {
	logger.Log(WARNING, fmt.Sprintf(format, v...))
}

func (logger *logger) Notice(v ...interface{}) {
	logger.Log(NOTICE, v...)
}

func (logger *logger) Noticef(format string, v ...interface{}) {
	logger.Log(NOTICE, fmt.Sprintf(format, v...))
}

func (logger *logger) Info(v ...interface{}) {
	logger.Log(INFO, v...)
}

func (logger *logger) Infof(format string, v ...interface{}) {
	logger.Log(INFO, fmt.Sprintf(format, v...))
}

func (logger *logger) Debug(v ...interface{}) {
	logger.Log(DEBUG, v...)
}

func (logger *logger) Debugf(format string, v ...interface{}) {
	logger.Log(DEBUG, fmt.Sprintf(format, v...))
}

// This logger can be user directly by just importing the package and using it
// the same way you would use the standard library's log package.
var stdout = New(WriterSink(os.Stdout, BasicFormat, BasicFields)).(*logger)

func Emergency(v ...interface{}) {
	stdout.Log(EMERGENCY, v...)
}

func Emergencyf(format string, v ...interface{}) {
	stdout.Log(EMERGENCY, fmt.Sprintf(format, v...))
}

func Alert(v ...interface{}) {
	stdout.Log(ALERT, v...)
}

func Alertf(format string, v ...interface{}) {
	stdout.Log(ALERT, fmt.Sprintf(format, v...))
}

func Critical(v ...interface{}) {
	stdout.Log(CRITICAL, v...)
}

func Criticalf(format string, v ...interface{}) {
	stdout.Log(CRITICAL, fmt.Sprintf(format, v...))
}

func Error(v ...interface{}) {
	stdout.Log(ERROR, v...)
}

func Errorf(format string, v ...interface{}) {
	stdout.Log(ERROR, fmt.Sprintf(format, v...))
}

func Warning(v ...interface{}) {
	stdout.Log(WARNING, v...)
}

func Warningf(format string, v ...interface{}) {
	stdout.Log(WARNING, fmt.Sprintf(format, v...))
}

func Notice(v ...interface{}) {
	stdout.Log(NOTICE, v...)
}

func Noticef(format string, v ...interface{}) {
	stdout.Log(NOTICE, fmt.Sprintf(format, v...))
}

func Info(v ...interface{}) {
	stdout.Log(INFO, v...)
}

func Infof(format string, v ...interface{}) {
	stdout.Log(INFO, fmt.Sprintf(format, v...))
}

func Debug(v ...interface{}) {
	stdout.Log(DEBUG, v...)
}

func Debugf(format string, v ...interface{}) {
	stdout.Log(DEBUG, fmt.Sprintf(format, v...))
}

// Standard library log functions

func (logger *logger) Fatalln(v ...interface{}) {
	logger.Log(CRITICAL, v...)
	os.Exit(1)
}

func (logger *logger) Fatalf(format string, v ...interface{}) {
	logger.Logf(CRITICAL, format, v...)
	os.Exit(1)
}

func (logger *logger) Panicln(v ...interface{}) {
	s := fmt.Sprint(v...)
	logger.Log(ERROR, s)
	panic(s)
}

func (logger *logger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	logger.Log(ERROR, s)
	panic(s)
}

func (logger *logger) Println(v ...interface{}) {
	logger.Log(INFO, v...)
}

func (logger *logger) Printf(format string, v ...interface{}) {
	logger.Logf(INFO, format, v...)
}

func Fatalln(v ...interface{}) {
	stdout.Log(CRITICAL, v...)
	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	stdout.Logf(CRITICAL, format, v...)
	os.Exit(1)
}

func Panicln(v ...interface{}) {
	s := fmt.Sprint(v...)
	stdout.Log(ERROR, s)
	panic(s)
}

func Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	stdout.Log(ERROR, s)
	panic(s)
}

func Println(v ...interface{}) {
	stdout.Log(INFO, v...)
}

func Printf(format string, v ...interface{}) {
	stdout.Logf(INFO, format, v...)
}
