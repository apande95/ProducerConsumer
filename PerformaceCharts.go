package main

import (
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

var (
	bar3DRangeColor = []string{
		"#313695", "#4575b4", "#74add1", "#abd9e9", "#e0f3f8",
		"#fee090", "#fdae61", "#f46d43", "#d73027", "#a50026",
	}

	xAxis = []int{}

	yAxis = []int{}
)

type pc struct {
	p int
	c int
	t float64
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func initVals() []pc {
	dat, err := os.ReadFile("performace.txt")
	check(err)
	datStr := string(dat)
	dataArr := strings.Split(datStr, "\n")
	var bar3Dpoints []pc
	var x_max int = 0
	var y_max int = 0
	for _, i := range dataArr {
		if len(i) > 0 {
			a := strings.Split(i, "#")
			j := new(pc)
			j.p, err = strconv.Atoi(a[0])
			check(err)
			x_max = Max(j.p, x_max)
			j.p -= 1
			j.c, err = strconv.Atoi(a[1])
			y_max = Max(j.c, y_max)
			j.c -= 1
			j.t, err = strconv.ParseFloat(strings.Replace(a[2], "s", "", -1), 64)
			check(err)
			bar3Dpoints = append(bar3Dpoints, *j)
		}
	}
	for i := 1; i <= x_max; i++ {
		xAxis = append(xAxis, i)
	}

	for i := 1; i <= y_max; i++ {
		yAxis = append(yAxis, i)
	}
	return bar3Dpoints
}

func genBar3dData(bar3Dpoints []pc) []opts.Chart3DData {
	ret := make([]opts.Chart3DData, 0)
	for _, d := range bar3Dpoints {
		ret = append(ret, opts.Chart3DData{
			Value: []interface{}{d.p, d.c, d.t},
		})
	}

	return ret
}

func bar3DBase() *charts.Bar3D {
	bar3d := charts.NewBar3D()
	bar3d.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Performance Evaluation"}),
		charts.WithVisualMapOpts(opts.VisualMap{
			Calculable: true,
			Max:        30,
			Range:      []float32{0, 30},
			InRange:    &opts.VisualMapInRange{Color: bar3DRangeColor},
		}),
		charts.WithGrid3DOpts(opts.Grid3D{
			BoxWidth: 200,
			BoxDepth: 80,
		}),
	)

	bar3Dpoints := initVals()

	bar3d.SetGlobalOptions(
		charts.WithXAxis3DOpts(opts.XAxis3D{Data: xAxis}),
		charts.WithYAxis3DOpts(opts.YAxis3D{Data: yAxis}),
	)
	bar3d.AddSeries("bar3d", genBar3dData(bar3Dpoints))
	return bar3d
}

type Bar3dExamples struct{}

func (Bar3dExamples) Examples() {
	page := components.NewPage()
	page.AddCharts(
		bar3DBase(),
	)

	f, err := os.Create("bar3d_1.html")
	check(err)
	page.Render(io.MultiWriter(f))
}

func main() {
	Bar3dExamples{}.Examples()
}
