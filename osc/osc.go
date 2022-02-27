package osc

import "math"

const (
	SampleRate = 48000

	baseFreq = 220
)

// Osc is an infinite Osc of 440 Hz sine wave.
type Osc struct {
	position  int64
	remaining []byte

	octaveFrequency int64
	frequency       int64
}

func NewOsc() *Osc {
	return &Osc{
		octaveFrequency: baseFreq,
		frequency:       baseFreq,
	}
}

func (o *Osc) MoveOctave(up bool) {
	if up {
		o.octaveFrequency *= 2
	} else {
		o.octaveFrequency /= 2
	}
}

func (o *Osc) MoveNote(note int64) {
	offset := o.octaveFrequency / 12
	o.frequency = o.octaveFrequency + offset*note
}

// Read is io.Reader's Read.
//
// Read fills the data with sine wave samples.
func (o *Osc) Read(buf []byte) (int, error) {
	if len(o.remaining) > 0 {
		n := copy(buf, o.remaining)
		o.remaining = o.remaining[n:]
		return n, nil
	}

	var origBuf []byte
	if len(buf)%4 > 0 {
		origBuf = buf
		buf = make([]byte, len(origBuf)+4-len(origBuf)%4)
	}

	length := int64(SampleRate / o.frequency)
	p := o.position / 4
	for i := 0; i < len(buf)/4; i++ {
		const max = 32767
		b := int16(math.Sin(2*math.Pi*float64(p)/float64(length)) * max)
		buf[4*i] = byte(b)
		buf[4*i+1] = byte(b >> 8)
		buf[4*i+2] = byte(b)
		buf[4*i+3] = byte(b >> 8)
		p++
	}

	o.position += int64(len(buf))
	o.position %= length * 4

	if origBuf != nil {
		n := copy(origBuf, buf)
		o.remaining = buf[n:]
		return n, nil
	}
	return len(buf), nil
}

// Close is io.Closer's Close.
func (s *Osc) Close() error {
	return nil
}
