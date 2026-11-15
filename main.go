package main

import (
	"fmt"

	"github.com/yeferson59/gofinance/v2/finance/annuities"
	"github.com/yeferson59/gofinance/v2/money"
)

func main() {
	p, _ := annuities.NewAnnuity().EffectiveAnnualRate(0.6666).Present(645671.25, money.COP).Periods(1).Monthly().MustBuild().PrincipalFuture()
	p2, _ := annuities.NewAnnuity().EffectiveAnnualRate(0.10).Present(645671.25, money.COP).Periods(1).Monthly().MustBuild().PrincipalFuture()

	fmt.Println(p.Sub(p2))
}
