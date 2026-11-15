package decimal

import (
	"math/big"
	"math/rand/v2"
	"testing"
)

// randPaddedDecimal returns a value whose coefficient ends in many zeros at a
// high scale — the shape a quotient takes (0.5 is 0.5000000000000000000) and
// the one Add and Mul's overflow retry exists for.
func randPaddedDecimal(r *rand.Rand) Decimal {
	sig := r.Uint64() >> r.IntN(64)
	zeros := uint8(r.IntN(int(maxScale) + 1))

	coef, ok := u128FromU64(sig).Mul(pow10[zeros])
	if !ok {
		coef = u128FromU64(sig)
	}

	return Decimal{newDec(r.IntN(2) == 0, coef, uint8(r.IntN(int(maxScale)+1)))}
}

// randLargeDecimal returns a value with a big integer part and no fraction.
func randLargeDecimal(r *rand.Rand) Decimal {
	coef := u128{hi: r.Uint64() >> r.IntN(64), lo: r.Uint64()}

	return Decimal{newDec(r.IntN(2) == 0, coef, 0)}
}

// Whenever the retry lets Add or Mul succeed, the answer must be exactly the
// one math/big gives (Mul truncated to 19 places, as always). The test also
// checks that the retry path actually ran, so it can't pass vacuously.
func TestOverflowRetryMatchesBig(t *testing.T) {
	r := rand.New(rand.NewPCG(3, 4))
	retried := 0

	for range 50000 {
		a, b := randPaddedDecimal(r), randLargeDecimal(r)
		if r.IntN(2) == 0 {
			a, b = b, a
		}

		// Did the un-retried path overflow?
		_, firstErr := a.value.add(b.value, false)

		sum, err := a.TryAdd(b)
		if err == nil {
			if firstErr != nil {
				retried++
			}

			want := new(big.Rat).Add(bigFromDecimal(t, a), bigFromDecimal(t, b))
			if got := bigFromDecimal(t, sum); got.Cmp(want) != 0 {
				t.Fatalf("%s + %s = %s, want %s", a, b, sum, want.FloatString(19))
			}
		}

		prod, err := a.TryMul(b)
		if err == nil {
			exact := new(big.Rat).Mul(bigFromDecimal(t, a), bigFromDecimal(t, b))
			if got := bigFromDecimal(t, prod); got.Cmp(truncateRat(exact, int(maxScale))) != 0 {
				t.Fatalf("%s × %s = %s, want %s", a, b, prod, exact.FloatString(25))
			}
		}
	}

	if retried == 0 {
		t.Fatal("no case exercised the overflow retry")
	}
}

// truncateRat truncates x toward zero to places fractional digits.
func truncateRat(x *big.Rat, places int) *big.Rat {
	scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(places)), nil)
	scaled := new(big.Int).Mul(x.Num(), scale)
	scaled.Quo(scaled, x.Denom()) // Quo truncates toward zero

	return new(big.Rat).SetFrac(scaled, scale)
}
