package backend_test

import (
	"fmt"
	"testing"

	"github.com/divtxt/raft"

	"github.com/divtxt/lockd/backend"
	"github.com/divtxt/lockd/util"
	"github.com/divtxt/raft/testhelpers"
)

const sampleAllValidChars = "$-.0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz"
const sampleNihongo = "日本語"
const sampleInvalidUtf8 = "\xbd\xb2\x3d\xbc\x20\xe2\x8c\x98"

func TestInMemoryBackend(t *testing.T) {

	var lb backend.LockBackend = backend.NewInMemoryBackend(0)

	lb_IsLocked := func(name string, expected bool) {
		if actual := lb.IsLocked(name); actual != expected {
			panic(fmt.Sprintf("%v", actual))
		}
	}
	lb_Lock := func(li raft.LogIndex, name string, expected bool) {
		if actual := lb.Lock(li, name); actual != expected {
			panic(fmt.Sprintf("%v", actual))
		}
		if actual := lb.GetLastApplied(); actual != li {
			panic(fmt.Sprintf("%v", actual))
		}
	}
	lb_Unlock := func(li raft.LogIndex, name string, expected bool) {
		if actual := lb.Unlock(li, name); actual != expected {
			panic(fmt.Sprintf("%v", actual))
		}
		if actual := lb.GetLastApplied(); actual != li {
			panic(fmt.Sprintf("%v", actual))
		}
	}

	lb_IsLocked("foo", false)
	lb_IsLocked("bar", false)

	lb_Lock(1, "foo", true)
	lb_IsLocked("foo", true)
	lb_IsLocked("bar", false)

	lb_Lock(2, "foo", false)
	lb_Lock(3, "bar", true)
	lb_IsLocked("bar", true)

	lb_Unlock(4, "bar", true)
	lb_Unlock(5, "foo", true)
	// can skip indexes
	lb_Unlock(8, "bar", false)

	// all printable ascii
	lb_IsLocked(util.NameValidChars, false)
	lb_Lock(11, util.NameValidChars, true)
	lb_IsLocked(util.NameValidChars, true)
	lb_Unlock(12, util.NameValidChars, true)
	lb_IsLocked(util.NameValidChars, false)

	// bad names should error
	testhelpers.TestHelper_ExpectPanic(
		t,
		func() { lb.IsLocked("") },
		"Name is empty string",
	)
	testhelpers.TestHelper_ExpectPanic(
		t,
		func() { lb.Lock(13, "") },
		"Name is empty string",
	)
	testhelpers.TestHelper_ExpectPanic(
		t,
		func() { lb.IsLocked(sampleNihongo) },
		"Name contains non-printable/non-ascii byte 230 at offset 0",
	)
	testhelpers.TestHelper_ExpectPanic(
		t,
		func() { lb.Unlock(13, sampleInvalidUtf8) },
		"Name contains non-printable/non-ascii byte 189 at offset 0",
	)
}
