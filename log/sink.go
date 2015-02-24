package log

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

import (
	"fmt"
	"io"
	"sync"
)

// Sink is the interface that wraps the basic Log method. The sink receives the
// Fields from the Logger.
type Sink interface {
	Log(Fields)
}

type writerSink struct {
	writer io.Writer
	format string
	fields []string
	mux    sync.Mutex
}

func (sink *writerSink) Log(fields Fields) {
	sink.mux.Lock()
	defer sink.mux.Unlock()
	vals := make([]interface{}, len(sink.fields))
	for i, field := range sink.fields {
		if fn, ok := fields[field]; ok {
			vals[i] = fn()
		} else {
			vals[i] = "???"
		}
	}
	fmt.Fprintf(sink.writer, sink.format, vals...)
}

// WriterSink creates a new sink that writes log messages to w.
func WriterSink(w io.Writer, format string, fields []string) Sink {
	return &writerSink{
		writer: w,
		format: format,
		fields: fields,
	}
}

type filter struct {
	priority Priority
	target   Sink
}

func (f *filter) Log(fields Fields) {
	if fields["priority"]().(Priority) <= f.priority {
		f.target.Log(fields)
	}
}

// Filter wraps the sink with leveled logging. A sink wrapped with this method
// will ony write if the priority is equal to or less than p.
func Filter(p Priority, s Sink) Sink {
	return &filter{
		priority: p,
		target:   s,
	}
}
