package evo

// A Genome is the encoded neural network and its last result when applied in evaluation.
// For performance reasons, helpers should keep the nodes (by ID) and conns (by source and then target IDs) sorted though this is not required.
type Genome struct {
	ID      int64     // The genome's unique identifier
	Species int       // The ID of the species
	Age     int       // Number of generations genome has been alive
	Fitness float64   // The genome's latest fitness score
	Novelty float64   // The genome's latest novelty score, if any
	Solved  bool      // True if the genome produced a solution in the last evaluation
	Traits  []float64 // Additional information, encoded as floats, that will be passed to the evaluation function
	Encoded Substrate // The encoded neural network layout
	Decoded Substrate // The decoded neural network layout
}

// Complexity returns the number of nodes and connections in the genome
func (g Genome) Complexity() int { return g.Encoded.Complexity() }

// A Phenome is the instatiated genome to be sent to the evaluator
type Phenome struct {
	ID      int64     // The unique identifier of the genome
	Traits  []float64 // Any additional information, specific to this genome, to be passed to the evaluation function. This is optional.
	Network           // The neural network made from the genome's encoding
}

// Result describes the outcome of running the evaluation. ID and fitness are required properties. If an error occurs in the evaluation, this should be returned with the result.
type Result struct {
	ID       int64       // The unique ID of the genome from which the phenome was made
	Solved   bool        // True if the network provided a winning solution
	Fitness  float64     // A positive value indicating the fitness of this network after evaluation
	Novelty  float64     // An optional value indicating the novelty of this network's decisions during evaluation
	Behavior interface{} // An optional slice describing the novelty of the network's decisions
}
