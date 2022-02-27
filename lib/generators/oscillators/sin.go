package oscillators

import (
	"errors"
	"math"

	"github.com/malparty/synth-xu/lib/generators"

	"github.com/faiface/beep"
)

// create streamer which will produce infinite sinusoid tone with the given frequency
// use other wrappers of this package to change amplitude or add time limit
// sampleRate must be at least two times grater then frequency, otherwise this function will return an error
func SinTone(sr beep.SampleRate, freq int) (beep.Streamer, error) {
	if int(sr)/freq < 2 {
		return nil, errors.New("faiface beep tone generator: samplerate must be at least 2 times grater then frequency")
	}
	r := &generators.OscStream{
		OscFunc: sinFunc,
	}
	r.Stat = 0.0
	srf := float64(sr)
	ff := float64(freq)
	steps := srf / ff
	r.Delta = 1.0 / steps
	return r, nil
}

func sinFunc(stat float64, _ float64) float64 {
	return math.Sin(stat * 2.0 * math.Pi)
}
