package log

import "github.com/sirupsen/logrus"

var debugLogger *logrus.Logger

func init() {
	debugLogger, _ = GetLogger("debug", nil)
}

func Debug(format string, args ...interface{}) {
	debugLogger.Debugf(format, args...)
}
