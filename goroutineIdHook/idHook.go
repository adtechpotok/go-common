package goroutineIdHook

import (
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

type IdHook struct {
	daemonName string
}

func (m *IdHook) Fire(entry *logrus.Entry) error {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idGoroutine := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	entry.Data["grId"] = idGoroutine
	entry.Data["name"] = m.daemonName

	return nil
}

func (m *IdHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func New(name string) *IdHook {
	return &IdHook{name}
}
