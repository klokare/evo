package watcher

import (
	"fmt"
	"sort"

	"github.com/klokare/evo"
)

type Console struct{}

func (h *Console) Watch(p evo.Population) error {
	sort.Sort(sort.Reverse(p.Genomes))
	b := p.Genomes[0]
	fmt.Printf("Generation %d has %d species. Best is %d with complexity %d and fitness %f\n", p.Generation, len(p.Species), b.ID, b.Complexity(), b.Fitness)
	return nil
}