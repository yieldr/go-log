package logrotate

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
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/yieldr/go-log/log"
)

// Logrotate is a special case of sink which writes to a file and is capable of
// rotating that file when certain conditions are met.
type Logrotate struct {
	file        *os.File
	buf         *bufio.Writer
	filename    string
	format      string
	interval    time.Duration
	fields      []string
	err         chan error
	stop        chan bool
	mux         sync.Mutex
	subscribers []Subcriber
}

func (l *Logrotate) open() error {
	file, err := os.OpenFile(l.filename, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		file, err = os.Create(l.filename)
		if err != nil {
			return fmt.Errorf("log: unable to open or create file %s", l.filename)
		}
	}
	l.file = file
	l.buf = bufio.NewWriter(l.file)
	return nil
}

func (l *Logrotate) flush() error {
	return l.buf.Flush()
}

func (l *Logrotate) close() error {
	if err := l.flush(); err != nil {
		return err
	}
	return l.file.Close()
}

func (l *Logrotate) reload() error {
	for _, subcriber := range l.subscribers {
		if err := subcriber.OnPreReload(); err != nil {
			return err
		}
	}

	if err := l.close(); err != nil {
		return err
	}
	if err := l.open(); err != nil {
		return err
	}

	for _, subcriber := range l.subscribers {
		if err := subcriber.OnPostReload(); err != nil {
			return err
		}
	}
	return nil
}

var dateFormat = "2006-01-02T150405"

// DateFormat sets the date format to be used as the extension to rotated files.
func DateFormat(s string) {
	dateFormat = s
}

func (l *Logrotate) rotate(t time.Time) error {
	for _, subcriber := range l.subscribers {
		if err := subcriber.OnPreRotate(t); err != nil {
			return err
		}
	}

	if err := l.close(); err != nil {
		return err
	}

	target := l.filename + "." + t.Format(dateFormat)
	if err := os.Rename(l.filename, target); err != nil {
		return err
	}

	for _, subcriber := range l.subscribers {
		if err := subcriber.OnPostRotate(target); err != nil {
			return err
		}
	}

	if err := l.open(); err != nil {
		return err
	}
	return nil
}

// Log satisfies the log.Sink interface so it can be supplied as an argument to
// log.New(). It writes the log to the internal buffer, using the format and
// fields.
func (l *Logrotate) Log(fields log.Fields) {
	l.mux.Lock()
	defer l.mux.Unlock()
	vals := make([]interface{}, len(l.fields))
	for i, field := range l.fields {
		if fn, ok := fields[field]; ok {
			vals[i] = fn()
		} else {
			vals[i] = "???"
		}
	}
	fmt.Fprintf(l.buf, l.format, vals...)
}

// Write writes p to the internal buffer.
func (l *Logrotate) Write(p []byte) (int, error) {
	l.mux.Lock()
	defer l.mux.Unlock()
	return l.buf.Write(p)
}

// Flush empties the contents of the internal buffer to the output file.
func (l *Logrotate) Flush() error {
	l.mux.Lock()
	defer l.mux.Unlock()
	return l.flush()
}

// Reload flushes the internal buffer to the output file, closes the file and re
// opens the file.
func (l *Logrotate) Reload() error {
	l.mux.Lock()
	defer l.mux.Unlock()
	return l.reload()
}

// Rotate performs a file rotation similar to that of the logrotate(8) utility.
func (l *Logrotate) Rotate() error {
	l.mux.Lock()
	defer l.mux.Unlock()
	return l.rotate(time.Now())
}

// Close flushes the internal buffer to the output file and closes the file.
func (l *Logrotate) Close() error {
	l.mux.Lock()
	defer l.mux.Unlock()
	return l.close()
}

// Run will block and call Rotate or Flush on their respective intervals. If an
// error occurs during those operations should be handled by callers using the
// error channel returned by the Error method. This method will only return once
// the Stop method is called.
func (l *Logrotate) Run() {
	rotate := time.NewTicker(l.interval)
	flush := time.NewTicker(time.Second * 3)
	for {
		select {
		case <-rotate.C:
			if err := l.Rotate(); err != nil {
				l.err <- err
			}
		case <-flush.C:
			if err := l.Flush(); err != nil {
				l.err <- err
			}
		case <-l.stop:
			rotate.Stop()
			flush.Stop()
			return
		}
	}
}

// Done returns a channel which will receive when the Run method has returned.
func (l *Logrotate) Done() <-chan bool {
	return l.stop
}

// Error returns a channel which will receive an error if one was encountered
// during a rotate or flush operation. It's up to the caller to handle the error
// so Logrotate will keep running until Stop is called.
func (l *Logrotate) Error() <-chan error {
	return l.err
}

// Stop ends the execution of Run.
func (l *Logrotate) Stop() {
	l.stop <- true
}

// AddSubscriber appends a subscriber to handle events.
// The user can add multiple subscribers, each subcriber
// will be called in the order as it is appended.
func (l *Logrotate) AddSubscriber(sub Subcriber) {
	l.subscribers = append(l.subscribers, sub)
}

// New returns a new Logrotate using the supplied arguments.
func New(file string, interval time.Duration, format string, fields []string) (*Logrotate, error) {
	l := &Logrotate{
		filename: file,
		format:   format,
		interval: interval,
		fields:   fields,
		err:      make(chan error),
		stop:     make(chan bool),
	}
	return l, l.open()
}

// Subcriber defines event handler functions inside logrotate:
// OnPreRotate is called BEFORE the file is rotated
// OnPostRotate is called AFTER the file is rotated
// OnPreReload is called BEFORE logrotate is reloaded
// OnPostReload is called AFTER logrotate is reloaded
type Subcriber interface {
	OnPreRotate(t time.Time) error
	OnPostRotate(rotateFilename string) error
	OnPreReload() error
	OnPostReload() error
}
