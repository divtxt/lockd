package statemachine_test

import (
	"github.com/divtxt/lockd/statemachine"
	"testing"
)

type f_name_expect_bool func(string, bool)

func TestInMemoryLSM(t *testing.T) {

	var lsm statemachine.LockStateMachine = statemachine.NewInMemoryLSM()

	f_withErrCheck := func(f func(name string) (bool, error)) f_name_expect_bool {
		return func(name string, expected bool) {
			actual, err := f(name)
			if err != nil {
				t.Fatal(err)
			}
			if actual != expected {
				t.Fatal(actual)
			}
		}
	}

	lsm_IsLocked := f_withErrCheck(lsm.IsLocked)
	lsm_Lock := f_withErrCheck(lsm.Lock)
	lsm_Unlock := f_withErrCheck(lsm.Unlock)

	lsm_IsLocked("foo", false)
	lsm_IsLocked("bar", false)

	lsm_Lock("foo", true)
	lsm_IsLocked("foo", true)
	lsm_IsLocked("bar", false)

	lsm_Lock("foo", false)
	lsm_Lock("bar", true)
	lsm_IsLocked("bar", true)

	lsm_Unlock("bar", true)
	lsm_Unlock("foo", true)
	lsm_Unlock("bar", false)
}
