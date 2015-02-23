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
	"bytes"
	"fmt"
	"testing"
)

func TestSink(t *testing.T) {
	var buf bytes.Buffer
	sink := WriterSink(&buf, SyslogFormat, SyslogFields)
	sink.Log(Fields{
		"priority": func() interface{} { return INFO },
		"message":  func() interface{} { return "hello!" },
	})
	if buf.String() != fmt.Sprintf(SyslogFormat, INFO, "hello!") {
		t.Errorf("unexpected output. %s", buf.String())
	}
}
