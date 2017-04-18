package testcases

import (
	"log"
)

func assert(f func(string) (bool, error), name string, expected bool) {
	actual, err := f(name)
	if err != nil {
		log.Panic(err)
	}
	if actual != expected {
		log.Panicf("Actual: %v != Expected: %v", actual, expected)
	}
}

func assertError(f func(string) (bool, error), name string, expectedErr string) {
	_, err := f(name)
	if err.Error() != expectedErr {
		log.Panicf("Actual error: '%s' != Expected error: '%s'", err, expectedErr)
	}
}
