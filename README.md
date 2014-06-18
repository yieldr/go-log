# go-log

A simple logging library for Go

### Examples

#### Default

This example uses the default log to log to standard out:

```go
package main

import (
	"github.com/yieldr/go-log/log"
)

func main() {
	log.Info("Hello World.")
	log.Error("There's nothing more to this program.")
}
```

#### Using Sinks and Formats

```go

package main

import (
	"github.com/yieldr/go-log/log"
	"os"
)

func main() {
	l := log.NewSimple(
		log.WriterSink(os.Stderr,
			"%s: %s[%d] %s\n",
			[]string{"priority", "executable", "pid", "message"}))
	l.Info("Here's a differently formatted log message.")
}
```

#### File Sink

This example only logs messages with priority `PriErr` and greater.

```go
package main

import (
	"github.com/yieldr/go-log/log"
	"os"
)

func main() {
	s, err := log.FileSink("/var/log/app.log", log.BasicFormat, log.BasicFields)
	if err != nil {
		panic(err)
	}
	l := log.NewSimple(s)
	l.Info("This will be written to file, in about .")
}
```

### Fields

The following fields are available for use in all sinks:

```go
"prefix"       string              // static field available to all sinks
"seq"          uint64              // auto-incrementing sequence number
"start_time"   time.Time           // start time of the log
"time"         string              // formatted time of log entry
"full_time"    time.Time           // time of log entry
"rtime"        time.Duration       // relative time of log entry since started
"pid"          int                 // process id
"executable"   string              // executable filename
```

In addition, if `verbose=true` is passed to `New()`, the following (somewhat expensive) runtime fields are also available:

```go
"funcname"     string              // function name where the log function was called
"lineno"       int                 // line number where the log function was called
"pathname"     string              // full pathname of caller
"filename"     string              // filename of caller
```

### Logging functions

All these functions can also be called directly to use the default log.

```go
func (*Logger) Log(priority Priority, v ...interface)
func (*Logger) Logf(priority Priority, format string, v ...interface{})
func (*Logger) Emergency(v ...interface)
func (*Logger) Emergencyf(format string, v ...interface{})
func (*Logger) Alert(v ...interface)
func (*Logger) Alertf(format string, v ...interface{})
func (*Logger) Critical(v ...interface)
func (*Logger) Criticalf(format string, v ...interface{})
func (*Logger) Error(v ...interface)
func (*Logger) Errorf(format string, v ...interface{})
func (*Logger) Warning(v ...interface)
func (*Logger) Warningf(format string, v ...interface{})
func (*Logger) Notice(v ...interface)
func (*Logger) Noticef(format string, v ...interface{})
func (*Logger) Info(v ...interface)
func (*Logger) Infof(format string, v ...interface{})
func (*Logger) Debug(v ...interface)
func (*Logger) Debugf(format string, v ...interface{})
```

### Acknowledgements

This package is a fork of [coreos/go-log](https://github.com/coreos/go-log), which in turn is inspired by [ccding/go-logging](https://github.com/ccding/go-logging).