package evo

// Crosser creates a new child from the parents through crossover (or cloning if there is only one parent). The crosser is not responsible for mutation or for assigning the genome an ID or to a species.
type Crosser interface {
	Cross(parents ...Genome) (child Genome, err error)
}

// Evaluator utilises the network provided and returns its fitness (or error) as a result
type Evaluator interface {
	Evaluate(Phenome) (Result, error)
}

// Matrix descibes data organised as a matrix. It mimics a subset of the signature of [gonum's mat.Matrix](https://godoc.org/gonum.org/v1/gonum/mat) which allows directly passing matrices from that package as inputs as well as any other type that implements it, such as [sparse](https://godoc.org/github.com/james-bowman/sparse). Network implementations, however, may expect a specific type and throw an error if they cannot convert to the desired type.
type Matrix interface {

	// Dims returns the dimensions of a Matrix.
	Dims() (r, c int)

	// At returns the value of a matrix element at row i, column j.
	// It will panic if i or j are out of bounds for the matrix.
	At(i, j int) float64
}

// Mutator changes the genome's encoding (nodes, conns, or traits)
type Mutator interface {
	Mutate(*Genome) error
}

// Populator provides a popluation from which the experiment will begin
type Populator interface {
	Populate() (Population, error)
}

// Seeder provides an unitialised genome from which to construct a new population
type Seeder interface {
	Seed() (Genome, error)
}

// Searcher processes each phenome through the evaluator and returns the result
type Searcher interface {
	Search(Evaluator, []Phenome) ([]Result, error)
}

// Selector examines a population returns the current genomes who will continue and those that will become parents
type Selector interface {
	Select(Population) (continuing []Genome, parents [][]Genome, err error)
}

// Speciator assigns the population's genomes to a species, creating and destroying species as necessary.
type Speciator interface {
	Speciate(*Population) error
}

// Transcriber creates the decoded substrate from the encoded one.
type Transcriber interface {
	Transcribe(Substrate) (Substrate, error)
}

// Translator creates a new network from defintions contained in the nodes and connections
type Translator interface {
	Translate(Substrate) (Network, error)
}

// Mutators collection which acts as a single mutator. Component mutators will be called in order
// until the complexity of the genome changes.
type Mutators []Mutator

// Mutate the genome with the composite mutators
func (m Mutators) Mutate(g *Genome) error {

	// Record the starting complexity
	n := g.Complexity()

	// Iterate the mutators in order
	for _, x := range m {

		// Use the current mutator on the genome
		if err := x.Mutate(g); err != nil {
			return err
		}

		// The complexity has changed so do not continue with the remaining mutations
		if g.Complexity() != n {
			break
		}
	}
	return nil
}
