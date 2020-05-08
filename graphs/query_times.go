package main

import (
	"os"
	"encoding/json"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
	"gonum.org/v1/plot/plotter"
	"image/color"
	"strings"
	"sort"
)

type GroupAndBarSorter struct {
	groups []string 
	bars [][]float64
	keyBar int
}

func (gbs *GroupAndBarSorter) Len() int {
	return len(gbs.bars[gbs.keyBar])
}

func (gbs *GroupAndBarSorter) Less(i, j int) bool {
	return gbs.bars[gbs.keyBar][i] < gbs.bars[gbs.keyBar][j]
}

func (gbs *GroupAndBarSorter) Swap(i, j int) {
	gbs.groups[i], gbs.groups[j] = gbs.groups[j], gbs.groups[i]
	for _, bar := range gbs.bars {
		bar[i], bar[j] = bar[j], bar[i]
	}
}

type TestResult struct {
	Dir string
	Nfrags int 
	IndexNS int64 
	DataNS []int64
}

var Green = color.RGBA{G: 255, A: 255}
var Red = color.RGBA{R: 255, A: 255}
var colors = []color.Color{
	color.RGBA{R: 128, G: 0,   B: 0,   A: 0},
	color.RGBA{R: 0,   G: 128, B: 0,   A: 0},
	color.RGBA{R: 0,   G: 0,   B: 128, A: 0},
}

var intervals = []float64{0, 1, 2, 3, 5, 10, 100, 1000}

func min(a, b int) int {
	if a < b {
		return a 
	} else {
		return b
	}
}

func toFloat64Slice(l []int64) []float64 {
	ret := make([]float64, len(l))
	for i, v := range l {
		ret[i] = float64(v)
	}
	return ret
}

func main() {
	files := os.Args[1:]

	file, err := os.Open(files[0])
	if err != nil {
		panic(err)
	}
	fileReader := json.NewDecoder(file)
	groups := make([]string, 0)
	bars := make([][]float64, len(intervals))
	for i, _ := range bars {
		bars[i] = make([]float64, 0)
	}
	for {
		var test TestResult
		err := fileReader.Decode(&test)
		if err != nil {
			break
		}
		groups = append(groups, test.Dir)
		for i, v := range test.DataNS {
			bars[i] = append(bars[i], float64(v)/float64(test.DataNS[len(intervals)-1]))
		}
	}

	// clean up groups
	for i, group := range groups {
		items := strings.Split(group, "/")
		items = items[3:]
		groups[i] = strings.Join(items, "/")
		if len(groups[i]) > 12 {
			groups[i] = groups[i][:12]
		}
	}

	// sourt by last interval
	keyBar := 0
	sorter := GroupAndBarSorter{groups: groups, bars: bars, keyBar: keyBar}
	sort.Sort(&sorter)


	across := 13
	width := 1
	height := 10

	plots := make([][]*plot.Plot, height)
	for row := 0; row < height; row++ {
		plots[row] = make([]*plot.Plot, width)
		for col := 0; col < width; col++ {
			ubound := min(row*across+across, len(bars[keyBar]))
			values := plotter.Values(bars[keyBar][row*across:ubound])
			redValues := make(plotter.Values, len(values))
			for i, value := range values {
				if value > 1 {
					redValues[i] = value 
				} else {
					redValues[i] = 0
				}
			}
			bar, err := plotter.NewBarChart(values, 0.5*vg.Centimeter)
			redbar, err := plotter.NewBarChart(redValues, 0.5*vg.Centimeter)
			bar.Color = Green
			redbar.Color = Red
			if err != nil {
				panic(err)
			}

			p, err := plot.New()
			if err != nil {
				panic(err)
			}
			p.Add(bar, redbar)
			p.Y.Max = 2
			p.Y.Min = 0
			p.NominalX(groups[row*across:ubound]...)
			plots[row][col] = p
		}
	}


	img := vgimg.New(vg.Inch*10, vg.Inch*11)
	canvas := draw.New(img)

	tiling := draw.Tiles{
		Rows: height,
		Cols: width,
		PadX:      vg.Millimeter,
		PadY:      vg.Millimeter,
		PadTop:    vg.Points(2),
		PadBottom: vg.Points(2),
		PadLeft:   vg.Points(2),
		PadRight:  vg.Points(2),
	}

	canvases := plot.Align(plots, tiling, canvas)
	for i, plotRow := range plots {
		for j, plot := range plotRow {
			if plot != nil {
				plot.Draw(canvases[i][j])
			}
		}
	}

	out, err := os.Create("example0.png")
	if err != nil {
		panic(err)
	}
	defer out.Close()
	png := vgimg.PngCanvas{Canvas: img}
	if _, err := png.WriteTo(out); err != nil {
		panic(err)
	}
}