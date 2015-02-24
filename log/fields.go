package log

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

import (
	"os"
	"sync/atomic"
	"time"
)

// Fields is used by the logger to supply to its sinks. It contains information
// such as the date, sequence, pid and so on.
type Fields map[string]fieldFn

type fieldFn func() interface{}

func (l *logger) prefixFn() interface{} {
	return l.prefix
}

var dateFormat = time.Stamp

// DateFormat sets the format to use when formatting the time.
func DateFormat(f string) {
	dateFormat = f
}

func (l *logger) createdFn() interface{} {
	return l.created.Format(dateFormat)
}

func (l *logger) timeFn() interface{} {
	return time.Now().Format(dateFormat)
}

func (l *logger) elapsedFn() interface{} {
	return time.Since(l.created)
}

func (l *logger) seqFn() interface{} {
	return atomic.AddUint64(&l.seq, 1)
}

func (l *logger) pidFn() interface{} {
	return os.Getpid()
}
