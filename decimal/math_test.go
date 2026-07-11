package decimal

import (
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"testing"
)

// refPrec is the math/big precision used by the reference implementations
// below: comfortably beyond fp120's ~36 significant digits.
const refPrec = 400

func refFloat(f float64) *big.Float {
	return big.NewFloat(f).SetPrec(refPrec)
}

// refLnSeries computes ln(x) for x near 1 via 2*atanh((x-1)/(x+1)).
func refLnSeries(x *big.Float) *big.Float {
	one := refFloat(1)
	num := new(big.Float).SetPrec(refPrec).Sub(x, one)
	den := new(big.Float).SetPrec(refPrec).Add(x, one)
	z := new(big.Float).SetPrec(refPrec).Quo(num, den)
	zsq := new(big.Float).SetPrec(refPrec).Mul(z, z)
	sum := new(big.Float).SetPrec(refPrec).Set(z)
	term := new(big.Float).SetPrec(refPrec).Set(z)
	t := new(big.Float).SetPrec(refPrec)
	old := new(big.Float).SetPrec(refPrec)

	for k := int64(3); ; k += 2 {
		term.Mul(term, zsq)
		t.Quo(term, refFloat(float64(k)))
		old.Set(sum)
		sum.Add(sum, t)
		if old.Cmp(sum) == 0 {
			break
		}
	}

	return sum.Mul(sum, refFloat(2))
}

var refLn2 = refLnSeries(refFloat(2))

// refLn computes ln(x) for any x > 0 by binary normalization plus
// refLnSeries on the mantissa.
func refLn(x *big.Float) *big.Float {
	m := new(big.Float).SetPrec(refPrec).Set(x)
	exp := m.MantExp(m) // m in [0.5, 1)

	r := refLnSeries(m)
	k := new(big.Float).SetPrec(refPrec).Mul(refLn2, refFloat(float64(exp)))

	return r.Add(r, k)
}

// refExp computes e^x by reduction x = k*ln2 + s, |s| <= ln2/2, plus the
// Taylor series for e^s.
func refExp(x *big.Float) *big.Float {
	q := new(big.Float).SetPrec(refPrec).Quo(x, refLn2)
	qi, _ := q.Int(nil)
	k := qi.Int64()

	kl := new(big.Float).SetPrec(refPrec).Mul(refLn2, refFloat(float64(k)))
	s := new(big.Float).SetPrec(refPrec).Sub(x, kl)

	sum := refFloat(1)
	term := refFloat(1)
	old := new(big.Float).SetPrec(refPrec)

	for i := int64(1); ; i++ {
		term.Mul(term, s)
		term.Quo(term, refFloat(float64(i)))
		old.Set(sum)
		sum.Add(sum, term)
		if old.Cmp(sum) == 0 {
			break
		}
	}

	// e^x = e^s * 2^k
	pow2 := new(big.Float).SetPrec(refPrec).SetMantExp(refFloat(1), int(k))

	return sum.Mul(sum, pow2)
}

func refFromString(t *testing.T, s string) *big.Float {
	t.Helper()
	f, _, err := big.ParseFloat(s, 10, refPrec, big.ToNearestEven)
	if err != nil {
		t.Fatalf("big.ParseFloat(%q): %v", s, err)
	}
	return f
}

// assertClose checks that got (a decimal128 result) is within tol of want.
func assertClose(t *testing.T, label string, got decimal128, want *big.Float, tol *big.Float) {
	t.Helper()
	g := refFromString(t, got.String())
	diff := new(big.Float).SetPrec(refPrec).Sub(g, want)
	diff.Abs(diff)
	if diff.Cmp(tol) > 0 {
		t.Fatalf("%s = %s, want %s (diff %s > tol %s)",
			label, got.String(), want.Text('g', 25), diff.Text('g', 5), tol.Text('g', 5))
	}
}

