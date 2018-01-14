package evo

import "context"

// Configurer provides a consistent way to configure one or more helpers
type Configurer interface {
	Configure(...interface{}) error
}

// Crosser creates a new child from the parents through crossover (or cloning if there is only one parent). The crosser is not responsible for mutation or for assigning the genome an ID or to a species.
type Crosser interface {
	Cross(ctx context.Context, parents ...Genome) (child Genome, err error)
}

// Evaluator utilises the network provided and returns its fitness (or error) as a result
type Evaluator interface {
	Evaluate(context.Context, Phenome) (Result, error)
}

// Mutator changes the genome's encoding (nodes, conns, or traits)
type Mutator interface {
	Mutate(context.Context, *Genome) error
}

// Populator provides a popluation from which the experiment will begin
type Populator interface {
	Populate(context.Context) (Population, error)
}

// Searcher processes each phenome through the evaluator and returns the result
type Searcher interface {
	Search(context.Context, Evaluator, []Phenome) ([]Result, error)
}

// Selector examines a population returns the current genomes who will continue and those that will become parents
type Selector interface {
	Select(context.Context, Population) (continuing []Genome, parents [][]Genome, err error)
}

// Speciator assigns the population's genomes to a species, creating and destroying species as necessary.
type Speciator interface {
	Speciate(context.Context, *Population) error
}

// Translator creates a new network from defintions contained in the nodes and connections
type Translator interface {
	Translate(context.Context, Substrate) (Network, error)
}

// Transcriber decodes the genome, returning the nodes and connections to be used to create the network
type Transcriber interface {
	Transcribe(context.Context, Substrate) (Substrate, error)
}
