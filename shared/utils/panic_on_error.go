package utils

import "fmt"

// Panics if error is not nil with provided message and error info.
func PanicOnError(err error, msg string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %v", msg, err.Error()))
	}
}
