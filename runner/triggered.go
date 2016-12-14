package runner

import (
	"sync"
)

// A TriggeredRunner runs a function in a goroutine every time it is triggered.
type TriggeredRunner struct {
	f       func()
	trigger chan struct{}
	wg      sync.WaitGroup
}

// NewTriggeredRunner creates a TriggeredRunner for the given function.
//
// A goroutine is started that will run the given function when triggered.
func NewTriggeredRunner(f func()) *TriggeredRunner {
	tr := &TriggeredRunner{
		f,
		make(chan struct{}, 1),
		sync.WaitGroup{},
	}
	tr.wg.Add(1)
	go tr.run()
	return tr
}

func (tr *TriggeredRunner) run() {
	defer tr.wg.Done()
	for {
		select {
		case _, ok := <-tr.trigger:
			if !ok {
				return
			}
			tr.f()
		}
	}
}

// TriggerRun triggers a run of the wrapped function in the goroutine.
//
// Returns immediately without waiting for the function to run.
//
// If the function is currently running, this trigger will be pending and result in the another
// run once the current run completes.
// Multiple such pending triggers will be collapsed into a single run.
//
// Will panic if StopSync() has been called.
func (tr *TriggeredRunner) TriggerRun() {
	select {
	case tr.trigger <- struct{}{}:
	default: // avoid blocking
	}
}

// StopSync will stop the goroutine, waiting for any current run to complete.
//
// Will panic if called more than once.
func (tr *TriggeredRunner) StopSync() {
	close(tr.trigger)
	tr.wg.Wait()
}

// TestHelperFakeRestart is meant for testing use only.
//
// It resets the state so that TriggerRun() calls will succeed but does not start a new goroutine
// to actually run the function when triggered.
// Use TestHelperRunOnceIfTriggerPending() to actually run the function.
//
// You should have called StopSync() before using this.
func (tr *TriggeredRunner) TestHelperFakeRestart() {
	tr.trigger = make(chan struct{}, 1)
}

// TestHelperRunOnceIfTriggerPending is meant for testing use only.
//
// It runs the function if a trigger pending, and returns a value indicating if it ran.
//
// You should be using this with TestHelperFakeRestart().
func (tr *TriggeredRunner) TestHelperRunOnceIfTriggerPending() bool {
	select {
	case <-tr.trigger:
		tr.f()
		return true
	default:
		return false
	}
}
