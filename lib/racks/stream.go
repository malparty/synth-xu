package racks

import (
	"errors"
	"math"

	"github.com/malparty/synth-xu/lib/modules"
)

// create stream which will produce infinite osciator tone with the given frequency
// use other wrappers of this package to change amplitude or add time limit
// sampleRate must be at least two times grater then frequency, otherwise this function will return an error
type OscStream struct {
	Stat       float64 // progress from 0 to 1
	Delta      float64 // space between two calculation
	StreamFunc modules.ModuleFunction
}

func (c *OscStream) nextSample() float64 {
	r := c.StreamFunc(c.Stat, c.Delta)
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
