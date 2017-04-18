package util_test

import (
	"strings"
	"testing"

	. "github.com/divtxt/lockd/util"
)

const sampleNihongo = "日本語"
const sampleInvalidUtf8 = "\xbd\xb2\x3d\xbc\x20\xe2\x8c\x98"

func TestIsValidLockName(t *testing.T) {

	// -- valid names:

	// random ascii string
	if e := IsValidLockName("Hello123"); e != "" {
		t.Error(e)
	}

	// all valid chars
	if e := IsValidLockName(NameValidChars); e != "" {
		t.Error(e)
	}

	// equal to max length
	var longestName = strings.Repeat("a", NameLenMaxBytes)
	if e := IsValidLockName(longestName); e != "" {
		t.Error(e)
	}

	// -- invalid names

	// blank string
	if e := IsValidLockName(""); e != "Name is empty string" {
		t.Error(e)
	}

	// ascii control character
	if e := IsValidLockName("a\nb\n"); e != "Name contains non-printable/non-ascii byte 10 at offset 1" {
		t.Error(e)
	}

	// unallowed printable ascii character
	if e := IsValidLockName("A++"); e != "Name contains non-printable/non-ascii byte 43 at offset 1" {
		t.Error(e)
	}

	// non-ascii unicode
	if e := IsValidLockName(sampleNihongo); e != "Name contains non-printable/non-ascii byte 230 at offset 0" {
		t.Error(e)
	}

	// non-ascii invalid utf8
	if e := IsValidLockName(sampleInvalidUtf8); e != "Name contains non-printable/non-ascii byte 189 at offset 0" {
		t.Error(e)
	}

	// exceeding max length
	if e := IsValidLockName(longestName + "a"); e != "Name is too long (129 bytes > max of 128)" {
		t.Error(e)
	}

}
