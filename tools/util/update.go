package util

func UpdateIntSlice(sli *[]int, val int) {
	if val != 0 {
		*sli = append(*sli, val)
	}
}

func UpdateFloatSlice(sli *[]float64, val float64) {
	if val != 0.0 {
		*sli = append(*sli, val)
	}
}
