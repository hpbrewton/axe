package main

import (
	"os"
	"math/big"
	"encoding/json"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"image/color"
)

type TestResult struct {
	Name string `json:"name"` // name of the test
	Nsize int `json:"nsize"`// size of test 
	IndexNS int64 `json:"index"`// time to index in nano seconds 
	QueriesNS []int64 `json:"queries"` // nrepeat * query length
}

var colors = []color.Color{
	color.RGBA{R: 128, G: 0,   B: 0,   A: 0},
	color.RGBA{R: 0,   G: 128, B: 0,   A: 0},
	color.RGBA{R: 0,   G: 0,   B: 128, A: 0},
}

func main() {
	files := os.Args[1:]


	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	for i, file := range files {
		file, err := os.Open(file)
		if err != nil {
			panic(err)
		}

		xys := make(plotter.XYs, 0)
		decoder := json.NewDecoder(file)
		for {
			var test TestResult
			err := decoder.Decode(&test)
			if err != nil {
				break
			}
			sum := big.NewFloat(1)
			for _, queryTime := range test.QueriesNS {
				sum.Add(sum, big.NewFloat(float64(queryTime)))
			}
			sum.Mul(big.NewFloat(1.0/10), sum)
			f, _ := sum.Float64()
			if true { //f < 25000000 {
				xys = append(xys, plotter.XY{float64(test.Nsize), f})
			}
		}
		scatter, err := plotter.NewScatter(xys)
		scatter.GlyphStyle.Color = colors[i]
		p.Add(scatter)
	}
	p.Save(600, 600, "example.png")
}