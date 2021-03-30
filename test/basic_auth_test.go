package test

import (
	"encoding/base64"
	"log"
	"reflect"
	"testing"
	"unsafe"
)

func TestBasicAuthInGin(t *testing.T) {
	log.Print(authorizationHeader("foo", "bar"))
}

func authorizationHeader(user, password string) string {
	base := user + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString(StringToBytes(base))
}

func StringToBytes(s string) (b []byte) {
	sh := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data, bh.Len, bh.Cap = sh.Data, sh.Len, sh.Len
	return b
}
