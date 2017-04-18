package statemachine_test

import (
	"testing"

	"github.com/divtxt/lockd/statemachine"
	"github.com/divtxt/lockd/util"
	"github.com/divtxt/raft/testhelpers"
)

const sampleAllValidChars = "$-.0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz"
const sampleNihongo = "日本語"
const sampleInvalidUtf8 = "\xbd\xb2\x3d\xbc\x20\xe2\x8c\x98"

func TestInMemoryLSM(t *testing.T) {

	var lsm statemachine.LockStateMachine = statemachine.NewInMemoryLSM()

	type f_name_expect_bool func(string, bool)
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
	lsm_IsLocked(util.NameValidChars, false)
	lsm_Lock(util.NameValidChars, true)
	lsm_IsLocked(util.NameValidChars, true)
	lsm_Unlock(util.NameValidChars, true)
	lsm_IsLocked(util.NameValidChars, false)

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
