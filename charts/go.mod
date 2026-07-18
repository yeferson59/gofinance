module github.com/yeferson59/gofinance/charts

go 1.26.5

require (
	github.com/go-echarts/go-echarts/v2 v2.7.2
	github.com/yeferson59/gofinance/v2 v2.0.0
)

// During development inside this repository the parent module is used
// directly; consumers outside the repo resolve the require above.
replace github.com/yeferson59/gofinance/v2 => ../
