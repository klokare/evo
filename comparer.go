package evo

// Comparer determines the relative order of two genomes
type Comparer interface {
	Compare(a, b Genome) int8
}

// ByFitness compares two genomes using a rubric that favours fitness
type ByFitness struct{}

// Compare two genomes for order based on fitness, from least fit to most fit.
func (c ByFitness) Compare(a, b Genome) int8 {
	switch {
	case a.Solved && !b.Solved:
		return 1
	case b.Solved && !a.Solved:
		return -1
	default:
		switch {
		case a.Fitness < b.Fitness:
			return -1
		case b.Fitness < a.Fitness:
			return 1
		default:
			switch {
			case a.Encoded.Complexity() > b.Encoded.Complexity():
				return -1
			case b.Encoded.Complexity() > a.Encoded.Complexity():
				return 1
			default:
				switch {
				case a.Decoded.Complexity() > b.Decoded.Complexity():
					return -1
				case b.Decoded.Complexity() > a.Decoded.Complexity():
					return 1
				default:
					switch {
					case a.ID > b.ID:
						return -1
					case b.ID > a.ID:
						return 1
					default:
						return 0
					}
				}
			}
		}
	}
}

// WithByFitness returns an option that set the experiment's comparer to a ByFitness helper.
func WithByFitness() Option {
	return func(e *Experiment) error {
		e.Comparer = ByFitness{}
		return nil
	}
}

// ByNovelty compares two genomes, favouring the more novel one
type ByNovelty struct{}

// Compare two genomes for order based on fitness, from least novel to most novel.
// TODO: think through considering more complexity as being related to higher potenial novelty
func (c ByNovelty) Compare(a, b Genome) int8 {
	switch {
	case a.Solved && !b.Solved:
		return 1
	case b.Solved && !a.Solved:
		return -1
	default:
		switch {
		case a.Novelty < b.Novelty:
			return -1
		case b.Novelty < a.Novelty:
			return 1
		default:
			switch {
			case a.Encoded.Complexity() > b.Encoded.Complexity():
				return 1
			case b.Encoded.Complexity() > a.Encoded.Complexity():
				return -1
			default:
				switch {
				case a.Decoded.Complexity() > b.Decoded.Complexity():
					return 1
				case b.Decoded.Complexity() > a.Decoded.Complexity():
					return -1
				default:
					switch {
					case a.ID < b.ID:
						return 1
					case b.ID < a.ID:
						return -1
					default:
						return 0
					}
				}
			}
		}
	}
}

// WithByNovelty returns an option that set the experiment's comparer to a ByNovelty helper.
func WithByNovelty() Option {
	return func(e *Experiment) error {
		e.Comparer = ByNovelty{}
		return nil
	}
}
