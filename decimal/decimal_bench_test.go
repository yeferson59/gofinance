package decimal

import (
	"testing"
)

func BenchmarkDecimalAdd(b *testing.B) {
	x := MustFromString("12345.6789")
	y := MustFromString("987.654321")
	b.ReportAllocs()
	for b.Loop() {
		_ = x.Add(y)
	}
}

func BenchmarkDecimalMul(b *testing.B) {
	x := MustFromString("12345.6789")
	y := MustFromString("987.654321")
	b.ReportAllocs()
	for b.Loop() {
		_ = x.Mul(y)
	}
}

func BenchmarkDecimalDiv(b *testing.B) {
	x := MustFromString("12345.6789")
	y := MustFromString("987.654321")
	b.ReportAllocs()
	for b.Loop() {
		_ = x.MustDiv(y)
	}
}

func BenchmarkDecimalCmp(b *testing.B) {
	x := MustFromString("12345.6789")
	y := MustFromString("987.654321")
	b.ReportAllocs()
	for b.Loop() {
		_ = x.Cmp(y)
	}
}

func BenchmarkDecimalParse(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = MustFromString("12345.6789")
	}
}

func BenchmarkDecimalString(b *testing.B) {
	x := MustFromString("12345.6789")
	b.ReportAllocs()
	for b.Loop() {
		_ = x.String()
	}
}

func BenchmarkDecimalLn(b *testing.B) {
	x := MustFromString("12345.6789")
	b.ReportAllocs()
	for b.Loop() {
		_, _ = x.Ln()
	}
}

func BenchmarkDecimalLog10(b *testing.B) {
	x := MustFromString("12345.6789")
	b.ReportAllocs()
	for b.Loop() {
		_, _ = x.Log10()
	}
}

func BenchmarkDecimalPowInt(b *testing.B) {
	x := MustFromString("1.005")
	n := MustFromString("360")
	b.ReportAllocs()
	for b.Loop() {
		_, _ = x.Pow(n)
	}
}

func BenchmarkDecimalPowFrac(b *testing.B) {
	x := MustFromString("1.05")
	n := MustFromString("2.5")
	b.ReportAllocs()
	for b.Loop() {
		_, _ = x.Pow(n)
	}
}
