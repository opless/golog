package util

import "fmt"
import "os"

type PanicThrowerFunc func(err string)

var ErrorThrower PanicThrowerFunc = nil

func MaybePanic(err error) {
	if err != nil {
		if ErrorThrower == nil {
			panic(err)
		} else {
			ErrorThrower(err)
		}
	}
}

func Debugging() bool {
	return os.Getenv("GOLOG_DEBUG") != ""
}

func Debugf(format string, args ...interface{}) {
	if Debugging() {
		fmt.Fprintf(os.Stderr, format, args...)
	}
}
