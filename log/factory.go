package log

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
	"sync"
)

var (
	// qualified package name, cached at first use
	logPackage string

	// Positions in the call stack when tracing to report the calling method
	minimumCallerDepth = 1

	// Used for caller information initialisation
	callerInitOnce sync.Once
)

const (
	maximumCallerDepth    int = 25
	knownLogFactoryFrames int = 5
)

type LogFactory struct {
	Loggers       map[string]*logrus.Logger
	defaultFilter Filter
}

type Filter interface {
	Customize(logger *logrus.Logger, name string) error
}

func (factory *LogFactory) GetLogger(name string, filter Filter) (*logrus.Logger, error) {
	if value, exists := factory.Loggers[name]; exists {
		return value, nil
	}
	logger := logrus.New()

	if filter == nil {
		filter = factory.defaultFilter
	}
	if filter == nil {
		return nil, errors.New("No filter specified")
	}

	if err := filter.Customize(logger, name); err != nil {
		return nil, err
	}

	// need to save the logger to map after creating a new logger.
	factory.Loggers[name] = logger

	return logger, nil
}

func getLogFactoryCaller() *runtime.Frame {
	// cache this package's fully-qualified name
	callerInitOnce.Do(func() {
		pcs := make([]uintptr, maximumCallerDepth)
		_ = runtime.Callers(1, pcs)

		// dynamic get the package name and the minimum caller depth
		for i := 0; i < maximumCallerDepth; i++ {
			funcName := runtime.FuncForPC(pcs[i]).Name()
			if strings.Contains(funcName, "getLogFactoryCaller") {
				logPackage = getPackageName(funcName)
				break
			}
		}

		minimumCallerDepth = knownLogFactoryFrames
	})

	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)

		// If the caller isn't part of this package, we're done
		if pkg != logPackage {
			return &f //nolint:scopelint
		}
	}

	// if we got here, we failed to find the caller's context
	return nil
}

// getPackageName reduces a fully qualified function name to the package name
// There really ought to be to be a better way...
func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}

	return f
}
