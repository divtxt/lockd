package testcases

import (
	"strings"

	"fmt"

	"github.com/divtxt/lockd/lockd_client"
	"github.com/divtxt/lockd/util"
)

const sampleNihongo = "日本語"
const sampleInvalidUtf8 = "\xbd\xb2\x3d\xbc\x20\xe2\x8c\x98"

func BadLockNamesTest() {
	fmt.Println("BadLockNamesTest")

	lc := lockd_client.NewLockdClient()

	// all printable ascii - implicitly tests for "/" in the name
	assert(lc.IsLocked, util.NameValidChars, false)
	assert(lc.Lock, util.NameValidChars, true)
	assert(lc.IsLocked, util.NameValidChars, true)
	assert(lc.Unlock, util.NameValidChars, true)
	assert(lc.IsLocked, util.NameValidChars, false)

	// longest name
	longestName := strings.Repeat("a", util.NameLenMaxBytes)
	assert(lc.IsLocked, longestName, false)

	// // bad names should panic
	// testhelpers.TestHelper_ExpectPanic(
	// 	t,
	// 	func() { lsm.IsLocked("") },
	// 	"Name is empty string",
	// )
	// testhelpers.TestHelper_ExpectPanic(
	// 	t,
	// 	func() { lsm.Lock("") },
	// 	"Name is empty string",
	// )
	// testhelpers.TestHelper_ExpectPanic(
	// 	t,
	// 	func() { lsm.IsLocked(sampleNihongo) },
	// 	"Name contains non-printable/non-ascii byte 230 at offset 0",
	// )
	// testhelpers.TestHelper_ExpectPanic(
	// 	t,
	// 	func() { lsm.Unlock(sampleInvalidUtf8) },
	// 	"Name contains non-printable/non-ascii byte 189 at offset 0",
	// )

	// assert(lc.Lock, "", true)
	// assert(lc.IsLocked, "foo", true)

	// // Dup lock should fail
	// assert(lc.Lock, "foo", false)
	// assert(lc.IsLocked, "foo", true)

	// // Lock another entry should work
	// assert(lc.Lock, "bar", true)
	// assert(lc.IsLocked, "bar", true)

	// // Unlock entries
	// assert(lc.Unlock, "bar", true)
	// assert(lc.Unlock, "foo", true)
	// assert(lc.IsLocked, "foo", false)
	// assert(lc.IsLocked, "bar", false)

	// // Dup unlock should fail
	// assert(lc.Unlock, "bar", false)
	// assert(lc.IsLocked, "bar", false)

}
