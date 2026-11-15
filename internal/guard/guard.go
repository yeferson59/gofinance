// Package guard upholds the library's error contract: a function that returns
// an error never panics.
//
// The finance packages compute with the decimal engine's panicking helpers
// (Add, Mul, MustDiv, Money.Add, ...) because threading an error through every
// intermediate step of a formula would bury it. Those helpers panic with an
// error value — numeric overflow, division by zero, a currency mismatch — on
// inputs that are valid by type but out of the engine's range. Every exported
// function that returns an error defers Recover, which turns such a panic back
// into that function's error, the same way encoding/json reports errors from
// deep inside its recursive encoder.
//
// Only panics carrying an error that is not a runtime.Error are recovered.
// Anything else — a nil dereference, an index out of range, a panic with a
// non-error value — is a bug rather than an out-of-range input, and keeps
// propagating.
package guard

import (
	"errors"
	"runtime"
)

// Recover converts a panic carrying an error value into *err. It must be
// deferred directly, so that its own recover call stops the panic:
//
//	func F() (_ decimal.Decimal, err error) {
//		defer guard.Recover(&err)
//		...
//	}
func Recover(err *error) {
	r := recover()
	if r == nil {
		return
	}

	if e := recoverable(r); e != nil {
		*err = e
		return
	}

	panic(r)
}

// recoverable returns the error a recovered panic value carries when Recover
// should turn it into an error, and nil when the panic must keep propagating.
func recoverable(r any) error {
	e, ok := r.(error)
	if !ok {
		return nil
	}

	var runtimeErr runtime.Error
	if errors.As(e, &runtimeErr) {
		return nil
	}

	return e
}
