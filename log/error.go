package log

import "github.com/sirupsen/logrus"

var errorLogger *logrus.Logger

func init() {
	errorLogger, _ = GetLogger("error", nil)
}

func Error(format string, args ...interface{}) {
	errorLogger.Errorf(format, args...)
}
