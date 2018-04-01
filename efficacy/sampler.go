package efficacy

import (
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/klokare/evo"
	"github.com/klokare/evo/internal/float"
)

// Field is the type of field being collected
type Field byte

// Known fields
const (
	BatchNumber Field = iota
	Began
	Ended
	Solved
	Generations
	Evaluations
	Seconds
	Fitness
	Novelty
	Encoded
	EncodedNodes
	EncodedConns
	Decoded
	DecodedNodes
	DecodedConns
)

// String returns a text description of the field
func (f Field) String() string {
	switch f {
	case BatchNumber:
		return "batch number"
	case Began:
		return "time began"
	case Ended:
		return "time ended"
	case Solved:
		return "solved"
	case Generations:
		return "generations"
	case Evaluations:
		return "evaluations"
	case Seconds:
		return "seconds"
	case Fitness:
		return "fitness"
	case Novelty:
		return "novelty"
	case Encoded:
		return "encoded complexity"
	case EncodedNodes:
		return "encoded nodes"
	case EncodedConns:
		return "encoded conns"
	case Decoded:
		return "decoded complexity"
	case DecodedNodes:
		return "decoded nodes"
	case DecodedConns:
		return "decoded conns"
	default:
		return "unknown field"
	}
}

// Method is the aggregrate function
type Method byte

// Known methods
const (
	Min Method = iota
	Max
	Mean
	Median
	Best // values from the most fit genome
)

// String provides the text description of the method
func (m Method) String() string {
	switch m {
	case Min:
		return "min"
	case Max:
		return "max"
	case Mean:
		return "mean"
	case Median:
		return "Median"
	case Best:
		return "best"
	default:
		return "unknown method"
	}
}

// Sample is the result of an evaluation run
type Sample struct {
	BatchNumber int
	Began       time.Time
	Ended       time.Time
	Solved      bool
	Generations int
	Evaluations int
	Seconds     float64
	Values      map[Field]map[Method]float64
}

// Sampler records the results of an experiment's run
type Sampler struct {
	w   io.WriteCloser
	enc *json.Encoder
}

// NewSampler creates a new efficacy sampler
func NewSampler(filename string) (*Sampler, error) {
	var err error
	s := new(Sampler)
	if s.w, err = os.Create(filename); err != nil {
		return nil, err
	}
	s.enc = json.NewEncoder(s.w)
	return s, nil
}

// Close the underlying writer returning any error from that action
func (s *Sampler) Close() error { return s.w.Close() }

// Record appends the sample to the file.
func (s *Sampler) Record(sample Sample) error {
	return s.enc.Encode(sample)
}

// Callbacks returns a new listener callback function for the given batch number.
func (s *Sampler) Callbacks(num int) (start, complete evo.Callback) {

	// Create a new sample
	sample := &Sample{
		BatchNumber: num,
		Began:       time.Now(),
		Values:      make(map[Field]map[Method]float64, 10),
	}

	// Callback which starts the timer
	start = func(evo.Population) error {
		sample.Began = time.Now()
		return nil
	}

	// Callback which stops the timer and writes the record
	complete = func(pop evo.Population) error {

		// Complete the sample
		sample.Ended = time.Now()
		sample.Seconds = float64(sample.Ended.Sub(sample.Began)) / float64(time.Second)
		sample.Generations = pop.Generation
		sample.Evaluations = pop.Generation * len(pop.Genomes) // TODO: should this be changed to an OnIterations listener that sums the number of genomes per iteration?

		// Identify the best genome
		genomes := make([]evo.Genome, len(pop.Genomes)) // work with a copy so as not to affect other callbacks
		copy(genomes, pop.Genomes)
		evo.SortBy(genomes, evo.BySolved, evo.ByFitness, evo.ByComplexity, evo.ByAge)
		best := genomes[len(genomes)-1]
		sample.Solved = best.Solved

		// Create the data
		n := len(pop.Genomes)
		fit := make([]float64, n)
		nov := make([]float64, n)
		enc := make([]float64, n)
		encn := make([]float64, n)
		encc := make([]float64, n)
		dec := make([]float64, n)
		decn := make([]float64, n)
		decc := make([]float64, n)

		for i, g := range genomes {
			fit[i] = g.Fitness
			nov[i] = g.Novelty
			encn[i] = float64(len(g.Encoded.Nodes))
			encc[i] = float64(len(g.Encoded.Conns))
			enc[i] = encn[i] + encc[i]
			decn[i] = float64(len(g.Decoded.Nodes))
			decc[i] = float64(len(g.Decoded.Conns))
			dec[i] = decn[i] + decc[i]
		}

		// Store the data
		sample.Values[Fitness] = map[Method]float64{
			Min:    float.Min(fit),
			Max:    float.Max(fit),
			Mean:   float.Mean(fit),
			Median: float.Median(fit),
			Best:   best.Fitness,
		}

		sample.Values[Novelty] = map[Method]float64{
			Min:    float.Min(nov),
			Max:    float.Max(nov),
			Mean:   float.Mean(nov),
			Median: float.Median(nov),
			Best:   best.Novelty,
		}

		sample.Values[Encoded] = map[Method]float64{
			Min:    float.Min(enc),
			Max:    float.Max(enc),
			Mean:   float.Mean(enc),
			Median: float.Median(enc),
			Best:   float64(best.Encoded.Complexity()),
		}

		sample.Values[EncodedNodes] = map[Method]float64{
			Min:    float.Min(encn),
			Max:    float.Max(encn),
			Mean:   float.Mean(encn),
			Median: float.Median(encn),
			Best:   float64(len(best.Encoded.Nodes)),
		}

		sample.Values[EncodedConns] = map[Method]float64{
			Min:    float.Min(encc),
			Max:    float.Max(encc),
			Mean:   float.Mean(encc),
			Median: float.Median(encc),
			Best:   float64(len(best.Encoded.Conns)),
		}

		sample.Values[Decoded] = map[Method]float64{
			Min:    float.Min(dec),
			Max:    float.Max(dec),
			Mean:   float.Mean(dec),
			Median: float.Median(dec),
			Best:   float64(best.Decoded.Complexity()),
		}

		sample.Values[DecodedNodes] = map[Method]float64{
			Min:    float.Min(decn),
			Max:    float.Max(decn),
			Mean:   float.Mean(decn),
			Median: float.Median(decn),
			Best:   float64(len(best.Decoded.Nodes)),
		}

		sample.Values[DecodedConns] = map[Method]float64{
			Min:    float.Min(decc),
			Max:    float.Max(decc),
			Mean:   float.Mean(decc),
			Median: float.Median(decc),
			Best:   float64(len(best.Decoded.Conns)),
		}

		// Record the sample
		return s.Record(*sample)
	}
	return
}
