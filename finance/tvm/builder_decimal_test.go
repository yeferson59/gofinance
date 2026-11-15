package tvm

import (
	"errors"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/v2/decimal"
)

func TestDecimalSettersMatchFloatSetters(t *testing.T) {
	want, err := NewTVM().PV(300000).Rate(0.005).N(360).SolvePMT()
	require.NoError(t, err)

	got, err := NewTVM().
		PVDecimal(decimal.MustFromString("300000")).
		RateDecimal(decimal.MustFromString("0.005")).
		NDecimal(decimal.MustFromString("360")).
		SolvePMT()
	require.NoError(t, err)

	assert.True(t, want.Equal(got), "float %s, decimal %s", want, got)

	fv, err := NewTVM().PMTDecimal(decimal.MustFromString("-100")).RateDecimal(decimal.MustFromString("0.01")).NDecimal(decimal.MustFromString("12")).SolveFV()
	require.NoError(t, err)

	fvFloat, err := NewTVM().PMT(-100).Rate(0.01).N(12).SolveFV()
	require.NoError(t, err)

	assert.True(t, fv.Equal(fvFloat), "PMTDecimal %s, PMT %s", fv, fvFloat)

	pv, err := NewTVM().FVDecimal(decimal.MustFromString("1000")).RateDecimal(decimal.MustFromString("0.01")).NDecimal(decimal.MustFromString("12")).SolvePV()
	require.NoError(t, err)

	pvFloat, err := NewTVM().FV(1000).Rate(0.01).N(12).SolvePV()
	require.NoError(t, err)

	assert.True(t, pv.Equal(pvFloat), "FVDecimal %s, FV %s", pv, pvFloat)
}

// Setters given a float with no decimal value report it through every solver
// instead of panicking.
func TestInvalidFloatInputsReturnErrors(t *testing.T) {
	nan := math.NaN()

	configs := map[string]Config{
		"N":    NewTVM().N(nan).Rate(0.01).PV(1000),
		"Rate": NewTVM().N(12).Rate(math.Inf(1)).PV(1000),
		"PV":   NewTVM().N(12).Rate(0.01).PV(nan),
		"PMT":  NewTVM().N(12).Rate(0.01).PMT(nan),
		"FV":   NewTVM().N(12).Rate(0.01).FV(math.Inf(-1)),
	}

	for name, config := range configs {
		t.Run(name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				for solver, solve := range map[string]func() (decimal.Decimal, error){
					"SolveFV": config.SolveFV, "SolvePV": config.SolvePV, "SolvePMT": config.SolvePMT,
					"SolveN": config.SolveN, "SolveRate": config.SolveRate,
				} {
					_, err := solve()
					assert.True(t, errors.Is(err, decimal.ErrInvalidFormat), "%s error = %v", solver, err)
				}
			})
		})
	}
}
