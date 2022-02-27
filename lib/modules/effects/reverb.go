package effects

import (
	"fmt"

	"github.com/malparty/synth-xu/lib/modules"
)

type Reverb struct {
	MixRate  float64
	FadeRate float64
	DelayMs  int

	buffer  []float64
	bufferb []float64
	bufferc []float64
	bufferd []float64

	bufferLen  int
	bufferbLen int
	buffercLen int
	bufferdLen int

	currentIndex  int
	currentIndexb int
	currentIndexc int
	currentIndexd int
}

func (r *Reverb) SetDelay(delay int) {
	r.DelayMs = delay

	fmt.Printf("SET DELAYS: %d \n", delay)

	r.resetBufferSize()
}

func (r *Reverb) GetModuleFunc() modules.ModuleFunction {
	r.currentIndex = 0
	r.buffer = []float64{}

	r.resetBufferSize()

	return func(stat float64, delta float64) (reverbLevel float64) {
		if r.currentIndex >= r.bufferLen {
			r.currentIndex = 0
		}
		if r.currentIndexb >= r.bufferbLen {
			r.currentIndexb = 0
		}
		if r.currentIndexc >= r.buffercLen {
			r.currentIndexc = 0
		}
		if r.currentIndexd >= r.bufferdLen {
			r.currentIndexd = 0
		}

		reverbLevel = r.buffer[r.currentIndex] * r.FadeRate / 100
		reverbLevelb := r.bufferb[r.currentIndexb] * r.FadeRate / 100
		reverbLevelc := r.bufferc[r.currentIndexc] * r.FadeRate / 100
		reverbLeveld := r.bufferd[r.currentIndexd] * r.FadeRate / 100

		r.buffer[r.currentIndex] = reverbLeveld + stat*r.MixRate/100
		r.bufferb[r.currentIndexb] = reverbLevel + stat*r.MixRate/100
		r.bufferc[r.currentIndexc] = reverbLevelb + stat*r.MixRate/100
		r.bufferd[r.currentIndexd] = reverbLevelc + stat*r.MixRate/100
		r.currentIndex++
		r.currentIndexb++
		r.currentIndexc++
		r.currentIndexd++

		return reverbLeveld + stat*r.MixRate/100
	}
}

func (r *Reverb) resetBufferSize() {
	r.buffer = []float64{}

	bufferSize := r.DelayMs

	r.buffer = r.resizeBuffer(bufferSize)
	r.bufferLen = bufferSize

	r.bufferb = r.resizeBuffer(bufferSize + 5)
	r.bufferbLen = bufferSize + 5

	r.bufferc = r.resizeBuffer(bufferSize + 3)
	r.buffercLen = bufferSize + 3

	r.bufferd = r.resizeBuffer(bufferSize - 5)
	r.bufferdLen = bufferSize - 5
}

func (r *Reverb) resizeBuffer(size int) []float64 {
	buffer := []float64{}

	sizeInt := size
	for i := 0; i < sizeInt+1; i++ {
		buffer = append(buffer, 0)
	}

	return buffer
}
