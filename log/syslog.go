// +build !windows

package log

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

import (
	"fmt"
	"log/syslog"
)

type syslogSink struct {
	w        *syslog.Writer
	priority Priority
	format   string
	fields   []string
}

func (sink *syslogSink) Log(fields Fields) {
	vals := make([]interface{}, len(sink.fields))
	for i, field := range sink.fields {
		if fn, ok := fields[field]; ok {
			vals[i] = fn()
		} else {
			vals[i] = "???"
		}
	}
	msg := fmt.Sprintf(sink.format, vals...)
	switch fields["priority"]().(Priority) {
	case EMERGENCY:
		sink.w.Emerg(msg)
	case ALERT:
		sink.w.Alert(msg)
	case CRITICAL:
		sink.w.Crit(msg)
	case ERROR:
		sink.w.Err(msg)
	case WARNING:
		sink.w.Warning(msg)
	case NOTICE:
		sink.w.Notice(msg)
	case INFO:
		sink.w.Info(msg)
	case DEBUG:
		sink.w.Debug(msg)
	default:
		sink.w.Err(msg)

	}
}

func (sink *syslogSink) Write(b []byte) (int, error) {
	return sink.w.Write(b)
}

func (sink *syslogSink) Close() error {
	return sink.w.Close()
}

// SyslogSink returns a sink that outputs to the local syslog daemon.
func SyslogSink(p Priority, tag, format string, fields []string) (*syslogSink, error) {
	prio := syslog.Priority(p) | syslog.LOG_USER
	w, err := syslog.New(prio, tag)
	if err != nil {
		return nil, err
	}
	return &syslogSink{w, p, format, fields}, nil
}
