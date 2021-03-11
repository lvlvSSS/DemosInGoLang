package log

import (
	"github.com/sirupsen/logrus"
	"go_tools/util/yaml"
)

type Filter interface {
	Customize(logger *logrus.Logger, name string) error
}

func GetLogger(name string, filter Filter) (*logrus.Logger, error) {
	logger := logrus.New()
	if filter == nil {
		filter = yaml.New()
	}
	if err := filter.Customize(logger, name); err != nil {
		return nil, err
	}
	return logger, nil
}