// lnTol returns the comparison tolerance for a result of magnitude |v|:
// one ulp at scale 19 plus a relative 1e-18 term.
func tolFor(v *big.Float) *big.Float {
	tol := new(big.Float).SetPrec(refPrec).Abs(v)
	tol.Mul(tol, refFloat(1e-18))
	return tol.Add(tol, refFloat(1e-19))
}

func TestFp120Constants(t *testing.T) {
	two120 := new(big.Float).SetPrec(refPrec).SetMantExp(refFloat(1), 120)

	check := func(name string, got u128, v *big.Float) {
		t.Helper()
		scaled := new(big.Float).SetPrec(refPrec).Mul(v, two120)
		scaled.Add(scaled, refFloat(0.5))
		i, _ := scaled.Int(nil)

		wantLo := new(big.Int).Mod(i, new(big.Int).Lsh(big.NewInt(1), 64)).Uint64()
		wantHi := new(big.Int).Rsh(i, 64).Uint64()
		if got.hi != wantHi || got.lo != wantLo {
			t.Errorf("%s = {hi: %#x, lo: %#x}, want {hi: %#x, lo: %#x}",
				name, got.hi, got.lo, wantHi, wantLo)
		}
	}

	ln2 := refLnSeries(refFloat(2))
	ln10 := refLn(refFloat(10))
	check("ln2Fp", ln2Fp, ln2)
	check("ln10Fp", ln10Fp, ln10)
	check("invLn2Fp", invLn2Fp, new(big.Float).SetPrec(refPrec).Quo(refFloat(1), ln2))
	check("invLn10Fp", invLn10Fp, new(big.Float).SetPrec(refPrec).Quo(refFloat(1), ln10))
	check("fpOne", fpOne, refFloat(1))
	check("fpTwo", fpTwo, refFloat(2))
}

func TestLnKnownValues(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"1", "0"},
		// ln(2) = 0.69314718055994530941723...
		{"2", "0.6931471805599453094"},
		// ln(10) = 2.30258509299404568401799...
		{"10", "2.302585092994045684"},
		{"0.5", "-0.6931471805599453094"},
		// ln(0.1) = -ln(10)
		{"0.1", "-2.302585092994045684"},
	}

	for _, tt := range tests {
		d := mustParseDec(t, tt.in)
		got, err := d.Ln()
		if err != nil {
			t.Fatalf("Ln(%s): %v", tt.in, err)
		}
		if got.String() != tt.want {
			t.Errorf("Ln(%s) = %s, want %s", tt.in, got.String(), tt.want)
		}
	}
}

func TestLnCrossCheckRandom(t *testing.T) {
	r := rand.New(rand.NewSource(1234))

	for range 500 {
		s := strings.TrimPrefix(randDecimalString(r, 12, 15), "-")
		a := mustParseDec(t, s)
		if a.IsZero() {
			continue
		}

		got, err := a.Ln()
		if err != nil {
			t.Fatalf("Ln(%s): %v", s, err)
		}

		want := refLn(refFromString(t, s))
		assertClose(t, "Ln("+s+")", got, want, tolFor(want))
	}
}

func TestLnExtremes(t *testing.T) {
	// Smallest positive decimal128 and a near-maximal one.
	small := mustParseDec(t, "0.0000000000000000001")
	got, err := small.Ln()
	if err != nil {
		t.Fatal(err)
	}
	want := refLn(refFromString(t, "1e-19"))
	assertClose(t, "Ln(1e-19)", got, want, tolFor(want))

	huge := newDec(false, u128{hi: ^uint64(0), lo: ^uint64(0)}, 0)
	got, err = huge.Ln()
	if err != nil {
		t.Fatal(err)
	}
	want = refLn(refFromString(t, huge.String()))
	assertClose(t, "Ln(2^128-1)", got, want, tolFor(want))

	// Values adjacent to 1 exercise the tiny-magnitude path.
	for _, s := range []string{"1.0000000000000000001", "0.9999999999999999999"} {
		got, err := mustParseDec(t, s).Ln()
		if err != nil {
			t.Fatal(err)
		}
		want := refLn(refFromString(t, s))
		assertClose(t, "Ln("+s+")", got, want, tolFor(refFloat(1)))
	}
}

