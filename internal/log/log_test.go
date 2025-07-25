package log

import (
	"testing"
)

func TestLog(t *testing.T) {
	defer func() {
		if err := Sync(); err != nil {
			Error(err)
		}
	}()
	Debug("test debug")
	Debugf("%s", "test debugf")

	Info("test info")
	Infof("%s", "test infof")

	Warn("test warn")
	Warnf("%s", "test warnf")

	Error("test error")
	Errorf("%s", "test errorf")
}

func TestPanicfLog(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			LogPanicf("%v", r)
		}
	}()
	panic("panictest")
}

func TestPanicLog(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			LogPanic(r)
		}
	}()
	panic("panictest")
}
