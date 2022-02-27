package modules

import (
	"errors"

	"github.com/faiface/beep"
)

type Voice struct {
	octaveFreq int
	noteFreq   int
	sampleRate float64
	osc        *OscStream
}

func NewVoice(sr beep.SampleRate, baseFreq int, oscFunc ModuleFunction) (*Voice, error) {
	if int(sr)/baseFreq < 2 {
		return nil, errors.New("faiface beep tone generator: samplerate must be at least 2 times grater then frequency")
	}
	osc := &OscStream{
		OscFunc: oscFunc,
		Stat:    0.0,
	}

	g := &Voice{
		sampleRate: float64(sr),
		osc:        osc,
		octaveFreq: baseFreq,
	}

	g.SetNote(0)

	return g, nil
}

func (g *Voice) GetOsc() *OscStream {
	return g.osc
}

func (g *Voice) GetOctaveFreq() int {
	return g.octaveFreq
}

func (g *Voice) OctaveFreqUp() {
	g.octaveFreq *= 2
}

func (g *Voice) OctaveFreqDown() {
	g.octaveFreq /= 2
}

func (g *Voice) GetFreq() int {
	return g.octaveFreq + g.noteFreq
}

func (g *Voice) SetNote(note int) {
	if note < 0 {
		return
	}

	semiToneOffset := float64(g.octaveFreq) / 12

	steps := g.sampleRate / (float64(g.octaveFreq) + (semiToneOffset * float64(note)))
	g.osc.Delta = 1.0 / steps

	g.noteFreq = note
}