func TestLogErrors(t *testing.T) {
	if _, err := decZero.Ln(); !errors.Is(err, ErrLogNonPositive) {
		t.Errorf("Ln(0): got %v, want ErrLogNonPositive", err)
	}
	if _, err := mustParseDec(t, "-3").Ln(); !errors.Is(err, ErrLogNonPositive) {
		t.Errorf("Ln(-3): got %v, want ErrLogNonPositive", err)
	}
	if _, err := mustParseDec(t, "-3").Log2(); !errors.Is(err, ErrLogNonPositive) {
		t.Errorf("Log2(-3): got %v, want ErrLogNonPositive", err)
	}
	if _, err := decZero.Log10(); !errors.Is(err, ErrLogNonPositive) {
		t.Errorf("Log10(0): got %v, want ErrLogNonPositive", err)
	}
	if _, err := mustParseDec(t, "8").Log(decOne); !errors.Is(err, ErrDivideByZero) {
		t.Errorf("Log(8, base 1): got %v, want ErrDivideByZero", err)
	}
	if _, err := mustParseDec(t, "8").Log(decZero); !errors.Is(err, ErrLogNonPositive) {
		t.Errorf("Log(8, base 0): got %v, want ErrLogNonPositive", err)
	}
}

func TestLogExactCases(t *testing.T) {
	cases := []struct {
		fn   func(decimal128) (decimal128, error)
		name string
		in   string
		want string
	}{
		{decimal128.Log10, "Log10", "1000", "3"},
		{decimal128.Log10, "Log10", "0.001", "-3"},
		{decimal128.Log10, "Log10", "100000000000000000000000000000000000000", "38"},
		{decimal128.Log2, "Log2", "8", "3"},
		{decimal128.Log2, "Log2", "0.25", "-2"},
		{decimal128.Log2, "Log2", "170141183460469231731687303715884105728", "127"}, // 2^127
		{decimal128.Ln, "Ln", "1", "0"},
	}

	for _, tt := range cases {
		got, err := tt.fn(mustParseDec(t, tt.in))
		if err != nil {
			t.Fatalf("%s(%s): %v", tt.name, tt.in, err)
		}
		if got.String() != tt.want {
			t.Errorf("%s(%s) = %s, want %s", tt.name, tt.in, got.String(), tt.want)
		}
	}
}

func TestLogBaseCrossCheck(t *testing.T) {
	// Exact cases.
	got, err := mustParseDec(t, "8").Log(mustParseDec(t, "2"))
	if err != nil || got.String() != "3" {
		t.Fatalf("Log(8, 2) = %s, %v; want 3", got.String(), err)
	}

	got, err = mustParseDec(t, "0.01").Log(mustParseDec(t, "10"))
	if err != nil || got.String() != "-2" {
		t.Fatalf("Log(0.01, 10) = %s, %v; want -2", got.String(), err)
	}

	// Random cross-check against ln(a)/ln(base).
	r := rand.New(rand.NewSource(77))
	for range 200 {
		as := strings.TrimPrefix(randDecimalString(r, 6, 8), "-")
		bs := strings.TrimPrefix(randDecimalString(r, 3, 4), "-")
		a, b := mustParseDec(t, as), mustParseDec(t, bs)
		if a.IsZero() || b.IsZero() || b.Equal(decOne) {
			continue
		}

		got, err := a.Log(b)
		if err != nil {
			t.Fatalf("Log(%s, %s): %v", as, bs, err)
		}

		want := new(big.Float).SetPrec(refPrec).Quo(refLn(refFromString(t, as)), refLn(refFromString(t, bs)))
		assertClose(t, "Log("+as+","+bs+")", got, want, tolFor(want))
	}
}

