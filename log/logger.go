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
	"os"
	"time"

	"bitbucket.org/kardianos/osext"
)

// The Option type can be passed to the New function to customise the logger.
type Option func(*Logger)

// Sink adds s to the logger's output sinks.
func Sink(s sink) Option {
	return func(l *Logger) {
		l.sinks = append(l.sinks, s)
	}
}

// Prefix sets the logger's prefix to s.
func Prefix(s string) Option {
	return func(l *Logger) {
		l.prefix = s
	}
}

// Logger is user-immutable immutable struct which can log to several outputs
type Logger struct {
	sinks   []Sink    // the sinks this logger will log to
	prefix  string    // static field available to all log sinks under this logger
	created time.Time // time when this logger was created
	seq     uint64    // sequential number of log message, starting at 1
	exec    string    // executable name
}

// New creates a new logger with the supplied options.
func New(options ...Option) *Logger {
	logger := new(Logger)
	logger.created = time.Now()
	logger.seq = 0
	logger.executable = executable()
	for _, option := range options {
		option(logger)
	}
	return logger
}

var defaultLogger *Logger

func init() {
	defaultLogger = New(Sink(WriterSink(os.Stdout, BasicFormat, BasicFields)))
}

func exec() string {
	path, err := osext.Executable()
	if err != nil {
		return "(UNKNOWN)"
	}
	return path.Base(path)
}
