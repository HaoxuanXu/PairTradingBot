package updater

func findMinMax(sli []float64) (float64, float64) {
	min := sli[0]
	max := sli[0]
	for _, value := range sli {
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}
	return min, max
}

func UpdatePriceRatioThreshold(longShortRatios, shortLongRatios []float64) float64 {
	minVal, _ := findMinMax(longShortRatios[1:])
	_, maxVal := findMinMax(shortLongRatios)

	return (minVal + maxVal) / 2.0
}
