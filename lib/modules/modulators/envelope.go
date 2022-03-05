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

	freq         float64 // current freq to assure time consistency across notes
	progress     float64
	currentLevel float64
	state        EnvelopeState
}

func (e *Envelope) ReleaseNote() {
	e.state = Release
	e.progress = 0
}

func (e *Envelope) TriggerNote() {
	e.state = Attack
	e.progress = 0
}

func (e *Envelope) SetFreq(freq float64) {
	if freq > 0 {
		e.freq = freq
	}
}

func (e *Envelope) GetModuleFunc() modules.ModuleFunction {
	// Init state is end of release (no sound)
	e.progress = e.Release
	e.state = Release

	return func(stat, delta float64) float64 {
		e.progress += delta / e.freq

		switch e.state {
		case Attack:
			if e.progress > e.Attack {
				// switch to decase
				e.progress = 0
				e.state = Decay

				return stat
			}

			e.currentLevel = e.progress / e.Attack

			return stat * e.currentLevel
		case Decay:
			if e.progress > e.Decay {
				// switch to Sustain
				e.progress = 0
				e.state = Sustain

				return stat
			}

			e.currentLevel = 1 - e.progress/e.Decay*e.Sustain

			return stat * e.currentLevel
		case Sustain:
			return stat * (1 - e.Sustain)
		case Release:
			if e.progress > e.Release {
				return 0
			}

			return stat * e.currentLevel * (1 - e.progress/e.Release)
		}

		return 0
	}
}
