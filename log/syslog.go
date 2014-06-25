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
//
// Author: Alex Kalyvitis <alex.kalyvitis@yieldr.com>
// Author: David Fisher <ddf1991@gmail.com>
// Based on previous package by: Cong Ding <dinggnu@gmail.com>
package log

import (
	"fmt"
	"log/syslog"
)

type syslogSink struct {
	w      *syslog.Writer
	format string
	fields []string
}

func (sink *syslogSink) Log(fields Fields) {
	vals := make([]interface{}, len(sink.fields))
	for i, field := range sink.fields {
		var ok bool
		vals[i], ok = fields[field]
		if !ok {
			vals[i] = "???"
		}
	}
	msg := fmt.Sprintf(sink.format, vals...)
	switch fields["priority"].(Priority) {
	case PriEmerg:
		sink.w.Emerg(msg)
	case PriAlert:
		sink.w.Alert(msg)
	case PriCrit:
		sink.w.Crit(msg)
	case PriErr:
		sink.w.Err(msg)
	case PriWarning:
		sink.w.Warning(msg)
	case PriNotice:
		sink.w.Notice(msg)
	case PriInfo:
		sink.w.Info(msg)
	case PriDebug:
		sink.w.Debug(msg)
	}
}

func (s *syslogSink) Write(b []byte) (int, error) {
	return s.w.Write(b)
}

func SyslogSink(p Priority, format string, fields []string) (*syslogSink, error) {
	prio := syslog.Priority(p) | syslog.LOG_USER
	w, err := syslog.New(prio, executableName())
	if err != nil {
		return nil, err
	}
	return &syslogSink{w, format, fields}, nil
}