// ratPow computes a^n exactly for integer n (n may be negative).
func ratPow(a *big.Rat, n int64) *big.Rat {
	res := new(big.Rat).SetInt64(1)
	base := new(big.Rat).Set(a)
	neg := n < 0
	if neg {
		n = -n
	}
	for ; n > 0; n >>= 1 {
		if n&1 == 1 {
			res.Mul(res, base)
		}
		base.Mul(base, base)
	}
	if neg {
		res.Inv(res)
	}
	return res
}

// ratStringHalfAway formats r rounded half away from zero to prec
// fractional digits, with trailing zeros trimmed (decimal128.String form).
func ratStringHalfAway(r *big.Rat, prec int) string {
	neg := r.Sign() < 0
	abs := new(big.Rat).Abs(r)

	scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(prec)), nil)
	num := new(big.Int).Mul(abs.Num(), scale)
	q, rem := new(big.Int).QuoRem(num, abs.Denom(), new(big.Int))
	if new(big.Int).Lsh(rem, 1).Cmp(abs.Denom()) >= 0 {
		q.Add(q, big.NewInt(1))
	}

	s := q.String()
	for len(s) <= prec {
		s = "0" + s
	}
	intPart, fracPart := s[:len(s)-prec], s[len(s)-prec:]
	fracPart = strings.TrimRight(fracPart, "0")

	out := intPart
	if fracPart != "" {
		out += "." + fracPart
	}
	if out == "0" {
		return "0"
	}
	if neg {
		out = "-" + out
	}
	return out
}

func TestPowIntegerCrossCheck(t *testing.T) {
	cases := []struct {
		base string
		exp  int64
	}{
		{"2", 10},
		{"1.05", 12},
		{"1.005", 12},
		{"0.5", 3},
		{"1.01", 100},
		{"1.000001", 10000},
		{"-2", 3},
		{"-2", 2},
		{"-1.5", 5},
		{"10", -2},
		{"3", -7},
		{"0.99", 250},
		{"123.456", 7},
	}

	for _, tt := range cases {
		a := mustParseDec(t, tt.base)
		e, err := decFromInt64(tt.exp, 0)
		if err != nil {
			t.Fatal(err)
		}

		got, err := a.Pow(e)
		if err != nil {
			t.Fatalf("Pow(%s, %d): %v", tt.base, tt.exp, err)
		}

		want := ratStringHalfAway(ratPow(toRat(t, tt.base), tt.exp), int(maxScale))
		if got.String() != want {
			t.Errorf("Pow(%s, %d) = %s, want %s", tt.base, tt.exp, got.String(), want)
		}
	}
}

func TestPowIntegerRandomCrossCheck(t *testing.T) {
	r := rand.New(rand.NewSource(4242))

	for range 300 {
		bs := randDecimalString(r, 2, 4)
		n := int64(r.Intn(41)) - 20 // [-20, 20]
		a := mustParseDec(t, bs)
		if a.IsZero() {
			continue
		}

		e, err := decFromInt64(n, 0)
		if err != nil {
			t.Fatal(err)
		}

		got, err := a.Pow(e)
		exact := ratPow(toRat(t, bs), n)

		// Skip combinations whose true magnitude exceeds the type's range.
		limit := new(big.Rat).SetInt(new(big.Int).Lsh(big.NewInt(1), 128))
		if new(big.Rat).Abs(exact).Cmp(limit) >= 0 {
			if err == nil {
				t.Fatalf("Pow(%s, %d): expected overflow, got %s", bs, n, got.String())
			}
			continue
		}
		if err != nil {
			t.Fatalf("Pow(%s, %d): %v", bs, n, err)
		}

		// Compare at the scale the result was produced at. Results whose
		// exact value fits in 38 significant digits must match exactly;
		// beyond that the 38-digit intermediates guarantee ~36 correct
		// significant digits, so allow a relative 1e-33.
		sc := int(got.scale)
		want := ratStringHalfAway(exact, sc)
		if got.String() != want {
			exactF := new(big.Float).SetPrec(refPrec).SetRat(exact)
			tol := new(big.Float).SetPrec(refPrec).Abs(exactF)
			tol.Mul(tol, refFloat(1e-33))
			assertClose(t, fmt.Sprintf("Pow(%s, %d)", bs, n), got, exactF, tol)
		}
	}
}

