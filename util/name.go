package util

import (
	"fmt"
)

// NameLenMaxBytes is the maximum number of bytes allowed in a lock name.
const NameLenMaxBytes = 128

// All valid chars
const NameValidChars = "$-.0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz"

// IsValidLockName checks if the given name is a valid lock name.
//
// Returns blank string if names is valid or string message describing why name is invalid.
//
// Lock names can only contain the following characters:
// - letters A-Z & a-z
// - digits 0-9
// - dollar sign ("$"), dash ("-"), underscores ("_") or period (".")
//
// Lock names also cannot be an empty string or be longer than NameLenMaxBytes bytes.
//
func IsValidLockName(name string) string {
	l := len(name)
	if l == 0 {
		return "Name is empty string"
	}
	if l > NameLenMaxBytes {
		return fmt.Sprintf("Name is too long (%v bytes > max of %v)", l, NameLenMaxBytes)
	}
	for i := 0; i < l; i++ {
		c := name[i]
		if !allowedChars[c] {
			return fmt.Sprintf("Name contains non-printable/non-ascii byte %v at offset %v", c, i)
		}
	}
	return ""
}

func calcAllowedCharsTable() [256]bool {
	var ac [256]bool
	allowChars := func(s string) {
		for i := 0; i < len(s); i++ {
			ac[s[i]] = true
		}
	}
	allowChars(NameValidChars)
	return ac
}

var allowedChars [256]bool = calcAllowedCharsTable()
