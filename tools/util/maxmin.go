package util

import "math"

func GetMaxInt(array []int) int {
	maxVal := int(math.Inf(-1))
	for _, val := range array {
		if val > maxVal {
			maxVal = val
		}
	}
	return maxVal
}

func GetMaxFloat(array []float64) float64 {
	maxVal := float64(math.Inf(-1))
	for _, val := range array {
		if val > maxVal {
			maxVal = val
		}
	}
	return maxVal
}

func GetAvgInt(array []int) int {
	sum := int(0)
	for i := 0; i < len(array); i++ {
		sum += array[i]
	}

	return int(sum / len(array))
}

func GetAvgFloat(array []float64) float64 {
	sum := float64(0.0)
	for i := 0; i < len(array); i++ {
		sum += array[i]
	}

	return float64(sum / float64(len(array)))
}
