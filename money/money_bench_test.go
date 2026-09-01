package money

import (
	json "encoding/json/v2"
	"testing"
)

func BenchmarkMoneyAdd(b *testing.B) {
	x := MustMoneyFromString("12345.68", USD)
	y := MustMoneyFromString("987.65", USD)

	b.ReportAllocs()
	for b.Loop() {
		_ = x.Add(y)
	}
}

func BenchmarkMoneyMarshalJSON(b *testing.B) {
	m := MustMoneyFromString("12345.67", USD)

	b.ReportAllocs()
	for b.Loop() {
		_, _ = json.Marshal(m)
	}
}

func BenchmarkMoneyUnmarshalJSON(b *testing.B) {
	data := []byte(`{"value":"12345.67","currency":"USD"}`)

	b.ReportAllocs()
	for b.Loop() {
		var m Money
		_ = json.Unmarshal(data, &m)
	}
}

func BenchmarkCurrencyMarshalJSON(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = json.Marshal(USD)
	}
}

func BenchmarkCurrencyUnmarshalJSON(b *testing.B) {
	data := []byte(`"USD"`)

	b.ReportAllocs()
	for b.Loop() {
		var c Currency
		_ = c.UnmarshalJSON(data)
	}
}

func BenchmarkMoneyMarshalBinary(b *testing.B) {
	m := MustMoneyFromString("12345.67", USD)
	buf := make([]byte, 0, 32)

	b.ReportAllocs()
	for b.Loop() {
		buf, _ = m.AppendBinary(buf[:0])
	}
	_ = buf
}

func BenchmarkMoneyUnmarshalBinary(b *testing.B) {
	data, err := MustMoneyFromString("12345.67", USD).MarshalBinary()
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	for b.Loop() {
		var m Money
		_ = m.UnmarshalBinary(data)
	}
}
