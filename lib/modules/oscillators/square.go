package oscillators

func SquareFunc(stat float64, _ float64) float64 {
	if stat < 0 {
		return -1
	}

	return 1
}
