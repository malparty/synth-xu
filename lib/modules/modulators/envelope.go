package modulators

import (
	"github.com/malparty/synth-xu/lib/modules"
)

type EnvelopeState string

const (
	Attack  EnvelopeState = "a"
	Decay   EnvelopeState = "d"
	Sustain EnvelopeState = "s"
	Release EnvelopeState = "r"
)

// All ADSR values shall be between 0 and 1
type Envelope struct {
	Attack  float64
	Decay   float64
	Sustain float64
	Release float64

	progress float64
	state    EnvelopeState
}

func (e *Envelope) ReleaseNote() {
	e.state = Release
	e.progress = 0
}

func (e *Envelope) TriggerNote() {
	e.state = Attack
	e.progress = 0
}

func (e *Envelope) GetModuleFunc() modules.ModuleFunction {
	// Init state is end of release (no sound)
	e.progress = e.Release
	e.state = Release

	return func(stat, delta float64) float64 {
		e.progress += delta / 1000

		switch e.state {
		case Attack:
			if e.progress > e.Attack {
				// switch to decase
				e.progress = 0
				e.state = Decay

				return stat
			}

			return stat * (e.progress / e.Attack)
		case Decay:
			if e.progress > e.Decay {
				// switch to Sustain
				e.progress = 0
				e.state = Sustain

				return stat
			}

			return stat * (1 - e.progress/e.Decay*e.Sustain)
		case Sustain:
			return stat * (1 - e.Sustain)
		case Release:
			if e.progress > e.Release {
				return 0
			}

			return stat * e.Sustain * (e.Release - e.progress)
		}

		return 0
	}
}
