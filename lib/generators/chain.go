package generators

type ChainGenerator struct {
	GeneratorFuncs []GeneratorFunction
}

func (g *ChainGenerator) ChainFunc(stat float64, delta float64) float64 {
	for _, funct := range g.GeneratorFuncs {
		stat = funct(stat, delta)
	}

	return stat
}
