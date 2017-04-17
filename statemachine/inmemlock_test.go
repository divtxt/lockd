package statemachine_test

import (
	"testing"

	"github.com/divtxt/lockd/statemachine"
	"github.com/divtxt/raft/testhelpers"
)

type f_name_expect_bool func(string, bool)

func TestInMemoryLSM(t *testing.T) {

	var lsm statemachine.LockStateMachine = statemachine.NewInMemoryLSM()

	fExpectResult := func(f func(name string) bool) f_name_expect_bool {
		return func(name string, expected bool) {
			actual := f(name)
			if actual != expected {
				t.Fatal(actual)
			}
		}
	}

	lsm_IsLocked := fExpectResult(lsm.IsLocked)
	lsm_Lock := fExpectResult(lsm.Lock)
	lsm_Unlock := fExpectResult(lsm.Unlock)

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

	// all printable ascii
	lsm_IsLocked(sampleAllPrintableAscii, false)
	lsm_Lock(sampleAllPrintableAscii, true)
	lsm_IsLocked(sampleAllPrintableAscii, true)
	lsm_Unlock(sampleAllPrintableAscii, true)
	lsm_IsLocked(sampleAllPrintableAscii, false)

	// bad names should panic
	testhelpers.TestHelper_ExpectPanic(
		t,
		func() { lsm.IsLocked("") },
		"Name is empty string",
	)
	testhelpers.TestHelper_ExpectPanic(
		t,
		func() { lsm.Lock("") },
		"Name is empty string",
	)
	testhelpers.TestHelper_ExpectPanic(
		t,
		func() { lsm.IsLocked(sampleNihongo) },
		"Name contains non-printable/non-ascii byte 230 at offset 0",
	)
	testhelpers.TestHelper_ExpectPanic(
		t,
		func() { lsm.Unlock(sampleInvalidUtf8) },
		"Name contains non-printable/non-ascii byte 189 at offset 0",
	)
}
