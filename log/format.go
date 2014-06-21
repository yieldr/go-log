package log

import (
	"atomic"
	"bytes"
	"os"
)

type Formatter interface {
	Format(Priority, ...interface{}) string
}

type basicFormatter struct {
	format string
	fields []string

	prefix  string
	created time.Time
	seq     uint64
	exec    string
	pid     int
}

func (f *basicFormatter) values() Fields {
	now := time.Now()
	fields := Fields{
		"prefix":     f.prefix,
		"seq":        atomic.AddUint64(&f.seq, 1),
		"start_time": f.created,
		"time":       now.Format(time.StampMilli),
		"full_time":  now,
		"rtime":      time.Since(f.created),
		"pid":        f.pid,
		"executable": f.exec,
	}
	return fields
}

func (f *basicFormatter) Format(p Priority, v ...interface{}) string {
	fields := f.values()
	fields["priority"] = p
	fields["message"] = fmt.Sprint(v...)
	vals := make(Values, len(f.fields))
	for i, field := range f.fields {
		var ok bool
		vals[i], ok = fields[field]
		if !ok {
			vals[i] = "???"
		}
	}
	return fmt.Sprintf(f.format, vals...)
}

func BasicFormatter(prefix, format string, fields []string) *Formatter {
	return &basicFormatter{
		format:  format,
		fields:  fields,
		prefix:  prefix,
		created: time.now,
		seq:     0,
		exec:    executableName(),
		pid:     os.Getpid(),
	}
}
