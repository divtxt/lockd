package statemachine

import (
	"fmt"
)

// NameLenMaxBytes is the maximum number of bytes allowed in a lock name.
const NameLenMaxBytes = 128

// IsValidLockName checks if the given name is a valid lock name.
//
// Returns blank string if names is valid or string message describing why name is invalid.
//
// Lock names must satisfy the following conditions:
// - must be printable ascii (byte values 32 to 126)
// - must not be empty string
// - must not be longer than NameLenMaxBytes bytes
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
		if c < 32 || c > 126 {
			return fmt.Sprintf("Name contains non-printable/non-ascii byte %v at offset %v", c, i)
		}
	}
	return ""
}
