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
	switch o.Type {
	case Saw:
		return SawFunc
	case Sin:
		return SinFunc
	case Square:
		return SquareFunc
	}

	return TriangleFunc
}
