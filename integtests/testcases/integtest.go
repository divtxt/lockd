package testcases

import (
	"fmt"

	"github.com/divtxt/lockd/lockd_client"
)

func SimpleIntegrationTest() {
	fmt.Println("SimpleIntegrationTest")

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
