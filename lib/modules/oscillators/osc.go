package oscillators

import (
	"github.com/malparty/synth-xu/lib/modules"
)

type OscType string

const (
	Saw      OscType = "saw"
	Sin      OscType = "sin"
	Square   OscType = "square"
	Triangle OscType = "triangle"
)

type Osc struct {
	Type OscType
}

func (o *Osc) GetModuleFunc() modules.ModuleFunction {
	return func(stat, delta float64) float64 {
		switch o.Type {
		case Saw:
			return SawFunc(stat, delta)
		case Sin:
			return SinFunc(stat, delta)
		case Square:
			return SquareFunc(stat, delta)
		}

		return TriangleFunc(stat, delta)
	}
}
