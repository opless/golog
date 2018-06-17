package util

import "fmt"
import "os"

// PanicThrowerFunc ...
type PanicThrowerFunc func(err error)

// ErrorThrower ...
var ErrorThrower PanicThrowerFunc = nil

// MaybePanic ...
func MaybePanic(err error) {
	if err != nil {
		fmt.Println("Panic", err)
		if ErrorThrower == nil {
			panic(err)
		} else {
			ErrorThrower(err)
		}
	}
}

// Debugging ...
func Debugging() bool {
	return os.Getenv("GOLOG_DEBUG") != ""
}

func Debugf(format string, args ...interface{}) {
	if Debugging() {
		fmt.Fprintf(os.Stderr, format, args...)
	}
}
