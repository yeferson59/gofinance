package decimal

import (
	json "encoding/json/v2"
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

func BenchmarkDecimalSqrt(b *testing.B) {
	x := MustFromString("12345.6789")
	b.ReportAllocs()
	for b.Loop() {
		_, _ = x.Sqrt()
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

func BenchmarkDecimalMarshalJSON(b *testing.B) {
	xs := []Decimal{
		MustFromString("12345.6789"),
		MustFromString("-0.001"),
		MustFromString("987.654321"),
	}
	b.ReportAllocs()
	for b.Loop() {
		_, _ = json.Marshal(xs)
	}
}

func BenchmarkDecimalUnmarshalJSON(b *testing.B) {
	data := []byte(`[12345.6789,-0.001,987.654321]`)
	b.ReportAllocs()
	for b.Loop() {
		var xs []Decimal
		_ = json.Unmarshal(data, &xs)
	}
}

func BenchmarkDecimalMarshalText(b *testing.B) {
	x := MustFromString("12345.6789")
	buf := make([]byte, 0, 32)
	b.ReportAllocs()
	for b.Loop() {
		buf, _ = x.AppendText(buf[:0])
	}
	_ = buf
}

func BenchmarkDecimalMarshalBinary(b *testing.B) {
	x := MustFromString("12345.6789")
	buf := make([]byte, 0, 32)
	b.ReportAllocs()
	for b.Loop() {
		buf, _ = x.AppendBinary(buf[:0])
	}
	_ = buf
}

func BenchmarkDecimalUnmarshalBinary(b *testing.B) {
	data, err := MustFromString("12345.6789").MarshalBinary()
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	for b.Loop() {
		var x Decimal
		_ = x.UnmarshalBinary(data)
	}
}
