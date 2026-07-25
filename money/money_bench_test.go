package money

import "testing"

func BenchmarkMoneyAdd(b *testing.B) {
	x := MustMoneyFromString("12345.68", USD)
	y := MustMoneyFromString("987.65", USD)

	b.ReportAllocs()
	for b.Loop() {
		_ = x.Add(y)
	}
}
