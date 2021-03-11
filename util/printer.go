package util

import (
	"runtime"
)

func GetCallerFuncDetails() (funcName string, fileName string, line int) {
	pc, file, l, _ := runtime.Caller(1)

	funcName = runtime.FuncForPC(pc).Name()
	fileName = file
	line = l
	return
}

func GetCallFuncName(depth int) string {
	pc, _, _, _ := runtime.Caller(depth)
	return runtime.FuncForPC(pc).Name()
}
