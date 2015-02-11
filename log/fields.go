// Copyright 2013 CoreOS, Inc.
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
	"sync/atomic"
	"time"
)

type Fields map[string]interface{}

func (logger *Logger) fields() Fields {
	fields := Fields{
		"prefix":     logger.prefixFn,      // static field available to all sinks
		"seq":        logger.seqFn,         // auto-incrementing sequence number
		"start_time": logger.createdFn,     // start time of the logger
		"time":       logger.timeMilliFn,   // formatted time of log entry
		"full_time":  logger.timeFn,        // time of log entry
		"rtime":      logger.timeElapsedFn, // relative time of log entry since started
		"pid":        logger.pidFn,         // process id
		"executable": logger.executableFn,  // executable filename
	}

	return fields
}

type fieldFn func() interface{}

func (l *Logger) prefixFn() interface{} {
	return l.prefix
}

func (l *Logger) seqFn() interface{} {
	return atomic.AddUint64(&logger.seq, 1)
}

func (l *Logger) createdFn() interface{} {
	return l.created
}

func (l *Logger) timeMilliFn() interface{} {
	return time.Now().Format(time.StampMilli)
}

func (l *Logger) timeFn() interface{} {
	return time.Now()
}

func (l *Logger) timeElapsedFn() interface{} {
	return time.Since(l.created)
}

func (l *Logger) pidFn() interface{} {
	return os.Getpid()
}

func (l *Logger) executableFn() interface{} {
	l.executable
}
