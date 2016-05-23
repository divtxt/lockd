package lockd

import (
	"testing"
)

func BlackboxTest_LockStatePersistence(t *testing.T, lsp LockStatePersistence) {
	var success bool
	var err error

	lock := func(n string) bool {
		success, err = lsp.Lock(n)
		if err != nil {
			t.Fatal(err)
		}
		return success
	}

	unlock := func(n string) bool {
		success, err = lsp.Unlock(n)
		if err != nil {
			t.Fatal(err)
		}
		return success
	}

	// Simple lock & unlock combinations
	if unlock("A") != false {
		t.Fatal()
	}

	if lock("A") != true {
		t.Fatal()
	}
	if lock("B") != true {
		t.Fatal()
	}
	if lock("A") != false {
		t.Fatal()
	}

	if unlock("A") != true {
		t.Fatal()
	}
	if unlock("B") != true {
		t.Fatal()
	}
	if unlock("B") != false {
		t.Fatal()
	}

}
