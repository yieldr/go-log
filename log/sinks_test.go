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
	"io/ioutil"
	"os"
	"sync"
	"testing"
)

const N = 10

var message = []byte{'m', 'e', 's', 's', 'a', 'g', 'e'}

func TestFileSink(t *testing.T) {
	var err error
	sinks := make([]*fileSink, N)
	for i := 0; i < N; i++ {
		file := fmt.Sprintf("/tmp/file_sink_test_%d.log", i)
		sinks[i], err = FileSink(file, BasicFormat, BasicFields)
		if err != nil {
			t.Error(err.Error())
			t.Fail()
		}
	}

	defer func() {
		for _, sink := range sinks {
			b, err := ioutil.ReadFile(sink.file.Name())
			if err != nil {
				t.Error(err.Error())
			}
			if len(b) != len(message)*N {
				t.Error("file content is not as much expected.")
			}
			os.Remove(sink.file.Name())
			sink.Close()
		}
	}()

	var write sync.WaitGroup
	var reload sync.WaitGroup

	// Concurrently write, flush and reload.
	for _, sink := range sinks {
		s := *sink
		write.Add(1)
		go func() {
			defer write.Done()
			for i := 0; i < N; i++ {
				n, err := s.Write(message)
				if err != nil {
					t.Error(err.Error())
				}
				if n == 0 {
					t.Error("no bytes written")
				}
				t.Log("WR", s.file.Name(), s.out.Buffered())
			}
			err = s.Flush()
			if err != nil {
				t.Error(err)
			}
			t.Log("FL", s.file.Name())
		}()
		reload.Add(1)
		go func() {
			defer reload.Done()
			err = s.Reload()
			if err != nil {
				t.Error(err.Error())
			}
			t.Log("RE", s.file.Name(), s.out.Buffered())
		}()
	}

	write.Wait()
	reload.Wait()
}