func TestPowFractionalCrossCheck(t *testing.T) {
	cases := [][2]string{
		{"4", "0.5"},
		{"9", "0.5"},
		{"2", "0.5"},
		{"27", "0.3333333333333333333"},
		{"1.05", "2.5"},
		{"100", "1.5"},
		{"0.25", "-0.5"},
		{"2", "10.25"},
		{"1.1", "-3.7"},
	}

	for _, tt := range cases {
		a, e := mustParseDec(t, tt[0]), mustParseDec(t, tt[1])

		got, err := a.Pow(e)
		if err != nil {
			t.Fatalf("Pow(%s, %s): %v", tt[0], tt[1], err)
		}

		x := new(big.Float).SetPrec(refPrec).Mul(refFromString(t, tt[1]), refLn(refFromString(t, tt[0])))
		want := refExp(x)
		assertClose(t, "Pow("+tt[0]+","+tt[1]+")", got, want, tolFor(want))
	}
}

func TestPowFractionalRandomCrossCheck(t *testing.T) {
	r := rand.New(rand.NewSource(99))

	for range 300 {
		bs := strings.TrimPrefix(randDecimalString(r, 3, 6), "-")
		es := randDecimalString(r, 1, 4)
		a, e := mustParseDec(t, bs), mustParseDec(t, es)
		if a.IsZero() {
			continue
		}

		got, err := a.Pow(e)
		if err != nil {
			// Overflow is legitimate for large bases with exponents near ±10.
			if errors.Is(err, ErrOverflow) {
				continue
			}
			t.Fatalf("Pow(%s, %s): %v", bs, es, err)
		}

		x := new(big.Float).SetPrec(refPrec).Mul(refFromString(t, es), refLn(refFromString(t, bs)))
		want := refExp(x)
		assertClose(t, "Pow("+bs+","+es+")", got, want, tolFor(want))
	}
}

func TestPowSpecialCases(t *testing.T) {
	one := decOne
	zero := decZero

	if got, err := zero.Pow(zero); err != nil || !got.Equal(one) {
		t.Errorf("0^0 = %s, %v; want 1", got.String(), err)
	}
	if got, err := zero.Pow(mustParseDec(t, "5")); err != nil || !got.IsZero() {
		t.Errorf("0^5 = %s, %v; want 0", got.String(), err)
	}
	if _, err := zero.Pow(mustParseDec(t, "-1")); !errors.Is(err, ErrDivideByZero) {
		t.Errorf("0^-1: got %v, want ErrDivideByZero", err)
	}
	if _, err := mustParseDec(t, "-2").Pow(mustParseDec(t, "0.5")); !errors.Is(err, ErrPowNegBase) {
		t.Errorf("(-2)^0.5: got %v, want ErrPowNegBase", err)
	}
	if got, err := mustParseDec(t, "-2").Pow(mustParseDec(t, "2.0")); err != nil || got.String() != "4" {
		// "2.0" still is an integer exponent after trimming.
		t.Errorf("(-2)^2.0 = %s, %v; want 4", got.String(), err)
	}
	if got, err := mustParseDec(t, "987.654").Pow(zero); err != nil || !got.Equal(one) {
		t.Errorf("987.654^0 = %s, %v; want 1", got.String(), err)
	}
	if got, err := one.Pow(mustParseDec(t, "123456789.123")); err != nil || got.String() != "1" {
		t.Errorf("1^big = %s, %v; want 1", got.String(), err)
	}
	if got, err := mustParseDec(t, "-1").Pow(mustParseDec(t, "12345678912345678917")); err != nil || got.String() != "-1" {
		t.Errorf("(-1)^odd = %s, %v; want -1", got.String(), err)
	}
}

