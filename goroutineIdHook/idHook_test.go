package goroutineIdHook

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFire(t *testing.T) {
	testAssert := assert.New(t)
	testCases := []struct {
		entry   *logrus.Entry
		message string
	}{
		{&logrus.Entry{Level: logrus.PanicLevel, Message: "1"}, "1"},
		{&logrus.Entry{Level: logrus.FatalLevel, Message: "1"}, "1"},
		{&logrus.Entry{Level: logrus.ErrorLevel, Message: "1"}, "1"},
		{&logrus.Entry{Level: logrus.WarnLevel, Message: "1"}, "1"},
		{&logrus.Entry{Level: logrus.InfoLevel, Message: "1"}, "1"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("entry:%s, message:%s", tc.entry.Level, tc.message), func(t *testing.T) {
			hook := &IdHook{}
			tc.entry.Data = make(map[string]interface{})
			hook.Fire(tc.entry)
			testAssert.NotEqual(tc.entry.Data["grId"], nil)
		})
	}
}

func TestLevels(t *testing.T) {
	hook := New("test")
	testAssert := assert.New(t)
	testAssert.Equal(logrus.AllLevels, hook.Levels())
}
