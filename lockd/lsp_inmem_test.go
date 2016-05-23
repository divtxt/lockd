package lockd

import (
	"testing"
)

func TestInMemoryLSP(t *testing.T) {
	lsp := NewInMemoryLSP()
	BlackboxTest_LockStatePersistence(t, lsp)
}
