package decimal

import (
	"math/rand/v2"
	"testing"
)

// trimNaive is the digit-at-a-time reference trimTrailingZeros replaced.
func trimNaive(a decimal128) decimal128 {
	coef, scale := a.coef, a.scale
	for scale > 0 && !coef.IsZero() {
		q, r := coef.QuoRem64(10)
		if r != 0 {
			break
		}

		coef, scale = q, scale-1
	}

	return decimal128{coef: coef, neg: a.neg, scale: scale}
}

// The binary search must agree with the naive loop for every mix of
// significant digits, trailing zeros and scale, including coefficients that
// fill the high word.
func TestTrimTrailingZerosMatchesNaive(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))

	for range 20000 {
		// A random significand times 10^zeros, at a random scale.
		sig := r.Uint64() >> r.IntN(64)
		zeros := uint8(r.IntN(20))
		coef, ok := u128FromU64(sig).Mul(pow10[zeros])
		if !ok {
			continue
		}

		if r.IntN(4) == 0 {
			// Push some cases into the high word.
			if c, ok := coef.Mul64(1 << 40); ok {
				coef = c
			}
		}

		a := decimal128{coef: coef, neg: r.IntN(2) == 0, scale: uint8(r.IntN(int(maxScale) + 1))}

		got, want := a.trimTrailingZeros(), trimNaive(a)
		if got != want {
			t.Fatalf("trim(%v at scale %d) = %v at scale %d, want %v at scale %d",
				a.coef, a.scale, got.coef, got.scale, want.coef, want.scale)
		}
	}
}
