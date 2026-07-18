module github.com/yeferson59/gofinance/examples

go 1.26.5

require (
	github.com/go-echarts/go-echarts/v2 v2.7.2
	github.com/yeferson59/gofinance v1.4.2
	github.com/yeferson59/gofinance/charts v0.0.0
)

replace (
	github.com/yeferson59/gofinance => ../
	github.com/yeferson59/gofinance/charts => ../charts
)
