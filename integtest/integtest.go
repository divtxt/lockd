package main

import (
	"log"

	"github.com/divtxt/lockd/lockd_client"
)

func main() {
	lc := lockd_client.NewLockdClient()

	// Initial state
	assert(lc.IsLocked, "foo", false)
	assert(lc.IsLocked, "bar", false)

	// Lock
	assert(lc.Lock, "foo", true)
	assert(lc.IsLocked, "foo", true)

	// Dup lock should fail
	assert(lc.Lock, "foo", false)
	assert(lc.IsLocked, "foo", true)

	// Lock another entry should work
	assert(lc.Lock, "bar", true)
	assert(lc.IsLocked, "bar", true)

	// Unlock entries
	assert(lc.Unlock, "bar", true)
	assert(lc.Unlock, "foo", true)
	assert(lc.IsLocked, "foo", false)
	assert(lc.IsLocked, "bar", false)

	// Dup unlock should fail
	assert(lc.Unlock, "bar", false)
	assert(lc.IsLocked, "bar", false)

}

func assert(f func(string) (bool, error), name string, expected bool) {
	actual, err := f(name)
	if err != nil {
		log.Panic(err)
	}
	if actual != expected {
		log.Panicf("Actual: %v != Expected: %v", actual, expected)
	}
}
