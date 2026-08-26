module github.com/yeferson59/gofinance/examples

go 1.27.0

require (
	github.com/go-echarts/go-echarts/v2 v2.7.2
	github.com/yeferson59/gofinance/charts v0.0.0
	github.com/yeferson59/gofinance/v2 v2.0.0
)

replace (
	github.com/yeferson59/gofinance/charts => ../charts
	github.com/yeferson59/gofinance/v2 => ../
)
