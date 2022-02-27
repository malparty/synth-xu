package effects

import "github.com/malparty/synth-xu/lib/modules"

type Limiter struct {
	Rate float64
}

func (l *Limiter) GetModuleFunc() modules.ModuleFunction {
	return func(stat float64, _ float64) float64 {
		return stat * l.Rate / 100
	}
}
