package float

import (
	"math"
	"sort"
)

// Min returns the minimum value in the slice
func Min(values []float64) float64 {
	n := len(values)
	switch n {
	case 0:
		return 0
	case 1:
		return values[0]
	default:
		x := values[0]
		for i := 1; i < n; i++ {
			if x > values[i] {
				x = values[i]
			}
		}
		return x
	}
}

// Max returns the maximum value in the slice
func Max(values []float64) float64 {
	n := len(values)
	switch n {
	case 0:
		return 0
	case 1:
		return values[0]
	default:
		x := values[0]
		for i := 1; i < n; i++ {
			if x < values[i] {
				x = values[i]
			}
		}
		return x
	}
}

func Range(values []float64) float64 {
	return Max(values) - Min(values)
}

func Sum(values []float64) float64 {
	s := 0.0
	for _, x := range values {
		s += x
	}
	return s
}

func Mean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	return Sum(values) / float64(len(values))
}

func Variance(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	m := Mean(values)
	v := 0.0
	for _, x := range values {
		v += (x - m) * (x - m)
	}
	return v
}

func Stdev(values []float64) float64 {
	v := Variance(values)
	return math.Sqrt(v / float64(len(values)))
}

func Median(values []float64) float64 {
	n := len(values)
	switch n {
	case 0:
		return 0
	case 1:
		return values[1]
	default:
		v2 := make([]float64, len(values)) // make a copy so we do not alter order of original slice
		copy(v2, values)
		sort.Float64s(v2)
		i := n / 2
		if n%2 == 0 {
			return (v2[i] + v2[i+1]) / 2.0
		}
		return v2[i]
	}
}

func Abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func Round(f float64, s int) float64 {
	shift := math.Pow(10, float64(s))
	return math.Floor(f*shift+.5) / shift
}
