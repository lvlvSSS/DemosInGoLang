package log

import (
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/zput/zxcTool/ztLog/zt_formatter"
	"io"
	"path"
	"runtime"
	"time"
)

func New(filter Filter) *LogFactory {
	if filter == nil {
		filter = &InnerFilter{
			Formatter: &zt_formatter.ZtFormatter{
				CallerPrettyfier: func(f *runtime.Frame) (string, string) {
					filename := path.Base(f.File)
					return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
				},
				Formatter: nested.Formatter{
					HideKeys:        true,
					FieldsOrder:     []string{"component", "category"},
					TimestampFormat: "2006-01-02 15:04:05.0000",
					ShowFullLevel:   true,
					NoColors:        true,
					NoFieldsColors:  true,
				},
			},
			Level:  logrus.InfoLevel,
			Writer: nil,
		}
	}
	return &LogFactory{
		Loggers:       make(map[string]*logrus.Logger),
		defaultFilter: filter,
	}
}

func init() {
	logrus.SetFormatter(&zt_formatter.ZtFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
		Formatter: nested.Formatter{
			HideKeys:        true,
			FieldsOrder:     []string{"component", "category"},
			TimestampFormat: "2006-01-02 15:04:05.0000",
			ShowFullLevel:   true,
			NoColors:        true,
			NoFieldsColors:  true,
		},
	})
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.InfoLevel)

	targetWriter, err := rotatelogs.New(
		"./logs/%Y%m%d/default.log",
		rotatelogs.WithMaxAge(-1),
		rotatelogs.WithRotationCount(7),
		rotatelogs.WithRotationTime(24*time.Hour))
	if err == nil {
		logrus.AddHook(lfshook.NewHook(lfshook.WriterMap{
			logrus.DebugLevel: targetWriter,
			logrus.InfoLevel:  targetWriter,
			logrus.WarnLevel:  targetWriter,
			logrus.ErrorLevel: targetWriter,
			logrus.FatalLevel: targetWriter,
			logrus.PanicLevel: targetWriter,
		}, logrus.StandardLogger().Formatter))
	}

}

type InnerFilter struct {
	Formatter logrus.Formatter
	Writer    io.Writer
	Level     logrus.Level
}

func (filter *InnerFilter) Customize(logger *logrus.Logger, name string) error {
	defer func() {
		if err := recover(); err != nil {
			logrus.SetReportCaller(true)
			logrus.Fatalf("Customize logger error: %[1]s", err)
		}
	}()

	logger.SetFormatter(filter.Formatter)
	logger.SetReportCaller(true)
	logger.SetLevel(filter.Level)
	if filter.Writer == nil {
		if targetWriter, err := rotatelogs.New(
			"./logs/%Y%m%d/default.log",
			rotatelogs.WithMaxAge(-1),
			rotatelogs.WithRotationCount(7),
			rotatelogs.WithRotationTime(24*time.Hour)); err == nil {
			filter.Writer = targetWriter
			goto directFormat
		} else {
			return err
		}

	}
directFormat:
	logger.AddHook(lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: filter.Writer,
		logrus.InfoLevel:  filter.Writer,
		logrus.WarnLevel:  filter.Writer,
		logrus.ErrorLevel: filter.Writer,
		logrus.FatalLevel: filter.Writer,
		logrus.PanicLevel: filter.Writer,
	}, filter.Formatter))

	return nil
}
