package oscillators

func SquareFunc(stat float64, delta float64) float64 {
	if stat+delta < -0.5 || stat+delta > 0.5 {
		return -1
	}

	return 1
}
