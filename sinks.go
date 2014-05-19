package logging

// Copyright 2013, CoreOS, Inc. All rights reserved.
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
// author: David Fisher <ddf1991@gmail.com>
// based on previous package by: Cong Ding <dinggnu@gmail.com>

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

const (
	flushInterval  = time.Second * 30
	reloadInterval = time.Hour * 1
)

type Sink interface {
	Log(Fields)
}

type ReloadSink interface {
	Sink
	Reload() error
	Flush() error
	Close() error
}

type nullSink struct{}

func (sink *nullSink) Log(fields Fields) {}

func NullSink() Sink {
	return &nullSink{}
}

type writerSink struct {
	lock   sync.Mutex
	out    io.Writer
	format string
	fields []string
}

func (sink *writerSink) Log(fields Fields) {
	vals := make([]interface{}, len(sink.fields))
	for i, field := range sink.fields {
		var ok bool
		vals[i], ok = fields[field]
		if !ok {
			vals[i] = "???"
		}
	}

	sink.lock.Lock()
	defer sink.lock.Unlock()
	fmt.Fprintf(sink.out, sink.format, vals...)
}

func WriterSink(out io.Writer, format string, fields []string) Sink {
	return &writerSink{
		out:    out,
		format: format,
		fields: fields,
	}
}

type priorityFilter struct {
	priority Priority
	target   Sink
}

func (filter *priorityFilter) Log(fields Fields) {
	// lower priority values indicate more important messages
	if fields["priority"].(Priority) <= filter.priority {
		filter.target.Log(fields)
	}
}

func PriorityFilter(priority Priority, target Sink) Sink {
	return &priorityFilter{
		priority: priority,
		target:   target,
	}
}

type fileSink struct {
	out    *bufio.Writer
	file   *os.File
	format string
	fields []string
	mux    sync.Mutex
}

func (sink *fileSink) open(name string) (err error) {
	sink.file, err = os.OpenFile(name, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		sink.file, err = os.Create(name)
		if err != nil {
			return errors.New("logging: unable to open or create file")
		}
	}
	return
}

func (sink *fileSink) close() error {
	err := sink.flush()
	if err != nil {
		return err
	}
	return sink.file.Close()
}

func (sink *fileSink) flush() error {
	return sink.out.Flush()
}

func (sink *fileSink) reload() (err error) {
	name := sink.file.Name()
	sink.close()
	err = sink.open(name)
	if err == nil {
		sink.out = bufio.NewWriter(sink.file)
	}
	return
}

func (sink *fileSink) daemon() {
	flush := time.NewTicker(flushInterval)
	reload := time.NewTicker(reloadInterval)
	for {
		select {
		case <-flush.C:
			sink.flush()
		case <-reload.C:
			sink.reload()
		}
	}
}

func (sink *fileSink) Log(fields Fields) {
	vals := make([]interface{}, len(sink.fields))
	for i, field := range sink.fields {
		var ok bool
		vals[i], ok = fields[field]
		if !ok {
			vals[i] = "???"
		}
	}

	sink.mux.Lock()
	defer sink.mux.Unlock()
	fmt.Fprintf(sink.out, sink.format, vals...)
}

// Closes and reopens the output file, in order to momentarily release it's file
// handle. Typically this functionality is combined with a SIGHUP system signal.
// Before reloading, the content of the buffer is flushed.
func (sink *fileSink) Reload() error {
	sink.mux.Lock()
	defer sink.mux.Unlock()
	return sink.reload()
}

// Flushes the buffer to disk.
func (sink *fileSink) Flush() error {
	sink.mux.Lock()
	defer sink.mux.Unlock()
	return sink.flush()
}

// Closes any open file handles used by the Sink.
func (sink *fileSink) Close() error {
	sink.mux.Lock()
	defer sink.mux.Unlock()
	return sink.close()
}

// Returns a new Sink able to buffer output and periodically flush to disk.
func FileSink(name string, format string, fields []string) ReloadSink {
	sink := &fileSink{
		format: format,
		fields: fields,
	}
	err := sink.open(name)
	if err != nil {
		panic(err)
	}
	sink.out = bufio.NewWriter(sink.file)
	go sink.daemon()
	return sink
}
