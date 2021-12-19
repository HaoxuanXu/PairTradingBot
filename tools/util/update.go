package util

func UpdateIntSlice(sli *[]int, val int) {
	*sli = append(*sli, val)
}

func UpdateFloatSlice(sli *[]float64, val float64) {
	*sli = append(*sli, val)
}