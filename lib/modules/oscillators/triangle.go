package oscillators

import "math"

// TODO: Test and validate — it's probably wrong
func TriangleFunc(stat float64, delta float64) float64 {
	if stat < 0 {
		_, r := math.Modf(2*stat - delta)
		return r*2.0 - 1.0
	}

	_, r := math.Modf(2*stat + delta)
	return r*2.0 - 1.0
}
