package statemachine_test

import (
	"strings"
	"testing"

	. "github.com/divtxt/lockd/statemachine"
)

const sampleAllPrintableAscii = " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[]^_`abcdefghijklmnopqrstuvwxyz{|}~"
const sampleNihongo = "日本語"
const sampleInvalidUtf8 = "\xbd\xb2\x3d\xbc\x20\xe2\x8c\x98"

func TestIsValidName(t *testing.T) {

	// -- valid names:

	// random ascii string
	if IsValidName("hello, world!") != "" {
		t.Error()
	}

	// all printable from 32 to 126 ascii
	if IsValidName(sampleAllPrintableAscii) != "" {
		t.Error()
	}

	// equal to max length
	var longestName = strings.Repeat("a", NameLenMaxBytes)
	if IsValidName(longestName) != "" {
		t.Error()
	}

	// -- invalid names

	// blank string
	if IsValidName("") == "" {
		t.Error()
	}

	// non-ascii unicode
	if IsValidName(sampleNihongo) == "" {
		t.Error()
	}

	// non-ascii invalid utf8
	if IsValidName(sampleInvalidUtf8) == "" {
		t.Error()
	}

	// exceeding max length
	if IsValidName(longestName+"a") == "" {
		t.Error()
	}

}

func TestIsPrintableASCII(t *testing.T) {
	// blank string
	if !IsPrintableASCII("") {
		t.Error()
	}

	// some random ascii strings
	if !IsPrintableASCII("abcABC123") {
		t.Error()
	}
	if !IsPrintableASCII("hello, world!") {
		t.Error()
	}

	// all printable from 32 to 126 ascii
	if !IsPrintableASCII(sampleAllPrintableAscii) {
		t.Error()
	}

	// non-ascii unicode
	if IsPrintableASCII(sampleNihongo) {
		t.Error()
	}

	// non-ascii invalid utf8
	if IsPrintableASCII(sampleInvalidUtf8) {
		t.Error()
	}
}
