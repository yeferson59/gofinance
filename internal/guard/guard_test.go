package guard

import (
	"errors"
	"testing"
)

var errOutOfRange = errors.New("numeric overflow")

func guarded(f func()) (err error) {
	defer Recover(&err)

	f()

	return nil
}

func TestRecoverTurnsErrorPanicIntoError(t *testing.T) {
	err := guarded(func() { panic(errOutOfRange) })
	if !errors.Is(err, errOutOfRange) {
		t.Fatalf("err = %v, want errOutOfRange", err)
	}
}

func TestRecoverLeavesNormalReturnAlone(t *testing.T) {
	if err := guarded(func() {}); err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
}

// Bugs keep panicking: recovering them would hide a defect behind an error
// the caller can't do anything about.
func TestRecoverRepanicsOnRuntimeError(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("a runtime error was swallowed")
		}
	}()

	_ = guarded(func() {
		var s []int
		_ = s[1]
	})
}

func TestRecoverRepanicsOnNonErrorValue(t *testing.T) {
	defer func() {
		if r := recover(); r != "boom" {
			t.Fatalf("recovered %v, want the original value", r)
		}
	}()

	_ = guarded(func() { panic("boom") })
}