func TestPowOverflowAndUnderflow(t *testing.T) {
	ten := mustParseDec(t, "10")
	two := mustParseDec(t, "2")

	// 10^38 is the largest exact power of ten decimal128 can hold.
	got, err := ten.Pow(mustParseDec(t, "38"))
	if err != nil {
		t.Fatalf("10^38: %v", err)
	}
	if want := "1" + strings.Repeat("0", 38); got.String() != want {
		t.Errorf("10^38 = %s, want %s", got.String(), want)
	}

	if _, err := ten.Pow(mustParseDec(t, "39")); !errors.Is(err, ErrOverflow) {
		t.Errorf("10^39: got %v, want ErrOverflow", err)
	}
	if _, err := two.Pow(mustParseDec(t, "500")); !errors.Is(err, ErrOverflow) {
		t.Errorf("2^500: got %v, want ErrOverflow", err)
	}

	// Way below the smallest representable value: rounds to zero.
	for _, e := range []string{"-25", "-500"} {
		got, err := ten.Pow(mustParseDec(t, e))
		if err != nil || !got.IsZero() {
			t.Errorf("10^%s = %s, %v; want 0", e, got.String(), err)
		}
	}

	// 10^-19 is still representable exactly.
	got, err = ten.Pow(mustParseDec(t, "-19"))
	if err != nil || got.String() != "0.0000000000000000001" {
		t.Errorf("10^-19 = %s, %v; want 0.0000000000000000001", got.String(), err)
	}
}

func TestSqrtExactSquares(t *testing.T) {
	cases := [][2]string{
		{"0", "0"},
		{"1", "1"},
		{"4", "2"},
		{"100", "10"},
		{"2.25", "1.5"},
		{"0.25", "0.5"},
		{"0.0625", "0.25"},
		{"152.2756", "12.34"},
		{"100000000000000000000000000000000000000", "10000000000000000000"}, // sqrt(10^38) = 10^19
	}

	for _, tt := range cases {
		got, err := mustParseDec(t, tt[0]).Sqrt()
		if err != nil {
			t.Fatalf("Sqrt(%s): %v", tt[0], err)
		}
		if got.String() != tt[1] {
			t.Errorf("Sqrt(%s) = %s, want %s", tt[0], got.String(), tt[1])
		}
	}
}

func TestSqrtKnownValues(t *testing.T) {
	cases := [][2]string{
		// sqrt(2) = 1.4142135623730950488016887...
		{"2", "1.4142135623730950488"},
		// sqrt(3) = 1.7320508075688772935274463...
		{"3", "1.7320508075688772935"},
		// sqrt(0.5) = 0.7071067811865475244008443...
		{"0.5", "0.7071067811865475244"},
		// sqrt(10) = 3.1622776601683793319988935... rounds up, zero trimmed
		{"10", "3.162277660168379332"},
		// smallest positive value: sqrt(1e-19) = 3.1622776601...e-10, whose
		// scale-19 coefficient rounds to 3162277660
		{"0.0000000000000000001", "0.000000000316227766"},
	}

	for _, tt := range cases {
		got, err := mustParseDec(t, tt[0]).Sqrt()
		if err != nil {
			t.Fatalf("Sqrt(%s): %v", tt[0], err)
		}
		if got.String() != tt[1] {
			t.Errorf("Sqrt(%s) = %s, want %s", tt[0], got.String(), tt[1])
		}
	}
}

func TestSqrtNegative(t *testing.T) {
	if _, err := mustParseDec(t, "-4").Sqrt(); !errors.Is(err, ErrSqrtNegative) {
		t.Errorf("Sqrt(-4): got %v, want ErrSqrtNegative", err)
	}
}

