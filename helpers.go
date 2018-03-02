package evo

// Configurer provides a consistent way to configure one or more helpers
type Configurer interface {
	Configure(...interface{}) error
}

// Crosser creates a new child from the parents through crossover (or cloning if there is only one parent). The crosser is not responsible for mutation or for assigning the genome an ID or to a species.
type Crosser interface {
	Cross(parents ...Genome) (child Genome, err error)
}

// Evaluator utilises the network provided and returns its fitness (or error) as a result
type Evaluator interface {
	Evaluate(Phenome) (Result, error)
}

// Mutator changes the genome's encoding (nodes, conns, or traits)
type Mutator interface {
	Mutate(*Genome) error
}

// Populator provides a popluation from which the experiment will begin
type Populator interface {
	Populate() (Population, error)
}

// Seeder provides the inital seed genome(s) that create the experiment's first evaluable population
type Seeder interface {
	Seed() ([]Genome, error)
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

// Translator creates a new network from defintions contained in the nodes and connections
type Translator interface {
	Translate(Substrate) (Network, error)
}

// Transcriber decodes the genome, returning the nodes and connections to be used to create the network
type Transcriber interface {
	Transcribe(Substrate) (Substrate, error)
}
