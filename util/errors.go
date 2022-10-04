package util

import (
	"fmt"
)

var panicFunc = func(v interface{}) {
	panic(v)
}

func PanicFunc(f func(v interface{})) {
	panicFunc = f
}

func Must(args ...interface{}) {
	for _, arg := range args {
		if err, ok := arg.(error); ok && err != nil {
			panicFunc(err)
			return
		}
	}
}

func selectVal[T any](rule func(T) bool, args ...interface{}) T {
	var val T
	found := false
	for _, arg := range args {
		if arg == nil {
			continue
		}
		if _, ok := arg.(T); ok && rule(arg.(T)) {
			val = arg.(T)
			found = true
			continue
		}
		Must(arg)
	}
	if !found {
		panic(fmt.Errorf("no value found"))
	}
	return val
}

func SelectNotNil[T any](args ...interface{}) T {
	return selectVal[T](func(arg T) bool {
		return true
	}, args...)
}

func SelectError(args ...interface{}) error {
	return selectVal[error](func(arg error) bool {
		if _, ok := arg.(error); ok {
			return true
		}
		return false
	}, args...)
}

func SelectString(args ...interface{}) string {
	return selectVal[string](func(str string) bool {
		return len(str) > 0
	}, args...)
}

func SelectAnyString(args ...interface{}) string {
	return selectVal[string](func(str string) bool {
		return true
	}, args...)
}

func SelectBool(args ...interface{}) bool {
	return selectVal[bool](func(arg bool) bool {
		return arg
	}, args...)
}

func SelectByteSlice(args ...interface{}) []byte {
	return selectVal[[]byte](func(arg []byte) bool {
		return len(arg) > 0
	}, args...)
}

func SelectAnyByteSlice(args ...interface{}) []byte {
	return selectVal[[]byte](func(arg []byte) bool {
		return true
	}, args...)
}
