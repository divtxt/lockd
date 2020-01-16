package testcases

import (
	"fmt"
	"strings"

	"github.com/divtxt/lockd/lockd_client"
	"github.com/divtxt/lockd/util"
)

const sampleNihongo = "日本語"
const sampleInvalidUtf8 = "\xbd\xb2\x3d\xbc\x20\xe2\x8c\x98"

func LockNamesTest() {
	fmt.Println("LockNamesTest")

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

	// empty string
	assertError(
		lc.Lock,
		"",
		"Unexpected response: 404 Not Found for POST /lock/",
	)

	// invalid chars
	assertError(
		lc.Lock,
		"A+B",
		"Bad Request: {\"error\":\"Name contains non-printable/non-ascii byte 43 at offset 1\"}\n",
	)

	// ascii control character
	assertError(
		lc.Lock,
		"hi\n",
		"Bad Request: {\"error\":\"Name contains non-printable/non-ascii byte 10 at offset 2\"}\n",
	)

	// non-ascii unicode
	assertError(
		lc.Lock,
		sampleNihongo,
		"Bad Request: {\"error\":\"Name contains non-printable/non-ascii byte 230 at offset 0\"}\n",
	)

	// non-ascii invalid utf8
	assertError(
		lc.Lock,
		sampleInvalidUtf8,
		"Bad Request: {\"error\":\"Name contains non-printable/non-ascii byte 189 at offset 0\"}\n",
	)

	// exceeding max length
	assertError(
		lc.Lock,
		longestName+"a",
		"Bad Request: {\"error\":\"Name is too long (129 bytes \\u003e max of 128)\"}\n",
	)
}
