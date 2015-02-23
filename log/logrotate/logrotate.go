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
package logrotate

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
	file     *os.File
	buf      *bufio.Writer
	filename string
	format   string
	interval time.Duration
	fields   []string
	err      chan error
	stop     chan bool
	mux      sync.Mutex
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
	if err := l.close(); err != nil {
		return err
	}
	if err := l.open(); err != nil {
		return err
	}
	return nil
}

var dateFormat = "2006-01-02T150405"

// DateFormat sets the date format to be used as the extension to rotated files.
func DateFormat(s string) {
	dateFormat = s
}

func (l *Logrotate) rotate(t time.Time) error {
	if err := l.close(); err != nil {
		return err
	}
	target := l.filename + "." + t.Format(dateFormat)
	if err := os.Rename(l.filename, target); err != nil {
		return err
	}
	if err := l.open(); err != nil {
		return err
	}
	return nil
}

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

func (l *Logrotate) Reload() error {
	l.mux.Lock()
	defer l.mux.Unlock()
	return l.reload()
}

func (l *Logrotate) Rotate() error {
	l.mux.Lock()
	defer l.mux.Unlock()
	return l.rotate(time.Now())
}

func (l *Logrotate) Close() error {
	l.mux.Lock()
	defer l.mux.Unlock()
	return l.close()
}

func (l *Logrotate) Run() {
	ticker := time.NewTicker(l.interval)
	for {
		select {
		case <-ticker.C:
			l.Rotate()
		case <-l.stop:
			ticker.Stop()
			return
		}
	}
}

func (l *Logrotate) Done() <-chan bool {
	return l.stop
}

func (l *Logrotate) Error() <-chan error {
	return l.err
}

func (l *Logrotate) Stop() {
	l.stop <- true
}

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
