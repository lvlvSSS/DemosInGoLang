package yaml

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

type ConfigEngine struct {
	dataMap map[interface{}]interface{}
}

func (engine *ConfigEngine) Load(file string) error {

	extension := filepath.Ext(file)
	switch extension {
	case ".yml", ".yaml":
		return engine.loadFromYaml(file)
	default:
		return errors.New(fmt.Sprintf("can not load the file[%s] ", file))
	}
}

func (engine *ConfigEngine) lazyInit() error {
	targetDir, _ := os.Getwd()
	return engine.Load(filepath.Join(targetDir, "resources", "config", "log4go.yaml"))
}

func (engine *ConfigEngine) loadFromYaml(file string) error {
	str, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(str, &engine.dataMap)
	if err != nil {
		return errors.New(fmt.Sprintf("can not parse the yaml file[%s]", file))
	}
	return nil
}

func (engine *ConfigEngine) Get(name string) (interface{}, error) {
	/*initialize the engine if necessary*/
	if engine.dataMap == nil || len(engine.dataMap) == 0 {
		if initErr := engine.lazyInit(); initErr != nil {
			return nil, initErr
		}
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("please specify the name")
	}

	keys := strings.Split(name, ".")
	datas := engine.dataMap
	var err error
	for index, key := range keys {
		data, ok := datas[key]
		if !ok {
			err = errors.New(fmt.Sprintf("the key[%s] not exists in data", key))
			break
		}

		if (index + 1) == len(keys) {
			return data, nil
		}
		fmt.Printf(" ---- %v -------\n", reflect.TypeOf(data).String())
		if reflect.TypeOf(data).String() == "map[interface {}]interface {}" {
			datas = data.(map[interface{}]interface{})
		}
	}

	return nil, err
}

func (engine *ConfigEngine) GetString(name string) (string, error) {
	value, err := engine.Get(name)
	if err != nil {
		return "", err
	}
	switch value.(type) {
	case string:
		return value.(string), nil
	case float64, int, bool:
		return fmt.Sprintf("%s", value), nil
	}
	return "", errors.New(fmt.Sprintf("the value of key[%s] in map can't convert to string", name))
}

func (engine *ConfigEngine) GetInt(name string) (int, error) {
	value, err := engine.Get(name)
	if err != nil {
		return 0, err
	}
	switch value.(type) {
	case string:
		return strconv.Atoi(value.(string))
	case float64:
		return int(value.(float64)), nil
	case int:
		return value.(int), nil
	case bool:
		b := value.(bool)
		if b {
			return 1, nil
		}
		return 0, nil
	}
	return 0, errors.New(fmt.Sprintf("the value of key[%s] in map can't convert to int", name))
}

func (engine *ConfigEngine) GetBool(name string) (bool, error) {
	value, err := engine.Get(name)
	if err != nil {
		return false, err
	}
	switch value.(type) {
	case bool:
		return value.(bool), nil
	case string:
		return strconv.ParseBool(value.(string))
	case float64:
		if value.(float64) == 0.0 {
			return false, nil
		}
		return true, nil
	case int:
		if value.(int) == 0 {
			return false, nil
		}
		return true, nil
	}
	return false, errors.New(fmt.Sprintf("the value of key[%s] in map can't convert to bool", name))

}

func (engine *ConfigEngine) GetFloat64(name string) (float64, error) {
	value, err := engine.Get(name)
	if err != nil {
		return 0, err
	}

	switch value.(type) {
	case float64:
		return value.(float64), nil
	case string:
		return strconv.ParseFloat(value.(string), 64)
	case int:
		return float64(value.(int)), nil
	case bool:
		if value.(bool) {
			return float64(1), nil
		}
		return float64(0), nil
	}
	return 0, errors.New(fmt.Sprintf("the value of key[%s] in map can't convert to float64", name))
}

func (engine *ConfigEngine) GetStruct(name string, s interface{}) (interface{}, error) {
	if engine == nil || len(engine.dataMap) == 0 {
		if initErr := engine.lazyInit(); initErr != nil {
			return nil, initErr
		}
	}

	if strings.TrimSpace(name) == "" {
		return engine.mapToStruct(engine.dataMap, s)
	}

	d, err := engine.Get(name)
	if err != nil {
		return nil, err
	}

	if reflect.TypeOf(d).String() != "map[interface {}]interface {}" {
		return nil, errors.New(fmt.Sprintf(
			"the value of the key[%s] in map is not map, so can't not convert to specified struct", name))
	}

	return engine.mapToStruct(d.(map[interface{}]interface{}), s)
}

func (engine *ConfigEngine) mapToStruct(m map[interface{}]interface{}, pointer interface{}) (interface{}, error) {
	for key, value := range m {
		if _, ok := key.(string); !ok {
			return nil, errors.New(fmt.Sprintf("the key[%v] need be converted to string", key))
		}
		engine.setField(pointer, key.(string), value)
	}
	return pointer, nil
}

/*
	The setField is to set the specified field of the struct.
	WARN: the field of struct could not be pointer.
*/
func (engine *ConfigEngine) setField(pointer interface{}, key string, value interface{}) error {
	/* the pointer is the pointer that point to the struct */
	realValue := reflect.Indirect(reflect.ValueOf(pointer))
	if realValue.Type().Kind() != reflect.Struct {
		return errors.New("can't set the non-struct field")
	}

	fieldValue := realValue.FieldByName(key)
	if !fieldValue.IsValid() {
		return errors.New(fmt.Sprintf("No such field[%s]", key))
	}
	if !fieldValue.CanSet() {
		return errors.New(fmt.Sprintf("field[%s] can not be set", key))
	}

	fieldValueType := fieldValue.Type()
	val := reflect.ValueOf(value)
	if fieldValueType.Kind() == reflect.Struct && val.Kind() == reflect.Map {
		switch rval := val.Interface(); rval.(type) {
		case map[interface{}]interface{}:
			for childKey, childValue := range rval.(map[interface{}]interface{}) {
				engine.setField(fieldValue.Addr().Interface(), childKey.(string), childValue)
			}
		case map[string]interface{}:
			for childKey, childValue := range rval.(map[string]interface{}) {
				engine.setField(fieldValue.Addr().Interface(), childKey, childValue)
			}
		}
	} else {
		if fieldValueType != val.Type() {
			return errors.New("Provided value type didn't match obj field type")
		}
		fieldValue.Set(val)
	}

	return nil
}
