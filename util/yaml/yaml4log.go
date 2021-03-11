package yaml

import (
	"github.com/sirupsen/logrus"
	"sync"
)

func (engine *ConfigEngine) Customize(logger *logrus.Logger, name string) error {
	return nil
}

var instance *ConfigEngine
var lock sync.Mutex

/*
	New is to create singleton pointer for ConfigEngine.
*/
func New() *ConfigEngine {
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()

		if instance == nil {
			tmp := ConfigEngine{}
			instance = &tmp
		}
	}
	return instance
}
