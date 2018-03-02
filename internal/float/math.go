package float

import "math"

func Min(values ...float64) float64 {
	if len(values) == 0 {
		return 0
	}
	m := values[0]
	for i := 1; i < len(values); i++ {
		if m > values[i] {
			m = values[i]
		}
	}
	return m
}

func Max(values ...float64) float64 {
	if len(values) == 0 {
		return 0
	}
	m := values[0]
	for i := 1; i < len(values); i++ {
		if m < values[i] {
			m = values[i]
		}
	}
	return m
}

func Range(values ...float64) float64 {
	return Max(values...) - Min(values...)
}

func Sum(values ...float64) float64 {
	s := 0.0
	for _, x := range values {
		s += x
	}
	return s
}

func Mean(values ...float64) float64 {
	if len(values) == 0 {
		return 0
	}
	return Sum(values...) / float64(len(values))
}

func Variance(values ...float64) float64 {
	if len(values) == 0 {
		return 0
	}
	m := Mean(values...)
	v := 0.0
	for _, x := range values {
		v += (x - m) * (x - m)
	}
	return v
}

func Stdev(values ...float64) float64 {
	v := Variance(values...)
	return math.Sqrt(v / float64(len(values)))
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
