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