// TestSqrtCrossCheckRandom checks correct rounding against big.Float.Sqrt:
// since exact ties are impossible, the result must always be within half
// an ulp of the true root.
func TestSqrtCrossCheckRandom(t *testing.T) {
	r := rand.New(rand.NewSource(31))
	halfUlp := refFromString(t, "0.5000001e-19")

	for range 2000 {
		s := strings.TrimPrefix(randDecimalString(r, 12, 15), "-")
		a := mustParseDec(t, s)

		got, err := a.Sqrt()
		if err != nil {
			t.Fatalf("Sqrt(%s): %v", s, err)
		}

		want := new(big.Float).SetPrec(refPrec).Sqrt(refFromString(t, s))
		assertClose(t, "Sqrt("+s+")", got, want, halfUlp)
	}
}

// TestSqrtCrossCheckSpec replicates the rounding spec with big.Int over
// extreme coefficients: coefficient = nearest integer to
// sqrt(coef * 10^(38-scale)).
func TestSqrtCrossCheckSpec(t *testing.T) {
	r := rand.New(rand.NewSource(32))

	for range 2000 {
		a := newDec(false, u128{hi: r.Uint64(), lo: r.Uint64()}, uint8(r.Intn(20)))

		got, err := a.Sqrt()
		if err != nil {
			t.Fatalf("Sqrt(%s): %v", a.String(), err)
		}

		n := new(big.Int).Mul(u128ToBig(a.coef),
			new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(38-a.scale)), nil))
		sq := new(big.Int).Sqrt(n)
		bound := new(big.Int).Mul(sq, new(big.Int).Add(sq, big.NewInt(1)))
		if n.Cmp(bound) > 0 {
			sq.Add(sq, big.NewInt(1))
		}

		want := newDec(false, u128FromBig(t, sq), maxScale)
		if !got.Equal(want) {
			t.Fatalf("Sqrt(%s) = %s, want %s", a.String(), got.String(), want.String())
		}
	}
}

func TestU128DecimalDigits(t *testing.T) {
	if got := u128Zero.decimalDigits(); got != 1 {
		t.Errorf("digits(0) = %d, want 1", got)
	}
	for i, p := range pow10 {
		if got := p.decimalDigits(); got != i+1 {
			t.Errorf("digits(10^%d) = %d, want %d", i, got, i+1)
		}
		if i > 0 {
			pm1, _ := p.Sub(u128One)
			if got := pm1.decimalDigits(); got != i {
				t.Errorf("digits(10^%d-1) = %d, want %d", i, got, i)
			}
		}
	}
	if got := (u128{hi: ^uint64(0), lo: ^uint64(0)}).decimalDigits(); got != 39 {
		t.Errorf("digits(2^128-1) = %d, want 39", got)
	}
}

func TestU256QuoRem64(t *testing.T) {
	// (2^64+5) * 10^19 / 10^19 must round-trip.
	x := u128{hi: 1, lo: 5}.MulFull(pow10[19])
	q, r := x.quoRem64(pow10[19].lo)
	if r != 0 || !q.carry.IsZero() || q.hi != 1 || q.lo != 5 {
		t.Fatalf("quoRem64 round-trip failed: q=%+v r=%d", q, r)
	}

	q, r = x.quoRem64(7)
	// cross-check with big.Int
	xb := new(big.Int).Lsh(big.NewInt(1), 64)
	xb.Add(xb, big.NewInt(5))
	xb.Mul(xb, new(big.Int).SetUint64(pow10[19].lo))
	qb, rb := new(big.Int).QuoRem(xb, big.NewInt(7), new(big.Int))

	gotQ := new(big.Int).Lsh(new(big.Int).SetUint64(q.carry.hi), 192)
	gotQ.Add(gotQ, new(big.Int).Lsh(new(big.Int).SetUint64(q.carry.lo), 128))
	gotQ.Add(gotQ, new(big.Int).Lsh(new(big.Int).SetUint64(q.hi), 64))
	gotQ.Add(gotQ, new(big.Int).SetUint64(q.lo))
	if gotQ.Cmp(qb) != 0 || rb.Uint64() != r {
		t.Fatalf("quoRem64 by 7 mismatch: got q=%s r=%d, want q=%s r=%s", gotQ, r, qb, rb)
	}
}
