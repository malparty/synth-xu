package racks

import (
	"errors"

	"github.com/faiface/beep"
)

type Voice struct {
	octaveFreq    int
	noteFreq      int
	sampleRate    float64
	oscStream     *OscStream
	ChainFunction *ModulesChain
}

func NewVoice(sr beep.SampleRate, baseFreq int, chain *ModulesChain) (*Voice, error) {
	if int(sr)/baseFreq < 2 {
		return nil, errors.New("faiface beep tone generator: samplerate must be at least 2 times grater then frequency")
	}
	osc := &OscStream{
		StreamFunc: chain.ChainFuncControlled,
		Stat:       0.0,
	}

	g := &Voice{
		ChainFunction: chain,
		sampleRate:    float64(sr),
		oscStream:     osc,
		octaveFreq:    baseFreq,
	}

	g.SetNote(0)

	return g, nil
}

func (v *Voice) GetOsc() *OscStream {
	return v.oscStream
}

func (v *Voice) GetOctaveFreq() int {
	return v.octaveFreq
}

func (v *Voice) OctaveFreqUp() {
	v.octaveFreq *= 2
}

func (v *Voice) OctaveFreqDown() {
	v.octaveFreq /= 2
}

func (v *Voice) GetFreq() int {
	return v.octaveFreq + v.noteFreq
}

func (v *Voice) SetNote(note int) {
	if note < 0 {
		return
	}

	semiToneOffset := float64(v.octaveFreq) / 12

	steps := v.sampleRate / (float64(v.octaveFreq) + (semiToneOffset * float64(note)))
	v.oscStream.Delta = 1.0 / steps

	v.noteFreq = note
}
