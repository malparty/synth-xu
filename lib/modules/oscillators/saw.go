package oscillators

import (
	"math"
)

func SawFunc(stat float64, delta float64) float64 {
	_, r := math.Modf(stat + delta + 0.5)
	return r*2.0 - 1.0
}
