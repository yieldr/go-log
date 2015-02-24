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
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/yieldr/go-log/log"
)

var tmpdir = fmt.Sprintf("%s%x/", os.TempDir(), rand.Int())

func TestRotate(t *testing.T) {
	err := os.MkdirAll(tmpdir, 0777)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	t.Logf("tmp dir: %s", tmpdir)

	sink, err := New(tmpdir+"test.log", time.Second*2, log.BasicFormat, log.BasicFields)
	if err != nil {
		t.Fatal(err)
	}
	defer sink.Close()

	log.New(sink).Info("hello!")

	// now rotate
	if err = sink.Rotate(); err != nil {
		t.Fatal(err)
	}

	// lets check if the rotated file exists
	files, err := ioutil.ReadDir(tmpdir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) < 2 {
		t.Fatalf("expected two files to be created")
	}

	file, err := os.Open(tmpdir + files[len(files)-1].Name())
	if err != nil {
		t.Fatal(err)
	}

	// compare the contents
	b, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	expected := fmt.Sprintf(log.BasicFormat, time.Now().Format(time.Stamp), log.INFO, "hello!")

	if !bytes.Equal([]byte(expected), b) {
		t.Fatalf("expected output to be %q, but intead got %q", expected, b)
	}
}

func TestRun(t *testing.T) {
	err := os.MkdirAll(tmpdir, 0777)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	t.Logf("tmp dir: %s", tmpdir)

	sink, err := New(tmpdir+"test.log", time.Second*2, log.BasicFormat, log.BasicFields)
	if err != nil {
		t.Fatal(err)
	}
	defer sink.Close()

	logger := log.New(sink)

	go sink.Run() // start logrotate in a separate goroutine.

	go func() {
		for i := 0; i < 40; i++ {
			logger.Infof("%x", rand.Int())
			time.Sleep(time.Millisecond * 250)
		}
		sink.Stop() // stop the Run() method and end the test.
	}()

	select {
	case err := <-sink.Error():
		t.Error(err)
	case <-sink.Done():
		sink.Close()
	}
}
