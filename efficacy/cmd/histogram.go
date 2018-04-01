package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image/color"
	"io"
	"log"
	"math"
	"os"
	"strings"

	"github.com/klokare/evo/efficacy"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func main() {

	// Create and parse the flags
	var (
		series1 = flag.String("series1", "", "name to use for series 1")
		series2 = flag.String("series2", "", "name to use for series 2")
		field   = flag.String("field", "", "data field to report")
		method  = flag.String("method", "mean", "method used for aggregation")
		solved  = flag.Bool("solved", false, "include only solved cases")
		output  = flag.String("output", "histogram.png", "filename for output chart")
		height  = flag.Float64("height", 9, "chart height in centimeters")
		width   = flag.Float64("width", 15, "chart weight in centimeters")
	)
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("you must specify at least one result file")
	} else if len(args) > 2 {
		log.Fatal("you may only specify one or two result files")
	}

	// Validate the field and method
	var err error
	var f efficacy.Field
	if f, err = validateField(*field); err != nil {
		log.Fatal(err)
	}

	var m efficacy.Method
	if m, err = validateMethod(*method); err != nil {
		log.Fatal(err)
	}

	// Load the data
	var s1name, s2name string
	var s1data, s2data []float64
	if s1name, s1data, err = makeSeries(*series1, args[0], f, m, *solved); err != nil {
		log.Fatal(err)
	}
	if len(args) > 1 {
		if s2name, s2data, err = makeSeries(*series2, args[1], f, m, *solved); err != nil {
			log.Fatal(err)
		}
	}

	// Make the chart
	if err := makeChart(*output, xlabel(f, m), s1name, s2name, s1data, s2data, *solved, vg.Length(*height), vg.Length(*width)); err != nil {
		log.Fatal(err)
	}
}

func validateField(name string) (efficacy.Field, error) {
	switch strings.ToLower(name) {
	case "generations":
		return efficacy.Generations, nil
	case "evaluations":
		return efficacy.Evaluations, nil
	case "seconds":
		return efficacy.Seconds, nil
	case "fitness":
		return efficacy.Fitness, nil
	case "novelty":
		return efficacy.Novelty, nil
	case "encoded":
		return efficacy.Encoded, nil
	case "encoded-nodes":
		return efficacy.EncodedNodes, nil
	case "encoded-conns":
		return efficacy.EncodedConns, nil
	case "decoded":
		return efficacy.Decoded, nil
	case "decoded-nodes":
		return efficacy.DecodedNodes, nil
	case "decoded-conns":
		return efficacy.DecodedConns, nil
	default:
		return 0, fmt.Errorf("unknown field name %s", name)
	}
}

func validateMethod(name string) (efficacy.Method, error) {
	switch strings.ToLower(name) {
	case "min":
		return efficacy.Min, nil
	case "max":
		return efficacy.Max, nil
	case "mean":
		return efficacy.Mean, nil
	case "median":
		return efficacy.Median, nil
	case "best":
		return efficacy.Best, nil
	default:
		return 0, fmt.Errorf("unknown method name %s", name)
	}
}

func xlabel(f efficacy.Field, m efficacy.Method) string {
	switch f {
	case efficacy.Generations, efficacy.Evaluations, efficacy.Seconds:
		return f.String()
	case efficacy.Fitness, efficacy.Novelty,
		efficacy.Encoded, efficacy.EncodedNodes, efficacy.EncodedConns,
		efficacy.Decoded, efficacy.DecodedNodes, efficacy.DecodedConns:
		return m.String() + " " + f.String()
	default:
		return "unknown x axis label"
	}
}

