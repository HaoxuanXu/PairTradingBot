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

func MaxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}
