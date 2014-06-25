package log

import "testing"

func TestSyslogSink(t *testing.T) {
	syslog, err := SyslogSink(PriDebug, BasicFormat, BasicFields)
	if err != nil {
		t.Error(err.Error())
	}
	logger := NewSimple(syslog)
	logger.Info("syslog is cuhl")
}
