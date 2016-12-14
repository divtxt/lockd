package runner_test

import (
	"fmt"
	"github.com/divtxt/lockd/runner"
	"testing"
	"time"
)

func TestTriggeredRunner(t *testing.T) {
	fRunning := false
	fRunCount := 0
	f := func() {
		fRunning = true
		time.Sleep(20 * time.Millisecond)
		fRunning = false
		fRunCount++
	}

	tr := runner.NewTriggeredRunner(f)
	if fRunning || fRunCount != 0 {
		t.Fatal()
	}

	// No first run without trigger
	time.Sleep(10 * time.Millisecond)
	if fRunning || fRunCount != 0 {
		t.Fatal()
	}

	// TriggerRun should run the function
	tr.TriggerRun() // Run #1
	time.Sleep(10 * time.Millisecond)
	if !fRunning {
		t.Fatal()
	}
	time.Sleep(20 * time.Millisecond)
	if fRunCount != 1 {
		t.Fatal()
	}

	// Extra TriggerRuns should return immediately and collapse into 1 run
	tr.TriggerRun() // Run #2
	time.Sleep(10 * time.Millisecond)
	if !fRunning {
		t.Fatal()
	}
	tr.TriggerRun() // Run #3
	tr.TriggerRun() // should be collapsed into pending run #3
	tr.TriggerRun() // should be collapsed into pending run #3
	time.Sleep(20 * time.Millisecond)
	if fRunCount != 2 {
		t.Fatal()
	}
	if !fRunning {
		t.Fatal()
	}
	time.Sleep(20 * time.Millisecond)
	if fRunning || fRunCount != 3 {
		t.Fatal()
	}

	// StopSync waits for current run to complete
	tr.TriggerRun() // Run #4
	time.Sleep(10 * time.Millisecond)
	if !fRunning {
		t.Fatal()
	}
	tr.StopSync()
	if fRunning || fRunCount != 4 {
		t.Fatal()
	}

	// TriggerRun & StopSync should now panic
	test_ExpectPanic(t, tr.TriggerRun, "send on closed channel")
	test_ExpectPanic(t, tr.StopSync, "close of closed channel")

	// TestHelperFakeRestart & TestHelperRunOnceIfTriggerPending
	tr.TestHelperFakeRestart()
	if tr.TestHelperRunOnceIfTriggerPending() {
		t.Fatal()
	}
	tr.TriggerRun() // Run #5
	time.Sleep(10 * time.Millisecond)
	if fRunning || fRunCount != 4 {
		t.Fatal()
	}
	if !tr.TestHelperRunOnceIfTriggerPending() {
		t.Fatal()
	}
	if fRunning || fRunCount != 5 {
		t.Fatal()
	}
}

func test_ExpectPanic(t *testing.T, f func(), expectedRecover string) {
	skipRecover := false
	defer func() {
		if !skipRecover {
			if r := recover(); fmt.Sprintf("%v", r) != expectedRecover {
				t.Fatal(fmt.Sprintf("Expected panic: %v; got: %v", expectedRecover, r))
			}
		}
	}()

	f()
	skipRecover = true
	t.Fatal(fmt.Sprintf("Expected panic: %v; got nothing!", expectedRecover))
}
