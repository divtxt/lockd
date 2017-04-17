package statemachine

// NameLenMaxBytes is the maximum number of bytes allowed in a lock name.
const NameLenMaxBytes = 128

// IsValidName checks if the given name is a valid lock name.
//
// Returns blank string if names is valid or message describing why lock name is invalid.
//
// The current implementation requires that a lock name must be:
// - non-empty string
// - at most NameLenMaxBytes bytes long
// - printable ascii as checked by IsPrintableASCII().
//
func IsValidName(name string) string {
	if name == "" {
		return "Name is blank"
	}
	if len(name) > NameLenMaxBytes {
		return "Name is too long"
	}
	if !IsPrintableASCII(name) {
		return "Name is not printable ASCII"
	}
	return ""
}

// IsPrintableASCII checks if the given string contains only printable ASCII characters.
//
// Printable ASCII is defined as byte values 32 to 126.
//
func IsPrintableASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < 32 || c > 126 {
			return false
		}
	}
	return true
}
