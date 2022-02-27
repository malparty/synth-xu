package generators

import (
	"errors"
	"math"

	"github.com/faiface/beep"
)

type GeneratorFunction func(stat float64, delta float64) float64

// create stream which will produce infinite osciator tone with the given frequency
// use other wrappers of this package to change amplitude or add time limit
// sampleRate must be at least two times grater then frequency, otherwise this function will return an error
type OscStream struct {
	Stat    float64 // progress from 0 to 1
	Delta   float64 // space between two calculation
	OscFunc GeneratorFunction
}

type Generator struct {
	octaveFreq int
	noteFreq   int
	sampleRate float64
	osc        *OscStream
}

func NewGenerator(sr beep.SampleRate, baseFreq int, oscFunc GeneratorFunction) (*Generator, error) {
	if int(sr)/baseFreq < 2 {
		return nil, errors.New("faiface beep tone generator: samplerate must be at least 2 times grater then frequency")
	}
	osc := &OscStream{
		OscFunc: oscFunc,
		Stat:    0.0,
	}

	g := &Generator{
		sampleRate: float64(sr),
		osc:        osc,
		octaveFreq: baseFreq,
	}

	g.SetNote(0)

	return g, nil
}

func (g *Generator) GetOsc() *OscStream {
	return g.osc
}

func (g *Generator) GetOctaveFreq() int {
	return g.octaveFreq
}

func (g *Generator) OctaveFreqUp() {
	g.octaveFreq *= 2
}

func (g *Generator) OctaveFreqDown() {
	g.octaveFreq /= 2
}

func (g *Generator) GetFreq() int {
	return g.octaveFreq + g.noteFreq
}

func (g *Generator) SetNote(note int) {
	if note < 0 {
		return
	}

	semiToneOffset := float64(g.octaveFreq) / 12

	steps := g.sampleRate / (float64(g.octaveFreq) + (semiToneOffset * float64(note)))
	g.osc.Delta = 1.0 / steps

	g.noteFreq = note
}

func (c *OscStream) nextSample() float64 {
	r := c.OscFunc(c.Stat, c.Delta)
	_, c.Stat = math.Modf(c.Stat + c.Delta)
	return r
}

func (c *OscStream) Stream(buf [][2]float64) (int, bool) {
	for i := 0; i < len(buf); i++ {
		s := c.nextSample()
		buf[i] = [2]float64{s, s}
	}
	return len(buf), true
}

func (c *OscStream) Err() error {
	return errors.New("error with OscStream")
}
