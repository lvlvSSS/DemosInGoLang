package test

import (
	"reflect"
	"testing"
)

type myMap map[interface{}]interface{}

func TestTypeAnotherName(t *testing.T) {
	map1 := make(myMap)

	map2 := make(map[interface{}]interface{})
	map2["b"] = 2
	map1["a"] = map2
	t.Log(reflect.TypeOf(map2).String())
	t.Log(reflect.TypeOf(myMap(map2)))

	t.Logf("%v", (map[interface{}]interface{})(map1))
}
