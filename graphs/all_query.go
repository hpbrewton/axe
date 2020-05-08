package main

import (
	"os"
	"encoding/json"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/plotter"
	"image/color"
	// "fmt"
)

type TestResult struct {
	Dir string
	Nfrags int 
	IndexNS int64 
	DataNS []int64
}

var cutoffs = []float64{0.000001, 1, 2, 3, 5, 10, 100, 1000}

func main() {
	file, err := os.Open("all2.dat")
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(file)
	var result TestResult
	decoder.Decode(&result)

	xys := make(plotter.XYs, len(cutoffs))
	for i, cutoff := range cutoffs{
		xys[i] = plotter.XY{cutoff, float64(result.DataNS[i])}
	}
	scatter, err := plotter.NewScatter(xys)
	scatter.GlyphStyle.Color = color.RGBA{G:255, A:255}
	scatter.GlyphStyle.Radius = 7.5
	scatter.GlyphStyle.Shape = &draw.PyramidGlyph{}
	p, err := plot.New()
	p.Y.Label.Text = "Time to do search (ns)"
	p.X.Label.Text = "Cuttoff of Search (log scale)"
	p.Add(scatter)
	p.Save(800, 800, "example.png")
}