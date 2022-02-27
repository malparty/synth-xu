package oscillators

import (
	"math"
)

func SinFunc(stat float64, _ float64) float64 {
	return math.Sin(stat * 2.0 * math.Pi)
}
