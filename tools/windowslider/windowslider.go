package windowslider

func SlideWindowInt(array []int, windowSize int) []int {
	arrayLength := len(array)
	offset := arrayLength - windowSize
	if offset > 0 {
		array = array[offset:]
	}
	return array
}

func SlideWindowFloat(array []float64, windowSize int) []float64 {
	arrayLength := len(array)
	offset := arrayLength - windowSize

	if offset > 0 {
		array = array[offset:]
	}
	return array
}
