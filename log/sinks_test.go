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
	"fmt"
	"os"
	"sync"
	"testing"
)

const length = 100

func TestFileSink(t *testing.T) {
	var err error
	sinks := make([]*fileSink, length)
	for i := 0; i < length; i++ {
		file := fmt.Sprintf("/tmp/file_sink_test_%d.log", i+1)
		t.Log(file)
		sinks[i], err = FileSink(file, BasicFormat, BasicFields)
		if err != nil {
			t.Error(err.Error())
			t.Fail()
		}
	}

	defer func() {
		for _, sink := range sinks {
			t.Log(sink.file.Name())
			os.Remove(sink.file.Name())
			sink.Close()
		}
	}()

	var write sync.WaitGroup
	var reload sync.WaitGroup

	// Concurrently write, flush and reload.
	for _, sink := range sinks {
		s := *sink
		t.Log(s.file.Name())
		write.Add(1)
		go func() {
			defer write.Done()
			n, err := s.Write([]byte("message\n"))
			if err != nil {
				t.Error(err.Error())
			}
			if n == 0 {
				t.Error("no bytes written")
			}
			err = s.Flush()
			if err != nil {
				t.Error(err)
				t.Log(s.file, s.out.Buffered())
			}
		}()
		reload.Add(1)
		go func() {
			defer reload.Done()
			err = s.Reload()
			if err != nil {
				t.Error(err.Error())
			}
		}()
	}

	write.Wait()
	reload.Wait()
}
