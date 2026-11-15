package gradients

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yeferson59/gofinance/v2/money"
)

// Builder inputs that have no decimal value are reported through Build's
// error instead of panicking inside the setter.
func TestInvalidFloatInputsReturnErrors(t *testing.T) {
	nan := math.NaN()

	arithmetic := map[string]ArithmeticConfig{
		"first payment": NewArithmeticSeries().FirstPayment(nan, money.USD).Gradient(10, money.USD).Rate(0.01).Periods(10),
		"gradient":      NewArithmeticSeries().FirstPayment(100, money.USD).Gradient(nan, money.USD).Rate(0.01).Periods(10),
		"rate":          NewArithmeticSeries().FirstPayment(100, money.USD).Gradient(10, money.USD).Rate(nan).Periods(10),
		"annual rate":   NewArithmeticSeries().FirstPayment(100, money.USD).Gradient(10, money.USD).AnnualRate(math.Inf(1)).Periods(10),
	}

	for name, config := range arithmetic {
		t.Run("arithmetic "+name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				_, err := config.Build()
				assert.Error(t, err)
			})
		})
	}

	geometric := map[string]GeometricConfig{
		"first payment": NewGeometricSeries().FirstPayment(nan, money.USD).GrowthRate(0.05).Rate(0.01).Periods(10),
		"growth rate":   NewGeometricSeries().FirstPayment(100, money.USD).GrowthRate(nan).Rate(0.01).Periods(10),
		"rate":          NewGeometricSeries().FirstPayment(100, money.USD).GrowthRate(0.05).Rate(nan).Periods(10),
		"annual rate":   NewGeometricSeries().FirstPayment(100, money.USD).GrowthRate(0.05).AnnualRate(math.Inf(-1)).Periods(10),
	}

	for name, config := range geometric {
		t.Run("geometric "+name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				_, err := config.Build()
				assert.Error(t, err)
			})
		})
	}
}
