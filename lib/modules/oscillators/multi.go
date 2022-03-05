package oscillators

import "github.com/malparty/synth-xu/lib/modules"

type MultiOsc struct {
	OscA       *Osc
	OscB       *Osc
	MixPercent int
}

func (o *MultiOsc) GetModuleFunc() modules.ModuleFunction {
	return func(stat, delta float64) float64 {
		statA := o.OscA.GetModuleFunc()(stat, delta)
		statB := o.OscB.GetModuleFunc()(stat, delta)

		return (statA*float64(o.MixPercent) + statB*float64(100-o.MixPercent)) / 100
	}
}