func makeSeries(name, filename string, fld efficacy.Field, met efficacy.Method, solved bool) (string, []float64, error) {

	// Load the data from the file
	var err error
	var f *os.File
	if f, err = os.Open(filename); err != nil {
		return "", nil, err
	}
	defer f.Close()

	dec := json.NewDecoder(f)

	var s efficacy.Sample
	data := make([]float64, 0, 100)
	cnt := 0
	for err == nil {
		if err = dec.Decode(&s); err != nil {
			break
		}
		switch fld {
		case efficacy.Generations:
			data = append(data, float64(s.Generations))
		case efficacy.Evaluations:
			data = append(data, float64(s.Evaluations))
		case efficacy.Seconds:
			data = append(data, s.Seconds)
		default:
			data = append(data, s.Values[fld][met])
		}
		if solved && !s.Solved {
			cnt++
		}
	}
	if err != io.EOF {
		return "", nil, err
	}

	// Set the name
	if name == "" {
		var info os.FileInfo
		if info, err = os.Stat(filename); err != nil {
			return "", nil, err
		}
		name = strings.Replace(info.Name(), ".json", "", 1)
	}
	if cnt == 1 {
		name = fmt.Sprintf("%s (1 failure)", name)
	} else if cnt > 1 {
		name = fmt.Sprintf("%s (%d failures)", name, cnt)
	}
	return name, data, nil
}

func makeChart(filename, xname string, s1name, s2name string, s1data, s2data []float64, solved bool, h, w vg.Length) error {

	// Note if using two series
	two := len(s2data) > 0

	// Create the plot
	p, err := plot.New()
	if err != nil {
		return err
	}
	p.Title.Text = "Histogram"
	p.X.Label.Text = xname
	p.Y.Label.Text = "density"

	// Determine the bins
	cnt, xmin, xmax := countAndRange(s1data, s2data)

	// Write series 1
	var h1, h2 *plotter.Histogram
	h1 = &plotter.Histogram{
		FillColor: color.NRGBA{R: 0, G: 0, B: 255, A: 128},
		LineStyle: draw.LineStyle{Color: color.Black, Width: vg.Points(1), Dashes: []vg.Length{}, DashOffs: 0},
	}
	h1.Bins, h1.Width = binsAndWidth(cnt, xmin, xmax, s1data)
	h1.Normalize(1)
	p.Add(h1)

	// Write series 2 if present
	if two {
		h2 = &plotter.Histogram{
			FillColor: color.NRGBA{R: 255, G: 0, B: 0, A: 128},
			LineStyle: draw.LineStyle{Color: color.Black, Width: vg.Points(1), Dashes: []vg.Length{}, DashOffs: 0},
		}
		h2.Bins, h2.Width = binsAndWidth(cnt, xmin, xmax, s2data)
		h2.Normalize(1)
		p.Add(h2)
	}

	// Add legend
	p.Legend.Add(s1name, h1)
	if two {
		p.Legend.Add(s2name, h2)
	}
	p.Legend.Top = true
	p.Legend.TextStyle.Font.Size = vg.Points(10)

	// Save the file
	return p.Save(w*vg.Centimeter, h*vg.Centimeter, filename)
}

func countAndRange(sets ...[]float64) (cnt int, min, max float64) {

	// Determine the actual count and min, max of the combined sets
	min = math.MaxFloat64
	max = -min
	for _, set := range sets {
		if len(set) == 0 {
			continue
		}
		if len(set) > cnt {
			cnt = len(set)
		}
		for _, x := range set {
			if min > x {
				min = x
			}
			if max < x {
				max = x
			}
		}
	}
	return
}

// TODO: incorporate a bin for solved. Can NaN be it's range?
func binsAndWidth(cnt int, xmin, xmax float64, data []float64) (bins []plotter.HistogramBin, w float64) {

	// Calculate the number of bins
	n := int(math.Ceil(math.Sqrt(float64(cnt))))
	bins = make([]plotter.HistogramBin, n)

	// Calculate the width and set the range for each bin
	w = (xmax - xmin) / float64(n)
	for i := range bins {
		bins[i].Min = xmin + float64(i)*w
		bins[i].Max = xmin + float64(i+1)*w
	}

	// Fill the bins
	for _, x := range data {
		bin := int((x - xmin) / w)
		if x == xmax {
			bin = n - 1
		}
		if bin < 0 || bin >= n {
			panic(fmt.Sprintf("%g, xmin=%g, xmax=%g, w=%g, bin=%d, n=%d\n",
				x, xmin, xmax, w, bin, n))
		}
		bins[bin].Weight += 1.0
	}
	return
}
